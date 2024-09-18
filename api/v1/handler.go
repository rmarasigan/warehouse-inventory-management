package v1

import (
	"log/slog"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func Handler(w http.ResponseWriter, r *http.Request, segment string) {
	defer log.Panic()
	var method = r.Method

	switch {
	case !IsValidMethod(method):
		log.Warn("unhandled method", slog.String("method", method))
		response.MethodNotAllowed(w, method)

		return

	case !IsValidPathMethod(method, segment):
		response.BadRequest(w, response.Response{Error: "provided path or method is invalid"})
		return

	default:
		switch segment {
		case Users:
			getUsers(w)

		case NewUser:
			createUser(w, r)

		case UpdateUser:
			updateUser(w, r)

		case DeleteUser:
			deleteUser(w, r)
		}

		return
	}
}
