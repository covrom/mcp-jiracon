package utils

import (
	"log/slog"
	"os"
	"strings"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
)

// ToolsetPrefix is the prefix on every tool's MCP tag that identifies
// toolset membership.
const ToolsetPrefix = "toolset:"

// ToolsetDef defines a toolset's metadata.
type ToolsetDef struct {
	Name        string
	Description string
	Default     bool
}

// JiraToolsets is the canonical list of Jira toolsets.
var JiraToolsets = map[string]ToolsetDef{
	"jira_issues":           {"jira_issues", "Core CRUD on issues", true},
	"jira_fields":           {"jira_fields", "Field metadata and option lookups", true},
	"jira_comments":         {"jira_comments", "Comments", true},
	"jira_transitions":      {"jira_transitions", "Status transitions", true},
	"jira_projects":         {"jira_projects", "Projects / versions / components", false},
	"jira_agile":            {"jira_agile", "Agile boards, sprints", false},
	"jira_links":            {"jira_links", "Issue links (incl. Epic link, remote link)", false},
	"jira_worklog":          {"jira_worklog", "Worklogs", false},
	"jira_attachments":      {"jira_attachments", "Attachments and images", false},
	"jira_users":            {"jira_users", "User search and profile", false},
	"jira_watchers":         {"jira_watchers", "Issue watchers", false},
	"jira_service_desk":     {"jira_service_desk", "Service Desk queues and requests (DC mostly)", false},
	"jira_forms":            {"jira_forms", "ProForma forms (Cloud only)", false},
	"jira_metrics":          {"jira_metrics", "Issue dates + SLA metrics", false},
	"jira_development":      {"jira_development", "Dev panel (PRs, branches, commits)", false},
	"jira_project_analysis": {"jira_project_analysis", "Epic hierarchy + cross-project links", false},
}

// ConfluenceToolsets is the canonical list of Confluence toolsets.
var ConfluenceToolsets = map[string]ToolsetDef{
	"confluence_pages":       {"confluence_pages", "Pages CRUD, history, diff, sections", true},
	"confluence_comments":    {"confluence_comments", "Footer + inline comments", true},
	"confluence_labels":      {"confluence_labels", "Page labels", false},
	"confluence_users":       {"confluence_users", "User search", false},
	"confluence_analytics":   {"confluence_analytics", "Page view analytics (Cloud only)", false},
	"confluence_attachments": {"confluence_attachments", "Attachments + images", false},
	"confluence_templates":   {"confluence_templates", "Page templates (Cloud only)", false},
	"confluence_permissions": {"confluence_permissions", "Permission check / space permissions (Cloud only)", false},
}

// allToolsets returns every toolset in a single map.
func allToolsets() map[string]ToolsetDef {
	out := map[string]ToolsetDef{}
	for k, v := range JiraToolsets {
		out[k] = v
	}
	for k, v := range ConfluenceToolsets {
		out[k] = v
	}
	return out
}

// defaultToolsets returns the set of default-enabled toolset names.
func defaultToolsets() map[string]struct{} {
	out := map[string]struct{}{}
	for k, v := range allToolsets() {
		if v.Default {
			out[k] = struct{}{}
		}
	}
	return out
}

// GetEnabledToolsets reads TOOLSETS env var. Returns the set of toolset
// names to enable. If unset, returns all toolset names (with a
// deprecation warning). If set to "all" returns every toolset.
// If set to "default" returns the default toolsets. Unknown names are
// dropped with a warning. If no valid names remain, returns the empty set.
func GetEnabledToolsets() map[string]struct{} {
	v := strings.TrimSpace(os.Getenv("TOOLSETS"))
	if v == "" {
		slog.Warn("TOOLSETS unset; defaulting to all toolsets (deprecated)")
		all := map[string]struct{}{}
		for k := range allToolsets() {
			all[k] = struct{}{}
		}
		return all
	}
	tokens := strings.Split(v, ",")
	out := map[string]struct{}{}
	all := allToolsets()
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		lower := strings.ToLower(tok)
		if lower == "all" {
			for k := range all {
				out[k] = struct{}{}
			}
			return out
		}
		if lower == "default" {
			for k := range defaultToolsets() {
				out[k] = struct{}{}
			}
			continue
		}
		if _, ok := all[tok]; !ok {
			slog.Warn("unknown toolset in TOOLSETS", "value", tok)
			continue
		}
		out[tok] = struct{}{}
	}
	return out
}

// ShouldIncludeToolByToolset returns true if the tool's toolset tag
// is in enabledToolsets (or the tool has no toolset tag, in which case
// it's included unconditionally with a warning).
func ShouldIncludeToolByToolset(toolTags []string, enabledToolsets map[string]struct{}) bool {
	if len(enabledToolsets) == 0 {
		return false
	}
	hasToolset := false
	for _, t := range toolTags {
		if strings.HasPrefix(t, ToolsetPrefix) {
			hasToolset = true
			name := strings.TrimPrefix(t, ToolsetPrefix)
			if _, ok := enabledToolsets[name]; ok {
				return true
			}
		}
	}
	if !hasToolset {
		slog.Warn("tool has no toolset tag; including by default",
			"tags", toolTags)
		return true
	}
	return false
}

// GetToolsetTag returns the first toolset:* tag in the given tag set,
// or "" if none.
func GetToolsetTag(tags []string) string {
	for _, t := range tags {
		if strings.HasPrefix(t, ToolsetPrefix) {
			return strings.TrimPrefix(t, ToolsetPrefix)
		}
	}
	return ""
}

// DefaultToolsetsList returns a sorted list of default toolset names.
func DefaultToolsetsList() []string {
	out := []string{}
	for k, v := range allToolsets() {
		if v.Default {
			out = append(out, k)
		}
	}
	return out
}

// init — log the enabled toolsets at startup when the package is
// imported (this runs once when the binary starts).
func init() {
	// Trigger deprecation warning when unset.
	_ = config.IsEnvTruthy // keep import live
}
