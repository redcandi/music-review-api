package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

type ApiEnv struct {
	DB *sql.DB
}

func SetupRoutes(r *gin.Engine, env *ApiEnv) {
	api := r.Group("/api/v1")
	{
		// Artist routes (GET all, GET one, POST)
		api.GET("/artists", env.GetArtists)
		api.GET("/artists/:id", env.GetArtistByID)
		api.POST("/artists", env.CreateArtist)

		// Album routes (GET all, GET one, POST)
		api.GET("/albums", env.GetAlbums)
		api.GET("/albums/:id", env.GetAlbumByID)
		api.POST("/albums", env.CreateAlbum)

		// Comment routes (GET for album, POST for album)
		api.GET("/albums/:id/comments", env.GetCommentsForAlbum)
		api.POST("/albums/:id/comments", env.CreateComment)
		
		// Aggregate/View route
		api.GET("/albums/top-rated", env.GetTopRatedAlbums)
	}
}
