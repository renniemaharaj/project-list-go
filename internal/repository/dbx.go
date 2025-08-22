package repository

import (
	"fmt"
	"os"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/joho/godotenv"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

// The repository struct
type repository struct {
	db *dbx.DB
	l  *logger.Logger
}

// Internal new repository
func newRepository(db *dbx.DB) *repository {
	return &repository{db: db, l: logger.New().Prefix("Repository")}
}

// NewRepository opens a Postgres database and returns a repository instance
func NewRepository() (*repository, error) {
	if dbx, err := GETDBX(); err == nil {
		return newRepository(dbx), nil
	}

	return nil, fmt.Errorf("couldn't load repository")
}

// Creates and returns an ozzo dbx connection
func GETDBX() (*dbx.DB, error) {
	l := logger.New().Prefix("Repository")
	// Load .env file
	if err := godotenv.Load(); err != nil {
		l.Fatal(err)
	}

	if dsn := os.Getenv("POSTGRE_DSN"); dsn != "" {
		db, err := dbx.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return nil, fmt.Errorf("couldn't open database connection")
}
