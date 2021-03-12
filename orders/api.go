package orders

import "github.com/rfranzoia/cupuama-go/config"

var app *config.AppConfig

// NewOrderAPI setups the configuration for orders
func NewOrderAPI(a *config.AppConfig) {
	app = a
}

// List list all orders
func List() []OrderItemsStatus {
	list, err := model.List(-1)
	if err != nil {
		return nil
	}
	return list
}

// Get retrive an order by its ID
func Get(orderID int64) OrderItemsStatus {
	list, err := model.List(orderID)
	if err != nil {
		return OrderItemsStatus{}
	}
	return list[0]
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
