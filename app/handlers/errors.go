package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"forum/app/application"
	"forum/app/templates"
)

// Opens a beautiful HTML 404 web page instead of the status 404 "Page not found"
func NotFound(app *application.Application, w http.ResponseWriter, r *http.Request) {
	app.ErrLog.Printf("wrong path: %s", r.URL.Path)

	w.WriteHeader(http.StatusNotFound) // Sets status code at 404
	if err:=templates.ExecuteError(w, r, http.StatusNotFound);err!=nil{
		app.ErrLog.Println(err)
		http.NotFound(w,r)
	}
}

func ServerError(app *application.Application, w http.ResponseWriter, r *http.Request, message string, err error) {
	app.ErrLog.Output(2, fmt.Sprintf("fail handling the page %s: %s: %s\n%v", r.URL.Path, message, err, debug.Stack()))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func ClientError(app *application.Application, w http.ResponseWriter, r *http.Request, errStatus int, logTexterr string) {
	app.ErrLog.Output(2, logTexterr)
	http.Error(w, "ERROR: "+http.StatusText(errStatus), errStatus)
}

func MethodNotAllowed(app *application.Application, w http.ResponseWriter, r *http.Request, allowedMethods ...string) {
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

func Forbidden(app *application.Application, w http.ResponseWriter, r *http.Request) {
	app.ErrLog.Printf("access was forbidden: %s", r.URL.Path)

	w.WriteHeader(http.StatusForbidden) // Sets status code at 403
	if err:=templates.ExecuteError(w, r, http.StatusForbidden);err!=nil{
		app.ErrLog.Println(err)
		http.Error(w, fmt.Sprintf("ERROR: %s. ", http.StatusText(http.StatusForbidden)), http.StatusForbidden)
	}
}
