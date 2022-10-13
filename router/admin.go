package router

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/handles/admin"
	"github.com/anyingiit/UnifyAuthCenter/handles/admin/login"
	"github.com/anyingiit/UnifyAuthCenter/handles/admin/userSession"
	"github.com/anyingiit/UnifyAuthCenter/middleware"
)

func registeAdmin(serverMux *http.ServeMux) {
	// 管理员
	serverMux.HandleFunc("/admin/login", middleware.ErrWrapper(login.LoginStaticPage, nil))
	serverMux.HandleFunc("/admin/login/handle", middleware.ErrWrapper(login.Handle, nil))
	serverMux.HandleFunc("/admin/status", middleware.ErrWrapper(admin.Status, nil))
	serverMux.HandleFunc("/admin/user_session/list", middleware.ErrWrapper(middleware.AuthorizationAdmin(userSession.List)()))
	serverMux.HandleFunc("/admin/user_session/delete", middleware.ErrWrapper(middleware.AuthorizationAdmin(userSession.Delete)()))
}
