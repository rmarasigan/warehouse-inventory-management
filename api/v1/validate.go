package v1

import (
	"net/http"
	"slices"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

const (
	users             string = "users"
	activateUser      string = users + "/activate"
	roles             string = "roles"
	storages          string = "storages"
	uoms              string = "uoms"
	currencies        string = "currencies"
	activateCurrency  string = currencies + "/activate"
	items             string = "items"
	transaction       string = "transactions"
	transactionNote   string = transaction + "/note"
	orderlinesNote    string = transaction + "/orderline-note"
	transactionCancel string = transaction + "/cancel"
)

func isValidPathMethod(method, segment string) bool {
	var valid = map[string][]string{
		users:             {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		activateUser:      {http.MethodPut},
		roles:             {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		storages:          {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		uoms:              {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		currencies:        {http.MethodGet},
		activateCurrency:  {http.MethodPut},
		items:             {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		transaction:       {http.MethodGet, http.MethodPost},
		transactionNote:   {http.MethodPut},
		orderlinesNote:    {http.MethodPut},
		transactionCancel: {http.MethodPut},
	}

	methods, exist := valid[segment]
	if exist {
		if slices.Contains(methods, method) {
			return true
		}

		log.Warn("invalid method for path", log.KVs(log.Map{"path": segment, "method": method}))
		return false
	}

	log.Warn("invalid path", log.KVs(log.Map{"path": segment, "method": method}))
	return false
}
