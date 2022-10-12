package router

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/handles/admin/login"
	"github.com/anyingiit/UnifyAuthCenter/middleware"
)

func registeAdmin(serverMux *http.ServeMux) {
	// 管理员
	serverMux.HandleFunc("/admin/login", middleware.ErrWrapper(login.LoginStaticPage))
	serverMux.HandleFunc("/admin/login/handle", middleware.ErrWrapper(login.Handle))
}
