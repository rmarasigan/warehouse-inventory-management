package mysql

import (
	"fmt"

	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

// RoleList retrieves a list of roles.
func ListRole() ([]schema.Role, error) { return FetchItems[schema.Role](RoleTable) }

// GetRole retrieves a specific role.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func GetRoleByID(id int) (schema.Role, error) {
	return RetrieveItemByField[schema.Role](RoleTable, "id", id)
}

func GetRoleByName(name string) (schema.Role, error) {
	return RetrieveItemByField[schema.Role](RoleTable, "name", name, "LOWER(?)")
}

// NewRole inserts a new role information into the 'role' table.
//
// Parameter:
//   - role: The role information that will be inserted.
func NewRole(role schema.Role) error {
	return InsertRecord(fmt.Sprintf("INSERT INTO %s (name) VALUES (:name);", RoleTable), role)
}

// UpdateRole updates/modifies the existing role information in the 'role'
// table.
//
// Parameter:
//   - role: The role information that will be modified.
func UpdateRole(role schema.Role) error {
	return UpdateRecordByID(RoleTable, role, []string{"name"})
}

// DeleteRole deletes existing role in the 'role' table.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func DeleteRole(id int) (int64, error) { return DeleteRecordByID(RoleTable, id) }

// RoleIDExists checks if a specific role id exists in the 'role' table.
//
// Parameter:
//   - id: The unique role id in the 'role' table.
func RoleIDExists(id int) (bool, error) {
	return exists(func() (schema.Role, error) { return GetRoleByID(id) })
}

// RoleNameExists checks if a specific role exists in the 'role' table.
//
// Parameter:
//   - name: The role name that will be checked.
func RoleNameExists(name string) (bool, error) {
	return exists(func() (schema.Role, error) { return GetRoleByName(name) })
}
