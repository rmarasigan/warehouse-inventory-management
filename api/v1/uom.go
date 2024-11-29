package v1

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func uomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUOMs(w, r)

	case http.MethodPost:
		createUOM(w, r)

	case http.MethodPut:
		updateUOM(w, r)

	case http.MethodDelete:
		deleteUOM(w, r)
	}
}

func getUOMs(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	uoms, err := getList(r, mysql.GetUOM, mysql.UOMList)
	if err != nil {
		log.Error(err.Error())
		response.InternalServer(w, nil)

		return
	}

	response.Success(w, uoms)
}

func createUOM(w http.ResponseWriter, r *http.Request) {
	defer func() {
		log.Panic()
		r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to read request body"})

		return
	}

	if len(body) == 0 {
		log.Error("missing request body", slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "request body cannot be empty"})

		return
	}

	ok, errors := validator.ValidateUOM(body)
	if !ok {
		log.Error(strings.Join(errors, ", "), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: strings.Join(errors, ", ")})

		return
	}

	data, err := apischema.NewUOM(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})

		return
	}

	uoms := convert.Schema(data, func(uom apischema.UOM) schema.UOM {
		return schema.UOM{
			ID:   uom.ID,
			Code: uom.Code,
			Name: uom.Name,
		}
	})

	for _, uom := range uoms {
		existing, err := mysql.UOMNameExists(uom.Name)
		if err != nil {
			log.Error(err.Error(), slog.Any("uom", uom), slog.Any("request", uoms), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to validate if role name exists", Details: uom})

			return
		}

		if !existing {
			err = mysql.NewUOM(uom)
			if err != nil {
				log.Error(err.Error(), slog.Any("uom", uom), slog.Any("request", uoms))
				response.InternalServer(w, response.Response{Error: "failed to create new uom", Details: uom})

				return
			}
		}
	}

	response.Created(w, nil)
}

func updateUOM(w http.ResponseWriter, r *http.Request) {
	defer func() {
		log.Panic()
		r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to read request body"})

		return
	}

	data, err := apischema.NewUOM(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "failed to unmarshal request body"})

		return
	}

	for _, uom := range data {
		uomID := uom.ID

		dbuoms, err := mysql.GetUOM(uomID)
		if err != nil {
			log.Error(err.Error(), slog.Any("uom", uom), slog.Any("request", uoms), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to retrieve uom"})

			return
		}

		for _, dbuom := range dbuoms {
			dbuom.UpdateValues(uom)

			err = mysql.UpdateUOM(dbuom)
			if err != nil {
				log.Error(err.Error(), slog.Any("id", uomID), slog.Any("uom", uom), slog.Any("request", uoms), slog.Any("path", r.URL.Path))
				response.InternalServer(w, response.Response{Error: "failed to update uom"})

				return
			}
		}
	}
}

func deleteUOM(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	affected, err := mysql.DeleteUOM(id)
	if err != nil {
		log.Error(err.Error(), slog.Any("id", id), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to delete uom"})

		return
	}

	response.Success(w, response.Response{Message: fmt.Sprintf("%d row(s) affected", affected)})
}
