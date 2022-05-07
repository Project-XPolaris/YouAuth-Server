package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccessToken *AccessToken
	Apps        []*App
}
