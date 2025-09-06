package utils

import (
	"CloudMusic/config"
	"CloudMusic/global"
	"CloudMusic/model"
	"github.com/meilisearch/meilisearch-go"
	"log"
)

func SyncDataBase() {
	meili := meilisearch.New(config.AppConfig.Meilisearch.Host)
	index := meili.Index("song")
	stats, _ := index.GetStats()
	if stats.NumberOfDocuments == 0 {
		// 首次全量同步
		var offset int
		const batch = 500
		for {
			var songs []model.SongDetailResp
			err := global.DB.
				Raw(`
			SELECT s.id,
			       s.name,
			       s.duration,
			       s.album_id as album_id,
			       a.name  AS artist_name,
			       al.name AS album_name,
			       al.cover AS album_cover
			FROM (
				SELECT id
				FROM songs
				ORDER BY RAND()
				LIMIT ?
				OFFSET ?
			) r
			JOIN songs  s ON s.id = r.id
			JOIN artists a ON a.id = s.artist_id
			JOIN albums al ON al.id = s.album_id
		`, batch, offset).
				Scan(&songs).Error
			if len(songs) == 0 {
				break
			}
			if err != nil {
				log.Printf("get songs===>%v", err)
				return
			}

			task, err := index.AddDocuments(&songs, nil)
			if err != nil {
				log.Fatal(err)
			}
			offset += len(songs)
			log.Printf("task 写入完成===>%v", task.TaskUID)
		}
		log.Println("全量同步完成！")
	} else {
		log.Println("搜索引擎数据已同步 已跳过...")
	}

}
func GetIndex() meilisearch.IndexManager {
	meili := meilisearch.New(config.AppConfig.Meilisearch.Host)
	index := meili.Index("song")
	return index
}
