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

	r.GET("/v1/animes", func(c *gin.Context) {
		err := routes.LoginHandler(c)
		if err != nil {
			log.Fatal(err)
		}
	})

	err := r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
