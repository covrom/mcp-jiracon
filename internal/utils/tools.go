package utils

import (
	"os"
	"strings"
)

// GetEnabledTools reads ENABLED_TOOLS (comma-separated) and returns
// the list of tool names that should be exposed. Returns nil when
// unset/empty (meaning "all tools").
func GetEnabledTools() []string {
	v := strings.TrimSpace(os.Getenv("ENABLED_TOOLS"))
	if v == "" {
		return nil
	}
	out := []string{}
	for _, name := range strings.Split(v, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}

// ShouldIncludeTool returns true if the named tool should be exposed
// given the ENABLED_TOOLS allowlist. When allowlist is nil (unset),
// all tools are included.
func ShouldIncludeTool(name string, enabled []string) bool {
	if len(enabled) == 0 {
		return true
	}
	for _, e := range enabled {
		if e == name {
			return true
		}
	}
	return false
}
