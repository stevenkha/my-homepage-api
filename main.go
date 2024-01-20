package main

import (
	"my-homepage-api/routes"

	log "github.com/ccpaging/log4go"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	version = "v1"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	rg := r.Group("/" + version)
	{
		routes.AnimeRoutes(rg.Group("/animes"))
		routes.MangaRoutes(rg.Group("/mangas"))
	}

	err := r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
