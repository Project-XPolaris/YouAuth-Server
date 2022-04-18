package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/projectxpolaris/youauth/config"
	"github.com/projectxpolaris/youauth/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
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
	encodePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	user := &database.User{Username: username, Password: string(encodePassword)}
	err = database.Instance.First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, InvalidateUsernameOrPassword
		}
		return "", nil, err
	}

	claims := &jwt.StandardClaims{
		Id:        user.Username,
		ExpiresAt: time.Now().Add(time.Duration(config.Instance.JWTConfig.AccessTokenExpire) * time.Second).Unix(),
		Issuer:    config.Instance.JWTConfig.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.Instance.JWTConfig.Secret))
	if err != nil {
		return "", nil, err
	}
	return ss, user, nil

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
	if accessTokenRecord.App == nil {
		return nil, InvalidateAppError
	}
	user := &database.User{Username: authClaim.Id}
	err = database.Instance.First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
