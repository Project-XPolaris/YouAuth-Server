package httpapi

import "github.com/projectxpolaris/youauth/database"

type BaseAppTemplate struct {
	Name   string `json:"name"`
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
}

func NewBaseAppTemplate(app *database.App) BaseAppTemplate {
	return BaseAppTemplate{
		Name:   app.Name,
		AppId:  app.AppId,
		Secret: app.Secret,
	}
}
