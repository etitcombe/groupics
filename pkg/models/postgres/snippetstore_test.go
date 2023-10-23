package postgres

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/etitcombe/groupics/pkg/models"
	"github.com/lib/pq"
)

func TestGetSnippet(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name        string
		snippetID   int
		wantSnippet *models.Snippet
		wantError   error
	}{
		{
			name:      "Valid ID",
			snippetID: 1,
			wantSnippet: &models.Snippet{
				ID:      1,
				Title:   "Snippet Title",
				Content: "Snippet Content Has More Text",
				Created: time.Date(2020, 10, 1, 10, 50, 0, 0, time.Local),
				Expires: time.Now().Add(7 * 24 * time.Hour),
			},
			wantError: nil,
		},
		{
			name:        "Zero ID",
			snippetID:   0,
			wantSnippet: nil,
			wantError:   models.ErrNoRecord,
		},
		{
			name:        "Non-existent ID",
			snippetID:   99,
			wantSnippet: nil,
			wantError:   models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := SnippetStore{db}

			snippet, err := m.Get(tt.snippetID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if tt.wantSnippet == nil {
				if !reflect.DeepEqual(snippet, tt.wantSnippet) {
					t.Errorf("want %+v; got %+v", tt.wantSnippet, snippet)
				}
			} else {
				if snippet.ID != tt.wantSnippet.ID {
					t.Errorf("want ID %v; got %v", tt.wantSnippet.ID, snippet.ID)
				}
				if snippet.Title != tt.wantSnippet.Title {
					t.Errorf("want Title %v; got %v", tt.wantSnippet.Title, snippet.Title)
				}
				if snippet.Content != tt.wantSnippet.Content {
					t.Errorf("want Content %v; got %v", tt.wantSnippet.Content, snippet.Content)
				}
				if !snippet.Created.Equal(tt.wantSnippet.Created) {
					t.Errorf("want Created %v; got %v", tt.wantSnippet.Created, snippet.Created)
				}
				if snippet.Expires.Year() != tt.wantSnippet.Expires.Year() ||
					snippet.Expires.Month() != tt.wantSnippet.Expires.Month() ||
					snippet.Expires.Day() != tt.wantSnippet.Expires.Day() {
					t.Errorf("want Expires %v; got %v", tt.wantSnippet.Expires, snippet.Expires)
				}
			}
		})
	}
}

func TestInsertSnippet(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name             string
		title            string
		content          string
		expiresIn        int
		wantError        error
		wantErrorCode    string
		wantErrorMessage string
	}{
		{
			"Valid",
			"Valid Title",
			"Valid content",
			7,
			nil,
			"",
			"",
		},
		{
			"Title Too Long",
			"Title Too Long Title Too Long Title Too Long Title Too Long Title Too Long Title Too Long Title Too Long Title Too Long Title Too Long Title Too Long Title Too Long",
			"Valid content",
			7,
			pq.Error{},
			"22001",
			"value too long for type character varying(100)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := SnippetStore{db}

			_, err := m.Insert(tt.title, tt.content, tt.expiresIn)

			if err != tt.wantError {
				var pqError *pq.Error
				if errors.As(err, &pqError) {
					if string(pqError.Code) != tt.wantErrorCode {
						t.Errorf("want %v; got %v", tt.wantErrorCode, pqError.Code)
					}
					if pqError.Message != tt.wantErrorMessage {
						t.Errorf("want %v; got %v", tt.wantErrorMessage, pqError.Message)
					}
				} else {
					t.Errorf("want %v; got %v", tt.wantError, err)
				}
			}
		})
	}
}

func TestLatestSnippet(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	s := SnippetStore{db}
	result, err := s.Latest()
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 3 {
		t.Errorf("want %d results; got %d", 3, len(result))
	}
}
