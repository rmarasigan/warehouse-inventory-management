package mysql

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/trail"
)

func retrieve[T any](query string, args ...any) (T, error) {
	var data T

	err := database.Get(&data, query, args...)
	if err != nil {
		trail.Warn("[retrieve] %s: %s", err.Error(), query)

		if errors.Is(err, sql.ErrNoRows) {
			return data, nil
		}

		return data, err
	}

	return data, nil
}

func fetch[T any](query string, args ...any) ([]T, error) {
	var list []T

	err := database.Select(&list, query, args...)
	if err != nil {
		trail.Warn("[fetch] %s: %s", err.Error(), query)
		return nil, err
	}

	return list, nil
}

func exists[T any](fn func() (T, error)) (bool, error) {
	result, err := fn()
	if err != nil {
		return false, err
	}

	value := reflect.ValueOf(result)
	if !value.IsValid() {
		return false, nil
	}

	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		// Verify if it has elements
		return value.Len() > 0, nil

	case reflect.Ptr, reflect.Interface:
		// Verify if it is not nil
		return !value.IsNil(), nil

	default:
		// Verify if the value is not zero value
		return !value.IsZero(), nil
	}
}

func delete[T any](query string, param T) (int64, error) {
	result, err := database.ExecContext(context.Background(), query, param)
	if err != nil {
		trail.Warn("[delete] %s: %s", err.Error(), query)
		return 0, err
	}

	return result.RowsAffected()
}
