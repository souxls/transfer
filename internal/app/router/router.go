package router

import (
	apiV1 "transfer/api/v1"
	"transfer/internal/app/middleware"
	"transfer/internal/app/pkg/logger"
	"transfer/internal/app/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	// "github.com/swaggo/gin-swagger"
)

var _ IRouter = (*Router)(nil)
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

type IRouter interface {
	Register(app *gin.Engine) error
}

type Router struct {
	UserAPI     *apiV1.User
	FileAPI     *apiV1.File
	UserService *service.User
}

func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

// Router 注册/api路由
func (a *Router) RegisterAPI(app *gin.Engine) {

	auth := middleware.Auth(a.UserService)
	err := middleware.Auth(a.UserService).MiddlewareInit()
	if err != nil {
		logger.Errorf("authMiddleware.MiddlewareInit() Error: %s", err)
	}

	// 设置默认默认路由
	app.NoRoute(auth.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		logger.Infof("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	g := app.Group("/api")

	v1 := g.Group("/v1")
	{
		// 注册 /api/v1
		v1.POST("register", a.UserAPI.Register)
		v1.POST("login", auth.LoginHandler)
		v1.GET("download/:id", a.FileAPI.Get)

		token := v1.Group("auth")
		{
			token.Use(auth.MiddlewareFunc())
			token.POST("refresh", auth.RefreshHandler)

		}

		user := v1.Group("user")
		{
			user.Use(auth.MiddlewareFunc())
			// 默认返回当前登录用户的信息
			user.GET("", a.UserAPI.Get)

		}

		f := v1.Group("files")
		{
			f.Use(auth.MiddlewareFunc())
			{

				f.GET("", a.FileAPI.Query)
				f.POST("", a.FileAPI.Create)

				f.GET(":id", a.FileAPI.Get)
				f.DELETE(":id", a.FileAPI.Delete)
				// 为文件添加用户授权
				f.POST(":id/users", a.FileAPI.CreateAuth)
				// 更新文件中用户授权信息
				f.PATCH(":id/users/:userid", a.FileAPI.UpdateAuth)

				// 根据短链接获取文件信息
				f.GET("short/:url", a.FileAPI.ByShortUrl)
				// 获取拥有下载权限的文件列表
				f.GET("downloads", a.FileAPI.ForDownload)
			}
		}
	}
}
