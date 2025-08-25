package repository

import (
	"context"
	"strings"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// InsertProjectByStruct will insert a project from project struct
func (r *repository) InsertProjectByStruct(ctx context.Context, p *entity.Project) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
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
	err := r.DB.Select().From("projects").Where(dbx.HashExp{"id": projectID}).One(&Project)
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
	err := r.DB.Select().
		From("projects").
		Where(dbx.In("id", args...)).
		All(&projects)

	return projects, err
}

// GetAllProjectIDS will list projects and return their IDs
func (r *repository) GetAllProjectIDS(ctx context.Context) ([]int, error) {
	var rows []struct {
		ID int `json:"id"`
	}

	err := r.DB.Select("p.id").
		From("projects p").
		OrderBy("id DESC").
		All(&rows)
	if err != nil {
		return nil, err
	}

	// Extract just the ints
	ids := make([]int, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}
	return ids, nil
}

// GetProjectIDSByPage will list all projects by page
func (r *repository) GetProjectIDSByPage(ctx context.Context, limit, offset int) ([]int, error) {
	var rows []struct {
		ID int `json:"id"`
	}

	err := r.DB.Select("p.id").
		From("projects p").
		OrderBy("id DESC").
		Limit(int64(limit)).
		Offset(int64(offset)).
		AndOrderBy("id DESC").
		All(&rows)
	if err != nil {
		return nil, err
	}

	// Extract just the ints
	ids := make([]int, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}
	return ids, nil
}

// GetProjectIDSBySearchQuery will use searchQuery and return matching project IDS
func (r *repository) GetProjectIDSBySearchQuery(ctx context.Context, searchQuery string) ([]int, error) {
	// Split `+` separated terms
	terms := strings.Split(searchQuery, "+")

	// Base query with joins
	q := r.DB.Select("p.id").
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
	var rows []struct {
		ID int `json:"id"`
	}
	err := q.All(&rows)
	if err != nil {
		return nil, err
	}

	// Extract just the ints
	ids := make([]int, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}
	return ids, nil
}

// UpdateProjectByStruct will update a project by project struct ID
func (r *repository) UpdateProjectByStruct(ctx context.Context, p *entity.Project) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
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
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Delete("projects", dbx.HashExp{"id": projectID}).Execute()
		return err
	})
}

// InsertProjectConsultantByStruct adds a consultant to project
func (r *repository) InsertProjectConsultantByStruct(ctx context.Context, projectConsultant entity.ConsultantProject) error {
	return r.UseTransaction(ctx, func(tx *dbx.Tx) error {
		_, err := tx.Insert("project_consultants", dbx.Params{
			"consultant_id": projectConsultant.ConsultantID,
			"project_id":    projectConsultant.ProjectID,
			"role":          projectConsultant.Role,
		}).Execute()
		return err
	})
}
