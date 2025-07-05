package api

/*
	Controller layer (HTTP handlers)
*/

import (
	"net/http"

	"github.com/gin-gonic/gin"

	handler "web-page-analyzer/service"
)

type AnalyzeRequest struct {
	URL string `json:"url" binding:"required,url"`
}

func Analyze(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing URL"})
		return
	}

	result, err := handler.AnalyzeURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
