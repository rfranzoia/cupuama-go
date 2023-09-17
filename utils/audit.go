package utils

// Audit validation for data creation and manipulation
type Audit struct {
	Deleted     bool   `json:"deleted"`
	DateCreated string `json:"dateCreated,omitempty"`
	DateUpdated string `json:"dateUpdated,omitempty"`
}
