package v1

import (
	"log/slog"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

const (
	users            string = "users"
	activateUser     string = users + "/activate"
	roles            string = "roles"
	storages         string = "storages"
	uoms             string = "uoms"
	currencies       string = "currencies"
	activateCurrency string = currencies + "/activate"
	items            string = "items"
)

func IsValidMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
		return true

	default:
		return false
	}
}

func IsValidPathMethod(method, segment string) bool {
	var valid = map[string][]string{
		users:            {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		activateUser:     {http.MethodPut},
		roles:            {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		storages:         {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		uoms:             {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		items:            {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		currencies:       {http.MethodGet},
		activateCurrency: {http.MethodPut},
	}

	methods, exist := valid[segment]
	if !exist {
		log.Warn("provided path segment is invalid", slog.String("path", segment))
		return exist
	}

	for _, value := range methods {
		if value == method {
			return true
		}
	}

	log.Warn("provided path segment or method is invalid",
		slog.String("path", segment),
		slog.String("method", method))

	return false
}
