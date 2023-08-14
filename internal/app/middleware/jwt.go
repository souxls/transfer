package middleware

import (
	"net/http"
	"time"

	"transfer/api/types/response"
	"transfer/api/types/schema"
	"transfer/internal/app/model"
	"transfer/internal/app/pkg/logger"
	"transfer/internal/app/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Auth(userService *service.User) *jwt.GinJWTMiddleware {
	ginJwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"),
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("Security.SecureKey")),
		Timeout:          viper.GetDuration("Security.TokenExpired") * time.Minute,
		Authenticator:    authenticator(userService),
		PayloadFunc:      payloadFunc(),
		LoginResponse:    loginResponse(),
		RefreshResponse:  refreshResponse(),
		IdentityHandler: func(c *gin.Context) interface{} {
			logger.Debug("start IdentiHandler ")
			claims := jwt.ExtractClaims(c)
			return &schema.UserInfo{
				Username: claims["id"].(string),
			}
		},
		IdentityKey:  "id",
		Authorizator: authorizator(userService),
		Unauthorized: func(c *gin.Context, code int, message string) {
			logger.Debug("start Unauthorized")
			response.RetMsg(c, 401, "未授权访问", "")
		},
		// SendAuthorization: true,
		SendCookie:    true,
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		// TODO: HTTPStatusMessageFunc:
	})

	return ginJwt
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	logger.Debug("start payloadFunc")
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*schema.UserInfo); ok {
			return jwt.MapClaims{
				"id": v.Username,
			}
		}
		return jwt.MapClaims{}
	}
}

func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		response.RetMsg(c, 200, "登录成功", schema.Token{
			Token:   token,
			Expired: expire.Format(time.RFC3339),
		})
	}
}

// 校验登陆用户是否合法并生成jwt token
func authenticator(userService *service.User) func(c *gin.Context) (interface{}, error) {
	logger.Debug("start authenticator")
	return func(c *gin.Context) (interface{}, error) {
		// 常规登录
		var login schema.LoginUser
		if err := c.ShouldBind(&login); err != nil {
			logger.Debugf("%s", err)
			return "", jwt.ErrMissingLoginValues
		}

		user := &model.User{
			Username: login.Username,
			Password: login.Password,
		}

		if err := userService.Login(user); err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		return &schema.UserInfo{
			Username: user.Username}, nil
	}
}

func authorizator(userService *service.User) func(data interface{}, c *gin.Context) bool {
	logger.Debug("start authorizator")
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*schema.UserInfo); ok && userService.CheckExist(&model.User{Username: v.Username}) {
			return true
		}

		return false
	}
}

func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		response.RetMsg(c, http.StatusOK, "", gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}
