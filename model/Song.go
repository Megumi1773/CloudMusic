package model

import "time"

type Song struct {
	Id       uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id,omitempty"`
	Name     string    `gorm:"varchar(100);not null;unique" json:"name"`
	ArtistId uint64    `gorm:"bigint(20);not null;" json:"artist_id,omitempty"`
	AlbumId  uint64    `gorm:"bigint(20);not null;" json:"album_id,omitempty"`
	Duration uint32    `gorm:"int(11);not null;" json:"duration"`
	Url      string    `gorm:"varchar(255);not null;" json:"url"`
	Lyric    string    `gorm:"type:longtext;not null;" json:"lyric"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`

	// 添加歌曲与歌手的多对一关系
	Artist Artist `gorm:"foreignKey:ArtistId" json:"-"`
	// 添加歌曲与专辑的多对一关系
	Album Album `gorm:"foreignKey:AlbumId" json:"-"`
	// 添加歌曲与歌单的多对多关系
	Playlists []Playlist `gorm:"many2many:playlist_songs;foreignKey:Id;joinForeignKey:SongID;References:ID;joinReferences:PlaylistID" json:"-"`
}

func (Song) TableName() string {
	return "songs"
}

type SongDetail struct {
	Song
	ArtistName string `gorm:"column:artist_name" json:"artist_name" json:"artistName,omitempty"`
	AlbumName  string `gorm:"column:album_name" json:"album_name,omitempty"`
	AlbumCover string `gorm:"column:album_cover" json:"album_cover,omitempty"`
}
type SongDetailResp struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Duration   uint32 `json:"duration"`
	AlbumId    uint64 `json:"album_id"`
	ArtistName string `json:"artist_name"`
	AlbumName  string `json:"album_name"`
	AlbumCover string `json:"album_cover"`
}
