package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

func ListUOM() ([]schema.UOM, error) {
	return fetch[schema.UOM]("SELECT * FROM unit_of_measurement;")
}

func GetUOM(id int) ([]schema.UOM, error) {
	return fetch[schema.UOM]("SELECT * FROM unit_of_measurement WHERE id = ?;", id)
}

func NewUOM(uom schema.UOM) error {
	query := `INSERT INTO unit_of_measurement (code, name)
						VALUES (:code, :name)`

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
	return delete("DELETE FROM unit_of_measurement WHERE id = ?;", id)
}

func UOMIDExists(id int) (bool, error) {
	return exists[schema.UOM]("SELECT * FROM unit_of_measurement WHERE id = ?", id)
}

func UOMNameExists(name string) (bool, error) {
	return exists[schema.UOM]("SELECT * FROM unit_of_measurement WHERE name = LOWER(?);", name)
}
