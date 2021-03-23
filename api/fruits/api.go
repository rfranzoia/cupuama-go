package fruits

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

type api struct {
	Service service
}

func NewAPI(app *config.AppConfig) *api {
	return &api{Service: NewFruitService(app)}
}

// NewOrderAPI setups the configuration for orders
func (api *api) RegisterRouting(g *echo.Group) {

	grf := g.Group("/v2/fruits")

	grf.GET("", api.Service.List)
	grf.GET("/:id", api.Service.Get)
	grf.POST("", api.Service.Create)
	grf.PUT("/:id", api.Service.Update)
	grf.DELETE("/:id", api.Service.Delete)

	//grf.OPTIONS("/:id", options)
}

func options(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.MessageJSON{
		Message: "",
	})
}
