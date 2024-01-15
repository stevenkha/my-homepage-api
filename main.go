package main

import (
	"log"
	"my-homepage-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	version := r.Group("/v1")
	{
		version.GET("/animes", func(c *gin.Context) {
			err := routes.AnimeHandler(c)
			if err != nil {
				log.Fatal(err)
			}
		})

		version.GET("/mangas", func(c *gin.Context) {
			err := routes.MangaHandler(c)
			if err != nil {
				log.Fatal(err)
			}
		})
	}

	err := r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
