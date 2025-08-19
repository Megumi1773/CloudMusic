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

// GetSongDetailExample 使用GORM关联关系的歌曲详情查询示例
func GetSongDetailExample(c *gin.Context) {
	songId := Respond.GetId(c)
	if songId == -1 {
		return
	}

	var song model.Song
	// 使用预加载功能加载关联的歌手和专辑信息
	err := global.DB.Preload("Artist").Preload("Album").First(&song, songId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌曲不存在")
			return
		}
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}

	// 构建响应数据
	filterData := struct {
		Id         int    `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
		Duration   uint32 `json:"duration,omitempty"`
		Lyric      string `json:"lyric,omitempty"`
		ArtistName string `json:"artistName,omitempty"`
		AlbumName  string `json:"albumName,omitempty"`
		AlbumCover string `json:"albumCover,omitempty"`
	}{
		Id:         song.Id,
		Name:       song.Name,
		Duration:   song.Duration,
		Lyric:      song.Lyric,
		ArtistName: song.Artist.Name,
		AlbumName:  song.Album.Name,
		AlbumCover: song.Album.Cover,
	}

	Respond.Resp.Success(c, "获取成功！", filterData)
}

// GetPlaylistDetailExample 获取歌单详情示例（包含歌单中的歌曲）
func GetPlaylistDetailExample(c *gin.Context) {
	playlistId := Respond.GetId(c)
	if playlistId == -1 {
		return
	}

	var playlist model.Playlist
	// 预加载歌单所属用户和歌单中的歌曲
	err := global.DB.Preload("User").Preload("Songs").First(&playlist, playlistId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌单不存在")
			return
		}
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}

	// 构建响应数据
	songs := make([]map[string]interface{}, 0)
	for _, song := range playlist.Songs {
		songs = append(songs, map[string]interface{}{
			"id":       song.Id,
			"name":     song.Name,
			"duration": song.Duration,
		})
	}

	result := map[string]interface{}{
		"id":          playlist.ID,
		"name":        playlist.Name,
		"description": playlist.Description,
		"cover":       playlist.Cover,
		"creator":     playlist.User.Nickname,
		"songs":       songs,
	}

	Respond.Resp.Success(c, "获取成功！", result)
}

// AddSongToPlaylistExample 添加歌曲到歌单示例
func AddSongToPlaylistExample(c *gin.Context) {
	playlistId := c.Param("id")
	var req struct {
		SongId int `json:"song_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	var playlist model.Playlist
	if err := global.DB.First(&playlist, playlistId).Error; err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "歌单不存在")
		return
	}

	var song model.Song
	if err := global.DB.First(&song, req.SongId).Error; err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "歌曲不存在")
		return
	}

	// 使用Association添加关联
	if err := global.DB.Model(&playlist).Association("Songs").Append(&song); err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "添加失败")
		return
	}

	Respond.Resp.Success(c, "添加成功", nil)
}

// RemoveSongFromPlaylistExample 从歌单中删除歌曲示例
func RemoveSongFromPlaylistExample(c *gin.Context) {
	playlistId := c.Param("id")
	songId := c.Param("songId")

	var playlist model.Playlist
	if err := global.DB.First(&playlist, playlistId).Error; err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "歌单不存在")
		return
	}

	var song model.Song
	if err := global.DB.First(&song, songId).Error; err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "歌曲不存在")
		return
	}

	// 使用Association删除关联
	if err := global.DB.Model(&playlist).Association("Songs").Delete(&song); err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}

	Respond.Resp.Success(c, "删除成功", nil)
}

// GetArtistWithSongsExample 获取歌手及其歌曲示例
func GetArtistWithSongsExample(c *gin.Context) {
	artistId := Respond.GetId(c)
	if artistId == -1 {
		return
	}

	var artist model.Artist
	// 预加载歌手的歌曲和专辑
	err := global.DB.Preload("Songs").Preload("Albums").First(&artist, artistId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌手不存在")
			return
		}
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}

	// 构建响应数据
	songs := make([]map[string]interface{}, 0)
	for _, song := range artist.Songs {
		songs = append(songs, map[string]interface{}{
			"id":       song.Id,
			"name":     song.Name,
			"duration": song.Duration,
		})
	}

	albums := make([]map[string]interface{}, 0)
	for _, album := range artist.Albums {
		albums = append(albums, map[string]interface{}{
			"id":          album.ID,
			"name":        album.Name,
			"cover":       album.Cover,
			"releaseTime": album.ReleaseTime,
		})
	}

	result := map[string]interface{}{
		"id":          artist.ID,
		"name":        artist.Name,
		"avatar":      artist.Avatar,
		"description": artist.Description,
		"songs":       songs,
		"albums":      albums,
	}

	Respond.Resp.Success(c, "获取成功！", result)
}
