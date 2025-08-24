package entity

import "time"

// Project table
type Project struct {
	ID                 int       `json:"ID"`
	ProjectedStartDate time.Time `json:"projectedStartDate"`
	StartDate          time.Time `json:"startDate"`
	ProjectedEndDate   time.Time `json:"projectedEndDate"`
	EndDate            time.Time `json:"endDate"`
	Number             string    `json:"number"`
	Name               string    `json:"name"`
	ManagerID          int       `json:"managerID"` // FK → consultants
	Description        string    `json:"description"`
}

// Project tags → separate many-to-many table
type ProjectTag struct {
	ID        int    `json:"ID"`
	ProjectID int    `json:"projectID"`
	Tag       string `json:"tag"`
}
