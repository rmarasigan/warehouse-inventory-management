package schema

import (
	"database/sql"
	"time"
)

type (
	Transaction struct {
		ID           int             `db:"id"`
		Reference    string          `db:"reference"`
		Orderlines   []Orderline     `db:"-"`
		Amount       sql.NullFloat64 `db:"amount"`
		Type         string          `db:"type"`
		IsCancelled  sql.NullBool    `db:"is_cancelled"`
		Note         sql.NullString  `db:"note"`
		CreatedBy    int             `db:"created_by"`
		UpdatedBy    sql.NullInt32   `db:"updated_by"`
		DateCreated  time.Time       `db:"date_created"`
		DateModified sql.NullTime    `db:"date_modified"`
	}

	Orderline struct {
		ID            int             `db:"id"`
		TransactionID int             `db:"transaction_id"`
		ItemID        int             `db:"item_id"`
		Quantity      int             `db:"quantity"`
		UnitPrice     sql.NullFloat64 `db:"unit_price"`
		TotalAmount   sql.NullFloat64 `db:"total_amount"`
		Note          sql.NullString  `db:"note"`
		IsVoided      sql.NullBool    `db:"is_voided"`
		CreatedBy     int             `db:"created_by"`
		UpdatedBy     sql.NullInt32   `db:"updated_by"`
		DateCreated   time.Time       `db:"date_created"`
		DateModified  sql.NullTime    `db:"date_modified"`
	}
)

func (t *Transaction) IsValidTransactionType() bool {
	return (t.Type == "inbound") || (t.Type == "outbound")
}
