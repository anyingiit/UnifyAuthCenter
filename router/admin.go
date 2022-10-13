package router

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/handles/admin"
	"github.com/anyingiit/UnifyAuthCenter/handles/admin/login"
	"github.com/anyingiit/UnifyAuthCenter/middleware"
)

func registeAdmin(serverMux *http.ServeMux) {
	// 管理员
	serverMux.HandleFunc("/admin/login", middleware.ErrWrapper(login.LoginStaticPage, nil))
	serverMux.HandleFunc("/admin/login/handle", middleware.ErrWrapper(login.Handle, nil))
	serverMux.HandleFunc("/admin/status", middleware.ErrWrapper(admin.Status, nil))
	serverMux.HandleFunc("/admin/session_list", middleware.ErrWrapper(middleware.AuthorizationAdmin(login.LoginStaticPage)()))
}
