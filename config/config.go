package config

import "database/sql"

// AppConfig holds the application config
type AppConfig struct {
	UseCache bool
	SQLCache map[string]string
	DB       *sql.DB
}
