package models

import (
	"errors"
	"time"
)

// ErrNoRecord represents the error when no data is found for a query.
var ErrNoRecord = errors.New("models: no matching record found")

// Snippet represent a snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
