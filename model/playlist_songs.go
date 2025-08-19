package model

// PlaylistSong 歌单歌曲关联表
type PlaylistSong struct {
	ID         uint64   `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement" json:"id"`
	PlaylistID uint     `gorm:"column:playlist_id;type:bigint(20);not null;uniqueIndex:idx_playlist_song" json:"playlist_id"`
	SongID     uint     `gorm:"column:song_id;type:bigint(20);not null;uniqueIndex:idx_playlist_song" json:"song_id"`
	Playlist   Playlist `gorm:"foreignKey:PlaylistID;references:ID" json:"-"`
	Song       Song     `gorm:"foreignKey:SongID;references:ID" json:"-"`
}
