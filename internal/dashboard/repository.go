package dashboard

import (
	"github.com/renniemaharaj/grouplogs/pkg/logger"
	"github.com/renniemaharaj/project-list-go/internal/database"
)

type Repository interface {
}

type repository struct {
	dbContext *database.DBContext
	l         *logger.Logger
}

func NewRepository(dbContext *database.DBContext, _l *logger.Logger) Repository {
	return &repository{dbContext, _l}
}
