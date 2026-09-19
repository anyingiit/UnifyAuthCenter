package authCenter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyingiit/UnifyAuthCenter/models"
	"github.com/anyingiit/UnifyAuthCenter/myErrors"
	"github.com/anyingiit/UnifyAuthCenter/utils"
	"github.com/google/uuid"
)

// The TOTP secret this handler validated against used to be a string
// literal on the line below, in a public repository. It now comes from
// utils.TOTPSecret(), which reads UNIFYAUTH_TOTP_SECRET -- see that
// function for why moving it out of the source is only half the problem.
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
	validated := utils.ValidateTOTP(query.Code, utils.TOTPSecret())
	if !validated {
		return myErrors.NewSimpleBadRequestError("auth center login failed", "code is invalid")
	}

	nowTime := time.Now()
	session := &models.Session{
		UUID:      uuid.New(),
		RoleId:    models.RoleUserId,
		CreatedAt: nowTime,                    // if CreatedAt value is empty time.Time obj, gorm will automatically set this to the current time
		ExpiredAt: nowTime.Add(time.Hour * 8), // Expires in 8 hours
	}
	result := session.Create()

	if result.Error != nil {
		return fmt.Errorf("create session failed, err: %s", result.Error.Error())
	}

	// 缓存不应存储有关客户端请求或服务器响应的任何内容，即不使用任何缓存。
	w.Header().Add("Cache-control", "no-store")
	w.Header().Add("Set-Cookie", fmt.Sprintf("uuid=%s; Path=/; Domain=anyingiit.com;", session.UUID.String()))
	w.WriteHeader(http.StatusOK)              // WriterHeader mast be after with w.Header().Xxx
	_, err := w.Write([]byte("auth success")) // if call the mothod before not call WriterHeader(...), then Status Code will be write 200
	if err != nil {
		return err
	}

	return nil
}
