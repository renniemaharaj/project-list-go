package repository

import (
	"fmt"
	"os"
	"sync"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/joho/godotenv"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

// repository struct with dbx singleton and logger
type repository struct {
	DB *dbx.DB
	l  *logger.Logger
}

var (
	singletonRepo *repository
	once          sync.Once
)

// Get returns a singleton repository instance
func Get() (*repository, error) {
	var err error
	once.Do(func() {
		db, dbErr := getSingletonDBX()
		if dbErr != nil {
			err = dbErr
			return
		}
		singletonRepo = &repository{DB: db, l: logger.New().Prefix("Repository")}
	})
	if singletonRepo == nil {
		return nil, fmt.Errorf("couldn't create repository: %w", err)
	}
	return singletonRepo, nil
}

// getSingletonDBX returns a singleton DB connection (Postgres)
func getSingletonDBX() (*dbx.DB, error) {
	// Load .env once
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	dsn := os.Getenv("POSTGRE_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("POSTGRE_DSN not set")
	}

	db, err := dbx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	return db, nil
}
