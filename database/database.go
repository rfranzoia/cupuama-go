package database

import (
	"database/sql"

	// no need to name this import
	_ "github.com/lib/pq"
)

var database *sql.DB
var err error

// GetConnection stabilish a connection with the database
func GetConnection() *sql.DB {

	connection := "host=localhost port=5432 user=cupuama dbname=cupuama password=Cupu4m4. sslmode=disable"
	database, err = sql.Open("postgres", connection)

	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	//fmt.Println("Successfully connected!")

	return database
}
