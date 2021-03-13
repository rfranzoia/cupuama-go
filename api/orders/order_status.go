package orders

// OrderStatus table definition
type OrderStatus struct {
	ID               int64
	Order            Orders
	Status           OrderStatusType
	StatusChangeDate string
}

type OrderStatusType struct {
	Value       int64
	Description string
}

var (
	OrderCreated          = OrderStatusType{Value: 0, Description: "order-created"}
	OrderConfirmed        = OrderStatusType{Value: 1, Description: "order-confirmed"}
	OrderInPreparation    = OrderStatusType{Value: 2, Description: "order-in-preparation"}
	OrderReadyForDelivery = OrderStatusType{Value: 3, Description: "order-ready-for-delivery"}
	OrderDispatched       = OrderStatusType{Value: 4, Description: "order-dispatched"}
	OrderDelivered        = OrderStatusType{Value: 5, Description: "order-delivered"}
	OrderCanceled         = OrderStatusType{Value: 9, Description: "order-canceled"}
)
