package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youauth/config"
	"net/http"
	"net/url"
)

var infoHandler haruka.RequestHandler = func(context *haruka.Context) {
	authUrl, err := url.Parse(config.Instance.JWTConfig.Url)
	if err != nil {
		AbortError(context, err, http.StatusBadRequest)
		return
	}
	context.JSON(haruka.JSON{
		"success": true,
		"name":    "YouAuth service",
		"auth": haruka.JSON{
			"type": "weblogin",
			"urls": haruka.JSON{
				"login":       "/login",
				"accessToken": "/oauth/token",
				"current":     "/auth/current",
			},
			"url": authUrl.String(),
		},
	})
}
