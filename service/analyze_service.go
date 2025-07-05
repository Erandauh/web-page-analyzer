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
	//processor "web-page-analyzer/process"
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

	// result := model.AnalysisResult{
	// 	HTMLVersion:    processor.DetectHTMLVersion(html),
	// 	Title:          doc.Find("title").Text(),
	// 	Headings:       processor.CountHeadings(doc),
	// 	LoginFormFound: processor.DetectLoginForm(doc),
	// }

	// internal, external, broken := processor.CountLinks(doc, parsedURL)
	// result.InternalLinks = internal
	// result.ExternalLinks = external
	// result.InaccessibleLinks = broken

	final := model.AnalysisResult{}
	if val, ok := result["html_version"].(string); ok {
		final.HTMLVersion = val
	}
	if val, ok := result["title"].(string); ok {
		final.Title = val
	}
	if val, ok := result["headings"].(map[string]int); ok {
		final.Headings = val
	}
	if val, ok := result["internal_links"].(int); ok {
		final.InternalLinks = val
	}
	if val, ok := result["external_links"].(int); ok {
		final.ExternalLinks = val
	}
	if val, ok := result["inaccessible_links"].(int); ok {
		final.InaccessibleLinks = val
	}
	if val, ok := result["login_form_found"].(bool); ok {
		final.LoginFormFound = val
	}
	return final, nil

	log.Printf("[INFO] Finished analysis for URL: %s", rawURL)
	return final, nil
}
