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

// Authenticate checks the for an active user with the given email and password.
func (s *UserStore) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM \"user\" WHERE email = $1 AND active = true"
	err := s.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

// Get returns a user based on id.
func (s *UserStore) Get(id int) (*models.User, error) {
	var u models.User
	stmt := "SELECT \"name\", email, created, active " +
		"FROM \"user\" WHERE id = $1"
	err := s.DB.QueryRow(stmt, id).Scan(&u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	u.ID = id
	return &u, nil
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
