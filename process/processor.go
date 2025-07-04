package processor

// # Core HTML parsing logic
// based on goquery

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func DetectHTMLVersion(html string) string {
	lower := strings.ToLower(html)

	switch {
	case strings.Contains(lower, "<!doctype html>"):
		return "HTML5"
	case strings.Contains(lower, "html 4.01 strict"):
		return "HTML 4.01 Strict"
	case strings.Contains(lower, "html 4.01 transitional"):
		return "HTML 4.01 Transitional"
	case strings.Contains(lower, "html 4.01 frameset"):
		return "HTML 4.01 Frameset"
	default:
		return "Unknown!"
	}
}

func CountHeadings(doc *goquery.Document) map[string]int {
	headings := map[string]int{
		"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
	}
	for tag := range headings {
		headings[tag] = doc.Find(tag).Length()
	}
	return headings
}

func CountLinks(doc *goquery.Document, base *url.URL) (int, int, int) {
	internal, external, broken := 0, 0, 0
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href == "" || strings.HasPrefix(href, "javascript:") {
			return
		}
		linkURL, err := url.Parse(href)
		if err != nil {
			return
		}
		fullURL := base.ResolveReference(linkURL).String()
		if linkURL.Host == "" || linkURL.Host == base.Host {
			internal++
		} else {
			external++
		}
		client := http.Client{Timeout: 3 * time.Second}
		resp, err := client.Head(fullURL)
		if err != nil || resp.StatusCode >= 400 {
			broken++
		}
	})
	return internal, external, broken
}

func DetectLoginForm(doc *goquery.Document) bool {
	found := false
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		hasPass := false
		s.Find("input").Each(func(i int, input *goquery.Selection) {
			t, _ := input.Attr("type")
			if strings.ToLower(t) == "password" {
				hasPass = true
			}
		})
		if hasPass {
			found = true
		}
	})
	return found
}
