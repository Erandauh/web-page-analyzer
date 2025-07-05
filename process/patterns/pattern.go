package patterns

/*
	Common interface for any pattern
	This gives you capability to just add new pattern when you need different behavior
*/

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Context struct {
	HTML     string
	URL      *url.URL
	Document *goquery.Document
}

type Pattern interface {
	Name() string
	Apply(ctx *Context, result map[string]any) error
}
