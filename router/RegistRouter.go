package router

import "net/http"

func RegisteRouters(serverMux *http.ServeMux) {
	registeAdmin(serverMux)
	registeAuthCenter(serverMux)
	registeInternal(serverMux)
	registeStatic(serverMux)
	registeTool(serverMux)
}
