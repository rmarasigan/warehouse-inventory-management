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

func itemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItems(w, r)

	case http.MethodPost:
		createItem(w, r)

	case http.MethodPut:
		updateItem(w, r)

	case http.MethodDelete:
		deleteItem(w, r)
	}
}

func getItems(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	list, err := getList(r, mysql.GetItemByID, mysql.ListItem)
	if err != nil {
		log.Error(err, "failed to retrieve items", log.KV("path", r.URL.Path))
		response.InternalServer(w, nil)

		return
	}

	items := convert.SchemaList(list, func(item schema.Item) apischema.Item {
		return apischema.Item{
			ID:           item.ID,
			Name:         item.Name,
			Description:  dbutils.GetString(item.Description),
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			UoMID:        item.UoMID,
			StockStatus:  item.StockStatus,
			StorageID:    item.StorageID,
			CreatedBy:    item.CreatedBy,
			DateCreated:  item.DateCreated,
			DateModified: dbutils.GetTime(item.DateModified),
		}
	})

	response.Success(w, items)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
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

	data, err := requestutils.Unmarshal(r.URL.Path, body, apischema.NewItem)
	if err != nil {
		response.BadRequest(w, response.Response{Error: "failed to unmarshal request body"})
		return
	}

	items := convert.SchemaList(data, func(item apischema.Item) schema.Item {
		return schema.Item{
			Name:        item.Name,
			Description: dbutils.SetString(item.Description),
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			UoMID:       item.UoMID,
			StockStatus: item.StockStatus,
			StorageID:   item.StorageID,
			CreatedBy:   item.CreatedBy,
		}
	})

	for _, item := range items {
		_, err := mysql.NewItemIfNotExists(item)
		if err != nil {
			log.Error(err, "failed to create item", log.KVs(log.Map{"item": item, "path": r.URL.Path}))
			response.InternalServer(w,
				response.Response{
					Error: err.Error(),
					Details: map[string]any{
						"request": data,
						"item":    item,
						"message": "failed to create item",
					},
				},
			)

			return
		}
	}

	response.Created(w, nil)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	data, err := requestutils.Unmarshal(r.URL.Path, body, apischema.NewItem)
	if err != nil {
		response.BadRequest(w, response.Response{Error: "failed to unmarshal request body"})
		return
	}

	items := convert.SchemaList(data, func(item apischema.Item) schema.Item {
		return schema.Item{
			ID:          item.ID,
			Name:        item.Name,
			Description: dbutils.SetString(item.Description),
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			StorageID:   item.StorageID,
			UoMID:       item.UoMID,
		}
	})

	for _, item := range items {
		err = mysql.UpdateItem(item)
		if err != nil {
			log.Error(err, "failed to update item", log.KVs(log.Map{"request": data, "item": item, "path": r.URL.Path}))

			response.InternalServer(w,
				response.Response{
					Error: err.Error(),
					Details: map[string]any{
						"request": data,
						"item":    item,
						"message": "failed to update item",
					},
				},
			)

			return
		}
	}

	response.Success(w, nil)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	affected, err := mysql.DeleteItem(id)
	if err != nil {
		log.Error(err, "failed to delete item", log.KVs(log.Map{"id": id, "path": r.URL.Path}))
		response.InternalServer(w, response.Response{Error: err.Error(), Details: "failed to delete item"})

		return
	}

	response.Success(w, response.Response{Message: fmt.Sprintf("%d rows(s) affected", affected)})
}
