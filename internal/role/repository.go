package role

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Repository interface {
	InsertConsultantRoleByStruct(ctx context.Context, consultantRole entity.ConsultantRole) error
}

type repository struct {
	dbContext *database.DBContext
	l         *logger.Logger
}

func NewRepository(dbContext *database.DBContext, _l *logger.Logger) Repository {
	return &repository{dbContext, _l}
}

// InsertConsultantRoleByStruct will insert a consultant role into consultant_roles table
func (r *repository) InsertConsultantRoleByStruct(ctx context.Context, consultantRole entity.ConsultantRole) error {
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("consultant_roles",
			dbx.Params{
				"consultant_id": consultantRole.ConsultantID,
				"role":          consultantRole.Role,
			}).Execute()
		return err
	})
}
