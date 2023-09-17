package api

import (
	"cupuama-go/config"
	"cupuama-go/service"
	"github.com/labstack/echo"
)

type OrderApi struct {
	service service.OrderService
}

func NewOrderAPI(app *config.AppConfig) *OrderApi {
	return &OrderApi{service: service.NewOrderService(app)}
}

// NewOrderAPI setups the configuration for orders	+
func (api *OrderApi) RegisterRouting(g *echo.Group) {

	gro := g.Group("/v2/orders")
	gro.GET("", api.service.List)
	gro.GET("/:id", api.service.Get)
	gro.POST("", api.service.Create)
	gro.PUT("/:id/status/:status", api.service.ChangeOrderStatus)
	gro.PUT("/:id/cancel", api.service.CancelOrder)
	gro.DELETE("/:id/items", api.service.DeleteOrderItems)
	gro.PUT("/:id", api.service.UpdateOrder)

}
