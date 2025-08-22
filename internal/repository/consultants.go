package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// =================== CONSULTANTS ===================

// CreateConsultant will insert a consultant into consultans table from consultant struct
func (r *repository) CreateConsultant(ctx context.Context, c *entity.Consultant) error {
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

// ListConsultants will get and return all consultants from consultants table
func (r *repository) ListConsultants(ctx context.Context) ([]entity.Consultant, error) {
	var list []entity.Consultant
	err := r.db.Select().From("consultants").All(&list)
	return list, err
}

// UpdateConsultant will update a consultant from consultants table
func (r *repository) UpdateConsultant(ctx context.Context, c *entity.Consultant) error {
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
