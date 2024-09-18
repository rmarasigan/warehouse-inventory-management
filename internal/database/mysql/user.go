package mysql

import (
	"fmt"

	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
)

// UserList retrieves a list of users.
func UserList() ([]apischema.User, error) {
	var (
		list  []schema.User
		query = "SELECT * FROM user;"
	)

	err := Select(&list, query)
	if err != nil {
		return nil, err
	}

	var users = convert.Schema(list, func(user schema.User) apischema.User {
		return apischema.User{
			ID:           user.ID,
			RoleID:       user.RoleID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email.String,
			LastLogin:    user.LastLogin.String,
			DateCreated:  user.DateCreated,
			DateModified: user.DateModified.String,
		}
	})

	return users, nil
}

// RetrieveUsers retrieve a user based on the passed id parameter.
//
// Parameter:
//   - id: The unique user id that will be fetched.
func RetrieveUsers(ids string) (*[]schema.User, error) {
	var list []schema.User
	query := fmt.Sprintf(`SELECT * FROM user WHERE id IN (%v)`, ids)

	err := Select(&list, query)
	if err != nil {
		return nil, err
	}

	return &list, nil
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
						SET role_id = :role_id, first_name = :first_name, last_name = :last_name, email = :email, password = :password
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

	result, err := Exec(query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// UserExists checks if a specific user exists in the 'user' table.
//
// Parameter:
//   - user: The user information that will be checked.
func UserExists(user schema.User) (bool, error) {
	var list []schema.User
	query := `SELECT * FROM user WHERE first_name = ? AND last_name = ? AND password = ?`

	err := Select(&list, query, user.FirstName, user.LastName, user.Password)
	if err != nil {
		return false, err
	}

	return (len(list) > 0), nil
}
