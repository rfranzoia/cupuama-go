package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/api/fruits"
	"github.com/rfranzoia/cupuama-go/api/products"
	"github.com/rfranzoia/cupuama-go/api/users"
	"github.com/rfranzoia/cupuama-go/config"
)

func RegisterAPIRoutes(g *echo.Group, app *config.AppConfig) {

	u := users.NewAPI(app)
	u.RegisterRouting(g)

	fruits.RegisterRouting(g, app)
	products.RegisterRouting(g, app)

}
