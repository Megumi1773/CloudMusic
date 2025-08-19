package model

import (
	"gorm.io/gorm"
	"time"
)

// Album  专辑
type Album struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(100);unique;not null" json:"name"`
	ArtistId    uint64    `gorm:"not null;type:bigint(20)" json:"artist_id"`
	Cover       string    `gorm:"type:varchar(255);default null" json:"cover"`
	Description string    `gorm:"type:longtext;default null" json:"description"`
	ReleaseTime time.Time `gorm:"type:timestamp;not null" json:"release_time"`

	// 添加专辑与歌手的多对一关系
	Artist Artist `gorm:"foreignKey:ArtistId" json:"-"`
	// 添加专辑与歌曲的一对多关系
	Songs []Song `gorm:"foreignKey:AlbumId" json:"-"`
}

type AlbumResp struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	ArtistName  string    `json:"artist_name"`
	Cover       string    `json:"cover"`
	Description string    `json:"description"`
	ReleaseTime time.Time `json:"release_time"`
}
