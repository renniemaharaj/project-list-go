package entity

import "time"

// TimeEntry table
type TimeEntry struct {
	ID           int       `json:"id"`
	Hours        float32   `json:"hours"` // fixed: use float not time.Time
	Title        string    `json:"title"`
	ConsultantID int       `json:"consultantID"` // FK → consultants
	Description  string    `json:"description"`
	ProjectID    int       `json:"projectID"` // FK → projects
	Type         string    `json:"type"`      // Debit or Credit
	EntryDate    time.Time `json:"entryDate"` // when it was logged
}
