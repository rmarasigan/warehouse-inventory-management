package mysql

import (
	"strings"

	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
)

// RoleList retrieves a list of roles.
func RoleList() ([]apischema.Role, error) {
	var (
		list  []schema.Role
		query = "SELECT * FROM role;"
	)

	err := Select(&list, query)
	if err != nil {
		return nil, err
	}

	var roles = convert.Schema(list, func(role schema.Role) apischema.Role {
		return apischema.Role{
			ID:   role.ID,
			Name: role.Name,
		}
	})

	return roles, nil
}

// NewRole inserts a new role information into the 'role' table.
//
// Parameter:
//   - role: The role information that will be inserted.
func NewRole(role schema.Role) error {
	query := `INSERT INTO role (id, name) VALUES (:id, :name);`

	_, err := NamedExec(query, role)

	return err
}

// RoleNameExists checks if a specific role exists in the 'role' table.
//
// Parameter:
//   - role: The role information that will be checked.
func RoleNameExists(role schema.Role) (bool, error) {
	var list []schema.Role
	query := `SELECT * FROM role WHERE name = LOWER(?);`

	err := Select(&list, query, strings.ToLower(role.Name))
	if err != nil {
		return false, err
	}

	return (len(list) > 0), nil
}
