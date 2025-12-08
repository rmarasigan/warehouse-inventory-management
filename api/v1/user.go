package v1

import (
	"net/http"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	dbutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/db_utils"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	requestutils "github.com/rmarasigan/warehouse-inventory-management/internal/utils/request_utils"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)

	case http.MethodPost:
		createUser(w, r)

	case http.MethodPut:
		if strings.HasSuffix(r.URL.Path, activateUser) {
			activateUserAccount(w, r)

		} else {
			updateUser(w, r)
		}

	case http.MethodDelete:
		deleteUser(w, r)
	}
}

// getUsers handles the HTTP request to retrieve a list of users. It writes
// the list of users to the HTTP response with an HTTP OK status. If an error
// occurs, it writes an HTTP Internal Server Error status.
func getUsers(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	list, err := getList(r, mysql.GetUserByID, mysql.ListUser)
	if err != nil {
		log.Error(err, "failed to retrieve users", log.KV("path", r.URL.Path))
		response.InternalServer(w, nil)
	}

	users := convert.SchemaList(list, func(user schema.User) apischema.User {
		return apischema.User{
			ID:           user.ID,
			RoleID:       user.RoleID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        dbutils.GetString(user.Email),
			LastLogin:    dbutils.GetString(user.LastLogin),
			Active:       user.Active,
			DateCreated:  user.DateCreated,
			DateModified: dbutils.GetTime(user.DateModified),
		}
	})

	response.Success(w, users)
}

// createUser handles the HTTP request to create a new user. It validates
// the request body, unmarshals it into a user object, checks if the user
// already exists in the database, and inserts it if not. If the request
// body is invalid or the user already exists, it writes an appropriate
// HTTP status response.
func createUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	validationErrors, err := requestutils.ValidateRequest(body, validator.ValidateStorage)
	if err != nil && len(validationErrors) > 0 {
		log.Error(err, validationErrors, log.KVs(log.Map{"request": string(body), "path": r.URL.Path}))
		response.BadRequest(w, response.Response{Error: err.Error(), Details: validationErrors})

		return
	}

	data, err := requestutils.Unmarshal(r.URL.Path, body, apischema.NewUser)
	if err != nil {
		response.BadRequest(w, response.Response{Error: "failed to unmarshal request body"})
		return
	}

	users := convert.SchemaList(data, func(user apischema.User) schema.User {
		return schema.User{
			RoleID:    user.RoleID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     dbutils.SetString(user.Email),
			Password:  user.Password,
			Active:    true,
		}
	})

	for _, user := range users {
		existing, err := mysql.UserExists(user)
		if err != nil {
			log.Error(err, "failed to validate if user exists",
				log.KVs(log.Map{
					"request": string(body),
					"user":    user,
					"path":    r.URL.Path,
				}),
			)
			response.InternalServer(w,
				response.Response{
					Error: err.Error(),
					Details: map[string]any{
						"user":    user,
						"message": "failed to validate if user exists",
					},
				})

			return
		}

		if !existing {
			_, err = mysql.NewUser(user)
			if err != nil {
				log.Error(err, "failed to create new user",
					log.KVs(log.Map{
						"user": user,
						"path": r.URL.Path,
					}),
				)
				response.InternalServer(w,
					response.Response{
						Error: err.Error(),
						Details: map[string]any{
							"request": data,
							"user":    user,
							"message": "failed to create new user",
						},
					},
				)

				return
			}
		}
	}

	response.Created(w, nil)
}

// updateUser handles the HTTP request to update/modify user information.
// It unmarshals the request body into a user object and updates the corresponding
// fields. If an error occurs, it responds with an HTTP Internal Server Error status.
func updateUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
		log.Panic()
	}()

	body, err := requestutils.ReadBody(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	data, err := requestutils.Unmarshal(r.URL.Path, body, apischema.NewUser)
	if err != nil {
		response.BadRequest(w, response.Response{Error: "failed to unmarshal request body"})
		return
	}

	users := convert.SchemaList(data, func(user apischema.User) schema.User {
		return schema.User{
			ID:        user.ID,
			RoleID:    user.RoleID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     dbutils.SetString(user.Email),
			Password:  user.Password,
		}
	})

	for _, user := range users {
		err = mysql.UpdateUser(user)
		if err != nil {
			log.Error(err, "failed to update user",
				log.KVs(log.Map{"request": data, "user": user, "path": r.URL.Path}))

			response.InternalServer(w,
				response.Response{
					Error: err.Error(),
					Details: map[string]any{
						"request": data,
						"user":    user,
						"message": "failed to update user",
					},
				},
			)

			return
		}
	}

	response.Success(w, nil)
}

func activateUserAccount(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	err = mysql.ActivateUser(id)
	if err != nil {
		log.Error(err, "failed to activate user account",
			log.KVs(log.Map{"id": id, "path": r.URL.Path}))

		response.InternalServer(w, response.Response{Error: "failed to activate user account"})

		return
	}

	response.Success(w, nil)
}

// deleteUser handles the HTTP request to delete a user. If an error occurs,
// it writes either HTTP Internal Server Error or Bad Request status.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	id, err := parameterID(r)
	if err != nil {
		response.BadRequest(w, response.Response{Error: err.Error()})
		return
	}

	err = mysql.DeleteUser(id)
	if err != nil {
		log.Error(err, "failed to delete user",
			log.KVs(log.Map{"id": id, "path": r.URL.Path}))

		response.InternalServer(w, response.Response{Error: err.Error(), Details: "failed to delete user"})

		return
	}

	response.Success(w, nil)
}
