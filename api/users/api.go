package users

import (
	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
)

type api struct {
	Service service
}

func NewAPI(app *config.AppConfig) *api {
	return &api{Service: New(app)}
}

// NewOrderAPI setups the configuration for orders
func (api *api) RegisterRouting(g *echo.Group) {

	gu := g.Group("/v2/users")
	gu.GET("", api.Service.List)
	gu.GET("/:login", api.Service.Get)
	gu.POST("", api.Service.Create)
	gu.PUT("/:login", api.Service.Update)
	gu.DELETE("/:login", api.Service.Delete)

}
