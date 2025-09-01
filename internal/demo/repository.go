package demo

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	fmtime "time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	internalConsultant "github.com/renniemaharaj/project-list-go/internal/consultant"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	internalProject "github.com/renniemaharaj/project-list-go/internal/project"
	internalRole "github.com/renniemaharaj/project-list-go/internal/role"
	internalStatus "github.com/renniemaharaj/project-list-go/internal/status"
	internalTag "github.com/renniemaharaj/project-list-go/internal/tag"
	"github.com/renniemaharaj/project-list-go/internal/time"
	"github.com/renniemaharaj/project-list-go/internal/utils"
)

type Repository interface {
	GenerateInsertDemoData(ctx context.Context) error
}

type repository struct {
	dbContext *database.DBContext
	l         *logger.Logger
}

func NewRepository(dbContext *database.DBContext, _l *logger.Logger) Repository {
	return &repository{dbContext, _l}
}

var (
	// use to get a random number, ((0 - max) + 1) of projects per consultant
	maxProjectPerConsultant = 5
)

// GenerateInsertDemoData inserts multiple demo consultants, projects, statuses, time entries,
// consultant roles, consultant projects, and project tags.
func (r *repository) GenerateInsertDemoData(ctx context.Context) error {
	// Check if projects already exist
	var count int
	if err := r.dbContext.Get().NewQuery("SELECT COUNT(*) FROM projects").Row(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // already seeded
	}

	rand.Seed(fmtime.Now().UnixNano())

	// --- Demo consultants ---
	consultants := []entity.Consultant{}

	contents, err := os.ReadFile("cfx_c.json")
	if err != nil {
		r.l.Fatal(err)
		panic("cfx_c.json containing consultants not found")
	}

	err = json.Unmarshal(contents, &consultants)
	if err != nil {
		r.l.Fatal(err)
		panic("Consultants hydrate error")
	}

	roles := []string{"administrator", "manager", "consultant"}

	if err := r.insertConsultants(ctx, consultants, roles); err != nil {
		return err
	}
	if err := r.generateProjectsForConsultants(ctx, consultants, roles); err != nil {
		return err
	}

	return nil
}

// Internal generateProject method coordinates project generations
func (r *repository) generateProject(ctx context.Context, c *entity.Consultant, index int, roles []string) error {
	tags := []string{"support", "implementation", "custom", "reports", "software"}

	project := &entity.Project{
		Name:        fmt.Sprintf("Demo Project %d (by %s)", index, c.FirstName),
		Number:      fmt.Sprintf("PRJ-%03d-%d", index, c.ID),
		Description: fmt.Sprintf("Auto-generated demo project %d for consultant %s", index, c.FirstName),
		ManagerID:   c.ID,
	}
	if err := internalProject.NewRepository(r.dbContext, r.l).InsertProjectByStruct(ctx, project); err != nil {
		return err
	}
	if err := r.dbContext.Get().NewQuery("SELECT id FROM projects WHERE number = {:number}").
		Bind(dbx.Params{"number": project.Number}).Row(&project.ID); err != nil {
		return err
	}

	if err := r.generateRandomStatuses(ctx, c, project); err != nil {
		return err
	}
	if err := r.assignConsultantsToProject(ctx, project, roles); err != nil {
		return err
	}
	if err := r.insertProjectTags(ctx, project, tags); err != nil {
		return err
	}
	allConsultants, _ := internalConsultant.NewRepository(r.dbContext, r.l).GetAllConsultants(ctx)
	if err := r.insertTimeEntries(ctx, project, allConsultants); err != nil {
		return err
	}

	return nil
}

// Internal helper function insertsConsultants into the consultants domain table
func (r *repository) insertConsultants(ctx context.Context, consultants []entity.Consultant, roles []string) error {
	for ci := range consultants {
		c := &consultants[ci]
		c.ProfilePicture = fmt.Sprintf("https://api.dicebear.com/7.x/lorelei/svg?seed=%s", url.QueryEscape(c.Email))

		if err := internalConsultant.NewRepository(r.dbContext, r.l).InsertConsultantByStruct(ctx, c); err != nil {
			return err
		}
		if err := r.dbContext.Get().NewQuery("SELECT id FROM consultants WHERE email = {:email}").
			Bind(dbx.Params{"email": c.Email}).Row(&c.ID); err != nil {
			return err
		}

		role := roles[rand.Intn(len(roles))]
		consultantRole := &entity.ConsultantRole{ConsultantID: c.ID, Role: role}
		if err := internalRole.NewRepository(r.dbContext, r.l).InsertConsultantRoleByStruct(ctx, *consultantRole); err != nil {
			return err
		}
	}
	return nil
}

// Internal generateProjectForConsultants method orchestrates generating x projects for each inserted consultant
func (r *repository) generateProjectsForConsultants(ctx context.Context, consultants []entity.Consultant, roles []string) error {
	for _, c := range consultants {
		projectCount := rand.Intn(maxProjectPerConsultant) + 1
		for pi := 0; pi < projectCount; pi++ {
			if err := r.generateProject(ctx, &c, pi+1, roles); err != nil {
				return err
			}
			r.l.Info(fmt.Sprintf("Seeded demo project %d for %s %s", pi, c.FirstName, c.LastName))
		}
	}
	return nil
}

// Internal Generate random statuses, generates random statuses for a given project
func (r *repository) generateRandomStatuses(ctx context.Context, c *entity.Consultant, project *entity.Project) error {
	allStatuses := []string{"planned", "active", "on-hold", "completed"}
	statusCount := rand.Intn(len(allStatuses)-1) + 1
	for i := 0; i < statusCount; i++ {
		st := allStatuses[rand.Intn(len(allStatuses))]
		status := &entity.ProjectStatus{
			Title:        st,
			Description:  fmt.Sprintf("Project %s is %s", project.Number, st),
			ConsultantID: c.ID,
			ProjectID:    project.ID,
		}
		if err := internalStatus.NewRepository(r.dbContext, r.l).InsertProjectStatusByStruct(ctx, status); err != nil {
			return err
		}
	}
	return nil
}

// Internal assignConsultantsToProject method assigns consultants to project
func (r *repository) assignConsultantsToProject(ctx context.Context, project *entity.Project, roles []string) error {
	allConsultants, err := internalConsultant.NewRepository(r.dbContext, r.l).GetAllConsultants(ctx)
	if err != nil {
		return err
	}
	assignedCount := rand.Intn(len(allConsultants)-1) + 1
	for i := 0; i < assignedCount; i++ {
		consultant := allConsultants[rand.Intn(len(allConsultants))]
		cp := &entity.ProjectConsultant{
			ProjectID:    project.ID,
			ConsultantID: consultant.ID,
			Role:         roles[rand.Intn(len(roles))],
		}
		if err := internalConsultant.NewRepository(r.dbContext, r.l).InsertProjectConsultantByStruct(ctx, *cp); err != nil {
			return err
		}
	}
	return nil
}

// Internal insertProjectTags inserts tags for project
func (r *repository) insertProjectTags(ctx context.Context, project *entity.Project, tags []string) error {
	tagCount := rand.Intn(len(tags)) + 1
	for i := 0; i < tagCount; i++ {
		tag := &entity.ProjectTag{ProjectID: project.ID, Tag: tags[rand.Intn(len(tags))]}
		if err := internalTag.NewRepository(r.dbContext, r.l).InsertProjectTagByStruct(ctx, *tag); err != nil {
			return err
		}
	}
	return nil
}

// Internal insertTimeEntries method inserts varying time entries into project
func (r *repository) insertTimeEntries(ctx context.Context, project *entity.Project, allConsultants []entity.Consultant) error {
	// choose random number of consultants (1 up to allConsultants length, capped by Max)
	numConsultants := rand.Intn(utils.MinMax(1, len(allConsultants), 5)) + 1
	rand.Shuffle(len(allConsultants), func(i, j int) { allConsultants[i], allConsultants[j] = allConsultants[j], allConsultants[i] })
	chosenConsultants := allConsultants[:numConsultants]

	createEntries := func(entryType string, consultant entity.Consultant, count int) error {
		for i := 0; i < count; i++ {
			hours := float32(rand.Intn(4)+1) + rand.Float32()

			// pick a future date relative to project.StartDate
			daysIntoFuture := rand.Intn(30) // within 30 days after project start
			entryDate := project.StartDate.AddDate(0, 0, daysIntoFuture)

			entry := &entity.TimeEntry{
				Title:        fmt.Sprintf("%s #%d", entryType, i+1),
				Description:  fmt.Sprintf("%s work for %s", entryType, project.Name),
				Hours:        hours,
				Type:         entryType,
				ConsultantID: consultant.ID,
				ProjectID:    project.ID,
				EntryDate:    entryDate,
			}
			if err := time.NewRepository(r.dbContext, r.l).InsertTimeEntryByStruct(ctx, entry); err != nil {
				return err
			}
		}
		return nil
	}

	// generate entries for each consultant in the chosen subset
	for _, consultant := range chosenConsultants {
		// random number of entries per consultant
		// debitCount := rand.Intn(utils.MinMax(1, 10, 5)) + 1  // up to 10 debits, default 5
		// creditCount := rand.Intn(utils.MinMax(1, 10, 5)) + 1 // up to 10 credits, default 5

		if err := createEntries("debit", consultant, 5); err != nil {
			return err
		}
		if err := createEntries("credit", consultant, 5); err != nil {
			return err
		}
	}

	return nil
}
