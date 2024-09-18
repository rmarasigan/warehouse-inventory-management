package v1

import (
	"log/slog"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

const (
	Users      string = "users"
	NewUser    string = "users/new"
	UpdateUser string = "users/update"
	DeleteUser string = "users/delete"
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
		Users:      {http.MethodGet},
		NewUser:    {http.MethodPost},
		UpdateUser: {http.MethodPut},
		DeleteUser: {http.MethodDelete},
	}

	methods, exist := valid[segment]
	if !exist {
		log.Warn("provided path segment is invalid", slog.String("path", segment))
		return exist
	}

	for _, value := range methods {
		return value == method
	}

	log.Warn("provided path segment or method is invalid",
		slog.String("path", segment),
		slog.String("method", method))

	return false
}
