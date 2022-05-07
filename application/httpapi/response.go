package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youauth/commons"
	"github.com/projectxpolaris/youauth/service"
	"github.com/projectxpolaris/youauth/youlog"
)

func AbortError(ctx *haruka.Context, err error, status int) {
	if apiError, ok := err.(*commons.APIError); ok {
		youlog.DefaultYouLogPlugin.Logger.Error(apiError.Err.Error())
		ctx.JSONWithStatus(haruka.JSON{
			"success": false,
			"err":     apiError.Desc,
			"code":    apiError.Code,
		}, status)
		return
	}
	// dispatch error
	switch err {
	case service.TokenExpired:
		ctx.JSONWithStatus(haruka.JSON{
			"success": false,
			"err":     "token expired",
			"code":    commons.TokenExpire,
		}, status)
		return
	}
	youlog.DefaultYouLogPlugin.Logger.Error(err.Error())
	ctx.JSONWithStatus(haruka.JSON{
		"success": false,
		"err":     err.(error).Error(),
		"code":    "9999",
	}, status)
}

func MakeSuccessResponse(context *haruka.Context) {
	context.JSON(haruka.JSON{
		"success": true,
	})
}
func MakeSuccessResponseWithData(context *haruka.Context, data interface{}) {
	context.JSON(haruka.JSON{
		"success": true,
		"data":    data,
	})
}

func MakeListResponse(context *haruka.Context, data interface{}, total int64, pageSize int, page int) {
	context.JSON(haruka.JSON{
		"success": true,
		"data": haruka.JSON{
			"total":    total,
			"pageSize": pageSize,
			"page":     page,
			"result":   data,
		},
	})
}
func RaiseErrorHtml(ctx *haruka.Context) {
	ctx.HTML("./templates/404.html", map[string]interface{}{})
}
