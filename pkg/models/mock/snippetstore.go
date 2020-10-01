package mock

import (
	"time"

	"github.com/etitcombe/groupics/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

// SnippetStore is a mock of a SnippetStore.
type SnippetStore struct{}

// Get returns a Snippet by its id.
func (s *SnippetStore) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// Insert inserts a record.
func (s *SnippetStore) Insert(title, content string, expires int) (int, error) {
	return 2, nil
}

// Latest returns the last 10 most recently created snippets.
func (s *SnippetStore) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
