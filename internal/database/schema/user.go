package schema

import (
	"database/sql"
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
