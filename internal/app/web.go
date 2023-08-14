package app

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"transfer/internal/app/router"
)

func InitGinEngine(r router.IRouter) *gin.Engine {
	app := gin.New()
	gin.SetMode(gin.ReleaseMode)
	// debug model
	if viper.GetViper().GetBool("Global.Debug") {
		gin.SetMode(gin.DebugMode)
	}

	if viper.GetViper().GetBool("Log.AccessLog") {
		app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %d \"%s\" %s\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Request.UserAgent(),
				param.Latency,
			)
		}))
	}
	app.Use(gin.Recovery())

	// 文件上传最多使用的内存，64M
	app.MaxMultipartMemory = 64 << 20
	// Router register
	r.Register(app)

	return app
}
