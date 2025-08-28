package database

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// UseTransaction provides an interface for transactions with defer rollback, error handling
// and transaction commit
func (dbContext *DBContext) UseTransaction(ctx context.Context, consume func(tx *dbx.Tx) error) error {
	_db, err := dbContext.Get()
	if err != nil {
		databaseLogger.Fatal(err)
	}

	tx, err := _db.WithContext(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // safe no-op if committed

	if err := consume(tx); err != nil {
		return err // rollback happens automatically
	}

	return tx.Commit()
}
