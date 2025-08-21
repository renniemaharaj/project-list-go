package entity

import "time"

// Status table
type Status struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateCreated time.Time `json:"dateCreated"`
}

// Project history mapping of project to status
type StatusHistory struct {
	ProjectID int       `json:"projectId"`
	StatusID  int       `json:"statusId"`
	ChangedAt time.Time `json:"changedAt"`
}
