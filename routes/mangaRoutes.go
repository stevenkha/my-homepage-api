package routes

import (
	"my-homepage-api/handlers"

	"github.com/gin-gonic/gin"
)

func MangaRoutes(g *gin.RouterGroup) {
	g.GET("", handlers.MangaHandler)
}
