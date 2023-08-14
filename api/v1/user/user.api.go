package user

import (
	"net/http"
	"transfer/api/types/response"
	"transfer/api/types/schema"
	"transfer/internal/app/model/user"
	"transfer/internal/app/pkg/logger"
	"transfer/internal/app/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var UserSet = wire.NewSet(wire.Struct(new(API), "*"))

type API struct {
	UserService        *service.User
	PermesssionService *service.Permession
}

func (a *API) Register(c *gin.Context) {
	userName := c.PostForm("Username")
	password := c.PostForm("Password")

	if userName == "" || password == "" {
		response.RetMsg(c, http.StatusBadRequest, "用户名或密码不能为空", "")
		return
	}

	user := &user.User{
		Username: userName,
		Password: password,
	}

	if err := a.UserService.Register(c, user); err != nil {
		logger.Errorf("user register error. %s", err)
		response.RetMsg(c, http.StatusBadRequest, err.Error(), "")
		return
	}
	response.RetMsg(c, http.StatusCreated, "账号注册成功", "")
}

func (a *API) Login(c *gin.Context) {

}

func (a *API) Get(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userName := claims["id"].(string)

	user := &user.User{
		Username: userName,
	}

	if err := a.UserService.ByName(user); err != nil {

		logger.Warnf("get user info error. %s", err.Error())
		response.RetMsg(c, http.StatusBadRequest, "用户信息获取失败", "")
		return
	}

	userInfo := schema.UserInfo{
		Username: user.Username,
		Realname: user.Realname,
		Rolename: a.UserService.Role(c, user.Roleid),
	}

	response.RetMsg(c, http.StatusOK, "", userInfo)
}
