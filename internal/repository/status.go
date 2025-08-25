package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// InsertProjectStatusByStruct will insert a new status relating to project id into status table
func (r *repository) InsertProjectStatusByStruct(ctx context.Context, s *entity.Status) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("project_statuses", dbx.Params{
			"title":         s.Title,
			"description":   s.Description,
			"project_id":    s.ProjectID,
			"consultant_id": s.ConsultantID,
		}).Execute()
		return err
	})
}

// GetStatusHistoryByProjectID will return all project_statuses relating to the projectID (history)
func (r *repository) GetStatusHistoryByProjectID(ctx context.Context, projectID int) ([]entity.Status, error) {
	var list []entity.Status
	err := r.DB.Select().
		From("project_statuses").
		Where(&dbx.HashExp{"project_id": projectID}).
		OrderBy("id DESC").
		All(&list)
	return list, err
}
