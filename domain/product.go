package domain

import "cupuama-go/utils"

// Products definition for products tables
type Products struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
	utils.Audit
}
