package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/etitcombe/groupics/pkg/forms"
	"github.com/etitcombe/groupics/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Not needed when using pat
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	snippets, err := app.snippetStore.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := homeViewModel{Snippets: snippets}
	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	// Not needed when using pat
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", createViewModel{Form: form})
		return
	}

	title := form.Get("title")
	content := form.Get("content")
	expires, _ := strconv.Atoi(form.Get("expires")) // We can ignore the error because we've already validated the value of the field.

	id, err := app.snippetStore.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Created successfully!")

	http.Redirect(w, r, fmt.Sprintf("/show/%d", id), http.StatusSeeOther)
}

func (app *application) createForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", createViewModel{Form: forms.New(nil)})
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "login.page.tmpl", createViewModel{Form: form})
		return
	}

	id, err := app.userStore.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or password is incorrect")
			app.render(w, r, "login.page.tmpl", createViewModel{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	http.Redirect(w, r, "/create", http.StatusSeeOther)
}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", createViewModel{Form: forms.New(nil)})
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")

	app.session.Put(r, "flash", "You&rsquo;ve been logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", createViewModel{Form: form})
		return
	}

	err = app.userStore.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", createViewModel{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", createViewModel{Form: forms.New(nil)})
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	err := app.parseTemplates()
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippetStore.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := showViewModel{Snippet: s}
	app.render(w, r, "show.page.tmpl", data)
}
