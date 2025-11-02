# Warehouse Inventory Management

A simple Warehouse Inventory Management System that is designed to help you manage and track inventory within a warehouse.

> [!WARNING]
>
> This project is for learning purposes and should not be used in production environments.

### Running the Application

To perform the **initial setup** and get the environment running, build the binary and run it with the `--db=init` flag:
```bash
dev@dev:~/warehouse-inventory-management$ go build
dev@dev:~/warehouse-inventory-management$ ./warehouse-inventory-management --db=init
```

You should see output similar to the following:
```bash
2025-11-02 11:03:47 AM | INFO   Initializing database...
2025-11-02 11:03:47 AM | OK     Successfully created inventory_management database...
2025-11-02 11:03:47 AM | OK     Successfully created role table...
2025-11-02 11:03:47 AM | OK     Successfully created users table...
2025-11-02 11:03:47 AM | OK     Successfully created unit_of_measurement table...
2025-11-02 11:03:47 AM | OK     Successfully created storage table...
2025-11-02 11:03:47 AM | OK     Successfully created currency table...
2025-11-02 11:03:47 AM | OK     Successfully created item table...
2025-11-02 11:03:47 AM | OK     Successfully created transactions table...
2025-11-02 11:03:47 AM | OK     Successfully created orderline table...
2025-11-02 11:03:47 AM | OK     Successfully imported items to role table...
2025-11-02 11:03:47 AM | OK     Successfully imported items to unit_of_measurement table...
2025-11-02 11:03:47 AM | OK     Successfully imported items to currency table...
2025-11-02 11:03:47 AM | INFO   Establishing MySQL connection...
2025-11-02 11:03:47 AM | OK     MySQL connection established.
2025-11-02 11:03:47 AM | INFO   Initializing Warehouse Inventory Management Server at 0.0.0.0:8080
                  ______   _____
   ___  ___  ___ /  /    \/     \
  /  / /  / /  //  /  __   __   /
 /  / /  / /  //  /  / /  / /  /
/  /_/  /_/  //  /  / /  / /  /
\_____,_____//__/__/ /__/ /__/
```

Once the database and tables are initialized, you can start the server normally without the `--db=init` flag:
```bash
dev@dev:~/warehouse-inventory-management$ go build
dev@dev:~/warehouse-inventory-management$ ./warehouse-inventory-management
```

This will:
* Connect to the existing database.
* Load cached configurations.
* Start the HTTP server on port 8080.

Expected output:
```bash
2025-11-02 12:49:16 PM | INFO   Establishing MySQL connection...
2025-11-02 12:49:16 PM | OK     MySQL connection established.
2025-11-02 12:49:16 PM | INFO   Initializing Warehouse Inventory Management Server at 0.0.0.0:8080
                  ______   _____
   ___  ___  ___ /  /    \/     \
  /  / /  / /  //  /  __   __   /
 /  / /  / /  //  /  / /  / /  /
/  /_/  /_/  //  /  / /  / /  /
\_____,_____//__/__/ /__/ /__/
```

> [!TIP]
>
> If you ever need to reset the database schema (e.g., for local testing), re-run the command with the `--db=init` flag: `./warehouse-inventory-management --db=init`

## Requirements
* **Go**: v1.24
* **MySQL**: 8.0.43

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

If the package fails to install or work on modern Python versions due to dependency/build issues (like `PyYAML`), try this workaround:
```bash
# Creates a new virtual environment named .venv in your current directory.
# This isolates your Python packages from the system-wide ones, so dependencies don't conflict.
$ python3 -m venv .venv

# Activates the virtual environment .venv.
# Your shell session now uses the Python and pip inside .venv, keeping your project dependencies separate.
$ source .venv/bin/activate

# Installs essential build tools inside the virtual environment.
$ pip install "cython<3.0.0" wheel setuptools

# Installs PyYAML version 5.4.1 specifically, because newer versions might cause build problems with this package.
# The --no-build-isolation flag tells pip to not create a separate isolated environment for building, so it uses the current environment and its installed packages (like the pinned cython and setuptools), which helps avoid some build errors.
$ pip install "pyyaml==5.4.1" --no-build-isolation

# Installs the openapi2jsonschema package.
$ pip install openapi2jsonschema

# Verify if installed correctly and is runnable.
$ openapi2jsonschema --help

# Exits the virtual environment, returning your shell session to the normal system Python environment.
$ deactivate
```

* Github Issue: [Fix installation on modern Python versions](https://github.com/instrumenta/openapi2jsonschema/issues/70)

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