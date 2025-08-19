package controllers

import (
	"CloudMusic/controllers/Respond"
	"CloudMusic/global"
	"CloudMusic/model"
	"CloudMusic/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Login user.POST("/login")
func Login(c *gin.Context) {
	var user model.User
	var loginUser model.User
	if err := c.ShouldBind(&loginUser); err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数有误")
		return
	}
	if loginUser.Username == "" || loginUser.Password == "" {
		Respond.Resp.Fail(c, http.StatusBadRequest, "用户名或密码不能为空")
		return
	}
	if res := global.DB.Where("username = ?", loginUser.Username).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "用户不存在")
			return
		} else {
			Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
			return
		}

	}
	if loginUser.Password != user.Password {
		Respond.Resp.Fail(c, http.StatusBadRequest, "密码错误!")
		return
	}
	token, err := utils.CreatToken(user.ID, user.Username)
	if err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	data := gin.H{
		"user":  user,
		"token": token,
	}
	Respond.Resp.Success(c, "登录成功", data)
}

// Register user.POST("/register")
func Register(c *gin.Context) {
	var user model.User
	var RegUser model.User

	if err := c.ShouldBind(&RegUser); err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if res := global.DB.Where("username = ?", RegUser.Username).First(&user); res.RowsAffected != 0 {
		Respond.Resp.Success(c, "用户已存在,请登录", nil)
		return
	}
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int()
	newUser := model.User{
		Username: RegUser.Username,
		Password: RegUser.Password,
		Nickname: "用户" + strconv.Itoa(randNum)[:6],
	}
	if res := global.DB.Create(&newUser).Error; res != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}

	Respond.Resp.Success(c, "注册成功", newUser)
}

// Logout user.POST("/logout")
func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	tokenBlackList := &model.TokenBlackList{
		Token: token,
	}
	global.DB.Create(&tokenBlackList)
	Respond.Resp.Success(c, "退出成功！", nil)
}

// GetUserInfoByUserId user.GET("/info/:userid") --获取指定用户信息
func GetUserInfoByUserId(c *gin.Context) {
	var userId int
	var err error
	userIdStr := c.Param("userid")
	userId, err = strconv.Atoi(userIdStr)
	if err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "id解析失败")
		return
	}
	var user model.User
	if err := global.DB.First(&user, "id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "用户不存在")
			return
		}
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	newUSer := gin.H{
		"username": user.Username,
		"nickname": user.Nickname,
		"email":    user.Email,
		"phone":    user.Phone,
		"avatar":   user.Avatar,
	}
	Respond.Resp.Success(c, "获取成功", newUSer)
}

// GetUserInfo user.GET("/info")
func GetUserInfo(c *gin.Context) {
	var user model.User
	username, _ := c.Get("username")
	if err := global.DB.First(&user, "username = ?", username).Error; err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	newUSer := gin.H{
		"username": user.Username,
		"nickname": user.Nickname,
		"email":    user.Email,
		"phone":    user.Phone,
		"avatar":   user.Avatar,
	}
	Respond.Resp.Success(c, "获取成功", newUSer)
}

// PutUserInfo user.PUT("/info")
func PutUserInfo(c *gin.Context) {
	var userinfo model.UserInfo
	var user model.User
	username, _ := c.Get("username")
	if err := c.ShouldBindJSON(&userinfo); err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if res := global.DB.Where("username = ?", username).First(&user); res.Error != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	//if userinfo.Nickname == user.Nickname && userinfo.Nickname != "" {
	//	Respond.Resp.Fail(c, http.StatusBadRequest, "新昵称不能与旧昵称相同")
	//	return
	//}
	//if userinfo.Phone == user.Phone && userinfo.Email != "" {
	//	Respond.Resp.Fail(c, http.StatusBadRequest, "新手机号不能与旧手机号相同")
	//	return
	//}
	//if userinfo.Email == user.Email && userinfo.Phone != "" {
	//	Respond.Resp.Fail(c, http.StatusBadRequest, "新邮箱不能与旧邮箱相同")
	//	return
	//}
	err := global.DB.Model(&user).Updates(userinfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "信息已存在")
			return
		} else {
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	Respond.Resp.Success(c, "更新成功", userinfo)
}
