package entity

// Consultant table
type Consultant struct {
	ID        int    `json:"ID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	// Roles          []string `json:"roles"` // stored in separate table
	ProfilePicture string `json:"profilePicture"`
}

// The consultant tole struct
type ConsultantRole struct {
	ID           int    `json:"ID"`
	ConsultantID int    `json:"consultantID"`
	Role         string `json:"role"`
}

// A consultant project struct
type ProjectConsultant struct {
	ID           int    `json:"ID"`
	ProjectID    int    `json:"projectID"`
	ConsultantID int    `json:"consultantID"`
	Role         string `json:"role"`
}
