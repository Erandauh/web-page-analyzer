package patterns

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestApply_HeadingCounterPattern_Successfully(t *testing.T) {
	type testCase struct {
		name     string
		html     string
		expected map[string]int
	}

	tests := []testCase{
		{
			name: "All heading levels",
			html: `
				<html><body>
					<h1>Header 1</h1>
					<h2>Header 2</h2>
					<h3>Header 3</h3>
					<h4>Header 4</h4>
					<h5>Header 5</h5>
					<h6>Header 6</h6>
				</body></html>`,
			expected: map[string]int{
				"h1": 1, "h2": 1, "h3": 1, "h4": 1, "h5": 1, "h6": 1,
			},
		},
		{
			name: "Some headings missing",
			html: `
				<html><body>
					<h1>Header 1</h1>
					<h2>Header 2</h2>
				</body></html>`,
			expected: map[string]int{
				"h1": 1, "h2": 1, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
			},
		},
		{
			name: "No headings",
			html: `<html><body><p>No headings here</p></body></html>`,
			expected: map[string]int{
				"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
			},
		},
	}

	pattern := &HeadingCounterPattern{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tc.html))
			assert.NoError(t, err)

			ctx := &Context{
				HTML:     tc.html,
				Document: doc,
			}

			result := make(map[string]any)

			err = pattern.Apply(ctx, result)
			assert.NoError(t, err)

			got := result[pattern.Name()].(map[string]int)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestHeadingCounterPatternName(t *testing.T) {
	type testCase struct {
		name     string
		pattern  *HeadingCounterPattern
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
		pattern:  &HeadingCounterPattern{},
		expected: "headings",
	})
}
