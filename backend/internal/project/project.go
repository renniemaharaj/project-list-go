package project

import (
	"time"
)

type Consultant struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`

	// "Administrator", "Manager" "Consultant"
	Roles          []string `json:"roles"`
	ProfilePicture string   `json:"profilePicture"`
}

type TimeEntry struct {
	Hours       time.Time  `json:"hours"`       // The amount of hours
	Title       string     `json:"title"`       // Entry title
	Description string     `json:"description"` // Entry reason
	Consultant  Consultant `json:"consultant"`  // Attached consultant
}

type TimeBudget struct {
	Assigned []TimeEntry `json:"assigned"` // Time debiting
	Entries  []TimeEntry `json:"entries"`  // Time crediting
}

type Status struct {
	Title       string `json:"title"`       // Active, Inactive, Idle
	Description string `json:"description"` // Status description
}

type Project struct {
	ID int `json:"id"` // The project ID
	// The project's projected start date
	ProjectedStartDate time.Time `json:"projectedStartDate"`
	// Project's actual start date
	StartDate time.Time `json:"startDate"`
	// The project's end date
	ProjectedEndDate time.Time `json:"projectedEndDate"`
	// The project's actual end date
	EndDate time.Time `json:"endDate"`
	// Higher level project number
	Number string `json:"number"`
	// The assigned project name
	Name string `json:"name"`
	// Hours assigned to project
	HoursAssigned float32 `json:"hoursAssigned"`
	// Time budgeting struct
	TimeBudget TimeBudget `json:"timeBudget"`
	// The project manager
	Manager Consultant `json:"manager"` // The project's manager
	// Project tags: configuration, optimization, training etc
	DescriptiveTags []string `json:"descriptiveTags"`
	// The project's description
	Description string `json:"description"`
	// Project's status history
	History []Status `json:"history"`
}
