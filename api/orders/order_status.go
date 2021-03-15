package orders

// OrderStatus table definition
type OrderStatus struct {
	ID               int64
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

	OrderStatusMap = map[int64]OrderStatusType{
		0: OrderCreated, 1: OrderConfirmed, 2: OrderInPreparation, 3: OrderReadyForDelivery,
		4: OrderDispatched, 5: OrderDelivered, 9: OrderCanceled,
	}
)

func (os *OrderStatus) equals(other OrderStatus) bool {
	if os.ID == other.ID || os.Status.Value == other.Status.Value {
		return true
	}
	return false
}

func (ost *OrderStatusType) equals(other OrderStatusType) bool {
	if ost.Value == other.Value {
		return true
	}

	return false
}
