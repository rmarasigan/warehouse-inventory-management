package requestutils

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func HasQueryParam(request *http.Request, key string) (string, bool) {
	ok := request.URL.Query().Has(key)
	if !ok {
		return "", false
	}

	value := strings.TrimSpace(request.URL.Query().Get(key))

	return value, true
}

func ReadBody(request *http.Request) ([]byte, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error(err, "failed to read request body", log.KV("path", request.URL.Path))
		return nil, err
	}

	if len(body) == 0 {
		emptyErr := errors.New("request body cannot be empty")
		log.Error(emptyErr, "missing request body", log.KV("path", request.URL.Path))

		return nil, emptyErr
	}

	return body, nil
}

func Unmarshal[T any](path string, body []byte, fn func([]byte) (T, error)) (T, error) {
	data, err := fn(body)
	if err != nil {
		log.Error(err, "failed to unmarshal request body",
			log.KVs(log.Map{"request": string(body), "path": path}))

		// Cannot return 'nil' because 'T' may not be nillable type.
		// We need to return a zero value for 'T' instead.
		var zero T
		return zero, err
	}

	return data, nil
}
