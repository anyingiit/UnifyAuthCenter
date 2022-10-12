package authCenter

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/anyingiit/UnifyAuthCenter/models"
	"github.com/anyingiit/UnifyAuthCenter/myErrors"
	"github.com/anyingiit/UnifyAuthCenter/utils"
	"github.com/google/uuid"
)

func Handle(w http.ResponseWriter, r *http.Request) error {
	type Query struct {
		Code string
	}
	query := Query{
		Code: r.FormValue("code"),
	}
	if query.Code == "" {
		return myErrors.NewSimpleBadRequestError("auth center login failed", "code is empty")
	}
	validated := utils.ValidateTOTP(query.Code, "REMOVED-SEE-README")
	if !validated {
		return myErrors.NewSimpleBadRequestError("auth center login failed", "code is invalid")
	}

	nowTime := time.Now()
	session := &models.Session{
		UUID:      uuid.New(),
		CreatedAt: nowTime,                    // if CreatedAt value is empty time.Time obj, gorm will automatically set this to the current time
		ExpiredAt: nowTime.Add(time.Hour * 8), // Expires in 8 hours
	}
	result := session.Create()

	if result.Error != nil {
		return fmt.Errorf("create session failed, err: %s", result.Error.Error())
	}

	t, err := template.ParseFiles("./template/login/success.tmpl")
	if err != nil {
		return err
	}
	// 缓存不应存储有关客户端请求或服务器响应的任何内容，即不使用任何缓存。
	w.Header().Add("Cache-control", "no-store")
	w.Header().Add("Set-Cookie", fmt.Sprintf("uuid=%s", session.UUID.String()))
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}
