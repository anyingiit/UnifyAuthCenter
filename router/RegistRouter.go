package router

import "net/http"

func RegisteRouters(serverMux *http.ServeMux) {
	registeAdmin(serverMux)
	registeAuthCenter(serverMux)
	registeStatic(serverMux)
	registeTool(serverMux)
}
