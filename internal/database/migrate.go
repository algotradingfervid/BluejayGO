package database

import (
	"database/sql" // Standard library SQL interface for database operations
	"fmt"          // String formatting for error messages and file path construction

	// golang-migrate/migrate/v4 is the main migration engine that orchestrates
	// the execution of migration files in order, tracks which migrations have
	// been applied, and handles version management.
	"github.com/golang-migrate/migrate/v4"

	// golang-migrate/migrate/v4/database/sqlite provides the SQLite-specific
	// database driver for the migration engine. It knows how to create the
	// schema_migrations table and execute SQLite-compatible SQL statements.
	"github.com/golang-migrate/migrate/v4/database/sqlite"

	// golang-migrate/migrate/v4/source/file enables reading migration files
	// from the filesystem. Imported with blank identifier to register the
	// "file://" source driver with the migration engine.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations executes all pending database migrations in order, bringing
// the database schema up to the latest version.
//
// This function:
//   - Reads migration files from the specified directory
//   - Checks which migrations have already been applied (via schema_migrations table)
//   - Executes only the pending "up" migrations in sequential order
//   - Updates the schema_migrations table to track applied versions
//   - Is idempotent: safe to run multiple times (skips already-applied migrations)
//
// Migration files should follow the naming convention:
//   {version}_{description}.up.sql   - for upgrades
//   {version}_{description}.down.sql - for rollbacks
//
// For example:
//   001_create_users_table.up.sql
//   001_create_users_table.down.sql
//   002_add_email_index.up.sql
//   002_add_email_index.down.sql
//
// Parameters:
//   - db: Active database connection (must be already initialized)
//   - migrationsPath: Filesystem path to the directory containing migration files
//     (e.g., "./db/migrations" or "/app/migrations")
//
// Returns:
//   - error: Any error encountered during migration, or nil if successful.
//     Returns nil if there are no pending migrations (ErrNoChange is suppressed)
//
// Example usage:
//
//	db, _ := InitDB(Config{Path: "./data/cms.db"})
//	err := RunMigrations(db, "./db/migrations")
//	if err != nil {
//	    log.Fatalf("Migration failed: %v", err)
//	}
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Create a SQLite-specific migration driver instance from the existing
	// database connection. This driver wraps our sql.DB connection and provides
	// the migrate library with SQLite-specific functionality like creating the
	// schema_migrations tracking table and executing migrations within transactions.
	//
	// Using WithInstance (instead of WithConnection) allows us to reuse our
	// pre-configured database connection with all its pragmas and settings.
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create a new migration instance that combines:
	//   1. Source: file:// URL pointing to migration files on disk
	//   2. Database driver: SQLite driver we just created
	//   3. Database name: "sqlite" (used for driver identification)
	//
	// The file:// prefix is required by the file source driver to know
	// it should read from the filesystem (vs. other sources like s3://, etc.)
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Execute all pending "up" migrations in order from oldest to newest.
	// The migrate library:
	//   1. Reads the current version from the schema_migrations table
	//   2. Finds all migration files with versions higher than current
	//   3. Executes each .up.sql file in order
	//   4. Updates schema_migrations after each successful migration
	//
	// If err is migrate.ErrNoChange, it means all migrations are already
	// applied, which is a success case (not an error). We suppress this
	// specific error to make the function idempotent.
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// RollbackMigration rolls back all applied database migrations, returning
// the database to an empty state (version 0).
//
// WARNING: This function is DESTRUCTIVE and will:
//   - Execute all "down" migrations in reverse order (newest to oldest)
//   - Drop tables, remove columns, delete data as specified in .down.sql files
//   - Leave the database in an empty or initial state
//   - Should only be used in development/testing, NEVER in production
//
// This is useful for:
//   - Completely resetting a development database
//   - Testing migration rollback logic
//   - Cleaning up test databases in automated tests
//
// For production, use targeted rollbacks (migrate to specific version) rather
// than rolling back everything. The golang-migrate library supports this via
// m.Migrate(targetVersion), but this function intentionally uses Down() to
// demonstrate the nuclear option.
//
// Parameters:
//   - db: Active database connection (must be already initialized)
//   - migrationsPath: Filesystem path to the directory containing migration files
//
// Returns:
//   - error: Any error encountered during rollback, or nil if successful.
//     Returns nil if there are no migrations to roll back (ErrNoChange is suppressed)
//
// Example usage:
//
//	// Development/testing only!
//	db, _ := InitDB(Config{Path: "./data/test.db"})
//	err := RollbackMigration(db, "./db/migrations")
//	if err != nil {
//	    log.Fatalf("Rollback failed: %v", err)
//	}
func RollbackMigration(db *sql.DB, migrationsPath string) error {
	// Create a SQLite-specific migration driver instance from the existing
	// database connection. This is identical to the setup in RunMigrations
	// because we need the same driver functionality for rollbacks.
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create a new migration instance with the file source and SQLite driver.
	// This setup is identical to RunMigrations, but we'll call Down() instead
	// of Up() to execute rollback migrations.
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Execute all "down" migrations in reverse order (newest to oldest).
	// The migrate library:
	//   1. Reads the current version from schema_migrations
	//   2. Finds all applied migrations from newest to oldest
	//   3. Executes each .down.sql file in reverse order
	//   4. Updates schema_migrations after each successful rollback
	//   5. Stops when it reaches version 0 (no migrations applied)
	//
	// If err is migrate.ErrNoChange, it means there are no migrations to
	// roll back (database is already at version 0). We suppress this error
	// to make the function idempotent.
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}
