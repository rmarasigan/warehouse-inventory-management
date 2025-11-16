package v1

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	dbutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/db_utils"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	requestutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/request_utils"
)

func transactionCancelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		cancelTransaction(w, r)
	}
}

func cancelTransaction(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	userIDParam, ok := requestutils.HasQueryParam(r, "user_id")
	if !ok {
		err := errors.New("missing 'user_id' from request query")
		log.Error(err, "query parameter 'user_id' is required", log.KV("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: err.Error()})

		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		log.Error(err, "failed to parse 'user_id' query parameter",
			log.KVs(map[string]any{
				"id":   userIDParam,
				"path": r.URL.Path,
			}),
		)
		response.BadRequest(w, response.Response{Error: "invalid 'user_id' value; must be an integer"})
		return
	}

	transaction, err := mysql.GetTransactionByID(id)
	if err != nil {
		log.Error(err, "failed to retrieve transaction",
			log.KVs(map[string]any{
				"id":   id,
				"path": r.URL.Path,
			}),
		)

		response.InternalServer(w,
			response.Response{
				Error:   err.Error(),
				Details: "failed to fetch transaction " + fmt.Sprint(id),
			},
		)

		return
	}

	type FailedOrderline struct {
		TransactionID int              `json:"transaction_id"`
		ItemID        int              `json:"item_id,omitempty"`
		Orderline     schema.Orderline `json:"orderline"`
		OrderlineID   int              `json:"orderline_id"`
		Reason        string           `json:"reason"`
	}

	var failed []FailedOrderline

	for _, orderline := range transaction.Orderlines {
		itemID := orderline.ItemID
		orderlineID := orderline.ID

		// Fetch the item information.
		item, err := mysql.GetItemByID(itemID)
		if err != nil {
			log.Warn("failed to fetch item id",
				log.KVs(map[string]any{
					"error":          err.Error(),
					"path":           r.URL.Path,
					"item_id":        itemID,
					"transaction_id": transaction.ID,
				}),
			)

			failed = append(failed, FailedOrderline{
				TransactionID: id,
				ItemID:        itemID,
				OrderlineID:   orderlineID,
				Reason:        "failed to fetch item: " + err.Error(),
			})

			continue
		}

		// Update the item quantity based on the transaction type.
		item.UpdateCancelledQuantity(transaction.Type, orderline.Quantity)
		err = mysql.UpdateItem(item)
		if err != nil {
			log.Warn("failed to update item quantity",
				log.KVs(map[string]any{
					"error":          err.Error(),
					"path":           r.URL.Path,
					"item_id":        item.ID,
					"transaction_id": transaction.ID,
				}),
			)

			failed = append(failed, FailedOrderline{
				TransactionID: id,
				ItemID:        itemID,
				OrderlineID:   orderlineID,
				Reason:        "failed to update item quantity: " + err.Error(),
			})

			continue
		}

		orderline.IsVoided = dbutils.SetBool(true)
		orderline.UpdatedBy = dbutils.SetInt(int32(userID))

		err = mysql.CancelOrderline(orderline)
		if err != nil {
			log.Warn("failed to cancel orderline",
				log.KVs(map[string]any{
					"error":          err.Error(),
					"path":           r.URL.Path,
					"item_id":        item.ID,
					"transaction_id": transaction.ID,
					"orderline_id":   orderline.ID,
				}),
			)

			failed = append(failed, FailedOrderline{
				TransactionID: id,
				ItemID:        itemID,
				OrderlineID:   orderlineID,
				Reason:        "failed to cancel orderline: " + err.Error(),
			})

			continue
		}
	}

	transaction.IsCancelled = dbutils.SetBool(true)
	transaction.UpdatedBy = dbutils.SetInt(int32(userID))

	err = mysql.CancelTransaction(transaction)
	if err != nil {
		log.Warn("failed to cancel transaction", slog.Any("error", err.Error()), slog.Int("transaction", transaction.ID), slog.Any("path", r.URL.Path))
		response.InternalServer(w,
			response.Response{
				Error: "failed to cancel transaction",
				Details: map[string]any{
					"transaction_id":   transaction.ID,
					"transaction_type": transaction.Type,
				},
			},
		)

		return
	}

	if len(failed) > 0 {
		err := errors.New("orderlines update not fully successful")
		log.Error(err, "failed to update orderlines", log.KV("errors", failed))
		response.MultiStatus(w,
			response.Response{
				Error: err.Error(),
				Details: map[string]any{
					"failed": failed,
				},
			},
		)

		return
	}

	response.Success(w,
		response.Response{
			Details: map[string]any{
				"transaction_id": transaction.ID,
				"is_cancelled":   true,
			},
		},
	)
}
