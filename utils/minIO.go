package utils

import (
	"CloudMusic/config"
	"github.com/minio/minio-go"
	"log"
	"net/url"
	"time"
)

var ioClint *minio.Client

func InitMinio() {
	var err error
	dns := config.AppConfig.MinIo.Url
	accKey := config.AppConfig.MinIo.AccessKey
	secKey := config.AppConfig.MinIo.SecretKey
	ioClint, err = minio.New(dns, accKey, secKey, false)
	if err != nil {
		log.Printf("link to minio failed, err:%v ===>", err)
	}
}

func GetPresUrl(bucketName string, objectName string) *url.URL {
	presignedURL, err := ioClint.PresignedGetObject(bucketName, objectName, time.Duration(config.AppConfig.MinIo.Expire), nil)
	if err != nil {
		log.Printf("GetPresUrl err:%v ===>", err)
	}
	return presignedURL
}
