package status

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Repository interface {
	InsertProjectStatusByStruct(ctx context.Context, s *entity.ProjectStatus) error
	GetStatusHistoryByProjectID(ctx context.Context, projectID int) ([]entity.ProjectStatus, error)
	GetStatusHistoryByProjectsIDS(ctx context.Context, projectIDS []int) ([]entity.ProjectStatus, error)
}

type repository struct {
	dbContext *database.DBContext
	logger    *logger.Logger
}

func NewRepository(dbContext *database.DBContext, _l *logger.Logger) Repository {
	return &repository{dbContext, _l}
}

// InsertProjectStatusByStruct will insert a new status relating to project id into status table
func (r *repository) InsertProjectStatusByStruct(ctx context.Context, s *entity.ProjectStatus) error {
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("project_statuses", dbx.Params{
			"title":         s.Title,
			"description":   s.Description,
			"project_id":    s.ProjectID,
			"consultant_id": s.ConsultantID,
		}).Execute()
		return err
	})
}

// GetStatusHistoryByProjectsIDS will return all project_statuses relating to the given projectIDs (history)
func (r *repository) GetStatusHistoryByProjectsIDS(ctx context.Context, projectIDS []int) ([]entity.ProjectStatus, error) {
	var list []entity.ProjectStatus

	// Convert []int -> []interface{} for dbx.In
	args := make([]interface{}, len(projectIDS))
	for i, id := range projectIDS {
		args[i] = id
	}

	err := r.dbContext.Get().WithContext(ctx).Select().
		From("project_statuses").
		Where(dbx.In("project_id", args...)).
		OrderBy("id DESC").
		All(&list)

	return list, err
}

// GetStatusHistoryByProjectID will return all project_statuses relating to the projectID (history)
func (r *repository) GetStatusHistoryByProjectID(ctx context.Context, projectID int) ([]entity.ProjectStatus, error) {
	var list []entity.ProjectStatus
	err := r.dbContext.Get().WithContext(ctx).Select().
		From("project_statuses").
		Where(&dbx.HashExp{"project_id": projectID}).
		OrderBy("id DESC").
		All(&list)
	return list, err
}
