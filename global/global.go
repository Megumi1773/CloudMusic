package global

import (
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	// DB 全局数据库对象
	DB *gorm.DB
	// Rdb 全局redis对象
	Rdb *redis.Client
	// D 全局邮箱发送器
	D *gomail.Dialer
)
