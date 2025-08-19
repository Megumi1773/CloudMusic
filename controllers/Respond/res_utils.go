package Respond

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetId(c *gin.Context) int {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Resp.Fail(c, http.StatusBadRequest, "请输入正确参数")
		c.Abort()
		return -1
	}
	return id
}
