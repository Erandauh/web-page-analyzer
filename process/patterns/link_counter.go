package patterns

/*
	This pattern analyzes the internal/external and broken links
*/

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type LinkCounterPattern struct {
	Client *http.Client
}

func init() {
	Register(&LinkCounterPattern{})
}

// func NewLinkCounterPattern(client *http.Client) *LinkCounterPattern {
// 	if client == nil {
// 		client = &http.Client{Timeout: 5 * time.Second}
// 	}
// 	return &LinkCounterPattern{
// 		Client: client,
// 	}
// }

func (p *LinkCounterPattern) Name() string {
	return "links"
}

func (p *LinkCounterPattern) Apply(ctx *Context, result map[string]any) error {

	log.Printf("[%s] Starting link analysis for URL: %s", p.Name(), ctx.URL.String())
	if p.Client == nil {
		p.Client = &http.Client{Timeout: 5 * time.Second}
		log.Printf("[%s] No HTTP client provided. Using default with 5s timeout", p.Name())
	}

	internal, external, broken := 0, 0, 0
	ctx.Document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href == "" || strings.HasPrefix(href, "javascript:") {
			log.Printf("[%s] Skipping invalid or JavaScript href: %s", p.Name(), href)
			return
		}

		linkURL, err := url.Parse(href)
		if err != nil {
			log.Printf("[%s] Failed to parse link: %s, error: %v", p.Name(), href, err)
			return
		}

		fullURL := ctx.URL.ResolveReference(linkURL)
		log.Printf("[%s] Checking link: %s", p.Name(), fullURL.String())
		resp, err := p.Client.Head(fullURL.String())

		if err != nil || (resp != nil && resp.StatusCode >= 400) {
			log.Printf("[%s] Broken link: %s (error: %v, status: %v)", p.Name(), fullURL.String(), err, getStatus(resp))
			broken++
			if resp != nil {
				resp.Body.Close()
			}
			return
		}
		if resp != nil {
			resp.Body.Close()
		}

		if fullURL.Hostname() == ctx.URL.Hostname() || fullURL.Hostname() == "" {
			internal++
			log.Printf("[%s] Internal link: %s", p.Name(), fullURL.String())
		} else {
			external++
			log.Printf("[%s] External link: %s", p.Name(), fullURL.String())
		}
	})

	result[p.Name()] = map[string]int{
		"internal": internal,
		"external": external,
		"broken":   broken,
	}

	log.Printf("[%s] Link analysis complete. Total: %d, Internal: %d, External: %d, Broken: %d",
		p.Name(), (internal + external + broken), internal, external, broken)

	return nil
}

func getStatus(resp *http.Response) int {
	if resp == nil {
		return 0
	}
	return resp.StatusCode
}
