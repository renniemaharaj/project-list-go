package entity

import "time"

// ProjectStatus struct
type ProjectStatus struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	ProjectID    int       `json:"projectID"`
	ConsultantID int       `json:"consultantID"`
	Description  string    `json:"description"`
	DateCreated  time.Time `json:"dateCreated"`
}
