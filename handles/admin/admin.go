package admin

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/utils"
)

func Admin(w http.ResponseWriter, r *http.Request) error {
	redirectTo := func(w http.ResponseWriter, redirectTargetUrl string) {
		w.Header().Add("Location", redirectTargetUrl)
		w.WriteHeader(http.StatusFound)
	}
	cookieString := r.Header.Get("Cookie")
	if cookieString == "" {
		redirectTo(w, "./admin/login")
		return nil
	}
	type Cookie struct {
		uuid string
	}
	cookie := Cookie{
		uuid: utils.ParseCookieValue(cookieString, "admin_password"),
	}
	if cookie.uuid == "" {
		redirectTo(w, "./admin/login")
		return nil
	}

	redirectTo(w, "./admin/status")
	return nil
}
