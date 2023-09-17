package database

import (
	"cupuama-go/logger"
	"github.com/jmoiron/sqlx"

	// no need to name this import
	_ "github.com/lib/pq"
)

var database *sqlx.DB
var err error

// GetConnection stabilish a connection with the database
func GetConnection() *sqlx.DB {
	logger.Log.Info("Connecting to database...")
	connection := "host=localhost port=5432 user=cupuama dbname=cupuama password=Cupu4m4. sslmode=disable"
	database, err = sqlx.Open("postgres", connection)

	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		logger.Log.Error("Error while connecting to the database " + err.Error())
		panic(err)
	}

	return database
}
