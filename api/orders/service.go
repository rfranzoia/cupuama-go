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

	return s.Get(c)
}

// CreateOrderStatus creates a new status for an order
func (s *service) ChangeOrderStatus(c echo.Context) error {

	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	statusID, _ := strconv.ParseInt(c.Param("status"), 10, 64)

	Status, isValid := OrderStatusMap[statusID]
	if !isValid {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("status informed is invalid"),
			Value:   statusID,
		})
	}

	if Status.equals(OrderCanceled) {
		return s.CancelOrder(c)
	}

	// if changing to OrderCanceled, then just call the appropriate CancelOrder method
	os := OrderStatus{
		Status: Status,
	}

	if err := model.CreateOrderStatus(orderID, os, nil); err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error adding status ´%s´ to order %d", Status.Description, orderID),
			Value:   err.Error(),
		})
	}

	return s.Get(c)

}

func (s *service) CancelOrder(c echo.Context) error {
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := model.CancelOrder(orderID); err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error canceling order %d", orderID),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Message: "Order successfully Canceled",
	})
}
