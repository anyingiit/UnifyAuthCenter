package authCenter

import (
	"html/template"
	"net/http"
)

func LoginStaticPage(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/authCenter/form.tmpl")
	if err != nil {
		return err
	}

	// 如果没有显式调用WriteHeader, 那么在第一次Write的时候, 将自动调用WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}
