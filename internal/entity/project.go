package entity

import "time"

// Project table
type Project struct {
	ID                 int        `json:"id"`
	ProjectedStartDate time.Time  `json:"projectedStartDate"`
	StartDate          time.Time  `json:"startDate"`
	ProjectedEndDate   time.Time  `json:"projectedEndDate"`
	EndDate            time.Time  `json:"endDate"`
	Number             string     `json:"number"`
	Name               string     `json:"name"`
	TimeBudget         TimeBudget `json:"timeBudget"`
	ManagerID          int        `json:"managerId"` // FK → consultants
	Description        string     `json:"description"`
}

// Project tags → separate many-to-many table
type ProjectTag struct {
	ProjectID int    `json:"projectId"`
	Tag       string `json:"tag"`
}
