package status

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Service interface {
	InsertProjectStatusByStruct(ctx context.Context, s *ProjectStatus) error
	GetStatusHistoryByProjectID(ctx context.Context, projectID int) ([]ProjectStatus, error)
}

// Service
type service struct {
	repo   Repository
	logger *logger.Logger
}

type ProjectStatus struct {
	entity.ProjectStatus
}

func NewService(repo Repository, logger *logger.Logger) Service {
	return &service{repo, logger}
}

func (s *service) GetStatusHistoryByProjectID(ctx context.Context, projectID int) ([]ProjectStatus, error) {
	statuses, err := s.repo.GetStatusHistoryByProjectID(ctx, projectID)
	if err != nil {
		return []ProjectStatus{}, err
	}
	results := []ProjectStatus{}
	for _, st := range statuses {
		results = append(results, ProjectStatus{st})
	}
	return results, nil
}

func (s *service) InsertProjectStatusByStruct(ctx context.Context, projectStatus *ProjectStatus) error {
	return s.repo.InsertProjectStatusByStruct(ctx, &projectStatus.ProjectStatus)
}
