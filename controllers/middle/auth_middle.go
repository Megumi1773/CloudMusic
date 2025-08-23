package middle

import (
	"CloudMusic/controllers/Respond"
	"CloudMusic/global"
	"CloudMusic/utils"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strings"
)

var ctx = context.Background()

// AuthMiddle token全局认证中间件
func AuthMiddle(c *gin.Context) {
	//拿到token
	token := c.GetHeader("Authorization")
	if token == "" {
		Respond.Resp.Fail(c, http.StatusUnauthorized, "未提供令牌")
		c.Abort()
		return
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	//处理token
	claims, err := utils.ParseToken(token)
	if err != nil {
		var msg string
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			log.Printf("签名有误:%v", err.Error())
			msg = "签名有误"
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			msg = "令牌过期"
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			msg = "令牌格式有误"
		}
		Respond.Resp.Fail(c, http.StatusUnauthorized, msg)
		c.Abort()
		return
	}
	//redis 处理 token
	tokenHash := fmt.Sprintf("%x", sha1.Sum([]byte(token)))
	redisKey := fmt.Sprintf("token:user:%d", uint(claims["userid"].(float64)))
	saved, err := global.Rdb.Get(ctx, redisKey).Result()
	switch {
	case errors.Is(redis.Nil, err):
		Respond.Resp.Fail(c, http.StatusUnauthorized, "令牌失效")
		c.Abort()
		return
	case err != nil:
		log.Printf("AuthMiddle redis get err: %v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统异常")
		c.Abort()
		return
	case saved != tokenHash:
		Respond.Resp.Fail(c, http.StatusUnauthorized, "令牌失效")
		c.Abort()
		return
	}
	//读取token信息添加到上下文
	if claims != nil {
		c.Set("username", claims["username"].(string))
		userIDFloat := claims["userid"].(float64)
		c.Set("userid", uint(userIDFloat))
	}
	c.Next()
}
