package consultant

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

type Repository interface {
	InsertConsultantByStruct(ctx context.Context, c *entity.Consultant) error
	GetConsultantDataByID(ctx context.Context, consultantID int) (*entity.Consultant, error)
	GetConsultantDataByIDS(ctx context.Context, consultantIDS []int) ([]entity.Consultant, error)
	GetConsultantsByProjectID(ctx context.Context, projectID int) ([]entity.Consultant, error)
	GetRelatedConsultantsByProjectID(ctx context.Context, projectID int) ([]entity.Consultant, error)
	GetRelatedConsultantsByProjectsIDS(ctx context.Context, projectIDs []int) ([]entity.ProjectConsultantLink, error)
	GetAllConsultants(ctx context.Context) ([]entity.Consultant, error)
	UpdateConsultantByStruct(ctx context.Context, c *entity.Consultant) error
	DeleteConsultantByID(ctx context.Context, consultantID int) error
	InsertProjectConsultantByStruct(ctx context.Context, projectConsultant entity.ProjectConsultant) error
}

type repository struct {
	dbContext *database.DBContext
	logger    *logger.Logger
}

func NewRepository(dbContext *database.DBContext, _l *logger.Logger) Repository {
	return &repository{dbContext, _l}
}

// InsertConsultantByStruct will insert a consultant into consultans table from consultant struct
func (r *repository) InsertConsultantByStruct(ctx context.Context, c *entity.Consultant) error {
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("consultants", dbx.Params{
			"first_name":      c.FirstName,
			"last_name":       c.LastName,
			"email":           c.Email,
			"profile_picture": c.ProfilePicture,
		}).Execute()
		return err
	})
}

// GetConsultantDataByIDS will get and return consultants by a list of IDs
func (r *repository) GetConsultantDataByIDS(ctx context.Context, consultantIDS []int) ([]entity.Consultant, error) {
	var list []entity.Consultant

	// Convert []int -> []interface{} for dbx.In
	args := make([]interface{}, len(consultantIDS))
	for i, id := range consultantIDS {
		args[i] = id
	}

	err := r.dbContext.Get().WithContext(ctx).Select().
		From("consultants").
		Where(dbx.In("id", args...)).
		OrderBy("id DESC").
		All(&list)

	return list, err
}

// GetConsultantDataByID will get and return consultant by id
func (r *repository) GetConsultantDataByID(ctx context.Context, consultantID int) (*entity.Consultant, error) {
	var c entity.Consultant
	err := r.dbContext.Get().WithContext(ctx).Select().From("consultants").Where(dbx.HashExp{"id": consultantID}).One(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetConsultantsByProjectID gets and returns project consultants
func (r *repository) GetConsultantsByProjectID(ctx context.Context, projectID int) ([]entity.Consultant, error) {
	var consultants []entity.Consultant

	err := r.dbContext.Get().WithContext(ctx).Select("c.*").
		From("project_consultants pc").
		InnerJoin("consultants c", dbx.NewExp("c.id = pc.consultant_id")).
		Where(dbx.HashExp{"pc.project_id": projectID}).
		All(&consultants)

	if err != nil {
		return nil, err
	}

	return consultants, nil
}

// GetRelatedConsultantsByProjectsIDS gets and returns project consultants for multiple project IDs
// Includes consultants linked via project_consultants and project_time_entries
func (r *repository) GetRelatedConsultantsByProjectsIDS(ctx context.Context, projectIDs []int) ([]entity.ProjectConsultantLink, error) {
	var consultants []entity.ProjectConsultantLink

	// Convert []int -> []interface{} for dbx.In
	args := make([]interface{}, len(projectIDs))
	for i, id := range projectIDs {
		args[i] = id
	}

	q1 := r.dbContext.Get().WithContext(ctx).
		Select("c.*", "pc.project_id").
		From("consultants c").
		InnerJoin("project_consultants pc", dbx.NewExp("c.id = pc.consultant_id")).
		Where(dbx.In("pc.project_id", args...))

	q2 := r.dbContext.Get().WithContext(ctx).
		Select("c.*", "te.project_id").
		From("consultants c").
		InnerJoin("project_time_entries te", dbx.NewExp("c.id = te.consultant_id")).
		Where(dbx.In("te.project_id", args...))

	sql := q1.Union(q2.Build())

	if err := sql.All(&consultants); err != nil {
		return nil, err
	}

	return consultants, nil
}

// GetRelatedConsultantsByProjectID gets and returns project consultants
// Includes consultants linked via project_consultants and project_time_entries
func (r *repository) GetRelatedConsultantsByProjectID(ctx context.Context, projectID int) ([]entity.Consultant, error) {
	var consultants []entity.Consultant

	q := r.dbContext.Get().WithContext(ctx).Select("c.*").
		Distinct(true).
		From("consultants c").
		InnerJoin("project_consultants pc", dbx.NewExp("c.id = pc.consultant_id")).
		InnerJoin("project_time_entries te", dbx.NewExp("c.id = te.consultant_id")).
		Where(dbx.Or(
			dbx.HashExp{"pc.project_id": projectID},
			dbx.HashExp{"te.project_id": projectID},
		))

	if err := q.All(&consultants); err != nil {
		return nil, err
	}

	return consultants, nil
}

// GetAllConsultants will get and return all consultants from consultants table
func (r *repository) GetAllConsultants(ctx context.Context) ([]entity.Consultant, error) {
	var list []entity.Consultant
	err := r.dbContext.Get().WithContext(ctx).Select().From("consultants").All(&list)
	return list, err
}

// UpdateConsultantByStruct will update a consultant from consultants table
func (r *repository) UpdateConsultantByStruct(ctx context.Context, c *entity.Consultant) error {
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
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
	// Delete will be done in a transaction which can be rolled back on returning error
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
		// 1. remove consultant from consultants tanle
		_, err := tx.Delete("consultants", dbx.HashExp{"id": consultantID}).Execute()
		if err != nil {
			return err
		}
		// 2. remove consultant time entries from project_time_entries
		_, err = tx.Delete("project_time_entries", dbx.HashExp{"consultant_id": consultantID}).Execute()
		if err != nil {
			return err
		}
		// 3. remove consultant statuses from project_statuses
		_, err = tx.Delete("project_statuses", dbx.HashExp{"consultant_id": consultantID}).Execute()
		if err != nil {
			return err
		}
		return nil
	})
}

// InsertProjectConsultantByStruct adds a consultant to project
func (r *repository) InsertProjectConsultantByStruct(ctx context.Context, projectConsultant entity.ProjectConsultant) error {
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("project_consultants", dbx.Params{
			"consultant_id": projectConsultant.ConsultantID,
			"project_id":    projectConsultant.ProjectID,
			"role":          projectConsultant.Role,
		}).Execute()
		return err
	})
}
