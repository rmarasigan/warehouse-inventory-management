package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

func ListStorage() ([]schema.Storage, error) {
	return fetch[schema.Storage]("SELECT * FROM storage;")
}

func GetStorage(id int) ([]schema.Storage, error) {
	return fetch[schema.Storage]("SELECT * FROM storage WHERE id = ?;", id)
}

func NewStorage(storage schema.Storage) error {
	query := `INSERT INTO storage (code, name, description)
						VALUES (:code, :name, :description);`

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
	return delete("DELETE FROM storage WHERE id = ?;", id)
}

func StorageIDExists(id int) (bool, error) {
	return entityExists(GetStorage, id)
}

func StorageNameExists(name string) (bool, error) {
	return exists[schema.Storage]("SELECT * FROM storage WHERE name = LOWER(?);", name)
}
