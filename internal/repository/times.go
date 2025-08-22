package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// =================== TIME ENTRIES ===================

// CreateTimeEntry will insert a time entry to time_entries table
func (r *repository) CreateTimeEntry(ctx context.Context, e *entity.TimeEntry) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("time_entries", dbx.Params{
			"hours":         e.Hours,
			"title":         e.Title,
			"description":   e.Description,
			"consultant_id": e.ConsultantID,
			"project_id":    e.ProjectID,
			"type":          e.Type,
			"entry_date":    e.EntryDate,
		}).Execute()
		return err
	})
}

// GetTimeEntryByID will return a specific time entry by ID
func (r *repository) GetTimeEntryByID(ctx context.Context, id int) (*entity.TimeEntry, error) {
	var e entity.TimeEntry
	err := r.db.Select().From("time_entries").Where(dbx.HashExp{"id": id}).One(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// ListTimeEntriesByProject will return all time entries for project
func (r *repository) ListTimeEntriesByProject(ctx context.Context, projectID int) ([]entity.TimeEntry, error) {
	var list []entity.TimeEntry
	err := r.db.Select().From("time_entries").Where(dbx.HashExp{"project_id": projectID}).All(&list)
	return list, err
}

// ListTimeEntriesByConsultant will return all time entries by consultant
func (r *repository) ListTimeEntriesByConsultant(ctx context.Context, consultantID int) ([]entity.TimeEntry, error) {
	var list []entity.TimeEntry
	err := r.db.Select().From("time_entries").Where(dbx.HashExp{"consultant_id": consultantID}).All(&list)
	return list, err
}

// UpdateTimeEntry will update a time entry by ID
func (r *repository) UpdateTimeEntry(ctx context.Context, e *entity.TimeEntry) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Update("time_entries", dbx.Params{
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

// DeleteTimeEntry will delete a time entry by ID
func (r *repository) DeleteTimeEntry(ctx context.Context, id int) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Delete("time_entries", dbx.HashExp{"id": id}).Execute()
		return err
	})
}
