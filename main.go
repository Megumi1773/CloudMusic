package main

import (
	"CloudMusic/config"
	"CloudMusic/router"
	"CloudMusic/utils"
	"log"
)

func init() {
	config.InitConfig()
	utils.InitMinio()
	utils.InitD()
	utils.SyncDataBase()
}

func main() {
	r := router.SetupRouter()
	log.Println("sever listening on port===>" + config.AppConfig.App.Port)
	_ = r.Run(config.AppConfig.App.Port)
}
