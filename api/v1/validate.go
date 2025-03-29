package v1

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

const (
	users              string = "users"
	activateUser       string = users + "/activate"
	roles              string = "roles"
	storages           string = "storages"
	uoms               string = "uoms"
	currencies         string = "currencies"
	activateCurrency   string = currencies + "/activate"
	items              string = "items"
	transaction        string = "transaction"
	transactionInbound string = transaction + "/inbound"
)

func IsValidPathMethod(method, segment string) bool {
	var valid = map[string][]string{
		users:              {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		activateUser:       {http.MethodPut},
		roles:              {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		storages:           {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		uoms:               {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		currencies:         {http.MethodGet},
		activateCurrency:   {http.MethodPut},
		items:              {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		transaction:        {http.MethodGet},
		transactionInbound: {http.MethodPost},
	}

	methods, exist := valid[segment]
	if exist {
		if slices.Contains(methods, method) {
			return true
		}

		log.Warn("invalid method for path", slog.String("path", segment), slog.String("method", method))
		return false
	}

	log.Warn("invalid path", slog.String("path", segment), slog.String("method", method))
	return false
}
