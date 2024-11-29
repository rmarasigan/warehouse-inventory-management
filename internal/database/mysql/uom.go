package mysql

import (
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

func UOMList() ([]schema.UOM, error) {
	var (
		list  []schema.UOM
		query = "SELECT * FROM unit_of_measurement;"
	)

	err := Select(&list, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func GetUOM(id int) ([]schema.UOM, error) {
	var (
		list  []schema.UOM
		query = "SELECT * FROM unit_of_measurement WHERE id = ?;"
	)

	err := Select(&list, query, id)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func NewUOM(uom schema.UOM) error {
	query := `INSERT INTO unit_of_measurement (id, code, name) VALUES (:id, :code, :name)`

	_, err := NamedExec(query, uom)

	return err
}

func UpdateUOM(uom schema.UOM) error {
	query := `UPDATE unit_of_measurement
						SET code = :code,
						name = :name
						WHERE id = :id;`

	_, err := NamedExec(query, uom)

	return err
}

func DeleteUOM(id int) (int64, error) {
	query := `DELETE FROM unit_of_measurement WHERE id = ?;`

	result, err := Exec(query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func UOMNameExists(name string) (bool, error) {
	var (
		list  []schema.UOM
		query = `SELECT * FROM unit_of_measurement WHERE name = LOWER(?);`
	)

	err := Select(&list, query, strings.ToLower(name))
	if err != nil {
		return false, err
	}

	return (len(list) > 0), nil
}
