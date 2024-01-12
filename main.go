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

	r.GET("/api/login", routes.LoginHandler)

	err := r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
