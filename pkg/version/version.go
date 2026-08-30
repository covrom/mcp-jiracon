// Package version exposes the build-time version string for the binary.
//
// The version is set via -ldflags at build time:
//
//	go build -ldflags "-X gitlab.services.mts.ru/rstsovan/mcp-atlassian/pkg/version.Version=v1.0.0"
package version

// Version is the build-time version string. Default "dev" for local builds.
var Version = "dev"
