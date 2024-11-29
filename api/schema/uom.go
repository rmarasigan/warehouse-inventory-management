package apischema

type UOM struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func NewUOM(data []byte) ([]UOM, error) {
	return unmarshal[UOM](data)
}
