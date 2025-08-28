package meta

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/consultant"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	"github.com/renniemaharaj/project-list-go/internal/project"
	"github.com/renniemaharaj/project-list-go/internal/status"
	"github.com/renniemaharaj/project-list-go/internal/time"
)

type Repository interface {
	GetProjectMetaByProjectID(ctx context.Context, projectID int) (*entity.ProjectMeta, error)
}

type repository struct {
	dbContext *database.DBContext
	l         *logger.Logger
}

func NewRepository(_db *database.DBContext, _l *logger.Logger) Repository {
	return &repository{_db, _l}
}

// GetProjectByID will get and return a project meta data by ID and (error or nil)
func (r *repository) GetProjectMetaByProjectID(ctx context.Context, projectID int) (*entity.ProjectMeta, error) {
	var projectMeta entity.ProjectMeta
	// first get time entries
	timeEntries, err := time.NewRepository(r.dbContext, r.l).GetTimeEntryHistoryByProjectID(ctx, projectID)
	if err != nil {
		r.l.Fatal(err)
	}
	projectMeta.TimeEntries = timeEntries
	// second get status history
	statusHistory, err := status.NewRepository(r.dbContext, r.l).GetStatusHistoryByProjectID(ctx, projectID)
	if err != nil {
		r.l.Fatal(err)
	}
	projectMeta.StatusHistory = statusHistory
	// third get project
	project, err := project.NewRepository(r.dbContext, r.l).GetProjectDataByID(ctx, projectID)
	if err != nil {
		r.l.Fatal(err)
	}
	// fourth get manager from project
	manager, err := consultant.NewRepository(r.dbContext, r.l).GetConsultantByID(ctx, project.ManagerID)
	if err != nil {
		r.l.Fatal(err)
	}
	projectMeta.Manager = *manager
	// fifth get project consultans
	projectConsultants, err := consultant.NewRepository(r.dbContext, r.l).GetRelatedConsultantsByProjectID(ctx, project.ID)
	if err != nil {
		r.l.Fatal(err)
	}
	projectMeta.Consultants = projectConsultants
	return &projectMeta, nil
}
