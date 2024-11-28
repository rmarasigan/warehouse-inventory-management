package v1

import (
	"database/sql"
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

func storageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getStorages(w, r)

	case http.MethodPost:
		createStorage(w, r)

	case http.MethodPut:
		updateStorage(w, r)

	case http.MethodDelete:
		deleteStorage(w, r)
	}
}

func getStorages(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	storages, err := getList(r, mysql.GetStorage, mysql.StorageList)
	if err != nil {
		log.Error(err.Error())
		response.InternalServer(w, nil)

		return
	}

	response.Success(w, storages)
}

func createStorage(w http.ResponseWriter, r *http.Request) {
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

	ok, errors := validator.ValidateStorage(body)
	if !ok {
		log.Error(strings.Join(errors, ", "), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: strings.Join(errors, ", ")})

		return
	}

	data, err := apischema.NewStorage(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})

		return
	}

	var storages = convert.Schema(data, func(storage apischema.Storage) schema.Storage {
		return schema.Storage{
			ID:          storage.ID,
			Code:        storage.Code,
			Name:        storage.Name,
			Description: sql.NullString{String: storage.Description, Valid: true},
		}
	})

	for _, storage := range storages {
		existing, err := mysql.StorageNameExists(storage.Name)
		if err != nil {
			log.Error(err.Error(), slog.Any("storage", storage), slog.Any("request", storages), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to validate if storage name exists", Details: storage})

			return
		}

		if !existing {
			err = mysql.NewStorage(storage)
			if err != nil {
				log.Error(err.Error(), slog.Any("storage", storage), slog.Any("request", storages))
				response.InternalServer(w, response.Response{Error: "failed to create a new storage", Details: storage})

				return
			}
		}
	}

	response.Created(w, nil)
}

func updateStorage(w http.ResponseWriter, r *http.Request) {
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

	storages, err := apischema.NewStorage(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("any", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})

		return
	}

	for _, storage := range storages {
		storageID := storage.ID

		results, err := mysql.GetStorage(storageID)
		if err != nil {
			log.Error(err.Error(), slog.Any("id", storageID), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to retrieve storage"})

			return
		}

		dbstorages := convert.Schema(results, func(result apischema.Storage) schema.Storage {
			return schema.Storage{
				ID:          result.ID,
				Code:        result.Code,
				Name:        result.Name,
				Description: sql.NullString{String: result.Description, Valid: true},
			}
		})

		for _, dbstorage := range dbstorages {
			dbstorage.UpdateValues(storage)

			err := mysql.UpdateStorage(dbstorage)
			if err != nil {
				log.Error(err.Error(), slog.Any("id", storageID), slog.Any("path", r.URL.Path))
				response.InternalServer(w, response.Response{Error: "failed to update storage"})

				return
			}
		}
	}

	response.Success(w, nil)
}

func deleteStorage(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	affected, err := mysql.DeleteStorage(id)
	if err != nil {
		log.Error(err.Error(), slog.Any("id", id), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to delete storage"})

		return
	}

	response.Success(w, response.Response{Message: fmt.Sprintf("%d row(s) affected", affected)})
}
