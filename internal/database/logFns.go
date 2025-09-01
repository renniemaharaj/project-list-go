package database

import (
	"context"
	"database/sql"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/renniemaharaj/grouplogs/pkg/logger"
)

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger *logger.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.SuccessF("DB query %s completed successfully in %d ms", sql, t.Milliseconds())
		} else {
			logger.ErrorF("DB query %s failed after %d ms: %v", sql, t.Milliseconds(), err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger *logger.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.SuccessF("DB exec %s completed successfully in %d ms", sql, t.Milliseconds())
		} else {
			logger.ErrorF("DB exec %s failed after %d ms: %v", sql, t.Milliseconds(), err)
		}
	}
}
