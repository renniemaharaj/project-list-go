package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// =================== PROJECT TAGS ===================

// AddProjectTagByProjectID, using a projectID will insert a new tag into project_tags table
func (r *repository) AddProjectTagByProjectID(ctx context.Context, projectID int, tag string) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("project_tags", dbx.Params{
			"project_id": projectID,
			"tag":        tag,
		}).Execute()
		return err
	})
}

// RemoveProjectTagByProjectID, using projectID && tag, will remove tag from project_tags table
func (r *repository) RemoveProjectTagByProjectID(ctx context.Context, projectID int, tag string) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Delete("project_tags", dbx.HashExp{
			"project_id": projectID,
			"tag":        tag,
		}).Execute()
		return err
	})
}

// GetProjectTagsByProjectID, from project_tags table, will return all tags with projectID
func (r *repository) GetProjectTagsByProjectID(ctx context.Context, projectID int) ([]string, error) {
	var tags []string
	err := r.db.Select("tag").
		From("project_tags").
		Where(dbx.HashExp{"project_id": projectID}).
		Column(&tags)
	return tags, err
}
