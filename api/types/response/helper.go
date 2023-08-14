package response

import "github.com/gin-gonic/gin"

type Message struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RetMsg(c *gin.Context, code int, msg string, data interface{}) {
	retMsg := Message{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	c.AbortWithStatusJSON(code, retMsg)
}
