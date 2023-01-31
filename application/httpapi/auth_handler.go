package httpapi

import (
	"errors"
	"fmt"
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youauth/config"
	"github.com/projectxpolaris/youauth/database"
	"github.com/projectxpolaris/youauth/service"
	"net/http"
	"net/url"
)

var loginHandler haruka.RequestHandler = func(context *haruka.Context) {
	appId := context.GetQueryString("client_id")
	app, err := service.GetAppWithAppId(appId)
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	redirectUrl := context.GetQueryString("redirect_url")
	if config.Instance.ExternalLoginPage != "" {
		url, err := url.Parse(config.Instance.ExternalLoginPage)
		if err != nil {
			RaiseErrorHtml(context)
			return
		}
		query := url.Query()
		query.Add("client_id", appId)
		query.Add("redirect_url", redirectUrl)
		url.RawQuery = query.Encode()

		http.Redirect(
			context.Writer,
			context.Request,
			url.String(),
			http.StatusFound,
		)
		return
	}
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
	u, err := url.Parse(requestBody.RedirectUrl)
	if err != nil {
		RaiseErrorHtml(context)
		return
	}
	qry := u.Query()
	qry.Set("code", authCode)
	u.RawQuery = qry.Encode()
	red := u.String()
	fmt.Println(red)
	resultUrl, _ := url.Parse("/login/success")
	resultQry := resultUrl.Query()
	resultQry.Set("redirect", red)
	resultUrl.RawQuery = resultQry.Encode()
	http.Redirect(context.Writer, context.Request, resultUrl.String(), http.StatusFound)
}

type GetOauthTokenData struct {
	AppId     string `json:"appId"`
	Code      string `json:"code"`
	Secret    string `json:"secret"`
	GrantType string `json:"grantType"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

var getOauthTokenHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody GetOauthTokenData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	var accessToken string
	var refreshToken string
	switch requestBody.GrantType {
	case "password":
		accessToken, refreshToken, err = service.GenerateAppTokenByPassword(requestBody.AppId, requestBody.Username, requestBody.Password)
	default:
		accessToken, refreshToken, err = service.GenerateAppToken(requestBody.Code)
		if err != nil {
			AbortError(context, err, http.StatusBadRequest)
			return
		}
	}
	template := NewBaseAppAuthTemplate(accessToken, refreshToken)
	context.JSON(template)
}

type GenerateTokenBody struct {
	Code         string `hsource:"form" hname:"code" json:"code"`
	GrantType    string `hsource:"form" hname:"grant_type" json:"grant_type"`
	ClientId     string `hsource:"form" hname:"client_id" json:"client_id"`
	Username     string `hsource:"form" hname:"username" json:"username"`
	Password     string `hsource:"form" hname:"password" json:"password"`
	RefreshToken string `hsource:"form" hname:"refresh_token" json:"refresh_token"`
}

var generateTokenHandler haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody GenerateTokenBody
	var err error
	// is form content type
	if context.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := context.Request.ParseForm()
		if err != nil {
			AbortError(context, err, http.StatusBadRequest)
			return
		}
		err = context.BindingInput(&requestBody)
		if err != nil {
			AbortError(context, err, http.StatusBadRequest)
			return
		}
	}
	// is json content type
	if context.Request.Header.Get("Content-Type") == "application/json" {
		err := context.ParseJson(&requestBody)
		if err != nil {
			AbortError(context, err, http.StatusBadRequest)
			return
		}
	}

	var accessToken string
	var refreshToken string
	switch requestBody.GrantType {
	case "password":
		accessToken, refreshToken, err = service.GenerateAppTokenByPassword(requestBody.ClientId, requestBody.Username, requestBody.Password)
		if err != nil {
			AbortError(context, err, http.StatusBadRequest)
			return
		}
	case "authorization_code":
		accessToken, refreshToken, err = service.GenerateAppToken(requestBody.Code)
		if err != nil {
			AbortError(context, err, http.StatusBadRequest)
			return
		}
	case "refresh_token":
		accessToken, refreshToken, err = service.RefreshToken(requestBody.RefreshToken)
		if err != nil {
			AbortError(context, err, http.StatusBadRequest)
			return
		}

	}
	template := NewBaseAppAuthTemplate(accessToken, refreshToken)
	context.JSON(template)
}

type RefreshOauthTokenData struct {
	RefreshToken string `json:"refreshToken"`
}

var refreshAccessToken haruka.RequestHandler = func(context *haruka.Context) {
	var requestBody RefreshOauthTokenData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	accessTokenString, refreshTokenString, err := service.RefreshToken(requestBody.RefreshToken)
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
	token, user, err := service.GenerateSelfToken(requestBody.Username, requestBody.Password)
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

var getUserListHandler haruka.RequestHandler = func(context *haruka.Context) {
	queryBuilder := service.UserQueryBuilder{}
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
	users, count, err := queryBuilder.GetDataAndCount()
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	data := NewUserTemplateList(users)
	MakeListResponse(context, data, count, queryBuilder.Page, queryBuilder.PageSize)
}

var deleteUserHandler = func(context *haruka.Context) {
	userId := context.GetPathParameterAsString("id")
	err := service.DeleteUser(userId)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(context)
}

type ChangePasswordData struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

var changePasswordHandler = func(context *haruka.Context) {
	var requestBody ChangePasswordData
	err := context.ParseJson(&requestBody)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	rawUser := context.Param["user"]
	if rawUser == nil {
		AbortError(context, errors.New("user not found"), http.StatusBadRequest)
		return
	}
	user := rawUser.(*database.User)
	err = service.ChangePassword(user.ID, requestBody.OldPassword, requestBody.NewPassword)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponse(context)
}

var generateAuthCodeHandler = func(context *haruka.Context) {
	rawUser := context.Param["user"]
	if rawUser == nil {
		AbortError(context, errors.New("user not found"), http.StatusBadRequest)
		return
	}
	user := rawUser.(*database.User)
	appId := context.GetQueryString("appid")
	authCode, err := service.LoginWithUser(user.ID, appId)
	if err != nil {
		AbortError(context, err, http.StatusInternalServerError)
		return
	}
	MakeSuccessResponseWithData(context, haruka.JSON{
		"authCode": authCode,
	})
}
