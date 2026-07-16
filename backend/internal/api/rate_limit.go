package api

import (
	"context"
	"os"
	"strconv"
	"log/slog"

	"arishem/internal/db"
	"github.com/google/uuid"
)

// GetMaxRunsPerDay returns the configured limit from environment or default
func GetMaxRunsPerDay() int {
	val := os.Getenv("MAX_RUNS_PER_DAY")
	if val == "" {
		return 10 // Default fallback limit
	}
	limit, err := strconv.Atoi(val)
	if err != nil {
		return 10
	}
	return limit
}

// CheckDailyLimitExceeded returns true if the organization has exceeded the daily runs limit
func CheckDailyLimitExceeded(ctx context.Context, orgID string) (bool, error) {
	// If it's a test environment where db is mock/nil, skip limit
	if db.GetPool() == nil {
		return false, nil
	}

	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return false, err
	}

	limit := GetMaxRunsPerDay()

	// Query total created scans + llm pentests in the last 24 hours
	query := `
		SELECT COUNT(*) FROM (
			SELECT id FROM scans WHERE org_id = $1 AND created_at >= NOW() - INTERVAL '1 day'
			UNION ALL
			SELECT id FROM llm_pentest_runs WHERE org_id = $1 AND created_at >= NOW() - INTERVAL '1 day'
		) as total
	`

	var count int
	err = db.GetPool().QueryRow(ctx, query, orgUUID).Scan(&count)
	if err != nil {
		if err.Error() == "closed pool" {
			return false, nil
		}
		slog.Error("Failed to check daily limit", slog.String("org_id", orgID), slog.Any("error", err))
		return false, err
	}

	if count >= limit {
		slog.Warn("Daily run limit exceeded", slog.String("org_id", orgID), slog.Int("count", count), slog.Int("limit", limit))
		return true, nil
	}

	return false, nil
}
