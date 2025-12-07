package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/trail"
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

func FetchItemsByFields[T any](table string, conditions map[string]any, args ...any) ([]T, error) {
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

	return fetch[T](query, args...)
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

// InsertRecord creates a new record in the specified table using the provided data and field names.
//
// Parameters:
//   - table: The name of the database where the record will be inserted.
//   - record: The record to be inserted.
//   - fields: A list of field names that specify which columns will be populated.
//
// Usage:
//
//	record := User{
//	  Name:  "John Doe",
//	  Email: "j.doe.email@example.com",
//	}
//
//	err := InsertRecord(TableName, record, "name", "email")
func InsertRecord(table string, record any, fields ...string) (int64, error) {
	if len(fields) == 0 {
		return 0, fmt.Errorf("must specify at least one field to perform insert operation")
	}

	var values []string

	// Build the VALUES clause dynamically
	for _, field := range fields {
		values = append(values, fmt.Sprintf(":%s", field))
	}

	// Construct the INSERT query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(fields, ", "),
		strings.Join(values, ", "),
	)

	result, err := database.NamedExec(query, record)
	if err != nil {
		trail.Error("[insert] %s: %s", err.Error(), query)
		return 0, err
	}

	return result.LastInsertId()
}

// InsertIfNotExists inserts a record into a table if a specific field value does not already exist.
//
// Parameters:
//   - table: The name of the database to update.
//   - record: The record to update, represented as a struct.
//   - uniqueField: The field to check for existence (e.g., "name").
//   - fields: The columns to insert (must include uniqueField).
func InsertIfNotExists(table string, record any, uniqueField string, fields ...string) (int64, error) {
	if len(fields) == 0 {
		return 0, fmt.Errorf("must specify at least one field")
	}

	if uniqueField == "" {
		return 0, fmt.Errorf("uniqueField must be specified")
	}

	// Build column names and placeholders
	columns := strings.Join(fields, ", ")

	var placeholders []string
	for _, field := range fields {
		placeholders = append(placeholders, ":"+field)
	}

	// Build the SQL query using WHERE NOT EXISTS
	query := fmt.Sprintf(
		`INSERT INTO %s (%s)
		 SELECT %s
		 WHERE NOT EXISTS (
		   SELECT 1 FROM %s WHERE %s = :%s
		 );`,
		table,
		columns,
		strings.Join(placeholders, ", "),
		table,
		uniqueField,
		uniqueField,
	)

	trail.Info("query: %v", query)

	result, err := database.NamedExec(query, record)
	if err != nil {
		trail.Error("[insert-if-not-exists] %s: %s", err.Error(), query)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		// Record already exists
		trail.Info("record already exist")
		return 0, nil
	}

	return result.LastInsertId()
}

// UpdateRecordByID updates a specific record in the given table by its ID. The fields are
// to be updated are dynamically specified in the 'fields' slice.
//
// Parameters:
//   - table: The name of the database to update.
//   - record: The record to update, represented as a struct.
//   - fields: A list of field names to be included in the UPDATE query.
//
// Usage:
//
//	record := User{
//	  ID:    1,
//	  Email: "new.email@example.com",
//	}
//
//	err := UpdateRecordByID(TableName, record, "email")
func UpdateRecordByID(table string, record any, fields ...string) error {
	if len(fields) == 0 {
		return fmt.Errorf("must specify at least one field to perform update operation")
	}

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
	if err != nil {
		trail.Error("[update] %s: %s", err.Error(), query)
		return err
	}

	return nil
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
func Exec(query string, args ...any) (sql.Result, error) {
	return database.ExecContext(context.Background(), query, args...)
}
