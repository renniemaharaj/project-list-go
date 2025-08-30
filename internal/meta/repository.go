package meta

import (
	"context"
	"fmt"
	"time"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/consultant"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	"github.com/renniemaharaj/project-list-go/internal/project"
	"github.com/renniemaharaj/project-list-go/internal/status"
	internalTime "github.com/renniemaharaj/project-list-go/internal/time"
)

type Repository interface {
	GetProjectMetaByProjectID(ctx context.Context, projectID int) (*entity.ProjectMeta, error)
	GetProjectsMetaByProjectIDS(ctx context.Context, projectIDs []int, light bool) (map[int]entity.ProjectMeta, []entity.Project, error)
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
	timeEntries, err := internalTime.NewRepository(r.dbContext, r.l).GetTimeEntryHistoryByProjectID(ctx, projectID)
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
	manager, err := consultant.NewRepository(r.dbContext, r.l).GetConsultantDataByID(ctx, project.ManagerID)
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

// GetProjectsMetaByProjectIDS will return meta data for multiple projects in batch.
// It reduces thousands of queries (per project) into 5 total batched queries and logs timings.
func (r *repository) GetProjectsMetaByProjectIDS(ctx context.Context, projectIDs []int, light bool) (map[int]entity.ProjectMeta, []entity.Project, error) {
	start := time.Now()
	r.l.Info(fmt.Sprintf("Starting GetProjectsMetaByProjectIDS for %d projects", len(projectIDs)))

	// --- 1. Batch fetch time entries ---
	callStart := time.Now()
	timeEntries, err := internalTime.NewRepository(r.dbContext, r.l).GetTimeEntryHistoryByProjectsIDS(ctx, projectIDs)
	if err != nil {
		return nil, nil, err
	}
	r.l.Info(fmt.Sprintf("Fetched %d time entries in %v", len(timeEntries), time.Since(callStart)))

	// --- 2. Batch fetch status history ---
	callStart = time.Now()
	statusHistory, err := status.NewRepository(r.dbContext, r.l).GetStatusHistoryByProjectsIDS(ctx, projectIDs)
	if err != nil {
		return nil, nil, err
	}
	r.l.Info(fmt.Sprintf("Fetched %d status history entries in %v", len(statusHistory), time.Since(callStart)))

	// --- 3. Batch fetch projects ---
	callStart = time.Now()
	projects, err := project.NewRepository(r.dbContext, r.l).GetProjectsDataByIDS(ctx, projectIDs)
	if err != nil {
		return nil, nil, err
	}
	r.l.Info(fmt.Sprintf("Fetched %d projects in %v", len(projects), time.Since(callStart)))

	// --- 4. Batch fetch related consultants ---
	callStart = time.Now()
	projectConsultants, err := consultant.NewRepository(r.dbContext, r.l).GetRelatedConsultantsByProjectsIDS(ctx, projectIDs)
	if err != nil {
		return nil, nil, err
	}
	r.l.Info(fmt.Sprintf("Fetched %d project consultants in %v", len(projectConsultants), time.Since(callStart)))

	// --- 5. Batch fetch managers ---
	callStart = time.Now()
	managerIDs := make([]int, 0, len(projects))
	for _, p := range projects {
		managerIDs = append(managerIDs, p.ManagerID)
	}
	managers, err := consultant.NewRepository(r.dbContext, r.l).GetConsultantDataByIDS(ctx, managerIDs)
	if err != nil {
		return nil, nil, err
	}
	r.l.Info(fmt.Sprintf("Fetched %d managers in %v", len(managers), time.Since(callStart)))

	// --- 6. Group results into maps for quick lookup ---
	timeMap := make(map[int][]entity.TimeEntry)
	for _, t := range timeEntries {
		timeMap[t.ProjectID] = append(timeMap[t.ProjectID], t)
	}

	statusMap := make(map[int][]entity.ProjectStatus)
	for _, s := range statusHistory {
		statusMap[s.ProjectID] = append(statusMap[s.ProjectID], s)
	}

	projectMap := make(map[int]entity.Project)
	for _, p := range projects {
		projectMap[p.ID] = p
	}

	managerMap := make(map[int]entity.Consultant)
	for _, m := range managers {
		managerMap[m.ID] = m
	}

	consultantsMap := make(map[int][]entity.Consultant)
	for _, c := range projectConsultants {
		consultantsMap[c.ProjectID] = append(consultantsMap[c.ProjectID], c.Consultant)
	}

	// --- 7. Construct ProjectMeta map ---
	projectMetas := make(map[int]entity.ProjectMeta, len(projectIDs))
	for _, pid := range projectIDs {
		p := projectMap[pid]

		projectMetas[pid] = entity.ProjectMeta{
			TimeEntries:   timeMap[pid],
			StatusHistory: statusMap[pid],
			Manager:       managerMap[p.ManagerID],
			Consultants:   consultantsMap[pid],
		}
	}

	r.l.Info(fmt.Sprintf("Completed GetProjectsMetaByProjectIDS for %d projects in total %v", len(projectIDs), time.Since(start)))
	return projectMetas, projects, nil
}
