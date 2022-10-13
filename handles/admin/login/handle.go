package login

import (
	"fmt"
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/myErrors"
	"github.com/anyingiit/UnifyAuthCenter/utils"
)

func Handle(w http.ResponseWriter, r *http.Request) error {
	type Query struct {
		password string
	}
	query := Query{
		password: r.FormValue("password"),
	}
	if query.password == "" {
		return myErrors.NewSimpleBadRequestError("admin login failed", "password is empty")
	}

	if !utils.IsValidateAdminiPassword(query.password) {
		return myErrors.NewSimpleAuthorizationError("admin login failed", "password invalid")
	}

	w.Header().Add("Cache-control", "no-store")
	w.Header().Add("Set-Cookie", fmt.Sprintf("admin_password=%s; Path=/admin;", query.password))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("admin auth success"))

	if err != nil {
		return err
	}
	return nil
}
