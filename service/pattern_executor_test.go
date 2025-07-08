package analyze_service

import (
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"web-page-analyzer/process/patterns"
)

func TestExecute_WithRealPattern(t *testing.T) {

	patterns.Clear()
	pattern := &patterns.HTMLVersionPattern{}
	patterns.Register(pattern)

	html := `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	assert.NoError(t, err)

	parsedURL, err := url.Parse("https://test.com")
	assert.NoError(t, err)

	ctx := &patterns.Context{
		URL:      parsedURL,
		HTML:     html,
		Document: doc,
	}

	result := make(map[string]any)

	Execute(ctx, result)

	val := result[pattern.Name()]

	assert.Equal(t, "HTML5", val)

	// Clean up registered patterns
	patterns.Clear()
}
