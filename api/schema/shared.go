package apischema

import "encoding/json"

type Shared struct {
	Note string `json:"note"`
}

func unmarshal[T any](data []byte) ([]T, error) {
	var (
		single T
		list   []T
	)

	err := json.Unmarshal(data, &list)
	if err == nil {
		return list, err
	}

	err = json.Unmarshal(data, &single)
	if err != nil {
		return nil, err
	}

	return []T{single}, nil
}
