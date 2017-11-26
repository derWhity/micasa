// Package repo contains the repository interfaces needed for MiCasa.
// The implementation of the repos can be found in the sub-packages
package repo

import (
	"errors"

	"github.com/derWhity/micasa/internal/models"
)

var (
	// ErrDuplicate is an error that is returned when a repo operation failed because of an already-existing
	// duplicate entry
	ErrDuplicate = errors.New("Duplicate record")
	// ErrNotExisting is an error that is returned when a repo is asked to return a specific item, but the item does
	// not exist inside the repo
	ErrNotExisting = errors.New("Record does not exist")
)

// UserRepo defines a repository that is able to store, query for and authenticate users
type UserRepo interface {
	// Create creates a new user
	Create(u *models.User) error
	// Update updates an existing user
	Update(u *models.User) error
	// Delete removes an existing user from the user storage
	Delete(id models.UserID) error
	// GetByID returns the user with the given ID
	GetByID(id models.UserID) (*models.User, error)
	// GetByCredentials returns the user which has the given username and password - this is used for login
	GetByCredentials(username string, password string) (*models.User, error)
	// Find searches for users matching the given search string - supports pagination
	Find(search string, offset uint, limit uint) ([]*models.User, error)
	// Check if the user exists
	Exists(id models.UserID) (bool, error)
}
