package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"forum/app/config"
	"forum/app/templates"
)

// Opens a beautiful HTML 404 web page instead of the status 404 "Page not found"
func NotFound(app *config.Application, w http.ResponseWriter, r *http.Request) {
	app.ErrLog.Printf("wrong path: %s", r.URL.Path)

	w.WriteHeader(http.StatusNotFound) // Sets status code at 404
	templates.ExecuteError(app, w, r, http.StatusNotFound)
}

func ServerError(app *config.Application, w http.ResponseWriter, r *http.Request, message string, err error) {
	app.ErrLog.Output(2, fmt.Sprintf("fail handling the page %s: %s: %s\n%v", r.URL.Path, message, err, debug.Stack()))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func ClientError(app *config.Application, w http.ResponseWriter, r *http.Request, errStatus int, logTexterr string) {
	app.ErrLog.Output(2, logTexterr)
	http.Error(w, "ERROR: "+http.StatusText(errStatus), errStatus)
}

func MethodNotAllowed(app *config.Application, w http.ResponseWriter, r *http.Request, allowedMethods ...string) {
	if allowedMethods == nil {
		panic("no methods is given to func MethodNotAllowed")
	}
	allowdeString := allowedMethods[0]
	for i := 1; i < len(allowedMethods); i++ {
		allowdeString += ", " + allowedMethods[i]
	}

	w.Header().Set("Allow", allowdeString)
	ClientError(app, w, r, http.StatusMethodNotAllowed, fmt.Sprintf("using the method %s to go to a page %s", r.Method, r.URL))
}

func Forbidden(app *config.Application, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.ErrLog.Printf("access was forbidden: %s", r.URL.Path)

		w.WriteHeader(http.StatusForbidden) // Sets status code at 404
		templates.ExecuteError(app, w, r, http.StatusForbidden)
	}
}
