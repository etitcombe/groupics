package postgres

import (
	"database/sql"
	"errors"

	"github.com/etitcombe/groupics/pkg/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// UserStore manages user data.
type UserStore struct {
	DB *sql.DB
}

// Authenticate checks the email and password.
func (s *UserStore) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get returns a user based on id.
func (s *UserStore) Get(id int) (*models.User, error) {
	return nil, nil
}

// Insert creates a record.
func (s *UserStore) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO \"user\" (name, email, hashed_password, created) " +
		"VALUES ($1, $2, $3, NOW())"
	_, err = s.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" { // https://www.postgresql.org/docs/current/errcodes-appendix.html
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}
