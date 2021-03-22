package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rfranzoia/cupuama-go/api"
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
