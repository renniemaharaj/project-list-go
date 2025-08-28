package database

import (
	"fmt"
	"os"
	"sync"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/joho/godotenv"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

// DBContext encapsulates a singleton DB connection
type DBContext struct {
	envVar string
	once   sync.Once
	DBX    *dbx.DB
	err    error
}

var (
	databaseLogger = logger.New().Prefix("Database Logger")

	// Automatic will automatically pick the best DB available
	Automatic = newAutomaticInstance()

	// Explicit DB instances
	Developer  = newInstance("DEV_POSTGRES_DSN")
	Docker     = newInstance("DOCKER_POSTGRES_DSN")
	Production = newInstance("PROD_POSTGRES_DSN")
	Testing    = newInstance("TEST_POSTGRES_DSN")
)

// newInstance creates a new dbInstance tied to an env var
func newInstance(envVar string) *DBContext {
	return &DBContext{envVar: envVar}
}

// newAutomaticInstance creates an adaptive DB context that tries multiple fallbacks
func newAutomaticInstance() *DBContext {
	return &DBContext{
		envVar: "ADAPTIVE", // fake label for clarity
	}
}

// Get automatically chooses database by trying Production -> Docker -> Developer
func (dbContext *DBContext) Get() (*dbx.DB, error) {
	// Special behavior for Adaptive
	if dbContext.envVar == "ADAPTIVE" {
		// Try Production
		if db, err := Production.GetManual(); err == nil {
			return db, nil
		}
		// Try Docker
		if db, err := Docker.GetManual(); err == nil {
			return db, nil
		}
		// Try Developer
		if db, err := Developer.GetManual(); err == nil {
			return db, nil
		}

		// Nothing worked -> panic
		panic("Adaptive DBContext: no valid DSN found (PROD_POSTGRES_DSN, DOCKER_POSTGRES_DSN, DEV_POSTGRES_DSN)")
	}

	// Default behavior (non-adaptive)
	return dbContext.GetManual()
}

// useDefault is the standard Get (factored out for clarity)
func (dbContext *DBContext) GetManual() (*dbx.DB, error) {
	dbContext.once.Do(func() {
		_ = godotenv.Load() // load .env safely

		dsn := os.Getenv(dbContext.envVar)
		if dsn == "" {
			dbContext.err = fmt.Errorf("%s not set", dbContext.envVar)
			return
		}

		dbContext.DBX, dbContext.err = dbx.Open("postgres", dsn)
		if dbContext.err != nil {
			dbContext.err = fmt.Errorf("failed to open DB (%s): %w", dbContext.envVar, dbContext.err)
		}
	})

	return dbContext.DBX, dbContext.err
}
