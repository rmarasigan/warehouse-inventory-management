package schema

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int            `db:"id"`
	RoleID       int            `db:"role_id"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"last_name"`
	Email        sql.NullString `db:"email"`
	Password     string         `db:"password"`
	LastLogin    sql.NullString `db:"last_login"`
	Active       bool           `db:"active"`
	DateCreated  time.Time      `db:"date_created"`
	DateModified sql.NullTime   `db:"date_modified"`
}
