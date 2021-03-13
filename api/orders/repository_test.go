package orders

import (
	"log"
	"testing"

	"github.com/rfranzoia/cupuama-go/api/fruits"
	"github.com/rfranzoia/cupuama-go/api/products"
	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

var ois OrderItemsStatus
var testOIS OrderItemsStatus

func init() {
	var app config.AppConfig

	// loads all queries into the application config cache
	qc, err := utils.CreateSQLCache("../queries/*.sql")
	if err != nil {
		log.Fatal("cannot create queries cache")
	}

	app.SQLCache = qc
	app.UseCache = false

	NewOrderAPI(&app)

	setupOrderItemsStatus()
}

func setupOrderItemsStatus() {
	orderItem1 := OrderItems{
		Product: products.Products{
			ID: 3,
		},
		Fruit: fruits.Fruits{
			ID: 5,
		},
		Quantity:  10,
		UnitPrice: 7.5,
	}

	orderItem2 := OrderItems{
		Product: products.Products{
			ID: 1,
		},
		Fruit: fruits.Fruits{
			ID: 5,
		},
		Quantity:  10,
		UnitPrice: 7.5,
	}

	testOIS = OrderItemsStatus{
		Order: Orders{
			TotalPrice: 150.0,
		},
		OrderItems: []OrderItems{
			orderItem1,
			orderItem2,
		},
	}
}

func TestList(t *testing.T) {
	orders, err := ois.List(-1)
	if err != nil {
		t.Errorf("error while listing orders %v", err)

	} else if len(orders) == 0 {
		t.Errorf("expected size of orders list should be greater than zero")
	}
}

func TestGetFoundOrder(t *testing.T) {
	order, err := ois.Get(8)
	if err != nil {
		t.Errorf("error while retrieving the order %v", err)

	} else if order.Order.ID == 0 {
		t.Errorf("expected order was not found")
	}
}

func TestGetNotFoundOrder(t *testing.T) {
	order, err := ois.Get(999999999)
	if err != nil {
		t.Errorf("error while retrieving the order %v", err)

	} else if order.Order.ID != 0 {
		t.Errorf("there wasn't supposed to exist any order with the ID = 999999999")
	}
}

func TestCreate(t *testing.T) {

	order, err := ois.Create(testOIS)
	if err != nil {
		t.Errorf("error while creating an order %v", err)
	}

	if order.Order.ID == 0 {
		t.Errorf("order was not created properly")
	}

	record, err := ois.Get(order.Order.ID)
	if err != nil {
		t.Errorf("cannot retrive the order because it was not created properly")

	} else if record.Order.ID != order.Order.ID {
		t.Errorf("retrieved created order differs")
	}

}

func TestCreateOrderStatus(t *testing.T) {

	order, err := ois.Create(testOIS)
	if err != nil {
		t.Errorf("fail create order for TestCreateOrderStatus %v", err)
	}

	os := OrderStatus{
		Order:  order.Order,
		Status: OrderConfirmed,
	}

	err = ois.CreateOrderStatus(os, nil)
	if err != nil {
		t.Errorf("fail create order status for TestCreateOrderStatus %v", err)
	}
}

func TestUpdateOrder(t *testing.T) {
	orderItem := OrderItems{
		Product: products.Products{
			ID: 1,
		},
		Fruit: fruits.Fruits{
			ID: 5,
		},
		Quantity:  5,
		UnitPrice: 7,
	}

	oi := []OrderItems{
		orderItem,
	}

	order, err := ois.Create(testOIS)
	if err != nil {
		t.Errorf("(TestUpdateOrder) error creating an order %v", err)
	}

	err = ois.UpdateOrder(order.Order.ID, oi)
	if err != nil {
		t.Errorf("(TestUpdateOrder) error updating order %v", err)
	}

	order, err = ois.Get(order.Order.ID)
	if err != nil {
		t.Errorf("(TestUpdateOrder) after update order.Get() has failed %v", err)

	} else if len(order.OrderItems) != 1 {
		t.Errorf("(TestUpdateOrder) update order fail %v", err)
	}
}

func TestDeleteItemFromOrder(t *testing.T) {

	orderItem := OrderItems{
		Product: products.Products{
			ID: 3,
		},
		Fruit: fruits.Fruits{
			ID: 5,
		},
		Quantity:  10,
		UnitPrice: 7.5,
	}

	order, err := ois.Create(testOIS)
	if err != nil {
		t.Errorf("(TestDeleteItemFromOrder) error creating an order %v", err)
	}

	err = ois.DeleteOrderItems(order.Order.ID, orderItem)
	if err != nil {
		t.Errorf("(TestDeleteItemFromOrder) error removing order item(s) %v", err)
	}

	order, err = ois.Get(order.Order.ID)
	if err != nil {
		t.Errorf("(TestDeleteItemFromOrder) after delete OrderItem order.Get() has failed %v", err)

	} else if len(order.OrderItems) != 1 {
		t.Errorf("(TestDeleteItemFromOrder) delete order item has failed %v", err)
	}
}

func TestCancelOrderCreated(t *testing.T) {
	err := ois.CancelOrder(33)
	if err != nil {
		t.Errorf("(TestCancelOrderCreated) cancel has failed: %v", err)
	}
}
