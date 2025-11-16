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
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	dbutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/db_utils"
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

	list, err := getList(r, mysql.GetStorageByID, mysql.ListStorage)
	if err != nil {
		log.Error(err, "failed to retrieve storages", log.KV("path", r.URL.Path))
		response.InternalServer(w, nil)

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

	ok, validationErrors := validator.ValidateStorage(body)
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

	data, err := apischema.NewStorage(body)
	if err != nil {
		log.Error(err, "failed to unmarshal request body",
			log.KVs(map[string]any{
				"request": string(body),
				"path":    r.URL.Path,
			}),
		)
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})

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
		existing, err := mysql.StorageNameExists(storage.Name)
		if err != nil {
			log.Error(err, "failed to validate if storage name exists",
				log.KVs(map[string]any{
					"request": string(body),
					"storage": storage,
					"path":    r.URL.Path,
				}),
			)
			response.InternalServer(w,
				response.Response{
					Error: err.Error(),
					Details: map[string]any{
						"storage": storage,
						"message": "failed to validate if storage name exists",
					},
				})

			return
		}

		if !existing {
			_, err = mysql.NewStorage(storage)
			if err != nil {
				log.Error(err, "failed to create storage",
					log.KVs(map[string]any{
						"storage": storage,
						"path":    r.URL.Path,
					}),
				)
				response.InternalServer(w,
					response.Response{
						Error: err.Error(),
						Details: map[string]any{
							"request": data,
							"storage": storage,
							"message": "failed to create storage",
						},
					},
				)

				return
			}
		}
	}

	response.Created(w, nil)
}

func updateStorage(w http.ResponseWriter, r *http.Request) {
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

	data, err := apischema.NewStorage(body)
	if err != nil {
		log.Error(err, "failed to unmarshal request body",
			log.KVs(map[string]any{
				"request": string(body),
				"path":    r.URL.Path,
			}),
		)
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})

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
		storageID := storage.ID

		existing, err := mysql.StorageIDExists(storageID)
		if err != nil {
			log.Error(err, "failed to validate if storage id exists",
				log.KVs(map[string]any{
					"request": string(body),
					"storage": storage,
					"path":    r.URL.Path,
				}),
			)
			response.InternalServer(w,
				response.Response{
					Error: err.Error(),
					Details: map[string]any{
						"storage": storage,
						"message": "failed to validate if storage id exists",
					},
				})

			return
		}

		if existing {
			err := mysql.UpdateStorage(storage)
			if err != nil {
				log.Error(err, "failed to update storage",
					log.KVs(map[string]any{
						"request": data,
						"storage": storage,
						"path":    r.URL.Path,
					}),
				)
				response.InternalServer(w,
					response.Response{
						Error: err.Error(),
						Details: map[string]any{
							"request": data,
							"storage": storage,
							"message": "failed to update storage",
						},
					},
				)

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
		log.Error(err, "failed to delete storage",
			log.KVs(map[string]any{
				"id":   id,
				"path": r.URL.Path,
			}),
		)
		response.InternalServer(w, response.Response{Error: err.Error(), Details: "failed to delete storage"})

		return
	}

	response.Success(w, response.Response{Message: fmt.Sprintf("%d row(s) affected", affected)})
}
