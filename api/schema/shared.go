package apischema

type Shared struct {
	UserID int32  `json:"user_id"`
	Note   string `json:"note"`
}

func NewNote(data []byte) (Shared, error) {
	shared, err := unmarshal[Shared](data)
	return shared[0], err
}
