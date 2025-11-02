package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

func GetOrderlineByID(id int) (schema.Orderline, error) {
	return RetrieveItemByField[schema.Orderline](OrderlineTable, "id", id)
}

func NewOrderline(transactionType string, orderline schema.Orderline) (int64, error) {
	var fields []string

	if transactionType == "inbound" {
		fields = []string{
			"transaction_id",
			"item_id",
			"quantity",
			"note",
			"is_voided",
			"created_by",
		}
	}

	if transactionType == "outbound" {
		fields = []string{
			"transaction_id",
			"item_id",
			"quantity",
			"unit_price",
			"total_amount",
			"note",
			"is_voided",
			"created_by",
		}
	}

	return InsertRecord(OrderlineTable, orderline, fields...)
}

func CancelOrderline(orderline schema.Orderline) error {
	return UpdateRecordByID(OrderlineTable, orderline, "is_voided", "updated_by")
}

func AddOrderlineNote(orderline schema.Orderline) error {
	return UpdateRecordByID(OrderlineTable, orderline, "note", "updated_by")
}
