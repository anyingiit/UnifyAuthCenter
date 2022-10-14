package authCenter

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/utils"
)

func AuthCenter(w http.ResponseWriter, r *http.Request) error {
	redirectTo := func(w http.ResponseWriter, redirectTargetUrl string) {
		w.Header().Add("Location", redirectTargetUrl)
		w.WriteHeader(http.StatusFound)
	}
	cookieString := r.Header.Get("Cookie")
	if cookieString == "" {
		redirectTo(w, "./auth_center/login")
		return nil
	}
	type Cookie struct {
		uuid string
	}
	cookie := Cookie{
		uuid: utils.ParseCookieValue(cookieString, "uuid"),
	}
	if cookie.uuid == "" {
		redirectTo(w, "./auth_center/login")
		return nil
	}

	redirectTo(w, "./auth_center/status")
	return nil
}
