// Package constants holds the default values used across all Jira
// and Confluence model conversions. Direct port of Python's
// src/mcp_atlassian/models/constants.py.
package constants

const (
	EmptyString = ""
	Unknown     = "Unknown"
	Unassigned  = "Unassigned"
	NoneValue   = "None"

	JiraDefaultID       = "0"
	JiraDefaultKey      = "UNKNOWN-0"
	JiraDefaultProject  = "0"
	ConfluenceDefaultID = "0"
	DefaultTimestamp    = "1970-01-01T00:00:00.000+0000"
)
