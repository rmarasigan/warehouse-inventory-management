# Warehouse Inventory Management

## API Validation
API validation schemas are generated from the [`openapi.yaml`](openapi.yaml) using the tool [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema).

**Installation**

To install [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema), run:
```bash
pip install openapi2jsonschema
```

**Generating JSON Schemas**

To generate the JSON schemas from the [`openapi.yaml`](openapi.yaml):
```bash
openapi2jsonschema openapi.yaml -o ./api/schema/validator/spec
```

This will create a set of JSON schemas in `./api/schema/validator/spec` directory.

## Reference
* [Go database/sql tutorial](http://go-database-sql.org/index.html)
* [Why even use *DB.exec() or prepared statements in Golang?](https://stackoverflow.com/questions/50664648/why-even-use-db-exec-or-prepared-statements-in-golang)