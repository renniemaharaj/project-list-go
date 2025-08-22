package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// =================== STATUS & HISTORY ===================

// CreateStatus will insert a new status relating to project id into status table
func (r *repository) CreateStatus(ctx context.Context, s *entity.Status) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("status", dbx.Params{
			"title":       s.Title,
			"description": s.Description,
			"project_id":  s.ProjectID,
		}).Execute()
		return err
	})
}

// GetStatusHistory will return all statuses relating to the projectID (hence, hystory)
func (r *repository) GetStatusHistory(ctx context.Context, projectID int) ([]entity.Status, error) {
	var list []entity.Status
	err := r.db.Select().From("status").Where(&dbx.HashExp{"project_id": projectID}).All(&list)
	return list, err
}
