package mysql

import (
	"strings"

	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
)

func StorageList() ([]apischema.Storage, error) {
	var (
		list  []schema.Storage
		query = "SELECT * FROM storage;"
	)

	err := Select(&list, query)
	if err != nil {
		return nil, err
	}

	var storages = convert.Schema(list, func(storage schema.Storage) apischema.Storage {
		return apischema.Storage{
			ID:          storage.ID,
			Code:        storage.Code,
			Name:        storage.Name,
			Description: storage.Description.String,
		}
	})

	return storages, nil
}

func GetStorage(id int) ([]apischema.Storage, error) {
	var (
		list  []schema.Storage
		query = `SELECT * FROM storage WHERE id = ?;`
	)

	err := Select(&list, query, id)
	if err != nil {
		return nil, err
	}

	var storages = convert.Schema(list, func(storage schema.Storage) apischema.Storage {
		return apischema.Storage{
			ID:          storage.ID,
			Code:        storage.Code,
			Name:        storage.Name,
			Description: storage.Description.String,
		}
	})

	return storages, nil
}

func NewStorage(storage schema.Storage) error {
	query := `INSERT INTO storage (id, code, name, description) VALUES (:id, :code, :name, :description)`

	_, err := NamedExec(query, storage)

	return err
}

func UpdateStorage(storage schema.Storage) error {
	query := `UPDATE storage
						SET code = :code,
						name = :name,
						description = :description
						WHERE id = :id;`

	_, err := NamedExec(query, storage)

	return err
}

func DeleteStorage(id int) (int64, error) {
	query := `DELETE FROM storage WHERE id = ?;`

	result, err := Exec(query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func StorageIDExists(id int) (bool, error) {
	storages, err := GetStorage(id)
	if err != nil {
		return false, err
	}

	return (len(storages) > 0), nil
}

func StorageNameExists(name string) (bool, error) {
	var (
		list  []schema.Storage
		query = `SELECT * FROM storage WHERE name = LOWER(?);`
	)

	err := Select(&list, query, strings.ToLower(name))
	if err != nil {
		return false, err
	}

	return (len(list) > 0), nil
}
