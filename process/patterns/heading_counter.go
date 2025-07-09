package patterns

import (
	"github.com/sirupsen/logrus"
)

/*
	This pattern analyzes the html doc version
*/

type HeadingCounterPattern struct{}

func init() {
	Register(&HeadingCounterPattern{})
}

func (p *HeadingCounterPattern) Name() string {
	return "headings"
}

func (p *HeadingCounterPattern) Apply(ctx *Context, result map[string]any) error {

	logrus.WithFields(logrus.Fields{
		"pattern": p.Name(),
		"url":     ctx.URL.String(),
	}).Info("Starting heading count")

	headings := map[string]int{
		"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
	}
	for tag := range headings {
		headings[tag] = ctx.Document.Find(tag).Length()
	}

	result[p.Name()] = headings
	logrus.WithFields(logrus.Fields{
		"pattern": p.Name(),
		"result":  headings,
	}).Info("Heading count completed")

	return nil
}
