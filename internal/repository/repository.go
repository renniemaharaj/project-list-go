package repository

import (
	"context"

	"github.com/renniemaharaj/project-list-go/internal/entity"

	_ "github.com/lib/pq"
)

// Repository definition
type Repository interface {
	// --- Consultant ---
	CreateConsultant(ctx context.Context, c *entity.Consultant) error
	GetConsultantByID(ctx context.Context, id int) (*entity.Consultant, error)
	ListConsultants(ctx context.Context) ([]entity.Consultant, error)
	UpdateConsultant(ctx context.Context, c *entity.Consultant) error
	DeleteConsultant(ctx context.Context, id int) error

	// --- Project ---
	CreateProject(ctx context.Context, p *entity.Project) error
	GetProjectByID(ctx context.Context, id int) (*entity.Project, error)
	ListProjects(ctx context.Context) ([]entity.Project, error)
	UpdateProject(ctx context.Context, p *entity.Project) error
	DeleteProject(ctx context.Context, id int) error

	// --- Project Tags ---
	AddProjectTag(ctx context.Context, projectID int, tag string) error
	RemoveProjectTag(ctx context.Context, projectID int, tag string) error
	GetProjectTags(ctx context.Context, projectID int) ([]string, error)

	// --- Status & History ---
	CreateStatus(ctx context.Context, s *entity.Status) error
	GetStatusHistory(ctx context.Context, projectID int) ([]entity.Status, error)

	// --- Time Entries ---
	CreateTimeEntry(ctx context.Context, e *entity.TimeEntry) error
	GetTimeEntryByID(ctx context.Context, id int) (*entity.TimeEntry, error)
	ListTimeEntriesByProject(ctx context.Context, projectID int) ([]entity.TimeEntry, error)
	ListTimeEntriesByConsultant(ctx context.Context, consultantID int) ([]entity.TimeEntry, error)
	UpdateTimeEntry(ctx context.Context, e *entity.TimeEntry) error
	DeleteTimeEntry(ctx context.Context, id int) error
}
