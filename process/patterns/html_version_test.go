package patterns

import (
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestApply_HTMLVersionPattern_Successfully(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5",
			html:     "<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>",
			expected: "HTML5",
		},
		{
			name: "html 4.01",
			html: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"
			        "http://www.w3.org/TR/html4/strict.dtd"><html><body></body></html>`,
			expected: "HTML 4.01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to create document: %v", err)
			}

			url, err := url.Parse("https://test.com")
			if err != nil {
				t.Fatalf("invalid URL: %v", err)
			}

			ctx := &Context{
				HTML:     tt.html,
				Document: doc,
				URL:      url,
			}

			result := make(map[string]any)

			p := &HTMLVersionPattern{}
			err = p.Apply(ctx, result)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			valStr, ok := result[p.Name()].(string)
			if !ok {
				t.Fatalf("Expected result to be string, got: %T", result[p.Name()])
			}

			if valStr != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, valStr)
			}
		})
	}
}

func TestHTMLVersionPatternName(t *testing.T) {
	type testCase struct {
		name     string
		pattern  *HTMLVersionPattern
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
		pattern:  &HTMLVersionPattern{},
		expected: "html_version",
	})
}
