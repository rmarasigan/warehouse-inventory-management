package v1

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err, "failed to read request body", log.KV("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to read request body"})

		return
	}

	if len(body) == 0 {
		errMsg := errors.New("request body cannot be empty")
		log.Error(errMsg, "missing request body", log.KV("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: errMsg.Error()})

		return
	}

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	ok, validationErrors := validator.ValidateNote(body)
	if !ok {
		errMsg := errors.New("invalid request body")
		log.Error(errMsg, strings.Join(validationErrors, ", "),
			log.KVs(map[string]any{
				"request": string(body),
				"path":    r.URL.Path,
			}),
		)
		response.BadRequest(w, response.Response{Error: errMsg.Error(), Details: strings.Join(validationErrors, ", ")})

		return
	}

	shared, err := apischema.NewNote(body)
	if err != nil {
		log.Error(err, "failed to unmarshal request body",
			log.KVs(map[string]any{
				"request": string(body),
				"path":    r.URL.Path,
			}),
		)
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})
	}

	orderline, err := mysql.GetOrderlineByID(id)
	if err != nil {
		log.Error(err, "failed to retrieve orderline",
			log.KVs(map[string]any{
				"id":      id,
				"path":    r.URL.Path,
				"request": string(body),
			}),
		)
		response.InternalServer(w,
			response.Response{
				Error:   err.Error(),
				Details: "failed to fetch orderline " + fmt.Sprint(id),
			},
		)

		return
	}

	orderline.Note = dbutils.SetString(shared.Note)
	orderline.UpdatedBy = dbutils.SetInt(shared.UserID)

	err = mysql.AddOrderlineNote(orderline)
	if err != nil {
		log.Error(err, "failed to add transaction note",
			log.KVs(map[string]any{
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
					"message":      "failed to add orderline note",
				},
			})
	}

	response.Success(w, nil)
}
