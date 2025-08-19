package controllers

import (
	"CloudMusic/controllers/Respond"
	"CloudMusic/global"
	"CloudMusic/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// GetSongDetail GET /api/songs/:id - 获取歌曲详情
func GetSongDetail(c *gin.Context) {
	var result model.SongDetail
	songId := Respond.GetId(c)
	if songId == -1 {
		return
	}
	err := global.DB.Table("songs").
		Select("songs.id,songs.name,songs.duration,songs.lyric, artists.name as artist_name, albums.name as album_name, albums.cover as album_cover").
		Joins("left join artists on artists.id=songs.artist_id").
		Joins("left join albums on albums.id=songs.album_id").
		Where("songs.id = ?", songId).
		First(&result).Error
	if err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	filterData := struct {
		Id         int    `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
		Duration   uint32 `json:"duration,omitempty"`
		Lyric      string `json:"lyric,omitempty"`
		ArtistName string `json:"artistName,omitempty"`
		AlbumName  string `json:"albumName,omitempty"`
		AlbumCover string `json:"albumCover,omitempty"`
	}{
		Id:         result.Song.Id,
		Name:       result.Song.Name,
		Duration:   result.Song.Duration,
		ArtistName: result.ArtistName,
		AlbumName:  result.AlbumName,
		AlbumCover: result.AlbumCover,
	}

	Respond.Resp.Success(c, "获取成功！", filterData)
}

// GetSongUrl GET /api/songs/url/:id - 获取歌曲播放地址
func GetSongUrl(c *gin.Context) {
	songId := Respond.GetId(c)
	if songId == -1 {
		return
	}
	var song model.Song
	if err := global.DB.Where("id = ?", songId).First(&song).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌曲不存在")
			return
		} else {
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	Respond.Resp.Success(c, "获取成功！", song.Url)
}

// GetLyric GET /api/songs/lyric/:id - 获取歌词
func GetLyric(c *gin.Context) {
	songId := Respond.GetId(c)
	if songId == -1 {
		return
	}
	var song model.Song
	if err := global.DB.Where("id = ?", songId).First(&song).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌曲不存在")
			return
		} else {
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	Respond.Resp.Success(c, "获取成功！", song.Lyric)
}
