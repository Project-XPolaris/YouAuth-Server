package httpapi

import "github.com/projectxpolaris/youauth/database"

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
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewBaseAppAuthTemplate(accessToken string, refreshToken string) BaseAppAuthTemplate {
	return BaseAppAuthTemplate{AccessToken: accessToken, RefreshToken: refreshToken}
}

type BaseTokenTemplate struct {
	Id       uint             `json:"id"`
	CreateAt string           `json:"createAt"`
	App      *BaseAppTemplate `json:"app,omitempty"`
}

func NewTokenTemplate(token *database.AccessToken) BaseTokenTemplate {
	template := BaseTokenTemplate{
		Id:       token.Model.ID,
		CreateAt: token.CreatedAt.Format(timeFormat),
	}
	if token.App != nil {
		appTemplate := NewBaseAppTemplateWithoutDetail(token.App)
		template.App = &appTemplate
	}
	return template
}

func NewTokenListTemplate(tokens []*database.AccessToken) []BaseTokenTemplate {
	tokenTemplates := make([]BaseTokenTemplate, 0)
	for _, token := range tokens {
		tokenTemplates = append(tokenTemplates, NewTokenTemplate(token))
	}
	return tokenTemplates
}
