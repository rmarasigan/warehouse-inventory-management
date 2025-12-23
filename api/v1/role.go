package v1

import (
	"fmt"
	"net/http"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	requestutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/request_utils"
)

func roleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRoles(w, r)

	case http.MethodPost:
		createRole(w, r)

	case http.MethodPut:
		updateRole(w, r)

	case http.MethodDelete:
		deleteRole(w, r)
	}
}

// getRoles handles the HTTP request to retrieve a list of role(s) or a specific role.
// It writes the list of role(s) to the HTTP response with an HTTP OK status. If an
// error occurs, it writes an HTTP Internal Server Error status.
func getRoles(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	list, err := getList(r, mysql.GetRoleByID, mysql.ListRole)
	if err != nil {
		log.Error(err, "failed to retrieve roles", log.KV("path", r.URL.Path))
		response.InternalServer(w, response.NewError(err, "failed to retrieve roles"))

		return
	}

	roles := convert.SchemaList(list, func(role schema.Role) apischema.Role {
		return apischema.Role{
			ID:   role.ID,
			Name: role.Name,
		}
	})

	response.Success(w, roles)
}

// createRole handles the HTTP request to create a new role. It validates
// the request body, unmarshals it into a role object, checks if the role
// already exists in the database, and inserts it if not. If the request
// body is invalid or the role already exist, it writes an appropriate HTTP
// status response.
func createRole(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
	if err != nil {
		response.BadRequest(w, response.NewError(err))
		return
	}

	validationErrors, err := requestutils.ValidateRequest(body, validator.ValidateRole)
	if err != nil && len(validationErrors) > 0 {
		log.Error(err, validationErrors, log.KVs(log.Map{"request": string(body), "path": r.URL.Path}))
		response.BadRequest(w, response.NewError(err, validationErrors))

		return
	}

	data, err := requestutils.Unmarshal(r.URL.Path, body, apischema.NewRole)
	if err != nil {
		response.BadRequest(w, response.NewError(err, "failed to unmarshal request body"))
		return
	}

	roles := convert.SchemaList(data, func(role apischema.Role) schema.Role {
		return schema.Role{
			Name: role.Name,
		}
	})

	for _, role := range roles {
		_, err = mysql.NewRoleIfNotExists(role)
		if err != nil {
			log.Error(err, "failed to create role",
				log.KVs(log.Map{"role": role, "path": r.URL.Path}))

			response.InternalServer(w, response.NewError(err,
				map[string]any{
					"request": data,
					"role":    role,
					"message": "failed to create role",
				}),
			)

			return
		}
	}

	response.Created(w, nil)
}

// updateRole handles the HTTP request to update/modify role information. It
// unmarshals the request body into a role object and updates the corresponding
// fields. If an error occurs, it responds with an HTTP Internal Server Error status.
func updateRole(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
	if err != nil {
		response.BadRequest(w, response.NewError(err))
		return
	}

	data, err := requestutils.Unmarshal(r.URL.Path, body, apischema.NewRole)
	if err != nil {
		response.BadRequest(w, response.NewError(err, "failed to unmarshal request body"))
		return
	}

	var roles = convert.SchemaList(data, func(role apischema.Role) schema.Role {
		return schema.Role{
			ID:   role.ID,
			Name: role.Name,
		}
	})

	for _, role := range roles {
		err = mysql.UpdateRole(role)
		if err != nil {
			log.Error(err, "failed to update role",
				log.KVs(log.Map{"request": data, "role": role, "path": r.URL.Path}))

			response.InternalServer(w, response.NewError(err,
				map[string]any{
					"request": data,
					"role":    role,
					"message": "failed to update role",
				}),
			)

			return
		}
	}

	response.Success(w, nil)
}

// deleteRole handles the HTTP request to delete a role. If an error occurs,
// it writes either HTTP Internal Server Error or Bad Request status.
func deleteRole(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.NewError(err))
		return
	}

	affected, err := mysql.DeleteRole(id)
	if err != nil {
		log.Error(err, "failed to delete role", log.KVs(log.Map{"id": id, "path": r.URL.Path}))
		response.InternalServer(w, response.NewError(err, "failed to delete role"))

		return
	}

	response.Success(w, response.New(fmt.Sprintf("%d row(s) affected", affected)))
}
