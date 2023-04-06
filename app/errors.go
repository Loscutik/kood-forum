package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
)

// Opens a beautiful HTML 404 web page instead of the status 404 "Page not found"
func (app *application) NotFound(w http.ResponseWriter, r *http.Request) {
	app.errLog.Printf("wrong path: %s", r.URL.Path)

	w.WriteHeader(http.StatusNotFound)                             // Sets status code at 404
	tm, _ := template.ParseFiles(TEMPLATES_PATH + "error404.html") // Opens the HTML web page
	err := tm.Execute(w, nil)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (app *application) ServerError(w http.ResponseWriter, r *http.Request, message string, err error) {
	app.errLog.Output(2, fmt.Sprintf("fail handling the page %s: %s: %s\n%s", r.URL.Path, message, err, debug.Stack()))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *application) ClientError(w http.ResponseWriter, r *http.Request, errStatus int, logTexterr string) {
	app.errLog.Output(2, logTexterr)
	http.Error(w, http.StatusText(errStatus), errStatus)
}

func (app *application) MethodNotAllowed(w http.ResponseWriter, r *http.Request, allowdeString string) {
	w.Header().Set("Allow", allowdeString)
	app.ClientError(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("using the method %s to go to a page %s", r.Method, r.URL))
}

func (app *application) Forbidden(w http.ResponseWriter, r *http.Request) {
	app.errLog.Printf("access was forbidden: %s", r.URL.Path)

	w.WriteHeader(http.StatusForbidden)                             // Sets status code at 404
	tm, _ := template.ParseFiles(TEMPLATES_PATH + "forbiden.html") // Opens the HTML web page
	err := tm.Execute(w, nil)
	if err != nil {
		app.ClientError(w, r,http.StatusForbidden, fmt.Sprintf("forbidden execute failed: %s",err))
	}
}