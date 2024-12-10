package apischema

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewRole(data []byte) ([]Role, error) {
	return unmarshal[Role](data)
}
