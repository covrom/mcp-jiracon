// Package utils provides cross-cutting helpers: logging, IO, media
// detection, etc. Direct port of Python's src/mcp_atlassian/utils/.

package utils

import (
	"log/slog"
	"os"
	"strings"
)

// MaskSensitive keeps the first keep and last keep characters and
// masks the middle with asterisks. Empty input returns "Not Provided".
func MaskSensitive(value string, keep ...int) string {
	if value == "" {
		return "Not Provided"
	}
	k := 4
	if len(keep) > 0 && keep[0] >= 0 {
		k = keep[0]
	}
	if len(value) <= 2*k {
		return strings.Repeat("*", len(value))
	}
	return value[:k] + strings.Repeat("*", len(value)-2*k) + value[len(value)-k:]
}

// SetupLogging configures the root slog logger with the given level.
// Output goes to stderr, format is text.
func SetupLogging(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// LogLevelFromEnv reads MCP_LOG_LEVEL or LOG_LEVEL, returns INFO if unset.
func LogLevelFromEnv() slog.Level {
	switch strings.ToUpper(os.Getenv("MCP_LOG_LEVEL")) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	}
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	}
	return slog.LevelInfo
}

// MaskSessionHeaders masks credentials in outbound HTTP headers for
// safe logging. Mirrors utils.logging.get_masked_session_headers.
func MaskSessionHeaders(headers map[string]string) map[string]string {
	out := make(map[string]string, len(headers))
	for k, v := range headers {
		lower := strings.ToLower(k)
		switch lower {
		case "authorization":
			if strings.HasPrefix(v, "Basic ") {
				out[k] = "Basic " + MaskSensitive(strings.TrimPrefix(v, "Basic "))
			} else if strings.HasPrefix(v, "Bearer ") {
				out[k] = "Bearer " + MaskSensitive(strings.TrimPrefix(v, "Bearer "))
			} else {
				out[k] = MaskSensitive(v)
			}
		case "cookie", "set-cookie", "proxy-authorization":
			out[k] = MaskSensitive(v)
		default:
			out[k] = v
		}
	}
	return out
}
