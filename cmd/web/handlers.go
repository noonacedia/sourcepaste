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

type userSignupFormValidations struct {
	Name     string
	Email    string
	Password string
	validators.Validator
}

type userLoginFormValidations struct {
	Email    string
	Password string
	validators.Validator
}

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It's OK")
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

func (app *application) userSignupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupFormValidations{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}
	form := userSignupFormValidations{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}
	form.CheckField(validators.NotBlank(form.Name), "name", "name cannot be empty")
	form.CheckField(validators.NotBlank(form.Email), "email", "email cannot be empty")
	form.CheckField(validators.EmailValid(form.Email, validators.EmailRegex), "email", "email should be like example@mail.com")
	form.CheckField(validators.NotBlank(form.Password), "password", "password cannot be empty")
	form.CheckField(validators.MinChars(form.Password, 8), "password", "password cannot have less than 8 chars")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}
	err = app.users.InsertUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "email is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successfull. Please log in")
	http.Redirect(w, r, "/users/login/", http.StatusSeeOther)
}

func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginFormValidations{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}
	form := userLoginFormValidations{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}
	form.CheckField(validators.NotBlank(form.Email), "email", "email cannot be empty")
	form.CheckField(validators.EmailValid(form.Email, validators.EmailRegex), "email", "this field must be a valid email address")
	form.CheckField(validators.NotBlank(form.Password), "password", "")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
