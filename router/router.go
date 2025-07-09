package router

// # Route definitions

import (
	api "web-page-analyzer/api/controller"

	_ "web-page-analyzer/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Enable CORS
	router.Use(CORSMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// v1 Api group
	v1 := router.Group("/v1")
	{
		// analysis
		v1.POST("/analyze", api.Analyze)
		// analysis async
		v1.POST("/analyze/async", api.AnalyzeAsync)
		// analysis job status
		v1.GET("/analyze/async/:id", api.GetAsyncAnalysisByID)
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
