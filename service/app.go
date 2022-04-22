package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/projectxpolaris/youauth/config"
	"github.com/projectxpolaris/youauth/database"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

var (
	InvalidateAppError = errors.New("invalidate app")
	AuthCodeExpire     = errors.New("auth code expire")
)

type AuthClaim struct {
	jwt.StandardClaims
	Type string `json:"type"`
}

func CreateApp(name string, callbackUrl string) (*database.App, error) {
	app := database.App{
		Name:     name,
		AppId:    xid.New().String(),
		Callback: callbackUrl,
	}
	claims := &jwt.StandardClaims{
		Id:        app.AppId,
		ExpiresAt: time.Now().Add(time.Duration(config.Instance.JWTConfig.AppTokenExpire) * time.Second).Unix(),
		Issuer:    config.Instance.JWTConfig.Issuer,
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.Instance.JWTConfig.Secret))
	if err != nil {
		return nil, err
	}
	app.Secret = ss

	err = database.Instance.Create(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}
func GetAppWithAppId(appId string) (*database.App, error) {
	app := database.App{}
	err := database.Instance.Where("app_id = ?", appId).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func LoginWithApp(appId string, username string, password string) (*database.User, string, error) {
	app := database.App{
		AppId: appId,
	}
	err := database.Instance.Where("app_id = ?", appId).First(&app).Error
	if err != nil {
		return nil, "", err
	}

	encodePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	user := &database.User{Username: username, Password: string(encodePassword)}
	err = database.Instance.First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", InvalidateUsernameOrPassword
		}
		return nil, "", err
	}

	authId := xid.New().String()
	authCode := database.AuthorizationCode{
		Code:   authId,
		AppId:  app.ID,
		UserId: user.ID,
	}
	err = database.Instance.Create(&authCode).Error
	if err != nil {
		return nil, "", err
	}
	return user, authId, nil
}

func GenerateAppToken(authCode string, appId string, secret string) (string, string, error) {
	authRecord := &database.AuthorizationCode{
		Code: authCode,
	}
	err := database.Instance.Where("code = ?", authCode).First(authRecord).Error
	if err != nil {
		return "", "", err
	}
	isAuthCodeExpire := authRecord.CreatedAt.Add(time.Duration(config.Instance.JWTConfig.AuthCodeExpires)*time.Second).Unix() < time.Now().Unix()
	if isAuthCodeExpire {
		return "", "", AuthCodeExpire
	}
	err = database.Instance.Preload("User").Preload("App").Where("code = ?", authCode).First(authRecord).Error
	if err != nil {
		return "", "", err
	}
	// check app is valid
	if authRecord.App.AppId != appId || authRecord.App.Secret != secret {
		return "", "", InvalidateAppError
	}
	_, accessTokenString, err := newJWTClaimsAndTokenString("access", authRecord.User.Username, authRecord.App.AppId)
	if err != nil {
		return "", "", err
	}

	_, refreshTokenString, err := newJWTClaimsAndTokenString("refresh", authRecord.User.Username, authRecord.App.AppId)
	if err != nil {
		return "", "", err
	}
	storeRefreshToken := &database.RefreshToken{
		UserId: authRecord.User.ID,
		Token:  refreshTokenString,
	}
	err = database.Instance.Create(&storeRefreshToken).Error
	if err != nil {
		return "", "", err
	}
	storeAccessToken := &database.AccessToken{
		TokenId:        accessTokenString,
		UserId:         authRecord.User.ID,
		RefreshTokenId: storeRefreshToken.ID,
		AppId:          authRecord.App.ID,
	}
	err = database.Instance.Create(storeAccessToken).Error
	if err != nil {
		return "", "", err
	}
	err = database.Instance.Save(storeRefreshToken).Error
	if err != nil {
		return "", "", err
	}
	// delete auth code
	//err = database.Instance.Unscoped().Delete(authRecord).Error
	return accessTokenString, refreshTokenString, nil
}

func RefreshToken(refreshToken string, secret string) (string, string, error) {
	refreshTokenRecord := &database.RefreshToken{}
	err := database.Instance.Preload("AccessToken").Preload("User").Preload("AccessToken.App").Where("token = ?", refreshToken).First(refreshTokenRecord).Error
	if err != nil {
		return "", "", err
	}
	// check app is valid

	if refreshTokenRecord.AccessToken.App.Secret != secret {
		return "", "", InvalidateAppError
	}
	_, accessTokenString, err := newJWTClaimsAndTokenString("access", refreshTokenRecord.User.Username, refreshTokenRecord.AccessToken.App.AppId)
	if err != nil {
		return "", "", err
	}

	_, refreshTokenString, err := newJWTClaimsAndTokenString("refresh", refreshTokenRecord.User.Username, refreshTokenRecord.AccessToken.App.AppId)
	if err != nil {
		return "", "", err
	}
	err = database.Instance.Unscoped().Delete(refreshTokenRecord.AccessToken).Error
	if err != nil {
		return "", "", err
	}
	err = database.Instance.Unscoped().Delete(refreshTokenRecord).Error
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}
func newJWTClaimsAndTokenString(claimsType string, id string, appId string) (*AuthClaim, string, error) {
	claims := newJWTClaims(claimsType, id, appId)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Instance.JWTConfig.Secret))
	if err != nil {
		return nil, "", err
	}
	return claims, tokenString, nil
}
func newJWTClaims(claimsType string, id string, appId string) *AuthClaim {
	var expire int64
	switch claimsType {
	case "access":
		expire = time.Now().Add(time.Duration(config.Instance.JWTConfig.AccessTokenExpire) * time.Second).Unix()
	case "refresh":
		expire = time.Now().Add(time.Duration(config.Instance.JWTConfig.RefreshTokenExpire) * time.Second).Unix()
	}
	accessTokenClaims := &AuthClaim{
		StandardClaims: jwt.StandardClaims{
			Id:        id,
			ExpiresAt: expire,
			Issuer:    config.Instance.JWTConfig.Issuer,
			IssuedAt:  time.Now().Unix(),
			Subject:   appId,
		},
		Type: claimsType,
	}
	return accessTokenClaims
}
