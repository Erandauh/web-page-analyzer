package patterns

/*
	This pattern analyzes the html doc version
*/

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type HTMLVersionPattern struct{}

func init() {
	Register(&HTMLVersionPattern{})
}

func (p *HTMLVersionPattern) Name() string {
	return "html_version"
}

func (p *HTMLVersionPattern) Apply(ctx *Context, result map[string]any) error {
	logrus.WithFields(logrus.Fields{
		"pattern": p.Name(),
		"url":     ctx.URL.String(),
	}).Info("Starting HTML version detection")

	html := strings.ToLower(ctx.HTML)

	switch {
	case strings.Contains(html, "<!doctype html>"):
		logrus.WithField("pattern", p.Name()).Info("Detected HTML5")
		result[p.Name()] = "HTML5"
	case strings.Contains(html, "html 4.01"):
		logrus.WithField("pattern", p.Name()).Info("Detected HTML 4.01")
		result[p.Name()] = "HTML 4.01"
	default:
		logrus.WithField("pattern", p.Name()).Warn("HTML version unknown")
		result[p.Name()] = "Unknown"
	}

	return nil
}
