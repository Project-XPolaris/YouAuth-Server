package httpapi

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/allentom/haruka"
	"github.com/allentom/haruka/middleware"
	"github.com/projectxpolaris/youauth/util"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var Logger = log.New().WithFields(log.Fields{
	"scope": "Application",
})

// 自定义静态文件处理器
type staticFileHandler struct {
	root http.FileSystem
}

func (h *staticFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !filepath.IsAbs(upath) {
		upath = filepath.Clean(upath)
	}

	// 设置正确的 MIME 类型
	ext := strings.ToLower(filepath.Ext(upath))
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	case ".woff":
		w.Header().Set("Content-Type", "font/woff")
	case ".woff2":
		w.Header().Set("Content-Type", "font/woff2")
	case ".ttf":
		w.Header().Set("Content-Type", "font/ttf")
	case ".eot":
		w.Header().Set("Content-Type", "application/vnd.ms-fontobject")
	}

	// 禁用缓存
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	http.FileServer(h.root).ServeHTTP(w, r)
}

func GetEngine() *haruka.Engine {
	e := haruka.NewEngine()
	e.UseCors(cors.AllowAll())
	e.UseMiddleware(middleware.NewLoggerMiddleware())
	e.UseMiddleware(middleware.NewPaginationMiddleware("page", "pageSize", 1, 20))
	e.UseMiddleware(&AuthMiddleware{})
	e.Router.GET("/login", loginHandler)
	e.Router.POST("/login/register", registerResultHandler)
	e.Router.GET("/register", registerHandler)
	e.Router.GET("/login/success", loginSuccessHandler)
	e.Router.POST("/login/oauth", oauthLoginHandler)
	e.Router.POST("/oauth/token", getOauthTokenHandler)
	e.Router.POST("/token", generateTokenHandler)
	e.Router.POST("/oauth/refresh", refreshAccessToken)
	e.Router.GET("/oauth/app", getAppHandler)
	e.Router.POST("/oauth/authcode", generateAuthCodeHandler)
	e.Router.GET("/auth/current", getCurrentUserHandler)
	e.Router.POST("/users/register", createUserHandler)
	e.Router.GET("/users", getUserListHandler)
	e.Router.DELETE("/user/appid:[0-9]+", deleteUserHandler)
	e.Router.POST("/user/auth", generateAuthHandler)
	e.Router.POST("/apps", createAppHandler)
	e.Router.GET("/apps", getAppListHandler)
	e.Router.POST("/my/password", changePasswordHandler)
	e.Router.DELETE("/app/{appid:[0-9|a-z|A-Z]+}", removeAppHandler)
	e.Router.GET("/info", infoHandler)
	if util.CheckFileExist("./dist") && util.FolderIsNotEmpty("./dist") && util.CheckFileExist("./dist/index.html") {
		e.Router.HandlerRouter.PathPrefix("/api").HandlerFunc(adminAPIReverse)
		e.Router.HandlerRouter.PathPrefix("/").Handler(spaHandler{
			staticPath: "./dist",
			indexPath:  "./dist/index.html",
		})
	}

	// 使用自定义的静态文件处理器
	e.Router.HandlerRouter.PathPrefix("/static/").Handler(&staticFileHandler{
		root: http.Dir("./static"),
	})

	return e
}
