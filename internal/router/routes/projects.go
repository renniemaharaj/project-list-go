package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/cache"
	"github.com/renniemaharaj/project-list-go/internal/entity"
	"github.com/renniemaharaj/project-list-go/internal/repository"
)

var (
	projectsLogger = logger.New().Prefix("Projects Router")
)

// Projects router, chi routing
func Projects(r chi.Router) {
	r.Get("/page/{pageNumber}", GetAllProjectIDSByPage)
	r.Get("/one/{projectID}", GetProjectsByID)
	r.Get("/meta/{projectID}", GetProjectMetaData)
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
		repos, err := repository.Get()
		if err != nil {
			return nil, err
		}
		return repos.GetProjectIDSBySearchQuery(r.Context(), searchQuery)
	})

	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(projects)
}

// GetProjectMetaData returns the meta data for a project by ID
func GetProjectMetaData(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectID")
	if projectIDStr == "" {
		http.Error(w, "projectID is required", http.StatusBadRequest)
		projectsLogger.Error("projectID was missing from request")
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid projectID", http.StatusBadRequest)
		projectsLogger.Fatal(err)
		return
	}

	projectMeta, err := cache.Use("projects:meta:"+projectIDStr, func() (*entity.ProjectMetaData, error) {
		repos, err := repository.Get()
		if err != nil {
			return nil, err
		}
		return repos.GetProjectMetaByProjectID(r.Context(), projectID)
	})

	if err != nil {
		http.Error(w, "Failed to fetch project meta", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(projectMeta)
}

// GetAllProjectIDSByPage returns projects paginated by page number
func GetAllProjectIDSByPage(w http.ResponseWriter, r *http.Request) {
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

	const pageSize = 10
	offset := pageNumber * pageSize

	projects, err := cache.Use("projects:page:"+pageNumberStr, func() ([]int, error) {
		repos, err := repository.Get()
		if err != nil {
			return nil, err
		}
		return repos.GetProjectIDSByPage(r.Context(), pageSize, offset)
	})

	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
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
		projectsLogger.Error("projectID was missing from request")
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid projectID", http.StatusBadRequest)
		projectsLogger.Fatal(err)
		return
	}

	project, err := cache.Use("projects:one:"+projectIDStr, func() (*entity.Project, error) {
		repos, err := repository.Get()
		if err != nil {
			return nil, err
		}
		return repos.GetProjectDataByID(r.Context(), projectID)
	})

	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		projectsLogger.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(project)
}
