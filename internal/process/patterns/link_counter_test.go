package patterns

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestLinkCounterPattern_Apply(t *testing.T) {

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// broken link (won’t respond)
	brokenLink := "http://localhost:9999/does-not-exist"

	html := `
	<html><body>
		<a href="` + server.URL + `/internal">Internal</a>
		<a href="` + brokenLink + `">Broken</a>
		<a href="javascript:void(0)">Invalid</a>
	</body></html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	assert.NoError(t, err)

	baseURL, err := url.Parse("http://127.0.0.1")
	assert.NoError(t, err)

	ctx := &Context{
		URL:      baseURL,
		HTML:     html,
		Document: doc,
	}

	result := make(map[string]any)

	pattern := &LinkCounterPattern{
		Client: server.Client(),
	}

	err = pattern.Apply(ctx, result)
	assert.NoError(t, err)

	// Validate
	linkCounts := result[pattern.Name()].(map[string]int)

	assert.Equal(t, 1, linkCounts["internal"], "internal link count")
	assert.Equal(t, 0, linkCounts["external"], "external link count")
	assert.Equal(t, 1, linkCounts["broken"], "broken link count")
}

func TestLinkCounterPatternName(t *testing.T) {
	type testCase struct {
		name     string
		pattern  *LinkCounterPattern
		expected string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.name, func(t *testing.T) {
			actualString := tc.pattern.Name()

			assert.Equal(t, tc.expected, actualString)
		})
	}

	validate(t, &testCase{
		name:     "should return correct pattern name",
		pattern:  &LinkCounterPattern{},
		expected: "links",
	})
}
