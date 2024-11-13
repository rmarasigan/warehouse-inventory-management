package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func parameterID(r *http.Request) (int, error) {
	idParam := strings.TrimSpace(r.URL.Query().Get("id"))

	if idParam == "" {
		log.Error("'id' is required", slog.Any("path", r.URL.Path))
		return 0, errors.New("missing 'id' query parameter")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error(err.Error(), slog.Any("id", idParam), slog.Any("path", r.URL.Path))
		return 0, errors.New("invalid 'id' value")
	}

	return id, nil
}
