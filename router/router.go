package router

// # Route definitions

import (
	api "web-page-analyzer/api/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Enable CORS
	router.Use(CORSMiddleware())

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// analysis
	router.POST("/analyze", api.Analyze)

	// analysis async
	router.POST("/analyze/async", api.AnalyzeAsync)

	// analysis job status
	router.GET("/analyze/async/:id", api.GetAsyncAnalysisByID)

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
