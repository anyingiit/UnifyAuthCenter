package middleware

import (
	"html/template"
	"log"
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/myErrors"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func ErrWrapper(handle appHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic: %v", r)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		err := handle(w, r)

		if err == nil {
			return
		}

		responseError := myErrors.NewMyError(err, http.StatusInternalServerError, "internal error", "internal error")

		if userError, ok := err.(myErrors.UserError); ok {
			responseError.StatusCode = userError.GetStatusCode()
			responseError.Event = userError.GetEvent()
			responseError.Reason = userError.GetReason()
		}

		InternalError := func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}

		if responseError.StatusCode == http.StatusInternalServerError {
			log.Printf("internal error: %s", responseError.Err.Error())
			InternalError(w)
			return
		}

		// user error
		t, err := template.ParseFiles("./template/error/user_error.tmpl")
		if err != nil {
			log.Printf("internal error: %s", responseError.Err.Error())
			InternalError(w)
			return
		}

		log.Printf("user error: %s", responseError.Err.Error())
		w.WriteHeader(responseError.StatusCode)
		err = t.Execute(w, struct {
			Message string
			Event   string
			Reason  string
		}{
			Message: http.StatusText(responseError.StatusCode),
			Event:   responseError.Event,
			Reason:  responseError.Reason,
		})
		if err != nil {
			log.Printf("internal error: %s", responseError.Err.Error())
			InternalError(w)
			return
		}
	}
}
