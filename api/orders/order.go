package orders

import "github.com/rfranzoia/cupuama-go/api/users"

// Orders table definition
type Orders struct {
	ID         int64
	OrderDate  string
	TotalPrice float64
	Audit      users.Audit
}

// OrderItemsStatus definition for the combination of Order + OrderItems + OrderStatus
type OrderItemsStatus struct {
	Order       Orders
	OrderStatus OrderStatus
	OrderItems  []OrderItems
}

var model *OrderItemsStatus
