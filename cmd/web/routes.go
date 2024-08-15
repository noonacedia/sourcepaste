package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.staticPath))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.Handle("GET /{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippets/{id}/{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetView)))
	mux.Handle("GET /snippets/{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreateForm)))
	mux.Handle("POST /snippets/{$}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetCreate)))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
