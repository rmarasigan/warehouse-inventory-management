package validator

import (
	"log/slog"
	"os"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	"github.com/xeipuuv/gojsonschema"
)

// The base path for the schema validator specification JSON files.
const basePath = "/api/schema/validator/spec/"

// isValid validates the input JSON against the specified JSON schema.
//
// Parameters:
//   - input:  A byte slice containing the JSON data to be validated.
//   - source: The path to the JSON schema file.
func isValid(input []byte, source string) (bool, []string) {
	var (
		errors   []string
		pwd, _   = os.Getwd()
		request  = gojsonschema.NewBytesLoader(input)
		document = gojsonschema.NewReferenceLoader("file://" + pwd + source)
	)

	schema, err := gojsonschema.NewSchema(document)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(input)), slog.String("source", source))
		return false, nil
	}

	result, err := schema.Validate(request)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(input)), slog.String("source", source))
		return false, nil
	}

	if result == nil || result.Valid() {
		return true, nil
	}

	for _, err := range result.Errors() {
		errors = append(errors, err.Description())
	}

	return false, errors
}

// ValidateUser validates the input JSON against the users schema.
//
// Returns:
//   - bool: 'true' if the input is valid, 'false' otherwise.
//   - []string: A list of error message if the validation fails.
func ValidateUser(input []byte) (bool, []string) {
	return isValid(input, basePath+"users.json")
}

// ValidateRole validates the input JSON against the roles schema.
//
// Returns:
//   - bool: 'true' if the input is valid, 'false' otherwise.
//   - []string: A list of error message if the validation fails.
func ValidateRole(input []byte) (bool, []string) {
	return isValid(input, basePath+"roles.json")
}

// ValidateStorage validates the input JSON against the storages schema.
//
// Returns:
//   - bool: 'true' if the input is valid, 'false' otherwise.
//   - []string: A list of error message if the validation fails.
func ValidateStorage(input []byte) (bool, []string) {
	return isValid(input, basePath+"storages.json")
}

// ValidateUOM validates the input JSON against the uoms schema.
//
// Returns:
//   - bool: 'true' if the input is valid, 'false' otherwise.
//   - []string: A list of error message if the validation fails.
func ValidateUOM(input []byte) (bool, []string) {
	return isValid(input, basePath+"uoms.json")
}
