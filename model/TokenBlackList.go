package model

import "gorm.io/gorm"

type TokenBlackList struct {
	gorm.Model
	Token string `gorm:"varchar(100);unique;not null" json:"token"`
}
