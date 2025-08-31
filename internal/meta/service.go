package meta

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	"github.com/renniemaharaj/project-list-go/internal/project"
)

type Service interface {
	GetProjectMetaByProjectID(ctx context.Context, projectID int) (*ProjectMeta, error)
	GetProjectsMetaByProjectIDS(ctx context.Context, projectIDs []int, light bool) (map[int]ProjectMeta, []project.Project, error)
}

// Service
type service struct {
	repo   Repository
	logger *logger.Logger
}

type ProjectMeta struct {
	entity.ProjectMeta
}

func NewService(repo Repository, logger *logger.Logger) Service {
	return &service{repo, logger}
}

func (s *service) GetProjectMetaByProjectID(ctx context.Context, projectID int) (*ProjectMeta, error) {
	md, err := s.repo.GetProjectMetaByProjectID(ctx, projectID)
	if err != nil {
		return &ProjectMeta{}, err
	}
	return &ProjectMeta{*md}, err
}

func (s *service) GetProjectsMetaByProjectIDS(ctx context.Context, projectIDs []int, light bool) (map[int]ProjectMeta, []project.Project, error) {
	projectMetas, projects, err := s.repo.GetProjectsMetaByProjectIDS(ctx, projectIDs, light)
	if err != nil {
		return nil, nil, err
	}
	resultMetas := make(map[int]ProjectMeta)
	for k, v := range projectMetas {
		resultMetas[k] = ProjectMeta{v}
	}
	projectData := make([]project.Project, len(projects))
	for i, p := range projects {
		projectData[i] = project.Project{p}
	}
	return resultMetas, projectData, nil
}
