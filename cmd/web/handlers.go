package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/etitcombe/groupics/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	snippets, err := app.snippetStore.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := homeViewModel{Snippets: snippets}
	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippetStore.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/show?id=%d", id), http.StatusSeeOther)
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	err := app.parseTemplates()
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
