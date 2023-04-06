package main

import "net/http"

func (app *application) routers() *http.ServeMux {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", app.homePageHandler)
	mux.HandleFunc("/signup", app.signupPageHandler)
	mux.HandleFunc("/signup/success", app.signupSuccessPageHandler)
	mux.HandleFunc("/login", app.signinPageHandler)
	mux.HandleFunc("/userinfo/", app.userPageHandler)
	mux.HandleFunc("/post/", app.postPageHandler)
	mux.HandleFunc("/addpost", app.addPostPageHandler)
	mux.HandleFunc("/post/create", app.postCreatorHandler)
	mux.HandleFunc("/liking", app.likingHandler)
	mux.HandleFunc("/logout", app.logoutHandler)

	fileServer := http.FileServer(http.Dir(STATIC_PATH))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	return mux
}
