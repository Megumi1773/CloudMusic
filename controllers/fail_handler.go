package controllers

import (
	"CloudMusic/controllers/Respond"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func fail(c *gin.Context, msg string, err error) {
	if err != nil {
		log.Printf("%s===>%v", msg, err)
	}
	Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
}
