package analyze_service

// # Business logic layer (analysis orchestration)

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"web-page-analyzer/model"
	processor "web-page-analyzer/process"
)

func AnalyzeURL(rawURL string) (model.AnalysisResult, error) {

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		return model.AnalysisResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.AnalysisResult{}, err
	}

	html := string(body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return model.AnalysisResult{}, err
	}

	parsedURL, _ := url.Parse(rawURL)
	result := model.AnalysisResult{
		HTMLVersion:    processor.DetectHTMLVersion(html),
		Title:          doc.Find("title").Text(),
		Headings:       processor.CountHeadings(doc),
		LoginFormFound: processor.DetectLoginForm(doc),
	}

	internal, external, broken := processor.CountLinks(doc, parsedURL)
	result.InternalLinks = internal
	result.ExternalLinks = external
	result.InaccessibleLinks = broken

	return result, nil
}
