package convert

func Schema[T any, U any](source T, fn func(T) U) U {
	return fn(source)
}

func SchemaList[T any, U any](source []T, fn func(T) U) []U {
	var destination []U

	for _, item := range source {
		destination = append(destination, fn(item))
	}

	if len(destination) == 0 {
		destination = make([]U, 0)
	}

	return destination
}
