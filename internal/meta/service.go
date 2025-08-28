package meta

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Service interface {
	GetProjectMetaByProjectID(ctx context.Context, projectID int) (*ProjectMeta, error)
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
