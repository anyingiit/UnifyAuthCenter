package utils

import "net/http"

func RedirectTo(w http.ResponseWriter, redirectTargetUrl string) {
	w.Header().Add("Location", redirectTargetUrl)
	w.WriteHeader(http.StatusFound)
}
