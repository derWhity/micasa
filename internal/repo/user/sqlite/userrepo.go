// Package sqlite provides a user repository that reads and writes user data from/to a SQLite database
package sqlite

import (
	"fmt"
	"strings"

	"github.com/derWhity/micasa/internal/models"
	"github.com/derWhity/micasa/internal/repo"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const (
	duplicateErrorPrefix = "UNIQUE constraint failed"
	insertQuery          = `INSERT INTO Users(userid, name, passwordHash, fullName) VALUES(?, ?, ?, ?)`
	deleteQuery          = `DELETE FROM Users WHERE userid = ?`
	existQuery           = `SELECT COUNT(*) AS count FROM Users WHERE userid = ?`
	getByIDQuery         = `SELECT
								(userid, name, passwordHash, fullName, createdAt, updatedAt)
							FROM
								Users
							WHERE
								userid = ?`
	updateQuery = `UPDATE
						Users
					SET
						name = ?,
						fullName = ?,
						passwordHash = ?,
						updatedAt = date('now')
					WHERE
						userid = ?`
)

// UserRepo provides a simple in-memory user storage
type UserRepo struct {
	db *sqlx.DB
}

// New creates a new user repository instance
func New(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func handleSqliteError(err error, defaultMessage string) error {
	if strings.HasPrefix(err.Error(), duplicateErrorPrefix) {
		return repo.ErrDuplicate
	}
	return errors.Wrap(err, defaultMessage)
}

// Create creates a new user
func (r *UserRepo) Create(u *models.User) error {
	u.ID = models.UserID(uuid.NewV4().String())
	// Fix the user's name to lowercase
	u.Name = strings.ToLower(u.Name)
	_, err := r.db.Exec(insertQuery, string(u.ID), u.Name, u.PasswordHash, u.FullName)
	if err != nil {
		return handleSqliteError(err, "Failed to insert user")
	}
	return nil
}

// Update updates an existing user
func (r *UserRepo) Update(u *models.User) error {
	// Fix the user's name to lowercase
	u.Name = strings.ToLower(u.Name)
	// Check if the user exists
	if exists, err := r.Exists(u.ID); err != nil || !exists {
		if !exists {
			return repo.ErrNotExisting
		}
		return err
	}
	// Update
	if _, err := r.db.Exec(updateQuery, u.Name, u.FullName, u.PasswordHash, u.ID); err != nil {
		return handleSqliteError(err, "Failed to update user")
	}
	return nil
}

// Delete removes an existing user from the user storage
func (r *UserRepo) Delete(id models.UserID) error {
	if _, err := r.db.Exec(deleteQuery, string(id)); err != nil {
		return errors.Wrapf(err, "Failed to delete user with ID #%d", id)
	}
	return nil
}

// Exists checks if the user with the given ID exists in the database
func (r *UserRepo) Exists(id models.UserID) (bool, error) {
	var num int64
	if err := r.db.Get(&num, existQuery, id); err != nil {
		return false, handleSqliteError(err, "Failed to check for user existence")
	}
	return num != 0, nil
}

// GetByID returns the user with the given ID
func (r *UserRepo) GetByID(id uint) (*models.User, error) {
	return nil, fmt.Errorf("Not implemented")
}

// GetByCredentials returns the user which has the given username and password - this is used for login
func (r *UserRepo) GetByCredentials(username string, password string) (*models.User, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Find searches for users matching the given search string - supports pagination
func (r *UserRepo) Find(search string, offset uint, limit uint) ([]*models.User, error) {
	return nil, fmt.Errorf("Not implemented")
}
