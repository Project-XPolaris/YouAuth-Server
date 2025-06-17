package config

import (
	"os"
	"strconv"

	"github.com/allentom/harukap/config"
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
	configer.SetDefault("addr", getEnvOrDefault("YOUAUTH_ADDR", ":8000"))
	configer.SetDefault("application", getEnvOrDefault("YOUAUTH_APPLICATION", "You Auth Service"))
	configer.SetDefault("instance", getEnvOrDefault("YOUAUTH_INSTANCE", "main"))

	// 从环境变量读取配置，如果环境变量存在则优先使用环境变量的值
	Instance = Config{
		JWTConfig: JWTConfig{
			Secret:             getEnvOrDefault("YOUAUTH_TOKEN_SECRET", configer.GetString("token.secret")),
			Issuer:             getEnvOrDefault("YOUAUTH_TOKEN_ISSUER", configer.GetString("token.issuer")),
			AccessTokenExpire:  getEnvInt64OrDefault("YOUAUTH_TOKEN_ACCESS_EXPIRES", configer.GetInt64("token.accessTokenExpiresIn")),
			RefreshTokenExpire: getEnvInt64OrDefault("YOUAUTH_TOKEN_REFRESH_EXPIRES", configer.GetInt64("token.refreshTokenExpiresIn")),
			AuthCodeExpires:    getEnvInt64OrDefault("YOUAUTH_TOKEN_AUTH_CODE_EXPIRES", configer.GetInt64("token.authCodeExpiresIn")),
			AppTokenExpire:     getEnvInt64OrDefault("YOUAUTH_TOKEN_APP_EXPIRES", configer.GetInt64("token.appTokenExpiresIn")),
			Url:                getEnvOrDefault("YOUAUTH_TOKEN_URL", configer.GetString("token.url")),
		},
		ExternalLoginPage: getEnvOrDefault("YOUAUTH_EXTERNAL_LOGIN_PAGE", configer.GetString("externalLoginPage")),
	}
}

// getEnvOrDefault 从环境变量获取字符串值，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt64OrDefault 从环境变量获取int64值，如果不存在或转换失败则返回默认值
func getEnvInt64OrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
