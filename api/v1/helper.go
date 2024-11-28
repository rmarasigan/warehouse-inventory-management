package v1

import (
	"errors"
	"fmt"
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

func getList[T any](r *http.Request, get func(id int) ([]T, error), list func() ([]T, error)) ([]T, error) {
	idParam := strings.TrimSpace(r.URL.Query().Get("id"))

	if len(idParam) > 0 {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return nil, err
		}

		// Fetch data for the given id
		result, err := get(id)
		if err != nil {
			return nil, fmt.Errorf("get: %w", err)
		}

		return result, nil

	}

	// Fetch all the data
	result, err := list()
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	return result, nil
}
