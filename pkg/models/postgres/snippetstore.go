package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/etitcombe/groupics/pkg/models"
)

// SnippetStore represent the actions that can be taken about Snippets.
type SnippetStore struct {
	DB *sql.DB
}

// Get returns a Snippet by its id.
func (s *SnippetStore) Get(id int) (*models.Snippet, error) {
	snippet := &models.Snippet{}

	stmt := "SELECT id, title, content, created, expires " +
		"FROM snippet " +
		"WHERE expires > NOW() AND id = $1"
	row := s.DB.QueryRow(stmt, id)

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return snippet, nil
}

// Insert inserts a record.
func (s *SnippetStore) Insert(title, content string, expires int) (int, error) {
	var id int

	stmt := fmt.Sprintf("INSERT INTO snippet (title, content, created, expires) "+
		"VALUES ($1, $2, NOW(), NOW() + INTERVAL '%d day') RETURNING id", expires)
	err := s.DB.QueryRow(stmt, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Latest returns the last 10 most recently created snippets.
func (s *SnippetStore) Latest() ([]*models.Snippet, error) {
	stmt := "SELECT id, title, content, created, expires " +
		"FROM snippet " +
		"WHERE expires > NOW() " +
		"ORDER BY created DESC LIMIT 10"
	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		snippet := &models.Snippet{}
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
