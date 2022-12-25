package authCenter

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/utils"
)

func AuthCenter(w http.ResponseWriter, r *http.Request) error {
	cookieString := r.Header.Get("Cookie")
	if cookieString == "" {
		utils.RedirectTo(w, "./auth_center/login")
		return nil
	}
	type Cookie struct {
		uuid string
	}
	cookie := Cookie{
		uuid: utils.ParseCookieValue(cookieString, "uuid"),
	}
	if cookie.uuid == "" {
		utils.RedirectTo(w, "./auth_center/login")
		return nil
	}

	utils.RedirectTo(w, "./auth_center/status")
	return nil
}
