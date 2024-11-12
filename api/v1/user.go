package v1

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/rmarasigan/warehouse-inventory-management/api/response"
	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
	"github.com/rmarasigan/warehouse-inventory-management/api/schema/validator"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w)

	case http.MethodPost:
		createUser(w, r)

	case http.MethodPut:
		updateUser(w, r)

	case http.MethodDelete:
		deleteUser(w, r)
	}
}

// getUsers handles the HTTP request to retrieve a list of users. It writes
// the list of users to the HTTP response with an HTTP OK status. If an error
// occurs, it writes an HTTP Internal Server Error status.
func getUsers(w http.ResponseWriter) {
	defer log.Panic()

	users, err := mysql.UserList()
	if err != nil {
		log.Error(err.Error())
		response.InternalServer(w, nil)
	}

	response.Success(w, users)
}

// createUser handles the HTTP request to create a new user. It validates
// the request body, unmarshals it into a user object, checks if the user
// already exists in the database, and inserts it if not. If the request
// body is invalid or the user already exists, it writes an appropriate
// HTTP status response.
func createUser(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to read request body"})
	}
	defer r.Body.Close()

	ok, errors := validator.ValidateUser(body)
	if !ok {
		if len(errors) > 0 {
			log.Error(strings.Join(errors, ", "), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
			response.BadRequest(w, response.Response{Error: strings.Join(errors, ", ")})
		}

		response.InternalServer(w, nil)
	}

	data, err := apischema.NewUser(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})
	}

	var users = convert.Schema(data, func(user apischema.User) schema.User {
		return schema.User{
			RoleID:      user.RoleID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       sql.NullString{String: user.Email, Valid: true}, // Valid is 'true' if String is not NULL
			Password:    user.Password,
			DateCreated: user.SetDateCreated(),
		}
	})

	for _, user := range users {
		existing, err := mysql.UserExists(user)
		if err != nil {
			log.Error(err.Error(), slog.Any("user", user), slog.Any("request", users), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to validate if user exists", Details: user})

			break
		}

		if !existing {
			err = mysql.NewUser(user)
			if err != nil {
				log.Error(err.Error(), slog.Any("user", user), slog.Any("request", users))
				response.InternalServer(w, response.Response{Error: "failed create new user account", Details: user})

				break
			}
		}
	}

	response.Created(w, nil)
}

// updateUser handles the HTTP request to update/modify user information.
// It unmarshals the request body into a user object and updates the corresponding
// fields. If an error occurs, it responds with an HTTP Internal Server Error status.
func updateUser(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error(), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to read request body"})
	}
	defer r.Body.Close()

	users, err := apischema.NewUser(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})
	}

	for _, user := range users {
		userID := fmt.Sprint(user.ID)

		dbusers, err := mysql.RetrieveUsers(userID)
		if err != nil {
			log.Error(err.Error(), slog.Any("id", userID), slog.Any("path", r.URL.Path))
			response.InternalServer(w, response.Response{Error: "failed to retrieve user account"})

			break
		}

		for _, dbuser := range *dbusers {
			dbuser.UpdateValues(user)

			err := mysql.UpdateUser(dbuser)
			if err != nil {
				log.Error(err.Error(), slog.Any("id", userID), slog.Any("path", r.URL.Path))
				response.InternalServer(w, response.Response{Error: "failed to update user account"})

				break
			}
		}
	}

	response.Success(w, nil)
}

// deleteUser handles the HTTP request to delete a user. If an error occurs,
// it writes either HTTP Internal Server Error or Bad Request status.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	defer log.Panic()

	userID := r.URL.Query().Get("id")
	if strings.TrimSpace(userID) == "" {
		log.Error("user 'id' is required", slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "missing user 'id' query parameter"})
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Error(err.Error(), slog.Any("id", id), slog.Any("path", r.URL.Path))
		response.BadRequest(w, response.Response{Error: "invalid user 'id' value"})
	}

	affected, err := mysql.DeleteUser(id)
	if err != nil {
		log.Error(err.Error(), slog.Any("id", id), slog.Any("path", r.URL.Path))
		response.InternalServer(w, response.Response{Error: "failed to delete user account"})
	}

	response.Success(w, response.Response{Message: fmt.Sprintf("%d row(s) affected", affected)})
}
