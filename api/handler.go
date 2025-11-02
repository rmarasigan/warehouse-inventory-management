package api

import (
	"net/http"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	v1 "github.com/rmarasigan/warehouse-inventory-management/api/v1"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	parts := strings.Split(path, "/")

	if len(parts) > 1 {
		var (
			version = parts[0]
			segment = strings.Join(parts[1:], "/")
		)

		switch version {
		case "v1":
			v1.Handler(w, r, segment)

		default:
			response.NotFound(w, response.Response{Message: "unrecognized version"})
		}
	}
}
