package router

import "net/http"

func RegisRouters(serverMux *http.ServeMux) {
	regisAdmin(serverMux)
	regisAuthCenter(serverMux)
	regisInternal(serverMux)
	regisStatic(serverMux)
	regisTool(serverMux)
}
