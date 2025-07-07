package patterns

/*
	This pattern analyzes the login forms
*/

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LoginDetectorPattern struct{}

func init() {
	Register(&LoginDetectorPattern{})
}

func (p *LoginDetectorPattern) Name() string {
	return "login_form_found"
}

func (p *LoginDetectorPattern) Apply(ctx *Context, result map[string]any) error {
	found := false
	ctx.Document.Find("form").Each(func(i int, s *goquery.Selection) {
		hasPass := false
		s.Find("input").Each(func(i int, input *goquery.Selection) {
			t, _ := input.Attr("type")
			if strings.ToLower(t) == "password" {
				hasPass = true
			}
		})
		if hasPass {
			found = true
		}
	})

	result[p.Name()] = found

	return nil

}
