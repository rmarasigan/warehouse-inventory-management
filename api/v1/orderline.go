package v1

import (
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	dbutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/db_utils"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	requestutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/request_utils"
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

	body, err := requestutils.ReadBody(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	validationErrors, err := requestutils.ValidateRequest(body, validator.ValidateItem)
	if err != nil && len(validationErrors) > 0 {
		log.Error(err, validationErrors, log.KVs(log.Map{"request": string(body), "path": r.URL.Path}))
		response.BadRequest(w, response.Response{Error: err.Error(), Details: validationErrors})

		return
	}

	shared, err := apischema.NewNote(body)
	if err != nil {
		log.Error(err, "failed to unmarshal request body",
			log.KVs(log.Map{"request": string(body), "path": r.URL.Path}))

		response.BadRequest(w, response.Response{Error: "failed to unmarshal request body"})
	}

	orderline := schema.Orderline{
		ID:        id,
		Note:      dbutils.SetString(shared.Note),
		UpdatedBy: dbutils.SetInt(shared.UserID),
	}

	err = mysql.UpdateOrderlineNote(orderline)
	if err != nil {
		log.Error(err, "failed to update orderline note",
			log.KVs(log.Map{
				"id":      id,
				"path":    r.URL.Path,
				"request": string(body),
			}),
		)

		response.InternalServer(w,
			response.Response{
				Error: err.Error(),
				Details: map[string]any{
					"orderline_id": id,
					"note":         shared.Note,
					"message":      "failed to update orderline note",
				},
			})
	}

	response.Success(w, nil)
}
