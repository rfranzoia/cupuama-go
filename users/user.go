package users

// Users definition for users tables
type Users struct {
	Login    string
	Password string
	Person   Person
	Audit    Audit
}

var model *Users
