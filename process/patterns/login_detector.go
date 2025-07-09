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
		hasEmailOrUser := false
		hasSubmit := false
		isLoginAction := false

		s.Find("input").Each(func(_ int, input *goquery.Selection) {
			inputType, _ := input.Attr("type")
			inputType = strings.ToLower(inputType)

			switch inputType {
			case "password":
				hasPass = true
			case "email", "text":
				nameAttr, _ := input.Attr("name")
				if val := strings.ToLower(nameAttr); strings.Contains(val, "email") || strings.Contains(val, "user") {
					hasEmailOrUser = true
				}
			case "submit":
				hasSubmit = true
			}
		})

		method, _ := s.Attr("method")
		if strings.ToUpper(method) == "POST" {
			isLoginAction = true
		}

		// final condition:
		if hasPass && isLoginAction && hasEmailOrUser && hasSubmit {
			formCount++
			found = true
			log.Printf("[%s] Login form candidate #%d found [pass: %v, user/email: %v, submit: %v, action: %v]", p.Name(), formCount, hasPass, hasEmailOrUser, hasSubmit, isLoginAction)
		}
	})

	log.Printf("[%s] Login forms detected: %d", p.Name(), formCount)

	result[p.Name()] = found
	log.Printf("[%s] Login detection result: %v", p.Name(), found)

	return nil

}
