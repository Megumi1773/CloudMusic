package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Port string `yaml:"port"`
	} `yaml:"app"`
	Database struct {
		Dns         string `yaml:"dns"`
		MaxInCons   int    `yaml:"max_in_cons"`
		MaxOpenCons int    `yaml:"max_open_cons"`
	} `yaml:"database"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	Jwt struct {
		Key string `yaml:"key"`
		Exp int    `yaml:"exp"`
	} `yaml:"jwt"`
	MinIo struct {
		Url       string `yaml:"url"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Expire    int    `yaml:"expire"`
	} `yaml:"minio"`
	Email struct {
		MailHost string `yaml:"mail_host"`
		MailPort int    `yaml:"mail_port"`
		MailUser string `yaml:"mail_user"`
		MailFrom string `yaml:"mail_from"`
		MailPass string `yaml:"mail_pass"`
	} `yaml:"email"`
}

var AppConfig *Config

func InitConfig() {
	//初始化读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	//读取
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	AppConfig = &Config{}
	//赋值
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("unable to decode config, %v", err)
	}

	initDb()
}
