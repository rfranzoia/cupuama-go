package fruits

import "github.com/rfranzoia/cupuama-go/utils"

// Fruits definition for fruits tables
type Fruits struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Initials string `json:"initials"`
	Harvest  string `json:"harvest"`
	utils.Audit
}

var model *Fruits
