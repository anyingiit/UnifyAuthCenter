package router

import (
	"github.com/anyingiit/UnifyAuthCenter/middleware"
	"github.com/anyingiit/UnifyAuthCenter/utils"
	"net/http"
)

func regisRoot(serverMux *http.ServeMux) {
	serverMux.HandleFunc("/", middleware.ErrWrapper(func(w http.ResponseWriter, r *http.Request) error {
		utils.RedirectTo(w, "/home")

		return nil
	}, nil))
}
