package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

func ListTransaction() ([]schema.Transaction, error) {
	transactions, err := FetchItems[schema.Transaction](TransactionTable)
	if err != nil {
		return nil, err
	}

	// Use index-based iteration to directly access and modify the original element
	for i := range transactions {
		orderlines, err := GetOrderlineByTransactionID(transactions[i].ID)
		if err != nil {
			return nil, err
		}

		// Changes are reflected in the original element
		transactions[i].Orderlines = orderlines
	}

	return transactions, nil
}

func GetTransactionByID(id int) (schema.Transaction, error) {
	transaction, err := RetrieveItemByField[schema.Transaction](TransactionTable, "id", id)
	if err != nil {
		return schema.Transaction{}, err
	}

	orderlines, err := GetOrderlineByTransactionID(transaction.ID)
	if err != nil {
		return schema.Transaction{}, err
	}

	transaction.Orderlines = append(transaction.Orderlines, orderlines...)

	return transaction, nil
}

func GetOrderlineByTransactionID(id int) ([]schema.Orderline, error) {
	condition := map[string]any{"transaction_id": "?"}
	return FetchItemsByFields[schema.Orderline](OrderlineTable, condition, id)
}

func NewTransaction(_type string, transaction schema.Transaction) (int64, error) {
	var fields []string

	if _type == "inbound" {
		fields = []string{
			"reference",
			"type",
			"note",
			"created_by",
		}
	}

	if _type == "outbound" {
		fields = []string{
			"reference",
			"type",
			"amount",
			"note",
			"created_by",
		}
	}

	return InsertRecord(TransactionTable, transaction, fields...)
}

func CancelTransaction(transaction schema.Transaction) error {
	return UpdateRecordByID(TransactionTable, transaction, "is_cancelled", "updated_by")
}

func AddTransactionNote(transaction schema.Transaction) error {
	return UpdateRecordByID(TransactionTable, transaction, "note", "updated_by")
}
