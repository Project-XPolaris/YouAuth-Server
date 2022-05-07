package config

import (
	"github.com/allentom/harukap/config"
	"os"
)

var DefaultConfigProvider *config.Provider

func InitConfigProvider() error {
	var err error
	customConfigPath := os.Getenv("YOUAUTH_CONFIG_PATH")
	DefaultConfigProvider, err = config.NewProvider(func(provider *config.Provider) {
		ReadConfig(provider)
	}, customConfigPath)
	return err
}

var Instance Config

type JWTConfig struct {
	Secret             string `json:"secret"`
	Issuer             string `json:"issuer"`
	AccessTokenExpire  int64
	RefreshTokenExpire int64
	AuthCodeExpires    int64
	AppTokenExpire     int64
	Url                string
}
type Config struct {
	JWTConfig         JWTConfig
	ExternalLoginPage string
}

func ReadConfig(provider *config.Provider) {
	configer := provider.Manager
	configer.SetDefault("addr", ":8000")
	configer.SetDefault("application", "You Auth Service")
	configer.SetDefault("instance", "main")

	Instance = Config{
		JWTConfig: JWTConfig{
			Secret:             configer.GetString("token.secret"),
			Issuer:             configer.GetString("token.issuer"),
			AccessTokenExpire:  configer.GetInt64("token.accessTokenExpiresIn"),
			RefreshTokenExpire: configer.GetInt64("token.refreshTokenExpiresIn"),
			AuthCodeExpires:    configer.GetInt64("token.authCodeExpiresIn"),
			AppTokenExpire:     configer.GetInt64("token.appTokenExpiresIn"),
			Url:                configer.GetString("token.url"),
		},
		ExternalLoginPage: configer.GetString("externalLoginPage"),
	}
}
