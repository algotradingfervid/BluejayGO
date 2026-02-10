// Package testutil provides common testing utilities for the Bluejay CMS test suite.
// It includes helpers for setting up test databases, running migrations, and managing
// test fixtures. These utilities ensure consistent test environment setup across
// unit tests, integration tests, and end-to-end tests.
package testutil

import (
	// Standard library imports
	"database/sql"  // SQL database interface for working with SQLite connections
	"os"            // File system operations for temp directories and file cleanup
	"path/filepath" // Cross-platform path manipulation for finding migration files
	"runtime"       // Runtime information used to locate project root from caller location
	"testing"       // Go testing framework providing test helpers and cleanup

	// Project imports
	"github.com/narendhupati/bluejay-cms/db/sqlc"          // sqlc generated database queries
	"github.com/narendhupati/bluejay-cms/internal/database" // Database initialization and migration runner
)

// SetupTestDB creates a fully initialized test database for use in test cases.
//
// This function:
//  1. Creates a temporary directory and SQLite database file
//  2. Initializes the database connection using the standard database.InitDB
//  3. Locates the migration files relative to the project root
//  4. Runs all migrations to create the schema
//  5. Returns a sqlc.Queries instance for type-safe database operations
//
// The function uses a file-based SQLite database rather than in-memory (:memory:)
// because golang-migrate requires a file path to properly track migration state.
//
// Parameters:
//   - t: The testing.T instance for the current test, used for marking this as a helper
//     and for automatic cleanup via t.TempDir()
//
// Returns:
//   - *sql.DB: The initialized database connection
//   - *sqlc.Queries: A sqlc query object bound to the test database
//   - func(): A cleanup function that closes the database and removes the temp file.
//     This should be deferred immediately after calling SetupTestDB.
//
// Example usage:
//
//	func TestMyFeature(t *testing.T) {
//	    db, queries, cleanup := testutil.SetupTestDB(t)
//	    defer cleanup()
//	    // Use queries to interact with test database
//	}
func SetupTestDB(t *testing.T) (*sql.DB, *sqlc.Queries, func()) {
	// Mark this as a test helper so failures point to the caller, not this function
	t.Helper()

	// Create a temporary directory that will be automatically cleaned up when the test ends.
	// We use a file-based database instead of :memory: because golang-migrate needs a
	// real file path to track which migrations have been applied.
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Initialize the database connection using the production database initialization code.
	// This ensures tests use the same connection settings as production.
	db, err := database.InitDB(database.Config{Path: dbPath})
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}

	// Locate the migrations directory by walking up from this file to the project root.
	// This approach works regardless of where tests are run from.
	migrationsPath := findMigrationsPath()
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		// Ensure we close the database if migrations fail to avoid resource leaks
		db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Create a sqlc.Queries instance for type-safe database operations in tests
	queries := sqlc.New(db)

	// Define the cleanup function to close database and remove temp files
	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, queries, cleanup
}

// findMigrationsPath locates the database migrations directory by walking up
// the directory tree from the current file to find the project root.
//
// This function uses runtime.Caller to get the current file's location, then
// walks up the directory tree until it finds a directory containing go.mod,
// which marks the project root. It then returns the path to db/migrations
// relative to that root.
//
// This approach ensures tests can locate migrations regardless of:
//   - The current working directory when tests are run
//   - Whether tests are run from the project root or a subdirectory
//   - The test execution environment (IDE, command line, CI/CD)
//
// Returns:
//   - string: The absolute path to the db/migrations directory, or "db/migrations"
//     as a fallback if the project root cannot be found (which should only happen
//     in misconfigured test environments)
//
// Algorithm:
//  1. Get the path of this source file using runtime.Caller(0)
//  2. Start from the directory containing this file
//  3. Walk up directories checking for go.mod in each
//  4. When go.mod is found, that directory is the project root
//  5. Return {project_root}/db/migrations
func findMigrationsPath() string {
	// Get the file path of the currently executing source file (this file).
	// runtime.Caller(0) returns: pc, file, line, ok
	// We only need the filename, so we discard other return values with _
	_, filename, _, _ := runtime.Caller(0)

	// Start from the directory containing this source file
	dir := filepath.Dir(filename)

	// Walk up the directory tree until we find go.mod (project root marker)
	for {
		// Check if go.mod exists in the current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found the project root - return the migrations path relative to it
			return filepath.Join(dir, "db", "migrations")
		}

		// Move up one directory level
		parent := filepath.Dir(dir)

		// If parent equals dir, we've reached the filesystem root without finding go.mod
		if parent == dir {
			break
		}

		dir = parent
	}

	// Fallback path if we couldn't find the project root.
	// This relative path will only work if tests are run from the project root.
	return "db/migrations"
}
