package users

import (
	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
)

type api struct {
	Service service
}

func NewAPI(app *config.AppConfig) *api {
	return &api{Service: NewUserService(app)}
}

// NewOrderAPI setups the configuration for orders
func (api *api) RegisterRouting(g *echo.Group) {

	gru := g.Group("/v2/users")
	gru.GET("", api.Service.List)
	gru.GET("/:login", api.Service.Get)
	gru.POST("", api.Service.Create)
	gru.PUT("/:login", api.Service.Update)
	gru.DELETE("/:login", api.Service.Delete)

}

func (api *api) RegisterLoginRouting(g *echo.Group) {

	// register login separatedly
	g.POST("/v2/login", api.Service.Login)

}
