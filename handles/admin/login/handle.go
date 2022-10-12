package login

import "net/http"

func Handle(w http.ResponseWriter, r *http.Request) error {
	type Query struct {
		code string
	}
	query := Query{
		code: r.FormValue("code"),
	}
	if query.code == "" {
		//TODO
	}

	//TODO

	return nil
}
