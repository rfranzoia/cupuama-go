package utils

// Audit validation for data creation and manipulation
type Audit struct {
	DateCreated string `json:"dateCreated"`
	Deleted     bool   `json:"deleted"`
	DateUpdated string `json:"dateUpdated"`
}
