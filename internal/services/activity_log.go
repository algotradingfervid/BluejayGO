package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type ActivityLogService struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewActivityLogService(queries *sqlc.Queries, logger *slog.Logger) *ActivityLogService {
	return &ActivityLogService{queries: queries, logger: logger}
}

// Log records an activity. It's fire-and-forget: errors are logged but not returned
// to avoid disrupting the main operation.
func (s *ActivityLogService) Log(ctx context.Context, userID int64, action, resourceType string, resourceID int64, resourceTitle, description string) {
	err := s.queries.CreateActivityLog(ctx, sqlc.CreateActivityLogParams{
		UserID:        sql.NullInt64{Int64: userID, Valid: userID > 0},
		Action:        action,
		ResourceType:  resourceType,
		ResourceID:    sql.NullInt64{Int64: resourceID, Valid: resourceID > 0},
		ResourceTitle: sql.NullString{String: resourceTitle, Valid: resourceTitle != ""},
		Description:   description,
	})
	if err != nil {
		s.logger.Error("failed to log activity", "error", err, "action", action, "resource", resourceType)
	}
}

// LogSimple logs an activity without a specific resource (e.g., login/logout).
func (s *ActivityLogService) LogSimple(ctx context.Context, userID int64, action, description string) {
	s.Log(ctx, userID, action, "system", 0, "", description)
}

// LogF is a convenience wrapper with fmt.Sprintf for description.
func (s *ActivityLogService) LogF(ctx context.Context, userID int64, action, resourceType string, resourceID int64, resourceTitle, format string, args ...interface{}) {
	s.Log(ctx, userID, action, resourceType, resourceID, resourceTitle, fmt.Sprintf(format, args...))
}
