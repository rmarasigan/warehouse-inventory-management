package mysql

import (
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
