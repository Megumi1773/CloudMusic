package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"varchar(100);not null;unique" json:"username,omitempty"`
	Password string `gorm:"varchar(100);not null;" json:"password,required"`
	Nickname string `gorm:"varchar(100);default:null;unique" json:"nickname,omitempty"`
	Email    string `gorm:"varchar(100);default:null;" json:"email,required"`
	Phone    string `gorm:"varchar(100);default:null;" json:"phone,omitempty"`
	Avatar   string `gorm:"varchar(100);default:https://tse2-mm.cn.bing.net/th/id/OIP-C.p6hdmBEvZCMwVcWDVnQr0QAAAA?o=7rm=3&rs=1&pid=ImgDetMain&o=7&rm=3;" json:"avatar,omitempty"`
	// 添加用户与歌单的一对多关系
	Playlists []Playlist `gorm:"foreignKey:UserId" json:"-"`
	// 添加用户与评论的一对多关系
	Comments []Comment `gorm:"foreignKey:UserId" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type UserInfo struct {
	Nickname string ` json:"nickName,omitempty"`
	Email    string ` json:"email,omitempty"`
	Phone    string ` json:"phone,omitempty"`
}

type RegUser struct {
	Email    string ` json:"email"`
	Password string ` json:"password"`
	Code     string ` json:"code"`
}
type GetRegCode struct {
	Email string `json:"email" binding:"required"`
}
