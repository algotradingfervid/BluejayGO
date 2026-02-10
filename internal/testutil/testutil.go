package testutil

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/database"
)

// SetupTestDB creates a temporary SQLite database with all migrations applied.
// Returns the db, queries, and a cleanup function.
func SetupTestDB(t *testing.T) (*sql.DB, *sqlc.Queries, func()) {
	t.Helper()

	// Create temp file for SQLite (in-memory doesn't work well with migrate)
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := database.InitDB(database.Config{Path: dbPath})
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}

	// Find migrations path relative to project root
	migrationsPath := findMigrationsPath()
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	queries := sqlc.New(db)
	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, queries, cleanup
}

func findMigrationsPath() string {
	// Walk up from the current file to find project root
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "db", "migrations")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Fallback
	return "db/migrations"
}
