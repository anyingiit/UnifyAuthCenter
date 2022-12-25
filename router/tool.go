package router

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/handles/tool/generationTOTP"
	"github.com/anyingiit/UnifyAuthCenter/middleware"
)

func regisTool(serverMux *http.ServeMux) {
	// 工具相关
	// 生成TOTP
	serverMux.HandleFunc("/tool/generation_TOTP", func(w http.ResponseWriter, r *http.Request) { // 重定向到欢迎页面
		w.Header().Add("Location", "/tool/generation_TOTP/welcome")
		w.WriteHeader(http.StatusFound)
	})
	serverMux.HandleFunc("/tool/generation_TOTP/welcome", middleware.ErrWrapper(generationTOTP.WelcomeStaticPage, nil)) // 欢迎页面
	serverMux.HandleFunc("/tool/generation_TOTP/form", middleware.ErrWrapper(generationTOTP.FormStaticPage, nil))       // 表单页面
	serverMux.HandleFunc("/tool/generation_TOTP/generation", middleware.ErrWrapper(generationTOTP.Handle, nil))         // 生成TOTP页面
}
