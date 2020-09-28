package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))
	mux.Get("/create", app.session.Enable(http.HandlerFunc(app.createForm)))
	mux.Post("/create", app.session.Enable(http.HandlerFunc(app.create)))
	mux.Get("/refresh", http.HandlerFunc(app.refresh))
	mux.Get("/show/:id", app.session.Enable(http.HandlerFunc(app.show)))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/create", app.create)
	// mux.HandleFunc("/refresh", app.refresh)
	// mux.HandleFunc("/show", app.show)

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
