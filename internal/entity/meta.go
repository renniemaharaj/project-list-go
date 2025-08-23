package entity

// ProjectMetaData struct encapsulates the necessary meta data of a poject
type ProjectMetaData struct {
	Manager       Consultant   `json:"manager"`
	TimeEntries   []TimeEntry  `json:"timeEntries"`
	StatusHistory []Status     `json:"statusHistory"`
	Consultants   []Consultant `json:"consultants"`
}
