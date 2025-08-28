package time

import (
	"context"

	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Service interface {
	InsertTimeEntryByStruct(ctx context.Context, e *TimeEntry) error
	GetTimeEntryByTimeEntryID(ctx context.Context, id int) (*TimeEntry, error)
	GetTimeEntryHistoryByProjectID(ctx context.Context, projectID int) ([]TimeEntry, error)
	GetTimeEntryHistoryByConsultantID(ctx context.Context, consultantID int) ([]TimeEntry, error)
	UpdateTimeEntryByStruct(ctx context.Context, e *TimeEntry) error
	DeleteTimeEntryByTimeEntryID(ctx context.Context, id int) error
}

// Service
type service struct {
	repo   Repository
	logger *logger.Logger
}

type TimeEntry struct {
	entity.TimeEntry
}

func NewService(repo Repository, logger *logger.Logger) Service {
	return &service{repo, logger}
}

func (s *service) InsertTimeEntryByStruct(ctx context.Context, timeEntry *TimeEntry) error {
	return s.repo.InsertTimeEntryByStruct(ctx, &timeEntry.TimeEntry)
}

func (s *service) GetTimeEntryByTimeEntryID(ctx context.Context, id int) (*TimeEntry, error) {
	timeEntry, err := s.repo.GetTimeEntryByTimeEntryID(ctx, id)
	if err != nil {
		return &TimeEntry{}, nil
	}
	return &TimeEntry{*timeEntry}, nil
}

func (s *service) GetTimeEntryHistoryByProjectID(ctx context.Context, projectID int) ([]TimeEntry, error) {
	timeEntries, err := s.repo.GetTimeEntryHistoryByProjectID(ctx, projectID)
	if err != nil {
		return []TimeEntry{}, err
	}
	results := []TimeEntry{}
	for _, te := range timeEntries {
		results = append(results, TimeEntry{te})
	}
	return results, nil
}

func (s *service) GetTimeEntryHistoryByConsultantID(ctx context.Context, consultantID int) ([]TimeEntry, error) {
	timeEntries, err := s.repo.GetTimeEntryHistoryByConsultantID(ctx, consultantID)
	if err != nil {
		return []TimeEntry{}, err
	}
	results := []TimeEntry{}
	for _, te := range timeEntries {
		results = append(results, TimeEntry{te})
	}
	return results, nil
}

func (s *service) UpdateTimeEntryByStruct(ctx context.Context, timeEntry *TimeEntry) error {
	return s.repo.UpdateTimeEntryByStruct(ctx, &timeEntry.TimeEntry)
}

func (s *service) DeleteTimeEntryByTimeEntryID(ctx context.Context, id int) error {
	return s.repo.DeleteTimeEntryByTimeEntryID(ctx, id)
}
