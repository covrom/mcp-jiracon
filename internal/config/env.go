// Package config provides env-var parsing primitives and shared config types.
//
// Direct port of Python's src/mcp_atlassian/utils/env.py and config.py.
package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Truthy values for IsEnvTruthy / IsEnvExtendedTruthy.
var truthyValues = map[string]bool{
	"true": true, "1": true, "yes": true, "y": true, "on": true,
}

// IsEnvTruthy reports whether the named env var is set to a truthy
// value: true, 1, yes (case-insensitive). If the var is unset/empty,
// returns the parsed default (or false when no default given).
func IsEnvTruthy(name string, defaultVal ...string) bool {
	v := os.Getenv(name)
	if v == "" {
		if len(defaultVal) > 0 {
			return strings.EqualFold(defaultVal[0], "true")
		}
		return false
	}
	return truthyValues[strings.ToLower(v)]
}

// IsEnvExtendedTruthy accepts the wider set true/1/yes/y/on.
func IsEnvExtendedTruthy(name string, defaultVal ...string) bool {
	v := os.Getenv(name)
	if v == "" {
		if len(defaultVal) > 0 {
			return strings.EqualFold(defaultVal[0], "true")
		}
		return false
	}
	return truthyValues[strings.ToLower(v)]
}

// IsEnvSSLVerify reports whether TLS verification should be enabled.
// Defaults to true. Only explicit false/0/no disables.
func IsEnvSSLVerify(name string, defaultVal ...string) bool {
	v := os.Getenv(name)
	if v == "" {
		if len(defaultVal) > 0 {
			return !truthyValues[strings.ToLower(defaultVal[0])]
		}
		return true
	}
	lower := strings.ToLower(v)
	if lower == "false" || lower == "0" || lower == "no" {
		return false
	}
	return true
}

// GetIntEnv parses an integer env var. Returns defaultVal if unset or
// unparseable.
func GetIntEnv(name string, defaultVal int) int {
	v := os.Getenv(name)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		slog.Warn("invalid int env var", "name", name, "value", v, "default", defaultVal)
		return defaultVal
	}
	return n
}

// GetFloatEnv parses a float env var. Returns defaultVal if unset or
// unparseable.
func GetFloatEnv(name string, defaultVal float64) float64 {
	v := os.Getenv(name)
	if v == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		slog.Warn("invalid float env var", "name", name, "value", v, "default", defaultVal)
		return defaultVal
	}
	return f
}

// GetCustomHeaders parses a KEY=VALUE,KEY=VALUE string into a map.
// Splits on the FIRST "=" only so values may contain "=".
func GetCustomHeaders(name string) map[string]string {
	v := os.Getenv(name)
	if v == "" {
		return nil
	}
	out := map[string]string{}
	for _, pair := range strings.Split(v, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		idx := strings.Index(pair, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(pair[:idx])
		val := strings.TrimSpace(pair[idx+1:])
		if key == "" {
			continue
		}
		out[key] = val
	}
	return out
}

// GetHeaderNames parses a comma-separated list of header names and
// returns a deduplicated slice preserving the first-seen casing.
func GetHeaderNames(name string) []string {
	v := os.Getenv(name)
	if v == "" {
		return nil
	}
	seen := map[string]struct{}{}
	out := []string{}
	for _, h := range strings.Split(v, ",") {
		h = strings.TrimSpace(h)
		if h == "" {
			continue
		}
		lower := strings.ToLower(h)
		if _, ok := seen[lower]; ok {
			continue
		}
		seen[lower] = struct{}{}
		out = append(out, h)
	}
	return out
}
