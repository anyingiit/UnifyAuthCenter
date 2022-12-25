package router

import "net/http"

func RegisRouters(serverMux *http.ServeMux) {
	regisRoot(serverMux)
	regisHome(serverMux)
	regisAdmin(serverMux)
	regisAuthCenter(serverMux)
	regisInternal(serverMux)
	regisStatic(serverMux)
	regisTool(serverMux)
}
