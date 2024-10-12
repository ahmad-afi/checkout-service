package utils

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type SQLTransaction struct {
	pg *sqlx.DB
}

func NewSQLTransaction(pg *sqlx.DB) SQLTransaction {
	return SQLTransaction{pg}
}

func (s *SQLTransaction) WrapperTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := s.pg.Beginx()
	if err != nil {
		return
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
