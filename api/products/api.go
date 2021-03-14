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

	grp := g.Group("/v2/products")
	grp.GET("", api.Service.List)
	grp.GET("/:id", api.Service.Get)
	grp.POST("", api.Service.Create)
	grp.PUT("/:id", api.Service.Update)
	grp.DELETE("/:id", api.Service.Delete)

}
