package v1

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
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
		log.Panic()
		r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to read request body"})
	}

	if len(body) == 0 {
		log.Error("missing request body", slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "request body cannot be empty"})

		return
	}

	idParam, ok := requestutils.HasQueryParam(r, "id")
	if !ok {
		err := "missing 'id' from request query"
		log.Error(err, slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: err})

		return
	}

	orderlineID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "the 'id' is invalid; must be an integer"})

		return
	}

	ok, errors := validator.ValidateNote(body)
	if !ok {
		log.Error(strings.Join(errors, ", "), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: strings.Join(errors, ", ")})

		return
	}

	orderline, err := mysql.GetOrderlineByID(orderlineID)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to fetch orderline " + idParam})

		return
	}

	var shared apischema.Shared
	err = json.Unmarshal(body, &shared)
	if err != nil {
		log.Error("failed to unmarshal request", slog.Any("error", err.Error()), slog.Int("id", orderlineID), slog.Any("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request"})

		return
	}

	orderline.Note = dbutils.SetString(shared.Note)
	orderline.UpdatedBy = dbutils.SetInt(shared.UserID)

	err = mysql.AddOrderlineNote(orderline)
	if err != nil {
		log.Error("failed to add transaction note", slog.Any("error", err.Error()), slog.Int("id", orderlineID), slog.Any("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w,
			response.Response{
				Error: "failed to add orderline note",
				Details: map[string]any{
					"note":         shared.Note,
					"orderline_id": orderlineID,
				},
			})
	}

	response.Success(w, nil)
}
