package domain

// OrderItems table definition
type OrderItems struct {
	ID        int64
	Product   Products
	Fruit     Fruits
	Quantity  int64
	UnitPrice float64
}
