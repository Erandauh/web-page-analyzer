package router

// # Route definitions

import (
	api "web-page-analyzer/api/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// analysis
	router.POST("/analyze", api.Analyze)

	return router
}
