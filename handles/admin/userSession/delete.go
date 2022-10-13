package userSession

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/models"
	"github.com/anyingiit/UnifyAuthCenter/myErrors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Delete(w http.ResponseWriter, r *http.Request) error {
	type Query struct {
		uuid string
	}
	query := Query{
		uuid: r.FormValue("uuid"),
	}
	if query.uuid == "" {
		return myErrors.NewSimpleBadRequestError("delete user session", "query uuid is empty")
	}

	UUID, err := uuid.Parse(query.uuid)

	if err != nil {
		return myErrors.NewSimpleBadRequestError("delete user session", "uuid type invalid")
	}

	session := models.Session{
		UUID: UUID,
	}

	result := session.First()

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return myErrors.NewSimpleBadRequestError("delete user session", "not find this uuid")
		}
		return fmt.Errorf("excute sql faild, err %s", result.Error)
	}

	result = session.Delete()

	if result.Error != nil {
		return fmt.Errorf("excute sql faild, err %s", result.Error)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete OK"))

	return nil
}
