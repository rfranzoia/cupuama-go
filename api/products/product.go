package products

import "github.com/rfranzoia/cupuama-go/api/users"

// Products definition for products tables
type Products struct {
	ID    int64
	Name  string
	Unit  string
	Audit users.Audit
}

var model *Products
