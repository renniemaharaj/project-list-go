package consultant

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// Consultant service interface
type Service interface {
	// Inserts a consultant from a consultant struct
	InsertConsultantByStruct(ctx context.Context, c *Consultant) error
	// Gets consultant by id
	GetConsultantByID(ctx context.Context, consultantID int) (*Consultant, error)
	// Gets all consultants who are explicitly attached to a project by projectID
	GetConsultantsByProjectID(ctx context.Context, projectID int) ([]Consultant, error)
	// Get all consultants who are explicitly attached and indirectly related to a project by...
	GetRelatedConsultantsByProjectID(ctx context.Context, projectID int) ([]Consultant, error)
	// Get all consultants from consultant table
	GetAllConsultants(ctx context.Context) ([]Consultant, error)
	// Updates a consultant by struct, struct must contain consultantID
	UpdateConsultantByStruct(ctx context.Context, c *Consultant) error
	// Deletes a consultant by consultantID
	DeleteConsultantByID(ctx context.Context, consultantID int) error
	// Inserts a consultant into project consultant table
	InsertProjectConsultantByStruct(ctx context.Context, projectConsultant ProjectConsultant) error
}

// Service layer consultant struct
type Consultant struct {
	entity.Consultant
}

type ProjectConsultant struct {
	entity.ProjectConsultant
}

// Service
type service struct {
	repo   Repository
	logger *logger.Logger
}

// NewService creates a new consultant service
func NewService(repo Repository, logger *logger.Logger) Service {
	return &service{repo, logger}
}

// Inserts a consultant from a consultant struct
func (s *service) InsertConsultantByStruct(ctx context.Context, c *Consultant) error {
	return s.repo.InsertConsultantByStruct(ctx, &c.Consultant)
}

// Gets consultant by id
func (s *service) GetConsultantByID(ctx context.Context, consultantID int) (*Consultant, error) {
	c, err := s.repo.GetConsultantByID(ctx, consultantID)
	if err != nil {
		return &Consultant{}, err
	}
	return &Consultant{*c}, nil
}

// Gets all consultants who are explicitly attached to a project by projectID
func (s *service) GetConsultantsByProjectID(ctx context.Context, projectID int) ([]Consultant, error) {
	cs, err := s.repo.GetConsultantsByProjectID(ctx, projectID)
	if err != nil {
		return []Consultant{}, err
	}
	result := []Consultant{}
	for _, c := range cs {
		result = append(result, Consultant{c})
	}
	return result, nil
}

// Get all consultants who are explicitly attached and indirectly related to a project by...
func (s *service) GetRelatedConsultantsByProjectID(ctx context.Context, projectID int) ([]Consultant, error) {
	cs, err := s.repo.GetRelatedConsultantsByProjectID(ctx, projectID)
	if err != nil {
		return []Consultant{}, err
	}
	result := []Consultant{}
	for _, c := range cs {
		result = append(result, Consultant{c})
	}
	return result, nil
}

// Get all consultants from consultant table
func (s *service) GetAllConsultants(ctx context.Context) ([]Consultant, error) {
	cs, err := s.repo.GetAllConsultants(ctx)
	if err != nil {
		return []Consultant{}, err
	}
	result := []Consultant{}
	for _, c := range cs {
		result = append(result, Consultant{c})
	}
	return result, nil
}

// Updates a consultant by struct, struct must contain consultantID
func (s *service) UpdateConsultantByStruct(ctx context.Context, c *Consultant) error {
	return s.repo.UpdateConsultantByStruct(ctx, &c.Consultant)
}

// Deletes a consultant by consultantID
func (s *service) DeleteConsultantByID(ctx context.Context, consultantID int) error {
	return s.repo.DeleteConsultantByID(ctx, consultantID)
}

func (s *service) InsertProjectConsultantByStruct(ctx context.Context, projectConsultant ProjectConsultant) error {
	return s.repo.InsertProjectConsultantByStruct(ctx, projectConsultant.ProjectConsultant)
}
