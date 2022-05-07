package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/projectxpolaris/youauth/database"
	"github.com/projectxpolaris/youauth/service"
	"net/http"
)

var getMyTokensHandler haruka.RequestHandler = func(ctx *haruka.Context) {
	user := ctx.Param["user"].(*database.User)
	queryBuilder := &service.TokenQueryBuilder{}
	err := ctx.BindingInput(queryBuilder)
	if err != nil {
		AbortError(ctx, err, http.StatusBadRequest)
		return
	}
	queryBuilder.UserId = user.ID
	if queryBuilder.PageSize == 0 {
		queryBuilder.PageSize = 20
	}
	if queryBuilder.Page == 0 {
		queryBuilder.Page = 1
	}
	queryBuilder.Preload = []string{"App"}
	tokens, count, err := queryBuilder.QueryWithCount()
	if err != nil {
		AbortError(ctx, err, http.StatusInternalServerError)
		return
	}
	data := NewTokenListTemplate(tokens)
	MakeListResponse(ctx, data, count, queryBuilder.PageSize, queryBuilder.Page)
}
