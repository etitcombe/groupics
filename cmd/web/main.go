package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/etitcombe/groupics/pkg/models"
	"github.com/etitcombe/groupics/pkg/models/postgres"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	session      *sessions.Session
	snippetStore interface {
		Get(int) (*models.Snippet, error)
		Insert(string, string, int) (int, error)
		Latest() ([]*models.Snippet, error)
	}
	userStore interface {
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
		Insert(string, string, string) error
	}
	templateCache map[string]*template.Template
	templateDir   string
}

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":4000", "HTTP network address")
	var dbHost, dbPassword string
	flag.StringVar(&dbHost, "dbhost", "localhost", "database host")
	flag.StringVar(&dbPassword, "dbpassword", "", "database password")
	var secret string
	flag.StringVar(&secret, "secret", "", "secret key (base64 URL-encoded)")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(fmt.Sprintf("dbname=groupics user=groupics password=%s host=%s port=5432 connect_timeout=10 sslmode=disable", dbPassword, dbHost))
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	if secret == "" {
		secret = "gIjJT1zsGLN_KKmuEM0o7Bp-wFy9W8YK5FZRnFh7BVM=" // got this value from rand.RememberToken in exp
	}
	secretBytes, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New(secretBytes)
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		session:      session,
		snippetStore: &postgres.SnippetStore{DB: db},
		userStore:    &postgres.UserStore{DB: db},
		templateDir:  "./ui/html/",
	}
	err = app.parseTemplates()
	if err != nil {
		errorLog.Fatal(err)
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		infoLog.Println("Starting server on", addr)
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		if errors.Is(err, http.ErrServerClosed) {
			infoLog.Println(err)
		} else {
			errorLog.Println(err)
		}
		infoLog.Println("Bye!")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	s := <-c

	infoLog.Printf("received signal %v; shutting down...\n", s)

	if err = srv.Shutdown(context.Background()); err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

/* https://pkg.go.dev/github.com/lib/pq
* dbname - The name of the database to connect to
* user - The user to sign in as
* password - The user's password
* host - The host to connect to. Values that start with / are for unix
  domain sockets. (default is localhost)
* port - The port to bind to. (default is 5432)
* sslmode - Whether or not to use SSL (default is require, this is not
  the default for libpq)
	* disable - No SSL
	* require - Always SSL (skip verification)
	* verify-ca - Always SSL (verify that the certificate presented by the
	server was signed by a trusted CA)
	* verify-full - Always SSL (verify that the certification presented by
	the server was signed by a trusted CA and the server host name
	matches the one in the certificate)
* fallback_application_name - An application_name to fall back to if one isn't provided.
* connect_timeout - Maximum wait for connection, in seconds. Zero or
  not specified means wait indefinitely.
*/
