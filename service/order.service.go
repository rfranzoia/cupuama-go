package service

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"cupuama-go/repository"
	"cupuama-go/utils"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type OrderService struct {
	app        *config.AppConfig
	repository repository.OrderRepository
}

// NewOrderAPI setups the configuration for orders
func NewOrderService(a *config.AppConfig) OrderService {
	return OrderService{
		app:        a,
		repository: repository.NewOrderRepository(a),
	}
}

// List list all orders
func (os *OrderService) List(c echo.Context) error {

	list, err := os.repository.List()
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: "error listing Orders",
			Value:   err.Error(),
		})
	}
	defer c.Request().Body.Close()

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: list,
	})

}

// Get retrive an order by its ID
func (os *OrderService) Get(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	order, err := os.repository.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error searching Order %d", id),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: order,
	})

}

// Create creates an order
func (os *OrderService) Create(c echo.Context) error {

	order := new(domain.OrderItemsStatus)

	if err := c.Bind(order); err != nil {
		logger.Log.Info("(CreateOrder:Bind) " + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating Order",
			Value:   err.Error(),
		})
	}

	id, err := os.repository.Create(order)
	if err != nil || id <= 0 {
		logger.Log.Info("(CreateOrder:Create)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating Order",
			Value:   err.Error(),
		})
	}

	o, err := os.repository.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error searching Order %d", id),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: o,
	})
}

// CreateOrderStatus creates a new status for an order
func (os *OrderService) ChangeOrderStatus(c echo.Context) error {

	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	statusID, _ := strconv.ParseInt(c.Param("status"), 10, 64)

	Status, isValid := domain.OrderStatusMap[statusID]
	if !isValid {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: "status informed is invalid",
			Value:   statusID,
		})
	}

	if Status.Equals(domain.OrderCanceled) {
		return os.CancelOrder(c)
	}

	// if changing to OrderCanceled, then just call the appropriate CancelOrder method
	status := domain.OrderStatus{
		Status: Status,
	}

	if err := os.repository.CreateOrderStatus(orderID, status, nil); err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error adding status ´%s´ to order %d", Status.Description, orderID),
			Value:   err.Error(),
		})
	}

	return os.Get(c)

}

func (os *OrderService) CancelOrder(c echo.Context) error {
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := os.repository.CancelOrder(orderID); err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("error canceling order %d", orderID),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Message: "Order successfully Canceled",
		Value:   "",
	})
}

func (os *OrderService) DeleteOrderItems(c echo.Context) error {
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var ois domain.OrderItemsStatus

	if err := c.Bind(&ois); err != nil {
		logger.Log.Info("(DeleteOrderItems:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error removing order items Order",
			Value:   err.Error(),
		})
	}

	if err := os.repository.DeleteOrderItems(orderID, ois.OrderItems); err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error removing order items Order",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Message: "Order Items successfully deleted",
		Value:   "",
	})
}

func (os *OrderService) UpdateOrder(c echo.Context) error {
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var order domain.OrderItemsStatus

	if err := c.Bind(&order); err != nil {
		logger.Log.Info("(UpdateOrder:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error updating order",
			Value:   err.Error(),
		})
	}

	if err := os.repository.UpdateOrder(orderID, order.OrderItems); err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error updating order",
			Value:   err.Error(),
		})
	}

	return os.Get(c)
}
