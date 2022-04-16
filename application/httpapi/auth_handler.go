package httpapi

import (
	"fmt"
	"github.com/allentom/haruka"
	"net/http"
	"youauth/service"
)

var loginHandler haruka.RequestHandler = func(context *haruka.Context) {
	appId := context.GetQueryString("appid")
	app, err := service.GetAppWithAppId(appId)
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	redirectUrl := context.GetQueryString("redirect")
	context.HTML("./templates/login.html", map[string]interface{}{
		"AppName":  app.Name,
		"Redirect": redirectUrl,
		"AppId":    appId,
	})
}
var registerHandler haruka.RequestHandler = func(context *haruka.Context) {
	context.HTML("./templates/register.html", map[string]interface{}{})
}

type RegisterUserForm struct {
	Username string `hsource:"form" hname:"username"`
	Password string `hsource:"form" hname:"password"`
}

var registerResultHandler haruka.RequestHandler = func(context *haruka.Context) {
	err := context.Request.ParseForm()
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	var requestBody RegisterUserForm
	err = context.BindingInput(&requestBody)
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	_, err = service.CreateUser(requestBody.Username, requestBody.Password)
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	http.Redirect(context.Writer, context.Request, "/login", http.StatusFound)
}
var loginSuccessHandler haruka.RequestHandler = func(context *haruka.Context) {
	redirectUrl := context.GetQueryString("redirect")
	context.HTML("./templates/success.html", map[string]interface{}{
		"Redirect": redirectUrl,
	})
}

type OauthLoginHandler struct {
	Username    string `hsource:"form" hname:"username"`
	Password    string `hsource:"form" hname:"password"`
	AppId       string `hsource:"form" hname:"appid"`
	RedirectUrl string `hsource:"form" hname:"redirect"`
}

var oauthLoginHandler haruka.RequestHandler = func(context *haruka.Context) {
	err := context.Request.ParseForm()
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	var requestBody OauthLoginHandler
	err = context.BindingInput(&requestBody)
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	_, authCode, err := service.LoginWithApp(requestBody.AppId, requestBody.Username, requestBody.Password)
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	http.Redirect(context.Writer, context.Request, fmt.Sprintf("/login/success?redirect=%s?code=%s", requestBody.RedirectUrl, authCode), http.StatusFound)
}

type GetOauthTokenData struct {
	AppId  string `json:"appId"`
	Code   string `json:"code"`
	Secret string `json:"secret"`
}

var getOauthTokenHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody GetOauthTokenData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}

	accessTokenString, refreshTokenString, err := service.GenerateAppToken(requestBody.Code, requestBody.AppId, requestBody.Secret)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	template := NewBaseAppAuthTemplate(accessTokenString, refreshTokenString)
	MakeSuccessResponseWithData(context, template)
}

type RefreshOauthTokenData struct {
	RefreshToken string `json:"refreshToken"`
	Secret       string `json:"secret"`
}

var refreshAccessToken haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody RefreshOauthTokenData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	accessTokenString, refreshTokenString, err := service.RefreshToken(requestBody.RefreshToken, requestBody.Secret)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	template := NewBaseAppAuthTemplate(accessTokenString, refreshTokenString)
	MakeSuccessResponseWithData(context, template)
}

type RegisterUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var createUserHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody RegisterUserData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	user, err := service.CreateUser(requestBody.Username, requestBody.Password)
	template := NewUserTemplate(user)
	MakeSuccessResponseWithData(context, template)
}

type GenerateAuthData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var generateAuthHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody GenerateAuthData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	token, user, err := service.GenerateToken(requestBody.Username, requestBody.Password)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	template := NewBaseUserAuthTemplate(token, user)
	MakeSuccessResponseWithData(context, template)
}

var getCurrentUserHandler haruka.RequestHandler = func(context *haruka.Context) {
	accessToken := context.GetQueryString("token")
	user, err := service.GetCurrentUser(accessToken)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	template := NewUserTemplate(user)
	MakeSuccessResponseWithData(context, template)
}
