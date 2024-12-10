package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

func ListStorage() ([]schema.Storage, error) {
	query := "SELECT * FROM storage;"
	return fetch[schema.Storage](query)
}

func GetStorage(id int) ([]schema.Storage, error) {
	query := "SELECT * FROM storage WHERE id = ?;"
	return fetch[schema.Storage](query, id)
}

func NewStorage(storage schema.Storage) error {
	query := "INSERT INTO storage (id, code, name, description) VALUES (:id, :code, :name, :description);"

	_, err := NamedExec(query, storage)

	return err
}

func UpdateStorage(storage schema.Storage) error {
	query := `UPDATE storage
						SET
							code = CASE
								WHEN :code = '' THEN code
								ELSE COALESCE(:code, code)
							END,
							name = CASE
								WHEN :name = '' THEN name
								ELSE COALESCE(:name, name)
							END,
							description = CASE
								WHEN :description = '' THEN description
								ELSE COALESCE(:description, description)
							END
						WHERE id = :id;`

	_, err := NamedExec(query, storage)

	return err
}

func DeleteStorage(id int) (int64, error) {
	query := "DELETE FROM storage WHERE id = ?;"
	return delete(query, id)
}

func StorageIDExists(id int) (bool, error) {
	storages, err := GetStorage(id)
	if err != nil {
		return false, err
	}

	return (len(storages) > 0), nil
}

func StorageNameExists(name string) (bool, error) {
	query := "SELECT * FROM storage WHERE name = LOWER(?);"
	return exists[schema.Storage](query, name)
}
