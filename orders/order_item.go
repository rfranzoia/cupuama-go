package orders

import (
	"github.com/rfranzoia/cupuama-go/fruits"
	"github.com/rfranzoia/cupuama-go/products"
)

// OrderItems table definition
type OrderItems struct {
	ID        int64
	Order     Orders
	Product   products.Products
	Fruit     fruits.Fruits
	Quantity  int64
	UnitPrice float64
}
