package v1

import (
	"log/slog"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func IsValidMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost:
		return true

	default:
		return false
	}
}

func IsValidPathMethod(method, segment string) bool {
	var valid = map[string][]string{
		"users":     {http.MethodGet},
		"users/new": {http.MethodPost},
	}

	methods, exist := valid[segment]
	if !exist {
		log.Warn("provided path segment is invalid", slog.String("path", segment))
		return exist
	}

	for _, value := range methods {
		return value == method
	}

	log.Warn("provided path segment and method is invalid",
		slog.String("path", segment),
		slog.String("method", method))

	return false
}
