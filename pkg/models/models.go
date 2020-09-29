package models

import (
	"errors"
	"time"
)

var (
	// ErrNoRecord represents the error when no data is found for a query.
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials represents the error when credentials are invalid.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail represents the error when the email is a duplicate.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

// Snippet represent a snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// User represnet an application user.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
