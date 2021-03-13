package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/api/users"
	"github.com/rfranzoia/cupuama-go/config"
)

func RegisterAPIRoutes(g *echo.Group, app *config.AppConfig) {
	users.RegisterRouting(g, app)
}
