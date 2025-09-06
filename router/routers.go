package router

import (
	"CloudMusic/controllers"
	"CloudMusic/controllers/middle"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	user := r.Group("/api/user") //登录 注册
	{
		user.POST("/login", controllers.EmaliLogin)           //登录
		user.POST("/register", controllers.RegisterWithEmail) //注册
		user.POST("/getcode", controllers.GetEmailCode)
	}
	authUser := r.Group("/api/user") // 用户
	authUser.Use(middle.AuthMiddle)  //token 中间件
	{
		authUser.POST("/logout", controllers.Logout)                   // 登出
		authUser.GET("/info", controllers.GetUserInfo)                 // 获取个人信息
		authUser.GET("/info/:userid", controllers.GetUserInfoByUserId) //获取指定用户的信息
		authUser.PUT("/info", controllers.PutUserInfo)                 // 更新个人信息
	}
	song := r.Group("/api/v1/songs") // 歌曲
	song.Use(middle.AuthMiddle)
	{
		song.GET("/:id", controllers.GetSongDetail) // 歌曲详细信息
		song.GET("url/:id", controllers.GetSongUrl) // 歌曲播放地址
		song.GET("lyric/:id", controllers.GetLyric) // 歌词
	}
	artists := r.Group("/api/artists")
	artists.Use(middle.AuthMiddle)
	{
		//获取歌手歌曲总数
		artists.GET("/:id/count", controllers.GetArtistSongsNum)
		// 获取歌手列表
		artists.GET("/", controllers.GetArtistsList)
		//获取歌手详情
		artists.GET("/:id", controllers.GetArtistsDetail)
		//获取歌手的歌曲
		artists.GET("/:id/songs", controllers.GetArtistsSongs)
		//获取歌手的专辑
		artists.GET("/:id/albums", controllers.GetArtistsAlbums)
	}
	albums := r.Group("/api/albums")
	albums.Use(middle.AuthMiddle)
	{
		albums.GET("/:id", controllers.GetAlbumDetail)
		albums.GET("/:id/songs", controllers.GetAlbumSongs)
	}

	playlists := r.Group("/api/playlists")
	playlists.Use(middle.AuthMiddle)
	{
		//	GET /api/playlists/info/:id - 获取歌单详情
		playlists.GET("/info/:id", controllers.GetPlayListInfo)
		//GET /api/playlists/ - 获取用户全部歌单
		playlists.GET("/list", controllers.GetPlayLists)
		//GET /api/playlists/list/:userid - 获取指定用户的所有歌单
		playlists.GET("/list/:userid", controllers.GetPlayListsByUserId)
		//GET /api/playlists/:id/songs - 获取歌单的歌曲
		playlists.GET("/:id/songs", controllers.GetPlayListSongs)
		//POST /api/playlists - 创建歌单
		playlists.POST("/", controllers.CreatePlayList)
		//PUT /api/playlists/:id - 更新歌单
		playlists.PUT("/:id", controllers.UpdatePlayList)
		//DELETE /api/playlists/:id - 删除歌单
		playlists.DELETE("/:id", controllers.DeletePlayList)
		//POST /api/playlists/:id/songs - 添加歌曲到歌单
		playlists.POST("/:id/songs", controllers.AddSongPlayList)
		//DELETE /api/playlists/:id/songs/:songId - 从歌单中删除歌曲
		playlists.DELETE("/:id/songs/:songId", controllers.DeleteSongPlayList)

		// GET /api/playlists/likesongids -获取所有喜欢的歌的id
		playlists.GET("/likesongids", controllers.GetLikeSongsId)
		// POST /api/playlists/toggle/like
		playlists.POST("/toggle/like", controllers.ToggleLikeSong)
		//	弃用
		// POST /api/playlists/like/:id -添加歌曲到我喜欢歌单
		//playlists.POST("/like/:id", controllers.LikeSong)
		// POST /api/playlists/unlike/:id -取消喜欢
		//playlists.POST("/unlike/:id", controllers.UnLikeSong)
	}

	everyday := r.Group("/api/everyday")
	everyday.Use(middle.AuthMiddle)
	{
		everyday.GET("/song", controllers.GetEveryDaySongList)
		everyday.GET("/mbulike", controllers.GetMayBeYouLike)
	}
	search := r.Group("/api/search")
	search.Use(middle.AuthMiddle)
	{
		search.GET("/", controllers.SearchHandler)
	}
	return r
}
