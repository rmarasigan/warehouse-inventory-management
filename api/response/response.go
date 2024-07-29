package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

type Response struct {
	Message string `json:"message"`
}

func response(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if data != nil {
		body, err := json.MarshalIndent(data, "", "  ")
		if err == nil {
			_, err := w.Write(body)
			if err != nil {
				log.Error("unable to write data", slog.Any("error", err), slog.Any("data", data))
				return
			}

			return
		}

		log.Error("unable to encode data", slog.Any("error", err), slog.Any("data", data))

		return
	}
}

func Success(w http.ResponseWriter, data interface{}) {
	response(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data interface{}) {
	response(w, http.StatusCreated, data)
}

func BadRequest(w http.ResponseWriter, data interface{}) {
	response(w, http.StatusBadRequest, data)
}

func NotFound(w http.ResponseWriter, data interface{}) {
	response(w, http.StatusNotFound, data)
}

func InternalServer(w http.ResponseWriter, data interface{}) {
	response(w, http.StatusInternalServerError, data)
}

func MethodNotAllowed(w http.ResponseWriter, method string) {
	response(w, http.StatusMethodNotAllowed, Response{Message: fmt.Sprintf("method '%s' is not supported", method)})
}

func NotImplemented(w http.ResponseWriter, data interface{}) {
	response(w, http.StatusNotImplemented, data)
}
