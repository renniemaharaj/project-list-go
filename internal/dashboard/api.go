package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/cache"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	internalMeta "github.com/renniemaharaj/project-list-go/internal/meta"
	internalProject "github.com/renniemaharaj/project-list-go/internal/project"
	"github.com/renniemaharaj/project-list-go/internal/utils"
)

// Metrics struct for dashboard workers
type metrics struct {
	Completed            int
	Active               int
	Idle                 int
	OutOfBudget          int
	TotalDebit           float64
	TotalCredit          float64
	CreditOverDebitRatio float64 // valid only if HasRatio
	HasRatio             bool
}

var (
	dashboardLogger = logger.New().Prefix("Dash Router")
	// Prevents duplicate recomputes of the same dashboard concurrently.
	dashboardBuildLock sync.Mutex
)

// Internal cache projects function to reuse dashboard-fetched projects
func proactivelyCacheProjectData(projects []entity.Project) {
	for _, p := range projects {
		// Cache each project under its own key
		_, _ = cache.Use(fmt.Sprintf("projects:one:%d", p.ID), func() (entity.Project, error) {
			// Return the already-fetched value; no DB call here.
			return p, nil
		})
	}
}

func Dashboard(r chi.Router) {
	r.Get("/", GetMetricsDashboard)
}

func GetMetricsDashboard(w http.ResponseWriter, r *http.Request) {
	// Guard: ensures that waiting clients must hit cache instead of unnecessarily computing
	dashboardBuildLock.Lock()
	defer dashboardBuildLock.Unlock()

	ctx := r.Context()
	// We wrap dashboard compute in a use cache interface which auto caches return values
	dashboardMetrics, err := cache.Use("metrics_dashboard", func() (entity.MetricsDashboard, error) {

		// 1) Fetch project IDs (cheap)
		// repository does not offer get all projects in complete (more performant to return ids)
		projectIDS, err := internalProject.NewRepository(database.Automatic, dashboardLogger).GetAllProjectIDS(ctx)
		if err != nil {
			return entity.MetricsDashboard{}, err
		}
		if len(projectIDS) == 0 {
			return entity.MetricsDashboard{}, nil
		}

		projectMetas, projects, err := internalMeta.NewRepository(database.Automatic, dashboardLogger).GetProjectsMetaByProjectIDS(ctx, projectIDS, true)
		// 2) Fetch project rows in bulk (cheap relative to per-project meta)
		// projects, err := internalProject.NewRepository(database.Automatic, dashboardLogger).GetProjectsDataByIDS(ctx, projectIDS)
		if err != nil {
			return entity.MetricsDashboard{}, err
		}
		// 2a) Proactively cache each project record (async fire-and-forget)
		go proactivelyCacheProjectData(projects)

		const idleThreshold = 7 * 24 * time.Hour
		now := time.Now()
		// 3) Pre-compute ending soon from projects (no need to involve workers)
		result := entity.MetricsDashboard{
			Projects: len(projects),
		}
		for _, p := range projects {
			if p.EndDate.After(now) && p.EndDate.Before(now.Add(idleThreshold)) {
				result.EndingSoon++
			}
		}

		// 4) Fan-out: fetch meta + compute per-project partials in a worker pool
		jobs := make(chan int)        // project IDs
		results := make(chan metrics) // per-project partials
		errs := make(chan error, 1)   // first error wins

		// Prevents < 4 and uncaps
		workerCount := utils.MinMax(200, runtime.NumCPU()*4, 0)
		worker := func(ctx context.Context) {
			for id := range jobs {
				// Use cache.Use around meta fetch
				meta, metaError := cache.Use(fmt.Sprintf("projects:meta:%d", id), func() (entity.ProjectMeta, error) {
					if pm, ok := projectMetas[id]; ok {
						return pm, nil
					}

					pm, err := internalMeta.NewRepository(database.Automatic, dashboardLogger).GetProjectMetaByProjectID(ctx, id)
					return *pm, err
				})
				if metaError != nil {
					errs <- metaError
					return
				}

				var p metrics
				// Latest status assumed at index 0 if present
				if len(meta.StatusHistory) > 0 {
					switch meta.StatusHistory[0].Title {
					case "completed":
						p.Completed = 1
					case "active":
						p.Active = 1
					}
					if now.Sub(meta.StatusHistory[0].DateCreated) > idleThreshold {
						p.Idle = 1
					}
				}

				// Debit/Credit & OOB
				var debit, credit float64
				for _, t := range meta.TimeEntries {
					switch t.Type {
					case "debit":
						debit += float64(t.Hours)
					case "credit":
						credit += float64(t.Hours)
					}
				}
				p.TotalDebit = debit
				p.TotalCredit = credit
				if credit > debit {
					p.OutOfBudget = 1
				}
				if debit > 0 {
					p.HasRatio = true
					p.CreditOverDebitRatio = credit / debit
				}

				// send partial
				select {
				case results <- p:
				case <-ctx.Done():
					return
				}
			}
		}

		// Start workers
		for i := 0; i < workerCount; i++ {
			go worker(ctx)
		}

		// Feed jobs
		go func() {
			defer close(jobs)
			for _, p := range projects {
				select {
				case jobs <- p.ID:
				case <-ctx.Done():
					return
				}
			}
		}()

		// Aggregate (single goroutine, no locks)
		var (
			totalRatios float64
			ratioCount  int
			received    int
		)

		for received < len(projects) {
			select {
			case err := <-errs:
				// Drain workers by cancel if you wrap ctx with context.WithCancel here.
				return entity.MetricsDashboard{}, err
			case p := <-results:
				received++
				result.Completed += p.Completed
				result.Active += p.Active
				result.Idle += p.Idle
				result.OutOfBudget += p.OutOfBudget
				result.TotalDebit += p.TotalDebit
				result.TotalCredit += p.TotalCredit

				if p.HasRatio {
					totalRatios += p.CreditOverDebitRatio
					ratioCount++
				}
			case <-ctx.Done():
				return entity.MetricsDashboard{}, ctx.Err()
			}
		}

		if ratioCount > 0 {
			result.AverageCreditOverDebit = float32(totalRatios / float64(ratioCount))
		}

		return result, nil
	})

	if err != nil {
		http.Error(w, "Failed to fetch dashboard metrics", http.StatusInternalServerError)
		dashboardLogger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dashboardMetrics)
}
