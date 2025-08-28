package project

import (
	"context"
	"strings"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	internalIDRows "github.com/renniemaharaj/project-list-go/internal/idRow"
)

type Repository interface {
	InsertProjectByStruct(ctx context.Context, p *entity.Project) error
	GetProjectDataByID(ctx context.Context, projectID int) (*entity.Project, error)
	GetProjectsDataByIDS(ctx context.Context, ids []int) ([]entity.Project, error)
	GetAllProjectIDS(ctx context.Context) ([]int, error)
	GetProjectIDSByPage(ctx context.Context, limit, offset int) ([]int, error)
	GetProjectIDSBySearchQuery(ctx context.Context, searchQuery string) ([]int, error)
	UpdateProjectByStruct(ctx context.Context, p *entity.Project) error
	DeleteProjectByID(ctx context.Context, projectID int) error
}

type repository struct {
	dbContext *database.DBContext
	logger    *logger.Logger
}

func NewRepository(dbContext *database.DBContext, logger *logger.Logger) Repository {
	return &repository{dbContext, logger}
}

// InsertProjectByStruct will insert a project from project struct
func (r *repository) InsertProjectByStruct(ctx context.Context, p *entity.Project) error {
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("projects", dbx.Params{
			"projected_start_date": p.ProjectedStartDate,
			"start_date":           p.StartDate,
			"projected_end_date":   p.ProjectedEndDate,
			"end_date":             p.EndDate,
			"number":               p.Number,
			"name":                 p.Name,
			"manager_id":           p.ManagerID,
			"description":          p.Description,
		}).Execute()
		return err
	})
}

// GetProjectDataByID will get and return a project by ID and (error or nil)
func (r *repository) GetProjectDataByID(ctx context.Context, projectID int) (*entity.Project, error) {
	var Project entity.Project
	err := r.dbContext.DBX.WithContext(ctx).Select().From("projects").Where(dbx.HashExp{"id": projectID}).One(&Project)
	if err != nil {
		return nil, err
	}
	return &Project, nil
}

// GetProjectsDataByIDS gets projects by a slice of int project IDs and (error or nil)
func (r *repository) GetProjectsDataByIDS(ctx context.Context, ids []int) ([]entity.Project, error) {
	if len(ids) == 0 {
		return []entity.Project{}, nil
	}

	// convert []int → []interface{} inline
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	var projects []entity.Project
	err := r.dbContext.DBX.WithContext(ctx).Select().
		From("projects").
		Where(dbx.In("id", args...)).
		All(&projects)

	return projects, err
}

// GetAllProjectIDS will list projects and return their IDs
func (r *repository) GetAllProjectIDS(ctx context.Context) ([]int, error) {
	idRows := internalIDRows.IDRows{}

	err := r.dbContext.DBX.WithContext(ctx).Select("p.id").
		From("projects p").
		OrderBy("id DESC").
		All(&idRows)
	if err != nil {
		return nil, err
	}

	// Extract just the ints
	ids := idRows.ToIntSlice()
	return ids, nil
}

// GetProjectIDSByPage will list all projects by page
func (r *repository) GetProjectIDSByPage(ctx context.Context, limit, offset int) ([]int, error) {
	idRows := internalIDRows.IDRows{}

	err := r.dbContext.DBX.WithContext(ctx).Select("p.id").
		From("projects p").
		OrderBy("id DESC").
		Limit(int64(limit)).
		Offset(int64(offset)).
		AndOrderBy("id DESC").
		All(&idRows)
	if err != nil {
		return nil, err
	}

	// Extract just the ints
	ids := idRows.ToIntSlice()
	return ids, nil
}

// GetProjectIDSBySearchQuery will use searchQuery and return matching project IDS
func (r *repository) GetProjectIDSBySearchQuery(ctx context.Context, searchQuery string) ([]int, error) {
	// Split `+` separated terms
	terms := strings.Split(searchQuery, "+")

	// Base query with joins
	q := r.dbContext.DBX.WithContext(ctx).Select("p.id").
		Distinct(true).
		From("projects p").
		InnerJoin("project_consultants pc", dbx.NewExp("pc.project_id = p.id")).
		InnerJoin("consultants c", dbx.NewExp("c.id = pc.consultant_id")).
		InnerJoin("project_time_entries te", dbx.NewExp("te.project_id = p.id")).
		InnerJoin("project_statuses s", dbx.NewExp("s.project_id = p.id")).
		InnerJoin("project_tags tg", dbx.NewExp("tg.project_id = p.id"))

	// For each term, build an OR condition across searchable fields
	// Search query keywords are paramterized to avoid injection
	for _, term := range terms {
		t := "%" + term + "%"
		q.AndWhere(dbx.Or(
			dbx.NewExp("p.name ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("p.description ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("c.first_name ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("c.last_name ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("c.email ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("te.title ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("te.description ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("s.title ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("s.description ILIKE {:t}", dbx.Params{"t": t}),
			dbx.NewExp("tg.tag ILIKE {:t}", dbx.Params{"t": t}),
		))
	}

	// Execute query
	idRows := internalIDRows.IDRows{}
	err := q.All(&idRows)
	if err != nil {
		return nil, err
	}

	// Extract just the ints
	ids := idRows.ToIntSlice()
	return ids, nil
}

// UpdateProjectByStruct will update a project by project struct ID
func (r *repository) UpdateProjectByStruct(ctx context.Context, p *entity.Project) error {
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Update("projects", dbx.Params{
			"projected_start_date": p.ProjectedStartDate,
			"start_date":           p.StartDate,
			"projected_end_date":   p.ProjectedEndDate,
			"end_date":             p.EndDate,
			"number":               p.Number,
			"name":                 p.Name,
			"manager_id":           p.ManagerID,
			"description":          p.Description,
		}, dbx.HashExp{"id": p.ID}).Execute()
		return err
	})
}

// DeleteProjectByID will delete a project by ID
func (r *repository) DeleteProjectByID(ctx context.Context, projectID int) error {
	return r.dbContext.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Delete("projects", dbx.HashExp{"id": projectID}).Execute()
		return err
	})
}
