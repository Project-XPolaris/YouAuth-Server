package httpapi

import "github.com/projectxpolaris/youauth/database"

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
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewBaseAppAuthTemplate(accessToken string, refreshToken string) BaseAppAuthTemplate {
	return BaseAppAuthTemplate{AccessToken: accessToken, RefreshToken: refreshToken}
}
