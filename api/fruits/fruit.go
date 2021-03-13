package fruits

import "github.com/rfranzoia/cupuama-go/api/users"

// Fruits definition for fruits tables
type Fruits struct {
	ID       int64
	Name     string
	Initials string
	Harvest  string
	Audit    users.Audit
}

var model *Fruits
