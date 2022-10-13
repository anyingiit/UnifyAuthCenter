package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyingiit/UnifyAuthCenter/models"
	"github.com/anyingiit/UnifyAuthCenter/myErrors"
	"github.com/anyingiit/UnifyAuthCenter/utils"
	"github.com/google/uuid"
)

// roleId is role.User or role.Admin
func Authorization(handle appHandler, parameteGeter func(*http.Request) (data interface{}, error error), validatar func(data interface{}) error) MiddlewareFunc {
	return MiddlewareFunc(func() (originHandle appHandler, middleware appHandler) {
		return handle, appHandler(func(w http.ResponseWriter, r *http.Request) error {
			parameterData, err := parameteGeter(r)
			if err != nil {
				return err
			}

			err = validatar(parameterData)
			if err != nil {
				return err
			}

			return nil
		})
	})
}

func AuthorizationUser(handle appHandler) MiddlewareFunc {
	type Data struct {
		UUID uuid.UUID
	}
	return Authorization(handle, func(r *http.Request) (interface{}, error) {
		type Cookie struct {
			uuid string
		}
		cookieStr := r.Header.Get("Cookie")
		if cookieStr == "" {
			return nil, myErrors.NewSimpleAuthorizationError("user unauthorized", "cookie is empty")
		}
		cookie := Cookie{
			uuid: utils.ParseCookieValue(cookieStr, "uuid"),
		}
		if cookie.uuid == "" {
			return nil, myErrors.NewSimpleAuthorizationError("user unauthorized", "uuid is empty")
		}
		UUID, err := uuid.Parse(cookie.uuid)
		if err != nil {
			return nil, myErrors.NewSimpleAuthorizationError("user unauthorized", "uuid type error")
		}

		return &Data{UUID: UUID}, nil
	}, func(i interface{}) error {
		data, ok := i.(*Data)
		if !ok {
			return fmt.Errorf("interface{} failed cover to *Data")
		}
		session := &models.Session{
			UUID: data.UUID,
		}
		result := session.First()
		if result.Error != nil {
			return myErrors.NewSimpleAuthorizationError("user unauthorized", "uuid invalid or expired")
		}

		if session.ExpiredAt.UnixNano() < time.Now().UnixNano() {
			return myErrors.NewSimpleAuthorizationError("user unauthorized", "uuid expired")
		}
		return nil
	})
}

func AuthorizationAdmin(handle appHandler) MiddlewareFunc {
	type Data struct {
		AdminPassword string
	}
	return Authorization(handle, func(r *http.Request) (interface{}, error) {
		type Cookie struct {
			admin_password string
		}
		cookieStr := r.Header.Get("Cookie")
		if cookieStr == "" {
			return nil, myErrors.NewSimpleAuthorizationError("admin unauthorized", "cookie is empty")
		}
		cookie := Cookie{
			admin_password: utils.ParseCookieValue(cookieStr, "admin_password"),
		}
		if cookie.admin_password == "" {
			return nil, myErrors.NewSimpleAuthorizationError("admin unauthorized", "admin_uuid is empty")
		}

		return &Data{AdminPassword: cookie.admin_password}, nil
	}, func(i interface{}) error {
		data, ok := i.(*Data)
		if !ok {
			return fmt.Errorf("interface{} failed cover to *Data")
		}

		if !utils.IsValidateAdminiPassword(data.AdminPassword) {
			return myErrors.NewSimpleAuthorizationError("admin unauthorized", "admin_password invalid")
		}

		return nil
	})
}

// TODO finish AuthorizationSystem
func AuthorizationSystem(handle appHandler) MiddlewareFunc {
	return Authorization(handle, func(r *http.Request) (interface{}, error) {
		return nil, fmt.Errorf("need finish parameteGeter")
	}, func(i interface{}) error {
		return fmt.Errorf("need finish validatar")
	})
}
