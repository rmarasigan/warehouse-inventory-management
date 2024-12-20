package apischema

import "time"

type Item struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Quantity     int       `json:"quantity"`
	UnitPrice    float64   `json:"unit_price"`
	UoMID        int       `json:"uom_id"`
	StockStatus  string    `json:"stock_status"`
	StorageID    int       `json:"storage_id"`
	CreatedBy    int       `json:"created_by"`
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified,omitempty"`
}

func NewItem(data []byte) ([]Item, error) {
	return unmarshal[Item](data)
}
