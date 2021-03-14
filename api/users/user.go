package users

import "github.com/rfranzoia/cupuama-go/utils"

// Users definition for users tables
type Users struct {
	Login    string
	Password string
	Person   Person
	Audit    utils.Audit
}

var model *Users
