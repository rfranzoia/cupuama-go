package orders

import (
	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
)

type api struct {
	Service service
}

func NewAPI(app *config.AppConfig) *api {
	return &api{Service: NewOrderService(app)}
}

// NewOrderAPI setups the configuration for orders	+
func (api *api) RegisterRouting(g *echo.Group) {

	gro := g.Group("/v2/orders")
	gro.GET("", api.Service.List)
	gro.GET("/:id", api.Service.Get)
	gro.POST("", api.Service.Create)
	gro.PUT("/:id/status/:status", api.Service.ChangeOrderStatus)
	gro.PUT("/:id/cancel", api.Service.CancelOrder)

	// gu.PUT("/:id", api.Service.Update)
	// gu.DELETE("/:id", api.Service.Delete)

}
