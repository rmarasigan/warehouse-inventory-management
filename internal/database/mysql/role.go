package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

// RoleList retrieves a list of roles.
func ListRole() ([]schema.Role, error) {
	query := "SELECT * FROM role;"
	return fetch[schema.Role](query)
}

// GetRole retrieves a specific role.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func GetRole(id int) ([]schema.Role, error) {
	query := "SELECT * FROM role WHERE id = ?;"
	return fetch[schema.Role](query, id)
}

// NewRole inserts a new role information into the 'role' table.
//
// Parameter:
//   - role: The role information that will be inserted.
func NewRole(role schema.Role) error {
	query := "INSERT INTO role (id, name) VALUES (:id, :name);"

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
						SET
							name = CASE
								WHEN :name = '' THEN name
								ELSE COALESCE(:name, name)
							END
						WHERE id = :id;`

	_, err := NamedExec(query, role)

	return err
}

// DeleteRole deletes existing role in the 'role' table.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func DeleteRole(id int) (int64, error) {
	query := "DELETE FROM role WHERE id = ?;"
	return delete(query, id)
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
	query := "SELECT * FROM role WHERE name = LOWER(?);"
	return exists[schema.Role](query, name)
}
