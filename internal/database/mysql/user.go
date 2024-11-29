package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

// UserList retrieves a list of users.
func UserList() ([]schema.User, error) {
	query := "SELECT * FROM user;"
	return fetch[schema.User](query)
}

// NewUser inserts new user information into the 'user' table.
//
// Parameter:
//   - user: The user information that will be inserted.
func NewUser(user schema.User) error {
	query := `INSERT INTO user (role_id, first_name, last_name, email, password, date_created)
							 VALUES (:role_id, :first_name, :last_name, :email, :password, :date_created)`

	_, err := NamedExec(query, user)

	return err
}

// UpdateUser updates/modifies the existing user information in the user
// table.
//
// Parameters:
//   - user: The user information that will be modified.
func UpdateUser(user schema.User) error {
	query := `UPDATE user
						SET
							role_id = CASE
								WHEN :role_id = '' THEN role_id
								ELSE COALESCE(:role_id, role_id)
							END,
							first_name = CASE
								WHEN :first_name = '' THEN first_name
								ELSE COALESCE(:first_name, first_name)
							END,
							last_name = CASE
								WHEN :last_name = '' THEN last_name
								ELSE COALESCE(:last_name, last_name)
							END,
							email = CASE
								WHEN :email = '' THEN email
								ELSE COALESCE(:email, email)
							END,
							password = CASE
								WHEN :password = '' THEN password
								ELSE COALESCE(:password, password)
							END
						WHERE id = :id`

	_, err := NamedExec(query, user)

	return err
}

// DeleteUser deletes existing user in the user table.
//
// Parameter:
//   - id: The unique user id that will be deleted.
func DeleteUser(id int) (int64, error) {
	query := `DELETE FROM user WHERE id = ?`
	return delete(query, id)
}

// UserExists checks if a specific user exists in the 'user' table.
//
// Parameter:
//   - user: The user information that will be checked.
func UserExists(user schema.User) (bool, error) {
	query := `SELECT * FROM user WHERE (first_name = ? AND last_name = ?`
	return exists[schema.User](query, user.FirstName, user.LastName, user.Password)
}

// UserIDExists checks if a specific user ID exists in the 'user' table.
//
// Parameter:
//   - id: The unique user id that will be checked.
func UserIDExists(id int) (bool, error) {
	query := `SELECT * FROM user WHERE id = ?;`
	return exists[schema.User](query, id)
}
