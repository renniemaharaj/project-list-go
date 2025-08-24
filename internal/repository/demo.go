package repository

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

// InsertSeededDemoData inserts multiple demo consultants, projects, statuses, time entries,
// consultant roles, consultant projects, and project tags.
func (r *repository) InsertSeededDemoData(ctx context.Context) error {
	// Check if projects already exist
	var count int
	if err := r.db.NewQuery("SELECT COUNT(*) FROM projects").Row(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // already seeded
	}

	rand.Seed(time.Now().UnixNano())

	// --- Demo consultants ---
	consultants := []entity.Consultant{
		{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
		{FirstName: "Hannah", LastName: "Doe", Email: "hannah.doe@example.com"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"},
		{FirstName: "Sarah", LastName: "Connor", Email: "s.connor@example.com"},
		{FirstName: "Mike", LastName: "Johnson", Email: "mike.j@example.com"},
		{FirstName: "Peter", LastName: "Connor", Email: "p.connor@example.com"},
		{FirstName: "Paul", LastName: "Connor", Email: "paul.connor@example.com"},
		{FirstName: "Trey", LastName: "White", Email: "trey.white@example.com"},
		{FirstName: "Rebecca", LastName: "Maharaj", Email: "rebecca.maharaj@example.com"},
	}

	roles := []string{"administrator", "manager", "consultant"}

	// Insert consultants + generate projects
	for ci := range consultants {
		c := &consultants[ci]
		c.ProfilePicture = fmt.Sprintf("https://api.dicebear.com/7.x/lorelei/svg?seed=%s", url.QueryEscape(c.Email))

		if err := r.InsertConsultantByStruct(ctx, c); err != nil {
			return err
		}
		if err := r.db.NewQuery("SELECT id FROM consultants WHERE email = {:email}").Bind(dbx.Params{"email": c.Email}).Row(&c.ID); err != nil {
			return err
		}

		// Assign random role to consultant
		role := roles[rand.Intn(len(roles))]
		consultantRole := &entity.ConsultantRole{ConsultantID: c.ID, Role: role}
		if err := r.InsertConsultantRoleByStruct(ctx, *consultantRole); err != nil {
			return err
		}
	}

	for _, c := range consultants {
		// Each consultant gets between 1 and 3 demo projects
		projectCount := rand.Intn(3) + 1
		for pi := 0; pi < projectCount; pi++ {
			if err := r.generateProject(ctx, &c, pi+1, roles); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateProject creates a project, random statuses, time entries, consultant projects, and project tags
func (r *repository) generateProject(ctx context.Context, c *entity.Consultant, index int, roles []string) error {
	tags := []string{"support", "implementation", "custom", "reports", "software"}

	project := &entity.Project{
		Name:        fmt.Sprintf("Demo Project %d (by %s)", index, c.FirstName),
		Number:      fmt.Sprintf("PRJ-%03d-%d", index, c.ID),
		Description: fmt.Sprintf("Auto-generated demo project %d for consultant %s", index, c.FirstName),
		ManagerID:   c.ID,
	}
	if err := r.InsertProjectByStruct(ctx, project); err != nil {
		return err
	}
	if err := r.db.NewQuery("SELECT id FROM projects WHERE number = {:number}").Bind(dbx.Params{"number": project.Number}).Row(&project.ID); err != nil {
		return err
	}

	// --- Random statuses ---
	allStatuses := []string{"planned", "active", "on-hold", "completed"}
	statusCount := rand.Intn(len(allStatuses)-1) + 1 // at least 1 status
	for i := 0; i < statusCount; i++ {
		st := allStatuses[rand.Intn(len(allStatuses))]
		status := &entity.Status{Title: st, Description: fmt.Sprintf("Project %s is %s", project.Number, st), ConsultantID: c.ID, ProjectID: project.ID}
		if err := r.InsertProjectStatusByStruct(ctx, status); err != nil {
			return err
		}
	}

	// --- Assign all consultants randomly to this project ---
	allConsultants, err := r.GetAllConsultants(ctx)
	if err != nil {
		return err
	}
	assignedCount := rand.Intn(len(allConsultants)-1) + 1
	for i := 0; i < assignedCount; i++ {
		consultant := allConsultants[rand.Intn(len(allConsultants))]
		cp := &entity.ConsultantProject{ProjectID: project.ID, ConsultantID: consultant.ID, Role: roles[rand.Intn(len(roles))]}
		if err := r.InsertProjectConsultantByStruct(ctx, *cp); err != nil {
			return err
		}
	}

	// --- Random project tags ---
	tagCount := rand.Intn(len(tags)) + 1
	for i := 0; i < tagCount; i++ {
		tag := &entity.ProjectTag{ProjectID: project.ID, Tag: tags[rand.Intn(len(tags))]}
		if err := r.InsertProjectTagByStruct(ctx, *tag); err != nil {
			return err
		}
	}

	// --- Random debit + credit time entries ---
	debitCount := rand.Intn(3) + 1
	creditCount := rand.Intn(3) + 1
	createEntries := func(entryType string, count int) error {
		for i := 0; i < count; i++ {
			hours := float32(rand.Intn(4)+1) + rand.Float32()
			entry := &entity.TimeEntry{Title: fmt.Sprintf("%s #%d", entryType, i+1), Description: fmt.Sprintf("%s work for %s", entryType, project.Name), Hours: hours, Type: entryType, ConsultantID: rand.Intn(len(allConsultants)) + 1, ProjectID: project.ID}
			if err := r.InsertTimeEntryByStruct(ctx, entry); err != nil {
				return err
			}
		}
		return nil
	}
	if err := createEntries("debit", debitCount); err != nil {
		return err
	}
	if err := createEntries("credit", creditCount); err != nil {
		return err
	}

	return nil
}
