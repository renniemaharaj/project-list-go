package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// InsertProjectTagByStruct will insert a project tag into project_tags
func (r *repository) InsertProjectTagByStruct(ctx context.Context, tag entity.ProjectTag) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("project_tags", dbx.Params{
			"project_id": tag.ProjectID,
			"tag":        tag.Tag,
		}).Execute()
		return err
	})
}

// GetProjectTagsByProjectID, from project_tags table, will return all tags with projectID
func (r *repository) GetProjectTagsByProjectID(ctx context.Context, projectID int) ([]string, error) {
	var tags []string
	err := r.DB.Select("tag").
		From("project_tags").
		Where(dbx.HashExp{"project_id": projectID}).OrderBy("id DESC").
		Column(&tags)
	return tags, err
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
