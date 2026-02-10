package database_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/narendhupati/bluejay-cms/internal/database"
)

func TestInitDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := database.InitDB(database.Config{Path: dbPath})
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer database.Close(db)

	if err := db.Ping(); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

func TestInitDB_InvalidPath(t *testing.T) {
	_, err := database.InitDB(database.Config{Path: "/nonexistent/dir/test.db"})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestRunMigrations(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := database.InitDB(database.Config{Path: dbPath})
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer database.Close(db)

	migrationsPath := findMigrationsDir(t)
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		t.Fatalf("RunMigrations failed: %v", err)
	}

	// Verify tables exist
	tables := []string{
		"settings", "admin_users", "product_categories", "blog_categories",
		"blog_authors", "industries", "partner_tiers", "whitepaper_topics",
		"products", "product_specs", "product_images", "product_features",
		"product_certifications", "product_downloads",
	}

	for _, table := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
		if err != nil {
			t.Errorf("failed to check table %s: %v", table, err)
		}
		if count == 0 {
			t.Errorf("table %s not found after migrations", table)
		}
	}
}

func TestRunMigrations_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := database.InitDB(database.Config{Path: dbPath})
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer database.Close(db)

	migrationsPath := findMigrationsDir(t)

	// Run twice - should not error
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		t.Fatalf("first RunMigrations failed: %v", err)
	}
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		t.Fatalf("second RunMigrations failed: %v", err)
	}
}

func findMigrationsDir(t *testing.T) string {
	t.Helper()
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
	t.Skip("could not find migrations directory")
	return ""
}
