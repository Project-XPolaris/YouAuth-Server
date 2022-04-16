package httpapi

import (
	"github.com/allentom/haruka"
	"github.com/allentom/haruka/middleware"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var Logger = log.New().WithFields(log.Fields{
	"scope": "Application",
})

func GetEngine() *haruka.Engine {
	e := haruka.NewEngine()
	e.UseCors(cors.AllowAll())
	e.UseMiddleware(middleware.NewLoggerMiddleware())
	e.UseMiddleware(middleware.NewPaginationMiddleware("page", "pageSize", 1, 20))
	e.Router.GET("/login", loginHandler)
	e.Router.POST("/login/register", registerResultHandler)
	e.Router.GET("/register", registerHandler)
	e.Router.GET("/login/success", loginSuccessHandler)
	e.Router.POST("/login/oauth", oauthLoginHandler)
	e.Router.POST("/oauth/token", getOauthTokenHandler)
	e.Router.POST("/oauth/refresh", refreshAccessToken)
	e.Router.GET("/auth/current", getCurrentUserHandler)
	e.Router.POST("/users", createUserHandler)
	e.Router.POST("/user/auth", generateAuthHandler)
	e.Router.POST("/apps", createAppHandler)
	e.Router.Static("/static", "./static")
	return e
}
