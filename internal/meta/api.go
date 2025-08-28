package meta

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
	metaLogger = logger.New().Prefix("Meta Logger")
)

func Meta(r chi.Router) {
	r.Get("/meta", GetProjectMetaData)
}

// GetProjectMetaData returns the meta data for a project by ID
func GetProjectMetaData(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectID")
	if projectIDStr == "" {
		http.Error(w, "projectID is required", http.StatusBadRequest)
		metaLogger.Error("projectID was missing from request")
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid projectID", http.StatusBadRequest)
		metaLogger.Fatal(err)
		return
	}

	projectMeta, err := cache.Use("projects:meta:"+projectIDStr, func() (entity.ProjectMeta, error) {
		md, err := NewRepository(database.Automatic, metaLogger).GetProjectMetaByProjectID(r.Context(), projectID)
		if err != nil {
			metaLogger.Fatal(err)
			return entity.ProjectMeta{}, err
		}
		return *md, err
	})

	if err != nil {
		http.Error(w, "Failed to fetch project meta", http.StatusInternalServerError)
		metaLogger.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(projectMeta)
}
