package repository

import (
	"context"

	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// GetProjectByID will get and return a project by ID and (error or nil)
func (r *repository) GetProjectMetaByProjectID(ctx context.Context, projectID int) (*entity.ProjectMetaData, error) {
	var projectMeta entity.ProjectMetaData
	// first get time entries
	timeEntries, err := r.GetTimeEntryHistory(ctx, projectID)
	if err != nil {
		r.l.Fatal(err)
	}
	projectMeta.TimeEntries = timeEntries
	// second get status history
	statusHistory, err := r.GetStatusHistoryByProjectID(ctx, projectID)
	if err != nil {
		r.l.Fatal(err)
	}
	projectMeta.StatusHistory = statusHistory
	// third get project
	project, err := r.GetProjectByID(ctx, projectID)
	if err != nil {
		r.l.Fatal(err)
	}
	// fourth get manager from project
	manager, err := r.GetConsultantByID(ctx, project.ManagerID)
	if err != nil {
		r.l.Fatal(err)
	}
	projectMeta.Manager.ID = manager.ID
	return &projectMeta, nil
}
