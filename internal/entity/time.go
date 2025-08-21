package entity

import "time"

// TimeEntry table
type TimeEntry struct {
	ID           int       `json:"id"`
	Hours        float32   `json:"hours"` // fixed: use float not time.Time
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ConsultantID int       `json:"consultantId"` // FK → consultants
	ProjectID    int       `json:"projectId"`    // FK → projects
	EntryDate    time.Time `json:"entryDate"`    // when it was logged
}

// TimeBudget is not a table itself, it's derived from queries
type TimeBudget struct {
	Assigned []TimeEntry `json:"assigned"`
	Entries  []TimeEntry `json:"entries"`
}
