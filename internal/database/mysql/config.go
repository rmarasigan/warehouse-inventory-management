package mysql

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rmarasigan/warehouse-inventory-management/internal/app/config"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/trail"
)

var (
	database *sqlx.DB
	mysql    string = "mysql"
)

// Connect opens a connection to MySQL using the specified DSN (data source name)
// with a default configuration for MaxOpenConns, MaxIdleConns and MaxLifetime.
// If the DSN username is not defined, it defaults to 'root'.
func Connect() {
	if database == nil {
		trail.Info("Establishing MySQL connection...")

		dbname, ok := config.GetCache(config.DBNameKey).(string)
		if !ok {
			panic(errors.New("expected string for database name but got a different type"))
		}

		dbuser, ok := config.GetCache(config.DBUserKey).(string)
		if !ok {
			panic(errors.New("expected string for database user but got a different type"))
		}

		dbpassword, ok := config.GetCache(config.DBPasswordKey).(string)
		if !ok {
			panic(errors.New("expected string for database password but got a different type"))
		}

		// Data Source Name
		var dsn = fmt.Sprintf("%s:%s@/%s?parseTime=true", dbuser, dbpassword, dbname)

		// Connects to the database and attempts a ping.
		db, err := sqlx.Connect(mysql, dsn)
		if err != nil {
			panic(err)
		}

		// Limit the number of connection used by the application.
		db.SetMaxOpenConns(5)

		// Recommended to be set the same as SetMaxOpenConns.
		db.SetMaxIdleConns(5)

		// Ensure connections are closed by the driver safely before
		// connection is closed by MySQL.
		db.SetConnMaxLifetime(time.Second * 270)

		database = db
		trail.OK("MySQL connection established.")
	}
}

// Close closes the database connection.
func Close() {
	if database != nil {
		err := database.Close()
		if err != nil {
			panic(err)
		}

		database = nil
		trail.Info("MySQL connection closed.")
	}
}
