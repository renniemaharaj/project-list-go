package demo

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

type Service interface {
	GenerateInsertDemoData(ctx context.Context) error
}

// Service
type service struct {
	repo   Repository
	logger *logger.Logger
}

func NewService(repo Repository, logger *logger.Logger) Service {
	return &service{repo, logger}
}

func (s *service) GenerateInsertDemoData(ctx context.Context) error {
	return s.repo.GenerateInsertDemoData(ctx)
}
