# Warehouse Inventory Management

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

**GET: Fetch List of User Information**
```bash
curl -X GET "http://0.0.0.0:8080/api/v1/users"
```

**POST: Create a New User**
```bash
curl -X POST "http://0.0.0.0:8080/api/v1/users/new" -H "Content-Type: application/json" -d '[{"role_id": 1, "first_name": "John", "last_name": "Doe", "email": "j.doe@example.com", "password": "john-doe-password"}, {"role_id": 1, "first_name": "Alice", "last_name": "Park", "password": "alice-password"}]'
```

**DELETE: Remove a User**
```bash
curl -X DELETE "http://0.0.0.0:8080/api/v1/users/delete?id=1"
```

**PUT: Update User Information**
```bash
curl -X PUT "http://0.0.0.0:8080/api/v1/users/update" -H "Content-Type: application/json" -d '[{"id": 16, "email": "j.doe@example.com", "password": "your-new-password"}]'
```

## Reference
* [gojsonschema](https://github.com/xeipuuv/gojsonschema)
* [OpenAPI Documentation](https://learn.openapis.org/)
* [Go database/sql tutorial](http://go-database-sql.org/index.html)
* [Illustrated guide to SQLX](https://jmoiron.github.io/sqlx/)
* [Why even use *DB.exec() or prepared statements in Golang?](https://stackoverflow.com/questions/50664648/why-even-use-db-exec-or-prepared-statements-in-golang)