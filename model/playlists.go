package model

import (
	"gorm.io/gorm"
	"time"
)

// Playlist 歌单
type Playlist struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name,omitempty"`
	UserId      uint   `gorm:"type:bigint(20);not null" json:"user_id"`
	Cover       string `gorm:"type:varchar(255);default null" json:"cover,omitempty"`
	Description string `gorm:"type:text;default null" json:"description,omitempty"`
	IsPublic    uint   `gorm:"type:tinyint(1);default 1" json:"is_public,omitempty"`
	Type        uint   `gorm:"type:tinyint(1);default 1" json:"type,omitempty"`
	//type userid 联合索引 唯一约束
	_ struct{} `gorm:"uniqueIndex:uidx_user_like,priority:1;constraint:CHECK (type IN (0,1))"`
	// 添加歌单与用户的多对一关系
	User User `gorm:"foreignKey:UserId" json:"-"`
	// 添加歌单与歌曲的多对多关系
	Songs []Song `gorm:"many2many:playlist_songs;foreignKey:ID;joinForeignKey:PlaylistID;References:Id;joinReferences:SongID" json:"-"`
}

type PlaylistResp struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	UserId      uint      `json:"user_id"`
	Nickname    string    `json:"nickname"`
	UserAvatar  string    `json:"user_avatar"`
	Description string    `json:"description,omitempty"`
	Cover       string    `json:"cover"`
	IsPublic    uint      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}

type PlayListRequest struct {
	Name        string `json:"name,omitempty"`
	UserId      uint   `json:"user_id,omitempty"`
	Cover       string `json:"cover,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublic    uint   `json:"is_public,omitempty"`
}
