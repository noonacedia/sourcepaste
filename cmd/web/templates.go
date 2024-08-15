package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/noonacedia/sourcepaste/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	year := time.Now().Year()
	flash := app.sessionManager.PopString(r.Context(), "flash")
	return &templateData{CurrentYear: year, Flash: flash}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
