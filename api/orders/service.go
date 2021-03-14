package orders

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

type service struct {
}

var app *config.AppConfig

// NewOrderAPI setups the configuration for orders
func NewOrderService(a *config.AppConfig) service {
	app = a
	return service{}
}

// List list all orders
func (s *service) List(c echo.Context) error {

	list, err := model.List(-1)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error listing Orders"),
			Value:   err.Error(),
		})
	}
	defer c.Request().Body.Close()

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: list,
	})

}

// Get retrive an order by its ID
func (s *service) Get(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	u, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error searching Order %d", id),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: u,
	})

}

// Create creates an order
func Create(ois OrderItemsStatus) {
	ois, err := model.Create(ois)
	if err != nil {
		// do something later
	}
}

// CreateOrderStatus creates a new status for an order
func CreateOrderStatus(os OrderStatus) {
	err := model.CreateOrderStatus(os, nil)
	if err != nil {
		// do something later
	}
}
