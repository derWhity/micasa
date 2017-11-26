package models

import (
	"fmt"
	"time"

	"github.com/elithrar/simple-scrypt"
)

// UserID is a string-based user ID
type UserID string

// User defines an user of the application and his/her permissions inside this application
type User struct {
	// Internal user ID
	ID UserID `db:"userid"`
	// The user name used to log-in
	Name string `db:"name"`
	// The hashed password for authentication
	PasswordHash string `db:"passwordHash"`
	// The full user name for display reasons
	FullName string `db:"fullName"`
	// Creation time
	CreatedAt time.Time `db:"createdAt"`
	// Last update time
	UpdatedAt time.Time `db:"updatedAt"`
}

// SetPassword sets a new password creating a password hash from the incoming password and storing it in the user's
// PasswordHash property
func (u *User) SetPassword(pass string) error {
	hash, err := scrypt.GenerateFromPassword([]byte(pass), scrypt.DefaultParams)
	if err != nil {
		return fmt.Errorf("SetPassword: Error during password hashing: %v", err)
	}
	// The library already uses a string encoding here - so there is no need to encode further
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword checks if the given password corresponds to the hash stored in the user struct.
// It returns an error if the password does not match or an error occurs when loading the password hash from the user
func (u *User) CheckPassword(pass string) error {
	return scrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pass))
}
