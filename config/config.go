package config

import (
	"github.com/jmoiron/sqlx"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache bool
	SQLCache map[string]string
	DB       *sqlx.DB
}
