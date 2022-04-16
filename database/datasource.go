package database

import (
	"github.com/allentom/harukap/datasource"
	"gorm.io/gorm"
)

var DefaultPlugin = &datasource.Plugin{
	OnConnected: func(db *gorm.DB) {
		Instance = db
		Instance.AutoMigrate(&User{}, &App{}, &AccessToken{}, &RefreshToken{}, &AuthorizationCode{})
	},
}
