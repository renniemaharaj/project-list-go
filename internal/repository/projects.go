package repository

import (
	"context"

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

// GetProjectByID will get and return a project by ID and (error or nil)
func (r *repository) GetProjectByID(ctx context.Context, projectID int) (*entity.Project, error) {
	var p entity.Project
	err := r.db.Select().From("projects").Where(dbx.HashExp{"id": projectID}).One(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetAllProjectsDesc will list projects and return projects
func (r *repository) GetAllProjectsDesc(ctx context.Context) ([]entity.Project, error) {
	var list []entity.Project
	err := r.db.Select().From("projects").OrderBy("id DESC").All(&list)
	return list, err
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
