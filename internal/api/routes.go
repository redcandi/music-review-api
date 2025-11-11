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
		api.POST("/signup", env.SignUp)
		api.POST("/login", env.Login)

		api.POST("/artists", env.CreateArtist)
		api.GET("/artists", env.GetArtists)

		api.POST("/albums", env.CreateAlbum)
		api.GET("/albums", env.GetAlbumsByRating) 
		api.GET("/albums/search", env.SearchAlbums)  
		api.GET("/albums/:id", env.GetAlbumDetails)

		api.POST("/albums/:id/comments", env.CreateComment)

		api.POST("/genres", env.CreateGenre)
		api.GET("/genres", env.GetGenres)
		api.POST("/albums/:id/genres", env.AddGenreToAlbum)

		api.GET("/users/:username/comments", env.GetCommentsByUser)
		api.DELETE("/users/:username", env.DeleteUser)
	}
}
