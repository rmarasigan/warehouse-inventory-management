package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	requestutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/request_utils"
)

func currencyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCurrencies(w, r)

	case http.MethodPut:
		updateCurrency(w, r)
	}
}

func getCurrencies(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	list, err := getList(r, mysql.GetCurrency, mysql.ListCurrency)
	if err != nil {
		log.Error(err, "failed to retrieve currency")
		response.InternalServer(w, response.NewError(err, "failed to retrieve currency"))

		return
	}

	currencies := convert.SchemaList(list, func(currency schema.Currency) apischema.Currency {
		return apischema.Currency{
			ID:     currency.ID,
			Code:   currency.Code,
			Symbol: currency.Symbol,
			Active: currency.Active,
		}
	})

	response.Success(w, currencies)
}

func updateCurrency(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	code, ok := requestutils.HasQueryParam(r, "code")
	if !ok {
		errMsg := errors.New("missing 'code' as the query parameter")
		log.Error(errMsg, "query parameter 'code' is required", log.KV("path", r.URL.Path))
		response.BadRequest(w, response.NewError(errMsg))

		return
	}

	err := mysql.ActivateCurrency(code)
	if err != nil {
		log.Error(err, "failed to activate currency", slog.Any("code", code))
		response.InternalServer(w, response.NewError(err, "failed to activate currency"))
	}

	response.Success(w, nil)
}
