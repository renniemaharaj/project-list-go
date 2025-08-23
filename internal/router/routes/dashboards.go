package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	"github.com/renniemaharaj/project-list-go/internal/repository"
)

var (
	dashboardLogger = logger.New().Prefix("Dash Router")
)

// Projects router, chi routing
func Dashboard(r chi.Router) {
	r.Get("/", GetDashboardMetrics)
}

func GetDashboardMetrics(w http.ResponseWriter, r *http.Request) {
	// Initialize repository
	repos, err := repository.NewRepository()
	if err != nil {
		http.Error(w, "Failed to initialize repository", 500)
		dashboardLogger.Fatal(err)
		return
	}

	dashboardMetrics := entity.DashboardMetrics{}
	// Fetch projects
	projects, err := repos.GetProjects(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch projects", 500)
		dashboardLogger.Fatal(err)
		return
	}
	dashboardMetrics.Projects = len(projects)

	// collect all metas & calculate ending soon
	projectMetas := []entity.ProjectMetaData{}
	for _, p := range projects {
		// Ending soon: project end date within next 7 days
		if p.EndDate.After(time.Now()) && p.EndDate.Before(time.Now().Add(7*24*time.Hour)) {
			dashboardMetrics.EndingSoon++
		}
		//collect metas
		projectMeta, err := repos.GetProjectMetaByProjectID(r.Context(), p.ID)
		if err != nil {
			http.Error(w, "Failed to fetch project metas", 500)
			dashboardLogger.Fatal(err)
			return
		}
		projectMetas = append(projectMetas, *projectMeta)
	}

	// idle threshold
	idleTime := 7 * 24 * time.Hour // consider idle if no status in last 7 days
	now := time.Now()

	var totalRatios float64
	var ratioCount int
	for _, meta := range projectMetas {
		// latest status assumed at index 0 (if sorted DESC by date or id)
		if len(meta.StatusHistory) > 0 {
			if meta.StatusHistory[0].Title == "completed" {
				dashboardMetrics.Completed++
			} else {
				dashboardMetrics.Active++
			}

			// idle check
			lastStatus := meta.StatusHistory[0]
			lastStatusDate := lastStatus.DateCreated
			if now.Sub(lastStatusDate) > idleTime {
				dashboardMetrics.Idle++
			}
		}

		// Debit / Credit
		var debit, credit float64
		for _, t := range meta.TimeEntries {
			switch t.Type {
			case "debit":
				debit += float64(t.Hours)
			case "credit":
				credit += float64(t.Hours)
			}
		}
		dashboardMetrics.TotalDebit += debit
		dashboardMetrics.TotalCredit += credit

		// increment out of budget
		if credit > debit {
			dashboardMetrics.OutOfBudget++
		}

		if debit > 0 {
			totalRatios += credit / debit
			ratioCount++
		}

	}

	// calculate average credit over debit
	if ratioCount > 0 {
		dashboardMetrics.AverageCreditOverDebit = float32(totalRatios / float64(ratioCount))
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboardMetrics)
}
