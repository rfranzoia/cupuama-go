package main

import (
	"cupuama-go/api"
	"cupuama-go/config"
	"cupuama-go/database"
	"cupuama-go/logger"
	"cupuama-go/utils"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	logger.Log.Info("Stating application server")
	db := database.GetConnection()
	defer db.Close()

	var app config.AppConfig

	// loads all queries into the application config cache
	qc, err := utils.CreateSQLCache()
	if err != nil {
		logger.Log.Fatal("cannot create queries cache")
	}

	app.SQLCache = qc
	app.UseCache = false
	app.DB = db

	e := ServerConfig(&app)
	e.Logger.Fatal(e.Start(":8080"))
}

func ServerConfig(app *config.AppConfig) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())

	// setup routes here
	g := e.Group("/cupuama-go/api")

	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${host}${path} ${latency_human}\n",
	}))

	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodHead},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	api.RegisterAPIRoutes(g, app)

	return e
}
