package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"web-page-analyzer/internal/persistance"
	mock_persistance "web-page-analyzer/internal/persistance/mocks"
	"web-page-analyzer/internal/process"
	"web-page-analyzer/model"
)

func TestGetAnalysisResultByID(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobID := "12345"
	expectedJob := &model.Job{
		ID:     jobID,
		Status: "COMPLETED",
	}
	executor := process.NewPatternExecutor()

	mockStore := mock_persistance.NewMockStore(ctrl)
	service := NewAnalyzerService(mockStore, executor)

	t.Run("should return job when found", func(t *testing.T) {
		mockStore.EXPECT().
			GetJob(jobID).
			Return(expectedJob, true)

		foundJob, err := service.GetAnalysisResultByID(jobID)
		assert.NoError(t, err)
		assert.Equal(t, jobID, foundJob.ID)
	})

	t.Run("should return error when job not found", func(t *testing.T) {
		mockStore.EXPECT().
			GetJob("non-existent-id").
			Return(nil, false)

		_, err := service.GetAnalysisResultByID("non-existent-id")
		assert.Error(t, err)
		assert.Equal(t, errors.New("job not found").Error(), err.Error())
	})
}

func TestAnalyzeURL(t *testing.T) {

	// no mocks here, let the code to flow through actual modules

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
			<head><title>Test Page</title></head>
			<body>
				<h1>Main Heading</h1>
				<h2>Sub Heading</h2>
				<form action="/login" method="POST">
					<input type="text" name="user">
					<input type="password" name="pass">
					<input type="submit" value="Login">
				</form>
				<a href="/internal">Internal</a>
				<a href="https://external.com">External</a>
			</body>
		</html>`
		w.Write([]byte(html))
	}))
	defer server.Close()

	store := persistance.NewMemoryStore()
	executor := process.NewPatternExecutor()
	service := NewAnalyzerService(store, executor)

	result, err := service.AnalyzeURL(server.URL)

	assert.NoError(t, err)

	assert.Equal(t, "HTML5", result.HTMLVersion)
	assert.Equal(t, "Test Page", result.Title)

	// headings
	expectedHeadings := map[string]int{
		"h1": 1, "h2": 1, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
	}
	assert.Equal(t, expectedHeadings, result.Headings)

	// login
	assert.True(t, result.LoginFormFound)

	// links
	assert.Contains(t, result.Links, "internal")
	assert.Contains(t, result.Links, "external")
	assert.Contains(t, result.Links, "broken")
}
