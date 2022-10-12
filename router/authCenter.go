package router

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/handles/authCenter"
	"github.com/anyingiit/UnifyAuthCenter/middleware"
)

func registeAuthCenter(serverMux *http.ServeMux) {
	// 验证中心
	serverMux.HandleFunc("/auth_center", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Location", "/auth_center/login")
		w.WriteHeader(http.StatusFound)
	})
	serverMux.HandleFunc("/auth_center/login", middleware.ErrWrapper(authCenter.LoginStaticPage))
	serverMux.HandleFunc("/auth_center/handle", middleware.ErrWrapper(authCenter.Handle))
	// http.HandleFunc("/login/success", loginPage)
	// http.HandleFunc("/login/failed", loginPage)
	serverMux.HandleFunc("/auth_center/status", middleware.ErrWrapper(authCenter.Status))
}
