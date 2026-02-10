// Package services provides business logic layer for the Bluejay CMS application.
// It encapsulates core functionality including activity logging, caching, product
// management, and file uploads. Services act as an intermediary between HTTP handlers
// and the database layer, implementing reusable business operations and enforcing
// application rules.
package services

import (
	// Standard library imports for context handling, database operations, formatting, and logging
	"context"      // Provides context for request-scoped values and cancellation
	"database/sql" // Provides SQL null types for handling nullable database fields
	"fmt"          // Used for string formatting in log descriptions
	"log/slog"     // Structured logging for recording errors and diagnostics

	// Internal application imports
	"github.com/narendhupati/bluejay-cms/db/sqlc" // Generated database query code from sqlc
)

// ActivityLogService provides functionality for recording user and system activities
// throughout the application. It tracks actions performed by users such as creating,
// updating, or deleting resources, enabling audit trails and activity monitoring.
type ActivityLogService struct {
	queries *sqlc.Queries // Database query interface for persisting activity logs
	logger  *slog.Logger  // Structured logger for recording service-level errors
}

// NewActivityLogService creates and initializes a new ActivityLogService instance.
// This service is responsible for recording all user and system activities in the
// application for auditing and monitoring purposes.
//
// Parameters:
//   - queries: Database query interface from sqlc for executing activity log operations
//   - logger: Structured logger for recording errors when activity logging fails
//
// Returns:
//   - *ActivityLogService: Initialized service ready to log activities
func NewActivityLogService(queries *sqlc.Queries, logger *slog.Logger) *ActivityLogService {
	return &ActivityLogService{queries: queries, logger: logger}
}

// Log records a complete activity event with full resource details. This is the core
// logging function used throughout the application to track user actions. It operates
// in a fire-and-forget manner: errors are logged but not returned to the caller to
// avoid disrupting the primary operation that triggered the activity.
//
// The function converts zero values to SQL NULL values for optional fields, ensuring
// proper database representation when user ID, resource ID, or resource title are not
// applicable to the logged action.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout handling
//   - userID: ID of the user performing the action (0 for system actions)
//   - action: Type of action performed (e.g., "created", "updated", "deleted")
//   - resourceType: Category of resource affected (e.g., "product", "category", "user")
//   - resourceID: ID of the specific resource affected (0 if not applicable)
//   - resourceTitle: Human-readable title of the resource (empty string if not applicable)
//   - description: Detailed description of what occurred
func (s *ActivityLogService) Log(ctx context.Context, userID int64, action, resourceType string, resourceID int64, resourceTitle, description string) {
	// Create the activity log entry with nullable fields properly handled
	err := s.queries.CreateActivityLog(ctx, sqlc.CreateActivityLogParams{
		UserID:        sql.NullInt64{Int64: userID, Valid: userID > 0},           // NULL if userID is 0 (system action)
		Action:        action,                                                      // Required: type of action performed
		ResourceType:  resourceType,                                                // Required: category of resource
		ResourceID:    sql.NullInt64{Int64: resourceID, Valid: resourceID > 0},    // NULL if resourceID is 0
		ResourceTitle: sql.NullString{String: resourceTitle, Valid: resourceTitle != ""}, // NULL if empty string
		Description:   description,                                                 // Required: detailed description
	})

	// Log errors internally but don't propagate them to maintain fire-and-forget behavior.
	// This ensures that activity logging failures don't disrupt the main application flow.
	if err != nil {
		s.logger.Error("failed to log activity", "error", err, "action", action, "resource", resourceType)
	}
}

// LogSimple logs a simplified activity event without a specific resource association.
// This is a convenience function for recording system-level or user-level actions that
// don't relate to a particular resource (e.g., user login, user logout, session timeout).
//
// The function automatically sets resourceType to "system", resourceID to 0, and
// resourceTitle to empty string, making it ideal for authentication and session events.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout handling
//   - userID: ID of the user performing the action (0 for system actions)
//   - action: Type of action performed (e.g., "login", "logout", "session_expired")
//   - description: Detailed description of what occurred
func (s *ActivityLogService) LogSimple(ctx context.Context, userID int64, action, description string) {
	// Delegate to the full Log function with system-level defaults
	s.Log(ctx, userID, action, "system", 0, "", description)
}

// LogF is a convenience wrapper around Log that accepts a format string and arguments
// for building the description dynamically. This is useful when activity descriptions
// need to include dynamic values like IDs, names, or timestamps.
//
// The function uses fmt.Sprintf internally to construct the description, allowing for
// printf-style formatting with type-safe arguments.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout handling
//   - userID: ID of the user performing the action (0 for system actions)
//   - action: Type of action performed (e.g., "created", "updated", "deleted")
//   - resourceType: Category of resource affected (e.g., "product", "category", "user")
//   - resourceID: ID of the specific resource affected (0 if not applicable)
//   - resourceTitle: Human-readable title of the resource (empty string if not applicable)
//   - format: Printf-style format string for the description
//   - args: Arguments to be interpolated into the format string
func (s *ActivityLogService) LogF(ctx context.Context, userID int64, action, resourceType string, resourceID int64, resourceTitle, format string, args ...interface{}) {
	// Build the formatted description and delegate to the full Log function
	s.Log(ctx, userID, action, resourceType, resourceID, resourceTitle, fmt.Sprintf(format, args...))
}
