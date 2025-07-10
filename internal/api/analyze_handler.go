package api

/*
	Controller layer (HTTP handlers)
*/

import (
	"net/http"

	"github.com/gin-gonic/gin"

	service "web-page-analyzer/internal/service"
	"web-page-analyzer/model"
)

type AnalyzerController struct {
	Service *service.AnalyzerService
}

func NewAnalyzerController(svc *service.AnalyzerService) *AnalyzerController {
	return &AnalyzerController{Service: svc}
}

// Analyze godoc
// @Summary      Analyze a webpage
// @Description  Analyze a webpage synchronously
// @Tags         analysis
// @Accept       json
// @Produce      json
// @Param        url  body  model.AnalyzeRequest  true  "URL AnalyzeRequest"
// @Success      200  {object}  model.AnalysisResult
// @Failure      400  {object}  map[string]string
// @Router       /v1/analyze [post]
func (ac *AnalyzerController) Analyze(c *gin.Context) {
	var req model.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing URL"})
		return
	}

	result, err := ac.Service.AnalyzeURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AnalyzeAsync godoc
// @Summary      Analyze a webpage asynchronously
// @Description  Submits a URL for analysis and returns a job ID to poll the result later.
// @Tags         analysis
// @Accept       json
// @Produce      json
// @Param        url  body  model.AnalyzeRequest  true  "URL to be analyzed"
// @Success      202  {object}  model.Job
// @Failure      400  {object}  map[string]string
// @Router       /v1/analyze/async [post]
func (ac *AnalyzerController) AnalyzeAsync(c *gin.Context) {
	var req model.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing URL"})
		return
	}

	done := make(chan model.AnalysisResult)
	errChan := make(chan error)

	job := ac.Service.AnalyzeURLAsync(req.URL, done, errChan)

	c.JSON(http.StatusAccepted, job)
}

// GetAsyncAnalysisByID godoc
// @Summary      Get analysis result by job ID
// @Description  Fetch the result of an async analysis job using the provided job ID.
// @Tags         analysis
// @Produce      json
// @Param        id   path      string  true  "Job ID"
// @Success      200  {object}  model.Job
// @Failure      404  {object}  map[string]string
// @Router       /v1/analyze/async/{id} [get]
func (ac *AnalyzerController) GetAsyncAnalysisByID(c *gin.Context) {
	jobID := c.Param("id")

	job, err := ac.Service.GetAnalysisResultByID(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}
