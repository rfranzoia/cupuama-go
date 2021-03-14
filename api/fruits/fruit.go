package fruits

import "github.com/rfranzoia/cupuama-go/utils"

// Fruits definition for fruits tables
type Fruits struct {
	ID       int64
	Name     string
	Initials string
	Harvest  string
	Audit    utils.Audit
}

var model *Fruits
