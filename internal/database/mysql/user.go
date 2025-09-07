package mysql

import (
	"fmt"

	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

// ListUser retrieves a list of users.
func ListUser() ([]schema.User, error) { return FetchItems[schema.User](UserTable) }

func GetUserByID(id int) (schema.User, error) {
	return RetrieveItemByField[schema.User](UserTable, "id", id)
}

func GetUserByName(firstName, lastName string) (schema.User, error) {
	conditions := map[string]any{
		"first_name": "?",
		"last_name":  "?",
	}

	return RetrieveItemByFields[schema.User](UserTable, conditions, firstName, lastName)
}

// NewUser inserts new user information into the 'user' table.
//
// Parameter:
//   - user: The user information that will be inserted.
func NewUser(user schema.User) (int64, error) {
	return InsertRecord(
		UserTable,
		user,
		"role_id",
		"first_name",
		"last_name",
		"email",
		"password",
		"is_active",
	)
}

// UpdateUser updates/modifies the existing user information in the user
// table.
//
// Parameters:
//   - user: The user information that will be modified.
func UpdateUser(user schema.User) error {
	return UpdateRecordByID(
		UserTable,
		user,
		"role_id",
		"first_name",
		"last_name",
		"email",
		"password",
	)
}

func ActivateUser(id int) error {
	query := fmt.Sprintf("UPDATE %s SET is_active = true WHERE id = ?", UserTable)
	_, err := Exec(query, id)

	return err
}

// DeleteUser updates the existing user 'is_active' field as 'false' in the user table.
//
// Parameter:
//   - id: The unique user id that will be deactivated.
func DeleteUser(id int) error {
	query := fmt.Sprintf("UPDATE %s SET is_active = false WHERE id = ?", UserTable)
	_, err := Exec(query, id)

	return err
}

// UserExists checks if a specific user exists in the 'user' table.
//
// Parameter:
//   - user: The user information that will be checked.
func UserExists(user schema.User) (bool, error) {
	return exists(func() (schema.User, error) { return GetUserByName(user.FirstName, user.LastName) })
}

// UserIDExists checks if a specific user ID exists in the 'user' table.
//
// Parameter:
//   - id: The unique user id that will be checked.
func UserIDExists(id int) (bool, error) {
	return exists(func() (schema.User, error) { return GetUserByID(id) })
}
