package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	"github.com/renniemaharaj/project-list-go/internal/repository"
)

var (
	projectsLogger = logger.New().Prefix("Projects Router")
)

// Projects router, chi routing
func Projects(r chi.Router) {
	r.Get("/all/page/{pageNumber}", GetAllProjects)
	r.Get("/one/{projectID}", GetProjectsByID)
	r.Get("/meta/one/{projectID}", GetProjectMetaData)
	r.Get("/metas/all", GetProjectMetas)
}

// GetProjectMetas gets all project metas
func GetProjectMetas(w http.ResponseWriter, r *http.Request) {
	// Initialize repository
	repos, err := repository.Get()
	if err != nil {
		http.Error(w, "Failed to initialize repository", 500)
		projectsLogger.Fatal(err)
		return
	}

	// Fetch projects
	projects, err := repos.GetAllProjectsDesc(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch projects", 500)
		projectsLogger.Fatal(err)
		return
	}

	projectMetas := []entity.ProjectMetaData{}
	for _, p := range projects {
		projectMeta, err := repos.GetProjectMetaByProjectID(r.Context(), p.ID)
		if err != nil {
			http.Error(w, "Failed to fetch project metas", 500)
			projectsLogger.Fatal(err)
			return
		}
		projectMetas = append(projectMetas, *projectMeta)
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projectMetas)
}

// GetProjectMetaData returns the meta data for a project by ID
func GetProjectMetaData(w http.ResponseWriter, r *http.Request) {
	// Get projectID from URL
	projectIDStr := chi.URLParam(r, "projectID")
	if projectIDStr == "" {
		http.Error(w, "projectID is required", http.StatusBadRequest)
		projectsLogger.Error("projectID was missing from request")
		return
	}

	// Convert string to int
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid projectID", http.StatusBadRequest)
		projectsLogger.Fatal(err)
		return
	}

	// Initialize repository
	repos, err := repository.Get()
	if err != nil {
		http.Error(w, "Failed to initialize repository", 500)
		projectsLogger.Fatal(err)
		return
	}

	// Fetch projectMeta by ID
	projectMeta, err := repos.GetProjectMetaByProjectID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
		return
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projectMeta)
}

// GetAllProjects returns projects paginated by page number
func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	// Initialize repository
	repos, err := repository.Get()
	if err != nil {
		http.Error(w, "Failed to initialize repository", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
		return
	}

	// Parse page number from URL
	pageNumberStr := chi.URLParam(r, "pageNumber")
	if pageNumberStr == "" {
		http.Error(w, "pageNumber is required", http.StatusBadRequest)
		projectsLogger.Error("pageNumber was missing from request")
		return
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber < 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	// Pagination parameters
	const pageSize = 5
	offset := pageNumber * pageSize

	// Fetch projects with limit & offset
	projects, err := repos.GetProjectsPage(r.Context(), pageSize, offset)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
		return
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GetProjectsByID returns a single project by ID
func GetProjectsByID(w http.ResponseWriter, r *http.Request) {
	// Get projectID from URL
	projectIDStr := chi.URLParam(r, "projectID")
	if projectIDStr == "" {
		http.Error(w, "projectID is required", http.StatusBadRequest)
		projectsLogger.Error("projectID was missing from request")
		return
	}

	// Convert string to int
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid projectID", http.StatusBadRequest)
		projectsLogger.Fatal(err)
		return
	}

	// Initialize repository
	repos, err := repository.Get()
	if err != nil {
		http.Error(w, "Failed to initialize repository", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
		return
	}

	// Fetch project by ID
	project, err := repos.GetProjectByID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
		return
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
