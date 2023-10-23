package postgres

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/etitcombe/groupics/pkg/models"
	"github.com/lib/pq"
)

func TestAuthenticateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name      string
		email     string
		password  string
		wantID    int
		wantError error
	}{
		{
			name:      "Valid Password",
			email:     "alice@example.com",
			password:  "olafbuddy",
			wantID:    1,
			wantError: nil,
		},
		{
			name:      "Invalid Password",
			email:     "alice@example.com",
			password:  "invalid",
			wantID:    0,
			wantError: models.ErrInvalidCredentials,
		},
		{
			name:      "Invalid Email",
			email:     "nope@example.com",
			password:  "invalid",
			wantID:    0,
			wantError: models.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserStore{db}

			id, err := m.Authenticate(tt.email, tt.password)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if id != tt.wantID {
				t.Errorf("want ID %v; got %v", tt.wantID, id)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2020, 10, 1, 10, 50, 0, 0, time.Local),
				Active:  true,
			},
			wantError: nil,
		},
		{
			name:      "Zero ID",
			userID:    0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Non-existent ID",
			userID:    2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize a connection pool to our test database, and defer a
			// call to the teardown function, so it is always run immediately
			// before this sub-test returns.
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserStore{db}

			user, err := m.Get(tt.userID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if tt.wantUser == nil {
				if !reflect.DeepEqual(user, tt.wantUser) {
					t.Errorf("want %+v; got %+v", tt.wantUser, user)
				}
			} else {
				if user.ID != tt.wantUser.ID {
					t.Errorf("want ID %v; got %v", tt.wantUser.ID, user.ID)
				}
				if user.Name != tt.wantUser.Name {
					t.Errorf("want Name %v; got %v", tt.wantUser.Name, user.Name)
				}
				if user.Email != tt.wantUser.Email {
					t.Errorf("want Email %v; got %v", tt.wantUser.Email, user.Email)
				}
				if !user.Created.Equal(tt.wantUser.Created) {
					t.Errorf("want Created %v; got %v", tt.wantUser.Created, user.Created)
				}
			}
		})
	}
}

func TestInsertUser(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}

	tests := []struct {
		name             string
		email            string
		password         string
		wantError        error
		wantErrorCode    string
		wantErrorMessage string
	}{
		{
			"Elsa Bedelsa",
			"elsa@arendelle.com",
			"olafbuddy",
			nil,
			"",
			"",
		},
		{
			"Alice Jones",
			"alice@example.com",
			"pass",
			models.ErrDuplicateEmail,
			"",
			"",
		},
		{
			"Email Too Long",
			"This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long This user name is too long",
			"pass",
			pq.Error{},
			"22001",
			"value too long for type character varying(255)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			s := UserStore{db}
			err := s.Insert(tt.name, tt.email, tt.password)

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
