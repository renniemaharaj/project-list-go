package project

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Service interface {
	InsertProjectByStruct(ctx context.Context, project *Project) error
	GetProjectDataByID(ctx context.Context, projectID int) (*Project, error)
	GetProjectsDataByIDS(ctx context.Context, ids []int) ([]Project, error)
	GetAllProjectIDS(ctx context.Context) ([]int, error)
	GetProjectIDSByPage(ctx context.Context, limit, offset int) ([]int, error)
	GetProjectIDSBySearchQuery(ctx context.Context, searchQuery string, limit, offset int) ([]int, error)
	UpdateProjectByStruct(ctx context.Context, project *Project) error
	DeleteProjectByID(ctx context.Context, projectID int) error
}

// Service
type service struct {
	repo   Repository
	logger *logger.Logger
}

type Project struct {
	entity.Project
}

// NewService creates a new consultant service
func NewService(repo Repository, logger *logger.Logger) Service {
	return &service{repo, logger}
}

func (s *service) InsertProjectByStruct(ctx context.Context, project *Project) error {
	return s.repo.InsertProjectByStruct(ctx, &project.Project)
}

func (s *service) GetProjectDataByID(ctx context.Context, projectID int) (*Project, error) {
	project, err := s.repo.GetProjectDataByID(ctx, projectID)
	if err != nil {
		return &Project{}, err
	}
	return &Project{*project}, nil
}

func (s *service) GetProjectsDataByIDS(ctx context.Context, ids []int) ([]Project, error) {
	projects, err := s.repo.GetProjectsDataByIDS(ctx, ids)
	if err != nil {
		return []Project{}, err
	}
	results := []Project{}
	for _, p := range projects {
		results = append(results, Project{p})
	}
	return results, nil
}

func (s *service) GetAllProjectIDS(ctx context.Context) ([]int, error) {
	return s.repo.GetAllProjectIDS(ctx)
}

func (s *service) GetProjectIDSByPage(ctx context.Context, limit, offset int) ([]int, error) {
	return s.repo.GetProjectIDSByPage(ctx, limit, offset)
}

func (s *service) GetProjectIDSBySearchQuery(ctx context.Context, searchQuery string, limit, offset int) ([]int, error) {
	return s.repo.GetProjectIDSBySearchQuery(ctx, searchQuery, limit, offset)
}

func (s *service) UpdateProjectByStruct(ctx context.Context, project *Project) error {
	return s.repo.UpdateProjectByStruct(ctx, &project.Project)
}

func (s *service) DeleteProjectByID(ctx context.Context, projectID int) error {
	return s.repo.DeleteProjectByID(ctx, projectID)
}
