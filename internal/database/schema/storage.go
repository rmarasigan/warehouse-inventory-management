package schema

import (
	"database/sql"
)

type Storage struct {
	ID          int            `db:"id"`
	Code        string         `db:"code"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}
