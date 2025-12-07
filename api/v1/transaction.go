package v1

import (
	"errors"
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

func transactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTransactions(w, r)

	case http.MethodPost:
		createTransaction(w, r)
	}
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()
	list, err := getList(r, mysql.GetTransactionByID, mysql.ListTransaction)
	if err != nil {
		log.Error(err, "failed to retrieve transactions", log.KV("path", r.URL.Path))
		response.InternalServer(w, nil)

		return
	}

	transactions := convert.SchemaList(list,
		func(transaction schema.Transaction) apischema.Transaction {
			orderlines := convert.SchemaList(transaction.Orderlines,
				func(orderline schema.Orderline) apischema.Orderline {
					return apischema.Orderline{
						ID:            orderline.ID,
						TransactionID: orderline.TransactionID,
						ItemID:        orderline.ItemID,
						Quantity:      orderline.Quantity,
						UnitPrice:     dbutils.GetFloat(orderline.UnitPrice),
						TotalAmount:   dbutils.GetFloat(orderline.TotalAmount),
						Note:          dbutils.GetString(orderline.Note),
						IsVoided:      dbutils.GetBool(orderline.IsVoided),
						CreatedBy:     orderline.CreatedBy,
						UpdatedBy:     dbutils.GetAsInt(orderline.UpdatedBy),
						DateCreated:   orderline.DateCreated,
						DateModified:  dbutils.GetTime(orderline.DateModified),
					}
				},
			)

			return apischema.Transaction{
				ID:           transaction.ID,
				Reference:    transaction.Reference,
				Orderlines:   orderlines,
				Note:         dbutils.GetString(transaction.Note),
				CreatedBy:    transaction.CreatedBy,
				UpdatedBy:    dbutils.GetAsInt(transaction.UpdatedBy),
				DateCreated:  transaction.DateCreated,
				DateModified: dbutils.GetTime(transaction.DateModified),
			}
		},
	)

	response.Success(w, transactions)
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
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

	data, err := apischema.NewTransaction(body)
	if err != nil {
		log.Error(err, "failed to unmarshal request body",
			log.KVs(log.Map{"request": string(body), "path": r.URL.Path}))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})

		return
	}

	transaction := convert.Schema(data,
		func(trans apischema.Transaction) schema.Transaction {
			var amount float64

			orderlines := convert.SchemaList(data.Orderlines,
				func(orderline apischema.Orderline) schema.Orderline {
					amount += orderline.TotalAmount

					return schema.Orderline{
						ItemID:      orderline.ItemID,
						Quantity:    orderline.Quantity,
						UnitPrice:   dbutils.SetFloat(orderline.UnitPrice),
						TotalAmount: dbutils.SetFloat(orderline.TotalAmount),
						Note:        dbutils.SetString(orderline.Note),
						CreatedBy:   data.CreatedBy,
					}
				})

			return schema.Transaction{
				Reference:  data.GenerateReference(),
				Orderlines: orderlines,
				Amount:     dbutils.SetFloat(amount),
				Type:       data.Type,
				Note:       dbutils.SetString(data.Note),
				CreatedBy:  data.CreatedBy,
			}
		})

	if !transaction.IsValidTransactionType() {
		err := errors.New("transaction '" + transaction.Type + "' is not implemented")
		log.Error(err, "invalid transaction type",
			log.KVs(log.Map{"request": data, "path": r.URL.Path}))

		response.NotImplemented(w,
			response.Response{
				Error: err.Error(),
				Details: map[string]any{
					"request": data,
					"message": "invalid transaction type",
				},
			},
		)

		return
	}

	transactionType := transaction.Type
	lastInsertID, err := mysql.NewTransaction(transactionType, transaction)
	if err != nil {
		log.Error(err, "failed to create transaction",
			log.KVs(log.Map{"request": data, "path": r.URL.Path}))

		response.InternalServer(w,
			response.Response{
				Error: err.Error(),
				Details: map[string]any{
					"request": data,
					"message": "failed to create transaction",
				},
			},
		)

		return
	}

	for _, orderline := range transaction.Orderlines {
		orderline.TransactionID = int(lastInsertID)

		// Fetch the item information.
		item, err := mysql.GetItemByID(orderline.ItemID)
		if err != nil {
			log.Error(err, "failed to fetch orderline item",
				log.KVs(log.Map{"request": data, "orderline": orderline, "path": r.URL.Path}))

			response.InternalServer(w,
				response.Response{
					Error: "failed to fetch orderline item",
					Details: map[string]any{
						"request":          data,
						"item_id":          orderline.ItemID,
						"transaction_id":   lastInsertID,
						"transaction_type": transactionType,
					},
				},
			)

			return
		}

		// Create a new orderline for the said transaction.
		_, err = mysql.NewOrderline(transactionType, orderline)
		if err != nil {
			log.Error(err, "failed to create a new orderline",
				log.KVs(log.Map{
					"request":     data,
					"orderline":   orderline,
					"path":        r.URL.Path,
					"transaction": transaction,
				}),
			)
			response.InternalServer(w,
				response.Response{
					Error: "failed to create a new orderline",
					Details: map[string]any{
						"request":          data,
						"item_id":          item.ID,
						"transaction_id":   lastInsertID,
						"transaction_type": transaction.Type,
					},
				},
			)

			return
		}

		item.UpdateQuantity(transactionType, orderline.Quantity)

		err = mysql.UpdateItem(item)
		if err != nil {
			log.Error(err, "failed to update item",
				log.KVs(log.Map{
					"request": data,
					"item":    item,
					"path":    r.URL.Path,
				}),
			)

			response.InternalServer(w,
				response.Response{
					Error: "failed to update item",
					Details: map[string]any{
						"request":          data,
						"item_id":          item.ID,
						"transaction_id":   lastInsertID,
						"transaction_type": transactionType,
					},
				},
			)

			return
		}
	}

	response.Success(w, nil)
}
