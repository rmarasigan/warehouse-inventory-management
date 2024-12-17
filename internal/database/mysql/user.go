package mysql

import (
	"fmt"

	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

// UserList retrieves a list of users.
func UserList() ([]schema.User, error) { return FetchItems[schema.User](UserTable) }

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
func NewUser(user schema.User) error {
	query := fmt.Sprintf(`INSERT INTO %s (role_id, first_name, last_name, email, password, date_created)
						VALUES (:role_id, :first_name, :last_name, :email, :password, :date_created);`, UserTable)

	return InsertRecord(query, user)
}

// UpdateUser updates/modifies the existing user information in the user
// table.
//
// Parameters:
//   - user: The user information that will be modified.
func UpdateUser(user schema.User) error {
	args := []string{"role_id", "first_name", "last_name", "email", "password"}
	return UpdateRecordByID(UserTable, user, args)
}

// DeleteUser deletes existing user in the user table.
//
// Parameter:
//   - id: The unique user id that will be deleted.
func DeleteUser(id int) (int64, error) { return DeleteRecordByID(UserTable, id) }

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
