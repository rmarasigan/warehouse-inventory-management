package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rmarasigan/warehouse-inventory-management/internal/app/config"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/convert"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/trail"
)

func Initialize() {
	log.Init()
	defer log.Panic()

	ctx := context.Background()

	// Load application configuration from YAML file.
	cfg, err := config.Load("wim-config.yaml")
	if err != nil {
		trail.Warn("failed to load app configuration")
		panic(err)
	}

	// Build DSN (Data Source Name) to connect to MySQL server without selecting a database yet.
	// We leave the database empty (trailing '/') so we can create it if it doesn't exist.
	//
	// DSN Options:
	//  - charset=utf8mb4: Ensures proper UTF-8 support, including emojis and special characters.
	//  - parseTime=True: Automatically converts MySQL DATE, DATETIME, and TIMESTAMP columns into Go's time.Time.
	//  - loc=UTC: Interprets MySQL datetime values in UTC for consistency across servers and environments.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.DatabaseUser(),
		cfg.DatabasePassword(),
		cfg.DatabaseHost(),
		cfg.DatabasePort(),
	)

	// Connect to the MySQL server.
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		trail.Warn("failed to create connection to database")
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err, "failed to close database connection")
		}
	}()

	// Create the database if it doesn't exist.
	err = database(ctx, db, cfg.DatabaseName())
	if err != nil {
		panic(err)
	}

	// Create all tables in the correct order to satisfy foreign key constraints.
	for _, tableName := range tablesOrder {
		query := databaseTables[tableName]

		// Each table is created if it doesn't already exist.
		err := tables(ctx, db, tableName, query)
		if err != nil {
			panic(err)
		}
	}

	// Insert roles from configuration (array of strings)
	if len(cfg.Role()) > 0 {
		for _, roleName := range cfg.Role() {
			role := schema.Role{Name: roleName}

			err := insert(db, role, roleInsert, "role")
			if err != nil {
				panic(err)
			}
		}

		trail.OK("Successfully imported items to role table...")
	}

	// Insert UnitOfMeasurement from config.
	if len(cfg.UnitOfMeasurement()) > 0 {
		list := convert.SchemaList(cfg.UnitOfMeasurement(),
			func(uom config.UnitOfMeasurement) schema.UOM {
				return schema.UOM{
					Code: uom.Code,
					Name: uom.Name,
				}
			},
		)

		err := insertList(db, list, uomInsert, "unit_of_measurement")
		if err != nil {
			panic(err)
		}
	}

	// Insert Currency from config.
	if len(cfg.Currency()) > 0 {
		list := convert.SchemaList(cfg.Currency(),
			func(currency config.Currency) schema.Currency {
				return schema.Currency{
					Code:   currency.Code,
					Symbol: currency.Symbol,
					Active: currency.Active,
				}
			},
		)

		err := insertList(db, list, currencyInsert, "currency")
		if err != nil {
			panic(err)
		}
	}
}
