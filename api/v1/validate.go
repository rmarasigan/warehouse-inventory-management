package v1

import (
	"log/slog"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

const (
	users string = "users"
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
		users: {http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
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
