package apischema

type Storage struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func NewStorage(data []byte) ([]Storage, error) {
	return unmarshal[Storage](data)
}
