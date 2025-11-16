package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Details any    `json:"details,omitempty"`
}

func response(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if data != nil {
		body, err := json.MarshalIndent(data, "", "  ")
		if err == nil {
			_, err := w.Write(body)
			if err != nil {
				log.Error(err, "unable to write data", log.KV("data", data))
				return
			}

			return
		}

		log.Error(err, "unable to encode data", log.KV("data", data))

		return
	}
}

func Success(w http.ResponseWriter, data any) {
	response(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data any) {
	response(w, http.StatusCreated, data)
}

func MultiStatus(w http.ResponseWriter, data any) {
	response(w, http.StatusMultiStatus, data)
}

func BadRequest(w http.ResponseWriter, data any) {
	response(w, http.StatusBadRequest, data)
}

func NotFound(w http.ResponseWriter, data any) {
	response(w, http.StatusNotFound, data)
}

func InternalServer(w http.ResponseWriter, data any) {
	response(w, http.StatusInternalServerError, data)
}

func MethodNotAllowed(w http.ResponseWriter, method string) {
	response(w, http.StatusMethodNotAllowed, Response{Message: fmt.Sprintf("method '%s' is not supported", method)})
}

func NotImplemented(w http.ResponseWriter, data any) {
	response(w, http.StatusNotImplemented, data)
}
