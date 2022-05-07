package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/projectxpolaris/youauth/config"
	"github.com/projectxpolaris/youauth/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var InvalidateUsernameOrPassword = errors.New("invalid username or password")
var InvalidateTokenType = errors.New("invalid token type")
var TokenExpired = errors.New("token expired")

func CreateUser(Username string, Password string) (*database.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &database.User{Username: Username, Password: string(hashedPassword)}

	err = database.Instance.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GenerateToken(username string, password string) (string, *database.User, error) {
	user := &database.User{Username: username}
	err := database.Instance.Where("username = ?", username).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, InvalidateUsernameOrPassword
		}
		return "", nil, err
	}

	encryptionErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encryptionErr != nil {
		return "", nil, InvalidateUsernameOrPassword
	}
	_, accessTokenString, err := newJWTClaimsAndTokenString("access", username, "self")
	if err != nil {
		return "", nil, err
	}
	storeAccessToken := &database.AccessToken{
		TokenId: accessTokenString,
		UserId:  user.ID,
	}
	err = database.Instance.Create(storeAccessToken).Error
	if err != nil {
		return "", nil, err
	}
	return accessTokenString, user, nil

}

func ParseToken(tokenString string) (*AuthClaim, error) {
	claims := AuthClaim{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.Instance.JWTConfig.Secret), nil
	})
	if err != nil {
		if jwtErr, ok := err.(*jwt.ValidationError); ok {
			if jwtErr.Errors == jwt.ValidationErrorExpired {
				return nil, TokenExpired
			}
		}
		return nil, err
	}
	if !token.Valid {
		return nil, InvalidateTokenType
	}
	return &claims, nil
}
func GetCurrentUser(accessToken string) (*database.User, error) {
	authClaim, err := ParseToken(accessToken)
	if err != nil {
		return nil, err
	}
	accessTokenRecord := &database.AccessToken{}
	err = database.Instance.Preload("App").Where("token_id = ?", accessToken).First(accessTokenRecord).Error
	if err != nil {
		return nil, err
	}
	// check app is valid
	if accessTokenRecord.App == nil && authClaim.Subject != "self" {
		return nil, InvalidateAppError
	}
	user := &database.User{}
	err = database.Instance.Where("username = ?", authClaim.Id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ParseAuthToken(tokenString string) (*database.AccessToken, error) {
	accessToken := &database.AccessToken{}
	err := database.Instance.Preload("User").Where("token_id = ?", tokenString).First(accessToken).Error
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}
