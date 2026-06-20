package postgres

import (
	"context"
	"fmt"

	"github.com/cmiski/sawako/services/auth/internal/authentication"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(
	pool *pgxpool.Pool,
) authentication.TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

func (m *TransactionManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"postgres: begin transaction: %w",
			err,
		)
	}

	txCtx := context.WithValue(
		ctx,
		txKey{},
		tx,
	)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf(
				"postgres: rollback transaction: %v (original: %w)",
				rbErr,
				err,
			)
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"postgres: commit transaction: %w",
			err,
		)
	}

	return nil
}

func querier(
	ctx context.Context,
	pool *pgxpool.Pool,
) pgxQuerier {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if ok {
		return tx
	}

	return pool
}

type pgxQuerier interface {
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)

	QueryRow(
		ctx context.Context,
		sql string,
		args ...any,
	) pgx.Row
}
