package tag

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Service interface {
	InsertProjectTagByStruct(ctx context.Context, projectTag ProjectTag) error
}

// Service
type service struct {
	repo   Repository
	logger *logger.Logger
}

type ProjectTag struct {
	entity.ProjectTag
}

func NewService(repo Repository, logger *logger.Logger) Service {
	return &service{repo, logger}
}

func (s *service) InsertProjectTagByStruct(ctx context.Context, projectTag ProjectTag) error {
	return s.repo.InsertProjectTagByStruct(ctx, projectTag.ProjectTag)
}
