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

func (app *application) MustMethods (h http.Handler, allowedMethods ...string) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		if checkMethods(r, allowedMethods...){
			app.MethodNotAllowed(w, r, http.MethodPost)
		}
		
	})

}