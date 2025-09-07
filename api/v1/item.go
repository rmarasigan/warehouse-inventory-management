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
	dbutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/db_utils"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
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
		log.Error(err.Error())
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

	ok, errors := validator.ValidateItem(body)
	if !ok {
		log.Error(strings.Join(errors, ", "), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: strings.Join(errors, ", ")})

		return
	}

	data, err := apischema.NewItem(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})

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
		existing, err := mysql.ItemNameExists(item.Name)
		if err != nil {
			log.Error(err.Error(), slog.Any("item", item), slog.Any("request", items), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to validate if item name exists", Details: item})

			return
		}

		if !existing {
			_, err = mysql.NewItem(item)
			if err != nil {
				log.Error(err.Error(), slog.Any("item", item), slog.Any("request", items))
				response.InternalServer(w, response.Response{Error: "failed to create new item", Details: item})

				return
			}
		}
	}

	response.Created(w, nil)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
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

	data, err := apischema.NewItem(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})

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
		itemID := item.ID

		existing, err := mysql.ItemIDExists(itemID)
		if err != nil {
			log.Error(err.Error(), slog.Any("item", item), slog.Any("request", items), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to validate if item id exists", Details: item})

			return
		}

		if existing {
			err = mysql.UpdateItem(item)
			if err != nil {
				log.Error(err.Error(), slog.Any("id", itemID), slog.Any("item", item), slog.Any("request", items), slog.Any("path", r.URL.Path))
				response.InternalServer(w, response.Response{Error: "failed to update item"})

				return
			}
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
		log.Error(err.Error(), slog.Any("id", id), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to delete item"})

		return
	}

	response.Success(w, response.Response{Message: fmt.Sprintf("%d rows(s) affected", affected)})
}
