package config

import (
	"CloudMusic/global"
	"CloudMusic/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func initDb() {
	//打开链接
	db, err := gorm.Open(mysql.Open(AppConfig.Database.Dns), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("failed to connect database, %v", err)
	}
	//设置数据库连接池 idle空闲 open普通链接 lifetime存活时间
	hsqldb, _ := db.DB()
	hsqldb.SetMaxIdleConns(AppConfig.Database.MaxInCons)
	hsqldb.SetMaxOpenConns(AppConfig.Database.MaxOpenCons)
	hsqldb.SetConnMaxLifetime(time.Hour)

	global.DB = db
	DbMigrate()
	global.Rdb = redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr,
		Password: AppConfig.Redis.Password,
	})

}

// DbMigrate 迁移数据库 对于model数据模型
func DbMigrate() {
	if err := global.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate database, %v", err)
	}
	if err := global.DB.AutoMigrate(&model.Song{}); err != nil {
		log.Fatalf("failed to migrate database, %v", err)
	}
	if err := global.DB.AutoMigrate(&model.Playlist{}); err != nil {
		log.Fatalf("failed to migrate database, %v", err)
	}
	if err := global.DB.AutoMigrate(&model.Artist{}); err != nil {
		log.Fatalf("failed to migrate database, %v", err)
	}
	if err := global.DB.AutoMigrate(&model.Album{}); err != nil {
		log.Fatalf("failed to migrate database, %v", err)
	}
	if err := global.DB.AutoMigrate(&model.Comment{}); err != nil {
		log.Fatalf("failed to migrate database, %v", err)
	}
}
