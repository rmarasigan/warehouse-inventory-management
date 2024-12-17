package mysql

import (
	"fmt"

	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

func ListUOM() ([]schema.UOM, error) { return FetchItems[schema.UOM](UoMTable) }

func GetUOMByID(id int) (schema.UOM, error) {
	return RetrieveItemByField[schema.UOM](UoMTable, "id", id)
}

func GetUOMByName(name string) (schema.UOM, error) {
	return RetrieveItemByField[schema.UOM](UoMTable, "name", name, "LOWER(?)")
}

func NewUOM(uom schema.UOM) error {
	return InsertRecord(fmt.Sprintf("INSERT INTO %s (code, name) VALUES (:code, :name)", UoMTable), uom)
}

func UpdateUOM(uom schema.UOM) error {
	return UpdateRecordByID(UoMTable, uom, []string{"code", "name"})
}

func DeleteUOM(id int) (int64, error) { return DeleteRecordByID(UoMTable, id) }

func UOMIDExists(id int) (bool, error) {
	return exists(func() (schema.UOM, error) { return GetUOMByID(id) })
}

func UOMNameExists(name string) (bool, error) {
	return exists(func() (schema.UOM, error) { return GetUOMByName(name) })
}
