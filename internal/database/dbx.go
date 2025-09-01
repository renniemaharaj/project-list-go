package database

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/joho/godotenv"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

// DBContext encapsulates a singleton DB connection
type DBContext struct {
	envVar string
	dbx    *dbx.DB
	err    error
}

var (
	databaseLogger = logger.NewLogger().Prefix("Database Logger")

	// Automatically resolves to the best available DB trying Production -> Docker -> Developer
	Automatic = newInstance("AUTOMATIC")

	// Explicit DB instances in order of priority
	Production = newInstance("PROD_POSTGRES_DSN")
	Docker     = newInstance("DOCKER_POSTGRES_DSN")
	Developer  = newInstance("DEV_POSTGRES_DSN")
	Testing    = newInstance("TEST_POSTGRES_DSN")
)

// newInstance creates a new dbInstance tied to an env var
func newInstance(envVar string) *DBContext {
	return &DBContext{envVar: envVar}
}

// Get returns the dbx db of the db context
func (dbContext *DBContext) Get() *dbx.DB {
	return dbContext.dbx
}

// EnvVar returns the environment variable name used by this db context
func (dbContext *DBContext) EnvVar() string {
	return dbContext.envVar
}

// Resolve automatically chooses database by trying Production -> Docker -> Developer
func (dbContext *DBContext) Resolve() (*dbx.DB, error) {
	if dbContext.envVar == "AUTOMATIC" {
		// Try in order, but assign result back into Automatic (self)
		candidates := []*DBContext{Production, Docker, Developer}

		for _, candidate := range candidates {
			if db, err := candidate.GetManual(); err == nil {
				// Mutate Automatic to hold this DB
				dbContext.dbx = candidate.dbx
				dbContext.err = candidate.err
				dbContext.envVar = candidate.envVar
				databaseLogger.Warning(fmt.Sprintf("Automatic bound to %s", candidate.envVar))
				return db, nil
			}
		}

		panic("Adaptive DBContext: no valid DSN found (PROD_POSTGRES_DSN, DOCKER_POSTGRES_DSN, DEV_POSTGRES_DSN)")
	}

	return dbContext.GetManual()
}

// GetManual tries to open and test a DB connection for this context
func (dbContext *DBContext) GetManual() (*dbx.DB, error) {
	if dbContext.dbx == nil {
		_ = godotenv.Load()

		dsn := os.Getenv(dbContext.envVar)
		if dsn == "" {
			dbContext.err = fmt.Errorf("%s not set", dbContext.envVar)
			databaseLogger.Warning(dbContext.err.Error())
			return nil, dbContext.err
		}

		db, err := dbx.Open("postgres", dsn)
		if err != nil {
			dbContext.err = fmt.Errorf("failed to open DB (%s): %w", dbContext.envVar, err)
			databaseLogger.Fatal(dbContext.err)
			return nil, dbContext.err
		}

		// Test connection
		if pingErr := db.DB().Ping(); pingErr != nil {
			dbContext.err = fmt.Errorf("failed to ping DB (%s): %w", dbContext.envVar, pingErr)
			databaseLogger.Warning(dbContext.err.Error())
			return nil, dbContext.err
		}

		dbContext.dbx = db
		dbContext.err = nil
	}

	return dbContext.dbx, dbContext.err
}
