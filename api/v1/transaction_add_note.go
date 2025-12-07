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

func transactionAddNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		transactionNote(w, r)
	}
}

func transactionNote(w http.ResponseWriter, r *http.Request) {
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

	validationErrors, err := requestutils.ValidateRequest(body, validator.ValidateStorage)
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

	transaction := schema.Transaction{
		ID:        id,
		Note:      dbutils.SetString(shared.Note),
		UpdatedBy: dbutils.SetInt(shared.UserID),
	}

	err = mysql.UpdateTransactionNote(transaction)
	if err != nil {
		log.Error(err, "failed to add transaction note",
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
					"transaction_id": id,
					"note":           shared.Note,
					"message":        "failed to add transaction note",
				},
			})

		return
	}

	response.Success(w, nil)
}
