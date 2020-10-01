package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", app.session.Enable(app.authenticate(app.noSurf(http.HandlerFunc(app.home)))))
	mux.Get("/create", app.session.Enable(app.authenticate(app.requireAuthentication(app.noSurf(http.HandlerFunc(app.createForm))))))
	mux.Post("/create", app.session.Enable(app.authenticate(app.requireAuthentication(app.noSurf(http.HandlerFunc(app.create))))))
	mux.Get("/login", app.session.Enable(app.noSurf(http.HandlerFunc(app.loginForm))))
	mux.Post("/login", app.session.Enable(app.noSurf(http.HandlerFunc(app.login))))
	mux.Post("/logout", app.session.Enable(app.authenticate(app.requireAuthentication(app.noSurf(http.HandlerFunc(app.logout))))))
	mux.Get("/ping", http.HandlerFunc(ping))
	mux.Get("/refresh", app.session.Enable(app.authenticate(app.requireAuthentication(app.noSurf(http.HandlerFunc(app.refresh))))))
	mux.Get("/show/:id", app.session.Enable(app.noSurf(http.HandlerFunc(app.show))))
	mux.Get("/signup", app.session.Enable(app.noSurf(http.HandlerFunc(app.signupForm))))
	mux.Post("/signup", app.session.Enable(app.noSurf(http.HandlerFunc(app.signup))))

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
