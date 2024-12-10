package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

func ListUOM() ([]schema.UOM, error) {
	query := "SELECT * FROM unit_of_measurement;"
	return fetch[schema.UOM](query)
}

func GetUOM(id int) ([]schema.UOM, error) {
	query := "SELECT * FROM unit_of_measurement WHERE id = ?;"
	return fetch[schema.UOM](query, id)
}

func NewUOM(uom schema.UOM) error {
	query := `INSERT INTO unit_of_measurement (id, code, name) VALUES (:id, :code, :name)`

	_, err := NamedExec(query, uom)

	return err
}

func UpdateUOM(uom schema.UOM) error {
	query := `UPDATE unit_of_measurement
						SET
							code = CASE
								WHEN :code = '' THEN code
								ELSE COALESCE(:code, code)
							END,
							name = CASE
								WHEN :name = '' THEN name
								ELSE COALESCE(:name, name)
							END
						WHERE id = :id;`

	_, err := NamedExec(query, uom)

	return err
}

func DeleteUOM(id int) (int64, error) {
	query := `DELETE FROM unit_of_measurement WHERE id = ?;`
	return delete(query, id)
}

func UOMIDExists(id int) (bool, error) {
	query := `SELECT * FROM unit_of_measurement WHERE id = ?`
	return exists[schema.UOM](query, id)
}

func UOMNameExists(name string) (bool, error) {
	query := `SELECT * FROM unit_of_measurement WHERE name = LOWER(?);`
	return exists[schema.UOM](query, name)
}
