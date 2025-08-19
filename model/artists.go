package model

import "gorm.io/gorm"

// Artist  歌手
type Artist struct {
	gorm.Model
	Name        string `gorm:"varchar(100);unique;not null" json:"name"`
	Nickname    string `gorm:"varchar(100);not null" json:"nickname"`
	Avatar      string `gorm:"varchar(100);not null" json:"avatar"`
	Description string `gorm:"varchar(100);not null" json:"description"`

	// 添加歌手与歌曲的一对多关系
	Songs []Song `gorm:"foreignKey:ArtistId" json:"-"`
	// 添加歌手与专辑的一对多关系
	Albums []Album `gorm:"foreignKey:ArtistId" json:"-"`
}

type ArtistResp struct {
	Id          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Description string `json:"description,omitempty"`
}
