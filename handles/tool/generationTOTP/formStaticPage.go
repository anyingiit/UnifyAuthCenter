package generationTOTP

import (
	"html/template"
	"net/http"
)

func FormStaticPage(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/tool/generation_TOTP/form.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}
