package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.staticPath))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /home", app.home)
	mux.HandleFunc("GET /snippet/{id}", app.snippetView)
	mux.HandleFunc("POST /snippet", app.snippetCreate)

	return mux
}
