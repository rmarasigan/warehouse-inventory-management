package apischema

import (
	"time"

	"github.com/google/uuid"
)

type (
	Transaction struct {
		ID           int         `json:"id"`
		Reference    string      `json:"reference"`
		Orderlines   []Orderline `json:"orderlines"`
		Amount       float64     `json:"amount"`
		Type         string      `json:"type,omitempty"`
		IsCancelled  bool        `json:"is_cancelled"`
		Note         string      `json:"note,omitempty"`
		CreatedBy    int         `json:"created_by"`
		UpdatedBy    int         `json:"updated_by,omitempty"`
		DateCreated  time.Time   `json:"date_created"`
		DateModified time.Time   `json:"date_modified,omitzero"`
	}

	Orderline struct {
		ID            int       `json:"id"`
		TransactionID int       `json:"transaction_id"`
		ItemID        int       `json:"item_id"`
		Quantity      int       `json:"quantity"`
		UnitPrice     float64   `json:"unit_price"`
		TotalAmount   float64   `json:"total_amount"`
		Note          string    `json:"note,omitempty"`
		IsVoided      bool      `json:"is_voided"`
		CreatedBy     int       `json:"created_by"`
		UpdatedBy     int       `json:"updated_by,omitempty"`
		DateCreated   time.Time `json:"date_created"`
		DateModified  time.Time `json:"date_modified,omitzero"`
	}
)

func NewTransaction(data []byte) (Transaction, error) {
	transactions, err := unmarshal[Transaction](data)
	if len(transactions) == 1 {
		return transactions[0], err
	}

	return Transaction{}, err
}

func (t Transaction) GenerateReference() string {
	reference, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}

	return reference.String()
}
