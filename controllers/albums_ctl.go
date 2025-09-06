package controllers

import (
	"CloudMusic/controllers/Respond"
	"CloudMusic/global"
	"CloudMusic/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// GetAlbumDetail GetAlbums GET /api/albums/:id - 获取专辑详情
func GetAlbumDetail(c *gin.Context) {
	id := Respond.GetId(c)
	var album model.Album
	if err := global.DB.Where("id = ?", id).First(&album).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "专辑未找到")
			return
		} else {
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	var artist model.Artist
	global.DB.Where("id = ?", album.ArtistId).First(&artist)
	res := model.AlbumResp{
		Id:          album.ID,
		Name:        album.Name,
		ArtistName:  artist.Name,
		ArtistCover: artist.Avatar,
		Cover:       album.Cover,
		Description: album.Description,
		ReleaseTime: album.ReleaseTime,
	}
	Respond.Resp.Success(c, "获取成功", res)
}

// GetAlbumSongs GET /api/albums/:id/songs - 获取专辑的歌曲
func GetAlbumSongs(c *gin.Context) {
	id := Respond.GetId(c)
	var albumCount int64
	if err := global.DB.Model(&model.Album{}).Where("id = ?", id).Count(&albumCount).Error; err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	if albumCount == 0 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "专辑未找到")
		return
	}

	var res []model.SongDetailResp
	err := global.DB.Table("songs").
		Select("songs.id, songs.name, songs.duration, artists.name as artist_name,albums.id as album_id,albums.name as album_name,albums.cover as album_cover").
		Joins("left join albums on songs.album_id = albums.id").
		Joins("left join artists on songs.artist_id = artists.id").
		Where("albums.id = ?", id).
		Scan(&res).Error
	if err != nil {
		log.Printf("专辑查歌===> %v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "获取成功", res)
}
