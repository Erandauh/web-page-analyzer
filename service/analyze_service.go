package analyze_service

/*
	Business logic layer (analysis orchestration)
*/
import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"

	"web-page-analyzer/model"
	"web-page-analyzer/persistance"
	"web-page-analyzer/process/patterns"
)

// analyze web url async
func AnalyzeURLAsync(rawURL string, done chan<- model.AnalysisResult, errChan chan<- error) model.Job {

	log := logrus.WithField("url", rawURL)
	log.Info("Starting async analysis")

	job := persistance.DefaultStore.CreateJob(rawURL)
	log = log.WithField("job_id", job.ID)
	log.Info("Job created")

	go func() {
		log.Info("Running analysis")
		result, err := AnalyzeURL(rawURL)
		if err != nil {
			log.WithError(err).Warn("Analysis failed")
			persistance.DefaultStore.CompleteJob(job.ID, result, err)
			errChan <- err
			return
		}

		log.Info("Analysis completed successfully")
		persistance.DefaultStore.CompleteJob(job.ID, result, nil)
		done <- result
	}()

	return job
}

// analyze web url
func AnalyzeURL(rawURL string) (model.AnalysisResult, error) {

	log := logrus.WithField("url", rawURL)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		log.WithError(err).Error("Failed to fetch URL")
		return model.AnalysisResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read response body")
		return model.AnalysisResult{}, err
	}

	html := string(body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.WithError(err).Error("Failed to parse HTML")
		return model.AnalysisResult{}, err
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.WithError(err).Error("Failed to parse URL")
		return model.AnalysisResult{}, err
	}
	log.WithField("parsed_url", parsedURL.String()).Info("Parsed URL successfully")

	ctx := &patterns.Context{
		HTML:     html,
		URL:      parsedURL,
		Document: doc,
	}

	result := make(map[string]any)

	log.Info("Executing analysis patterns...")
	Execute(ctx, result)
	log.Info("Pattern execution completed")

	finalResult := toAnalysisResult(result, rawURL, doc)
	log.Info("AnalysisResult constructed")

	return finalResult, nil
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

	logrus.WithField("url", rawURL).Info("Finished analysis for URL")

	return final
}
