package role

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Service interface {
	InsertConsultantRoleByStruct(ctx context.Context, consultantRole ConsultantRole) error
}

// Service
type service struct {
	repo   Repository
	logger *logger.Logger
}

type ConsultantRole struct {
	entity.ConsultantRole
}

func NewService(repo Repository, logger *logger.Logger) Service {
	return &service{repo, logger}
}

func (s *service) InsertConsultantRoleByStruct(ctx context.Context, consultantRole ConsultantRole) error {
	return s.repo.InsertConsultantRoleByStruct(ctx, consultantRole.ConsultantRole)
}
