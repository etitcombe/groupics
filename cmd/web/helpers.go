package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) addDefaultData(r *http.Request, data interface{}) interface{} {
	switch vm := data.(type) {
	default:
		app.infoLog.Println("addDefaultData: data is not a viewModel, yo")
		return data
	case createViewModel:
		vm.CSRFToken = nosurf.Token(r)
		vm.Flash = template.HTML(app.session.PopString(r, "flash"))
		vm.IsAuthenticated = app.isAuthenticated(r)
		vm.Year = time.Now().Year()
		return vm
	case homeViewModel:
		vm.CSRFToken = nosurf.Token(r)
		vm.Flash = template.HTML(app.session.PopString(r, "flash"))
		vm.IsAuthenticated = app.isAuthenticated(r)
		vm.Year = time.Now().Year()
		return vm
	case showViewModel:
		vm.CSRFToken = nosurf.Token(r)
		vm.Flash = template.HTML(app.session.PopString(r, "flash"))
		vm.IsAuthenticated = app.isAuthenticated(r)
		vm.Year = time.Now().Year()
		return vm
	case viewModel:
		vm.CSRFToken = nosurf.Token(r)
		vm.Flash = template.HTML(app.session.PopString(r, "flash"))
		vm.IsAuthenticated = app.isAuthenticated(r)
		vm.Year = time.Now().Year()
		return vm
	}
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) parseTemplates() error {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(app.templateDir, "*.page.tmpl"))
	if err != nil {
		return err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return err
		}
		ts, err = ts.ParseGlob(filepath.Join(app.templateDir, "*.layout.tmpl"))
		if err != nil {
			return err
		}
		ts, err = ts.ParseGlob(filepath.Join(app.templateDir, "*.partial.tmpl"))
		if err != nil {
			return err
		}

		cache[name] = ts
	}

	app.templateCache = cache
	return nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := bytes.Buffer{}

	err := ts.Execute(&buf, app.addDefaultData(r, data))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
