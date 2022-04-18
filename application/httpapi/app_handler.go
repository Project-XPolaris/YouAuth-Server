package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youauth/service"
	"net/http"
)

type CreateAppData struct {
	Name     string `json:"name"`
	Callback string `json:"callback"`
}

var createAppHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody CreateAppData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	app, err := service.CreateApp(requestBody.Name, requestBody.Callback)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	template := NewBaseAppTemplate(app)
	MakeSuccessResponseWithData(context, template)
}
