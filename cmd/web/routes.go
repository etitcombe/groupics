package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/create", http.HandlerFunc(app.createForm))
	mux.Post("/create", http.HandlerFunc(app.create))
	mux.Get("/show/:id", http.HandlerFunc(app.show))

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
