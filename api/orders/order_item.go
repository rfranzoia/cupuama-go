package orders

import (
	"github.com/rfranzoia/cupuama-go/api/fruits"
	"github.com/rfranzoia/cupuama-go/api/products"
)

// OrderItems table definition
type OrderItems struct {
	ID        int64
	Product   products.Products
	Fruit     fruits.Fruits
	Quantity  int64
	UnitPrice float64
}
