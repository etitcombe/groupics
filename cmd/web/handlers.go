package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expiresValue := r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	var expires int

	if strings.TrimSpace(expiresValue) == "" {
		errors["expires"] = "This field cannot be blank"
	} else {
		expires, err = strconv.Atoi(expiresValue)
		if err != nil {
			errors["expires"] = "This field must be an integer"
		} else {
			if expires != 365 && expires != 7 && expires != 1 {
				errors["expires"] = "This field is invalid (only 365, 7, or 1 allowed)"
			}
		}
	}

	if len(errors) > 0 {
		vm := createViewModel{
			FormData:   r.PostForm,
			FormErrors: errors,
		}
		app.render(w, r, "create.page.tmpl", vm)
		return
	}

	id, err := app.snippetStore.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/show/%d", id), http.StatusSeeOther)
}

func (app *application) createForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", createViewModel{})
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	err := app.parseTemplates()
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
