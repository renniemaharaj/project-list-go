package entity

// Consultant table
type Consultant struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	// Roles          []string `json:"roles"` // stored in separate table
	ProfilePicture string `json:"profilePicture"`
}
