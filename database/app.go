package database

import "gorm.io/gorm"

type App struct {
	gorm.Model
	AppId    string
	Name     string
	Callback string
	Secret   string
}
