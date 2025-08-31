package meta

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/cache"
	"github.com/renniemaharaj/project-list-go/internal/database"
)

var (
	metaLogger = logger.New().Prefix("Meta Logger")
)

func Meta(r chi.Router) {
	r.Get("/{projectID}", GetProjectMetaByProjectID)
}

// GetProjectMetaByProjectID returns the meta data for a project by ID
func GetProjectMetaByProjectID(w http.ResponseWriter, r *http.Request) {
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

	projectMeta, err := cache.Use("projects:meta:"+projectIDStr, func() (*ProjectMeta, error) {
		md, err := NewService(NewRepository(database.Automatic, metaLogger), metaLogger).GetProjectMetaByProjectID(r.Context(), projectID)
		if err != nil {
			metaLogger.Fatal(err)
			return &ProjectMeta{}, err
		}
		return md, err
	})

	if err != nil {
		http.Error(w, "Failed to fetch project meta", http.StatusInternalServerError)
		metaLogger.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(projectMeta)
}
