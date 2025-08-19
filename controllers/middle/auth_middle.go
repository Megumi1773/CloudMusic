package middle

import (
	"CloudMusic/controllers/Respond"
	"CloudMusic/global"
	"CloudMusic/model"
	"CloudMusic/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

// AuthMiddle token全局认证中间件
func AuthMiddle(c *gin.Context) {
	//token黑名单
	var blackToken *model.TokenBlackList
	//拿到token
	token := c.GetHeader("Authorization")
	if token == "" {
		Respond.Resp.Fail(c, http.StatusUnauthorized, "未提供令牌")
		c.Abort()
		return
	}
	//从数据库读取黑名单token 后续会更换成Redis
	token = strings.Replace(token, "Bearer ", "", 1)
	if res := global.DB.Where("token = ?", token).First(&blackToken); res.RowsAffected != 0 {
		Respond.Resp.Fail(c, http.StatusUnauthorized, "令牌已失效")
		c.Abort()
		return
	}
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
	//读取token信息添加到上下文
	if claims != nil {
		c.Set("username", claims["username"].(string))
		userIDFloat := claims["userid"].(float64)
		c.Set("userid", uint(userIDFloat))
	}
	c.Next()
}
