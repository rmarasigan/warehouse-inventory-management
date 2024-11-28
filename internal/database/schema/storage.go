package schema

import (
	"database/sql"
	"strings"

	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
)

type Storage struct {
	ID          int            `db:"id"`
	Code        string         `db:"code"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}

func (s *Storage) UpdateValues(storage apischema.Storage) {
	if strings.TrimSpace(storage.Code) != "" {
		s.Code = storage.Code
	}

	if strings.TrimSpace(storage.Name) != "" {
		s.Name = storage.Name
	}

	if strings.TrimSpace(storage.Description) != "" {
		// Valid is 'true' if String is not NULL
		s.Description = sql.NullString{String: storage.Description, Valid: true}
	}
}
