package mysql

import "github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"

func ListUOM() ([]schema.UOM, error) { return FetchItems[schema.UOM](UoMTable) }

func GetUOMByID(id int) (schema.UOM, error) {
	return RetrieveItemByField[schema.UOM](UoMTable, "id", id)
}

func GetUOMByName(name string) (schema.UOM, error) {
	return RetrieveItemByField[schema.UOM](UoMTable, "name", name, "LOWER(?)")
}

func NewUOM(uom schema.UOM) (int64, error) { return InsertRecord(UoMTable, uom, "code", "name") }

func UpdateUOM(uom schema.UOM) error { return UpdateRecordByID(UoMTable, uom, "code", "name") }

func DeleteUOM(id int) (int64, error) { return DeleteRecordByID(UoMTable, id) }

func UOMIDExists(id int) (bool, error) {
	return exists(func() (schema.UOM, error) { return GetUOMByID(id) })
}

func UOMNameExists(name string) (bool, error) {
	return exists(func() (schema.UOM, error) { return GetUOMByName(name) })
}
