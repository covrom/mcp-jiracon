package utils

import (
	"fmt"
	"runtime/debug"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/pkg/version"
)

// DefaultUserAgent returns the User-Agent string sent on outbound
// requests. Includes the build-time version. Mirrors utils.user_agent.
func DefaultUserAgent() string {
	v := version.Version
	if v == "" {
		// Best-effort version detection from debug.ReadBuildInfo.
		if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" {
			v = bi.Main.Version
		} else {
			v = "0.0.0"
		}
	}
	return fmt.Sprintf("mcp-atlassian-go/%s", v)
}
