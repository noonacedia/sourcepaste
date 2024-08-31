package main

import (
	"net/http"

	"github.com/noonacedia/sourcepaste/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/{filepath...}", fileServer)
	mux.HandleFunc("GET /ping/{$}", app.ping)
	mux.Handle("GET /{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.home)))))
	mux.Handle("GET /about/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.about)))))

	mux.Handle("GET /snippets/{id}/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.snippetView)))))
	mux.Handle("GET /snippets/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.snippetCreateForm))))))
	mux.Handle("POST /snippets/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.snippetCreate))))))

	mux.Handle("GET /users/signup/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.userSignupForm)))))
	mux.Handle("POST /users/signup/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.userSignup)))))
	mux.Handle("GET /users/login/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.userLoginForm)))))
	mux.Handle("POST /users/login/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(http.HandlerFunc(app.userLogin)))))
	mux.Handle("POST /users/logout/{$}", app.sessionManager.LoadAndSave(noSurf(app.authenticate(app.requireAuthentication(http.HandlerFunc(app.userLogout))))))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
