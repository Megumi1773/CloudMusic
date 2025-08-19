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

// GetPlayListInfo GET /api/playlists/info/:id - 获取歌单详情
func GetPlayListInfo(c *gin.Context) {
	id := Respond.GetId(c)
	if id == -1 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var res model.PlaylistResp
	err := global.DB.Table("playlists").
		Select("playlists.id,playlists.name,users.nickname as nickname,users.avatar as user_avatar,playlists.description,playlists.cover,playlists.is_public,playlists.created_at as created_at").
		Joins("left join users on playlists.user_id = users.id").
		Where("playlists.id = ?", id).
		First(&res).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌单不存在")
			return
		}
		log.Printf("Joins get playlists err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "获取成功", res)
}

// GetPlayListsByUserId GET /api/playlists/list/:userid - 获取用户全部歌单
func GetPlayListsByUserId(c *gin.Context) {
	var userId int
	var err error
	userIdStr := c.Param("userid")
	userId, err = strconv.Atoi(userIdStr)
	if err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "id解析失败")
		return
	}
	var user model.User
	if err = global.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "用户不存在")
			return
		}
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	var res []model.PlaylistResp
	err = global.DB.Table("playlists").
		Select("playlists.id,playlists.name,users.nickname as nickname,users.avatar as user_avatar,playlists.cover,playlists.is_public,playlists.created_at").
		Joins("left join users on playlists.user_id = users.id").
		Where("playlists.user_id = ?", user.ID).
		Scan(&res).Error
	if err != nil {
		log.Printf("Joins get playlists by user all playlists err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "获取成功", res)
}

// GetPlayLists GET /api/playlists/list - 获取用户全部歌单
func GetPlayLists(c *gin.Context) {
	// 获取 id name cover is_public 即可
	userid, _ := c.Get("userid")
	var res []model.PlaylistResp
	err := global.DB.Table("playlists").
		Select("playlists.id,playlists.name,users.nickname as nickname,users.avatar as user_avatar,playlists.cover,playlists.is_public,playlists.created_at").
		Joins("left join users on playlists.user_id = users.id").
		Where("playlists.user_id = ?", userid).
		Scan(&res).Error
	if err != nil {
		log.Printf("Joins get playlists by user all playlists err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "获取成功", res)
}

// GetPlayListSongs GET /api/playlists/:id/songs - 获取歌单的歌曲
func GetPlayListSongs(c *gin.Context) {
	id := Respond.GetId(c)
	if id == -1 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var res []model.SongDetailResp
	err := global.DB.Table("playlist_songs").
		Select("songs.id,songs.name,songs.duration,artists.name as artist_name,albums.name as album_name,albums.cover as album_cover").
		Joins("left join songs on playlist_songs.song_id = songs.id").
		Joins("left join artists on songs.artist_id = artists.id").
		Joins("left join albums on songs.album_id = albums.id").
		Where("playlist_songs.playlist_id = ?", id).
		Scan(&res).Error
	if err != nil {

		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "获取成功", res)
}

// CreatePlayList POST /api/playlists - 创建歌单
func CreatePlayList(c *gin.Context) {
	var reqJson model.PlayListRequest
	userid, err := c.Get("userid")
	//fmt.Println(userid)
	//fmt.Println(c.Get("username"))
	if err == false {
		Respond.Resp.Fail(c, http.StatusUnauthorized, "为什么找不到你的id,你想想办法")
		return
	}
	if err := c.ShouldBind(&reqJson); err != nil {
		log.Printf("ShouldBindJSON err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	log.Printf("reqJson===>%v", reqJson)
	if reqJson.Name == "" {
		Respond.Resp.Fail(c, http.StatusBadRequest, "歌单名称不能为空")
		return
	}
	playlist := model.Playlist{
		Name:        reqJson.Name,
		UserId:      userid.(uint),
		Cover:       reqJson.Cover,
		Description: reqJson.Description,
		IsPublic:    1,
	}
	if reqJson.IsPublic != 0 {
		playlist.IsPublic = reqJson.IsPublic
	}
	if err := global.DB.Create(&playlist).Error; err != nil {
		log.Printf("Create playlists err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusBadRequest, "系统错误")
		return
	}
	Respond.Resp.Success(c, "创建成功！", reqJson.Name)
}

// UpdatePlayList PUT /api/playlists/:id - 更新歌单
func UpdatePlayList(c *gin.Context) {
	playlistId := Respond.GetId(c)
	if playlistId == -1 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数有误")
		return
	}
	userid, _ := c.Get("userid")
	var reqJson model.PlayListRequest
	if err := c.ShouldBindJSON(&reqJson); err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数有误")
		return
	}
	var playlist model.Playlist
	if err := global.DB.Where("user_id = ? AND id = ?", userid, playlistId).First(&playlist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "这不是你的歌单你要干什么？？！")
			return
		} else {
			log.Printf("put playlist --- get playlistinfo err:==>%v", err)
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	if err := global.DB.Model(&model.Playlist{}).Where("id = ?", playlistId).Updates(reqJson).Error; err != nil {
		log.Printf("put playlist --- Update playlistinfo err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "修改成功", reqJson)

}

// DeletePlayList DELETE /api/playlists/:id - 删除歌单
func DeletePlayList(c *gin.Context) {
	playlistId := Respond.GetId(c)
	if playlistId == -1 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	userid, _ := c.Get("userid")
	var playlist model.Playlist
	if err := global.DB.Where("user_id = ? AND id = ?", userid, playlistId).Unscoped().Delete(&playlist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "这不是你的歌单！！！")
			return
		} else {
			log.Printf("Delete playlist --- get playlistinfo err:==>%v", err)
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	Respond.Resp.Success(c, "删除成功", playlistId)
}

// AddSongPlayList POST /api/playlists/:id/songs - 添加歌曲到歌单
func AddSongPlayList(c *gin.Context) {

	playlistID := Respond.GetId(c)
	if playlistID == -1 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var req struct {
		SongIds []int64 `json:"song_ids"`
	}
	//Association Gorm关联方法
	if err := c.ShouldBind(&req); err != nil {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var playlist model.Playlist
	if err := global.DB.First(&playlist, playlistID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌单不存在")
			return
		} else {
			log.Printf("Get playlist err:==>%v", err)
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	//查询出已存在的歌曲ID  Pluck取出相应表的其中一列字段的所有值 打入切片
	var existSongIds []int64
	if err := global.DB.Table("playlist_songs").
		Where("playlist_id = ?", playlistID).
		Pluck("song_id", &existSongIds).Error; err != nil {
		log.Printf("Get existing song IDs err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	//制作map的已存在id 并设置值为true
	existingMap := make(map[int64]bool)
	for _, id := range existSongIds {
		existingMap[id] = true
	}
	//声明新歌id切片
	var newSongIds []int64
	//如果不存在 反转成true 打入新歌id切片
	for _, id := range req.SongIds {
		if !existingMap[id] {
			newSongIds = append(newSongIds, id)
		}
	}
	//新歌切片长度为0 返回所有歌曲已存在
	if len(newSongIds) == 0 {
		Respond.Resp.Success(c, "歌单中已经存在所有歌曲", nil)
		return
	}
	//查询歌曲存不存在
	var songs []model.Song
	if err := global.DB.Where("id IN (?)", newSongIds).Find(&songs).Error; err != nil {
		log.Printf("Get songs err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	//添加关联如歌曲歌单表 Association创建关联与数据模型表关联字段名一致
	if err := global.DB.Model(&playlist).Association("Songs").Append(&songs); err != nil {
		log.Printf("Association err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "添加成功！", req.SongIds)
}

// DeleteSongPlayList DELETE /api/playlists/:id/songs/:songId - 从歌单中删除歌曲
func DeleteSongPlayList(c *gin.Context) {
	playlistID := Respond.GetId(c)
	if playlistID == -1 {
		Respond.Resp.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	var playlist model.Playlist
	if err := global.DB.First(&playlist, playlistID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌单不存在")
			return
		} else {
			log.Printf("Get playlist err:==>%v", err)
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
			return
		}
	}
	deleteSongId := c.Param("songId")
	dId, _ := strconv.Atoi(deleteSongId)
	var song model.Song
	if err := global.DB.First(&song, dId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Respond.Resp.Fail(c, http.StatusBadRequest, "歌曲不存在")
			return
		} else {
			log.Printf("歌单删除歌曲-获取歌曲错误:==>%v", err)
			Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		}
	}
	// 删除有必要做歌单存不存在该歌曲的验证??? 存疑 可以去除
	var count int64
	if err := global.DB.Table("playlist_songs").Where("playlist_id = ? AND song_id = ?", playlist.ID, song.Id).Count(&count).Error; err != nil {
		log.Printf("Delete playlist_songs song By Get playlist_songs Count err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	if count == 0 {
		Respond.Resp.Success(c, "歌曲已经不存在于歌单", nil)
		return
	}
	//或许不必要使用关联删除 可以直接操作中间表(playlist_songs) 删除记录
	if err := global.DB.Model(&playlist).Association("Songs").Delete(&song); err != nil {
		log.Printf("Delete song err:==>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}
	Respond.Resp.Success(c, "删除成功", dId)
}
