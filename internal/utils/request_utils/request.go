package requestutils

import (
	"net/http"
	"strings"
)

func HasQueryParam(request *http.Request, key string) (string, bool) {
	ok := request.URL.Query().Has(key)
	if !ok {
		return "", false
	}

	value := strings.TrimSpace(request.URL.Query().Get(key))

	return value, true
}
