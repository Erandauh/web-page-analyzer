package analyze_service

/*
	Business logic layer (analysis orchestration)
*/
import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"web-page-analyzer/model"
	"web-page-analyzer/persistance"
	"web-page-analyzer/process/patterns"
)

// analyze web url async
func AnalyzeURLAsync(rawURL string, done chan<- model.AnalysisResult, errChan chan<- error) model.Job {

	job := persistance.DefaultStore.CreateJob(rawURL)

	go func() {
		result, err := AnalyzeURL(rawURL)
		if err != nil {
			persistance.DefaultStore.CompleteJob(job.ID, result, err)
			errChan <- err
			return
		}
		persistance.DefaultStore.CompleteJob(job.ID, result, nil)
		done <- result
	}()

	return job
}

// analyze web url
func AnalyzeURL(rawURL string) (model.AnalysisResult, error) {

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch URL: %s, error: %v", rawURL, err)
		return model.AnalysisResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body from URL: %s, error: %v", rawURL, err)
		return model.AnalysisResult{}, err
	}

	html := string(body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("[ERROR] Failed to parse HTML for URL: %s, error: %v", rawURL, err)
		return model.AnalysisResult{}, err
	}

	parsedURL, _ := url.Parse(rawURL)

	ctx := &patterns.Context{
		HTML:     html,
		URL:      parsedURL,
		Document: doc,
	}

	result := make(map[string]any)

	Execute(ctx, result)

	return toAnalysisResult(result, rawURL, doc), nil
}

// get analyze web url async job data
func GetAnalysisResultByID(jobID string) (model.Job, error) {
	job, ok := persistance.DefaultStore.GetJob(jobID)
	if !ok {
		return model.Job{}, errors.New("job not found")
	}

	return *job, nil
}

func toAnalysisResult(result map[string]any, rawURL string, doc *goquery.Document) model.AnalysisResult {
	final := model.AnalysisResult{}

	if val, ok := result["html_version"].(string); ok {
		final.HTMLVersion = val
	}
	if val, ok := result["headings"].(map[string]int); ok {
		final.Headings = val
	}
	if val, ok := result["links"].(map[string]int); ok {
		final.Links = val
	}
	if val, ok := result["login_form_found"].(bool); ok {
		final.LoginFormFound = val
	}
	final.Title = doc.Find("title").Text()

	log.Printf("[INFO] Finished analysis for URL: %s", rawURL)

	return final
}
