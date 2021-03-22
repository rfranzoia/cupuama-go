package utils

// Audit validation for data creation and manipulation
type Audit struct {
	DateCreated string `json:"dateCreated,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
	DateUpdated string `json:"dateUpdated,omitempty"`
}
