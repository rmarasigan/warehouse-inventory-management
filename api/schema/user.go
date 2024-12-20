package apischema

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	RoleID       int       `json:"role_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email,omitempty"`
	Password     string    `json:"password,omitempty"`
	LastLogin    string    `json:"last_login,omitempty"`
	Active       bool      `json:"active"`
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified,omitempty"`
}

func NewUser(data []byte) ([]User, error) {
	return unmarshal[User](data)
}
