package database

import "gorm.io/gorm"

type AuthorizationCode struct {
	gorm.Model
	Code   string
	AppId  uint
	UserId uint
	User   *User
	App    *App
}
type AccessToken struct {
	gorm.Model
	TokenId        string
	UserId         uint
	RefreshTokenId uint
	AppId          uint
	RefreshToken   *RefreshToken
	User           *User
	App            *App
}

type RefreshToken struct {
	gorm.Model
	Token       string
	UserId      uint
	AccessToken *AccessToken
	User        *User
}
