package apischema

import "time"

type Item struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	StorageID    int       `json:"storage_id"`
	Quantity     int       `json:"quantity"`
	UnitPrice    float64   `json:"unit_price"`
	UoMID        int       `json:"uom_id"`
	CreatedBy    int       `json:"created_by"`
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified,omitempty"`
}

func NewItem(data []byte) ([]Item, error) {
	return unmarshal[Item](data)
}
