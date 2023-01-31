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

func CreateApp(name string, callbackUrl string, userId uint) (*database.App, error) {
	app := database.App{
		Name:     name,
		AppId:    xid.New().String(),
		Callback: callbackUrl,
		UserId:   &userId,
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

	user := &database.User{Username: username}
	err = database.Instance.Where("username = ?", username).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", InvalidateUsernameOrPassword
		}
		return nil, "", err
	}
	encryptionErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encryptionErr != nil {
		return nil, "", InvalidateUsernameOrPassword
	}
	authId, err := GenerateAuthCode(user.ID, app.ID)
	if err != nil {
		return nil, "", err
	}
	return user, authId, nil
}
func LoginWithUser(userId uint, appId string) (string, error) {
	app := database.App{
		AppId: appId,
	}
	err := database.Instance.Where("app_id = ?", appId).First(&app).Error
	if err != nil {
		return "", err
	}
	return GenerateAuthCode(userId, app.ID)
}
func GenerateAuthCode(userId uint, appId uint) (string, error) {
	authId := xid.New().String()
	authCode := database.AuthorizationCode{
		Code:   authId,
		AppId:  &appId,
		UserId: &userId,
	}
	err := database.Instance.Create(&authCode).Error
	if err != nil {
		return "", err
	}
	return authId, nil
}

// GenerateAppTokenByPassword for login with username and password with appid
func GenerateAppTokenByPassword(appId string, username string, password string) (string, string, error) {
	app, err := GetAppByAppId(appId)
	if err != nil {
		return "", "", err
	}
	user := &database.User{Username: username}
	err = database.Instance.Where("username = ?", username).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", InvalidateUsernameOrPassword
		}
		return "", "", err
	}
	encryptionErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if encryptionErr != nil {
		return "", "", InvalidateUsernameOrPassword
	}
	_, accessTokenString, err := newJWTClaimsAndTokenString("access", username, "self")
	if err != nil {
		return "", "", err
	}
	_, refreshTokenString, err := newJWTClaimsAndTokenString("refresh", user.Username, app.AppId)
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}
func GenerateAppToken(authCode string) (string, string, error) {
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
	_, accessTokenString, err := newJWTClaimsAndTokenString("access", authRecord.User.Username, authRecord.App.AppId)
	if err != nil {
		return "", "", err
	}

	_, refreshTokenString, err := newJWTClaimsAndTokenString("refresh", authRecord.User.Username, authRecord.App.AppId)
	if err != nil {
		return "", "", err
	}

	// delete auth code
	//err = database.Instance.Unscoped().Delete(authRecord).Error
	return accessTokenString, refreshTokenString, nil
}

func RefreshToken(refreshToken string) (string, string, error) {
	refreshUserAuth, err := ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	// check app is valid
	_, accessTokenString, err := newJWTClaimsAndTokenString("access", refreshUserAuth.Id, refreshUserAuth.Subject)
	if err != nil {
		return "", "", err
	}
	_, refreshTokenString, err := newJWTClaimsAndTokenString("refresh", refreshUserAuth.Id, refreshUserAuth.Subject)
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

type AppQueryBuilder struct {
	Ids        []string `hsource:"query" hname:"ids"`
	NameSearch string   `hsource:"query" hname:"search"`
	Name       string   `hsource:"query" hname:"name"`
	Page       int      `hsource:"query" hname:"page"`
	PageSize   int      `hsource:"query" hname:"pageSize"`
	Order      string   `hsource:"query" hname:"order"`
	UserId     uint
}

func (b *AppQueryBuilder) GetDataAndCount() ([]*database.App, int64, error) {
	apps := make([]*database.App, 0)
	var count int64
	query := database.Instance.Model(&database.App{})
	if len(b.Ids) > 0 {
		query = query.Where("id in (?)", b.Ids)
	}
	if b.NameSearch != "" {
		query = query.Where("name like ?", "%"+b.NameSearch+"%")
	}
	if b.Name != "" {
		query = query.Where("name = ?", b.Name)
	}
	if b.Order != "" {
		query = query.Order(b.Order)
	}
	if b.UserId > 0 {
		query = query.Where("user_id = ?", b.UserId)
	}
	err := query.Offset((b.Page - 1) * b.PageSize).
		Limit(b.PageSize).
		Find(&apps).
		Offset(-1).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}
	return apps, count, nil
}

func RemoveAppByAppId(appId string, userId uint) error {
	// find app
	app := &database.App{}
	err := database.Instance.Where("app_id = ?", appId).First(app).Error
	if err != nil {
		return err
	}
	if *app.UserId != userId {
		return InvalidateAppError
	}
	err = database.Instance.Unscoped().Delete(&database.App{}, "app_id = ?", appId).Error
	if err != nil {
		return err
	}
	return nil
}
func GetAppByAppId(appId string) (*database.App, error) {
	app := &database.App{}
	err := database.Instance.Where("app_id = ?", appId).First(app).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}
