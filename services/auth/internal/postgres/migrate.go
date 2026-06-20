package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(
	ctx context.Context,
	pool *pgxpool.Pool,
	migrationsDir string,
) error {
	if err := ensureSchemaMigrationsTable(ctx, pool); err != nil {
		return err
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf(
			"postgres: read migrations dir: %w",
			err,
		)
	}

	var files []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		files = append(files, entry.Name())
	}

	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(file, ".sql")

		applied, err := isMigrationApplied(
			ctx,
			pool,
			version,
		)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		path := filepath.Join(
			migrationsDir,
			file,
		)

		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf(
				"postgres: read migration %s: %w",
				file,
				err,
			)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf(
				"postgres: begin migration %s: %w",
				version,
				err,
			)
		}

		if _, err := tx.Exec(
			ctx,
			string(sql),
		); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf(
				"postgres: apply migration %s: %w",
				version,
				err,
			)
		}

		if _, err := tx.Exec(
			ctx,
			`
			INSERT INTO schema_migrations (version)
			VALUES ($1)
			`,
			version,
		); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf(
				"postgres: record migration %s: %w",
				version,
				err,
			)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf(
				"postgres: commit migration %s: %w",
				version,
				err,
			)
		}
	}

	return nil
}

func ensureSchemaMigrationsTable(
	ctx context.Context,
	pool *pgxpool.Pool,
) error {
	_, err := pool.Exec(
		ctx,
		`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
		`,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: ensure schema_migrations table: %w",
			err,
		)
	}

	return nil
}

func isMigrationApplied(
	ctx context.Context,
	pool *pgxpool.Pool,
	version string,
) (bool, error) {
	var exists bool

	err := pool.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM schema_migrations
			WHERE version = $1
		)
		`,
		version,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf(
			"postgres: check migration %s: %w",
			version,
			err,
		)
	}

	return exists, nil
}
