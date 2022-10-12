package internalApi

import (
	"net/http"
	"strings"
	"time"

	"github.com/anyingiit/UnifyAuthCenter/models"
	"github.com/anyingiit/UnifyAuthCenter/myErrors"
	"github.com/google/uuid"
)

// 提供给服务器内部的状态检查
func AuthorizationUserAuthStatus(w http.ResponseWriter, r *http.Request) error {
	type Cookie struct {
		uuid string
	}
	cookieStr := r.Header.Get("Cookie")
	if cookieStr == "" {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "cookie is empty")
	}
	cookie := Cookie{
		uuid: strings.Split(cookieStr, "uuid=")[1],
	}
	if cookie.uuid == "" {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "uuid is empty")
	}
	UUID, err := uuid.Parse(cookie.uuid)
	if err != nil {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "uuid type error")
	}
	session := &models.Session{
		UUID: UUID,
	}
	result := session.First()
	if result.Error != nil {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "uuid invalid or expired")
	}

	if session.ExpiredAt.UnixNano() < time.Now().UnixNano() {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "uuid expired")
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("auth center authorized"))

	if err != nil {
		return err
	}

	return nil
}
