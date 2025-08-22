package repository

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// UseTransaction provides an interface for transactions with defer rollback, error handling
// and transaction commit
func (r *repository) UseTransaction(ctx context.Context, consume func(tx *dbx.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // safe no-op if committed

	if err := consume(tx); err != nil {
		return err // rollback happens automatically
	}

	return tx.Commit()
}
