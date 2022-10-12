package login

import (
	"html/template"
	"net/http"
)

func LoginStaticPage(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/admin/form.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}
