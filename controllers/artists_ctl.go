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
	"strconv"
)

// GetArtistsList GET /api/artists -获取歌手列表
func GetArtistsList(c *gin.Context) {
	var artists []model.ArtistResp
	err := global.DB.Model(&model.Artist{}).Limit(10).Scan(&artists).Error
	if err != nil {
		log.Printf("Get artists list ==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "获取成功", artists)
}

// GetArtistsDetail GET /api/artists/:id - 获取歌手详情
func GetArtistsDetail(c *gin.Context) {
	id := Respond.GetId(c)
	var artists model.Artist
	if err := global.DB.Where("id = ?", id).First(&artists).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusInternalServerError, "歌手迷失啦")
			return
		} else {
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	result := struct {
		Id          uint   `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Nickname    string `json:"nickname,omitempty"`
		Avatar      string `json:"avatar,omitempty"`
		Description string `json:"description,omitempty"`
	}{
		Id:          artists.ID,
		Name:        artists.Name,
		Nickname:    artists.Nickname,
		Avatar:      artists.Avatar,
		Description: artists.Description,
	}
	Respond.Resp.Success(c, "获取成功", result)
}

// GetArtistSongsNum  GET /api/artists/songscount/:id
func GetArtistSongsNum(c *gin.Context) {
	id := Respond.GetId(c)
	if id == -1 {
		return
	}
	var songsNum int64
	err := global.DB.Table("songs").Where("artist_id = ?", id).Count(&songsNum).Error
	if err != nil {
		log.Printf("Get songs num ==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统异常")
		return
	}
	Respond.Resp.Success(c, "获取成功", songsNum)
}

// GetArtistsSongs GET /api/artists/:id/songs - 获取歌手的歌曲
func GetArtistsSongs(c *gin.Context) {
	id := Respond.GetId(c)
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pagesize", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	var artistsCount int64
	if err := global.DB.Model(&model.Artist{}).Where("id = ?", id).Count(&artistsCount).Error; err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		log.Printf("歌手不存在===>%v", err.Error())
		return
	}
	if artistsCount == 0 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "歌手迷失啦")
		return
	}
	offset := (page - 1) * pageSize
	var result []model.SongDetailResp
	err := global.DB.Table("songs").
		Select("songs.id, songs.name, songs.duration, artists.name as artist_name,albums.cover as album_cover,albums.name as album_name").
		Joins("left join artists on songs.artist_id = artists.id").
		Joins("left join albums on  songs.album_id = albums.id").
		Where("songs.artist_id = ?", id).
		Limit(pageSize).
		Offset(offset).
		Scan(&result).Error
	if err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		log.Printf("join联查失败===>%v", err.Error())
		return
	}
	Respond.Resp.Success(c, "获取成功！", result)
}

// GetArtistsAlbums GET /api/artists/:id/albums - 获取歌手的专辑
func GetArtistsAlbums(c *gin.Context) {
	id := Respond.GetId(c)
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pagesize", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	var artistsCount int64
	if err := global.DB.Model(&model.Album{}).Where("id = ?", id).Count(&artistsCount).Error; err != nil {
		log.Printf("歌手不存在===>%v", err.Error())
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	if artistsCount == 0 {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "歌手迷失啦")
		return
	}
	offset := (page - 1) * pageSize

	var result []model.AlbumResp
	err := global.DB.Table("albums").
		Select("albums.id, albums.name,albums.cover,albums.description,albums.release_time,artists.name as artist_name").
		Joins("left join artists on albums.artist_id = artists.id").
		Where("albums.id = ?", id).
		Limit(pageSize).
		Offset(offset).
		Scan(&result).Error
	if err != nil {
		log.Printf("专辑列表join查询失败===>%v", err.Error())
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}

	Respond.Resp.Success(c, "获取成功！", result)
}
