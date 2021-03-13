package main

import (
	"log"

	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/database"
	"github.com/rfranzoia/cupuama-go/utils"
)

func main() {
	defer database.GetConnection().Close()

	var app config.AppConfig

	// loads all queries into the application config cache
	qc, err := utils.CreateSQLCache()
	if err != nil {
		log.Fatal("cannot create queries cache")
	}

	app.SQLCache = qc
	app.UseCache = false

}
