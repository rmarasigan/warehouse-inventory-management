package v1

import (
	"fmt"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	dbutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/db_utils"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	requestutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/request_utils"
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

	list, err := getList(r, mysql.GetStorageByID, mysql.ListStorage)
	if err != nil {
		log.Error(err, "failed to retrieve storages", log.KV("path", r.URL.Path))
		response.InternalServer(w, response.NewError(err, "failed to retrieve storages"))

		return
	}

	storages := convert.SchemaList(list, func(storage schema.Storage) apischema.Storage {
		return apischema.Storage{
			ID:          storage.ID,
			Code:        storage.Code,
			Name:        storage.Name,
			Description: dbutils.GetString(storage.Description),
		}
	})

	response.Success(w, storages)
}

func createStorage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
	if err != nil {
		response.BadRequest(w, response.NewError(err))
		return
	}

	validationErrors, err := requestutils.ValidateRequest(body, validator.ValidateStorage)
	if err != nil && len(validationErrors) > 0 {
		log.Error(err, validationErrors, log.KVs(log.Map{"request": string(body), "path": r.URL.Path}))
		response.BadRequest(w, response.NewError(err, validationErrors))

		return
	}

	data, err := requestutils.Unmarshal(r.URL.Path, body, apischema.NewStorage)
	if err != nil {
		response.BadRequest(w, response.NewError(err, "failed to unmarshal request body"))
		return
	}

	storages := convert.SchemaList(data, func(storage apischema.Storage) schema.Storage {
		return schema.Storage{
			ID:          storage.ID,
			Code:        storage.Code,
			Name:        storage.Name,
			Description: dbutils.SetString(storage.Description),
		}
	})

	for _, storage := range storages {
		_, err = mysql.NewStorageIfNotExists(storage)
		if err != nil {
			log.Error(err, "failed to create storage",
				log.KVs(log.Map{"storage": storage, "path": r.URL.Path}))

			response.InternalServer(w, response.NewError(err,
				map[string]any{
					"request": data,
					"storage": storage,
					"message": "failed to create storage",
				}),
			)

			return
		}
	}

	response.Created(w, nil)
}

func updateStorage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
	if err != nil {
		response.BadRequest(w, response.NewError(err))
		return
	}

	data, err := requestutils.Unmarshal(r.URL.Path, body, apischema.NewStorage)
	if err != nil {
		response.BadRequest(w, response.NewError(err, "failed to unmarshal request body"))
		return
	}

	storages := convert.SchemaList(data, func(storage apischema.Storage) schema.Storage {
		return schema.Storage{
			ID:          storage.ID,
			Code:        storage.Code,
			Name:        storage.Name,
			Description: dbutils.SetString(storage.Description),
		}
	})

	for _, storage := range storages {
		err := mysql.UpdateStorage(storage)
		if err != nil {
			log.Error(err, "failed to update storage",
				log.KVs(log.Map{"request": data, "storage": storage, "path": r.URL.Path}))

			response.InternalServer(w, response.NewError(err,
				map[string]any{
					"request": data,
					"storage": storage,
					"message": "failed to update storage",
				}),
			)

			return
		}
	}

	response.Success(w, nil)
}

func deleteStorage(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.NewError(err))
		return
	}

	affected, err := mysql.DeleteStorage(id)
	if err != nil {
		log.Error(err, "failed to delete storage", log.KVs(log.Map{"id": id, "path": r.URL.Path}))
		response.InternalServer(w, response.NewError(err, "failed to delete storage"))

		return
	}

	response.Success(w, response.New(fmt.Sprintf("%d row(s) affected", affected)))
}
