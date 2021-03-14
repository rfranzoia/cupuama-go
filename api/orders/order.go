package orders

import "github.com/rfranzoia/cupuama-go/utils"

// Orders table definition
type Orders struct {
	ID         int64
	OrderDate  string
	TotalPrice float64
	Audit      utils.Audit
}

// OrderItemsStatus definition for the combination of Order + OrderItems + OrderStatus
type OrderItemsStatus struct {
	Order       Orders
	OrderStatus OrderStatus
	OrderItems  []OrderItems
}

var model *OrderItemsStatus
