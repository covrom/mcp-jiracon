package utils

import (
	"log/slog"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
)

// ClampLimit applies the ATLASSIAN_MAX_PAGINATION_LIMIT ceiling to a
// requested limit. If the env var is unset or <=0, the cap is disabled.
func ClampLimit(requested int, context string) int {
	if requested <= 0 {
		return requested
	}
	cap := config.GetIntEnv("ATLASSIAN_MAX_PAGINATION_LIMIT", 0)
	if cap <= 0 {
		return requested
	}
	if requested > cap {
		slog.Info("pagination limit clamped",
			"context", context, "requested", requested, "cap", cap)
		return cap
	}
	return requested
}
