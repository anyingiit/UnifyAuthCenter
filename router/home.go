package router

import (
	"github.com/anyingiit/UnifyAuthCenter/middleware"
	"html/template"
	"net/http"
)

func regisHome(serverMux *http.ServeMux) {
	serverMux.HandleFunc("/home", middleware.ErrWrapper(func(w http.ResponseWriter, r *http.Request) error {
		t, err := template.ParseFiles("./template/home/home.tmpl")
		if err != nil {
			return err
		}
		err = t.Execute(w, nil)
		if err != nil {
			return err
		}

		return nil
	}, nil))
}
