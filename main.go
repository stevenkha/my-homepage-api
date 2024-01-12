package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	err := r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
