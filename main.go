package main

import (
	"CloudMusic/config"
	"CloudMusic/router"
	"CloudMusic/utils"
)

func init() {
	config.InitConfig()
	utils.InitMinio()
}

func main() {
	r := router.SetupRouter()
	_ = r.Run(config.AppConfig.App.Port)
}
