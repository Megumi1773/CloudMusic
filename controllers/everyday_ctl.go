package controllers

import (
	"CloudMusic/controllers/Respond"
	"CloudMusic/global"
	"CloudMusic/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetEveryDaySongList(c *gin.Context) {
	var res []model.SongDetailResp

	err := global.DB.
		Raw(`
			SELECT s.id,
			       s.name,
			       s.duration,
			       a.name  AS artist_name,
			       al.name AS album_name,
			       al.cover AS album_cover
			FROM (
				SELECT id
				FROM songs
				ORDER BY RAND()
				LIMIT 30
			) r
			JOIN songs  s ON s.id = r.id
			JOIN artists a ON a.id = s.artist_id
			JOIN albums al ON al.id = s.album_id
		`).
		Scan(&res).Error

	if err != nil {
		log.Printf("get everyday songs===>%v", err)
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}

	Respond.Resp.Success(c, "获取成功", res)
}
