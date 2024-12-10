package mysql

func fetch[T any](query string, args ...any) ([]T, error) {
	var list []T

	err := Select(&list, query, args...)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func exists[T any](query string, args ...any) (bool, error) {
	list, err := fetch[T](query, args...)
	if err != nil {
		return false, err
	}

	return (len(list) > 0), nil
}

func entityExists[T any, U any](fetchFn func(param U) ([]T, error), param U) (bool, error) {
	list, err := fetchFn(param)
	if err != nil {
		return false, err
	}

	return len(list) > 0, err
}

func delete[T any](query string, param T) (int64, error) {
	result, err := Exec(query, param)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
