package users

// Person basic person information
type Person struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dateOfBirth"`
}
