package generationTOTP

import (
	"html/template"
	"net/http"
)

func WelcomeStaticPage(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/tool/generation_TOTP/welcome.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}
