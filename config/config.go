package config

import "log"

// AppConfig holds the application config
type AppConfig struct {
	UseCache bool
	SQLCache map[string]string
	InfoLog  *log.Logger
}
