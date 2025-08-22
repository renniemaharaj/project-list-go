package repository

import (
	"context"
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// SeedDemoData inserts multiple demo consultants, projects, statuses, and time entries
func (r *repository) SeedDemoData(ctx context.Context) error {
	// Check if projects already exist
	var count int
	if err := r.db.NewQuery("SELECT COUNT(*) FROM projects").Row(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // already seeded
	}

	// --- Demo data ---
	consultants := []entity.Consultant{
		{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"},
	}

	projects := []struct {
		Name        string
		Number      string
		Description string
	}{
		{"Demo Project A", "PRJ-001", "Kickoff demo project A"},
		{"Demo Project B", "PRJ-002", "Kickoff demo project B"},
	}
	statuses := []string{"In Progress", "Completed"}

	// Slice of anonymous structs with Title, Type, and Hours
	timeEntries := []struct {
		Title string
		Type  string
		Hours float32
	}{
		{Title: "Assigned", Type: "debit", Hours: 1.5},
		{Title: "Work", Type: "credit", Hours: 2.5},
		{Title: "Review", Type: "debit", Hours: 1.0},
	}

	// --- Insert consultants ---
	for ci := range consultants {
		c := &consultants[ci]
		if err := r.CreateConsultant(ctx, c); err != nil {
			return err
		}
		if err := r.db.NewQuery(
			"SELECT id FROM consultants WHERE email = {:email}",
		).Bind(dbx.Params{"email": c.Email}).Row(&c.ID); err != nil {
			return err
		}

		// --- For each consultant, create projects ---
		for _, p := range projects {
			project := &entity.Project{
				Name:        fmt.Sprintf("%s (by %s)", p.Name, c.FirstName),
				Number:      fmt.Sprintf("%s-%d", p.Number, ci+1),
				Description: p.Description,
				ManagerID:   c.ID,
			}
			if err := r.CreateProject(ctx, project); err != nil {
				return err
			}
			if err := r.db.NewQuery(
				"SELECT id FROM projects WHERE number = {:number}",
			).Bind(dbx.Params{"number": project.Number}).Row(&project.ID); err != nil {
				return err
			}

			// --- Add statuses ---
			for _, st := range statuses {
				status := &entity.Status{
					Title:        st,
					Description:  fmt.Sprintf("Project %s is %s", project.Number, st),
					ConsultantID: c.ID,
					ProjectID:    project.ID,
				}
				if err := r.CreateStatus(ctx, status); err != nil {
					return err
				}
				if err := r.db.NewQuery(
					"SELECT id FROM statuses WHERE project_id = {:pid} AND consultant_id = {:cid} AND title = {:title}",
				).Bind(dbx.Params{"pid": project.ID, "cid": c.ID, "title": st}).Row(&status.ID); err != nil {
					return err
				}
			}

			// --- Add time entries ---
			for _, te := range timeEntries {
				entry := &entity.TimeEntry{
					Title:        te.Title,
					Description:  fmt.Sprintf("%s for %s", te.Title, project.Name),
					Hours:        te.Hours, // just vary a little
					Type:         te.Type,
					ConsultantID: c.ID,
					ProjectID:    project.ID,
				}
				if err := r.CreateTimeEntry(ctx, entry); err != nil {
					return err
				}
				if err := r.db.NewQuery(
					`SELECT id FROM time_entries 
					 WHERE title = {:title} AND consultant_id = {:cid} AND project_id = {:pid}`,
				).Bind(dbx.Params{"title": entry.Title, "cid": c.ID, "pid": project.ID}).Row(&entry.ID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
