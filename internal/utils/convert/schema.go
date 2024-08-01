package convert

func Schema[T any, U any](source []T, fn func(T) U) []U {
	var destination []U

	for _, item := range source {
		destination = append(destination, fn(item))
	}

	return destination
}
