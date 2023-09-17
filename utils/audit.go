package utils

// Audit validation for data creation and manipulation
type Audit struct {
	Deleted     bool   `json:"deleted,omitempty"`
	DateCreated string `json:"dateCreated,omitempty" db:"date_created"`
	DateUpdated string `json:"dateUpdated,omitempty" db:"date_updated"`
}
