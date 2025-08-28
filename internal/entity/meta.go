package entity

// ProjectMeta struct encapsulates the necessary meta data of a poject
type ProjectMeta struct {
	Manager       Consultant      `json:"manager"`
	TimeEntries   []TimeEntry     `json:"timeEntries"`
	StatusHistory []ProjectStatus `json:"statusHistory"`
	Consultants   []Consultant    `json:"consultants"`
}
