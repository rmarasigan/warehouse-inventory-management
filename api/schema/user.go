package apischema

import (
	"time"
)

type User struct {
	ID           int    `json:"id"`
	RoleID       int    `json:"role_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	LastLogin    string `json:"last_login,omitempty"`
	DateCreated  string `json:"date_created,omitempty"`
	DateModified string `json:"date_modified,omitempty"`
}

func (user User) SetDateCreated() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NewUser(data []byte) ([]User, error) {
	return unmarshal[User](data)
}
