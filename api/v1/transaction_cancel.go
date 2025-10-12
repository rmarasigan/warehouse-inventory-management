package v1

import (
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

	idParam, ok := requestutils.HasQueryParam(r, "id")
	if !ok {
		err := "missing 'id' from request query"
		log.Error(err, slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: err})

		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "the 'id' is invalid; must be an integer"})
		return
	}

	userIDParam, ok := requestutils.HasQueryParam(r, "user_id")
	if !ok {
		err := "missing 'user_id' from request query"
		log.Error(err, slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: err})

		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "the 'user_id' is invalid; must be an integer"})
		return
	}

	transaction, err := mysql.GetTransactionByID(id)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to fetch transaction " + idParam})
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
			log.Warn("failed to fetch item id", slog.Any("error", err.Error()), slog.Int("transaction_id", id), slog.Any("path", r.URL.Path))
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
			log.Warn("failed to update item quantity", slog.Any("error", err.Error()), slog.Int("transaction", id), slog.Int("item_id", item.ID), slog.Any("path", r.URL.Path))
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
			log.Warn("failed to cancel orderline", slog.Any("error", err.Error()), slog.Int("transaction", transaction.ID), slog.Int("item_id", item.ID), slog.Int("orderline_id", orderline.ID), slog.Any("path", r.URL.Path))
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
		log.Error("failed to process all orderline/s", slog.Any("errors", failed))
		response.MultiStatus(w,
			response.Response{
				Error: "failed to process all orderline/s",
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
