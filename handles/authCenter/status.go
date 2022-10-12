package authCenter

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/anyingiit/UnifyAuthCenter/models"
	"github.com/google/uuid"
)

// 提供给用户的状态检查
func Status(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/authCenter/status.tmpl")
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
		uuid string
	}
	cookieStr := r.Header.Get("Cookie")
	if cookieStr == "" {
		return isUnLogined(w, "cookie is empty")
	}
	cookie := Cookie{
		uuid: strings.Split(cookieStr, "uuid=")[1],
	}
	if cookie.uuid == "" {
		return isUnLogined(w, "uuid is empty")
	}
	UUID, err := uuid.Parse(cookie.uuid)
	if err != nil {
		return isUnLogined(w, "uuid type error")
	}
	session := &models.Session{
		UUID: UUID,
	}
	result := session.First()
	if result.Error != nil {
		return isUnLogined(w, "uuid invalid or expired")
	}

	if session.ExpiredAt.UnixNano() < time.Now().UnixNano() {
		return isUnLogined(w, "uuid expired")
	}

	return isLogined(w)
}
