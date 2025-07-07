package patterns

import (
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
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(tt.html))

			ctx := &Context{
				HTML:     tt.html,
				Document: doc,
			}

			result := make(map[string]any)

			p := &HTMLVersionPattern{}
			p.Apply(ctx, result)

			val := result[p.Name()]

			if val != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, val)
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
