package analyze_service

/*
	Business logic layer (analysis orchestration)
*/
import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"web-page-analyzer/model"
	"web-page-analyzer/process/patterns"
)

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
	return final, nil
}
