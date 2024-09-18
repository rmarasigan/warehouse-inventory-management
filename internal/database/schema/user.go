package schema

import (
	"database/sql"
	"strings"

	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
)

type User struct {
	ID           int            `db:"id"`
	RoleID       int            `db:"role_id"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"last_name"`
	Email        sql.NullString `db:"email"`
	Password     string         `db:"password"`
	LastLogin    sql.NullString `db:"last_login"`
	DateCreated  string         `db:"date_created"`
	DateModified sql.NullString `db:"date_modified"`
}

func (u *User) UpdateValues(user apischema.User) {
	if user.RoleID != 0 {
		u.RoleID = user.RoleID
	}

	if strings.TrimSpace(user.FirstName) != "" {
		u.FirstName = user.FirstName
	}

	if strings.TrimSpace(user.LastName) != "" {
		u.LastName = user.LastName
	}

	if strings.TrimSpace(user.Email) != "" {
		// Valid is 'true' if String is not NULL
		u.Email = sql.NullString{String: user.Email, Valid: true}
	}

	if strings.TrimSpace(user.Password) != "" {
		u.Password = user.Password
	}
}
