package apischema

import "encoding/json"

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewRole(data []byte) ([]Role, error) {
	var (
		role  Role
		roles []Role
	)

	err := json.Unmarshal(data, &roles)
	if err == nil {
		return roles, nil
	}

	err = json.Unmarshal(data, &role)
	if err != nil {
		return nil, err
	}

	return []Role{role}, nil
}
