package requestutils

import (
	"errors"
	"strings"
)

func ValidateRequest(input []byte, fn func([]byte) (bool, []string)) (string, error) {
	ok, validationErrors := fn(input)
	if !ok {
		return strings.Join(validationErrors, ", "), errors.New("invalid request body")
	}

	return "", nil
}
