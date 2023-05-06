package handlers

import (
	"net/http"

	"forum/app/application"
)

func checkMethods(r *http.Request, methods ...string) bool {
	for _, mth := range methods {
		if r.Method == mth {
			return true
		}
	}
	return false
}

/*
MustMethods wrapper makes sure that the request's method is allowed
*/
func MustMethods(app *application.Application, h http.Handler, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !checkMethods(r, allowedMethods...) {
			MethodNotAllowed(app, w, r, allowedMethods...)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func NotAuth(app *application.Application, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ses, err := checkLoggedin(app, w, r)
		if err != nil {
			// checkLoggedin has already written error status to w
			return
		}
		if ses.LoginStatus == loggedin {
			w.Header().Add("Location", "/")
			w.WriteHeader(204)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func Signs(app *application.Application, h http.HandlerFunc, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MustMethods(app, NotAuth(app,h), allowedMethods...).ServeHTTP(w, r)
	})
}
