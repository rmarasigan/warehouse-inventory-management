package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

// RoleList retrieves a list of roles.
func ListRole() ([]schema.Role, error) {
	return fetch[schema.Role]("SELECT * FROM role;")
}

// GetRole retrieves a specific role.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func GetRole(id int) ([]schema.Role, error) {
	return fetch[schema.Role]("SELECT * FROM role WHERE id = ?;", id)
}

// NewRole inserts a new role information into the 'role' table.
//
// Parameter:
//   - role: The role information that will be inserted.
func NewRole(role schema.Role) error {
	_, err := NamedExec("INSERT INTO role (name) VALUES (:name);", role)
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
	return delete("DELETE FROM role WHERE id = ?;", id)
}

// RoleIDExists checks if a specific role id exists in the 'role' table.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func RoleIDExists(id int) (bool, error) {
	return entityExists(GetRole, id)
}

// RoleNameExists checks if a specific role exists in the 'role' table.
//
// Parameter:
//   - name: The role name that will be checked.
func RoleNameExists(name string) (bool, error) {
	return exists[schema.Role]("SELECT * FROM role WHERE name = LOWER(?);", name)
}
