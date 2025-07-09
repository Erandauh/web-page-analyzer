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
	"github.com/sirupsen/logrus"
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

	logrus.WithFields(logrus.Fields{
		"pattern": p.Name(),
		"url":     ctx.URL.String(),
	}).Info("Starting link analysis")

	if p.Client == nil {
		p.Client = &http.Client{Timeout: 5 * time.Second}
		logrus.Warn("[%s] No HTTP client provided. Using default with 5s timeout", p.Name())
	}

	internal, external, broken := 0, 0, 0
	ctx.Document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href == "" || strings.HasPrefix(href, "javascript:") {
			logrus.WithField("href", href).Debug("Skipping invalid or JavaScript href")
			return
		}

		linkURL, err := url.Parse(href)
		if err != nil {
			logrus.WithField("href", href).WithError(err).Warn("Failed to parse link")
			return
		}

		fullURL := ctx.URL.ResolveReference(linkURL)
		logrus.WithField("resolved_url", fullURL.String()).Debug("Checking link")
		resp, err := p.Client.Head(fullURL.String())

		if err != nil || (resp != nil && resp.StatusCode >= 400) {
			logrus.WithFields(logrus.Fields{
				"link":   fullURL.String(),
				"status": getStatus(resp),
				"error":  err,
			}).Warn("Broken link")
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
			logrus.WithField("link", fullURL.String()).Debug("Internal link")
		} else {
			external++
			logrus.WithField("link", fullURL.String()).Debug("External link")
		}
	})

	result[p.Name()] = map[string]int{
		"internal": internal,
		"external": external,
		"broken":   broken,
	}

	logrus.WithFields(logrus.Fields{
		"internal": internal,
		"external": external,
		"broken":   broken,
		"total":    internal + external + broken,
	}).Info("Link analysis complete")

	return nil
}

func getStatus(resp *http.Response) int {
	if resp == nil {
		return 0
	}
	return resp.StatusCode
}
