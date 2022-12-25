package router

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/handles/internalApi"
	"github.com/anyingiit/UnifyAuthCenter/middleware"
)

func regisInternal(serverMux *http.ServeMux) {
	// 内部接口
	serverMux.HandleFunc("/internal/authorization_user_auth_status", middleware.ErrWrapper(middleware.AuthorizationInternal(internalApi.AuthorizationUserAuthStatus)()))
}
