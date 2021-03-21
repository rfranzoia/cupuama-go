package users

import "github.com/rfranzoia/cupuama-go/utils"

// Users definition for users tables
type Users struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Person
	utils.Audit
}

var model *Users
