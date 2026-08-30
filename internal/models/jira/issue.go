package jira

import (
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/constants"
)

// JiraIssue represents a Jira issue. The Fields map holds all
// issue-specific data (summary, status, assignee, project, etc.).
type JiraIssue struct {
	ID        string         `json:"id"`
	Key       string         `json:"key"`
	Self      string         `json:"self"`
	Fields    map[string]any `json:"fields"`
	Changelog map[string]any `json:"changelog,omitempty"`
}

// FromAPIResponse populates the issue from the REST response.
func (i *JiraIssue) FromAPIResponse(data map[string]any) error {
	i.ID = models.AsID(data["id"])
	i.Key = models.MapString(data, "key")
	i.Self = models.MapString(data, "self")
	if f := models.MapMap(data, "fields"); f != nil {
		i.Fields = f
	} else {
		i.Fields = map[string]any{}
	}
	i.Changelog = models.MapMap(data, "changelog")
	return nil
}

// ToSimplifiedDict returns a flat, JSON-friendly map with the most
// commonly-used fields extracted.
func (i *JiraIssue) ToSimplifiedDict() map[string]any {
	out := map[string]any{
		"id":   i.ID,
		"key":  i.Key,
		"self": i.Self,
	}
	if i.Fields == nil {
		return out
	}
	// Pull a few common top-level fields out so the LLM sees them
	// without navigating the nested Fields map.
	for _, k := range []string{
		"summary", "description", "status", "priority", "issuetype",
		"assignee", "reporter", "creator", "labels", "created", "updated",
		"duedate", "resolutiondate", "resolution", "parent", "subtasks",
		"fixVersions", "components", "project", "watches", "timeoriginalestimate",
		"timeestimate", "timespent", "workratio", "issuelinks",
	} {
		if v, ok := i.Fields[k]; ok && v != nil {
			out[k] = v
		}
	}
	if len(out) == 3 {
		out["fields"] = i.Fields
	}
	return out
}

// Summary returns the issue's summary field, or "" when missing.
func (i *JiraIssue) Summary() string {
	s, _ := i.Fields["summary"].(string)
	return s
}

// ProjectKey returns the issue's project key (e.g. "PROJ"), or "".
func (i *JiraIssue) ProjectKey() string {
	if proj, ok := i.Fields["project"].(map[string]any); ok {
		s, _ := proj["key"].(string)
		return s
	}
	return ""
}

// Init keeps the constants import live for go vet.
var _ = constants.JiraDefaultID
