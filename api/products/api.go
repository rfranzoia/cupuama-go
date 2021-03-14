package products

import (
	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
)

type api struct {
	Service service
}

func NewAPI(app *config.AppConfig) *api {
	return &api{Service: NewProductService(app)}
}

// NewOrderAPI setups the configuration for orders
func (api *api) RegisterRouting(g *echo.Group) {

	gu := g.Group("/v2/products")
	gu.GET("", api.Service.List)
	gu.GET("/:id", api.Service.Get)
	gu.POST("", api.Service.Create)
	gu.PUT("/:id", api.Service.Update)
	gu.DELETE("/:id", api.Service.Delete)

}
