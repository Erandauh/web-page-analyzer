package patterns

/*
	This pattern analyzes the internal/external and broken links
*/

import (
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

func (p *LinkCounterPattern) Name() string {
	return "links"
}

func (p *LinkCounterPattern) Apply(ctx *Context, result map[string]any) error {

	if p.Client == nil {
		p.Client = &http.Client{Timeout: 5 * time.Second}
	}

	internal, external, broken := 0, 0, 0
	ctx.Document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href == "" || strings.HasPrefix(href, "javascript:") {
			return
		}

		linkURL, err := url.Parse(href)
		if err != nil {
			return
		}

		fullURL := ctx.URL.ResolveReference(linkURL)
		resp, err := p.Client.Head(fullURL.String())

		if err != nil || (resp != nil && resp.StatusCode >= 400) {
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
		} else {
			external++
		}
	})

	result[p.Name()] = map[string]int{
		"internal": internal,
		"external": external,
		"broken":   broken,
	}

	return nil
}
