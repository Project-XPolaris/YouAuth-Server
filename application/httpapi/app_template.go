package httpapi

import "github.com/projectxpolaris/youauth/database"

type BaseAppTemplate struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	AppId  string `json:"appId,omitempty"`
	Secret string `json:"secret,omitempty"`
}

func NewBaseAppTemplate(app *database.App) BaseAppTemplate {
	return BaseAppTemplate{
		Id:     app.ID,
		Name:   app.Name,
		AppId:  app.AppId,
		Secret: app.Secret,
	}
}
func NewBaseAppTemplateWithoutDetail(app *database.App) BaseAppTemplate {
	return BaseAppTemplate{
		Id:   app.ID,
		Name: app.Name,
	}
}
func NewAppTemplateList(apps []*database.App) []BaseAppTemplate {
	appTemplates := make([]BaseAppTemplate, 0)
	for _, app := range apps {
		appTemplates = append(appTemplates, NewBaseAppTemplate(app))
	}
	return appTemplates
}
