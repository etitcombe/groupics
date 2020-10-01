package postgres

import (
	"reflect"
	"testing"
	"time"

	"github.com/etitcombe/groupics/pkg/models"
)

func TestGet(t *testing.T) {
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
