# Warehouse Inventory Management

A Warehouse Inventory Management system is designed to help you manage and track inventory within a warehouse.

### Running the Application

```bash
dev@dev:~/warehouse-inventory-management$ go build; ./warehouse-inventory-management
2024-08-08 2:59:39 PM | INFO    Establishing MySQL connection...
2024-08-08 2:59:39 PM | OK      MySQL connection established.
2024-08-08 2:59:39 PM | INFO    Initializing Warehouse Inventory Management Server at 0.0.0.0:8080
                  ______   _____
   ___  ___  ___ /  /    \/     \
  /  / /  / /  //  /  __   __   /
 /  / /  / /  //  /  / /  / /  /
/  /_/  /_/  //  /  / /  / /  /
\_____,_____//__/__/ /__/ /__/
```

## Requirements
### Software
* **Go**: v1.21
* **MySQL**: v8.0.35

### Configuration
The application uses a configuration file: [`wim-config.yaml`](wim-config.yaml). Ensure this file is properly set up before running the application.

## API Validation
API validation schemas are generated from the [`api-specification.yaml`](api-specification.yaml) using the tool [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema).

**Installation**

To install [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema), run:
```bash
pip install openapi2jsonschema
```

**Generating JSON Schemas**

To generate the JSON schemas from the [`api-specification.yaml`](api-specification.yaml):
```bash
openapi2jsonschema api-specification.yaml -o ./api/schema/validator/spec
```

This will create a set of JSON schemas in `./api/schema/validator/spec` directory.

### Users
Below are example `cURL` commands for interacting with the `/api/v1/users/` API.

**Local Endpoint**: `http://0.0.0.0:8080/api/v1/users`

#### Operations

**`GET`: Fetch List of User Information**
```bash
curl -X GET "http://0.0.0.0:8080/api/v1/users"
```

**`POST`: Create a New User**
```bash
curl -X POST "http://0.0.0.0:8080/api/v1/users" -H "Content-Type: application/json" -d '[{"role_id": 1, "first_name": "John", "last_name": "Doe", "email": "j.doe@example.com", "password": "john-doe-password"}, {"role_id": 1, "first_name": "Alice", "last_name": "Park", "password": "alice-password"}]'
```

**`DELETE`: Remove a User**
```bash
curl -X DELETE "http://0.0.0.0:8080/api/v1/users?id=1"
```

**`PUT`: Update User Information**
```bash
curl -X PUT "http://0.0.0.0:8080/api/v1/users" -H "Content-Type: application/json" -d '[{"id": 16, "email": "j.doe@example.com", "password": "your-new-password"}]'
```

### Roles
Below are example `cURL` commands for interacting with the `/api/v1/roles/` API.

**Local Endpoint**: `http://0.0.0.0:8080/api/v1/roles`

#### Operations

**`GET`: Fetch List of Role Information**
```bash
curl -X GET "http://0.0.0.0:8080/api/v1/roles"
```

**`POST`: Create a New Role**
```bash
curl -X POST "http://0.0.0.0:8080/api/v1/roles" -H "Content-Type: application/json" -d '[{"name": "Administrator"}]'
```

## Reference
* [gojsonschema](https://github.com/xeipuuv/gojsonschema)
* [OpenAPI Documentation](https://learn.openapis.org/)
* [Go database/sql tutorial](http://go-database-sql.org/index.html)
* [Illustrated guide to SQLX](https://jmoiron.github.io/sqlx/)
* [Why even use *DB.exec() or prepared statements in Golang?](https://stackoverflow.com/questions/50664648/why-even-use-db-exec-or-prepared-statements-in-golang)