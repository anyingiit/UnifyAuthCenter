package userSession

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/anyingiit/UnifyAuthCenter/models"
)

func List(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/admin/user_session_list.tmpl")
	if err != nil {
		return err
	}

	sessions := models.Sessions{}

	result := sessions.Find()

	if result.Error != nil {
		return fmt.Errorf("excute sessions Find has err: %s", result.Error)
	}

	type Data struct {
		UUID       string
		ExpiredAt  string
		LogoutHref string
	}

	locatin, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return fmt.Errorf("get Location 'Asia/Shanghai' failed, err %s", err)
	}
	w.Header().Add("Cache-control", "no-store")
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, func() (result []Data) {
		for _, item := range sessions {
			uuid := item.UUID.String()
			result = append(result, Data{
				UUID:       uuid,
				ExpiredAt:  item.ExpiredAt.In(locatin).String(),
				LogoutHref: "./delete?uuid=" + uuid,
			})
		}
		return result
	}())
	if err != nil {
		return err
	}

	return nil
}
