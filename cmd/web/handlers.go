package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/noonacedia/sourcepaste/internal/models"
	"github.com/noonacedia/sourcepaste/internal/validators"
)

type snippetCreateFormValidations struct {
	Title   string
	Content string
	Expires int
	validators.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	snippetId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || snippetId < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(snippetId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateFormValidations{}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := snippetCreateFormValidations{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	form.CheckField(validators.NotBlank(form.Title), "title", "title shouldn't be empty")
	form.CheckField(validators.MaxChars(form.Title, 100), "title", "title cannot be longer than 100 chars")
	form.CheckField(validators.NotBlank(form.Title), "content", "content shouldn't be empty")
	form.CheckField(validators.PermittedInt([]int{1, 7, 365}, form.Expires), "expires", "expires int should be 1, 7 or 365")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}
	id, err := app.snippets.Insert(
		form.Title,
		form.Content,
		form.Expires,
	)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}
