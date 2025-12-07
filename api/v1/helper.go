package v1

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	requestutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/request_utils"
)

func parameterID(r *http.Request) (int, error) {
	idParam, ok := requestutils.HasQueryParam(r, "id")
	if !ok {
		errMsg := errors.New("missing 'id' in the request query parameter")
		log.Error(errMsg, "query parameter 'id' is required", log.KV("path", r.URL.Path))

		return 0, errMsg
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Error(err, "failed to parse 'id' query parameter",
			log.KVs(log.Map{
				"id":   idParam,
				"path": r.URL.Path,
			}),
		)

		return 0, errors.New("invalid 'id' value")
	}

	return id, nil
}

func getList[T any](r *http.Request, get func(id int) (T, error), list func() ([]T, error)) ([]T, error) {
	// Check if the "id" parameter is provided.
	idParam, ok := requestutils.HasQueryParam(r, "id")
	if !ok {
		// Fetch all the data
		result, err := list()
		if err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}

		// Ensure result is not nil
		if result == nil {
			result = []T{}
		}

		return result, nil
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, err
	}

	// Fetch data for the given ID
	item, err := get(id)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	// Return an empty array if it is zero value
	if reflect.ValueOf(item).IsZero() {
		return []T{}, nil
	}

	// Wrap the single item in an array
	return []T{item}, nil
}
