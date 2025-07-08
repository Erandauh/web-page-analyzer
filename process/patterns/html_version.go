package patterns

/*
	This pattern analyzes the html doc version
*/

import (
	"log"
	"strings"
)

type HTMLVersionPattern struct{}

func init() {
	Register(&HTMLVersionPattern{})
}

func (p *HTMLVersionPattern) Name() string {
	return "html_version"
}

func (p *HTMLVersionPattern) Apply(ctx *Context, result map[string]any) error {
	log.Printf("[%s] Starting HTML version detection for URL: %s", p.Name(), ctx.URL.String())
	html := strings.ToLower(ctx.HTML)

	switch {
	case strings.Contains(html, "<!doctype html>"):
		log.Printf("[%s] Detected: HTML5", p.Name())
		result[p.Name()] = "HTML5"
	case strings.Contains(html, "html 4.01"):
		log.Printf("[%s] Detected: HTML 4.01", p.Name())
		result[p.Name()] = "HTML 4.01"
	default:
		log.Printf("[%s] HTML version unknown", p.Name())
		result[p.Name()] = "Unknown"
	}

	return nil
}
