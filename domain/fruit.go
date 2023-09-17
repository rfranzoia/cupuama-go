package domain

import "cupuama-go/utils"

// Fruits definition for fruits tables
type Fruits struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	Initials string `json:"initials"`
	Harvest  string `json:"harvest,omitempty"`
	utils.Audit
}
