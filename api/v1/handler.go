package v1

import (
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func Handler(w http.ResponseWriter, r *http.Request, segment string) {
	defer log.Panic()
	var method = r.Method

	// Validate the method and path
	if !IsValidPathMethod(method, segment) {
		response.BadRequest(w, response.Response{Error: "provided path or method is invalid"})
		return
	}

	handlers := map[string]func(http.ResponseWriter, *http.Request){
		users:            userHandler,
		activateUser:     userHandler,
		roles:            roleHandler,
		storages:         storageHandler,
		uoms:             uomHandler,
		currencies:       currencyHandler,
		activateCurrency: currencyHandler,
		items:            itemHandler,
		transaction:      transactionHandler,
	}

	// Handle the request if the segment is valid
	handler, exists := handlers[segment]
	if exists {
		handler(w, r)
		return
	}

	response.NotFound(w, response.Response{Message: "unrecognized path"})
}
