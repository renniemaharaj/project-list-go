package database

import (
	"context"
	"database/sql"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// QueryDBLogFunc returns a logging function that can be used to log SQL queries.
func QueryDBLogFunc() dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			databaseLogger.SuccessF("DB query %s completed successfully in %d ms", sql, t.Milliseconds())
		} else {
			databaseLogger.ErrorF("DB query %s failed after %d ms: %v", sql, t.Milliseconds(), err)
		}
	}
}

// ExecDBLogFunc returns a logging function that can be used to log SQL executions.
func ExecDBLogFunc() dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			databaseLogger.SuccessF("DB exec %s completed successfully in %d ms", sql, t.Milliseconds())
		} else {
			databaseLogger.ErrorF("DB exec %s failed after %d ms: %v", sql, t.Milliseconds(), err)
		}
	}
}
