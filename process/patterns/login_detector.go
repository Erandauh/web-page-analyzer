package patterns

/*
	This pattern analyzes the login forms
*/

import (
	"log"
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

	log.Printf("[%s] Starting login form detection", p.Name())
	found := false
	formCount := 0

	ctx.Document.Find("form").Each(func(_ int, s *goquery.Selection) {

		hasPass := false
		s.Find("input").Each(func(_ int, input *goquery.Selection) {
			t, _ := input.Attr("type")
			if strings.ToLower(t) == "password" {
				hasPass = true
			}
		})
		if hasPass {
			formCount++
			found = true
		}
	})

	log.Printf("[%s] Login forms detected: %d", p.Name(), formCount)

	result[p.Name()] = found
	log.Printf("[%s] Login detection result: %v", p.Name(), found)

	return nil

}
