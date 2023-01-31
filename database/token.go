package database

import "gorm.io/gorm"

type AuthorizationCode struct {
	gorm.Model
	Code   string
	AppId  *uint
	UserId *uint
	User   *User
	App    *App
}
