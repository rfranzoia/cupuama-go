package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/api/fruits"
	"github.com/rfranzoia/cupuama-go/api/orders"
	"github.com/rfranzoia/cupuama-go/api/products"
	"github.com/rfranzoia/cupuama-go/api/users"
	"github.com/rfranzoia/cupuama-go/config"
)

func RegisterAPIRoutes(g *echo.Group, app *config.AppConfig) {

	u := users.NewAPI(app)
	u.RegisterRouting(g)

	p := products.NewAPI(app)
	p.RegisterRouting(g)

	f := fruits.NewAPI(app)
	f.RegisterRouting(g)

	o := orders.NewAPI(app)
	o.RegisterRouting(g)

}
