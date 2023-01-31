package httpapi

import (
	"github.com/projectxpolaris/youauth/config"
	"github.com/projectxpolaris/youauth/database"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type BaseUserTemplate struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

func NewUserTemplate(user *database.User) BaseUserTemplate {
	return BaseUserTemplate{
		Id:       user.Model.ID,
		Username: user.Username,
	}
}
func NewUserTemplateList(users []*database.User) []BaseUserTemplate {
	userTemplates := make([]BaseUserTemplate, 0)
	for _, user := range users {
		userTemplates = append(userTemplates, NewUserTemplate(user))
	}
	return userTemplates
}

type BaseUserAuthTemplate struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func NewBaseUserAuthTemplate(token string, user *database.User) BaseUserAuthTemplate {
	return BaseUserAuthTemplate{
		Id:       user.Model.ID,
		Username: user.Username,
		Token:    token,
	}
}

type BaseAppAuthTemplate struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func NewBaseAppAuthTemplate(accessToken string, refreshToken string) BaseAppAuthTemplate {
	return BaseAppAuthTemplate{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().Add(time.Duration(config.Instance.JWTConfig.AccessTokenExpire) * time.Second).Unix(),
		TokenType:    "Bearer",
	}
}

type BaseTokenTemplate struct {
	Id       uint             `json:"id"`
	CreateAt string           `json:"createAt"`
	App      *BaseAppTemplate `json:"app,omitempty"`
}
