package v1

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
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
		log.Error(err.Error())
		response.InternalServer(w, nil)

		return
	}

	currencies := convert.Schema(list, func(currency schema.Currency) apischema.Currency {
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

	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if strings.TrimSpace(code) == "" {
		log.Error("missing 'code' as the query parameter", slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "missing 'code' as the query parameter"})
	}

	err := mysql.ActivateCurrency(code)
	if err != nil {
		log.Error(err.Error(), slog.Any("code", code))
		response.InternalServer(w, nil)
	}

	response.Success(w, nil)
}
