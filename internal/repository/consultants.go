package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// InsertConsultantByStruct will insert a consultant into consultans table from consultant struct
func (r *repository) InsertConsultantByStruct(ctx context.Context, c *entity.Consultant) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("consultants", dbx.Params{
			"first_name":      c.FirstName,
			"last_name":       c.LastName,
			"email":           c.Email,
			"profile_picture": c.ProfilePicture,
		}).Execute()
		return err
	})
}

// GetConsultantByID will get and return consultant by id
func (r *repository) GetConsultantByID(ctx context.Context, consultantID int) (*entity.Consultant, error) {
	var c entity.Consultant
	err := r.db.Select().From("consultants").Where(dbx.HashExp{"id": consultantID}).One(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetConsultantsByProjectID gets and returns project consultants
func (r *repository) GetConsultantsByProjectID(ctx context.Context, projectID int) ([]entity.Consultant, error) {
	var pc = []entity.ConsultantProject{}
	err := r.db.Select().From("project_consultants").Where(dbx.HashExp{"project_id": projectID}).All(&pc)
	if err != nil {
		return nil, err
	}

	var consultants = []entity.Consultant{}
	for _, c := range pc {
		consultant, err := r.GetConsultantByID(ctx, c.ConsultantID)
		if err != nil {
			return nil, err
		}
		consultants = append(consultants, *consultant)
	}
	return consultants, nil
}

// GetRelatedConsultantsByProjectID gets and returns project consultants
// Includes consultants linked via project_consultants and project_time_entries
func (r *repository) GetRelatedConsultantsByProjectID(ctx context.Context, projectID int) ([]entity.Consultant, error) {
	uniqueIDs := make(map[int]struct{})

	// 1. Collect consultants from project_consultants
	var pc []entity.ConsultantProject
	if err := r.db.Select().From("project_consultants").
		Where(dbx.HashExp{"project_id": projectID}).All(&pc); err != nil {
		return nil, err
	}
	for _, c := range pc {
		uniqueIDs[c.ConsultantID] = struct{}{}
	}

	// 2. Collect consultants from time_entries
	var tes []entity.TimeEntry
	if err := r.db.Select("consultant_id").
		From("project_time_entries").
		Where(dbx.HashExp{"project_id": projectID}).All(&tes); err != nil {
		return nil, err
	}
	for _, te := range tes {
		uniqueIDs[te.ConsultantID] = struct{}{}
	}

	// Fetch consultant details and merge
	var consultants []entity.Consultant
	for id := range uniqueIDs {
		consultant, err := r.GetConsultantByID(ctx, id)
		if err != nil {
			return nil, err
		}
		consultants = append(consultants, *consultant)
	}

	return consultants, nil
}

// GetAllConsultants will get and return all consultants from consultants table
func (r *repository) GetAllConsultants(ctx context.Context) ([]entity.Consultant, error) {
	var list []entity.Consultant
	err := r.db.Select().From("consultants").All(&list)
	return list, err
}

// UpdateConsultantByStruct will update a consultant from consultants table
func (r *repository) UpdateConsultantByStruct(ctx context.Context, c *entity.Consultant) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Update("consultants", dbx.Params{
			"first_name":      c.FirstName,
			"last_name":       c.LastName,
			"email":           c.Email,
			"profile_picture": c.ProfilePicture,
		}, dbx.HashExp{"id": c.ID}).Execute()
		return err
	})
}

// DeleteConsultantByID will delete a consultant by id from consultants table
func (r *repository) DeleteConsultantByID(ctx context.Context, consultantID int) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Delete("consultants", dbx.HashExp{"id": consultantID}).Execute()
		return err
	})
}

// InsertConsultantRoleByStruct will delete a consultant by id from consultants table
func (r *repository) InsertConsultantRoleByStruct(ctx context.Context, consultantRole entity.ConsultantRole) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("consultant_roles",
			dbx.Params{
				"consultant_id": consultantRole.ConsultantID,
				"role":          consultantRole.Role,
			}).Execute()
		return err
	})
}
