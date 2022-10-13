package admin

import (
	"html/template"
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/utils"
)

func Status(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/admin/login/status.tmpl")
	if err != nil {
		return err
	}

	isLogined := func(w http.ResponseWriter) error {
		// 如果没有显式调用WriteHeader, 那么在第一次Write的时候, 将自动调用WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		err = t.Execute(w, struct {
			Logined         string
			UnLoginedReason string
		}{
			Logined:         "true",
			UnLoginedReason: "",
		})
		if err != nil {
			return err
		}

		return nil
	}

	isUnLogined := func(w http.ResponseWriter, reason string) error {
		w.WriteHeader(http.StatusOK)
		err = t.Execute(w, struct {
			Logined         string
			UnLoginedReason string
		}{
			Logined:         "false",
			UnLoginedReason: reason,
		})
		if err != nil {
			return err
		}

		return nil
	}
	type Cookie struct {
		admin_password string
	}
	cookieStr := r.Header.Get("Cookie")
	if cookieStr == "" {
		return isUnLogined(w, "cookie is empty")
	}
	cookie := Cookie{
		admin_password: utils.ParseCookieValue(cookieStr, "admin_password"),
	}
	if cookie.admin_password == "" {
		return isUnLogined(w, "admin_password is empty")
	}

	if !utils.IsValidateAdminiPassword(cookie.admin_password) {
		return isUnLogined(w, "admin_password invalid")
	}

	return isLogined(w)
}
