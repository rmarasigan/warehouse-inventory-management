package v1

import (
	"net/http"

	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	dbutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/db_utils"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func orderlineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		orderlineNote(w, r)
	}
}

func orderlineNote(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	updateNote(w, r,
		func(id int, shared apischema.Shared) any {
			return schema.Orderline{
				ID:        id,
				Note:      dbutils.SetString(shared.Note),
				UpdatedBy: dbutils.SetInt(shared.UserID),
			}
		},
		func(record any) error {
			return mysql.UpdateOrderlineNote(record.(schema.Orderline))
		},
	)
}
