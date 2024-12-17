package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// FetchItems retrieves all records from the specified table and returns them as a slice of the specified
// type.
//
// Parameters:
//   - table: The name of the database to query.
//
// Usage:
//
//	data, err := FetchItems[schema.MySchema](TableName)
func FetchItems[T any](table string) ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s;", table)
	return fetch[T](query)
}

// RetrieveItemByField is a generic function to retrieve single record from a database by a single field
// with an optional transformation. If no transformation is provided, the default is no transformation ("?").
//
// Usage:
//
//	RetrieveItemByField[schema.MySchema](TableName, "field_name", field_value)
//	RetrieveItemByField[schema.MySchema](TableName, "field_name", field_value, "LOWER(?)")
func RetrieveItemByField[T any](table string, field string, value any, transformation ...string) (T, error) {
	// Default transformation to "?" if none provided
	trans := "?"
	if len(transformation) > 0 {
		trans = transformation[0]
	}

	return RetrieveItemByFields[T](table, map[string]any{field: trans}, value)
}

// RetrieveItemByFields retrieves a single record based on multiple conditions with optional transformations.
// Conditions map field names to transformations, defaulting to "?" if no transformation is provided.
//
// Usage:
//
//	conditions := map[string]any{
//	  "field_name_1": "?",        // Without transformation
//	  "field_name_2": "LOWER(?)", // With transformation
//	}
//
//	data, err := RetrieveItemByFields[schema.MySchema](TableName, conditions, fieldOne, fieldTwo)
func RetrieveItemByFields[T any](table string, conditions map[string]any, args ...any) (T, error) {
	var whereClauses []string

	// Build the WHERE clause dynamically
	for field, transformation := range conditions {
		if transformation == "" {
			// Default to no transformation
			transformation = "?"
		}

		whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", field, transformation))
	}

	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s;",
		table,
		strings.Join(whereClauses, " AND "),
	)

	return retrieve[T](query, args...)
}

func InsertRecord[T any](query string, record T) error {
	_, err := database.NamedExec(query, record)
	return err
}

// UpdateRecordByID updates a specific record in the given table by its ID. The fields are
// to be updated are dynamically specified in the 'fields' slice.
//
// Parameters:
//   - table: The name of the database to update.
//   - record: The record to update, represented as a struct.
//   - fields: A slice of field names to be included in the UPDATE query.
//
// Usage:
//
//	record := User{
//	  ID:    1,
//	  Email: "new.email@example.com",
//	}
//
//	err := UpdateRecordByID(TableName, record, []string{"email"})
func UpdateRecordByID(table string, record any, fields []string) error {
	var setClause []string

	// Build the SET clause dynamically for fields to update
	for _, field := range fields {
		clause := fmt.Sprintf("%s = CASE WHEN :%s = '' THEN %s ELSE COALESCE(:%s, %s) END", field, field, field, field, field)
		setClause = append(setClause, clause)
	}

	// Construct the UPDATE query
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = :id;",
		table,
		strings.Join(setClause, ", "),
	)

	_, err := database.NamedExec(query, record)

	return err
}

func DeleteRecordByID(table string, id int) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?;", table)
	return delete(query, id)
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
