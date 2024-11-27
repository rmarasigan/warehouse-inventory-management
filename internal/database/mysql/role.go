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

// GetRole retrieves a specific role.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func GetRole(id int) ([]apischema.Role, error) {
	var list []schema.Role
	query := `SELECT * FROM role WHERE id = ?;`

	err := Select(&list, query, id)
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

// UpdateRole updates/modifies the existing role information in the 'role'
// table.
//
// Parameter:
//   - role: The role information that will be modified.
func UpdateRole(role schema.Role) error {
	query := `UPDATE role
						SET name = :name
						WHERE id = :id;`

	_, err := NamedExec(query, role)

	return err
}

// DeleteRole deletes existing role in the 'role' table.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func DeleteRole(id int) (int64, error) {
	query := `DELETE FROM role WHERE id = ?`

	result, err := Exec(query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// RoleIDExists checks if a specific role id exists in the 'role' table.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func RoleIDExists(id int) (bool, error) {
	roles, err := GetRole(id)
	if err != nil {
		return false, err
	}

	return (len(roles) > 0), nil
}

// RoleNameExists checks if a specific role exists in the 'role' table.
//
// Parameter:
//   - name: The role name that will be checked.
func RoleNameExists(name string) (bool, error) {
	var list []schema.Role
	query := `SELECT * FROM role WHERE name = LOWER(?);`

	err := Select(&list, query, strings.ToLower(name))
	if err != nil {
		return false, err
	}

	return (len(list) > 0), nil
}
