package Respond

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一全局响应类
type Response struct {
	Code int
	Data interface{}
}

var Resp *Response

func (*Response) Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": msg,
		"data":    data,
	})
}
func (*Response) Fail(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    gin.H{},
	})
}
