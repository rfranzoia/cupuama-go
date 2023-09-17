package api

import (
	"cupuama-go/config"

	"github.com/labstack/echo"
)

func RegisterAPIRoutes(g *echo.Group, app *config.AppConfig) {

	NewUserAPI(app).RegisterRouting(g)
	NewProductAPI(app).RegisterRouting(g)
	NewFruitAPI(app).RegisterRouting(g)

	o := NewOrderAPI(app)
	o.RegisterRouting(g)

}
