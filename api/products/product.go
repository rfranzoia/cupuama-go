package products

import "github.com/rfranzoia/cupuama-go/utils"

// Products definition for products tables
type Products struct {
	ID    int64
	Name  string
	Unit  string
	Audit utils.Audit
}

var model *Products
