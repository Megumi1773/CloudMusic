package utils

import (
	"CloudMusic/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreatToken(userId uint, username string) (string, error) {
	//创建token
	//塞入自定义Claims信息
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":   userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.AppConfig.Jwt.Exp)).Unix(),
	})
	//签名成token
	return token.SignedString([]byte(config.AppConfig.Jwt.Key))
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.Jwt.Key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token is invalid")
}
