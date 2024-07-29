package v1

import (
	"log/slog"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func Handler(w http.ResponseWriter, r *http.Request, segment string) {
	var method = r.Method

	switch {
	case !IsValidMethod(method):
		log.Warn("unhandled method", slog.String("method", method))
		response.MethodNotAllowed(w, method)

		return

	case !IsValidPathMethod(method, segment):
		log.Warn("invalid path and method", slog.String("path", segment), slog.String("method", method))
		response.MethodNotAllowed(w, method)

		return

	default:
		return
	}
}
