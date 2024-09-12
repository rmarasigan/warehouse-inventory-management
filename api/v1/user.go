package v1

import (
	"database/sql"
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

// getUsers handles the HTTP request to retrieve a list of users. It writes
// the list of users to the HTTP response with an HTTP OK status. If an error
// occurs, it writes an HTTP Internal Server Error status.
func getUsers(w http.ResponseWriter) {
	defer func() {
		log.Panic()
		mysql.Close()
	}()

	mysql.Connect()
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
	defer func() {
		log.Panic()
		mysql.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		response.InternalServer(w, response.Response{Error: "failed to read request body"})
	}
	defer r.Body.Close()

	ok, errors := validator.ValidateUser(body)
	if !ok {
		if len(errors) > 0 {
			log.Error(strings.Join(errors, ", "), slog.String("request", string(body)))
			response.BadRequest(w, response.Response{Error: strings.Join(errors, ", ")})
		}

		response.InternalServer(w, nil)
	}

	data, err := apischema.NewUser(body)
	if err != nil {
		log.Error(err.Error(), slog.String("request", string(body)))
		response.InternalServer(w, response.Response{Error: "failed to unmarshal request body"})
	}

	var users = convert.Schema(data, func(user apischema.User) schema.User {
		return schema.User{
			RoleID:      user.RoleID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       sql.NullString{String: user.Email},
			Password:    user.Password,
			DateCreated: user.SetDateCreated(),
		}
	})

	mysql.Connect()
	for _, user := range users {
		existing, err := mysql.UserExists(user)
		if err != nil {
			log.Error(err.Error(), slog.Any("user", user), slog.Any("request", users))
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

	response.Success(w, nil)
}
