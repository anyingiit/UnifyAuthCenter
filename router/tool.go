package router

import (
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/handles/tool/generationTOTP"
	"github.com/anyingiit/UnifyAuthCenter/middleware"
)

func registeTool(serverMux *http.ServeMux) {
	// 工具相关
	// 生成TOTP
	serverMux.HandleFunc("/tool/generation_TOTP", func(w http.ResponseWriter, r *http.Request) { // 重定向到欢迎页面
		w.Header().Add("Location", "/tool/generation_TOTP/welcome")
		w.WriteHeader(http.StatusFound)
	})
	serverMux.HandleFunc("/tool/generation_TOTP/welcome", middleware.ErrWrapper(generationTOTP.WelcomeStaticPage)) // 欢迎页面
	serverMux.HandleFunc("/tool/generation_TOTP/form", middleware.ErrWrapper(generationTOTP.FormStaticPage))       // 表单页面
	serverMux.HandleFunc("/tool/generation_TOTP/generation", middleware.ErrWrapper(generationTOTP.Handle))         // 生成TOTP页面
}
