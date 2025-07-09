package patterns

/*
	This pattern analyzes the login forms
*/

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type LoginDetectorPattern struct{}

func init() {
	Register(&LoginDetectorPattern{})
}

func (p *LoginDetectorPattern) Name() string {
	return "login_form_found"
}

func (p *LoginDetectorPattern) Apply(ctx *Context, result map[string]any) error {

	logrus.WithFields(logrus.Fields{
		"pattern": p.Name(),
		"url":     ctx.URL.String(),
	}).Info("Starting login form detection")

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
			logrus.WithFields(logrus.Fields{
				"has_password":   hasPass,
				"has_email_user": hasEmailOrUser,
				"has_submit":     hasSubmit,
				"is_post":        isLoginAction,
			}).Info("Login form candidate found")
		}
	})

	logrus.WithField("count", formCount).Info("Login forms detected")

	result[p.Name()] = found
	logrus.WithField("login_detected", found).Info("Login detection result")

	return nil

}
