package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youauth/database"
	"github.com/projectxpolaris/youauth/service"
	"net/http"
)

type CreateAppData struct {
	Name     string `json:"name"`
	Callback string `json:"callback"`
}

var createAppHandler haruka.RequestHandler = func(context *haruka.Context) {
	user := context.Param["user"].(*database.User)
	var requestBody CreateAppData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	app, err := service.CreateApp(requestBody.Name, requestBody.Callback, user.ID)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	template := NewBaseAppTemplate(app)
	MakeSuccessResponseWithData(context, template)
}

var getAppListHandler haruka.RequestHandler = func(context *haruka.Context) {
	user := context.Param["user"].(*database.User)
	queryBuilder := service.AppQueryBuilder{}
	err := context.BindingInput(&queryBuilder)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	if queryBuilder.Page < 1 {
		queryBuilder.Page = 1
	}
	if queryBuilder.PageSize < 1 {
		queryBuilder.PageSize = 20
	}
	queryBuilder.UserId = user.ID
	apps, count, err := queryBuilder.GetDataAndCount()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewAppTemplateList(apps)
	MakeListResponse(context, data, count, queryBuilder.Page, queryBuilder.PageSize)
}

var removeAppHandler haruka.RequestHandler = func(context *haruka.Context) {
	user := context.Param["user"].(*database.User)
	appId := context.GetPathParameterAsString("appid")
	err := service.RemoveAppByAppId(appId, user.ID)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(context)
}

var getAppHandler haruka.RequestHandler = func(context *haruka.Context) {
	appId := context.GetQueryString("appid")
	app, err := service.GetAppByAppId(appId)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	template := NewBaseAppTemplateWithoutDetail(app)
	MakeSuccessResponseWithData(context, template)
}
