package v1

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func roleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRoles(w)

	case http.MethodPost:
		createRole(w, r)
	}
}

// getRoles handles the HTTP request to retrieve a list of roles. It writes
// the list of roles to the HTTP response with an HTTP OK status. If an error
// occurs, it writes an HTTP Internal Server Error status.
func getRoles(w http.ResponseWriter) {
	defer log.Panic()

	roles, err := mysql.RoleList()
	if err != nil {
		log.Error(err.Error())
		response.InternalServer(w, nil)
	}

	response.Success(w, roles)
}

// createRole handles the HTTP request to create a new role. It validates
// the request body, unmarshals it into a role object, checks if the role
// already exists in the database, and inserts it if not. If the request
// body is invalid or the role already exist, it writes an appropriate HTTP
// status response.
func createRole(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to read request body"})
	}
	defer r.Body.Close()

	ok, errors := validator.ValidateRole(body)
	if !ok {
		if len(errors) > 0 {
			log.Error(strings.Join(errors, ", "), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
			response.BadRequest(w, response.Response{Error: strings.Join(errors, ", ")})
		}

		response.InternalServer(w, nil)
	}

	data, err := apischema.NewRole(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})
	}

	var roles = convert.Schema(data, func(role apischema.Role) schema.Role {
		return schema.Role{
			ID:   role.ID,
			Name: role.Name,
		}
	})

	for _, role := range roles {
		existing, err := mysql.RoleNameExists(role)
		if err != nil {
			log.Error(err.Error(), slog.Any("role", role), slog.Any("request", roles), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to validate if role name exists", Details: role})

			break
		}

		if !existing {
			err = mysql.NewRole(role)
			if err != nil {
				log.Error(err.Error(), slog.Any("role", role), slog.Any("request", roles))
				response.InternalServer(w, response.Response{Error: "failed to create new role", Details: role})

				break
			}
		}
	}

	response.Created(w, nil)
}
