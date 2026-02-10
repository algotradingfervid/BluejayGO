// Package database provides SQLite database initialization, configuration,
// and connection management for the Bluejay CMS application.
//
// This package handles:
//   - Database connection setup with optimized SQLite pragmas
//   - Connection pool configuration for single-writer SQLite usage
//   - Database migration execution and rollback
//   - WAL (Write-Ahead Logging) mode configuration for better concurrency
//
// The package is designed to work with modernc.org/sqlite, a pure-Go SQLite
// driver that doesn't require CGO, making builds simpler and more portable.
package database

import (
	"database/sql" // Standard library SQL interface for database operations
	"fmt"          // String formatting for error messages and DSN construction

	// modernc.org/sqlite is a pure-Go SQLite driver (no CGO required)
	// Imported with blank identifier to register the "sqlite" driver
	_ "modernc.org/sqlite"
)

// Config holds the configuration parameters for database initialization.
// Currently only contains the file path, but can be extended with additional
// settings like connection pool sizes, timeouts, or pragma configurations.
type Config struct {
	// Path is the filesystem path to the SQLite database file.
	// If the file doesn't exist, SQLite will create it automatically.
	// Use ":memory:" for an in-memory database (useful for testing).
	Path string
}

// InitDB initializes and configures a SQLite database connection with
// production-ready settings optimized for the CMS use case.
//
// The function sets up the database with the following configurations:
//   - WAL (Write-Ahead Logging) mode for better concurrency
//   - Foreign key constraints enabled for referential integrity
//   - Busy timeout of 5 seconds to handle concurrent access
//   - NORMAL synchronous mode for balanced safety and performance
//   - Cache size of 2000 pages (typically ~8MB) for better query performance
//   - Single connection pool to respect SQLite's single-writer limitation
//
// Parameters:
//   - cfg: Configuration containing the database file path
//
// Returns:
//   - *sql.DB: Configured database connection ready for use
//   - error: Any error encountered during initialization
//
// Example usage:
//
//	db, err := InitDB(Config{Path: "./data/cms.db"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer Close(db)
func InitDB(cfg Config) (*sql.DB, error) {
	// Construct the Data Source Name (DSN) with SQLite-specific pragmas.
	// These pragmas are connection-level settings that configure SQLite's behavior:
	//
	// _journal_mode=WAL: Enables Write-Ahead Logging mode, which allows
	//   readers and writers to operate concurrently. WAL is recommended for
	//   web applications with moderate write activity.
	//
	// _busy_timeout=5000: Sets a 5-second timeout when the database is locked.
	//   SQLite will retry the operation for up to 5 seconds before returning
	//   a SQLITE_BUSY error. This helps handle concurrent access gracefully.
	//
	// _foreign_keys=on: Enables foreign key constraint checking. SQLite
	//   doesn't enable this by default for backward compatibility, but it's
	//   essential for maintaining referential integrity.
	//
	// _synchronous=NORMAL: Balances durability and performance. FULL would
	//   be safer but slower; NORMAL is sufficient for most applications when
	//   combined with WAL mode.
	//
	// _cache_size=2000: Sets the page cache to 2000 pages. With SQLite's
	//   default 4KB page size, this provides ~8MB of cache, improving query
	//   performance for frequently accessed data.
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on&_synchronous=NORMAL&_cache_size=2000", cfg.Path)

	// Open the database connection using the "sqlite" driver registered by
	// modernc.org/sqlite. This doesn't actually connect yet, just creates
	// the sql.DB instance.
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for SQLite's single-writer limitation.
	// SQLite only allows one concurrent writer, so we limit the pool to a
	// single connection to avoid "database is locked" errors.
	//
	// SetMaxOpenConns(1): Only one connection can be open at a time. This
	//   prevents connection pool contention and ensures serialized writes.
	db.SetMaxOpenConns(1)

	// SetMaxIdleConns(1): Keep one connection alive in the pool. This avoids
	//   the overhead of repeatedly opening/closing the database file.
	db.SetMaxIdleConns(1)

	// SetConnMaxLifetime(0): Connections never expire due to age. Since we
	//   only have one connection and it's long-lived, there's no benefit to
	//   rotating it.
	db.SetConnMaxLifetime(0)

	// SetConnMaxIdleTime(0): Idle connections are never closed. Again, with
	//   a single long-lived connection, we want to keep it open to avoid
	//   reconnection overhead.
	db.SetConnMaxIdleTime(0)

	// Execute a PRAGMA statement to ensure foreign keys are enabled.
	// While we set this in the DSN, explicitly executing it provides an
	// extra verification that the setting took effect. This is a safety
	// measure to catch any driver issues or configuration problems.
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Ping the database to verify the connection is actually working.
	// This forces the driver to establish a real connection and execute
	// a simple query, catching any connection issues early.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Close gracefully closes the database connection and releases all resources.
//
// This function should be called when shutting down the application to ensure:
//   - All pending transactions are completed or rolled back
//   - File handles are properly released
//   - WAL checkpoint is executed (writes buffered changes to main DB file)
//   - No database locks are left hanging
//
// It's safe to call Close on an already-closed database; subsequent calls
// will return an error but won't cause panics.
//
// Parameters:
//   - db: The database connection to close
//
// Returns:
//   - error: Any error encountered during closing (usually nil)
//
// Example usage:
//
//	db, _ := InitDB(Config{Path: "./data/cms.db"})
//	defer Close(db) // Ensures cleanup even if the program panics
func Close(db *sql.DB) error {
	return db.Close()
}
