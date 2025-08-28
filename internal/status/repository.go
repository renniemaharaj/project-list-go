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

// GetStatusHistoryByProjectID will return all project_statuses relating to the projectID (history)
func (r *repository) GetStatusHistoryByProjectID(ctx context.Context, projectID int) ([]entity.ProjectStatus, error) {
	var list []entity.ProjectStatus
	err := r.dbContext.DBX.WithContext(ctx).Select().
		From("project_statuses").
		Where(&dbx.HashExp{"project_id": projectID}).
		OrderBy("id DESC").
		All(&list)
	return list, err
}
