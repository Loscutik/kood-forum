package main

import (
	"net/http"

	"forum/app/application"
	"forum/app/handlers"
	"forum/app/templates"
)

func routers(app *application.Application) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", handlers.HomePageHandler(app))
	mux.Handle("/signup", handlers.Signs(app, handlers.SignupPageHandler(app), http.MethodPost))
	mux.Handle("/signup/success", handlers.SignupSuccessPageHandler(app))
	mux.Handle("/login", handlers.Signs(app, handlers.SigninPageHandler(app), http.MethodPost))
	mux.Handle("/userinfo/", handlers.UserPageHandler(app))
	mux.Handle("/settings", handlers.SettingsPageHandler(app))
	mux.Handle("/post/", handlers.PostPageHandler(app))
	mux.Handle("/addpost", handlers.AddPostPageHandler(app))
	mux.Handle("/post/create", handlers.PostCreatorHandler(app))
	mux.Handle("/liking", handlers.LikingHandler(app))
	mux.Handle("/logout", handlers.LogoutHandler(app))

	fileServer := http.FileServer(http.Dir(templates.STATIC_PATH))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	return mux
}
