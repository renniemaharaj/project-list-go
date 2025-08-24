package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// InsertTimeEntryByStruct will insert a time entry to project_time_entries table
func (r *repository) InsertTimeEntryByStruct(ctx context.Context, e *entity.TimeEntry) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("project_time_entries", dbx.Params{
			"hours":         e.Hours,
			"title":         e.Title,
			"description":   e.Description,
			"consultant_id": e.ConsultantID,
			"project_id":    e.ProjectID,
			"type":          e.Type,
		}).Execute()
		return err
	})
}

// GetTimeEntryByID will return a specific time entry by ID
func (r *repository) GetTimeEntryByID(ctx context.Context, id int) (*entity.TimeEntry, error) {
	var e entity.TimeEntry
	err := r.db.Select().From("project_time_entries").Where(dbx.HashExp{"id": id}).One(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetTimeEntryHistoryByProductID will return all time entries for project
func (r *repository) GetTimeEntryHistoryByProductID(ctx context.Context, projectID int) ([]entity.TimeEntry, error) {
	var list []entity.TimeEntry
	err := r.db.Select().
		From("project_time_entries").
		Where(dbx.HashExp{"project_id": projectID}).
		OrderBy("id DESC").
		All(&list)
	return list, err
}

// GetTimeEntryHistoryByConsultantID will return all time entries by consultant
func (r *repository) GetTimeEntryHistoryByConsultantID(ctx context.Context, consultantID int) ([]entity.TimeEntry, error) {
	var list []entity.TimeEntry
	err := r.db.Select().
		From("project_time_entries").
		Where(dbx.HashExp{"consultant_id": consultantID}).
		OrderBy("id DESC").
		All(&list)
	return list, err
}

// UpdateTimeEntryByStruct will update a time entry by ID
func (r *repository) UpdateTimeEntryByStruct(ctx context.Context, e *entity.TimeEntry) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Update("project_time_entries", dbx.Params{
			"hours":         e.Hours,
			"title":         e.Title,
			"description":   e.Description,
			"consultant_id": e.ConsultantID,
			"project_id":    e.ProjectID,
			"type":          e.Type,
			"entry_date":    e.EntryDate,
		}, dbx.HashExp{"id": e.ID}).Execute()
		return err
	})
}

// DeleteTimeEntryByID will delete a time entry by ID
func (r *repository) DeleteTimeEntryByID(ctx context.Context, id int) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Delete("project_time_entries", dbx.HashExp{"id": id}).Execute()
		return err
	})
}
