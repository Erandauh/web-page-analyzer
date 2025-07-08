package patterns

/*
	This pattern analyzes the html doc version
*/
import (
	"log"
)

type HeadingCounterPattern struct{}

func init() {
	Register(&HeadingCounterPattern{})
}

func (p *HeadingCounterPattern) Name() string {
	return "headings"
}

func (p *HeadingCounterPattern) Apply(ctx *Context, result map[string]any) error {

	log.Printf("[%s] Starting heading count for URL: %s", p.Name(), ctx.URL.String())

	headings := map[string]int{
		"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
	}
	for tag := range headings {
		headings[tag] = ctx.Document.Find(tag).Length()
	}

	result[p.Name()] = headings
	log.Printf("[%s] Heading count result: %+v", p.Name(), headings)

	return nil
}
