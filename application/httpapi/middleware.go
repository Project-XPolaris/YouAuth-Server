package httpapi

import (
	"errors"
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youauth/service"
	"net/http"
	"strings"
)

var NoAuthUrls = []string{
	"/login",
	"/login/register",
	"/register",
	"/login/success",
	"/login/oauth",
	"/oauth/token",
	"/oauth/refresh",
	"/auth/current",
	"/user/auth",
	"/info",
	"/users/register",
	"/oauth/app",
}

type AuthMiddleware struct {
}

func (m *AuthMiddleware) OnRequest(ctx *haruka.Context) {
	for _, path := range NoAuthUrls {
		if ctx.Request.URL.Path == path {
			return
		}
	}

	rawString := ctx.Request.Header.Get("Authorization")
	if len(rawString) == 0 {
		rawString = ctx.GetQueryString("token")
	}
	if len(rawString) > 0 {
		rawString = strings.Replace(rawString, "Bearer ", "", 1)
		token, err := service.ParseAuthToken(rawString)
		if err != nil {
			ctx.Interrupt()
			AbortError(ctx, err, http.StatusForbidden)
			return
		}
		ctx.Param["user"] = token.User
		return
	}
	AbortError(ctx, errors.New("auth failed"), http.StatusForbidden)
	ctx.Interrupt()
}
