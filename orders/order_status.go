package orders

// OrderStatus table definition
type OrderStatus struct {
	ID                int64
	Order             Orders
	Status            int64
	StatusChangeDate  string
	StatusDescription string
}

/*
	current valid status are:
		0 - order-created
		1 - order-confirmed
		2 - order-in-preparation
		3 - order-ready-for-delivery
		4 - order-dispatched
		5 - order-delivered
		9 - order-canceled

*/
