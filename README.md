# Warehouse Inventory Management

A simple Warehouse Inventory Management System that is designed to help you manage and track inventory within a warehouse.

> [!WARNING]
>
> This project is for learning purposes and should not be used in production environments.

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
* **Go**: v1.22
* **MySQL**: v8.0.35

### Configuration
> [!IMPORTANT]
>
> The application uses a configuration file: [`wim-config.yaml`](wim-config.yaml).
>
> **Ensure this file is properly set up before running the application.**

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

## Reference
* [gojsonschema](https://github.com/xeipuuv/gojsonschema)
* [`time.Time` support](https://github.com/go-sql-driver/mysql?tab=readme-ov-file#timetime-support)
  * [`parseTime`](https://github.com/go-sql-driver/mysql?tab=readme-ov-file#parsetime)
* [OpenAPI Documentation](https://learn.openapis.org/)
* [Go database/sql tutorial](http://go-database-sql.org/index.html)
* [Illustrated guide to SQLX](https://jmoiron.github.io/sqlx/)
* [Why even use *DB.exec() or prepared statements in Golang?](https://stackoverflow.com/questions/50664648/why-even-use-db-exec-or-prepared-statements-in-golang)