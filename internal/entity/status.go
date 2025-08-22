package entity

import "time"

// Status struct
type Status struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ProjectID   int       `json:"projectId"`
	Description string    `json:"description"`
	DateCreated time.Time `json:"dateCreated"`
}
