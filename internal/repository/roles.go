package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// InsertConsultantRoleByStruct will insert a consultant role into consultant_roles table
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
