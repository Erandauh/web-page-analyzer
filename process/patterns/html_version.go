package patterns

/*
	This pattern analyzes the html doc version
*/

import "strings"

type HTMLVersionPattern struct{}

func init() {
	Register(&HTMLVersionPattern{})
}

func (p *HTMLVersionPattern) Name() string {
	return "html_version"
}

func (p *HTMLVersionPattern) Apply(ctx *Context, result map[string]any) error {
	html := strings.ToLower(ctx.HTML)

	switch {
	case strings.Contains(html, "<!doctype html>"):
		result[p.Name()] = "HTML5"
	case strings.Contains(html, "html 4.01 strict"):
		result[p.Name()] = "HTML 4.01 Strict"
	case strings.Contains(html, "html 4.01 transitional"):
		result[p.Name()] = "HTML 4.01 Transitional"
	default:
		result[p.Name()] = "Unknown"
	}
	return nil
}
