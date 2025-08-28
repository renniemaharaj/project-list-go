package project

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/cache"
	"github.com/renniemaharaj/project-list-go/internal/database"
	"github.com/renniemaharaj/project-list-go/internal/entity"
)

var (
	projectLogger = logger.New().Prefix("Projects Router")
)

// ProjectHandler router, chi routing
func ProjectHandler(r chi.Router) {
	r.Get("/page/{pageNumber}", GetAllProjectIDSByPage)
	r.Get("/one/{projectID}", GetProjectsByID)
	r.Get("/search/{searchQuery}", GetProjectsBySearchQuery)
}

// Uses a search query to get matching projects
func GetProjectsBySearchQuery(w http.ResponseWriter, r *http.Request) {
	searchQuery := chi.URLParam(r, "searchQuery")
	if searchQuery == "" {
		http.Error(w, "search query required", http.StatusBadRequest)
		return
	}

	projects, err := cache.Use("projects:search:"+searchQuery, func() ([]int, error) {
		return NewRepository(database.Automatic, projectLogger).GetProjectIDSBySearchQuery(r.Context(), searchQuery)
	})

	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		projectLogger.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(projects)
}

// GetAllProjectIDSByPage returns projects paginated by page number
func GetAllProjectIDSByPage(w http.ResponseWriter, r *http.Request) {
	pageNumberStr := chi.URLParam(r, "pageNumber")
	if pageNumberStr == "" {
		http.Error(w, "pageNumber is required", http.StatusBadRequest)
		projectLogger.Error("pageNumber was missing from request")
		return
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber < 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	const pageSize = 10
	offset := pageNumber * pageSize

	projects, err := cache.Use("projects:page:"+pageNumberStr, func() ([]int, error) {
		return NewRepository(database.Automatic, projectLogger).GetProjectIDSByPage(r.Context(), pageSize, offset)
	})

	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		projectLogger.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(projects)
}

// GetProjectsByID returns a single project by ID
func GetProjectsByID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectID")
	if projectIDStr == "" {
		http.Error(w, "projectID is required", http.StatusBadRequest)
		projectLogger.Error("projectID was missing from request")
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid projectID", http.StatusBadRequest)
		projectLogger.Fatal(err)
		return
	}

	project, err := cache.Use("projects:one:"+projectIDStr, func() (*entity.Project, error) {
		return NewRepository(database.Automatic, projectLogger).GetProjectDataByID(r.Context(), projectID)
	})

	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		projectLogger.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(project)
}
