package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/tigorlazuardi/claw/lib/otel"
)

//go:embed *.sql
var embedMigrations embed.FS

// MigrateOption configures migration behavior
type MigrateOption func(*MigrateConfig)

// MigrateConfig holds migration configuration
type MigrateConfig struct {
	// Reset rollbacks all migrations and reapplies them.
	// WARNING: This option will DROP ALL DATA in the database.
	// DO NOT USE IN PRODUCTION ENVIRONMENTS.
	Reset bool
	// FS is the filesystem to use for migrations. If nil, uses embedded migrations.
	FS fs.FS
}

// WithReset enables the reset option.
// WARNING: This will rollback all migrations and reapply them, causing ALL DATA LOSS.
// This option should NEVER be used in production environments.
func WithReset(reset bool) MigrateOption {
	return func(config *MigrateConfig) {
		config.Reset = reset
	}
}

// WithFS sets a custom filesystem for migrations.
// If not provided, the embedded migrations filesystem will be used.
func WithFS(filesystem fs.FS) MigrateOption {
	return func(config *MigrateConfig) {
		config.FS = filesystem
	}
}

// Migrate executes database migrations using goose.
// It accepts a context, database connection, and optional MigrateOption functions.
//
// Available options:
//   - WithReset: WARNING - Rollbacks all migrations and reapplies them (DATA LOSS)
//   - WithFS: Use custom filesystem for migrations (default: embedded)
//
// Example usage:
//
//	err := Migrate(ctx, db)
//	err := Migrate(ctx, db, WithReset(true)) // WARNING: DATA LOSS
//	err := Migrate(ctx, db, WithFS(customFS))
func Migrate(ctx context.Context, db *sql.DB, options ...MigrateOption) error {
	ctx, span := otel.Start(ctx)
	defer span.End()

	config := &MigrateConfig{
		Reset: false,
		FS:    embedMigrations,
	}

	// Apply options
	for _, option := range options {
		option(config)
	}

	// Set the filesystem for goose
	goose.SetBaseFS(config.FS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Handle reset option
	if config.Reset {
		// Check if goose_db_version table exists before attempting reset
		var tableName string
		err := db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name='goose_db_version'").Scan(&tableName)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to check if goose_db_version table exists: %w", err)
		}

		// Only reset if the table exists (database has been initialized)
		if tableName == "goose_db_version" {
			if err := goose.ResetContext(ctx, db, "."); err != nil {
				return fmt.Errorf("failed to reset migrations: %w", err)
			}
		}
	}

	// Run migrations
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
