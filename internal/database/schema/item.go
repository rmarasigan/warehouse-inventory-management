package schema

import (
	"database/sql"
	"time"
)

type Item struct {
	ID           int            `db:"id"`
	Name         string         `db:"name"`
	Description  sql.NullString `db:"description"`
	Quantity     int            `db:"quantity"`
	UnitPrice    float64        `db:"unit_price"`
	UoMID        int            `db:"uom_id"`
	StockStatus  string         `db:"stock_status"`
	StorageID    int            `db:"storage_id"`
	CreatedBy    int            `db:"created_by"`
	DateCreated  time.Time      `db:"date_created"`
	DateModified sql.NullTime   `db:"date_modified"`
}

func (i *Item) UpdateQuantity(transactionType string, quantity int) {
	if transactionType == "inbound" {
		i.Quantity += quantity
	}

	if transactionType == "outbound" {
		i.Quantity -= quantity
	}
}

func (i *Item) UpdateCancelledQuantity(transactionType string, quantity int) {
	if transactionType == "inbound" {
		i.Quantity -= quantity
	}

	if transactionType == "outbound" {
		i.Quantity += quantity
	}
}
