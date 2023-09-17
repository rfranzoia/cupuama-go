package utils

// Person basic person information
type Person struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	DateOfBirth string `json:"dateOfBirth,omitempty"`
}