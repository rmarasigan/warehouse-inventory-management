package mysql

import "github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"

func ListStorage() ([]schema.Storage, error) { return FetchItems[schema.Storage](StorageTable) }

func GetStorageByID(id int) (schema.Storage, error) {
	return RetrieveItemByField[schema.Storage](StorageTable, "id", id)
}

func GetStorageByName(name string) (schema.Storage, error) {
	return RetrieveItemByField[schema.Storage](StorageTable, "name", name, "LOWER(?)")
}

func NewStorage(storage schema.Storage) (int64, error) {
	return InsertRecord(StorageTable, storage, "code", "name", "description")
}

func UpdateStorage(storage schema.Storage) error {
	return UpdateRecordByID(StorageTable, storage, "code", "name", "description")
}

func DeleteStorage(id int) (int64, error) { return DeleteRecordByID(StorageTable, id) }

func StorageIDExists(id int) (bool, error) {
	return exists(func() (schema.Storage, error) { return GetStorageByID(id) })
}

func StorageNameExists(name string) (bool, error) {
	return exists(func() (schema.Storage, error) { return GetStorageByName(name) })
}
