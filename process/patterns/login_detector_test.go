package patterns

import (
	"strings"
	"testing"

	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestLoginDetectorPattern_Apply(t *testing.T) {
	tests := []struct {
		name           string
		html           string
		expectedResult bool
	}{
		{
			name: "Login form working",
			html: `
				<html><body>
					<form action="/login">
						<input type="text" name="username">
						<input type="password" name="password">
						<input type="submit" value="Login">
					</form>
				</body></html>`,
			expectedResult: true,
		},
		{
			name: "Form without password",
			html: `
				<html><body>
					<form>
						<input type="text" name="email">
						<input type="submit">
					</form>
				</body></html>`,
			expectedResult: false,
		},
		{
			name: "No login form ",
			html: `
				<html><body>
					<div>No form here</div>
				</body></html>`,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tc.html))
			assert.NoError(t, err)

			ctx := &Context{
				Document: doc,
				URL:      &url.URL{Scheme: "http", Host: "test.com"},
				HTML:     tc.html,
			}

			result := make(map[string]any)

			p := &LoginDetectorPattern{}
			err = p.Apply(ctx, result)
			assert.NoError(t, err)

			val := result[p.Name()].(bool)
			assert.Equal(t, tc.expectedResult, val)
		})
	}
}

func TestLoginDetectorPatternName(t *testing.T) {
	type testCase struct {
		name     string
		pattern  *LoginDetectorPattern
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
		pattern:  &LoginDetectorPattern{},
		expected: "login_form_found",
	})
}
