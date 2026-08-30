// Package models contains the Go-typed representation of Atlassian
// REST responses. Direct port of Python's src/mcp_atlassian/models/.
//
// Every model implements FromAPIResponse (mutates receiver) and
// ToSimplifiedDict (returns map[string]any).
package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/constants"
)

// APIResponse is the contract every model implements.
type APIResponse interface {
	FromAPIResponse(data map[string]any) error
	ToSimplifiedDict() map[string]any
}

// FormatTimestamp converts an ISO 8601 timestamp to the canonical
// "YYYY-MM-DD HH:MM:SS TZ" form. Returns "" when input is empty;
// returns the original string when parsing fails.
func FormatTimestamp(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02T15:04:05-0700",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC().Format("2006-01-02 15:04:05 UTC")
		}
	}
	return s
}

// IsValidTimestamp returns true when s is a parseable timestamp.
func IsValidTimestamp(s string) bool {
	if s == "" {
		return false
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		if _, err := time.Parse(layout, s); err == nil {
			return true
		}
	}
	return false
}

// StringOrDefault returns v when non-empty, otherwise def.
func StringOrDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// MapString returns the string at key, or "" when missing.
func MapString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// MapAny returns the value at key, or nil when missing.
func MapAny(m map[string]any, key string) any {
	if v, ok := m[key]; ok {
		return v
	}
	return nil
}

// MapMap returns the value at key as a map, or nil when missing.
func MapMap(m map[string]any, key string) map[string]any {
	if v, ok := m[key]; ok {
		if mm, ok := v.(map[string]any); ok {
			return mm
		}
	}
	return nil
}

// MapSlice returns the value at key as a slice, or nil when missing.
func MapSlice(m map[string]any, key string) []any {
	if v, ok := m[key]; ok {
		if s, ok := v.([]any); ok {
			return s
		}
	}
	return nil
}

// MapInt returns the int64 at key, or 0 when missing.
func MapInt(m map[string]any, key string) int64 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int64:
		return x
	case int:
		return int64(x)
	case string:
		n, _ := strconv.ParseInt(x, 10, 64)
		return n
	}
	return 0
}

// MapBool returns the bool at key, or false when missing.
func MapBool(m map[string]any, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// AsID coerces a value to a string ID.
func AsID(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case nil:
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// JoinStrings joins s with sep, ignoring empty strings.
func JoinStrings(s []string, sep string) string {
	out := []string{}
	for _, x := range s {
		if x != "" {
			out = append(out, x)
		}
	}
	return strings.Join(out, sep)
}

var _ = constants.Unknown
var _ = constants.Unassigned
