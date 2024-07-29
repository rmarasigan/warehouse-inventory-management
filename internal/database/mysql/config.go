package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	cfg      config
	database *sql.DB
	mysql    string = "mysql"
)

type config struct {
	User         string
	Password     string
	DatabaseName string
}

// SetUser sets the database configuration username.
func SetUser(username string) {
	cfg.User = username
}

// SetPassword sets the database configuration password.
func SetPassword(password string) {
	cfg.Password = password
}

// SetDatabaseName sets the database configuration database name.
func SetDatabaseName(dbname string) {
	cfg.DatabaseName = dbname
}

// Connect opens a connection to MySQL using the specified DSN (data source name)
// with a default configuration for MaxOpenConns, MaxIdleConns and MaxLifetime.
// If the DSN username is not defined, it defaults to 'root'.
func Connect() {
	if database == nil {
		if strings.TrimSpace(cfg.User) == "" {
			cfg.User = "root"
		}

		if strings.TrimSpace(cfg.Password) == "" {
			panic("mysql password is not configured")
		}

		if strings.TrimSpace(cfg.DatabaseName) == "" {
			panic("mysql database name is not configured")
		}

		// Data Source Name
		var dsn = fmt.Sprintf("%s:%s@/%s", cfg.User, cfg.Password, cfg.DatabaseName)

		db, err := sql.Open(mysql, dsn)
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
	}
}
