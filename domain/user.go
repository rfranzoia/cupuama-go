package domain

import "cupuama-go/utils"

// Users definition for users tables
type Users struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	utils.Person
	utils.Audit
}
