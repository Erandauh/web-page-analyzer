package analyze_service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"web-page-analyzer/persistance"
)

func TestGetAnalysisResultByID(t *testing.T) {

	url := "www.test.com"
	job := persistance.DefaultStore.CreateJob(url) // Save to the store

	t.Run("should return job when found", func(t *testing.T) {
		foundJob, err := GetAnalysisResultByID(job.ID)
		assert.NoError(t, err)
		assert.Equal(t, job.ID, foundJob.ID)
	})

	t.Run("should return error when job not found", func(t *testing.T) {
		_, err := GetAnalysisResultByID("non-existent-id")
		assert.Error(t, err)
		assert.Equal(t, errors.New("job not found").Error(), err.Error())
	})
}

func TestAnalyzeURL(t *testing.T) {
	// Start mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
			<head><title>Test Page</title></head>
			<body>
				<h1>Main Heading</h1>
				<h2>Sub Heading</h2>
				<form>
					<input type="text" name="user">
					<input type="password" name="pass">
				</form>
				<a href="/internal">Internal</a>
				<a href="https://external.com">External</a>
			</body>
		</html>`
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Call the function under test
	result, err := AnalyzeURL(server.URL)

	assert.NoError(t, err)

	assert.Equal(t, "HTML5", result.HTMLVersion)
	assert.Equal(t, "Test Page", result.Title)

	// Headings
	expectedHeadings := map[string]int{
		"h1": 1, "h2": 1, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
	}
	assert.Equal(t, expectedHeadings, result.Headings)

	// Login form detection
	assert.True(t, result.LoginFormFound)

	// Links: we don’t assert exact link counts here, as link checking (HEAD requests) depends on the environment.
	assert.Contains(t, result.Links, "internal")
	assert.Contains(t, result.Links, "external")
	assert.Contains(t, result.Links, "broken") // might be 0
}
