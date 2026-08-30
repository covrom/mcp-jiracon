// Package jira provides Go-typed representations of Jira REST API
// responses. Direct port of Python's src/mcp_atlassian/models/jira/.
package jira

import (
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/constants"
)

// JiraUser represents a Jira user. The identifier field differs by
// deployment: Cloud uses accountId, DC uses name+key.
type JiraUser struct {
	AccountID    string            `json:"accountId,omitempty"`
	Username     string            `json:"username,omitempty"`
	Key          string            `json:"key,omitempty"`
	DisplayName  string            `json:"displayName,omitempty"`
	EmailAddress string            `json:"emailAddress,omitempty"`
	Active       bool              `json:"active,omitempty"`
	AvatarURLs   map[string]string `json:"avatarUrls,omitempty"`
	TimeZone     string            `json:"timeZone,omitempty"`
}

// FromAPIResponse populates u from a Jira user dict.
func (u *JiraUser) FromAPIResponse(data map[string]any) error {
	u.AccountID = models.MapString(data, "accountId")
	u.Username = models.MapString(data, "username")
	u.Key = models.MapString(data, "key")
	u.DisplayName = models.MapString(data, "displayName")
	u.EmailAddress = models.MapString(data, "emailAddress")
	u.Active = models.MapBool(data, "active")
	u.AvatarURLs = mapStringMap(data, "avatarUrls")
	u.TimeZone = models.MapString(data, "timeZone")
	return nil
}

// ToSimplifiedDict returns the JSON-serializable dict.
func (u *JiraUser) ToSimplifiedDict() map[string]any {
	out := map[string]any{}
	if u.AccountID != "" {
		out["accountId"] = u.AccountID
	}
	if u.Username != "" {
		out["username"] = u.Username
	}
	if u.Key != "" {
		out["key"] = u.Key
	}
	out["displayName"] = models.StringOrDefault(u.DisplayName, constants.Unknown)
	out["emailAddress"] = u.EmailAddress
	out["active"] = u.Active
	if len(u.AvatarURLs) > 0 {
		out["avatarUrls"] = u.AvatarURLs
	}
	if u.TimeZone != "" {
		out["timeZone"] = u.TimeZone
	}
	return out
}

// JiraStatus is a status with name + id; defaults to Unknown/0.
type JiraStatus struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FromAPIResponse populates from a status dict.
func (s *JiraStatus) FromAPIResponse(data map[string]any) error {
	s.ID = models.AsID(data["id"])
	s.Name = models.MapString(data, "name")
	if s.ID == "" {
		s.ID = constants.JiraDefaultID
	}
	if s.Name == "" {
		s.Name = constants.Unknown
	}
	return nil
}

// ToSimplifiedDict returns the JSON-serializable dict.
func (s *JiraStatus) ToSimplifiedDict() map[string]any {
	return map[string]any{"id": s.ID, "name": s.Name}
}

// JiraPriority mirrors a priority object.
type JiraPriority struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FromAPIResponse populates from a priority dict.
func (p *JiraPriority) FromAPIResponse(data map[string]any) error {
	p.ID = models.AsID(data["id"])
	p.Name = models.MapString(data, "name")
	if p.ID == "" {
		p.ID = constants.JiraDefaultID
	}
	if p.Name == "" {
		p.Name = constants.NoneValue
	}
	return nil
}

// ToSimplifiedDict returns the JSON-serializable dict.
func (p *JiraPriority) ToSimplifiedDict() map[string]any {
	return map[string]any{"id": p.ID, "name": p.Name}
}

// JiraIssueType is an issue type (Task, Bug, Story, etc).
type JiraIssueType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Subtask     bool   `json:"subtask,omitempty"`
}

// FromAPIResponse populates from an issue type dict.
func (t *JiraIssueType) FromAPIResponse(data map[string]any) error {
	t.ID = models.AsID(data["id"])
	t.Name = models.MapString(data, "name")
	t.Description = models.MapString(data, "description")
	t.Subtask = models.MapBool(data, "subtask")
	if t.ID == "" {
		t.ID = constants.JiraDefaultID
	}
	if t.Name == "" {
		t.Name = constants.Unknown
	}
	return nil
}

// ToSimplifiedDict returns the JSON-serializable dict.
func (t *JiraIssueType) ToSimplifiedDict() map[string]any {
	out := map[string]any{"id": t.ID, "name": t.Name, "subtask": t.Subtask}
	if t.Description != "" {
		out["description"] = t.Description
	}
	return out
}

// JiraComment represents a single Jira comment.
type JiraComment struct {
	ID      string    `json:"id"`
	Body    string    `json:"body,omitempty"`
	Author  *JiraUser `json:"author,omitempty"`
	Created string    `json:"created,omitempty"`
	Updated string    `json:"updated,omitempty"`
}

// FromAPIResponse populates from a comment dict.
func (c *JiraComment) FromAPIResponse(data map[string]any) error {
	c.ID = models.AsID(data["id"])
	c.Body = models.MapString(data, "body")
	c.Created = models.MapString(data, "created")
	c.Updated = models.MapString(data, "updated")
	if a := models.MapMap(data, "author"); a != nil {
		c.Author = &JiraUser{}
		c.Author.FromAPIResponse(a)
	}
	return nil
}

// ToSimplifiedDict returns the JSON-serializable dict.
func (c *JiraComment) ToSimplifiedDict() map[string]any {
	out := map[string]any{"id": c.ID, "body": c.Body, "created": c.Created, "updated": c.Updated}
	if c.Author != nil {
		out["author"] = c.Author.ToSimplifiedDict()
	}
	return out
}

// JiraWorklog represents a single Jira worklog entry.
type JiraWorklog struct {
	ID               string    `json:"id"`
	IssueID          string    `json:"issueId,omitempty"`
	TimeSpent        string    `json:"timeSpent,omitempty"`
	TimeSpentSeconds int64     `json:"timeSpentSeconds"`
	Started          string    `json:"started,omitempty"`
	Created          string    `json:"created,omitempty"`
	Updated          string    `json:"updated,omitempty"`
	Comment          string    `json:"comment,omitempty"`
	Author           *JiraUser `json:"author,omitempty"`
}

// FromAPIResponse populates from a worklog dict.
func (w *JiraWorklog) FromAPIResponse(data map[string]any) error {
	w.ID = models.AsID(data["id"])
	w.IssueID = models.AsID(data["issueId"])
	w.TimeSpent = models.MapString(data, "timeSpent")
	w.TimeSpentSeconds = models.MapInt(data, "timeSpentSeconds")
	w.Started = models.MapString(data, "started")
	w.Created = models.MapString(data, "created")
	w.Updated = models.MapString(data, "updated")
	w.Comment = models.MapString(data, "comment")
	if a := models.MapMap(data, "author"); a != nil {
		w.Author = &JiraUser{}
		w.Author.FromAPIResponse(a)
	}
	return nil
}

// ToSimplifiedDict returns the JSON-serializable dict.
func (w *JiraWorklog) ToSimplifiedDict() map[string]any {
	out := map[string]any{
		"id":               w.ID,
		"timeSpent":        w.TimeSpent,
		"timeSpentSeconds": w.TimeSpentSeconds,
		"started":          w.Started,
		"created":          w.Created,
		"updated":          w.Updated,
	}
	if w.IssueID != "" {
		out["issueId"] = w.IssueID
	}
	if w.Comment != "" {
		out["comment"] = w.Comment
	}
	if w.Author != nil {
		out["author"] = w.Author.ToSimplifiedDict()
	}
	return out
}

// mapStringMap is a small helper to coerce nested map[string]string.
func mapStringMap(m map[string]any, key string) map[string]string {
	v, ok := m[key]
	if !ok {
		return nil
	}
	if raw, ok := v.(map[string]any); ok {
		out := make(map[string]string, len(raw))
		for k, v := range raw {
			if s, ok := v.(string); ok {
				out[k] = s
			}
		}
		return out
	}
	return nil
}
