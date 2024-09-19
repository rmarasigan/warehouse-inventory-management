package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Default values for the server configuration.
const (
	DefaultPort            string = "8080"
	DefaultIP              string = "0.0.0.0"
	DefaultApplicationName string = "Warehouse Inventory Management"
	DefaultDatabaseName    string = "wim_db"
	DefaultDatabaseUser    string = "root"

	DBNameKey     string = "db_name"
	DBUserKey     string = "db_user"
	DBPasswordKey string = "db_password"
)

type config struct {
	Application Application `yaml:"application,omitempty"`
	MySQL       MySQL       `yaml:"mysql,omitempty"`
}

// Application holds the application-specific configuration details.
type Application struct {
	Name string `yaml:"name,omitempty"`
	Host string `yaml:"host,omitempty"`
	Port string `yaml:"port,omitempty"`
}

// MySQL holds the MySQL database-specific configuration details.
type MySQL struct {
	DatabaseName string `yaml:"database_name,omitempty"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
}

// Load reads the configuration from a YAML file at the given path and loads it.
// Returns an error if the file is not found, cannot be read, or contains invalid
// YAML.
//
// Parameter:
//   - file: It should include the full path and filename.
func Load(file string) (*config, error) {
	if strings.TrimSpace(file) == "" {
		return nil, errors.New("file path not provided")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file at '%s': %w", file, err)
	}

	var cfg config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML from file '%s': %w", file, err)
	}

	// Initialize a new cache
	NewCache()

	// Cache the database config values
	SetCache(DBNameKey, cfg.DatabaseName())
	SetCache(DBUserKey, cfg.DatabaseUser())
	SetCache(DBPasswordKey, cfg.DatabasePassword())

	return &cfg, nil
}

// AppName returns the application name from the configuration. It uses
// default value if the application name is not provided in the configuration.
func (cfg config) AppName() string {
	if strings.TrimSpace(cfg.Application.Name) == "" {
		cfg.Application.Name = DefaultApplicationName
	}

	return cfg.Application.Name
}

// DatabaseName returns the MySQL database name from the configuration.
// It uses default value if the database name is not provided in the
// configuration.
func (cfg config) DatabaseName() string {
	if strings.TrimSpace(cfg.MySQL.DatabaseName) == "" {
		cfg.MySQL.DatabaseName = DefaultDatabaseName
	}

	return cfg.MySQL.DatabaseName
}

// DatabaseUser returns the MySQL database username. It uses default
// value if the username is not provided in the configuration.
func (cfg config) DatabaseUser() string {
	if strings.TrimSpace(cfg.MySQL.Username) == "" {
		cfg.MySQL.Username = DefaultDatabaseUser
	}

	return cfg.MySQL.Username
}

// DatabasePassword returns the MySQL database password from the configuration.
func (cfg config) DatabasePassword() string {
	return cfg.MySQL.Password
}

// ServerAddress returns the server address in the format "host:port".
// It uses default values if the host or port is not provided in the configuration.
func (cfg config) ServerAddress() string {
	host := cfg.Application.Host
	port := cfg.Application.Port

	switch {
	case host != "" && port != "":
		return fmt.Sprintf("%s:%s", host, port)

	case host == "" && port != "":
		return fmt.Sprintf("%s:%s", DefaultIP, port)

	case host != "" && port == "":
		return fmt.Sprintf("%s:%s", host, DefaultPort)

	default:
		return fmt.Sprintf("%s:%s", DefaultIP, DefaultPort)
	}
}
