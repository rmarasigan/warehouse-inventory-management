package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/trail"
)

// NamedExec executes a named query using the provided arguments.
//
// Parameters:
//   - query: The MySQL query to execute.
//   - args: The arguments for the query.
//
// Usage:
//
//	type User {
//	  FirstName  string `db:"first_name"`
//	  LastName   string `db:"last_name"`
//	}
//	...
//	var user = User{FirstName: "John", LastName: "Doe"}
//	var query = `INSERT INTO table_name (first_name, last_name) VALUES (:first_name, :last_name)`
//
//	result, err := NamedExec(query, user)
func NamedExec(query string, args interface{}) (sql.Result, error) {
	return database.NamedExec(query, args)
}

// Exec executes a query using the provided arguments.
//
// Parameters:
//   - query: The MySQL query to execute.
//   - args: The arguments for the query.
//
// Usage:
//
//	var query = `DELETE FROM table_name WHERE id  = ?`
//
//	result, err := Exec(query, 1)
//	if err != nil {
//	  panic(err)
//	}
//
//	affected, err := result.RowsAffected()
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return database.ExecContext(context.Background(), query, args...)
}

// Get executes a query and scans the result into the destination.
//
// Parameters:
//   - destination: The destination where the result will be stored.
//   - query:       The MySQL query to execute.
//   - args:        The arguments for the query.
//
// Usage:
//
//	type User {
//	  FirstName  string `db:"first_name"`
//	  LastName   string `db:"last_name"`
//	}
//	...
//	Connect()
//	defer Close()
//
//	var user User
//	var query = `SELECT * FROM user WHERE id = ?;`
//
//	err := Get(&user, query, 1)
func Get(destination interface{}, query string, args ...interface{}) error {
	err := database.Get(destination, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			trail.Info(err.Error())
			return nil
		}

		log.Error(err.Error(), slog.String("operation", "Get"), slog.String("query", query), slog.Any("args", args))
		return err
	}

	return nil
}

// Select executes a query and scans the result into the destination slice.
//
// Parameters:
//   - destination: The destination where the result will be stored.
//   - query:       The MySQL query to execute.
//   - args:        The arguments for the query.
//
// Usage:
//
//	type User {
//	  FirstName  string `db:"first_name"`
//	  LastName   string `db:"last_name"`
//	}
//	...
//	Connect()
//	defer Close()
//
//	var users []User
//	var query = `SELECT * FROM user;`
//
//	err := Select(&users, query)
func Select(destination interface{}, query string, args ...interface{}) error {
	err := database.Select(destination, query, args...)
	if err != nil {
		log.Error(err.Error(), slog.String("operation", "Select"), slog.String("query", query), slog.Any("args", args))
		return err
	}

	return nil
}
