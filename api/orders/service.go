package orders

import (
	"fmt"
	"log"
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

	list, err := model.List()
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

	o, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error searching Order %d", id),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: o,
	})

}

// Create creates an order
func (s *service) Create(c echo.Context) error {

	order := new(OrderItemsStatus)

	if err := c.Bind(order); err != nil {
		log.Println("(CreateFruit:Bind)", err)
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating Order",
			Value:   err.Error(),
		})
	}

	id, err := model.Create(order)
	if err != nil || id <= 0 {
		log.Println("(CreateFruit:Create)", err)
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating Order",
			Value:   err.Error(),
		})
	}

	o, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error searching created Order %d", id),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: o,
	})
}

// CreateOrderStatus creates a new status for an order
func CreateOrderStatus(orderID int64, os OrderStatus) {
	err := model.CreateOrderStatus(orderID, os, nil)
	if err != nil {
		// do something later
	}
}