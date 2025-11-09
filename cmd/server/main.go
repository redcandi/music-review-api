package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"music-review-api/internal/api"
	"music-review-api/internal/db"
)


func main(){
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:",err)
	}
	defer database.Close()


	env := &api.ApiEnv{DB: database}

	//router setup
	r:=gin.Default()
	
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET","POST","PUT","DELETE"},
		AllowHeaders: []string{"Origin","Content-Type","Authorization"},
		AllowCredentials: true,
	}))

	api.SetupRoutes(r,env)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil{
		log.Fatal("Failed to start server:", err)
	}

}
