package controllers

import (
	"CloudMusic/config"
	"CloudMusic/controllers/Respond"
	"CloudMusic/global"
	"CloudMusic/model"
	"CloudMusic/utils"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

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

// EmaliLogin 新登入 token redis管理 邮箱验证
func EmaliLogin(c *gin.Context) {
	//邮箱 密码必填 其余omitempty
	var loginUser model.User
	if err := c.ShouldBind(&loginUser); err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数有误")
		return
	}
	if loginUser.Email == "" || loginUser.Password == "" {
		Respond.Resp.Fail(c, http.StatusBadRequest, "邮箱或密码不能为空")
		return
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@(163\.com|qq\.com|foxmail\.com|126\.com|gmail\.com|hotmail\.com)$`)
	isPassing := re.MatchString(loginUser.Email)
	if !isPassing {
		Respond.Resp.Fail(c, http.StatusBadRequest, "邮箱格式不正确")
		return
	}
	var user model.User
	if res := global.DB.Where("email = ?", loginUser.Email).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "用户不存在")
			return
		} else {
			Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	if user.Password != loginUser.Password {
		Respond.Resp.Fail(c, http.StatusBadRequest, "密码错误!")
		return
	}
	token, err := utils.CreatToken(user.ID, user.Username)
	if err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	redisKey := fmt.Sprintf("token:user:%d", user.ID)
	//不存完整token 存20b tokenHash
	tokenHash := fmt.Sprintf("%x", sha1.Sum([]byte(token)))
	if err := global.Rdb.Set(ctx, redisKey, tokenHash, time.Duration(config.AppConfig.Jwt.Exp)*time.Hour).Err(); err != nil {
		log.Printf("redis token set err ===>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
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

// RegisterWithEmail 新注册  邮箱注册
func RegisterWithEmail(c *gin.Context) {
	var RegUser model.RegUser
	var err error
	if err = c.ShouldBindJSON(&RegUser); err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if RegUser.Email == "" {
		Respond.Resp.Fail(c, http.StatusBadRequest, "邮箱地址不能为空")
		return
	}
	if RegUser.Code == "" {
		Respond.Resp.Fail(c, http.StatusBadRequest, "验证码不能为空")
		return
	}
	key := fmt.Sprintf("%s:reg:emailcode", RegUser.Email)
	var code string
	if code, err = global.Rdb.Get(ctx, key).Result(); err == nil {
		if errors.Is(err, redis.Nil) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "验证码过期")
			return
		} else if err != nil {
			log.Printf("redis get code ===>%v", err)
			Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
			return
		}
	}
	if code != RegUser.Code {
		Respond.Resp.Fail(c, http.StatusBadRequest, "验证码不正确")
		return
	}
	global.Rdb.Del(ctx, key)
	var userCont int64
	if err = global.DB.Model(&model.User{}).
		Where("email = ?", RegUser.Email).
		Count(&userCont).Error; err != nil {
		log.Printf("User count err ===>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	if userCont > 0 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "邮箱已被注册")
		return
	}
	if RegUser.Password == "" {
		Respond.Resp.Fail(c, http.StatusBadRequest, "密码不能为空")
		return
	}
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int()
	username := strings.Replace(RegUser.Email, "@qq.com", "", 1)
	nickname := "用户" + strconv.Itoa(randNum)[:6]
	newUser := model.User{
		Username: username,
		Password: RegUser.Password,
		Nickname: nickname,
		Email:    RegUser.Email,
	}
	//启动 注册事务
	tx := global.DB.Begin()
	//事务新建用户
	if err := tx.Create(&newUser).Error; err != nil {
		log.Printf("add new user err ===>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	//事务回窜用户信息 直接拿新用户的id
	//自动新建我喜欢歌单
	name := fmt.Sprintf("%s喜欢的音乐", nickname)
	//写名字为 {用户名字}喜欢的音乐 type 1 特殊类型 用户id 其他default
	if err := tx.Create(&model.Playlist{
		Name:   name,
		Type:   1,
		UserId: newUser.ID,
	}).Error; err != nil {
		tx.Rollback()
		log.Printf("create like-playlist err ===>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	tx.Commit()
	Respond.Resp.Success(c, "注册成功", newUser)
}
func GetEmailCode(c *gin.Context) {
	var email model.GetRegCode
	if err := c.ShouldBindJSON(&email); err != nil {
		log.Printf("reg mail bindjson err===>%v", err)
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var userCount int64
	if err := global.DB.Model(&model.User{}).Where("email = ?", email.Email).Count(&userCount).Error; err != nil {
		log.Printf("Get mail code user count err ===>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	if userCount > 0 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "用户已存在，请登录")
		return
	}
	key := fmt.Sprintf("%s:reg:emailcode", email.Email)
	code := fmt.Sprintf("%06d", rand.Intn(1000000)) // 6位数字验证码
	if err := global.Rdb.Set(ctx, key, code, 5*time.Minute).Err(); err != nil {
		log.Printf("redis set mail code err ===>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	if err := utils.SendCodeEmail(email.Email, code); err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "验证码发送失败,请联系管理员")
		return
	}
	Respond.Resp.Success(c, "发送成功", nil)
}

// Logout user.POST("/logout")
func Logout(c *gin.Context) {
	userID, _ := c.Get("userid")
	redisKey := fmt.Sprintf("token:user:%d", userID)
	if err := global.Rdb.Del(ctx, redisKey).Err(); err != nil {
		log.Printf("redis delete token err ===>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
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
