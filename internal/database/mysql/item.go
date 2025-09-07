package mysql

import "github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"

func ListItem() ([]schema.Item, error) { return FetchItems[schema.Item](ItemTable) }

func GetItemByID(id int) (schema.Item, error) {
	return RetrieveItemByField[schema.Item](ItemTable, "id", id)
}

func GetItemByName(name string) (schema.Item, error) {
	return RetrieveItemByField[schema.Item](ItemTable, "name", name, "LOWER(?)")
}

func NewItem(item schema.Item) (int64, error) {
	fields := []string{
		"name",
		"description",
		"quantity",
		"unit_price",
		"uom_id",
		"stock_status",
		"storage_id",
		"created_by",
	}

	return InsertRecord(ItemTable, item, fields...)
}

func UpdateItem(item schema.Item) error {
	fields := []string{
		"name",
		"description",
		"quantity",
		"quantity",
		"unit_price",
		"storage_id",
		"uom_id",
	}

	return UpdateRecordByID(ItemTable, item, fields...)
}

func DeleteItem(id int) (int64, error) { return DeleteRecordByID(ItemTable, id) }

func ItemIDExists(id int) (bool, error) {
	return exists(func() (schema.Item, error) { return GetItemByID(id) })
}

func ItemNameExists(name string) (bool, error) {
	return exists(func() (schema.Item, error) { return GetItemByName(name) })
}
