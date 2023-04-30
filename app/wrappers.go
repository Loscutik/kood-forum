package main

import "net/http"

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
func (app *application) MustMethods(h http.Handler, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !checkMethods(r, allowedMethods...) {
			app.MethodNotAllowed(w, r, allowedMethods...)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (app *application) NotAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ses, err := app.checkLoggedin(w, r)
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

func (app *application) Signs(h http.HandlerFunc, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.MustMethods(app.NotAuth(h),allowedMethods...).ServeHTTP(w,r)
	})
}
