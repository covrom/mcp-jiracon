// Extended Jira fetcher — watchers, transitions, worklog, links,
// comments, projects, agile, fields, SD. DC-only.
package jira

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	jira "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/jira"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/utils"
)

// --- Watchers -----------------------------------------------------------

func (f *Fetcher) GetIssueWatchers(ctx context.Context, issueKey string) (map[string]any, error) {
	var raw map[string]any
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/api/2/issue/%s/watchers", issueKey), &raw); err != nil {
		return nil, fmt.Errorf("jira: watchers %s: %w", issueKey, err)
	}
	return raw, nil
}

func (f *Fetcher) AddWatcher(ctx context.Context, issueKey, userID string) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s/watchers", issueKey)
	return f.client.Post(ctx, path, map[string]any{"name": userID}, nil)
}

func (f *Fetcher) RemoveWatcher(ctx context.Context, issueKey, username string) error {
	return f.client.Delete(ctx, fmt.Sprintf("/rest/api/2/issue/%s/watchers?username=%s", issueKey, username))
}

// --- Transitions --------------------------------------------------------

func (f *Fetcher) GetTransitions(ctx context.Context, issueKey string) ([]map[string]any, error) {
	var raw struct {
		Transitions []map[string]any `json:"transitions"`
	}
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/api/2/issue/%s/transitions", issueKey), &raw); err != nil {
		return nil, fmt.Errorf("jira: transitions %s: %w", issueKey, err)
	}
	return raw.Transitions, nil
}

func (f *Fetcher) TransitionIssue(ctx context.Context, issueKey, transitionID string, fields map[string]any, comment string) (*jira.JiraIssue, error) {
	body := map[string]any{"transition": map[string]any{"id": transitionID}}
	if len(fields) > 0 {
		body["fields"] = fields
	}
	if comment != "" {
		body["update"] = map[string]any{"comment": []map[string]any{{"add": map[string]any{"body": comment}}}}
	}
	var raw map[string]any
	if err := f.client.Post(ctx, fmt.Sprintf("/rest/api/2/issue/%s/transitions", issueKey), body, &raw); err != nil {
		return nil, fmt.Errorf("jira: transition %s: %w", issueKey, err)
	}
	issue := &jira.JiraIssue{}
	issue.FromAPIResponse(raw)
	return issue, nil
}

// --- Worklog ------------------------------------------------------------

func (f *Fetcher) GetWorklogs(ctx context.Context, issueKey string) ([]*jira.JiraWorklog, error) {
	var allWorklogs []*jira.JiraWorklog
	startAt := 0
	for {
		var raw struct {
			StartAt    int              `json:"startAt"`
			MaxResults int              `json:"maxResults"`
			Total      int              `json:"total"`
			Worklogs   []map[string]any `json:"worklogs"`
		}
		req, err := f.client.NewRequest(ctx, "GET", fmt.Sprintf("/rest/api/2/issue/%s/worklog", issueKey), nil)
		if err != nil {
			return nil, fmt.Errorf("jira: worklog %s: %w", issueKey, err)
		}
		q := req.URL.Query()
		q.Set("startAt", fmt.Sprintf("%d", startAt))
		req.URL.RawQuery = q.Encode()
		if err := f.client.Do(req, &raw); err != nil {
			return nil, fmt.Errorf("jira: worklog %s: %w", issueKey, err)
		}
		for _, w := range raw.Worklogs {
			wl := &jira.JiraWorklog{}
			wl.FromAPIResponse(w)
			allWorklogs = append(allWorklogs, wl)
		}
		startAt += len(raw.Worklogs)
		if startAt >= raw.Total || len(raw.Worklogs) == 0 {
			break
		}
	}
	return allWorklogs, nil
}

func (f *Fetcher) AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) (map[string]any, error) {
	return f.AddWorklogExtended(ctx, issueKey, timeSpent, comment, "", "", "")
}

// AddWorklogExtended adds a worklog entry with all supported optional fields.
func (f *Fetcher) AddWorklogExtended(ctx context.Context, issueKey, timeSpent, comment, started, originalEstimate, remainingEstimate string) (map[string]any, error) {
	body := map[string]any{"timeSpent": timeSpent}
	if comment != "" {
		body["comment"] = comment
	}
	if started != "" {
		body["started"] = started
	}
	if originalEstimate != "" {
		if body["adjustEstimate"] == nil {
			body["adjustEstimate"] = map[string]any{}
		}
		body["adjustEstimate"].(map[string]any)["newEstimate"] = originalEstimate
	}
	if remainingEstimate != "" {
		if body["adjustEstimate"] == nil {
			body["adjustEstimate"] = map[string]any{}
		}
		body["adjustEstimate"].(map[string]any)["newRemainingEstimate"] = remainingEstimate
	}
	var raw map[string]any
	if err := f.client.Post(ctx, fmt.Sprintf("/rest/api/2/issue/%s/worklog", issueKey), body, &raw); err != nil {
		return nil, fmt.Errorf("jira: add worklog %s: %w", issueKey, err)
	}
	return raw, nil
}

// --- Links --------------------------------------------------------------

func (f *Fetcher) GetIssueLinkTypes(ctx context.Context) ([]map[string]any, error) {
	var raw struct {
		IssueLinkTypes []map[string]any `json:"issueLinkTypes"`
	}
	if err := f.client.Get(ctx, "/rest/api/2/issueLinkType", &raw); err != nil {
		return nil, fmt.Errorf("jira: link types: %w", err)
	}
	return raw.IssueLinkTypes, nil
}

func (f *Fetcher) LinkIssueToEpic(ctx context.Context, issueKey, epicKey string) (*jira.JiraIssue, error) {
	body := map[string]any{"fields": map[string]any{"Epic Link": epicKey}}
	var raw map[string]any
	if err := f.client.Put(ctx, fmt.Sprintf("/rest/api/2/issue/%s", issueKey), body, &raw); err != nil {
		return nil, fmt.Errorf("jira: link epic %s: %w", issueKey, err)
	}
	issue := &jira.JiraIssue{}
	issue.FromAPIResponse(raw)
	return issue, nil
}

func (f *Fetcher) CreateIssueLink(ctx context.Context, linkType, inward, outward string) (map[string]any, error) {
	body := map[string]any{"type": map[string]any{"name": linkType}, "inwardIssue": map[string]any{"key": inward}, "outwardIssue": map[string]any{"key": outward}}
	var raw map[string]any
	if err := f.client.Post(ctx, "/rest/api/2/issueLink", body, &raw); err != nil {
		return nil, fmt.Errorf("jira: create link: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) CreateRemoteIssueLink(ctx context.Context, issueKey, url, title, summary, relationship, iconURL string) (map[string]any, error) {
	obj := map[string]any{"url": url, "title": title}
	if summary != "" {
		obj["summary"] = summary
	}
	if relationship != "" {
		obj["relationship"] = relationship
	}
	if iconURL != "" {
		obj["icon"] = map[string]any{"url16x16": iconURL}
	}
	var raw map[string]any
	if err := f.client.Post(ctx, fmt.Sprintf("/rest/api/2/issue/%s/remotelink", issueKey), map[string]any{"object": obj}, &raw); err != nil {
		return nil, fmt.Errorf("jira: remote link %s: %w", issueKey, err)
	}
	return raw, nil
}

func (f *Fetcher) RemoveIssueLink(ctx context.Context, linkID string) error {
	return f.client.Delete(ctx, fmt.Sprintf("/rest/api/2/issueLink/%s", linkID))
}

// --- Comments -----------------------------------------------------------

func (f *Fetcher) AddComment(ctx context.Context, issueKey, body string) (map[string]any, error) {
	var raw map[string]any
	if err := f.client.Post(ctx, fmt.Sprintf("/rest/api/2/issue/%s/comment", issueKey), map[string]any{"body": body}, &raw); err != nil {
		return nil, fmt.Errorf("jira: add comment %s: %w", issueKey, err)
	}
	return raw, nil
}

func (f *Fetcher) EditComment(ctx context.Context, issueKey, commentID, body string) (map[string]any, error) {
	var raw map[string]any
	if err := f.client.Put(ctx, fmt.Sprintf("/rest/api/2/issue/%s/comment/%s", issueKey, commentID), map[string]any{"body": body}, &raw); err != nil {
		return nil, fmt.Errorf("jira: edit comment: %w", err)
	}
	return raw, nil
}

// --- Projects / Versions / Components ----------------------------------

func (f *Fetcher) GetProjectIssues(ctx context.Context, projectKey string, start, limit int) (map[string]any, error) {
	var raw map[string]any
	req, err := f.client.NewRequest(ctx, "GET", "/rest/api/2/search", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("jql", fmt.Sprintf("project = %s", projectKey))
	q.Set("startAt", fmt.Sprintf("%d", start))
	q.Set("maxResults", fmt.Sprintf("%d", limit))
	q.Set("fields", strings.Join(defaultReadFields, ","))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (f *Fetcher) GetAllProjects(ctx context.Context, includeArchived bool) ([]map[string]any, error) {
	var raw struct {
		Values []map[string]any `json:"values"`
	}
	req, err := f.client.NewRequest(ctx, "GET", "/rest/api/2/project/search", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("maxResults", "100")
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: projects: %w", err)
	}
	if !includeArchived {
		filtered := make([]map[string]any, 0, len(raw.Values))
		for _, p := range raw.Values {
			archived, _ := p["archived"].(bool)
			if !archived {
				filtered = append(filtered, p)
			}
		}
		return filtered, nil
	}
	return raw.Values, nil
}

func (f *Fetcher) SearchProjects(ctx context.Context, query string, maxResults int) ([]map[string]any, error) {
	var raw struct {
		Projects []map[string]any `json:"projects"`
	}
	req, err := f.client.NewRequest(ctx, "GET", "/rest/api/2/project/picker", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("query", query)
	q.Set("maxResults", fmt.Sprintf("%d", maxResults))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: search projects: %w", err)
	}
	return raw.Projects, nil
}

func (f *Fetcher) GetProjectIssueTypes(ctx context.Context, projectKey string) ([]map[string]any, error) {
	var raw struct {
		Values []map[string]any `json:"values"`
	}
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/api/2/issue/createmeta/%s/issuetypes", projectKey), &raw); err != nil {
		return nil, fmt.Errorf("jira: issue types: %w", err)
	}
	return raw.Values, nil
}

func (f *Fetcher) GetCreateFields(ctx context.Context, projectKey, issueTypeID string) ([]map[string]any, error) {
	var raw struct {
		Fields []map[string]any `json:"fields"`
	}
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/api/2/issue/createmeta/%s/issuetype/%s", projectKey, issueTypeID), &raw); err != nil {
		return nil, fmt.Errorf("jira: create fields: %w", err)
	}
	return raw.Fields, nil
}

func (f *Fetcher) GetProjectVersions(ctx context.Context, projectKey string) ([]map[string]any, error) {
	var raw []map[string]any
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/api/2/project/%s/versions", projectKey), &raw); err != nil {
		return nil, fmt.Errorf("jira: versions: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) GetProjectComponents(ctx context.Context, projectKey string) ([]map[string]any, error) {
	var raw []map[string]any
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/api/2/project/%s/components", projectKey), &raw); err != nil {
		return nil, fmt.Errorf("jira: components: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) GetProjectFields(ctx context.Context, projectKey string) ([]map[string]any, error) {
	var raw struct {
		Values []map[string]any `json:"values"`
	}
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/api/2/issue/createmeta/%s/issuetypes", projectKey), &raw); err != nil {
		return nil, fmt.Errorf("jira: project fields: %w", err)
	}
	seen := map[string]struct{}{}
	out := []map[string]any{}
	for _, it := range raw.Values {
		fm, _ := it["fields"].(map[string]any)
		if fm == nil {
			continue
		}
		for fk, fv := range fm {
			if _, ok := seen[fk]; ok {
				continue
			}
			seen[fk] = struct{}{}
			fvm, ok := fv.(map[string]any)
			if !ok {
				continue
			}
			out = append(out, map[string]any{"field_id": fk, "name": fvm["name"], "required": fvm["required"], "custom": strings.HasPrefix(fk, "customfield_")})
		}
	}
	return out, nil
}

func (f *Fetcher) CreateVersion(ctx context.Context, projectKey, name, startDate, releaseDate, desc string) (map[string]any, error) {
	body := map[string]any{"name": name, "project": projectKey}
	if startDate != "" {
		body["startDate"] = startDate
	}
	if releaseDate != "" {
		body["releaseDate"] = releaseDate
	}
	if desc != "" {
		body["description"] = desc
	}
	var raw map[string]any
	if err := f.client.Post(ctx, "/rest/api/2/version", body, &raw); err != nil {
		return nil, fmt.Errorf("jira: create version: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) UpdateVersion(ctx context.Context, versionID, name, desc string, archived *bool, released *bool, startDate, releaseDate string) (map[string]any, error) {
	body := map[string]any{}
	if archived != nil {
		body["archived"] = *archived
	}
	if released != nil {
		body["released"] = *released
	}
	if name != "" {
		body["name"] = name
	}
	if desc != "" {
		body["description"] = desc
	}
	if startDate != "" {
		body["startDate"] = startDate
	}
	if releaseDate != "" {
		body["releaseDate"] = releaseDate
	}
	var raw map[string]any
	if err := f.client.Put(ctx, fmt.Sprintf("/rest/api/2/version/%s", versionID), body, &raw); err != nil {
		return nil, fmt.Errorf("jira: update version: %w", err)
	}
	return raw, nil
}

// --- Agile / Boards / Sprints ------------------------------------------

func (f *Fetcher) GetAgileBoards(ctx context.Context, boardName, projectKey, boardType string, start, limit int) ([]map[string]any, error) {
	var raw struct {
		Values []map[string]any `json:"values"`
	}
	req, err := f.client.NewRequest(ctx, "GET", "/rest/agile/1.0/board", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if boardName != "" {
		q.Set("name", boardName)
	}
	if projectKey != "" {
		q.Set("projectKeyOrId", projectKey)
	}
	if boardType != "" {
		q.Set("type", boardType)
	}
	q.Set("startAt", fmt.Sprintf("%d", start))
	q.Set("maxResults", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: boards: %w", err)
	}
	return raw.Values, nil
}

func (f *Fetcher) GetBoardIssues(ctx context.Context, boardID, jql string, start, limit int) (map[string]any, error) {
	var raw map[string]any
	req, err := f.client.NewRequest(ctx, "GET", fmt.Sprintf("/rest/agile/1.0/board/%s/issue", boardID), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("jql", jql)
	q.Set("startAt", fmt.Sprintf("%d", start))
	q.Set("maxResults", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: board issues: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) GetSprintsFromBoard(ctx context.Context, boardID, state string, start, limit int) ([]map[string]any, error) {
	var raw struct {
		Values []map[string]any `json:"values"`
	}
	req, err := f.client.NewRequest(ctx, "GET", fmt.Sprintf("/rest/agile/1.0/board/%s/sprint", boardID), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if state != "" {
		q.Set("state", state)
	}
	q.Set("startAt", fmt.Sprintf("%d", start))
	q.Set("maxResults", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: sprints: %w", err)
	}
	return raw.Values, nil
}

func (f *Fetcher) GetSprintIssues(ctx context.Context, sprintID string, start, limit int) (map[string]any, error) {
	var raw map[string]any
	req, err := f.client.NewRequest(ctx, "GET", fmt.Sprintf("/rest/agile/1.0/sprint/%s/issue", sprintID), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("startAt", fmt.Sprintf("%d", start))
	q.Set("maxResults", fmt.Sprintf("%d", limit))
	q.Set("fields", strings.Join(defaultReadFields, ","))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: sprint issues: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) CreateSprint(ctx context.Context, boardID, name, startDate, endDate, goal string) (map[string]any, error) {
	body := map[string]any{"name": name, "originBoardId": boardID}
	if startDate != "" {
		body["startDate"] = startDate
	}
	if endDate != "" {
		body["endDate"] = endDate
	}
	if goal != "" {
		body["goal"] = goal
	}
	var raw map[string]any
	if err := f.client.Post(ctx, "/rest/agile/1.0/sprint", body, &raw); err != nil {
		return nil, fmt.Errorf("jira: create sprint: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) UpdateSprint(ctx context.Context, sprintID, name, state, startDate, endDate, goal string) (map[string]any, error) {
	body := map[string]any{}
	if name != "" {
		body["name"] = name
	}
	if state != "" {
		body["state"] = state
	}
	if startDate != "" {
		body["startDate"] = startDate
	}
	if endDate != "" {
		body["endDate"] = endDate
	}
	if goal != "" {
		body["goal"] = goal
	}
	var raw map[string]any
	if err := f.client.Put(ctx, fmt.Sprintf("/rest/agile/1.0/sprint/%s", sprintID), body, &raw); err != nil {
		return nil, fmt.Errorf("jira: update sprint: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) AddIssuesToSprint(ctx context.Context, sprintID string, issueKeys []string) error {
	return f.client.Post(ctx, fmt.Sprintf("/rest/agile/1.0/sprint/%s/issue", sprintID), map[string]any{"issues": issueKeys}, nil)
}

func (f *Fetcher) MoveIssuesToBacklog(ctx context.Context, issueKeys []string) error {
	return f.client.Post(ctx, "/rest/agile/1.0/backlog/issue", map[string]any{"issues": issueKeys}, nil)
}

// --- Fields -------------------------------------------------------------

func (f *Fetcher) SearchFields(ctx context.Context, keyword string, limit int, refresh bool) ([]map[string]any, error) {
	_ = refresh
	var raw []map[string]any
	if err := f.client.Get(ctx, "/rest/api/2/field", &raw); err != nil {
		return nil, fmt.Errorf("jira: fields: %w", err)
	}
	if keyword != "" {
		// Score fields by fuzzy match and sort descending, then trim to limit.
		type scoredField struct {
			f    map[string]any
			scr  float64
		}
		list := make([]scoredField, len(raw))
		for i, field := range raw {
			name, _ := field["name"].(string)
			list[i] = scoredField{f: field, scr: fieldScore(name, keyword)}
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].scr > list[j].scr
		})
		raw = make([]map[string]any, 0, len(list))
		for _, s := range list {
			if s.scr > 0 {
				raw = append(raw, s.f)
			}
		}
	}
	if limit > 0 && limit < len(raw) {
		raw = raw[:limit]
	}
	return raw, nil
}

// fieldScore computes a fuzzy-match score in [0,1] between a field name and
// a search keyword. It favours exact matches, prefix matches, and token
// overlap, backed by a simple Levenshtein distance.
func fieldScore(name, keyword string) float64 {
	lower := strings.ToLower(name)
	kw := strings.ToLower(keyword)

	if lower == kw {
		return 1.0
	}
	if strings.HasPrefix(lower, kw) {
		return 0.9
	}
	if strings.Contains(lower, kw) {
		// Substring match — score decays with relative length difference.
		return 0.8 - (float64(len(lower)-len(kw)) / float64(len(lower)) * 0.4)
	}

	// Token-based scoring.
	nameTokens := strings.FieldsFunc(lower, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-' || r == '.' || r == '/'
	})
	kwTokens := strings.FieldsFunc(kw, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-' || r == '.' || r == '/'
	})
	if len(nameTokens) == 0 || len(kwTokens) == 0 {
		return 0
	}

	bestTotal := 0.0
	for _, nt := range nameTokens {
		best := 0.0
		for _, kt := range kwTokens {
			if nt == kt {
				best = 1.0
				break
			}
			if strings.HasPrefix(nt, kt) {
				best = maxFloat(best, 0.85)
				continue
			}
			if strings.Contains(nt, kt) {
				best = maxFloat(best, 0.6)
				continue
			}
			// Levenshtein similarity for short tokens.
			d := levenshteinDistance(nt, kt)
			maxLen := maxInt2(len(nt), len(kt))
			if maxLen > 0 {
				sim := 1.0 - float64(d)/float64(maxLen)
				if sim > 0.3 {
					best = maxFloat(best, sim*0.5)
				}
			}
		}
		bestTotal += best
	}
	if len(nameTokens) == 0 {
		return 0
	}
	return bestTotal / float64(len(nameTokens))
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func maxInt2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func levenshteinDistance(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	// Use one-row DP for memory efficiency.
	prev := make([]int, lb+1)
	cur := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		cur[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = minInt3(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev, cur = cur, prev
	}
	return prev[lb]
}

func minInt3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// GetFieldOptions returns the allowed option values for a custom field
// (select, multi-select, radio, checkbox, cascading select).
// Uses createmeta — the approach recommended for DC (Server) installations.
// Requires projectKey and issueType to resolve the correct field configuration.
func (f *Fetcher) GetFieldOptions(ctx context.Context, fieldID, projectKey, issueType string) ([]map[string]any, error) {
	// Step 1: get issue types for the project to resolve issueType name → ID.
	issueTypes, err := f.GetProjectIssueTypes(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("jira: field options: get issue types: %w", err)
	}
	var issueTypeID string
	for _, it := range issueTypes {
		name, _ := it["name"].(string)
		if strings.EqualFold(name, issueType) {
			issueTypeID, _ = it["id"].(string)
			break
		}
	}
	if issueTypeID == "" {
		return nil, fmt.Errorf("jira: field options: issue type %q not found in project %q", issueType, projectKey)
	}

	// Step 2: paginate through createmeta fields to find the matching fieldID.
	const pageSize = 50
	startAt := 0
	for {
		req, err := f.client.NewRequest(ctx, "GET",
			fmt.Sprintf("/rest/api/2/issue/createmeta/%s/issuetypes/%s", projectKey, issueTypeID), nil)
		if err != nil {
			return nil, fmt.Errorf("jira: field options: %w", err)
		}
		q := req.URL.Query()
		q.Set("startAt", fmt.Sprintf("%d", startAt))
		q.Set("maxResults", fmt.Sprintf("%d", pageSize))
		req.URL.RawQuery = q.Encode()

		var raw map[string]any
		if err := f.client.Do(req, &raw); err != nil {
			return nil, fmt.Errorf("jira: field options: %w", err)
		}

		values, _ := raw["values"].([]any)
		total, _ := raw["total"].(float64)

		for _, v := range values {
			entry, ok := v.(map[string]any)
			if !ok {
				continue
			}
			fid, _ := entry["fieldId"].(string)
			if fid != fieldID {
				continue
			}
			allowedValues, _ := entry["allowedValues"].([]any)
			return convertFieldOptions(allowedValues), nil
		}

		startAt += len(values)
		if float64(startAt) >= total || len(values) == 0 {
			break
		}
	}
	return nil, nil
}

// convertFieldOptions transforms raw allowedValues from createmeta into
// a simplified slice with id, value, disabled, and child_options fields.
func convertFieldOptions(raw []any) []map[string]any {
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		opt := map[string]any{
			"id":    fmt.Sprintf("%v", m["id"]),
			"value": m["value"],
		}
		if disabled, _ := m["disabled"].(bool); disabled {
			opt["disabled"] = true
		}
		if children, _ := m["cascadingOptions"].([]any); len(children) > 0 {
			childOpts := make([]map[string]any, 0, len(children))
			for _, c := range children {
				cm, ok := c.(map[string]any)
				if !ok {
					continue
				}
				childOpts = append(childOpts, map[string]any{
					"id":    fmt.Sprintf("%v", cm["id"]),
					"value": cm["value"],
				})
			}
			if len(childOpts) > 0 {
				opt["child_options"] = childOpts
			}
		}
		out = append(out, opt)
	}
	return out
}

// --- Search users -------------------------------------------------------

func (f *Fetcher) SearchAssignableUsers(ctx context.Context, query, projectKey, issueKey string, limit int) ([]*jira.JiraUser, error) {
	req, err := f.client.NewRequest(ctx, "GET", "/rest/api/2/user/assignable/search", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("query", query)
	if projectKey != "" {
		q.Set("project", projectKey)
	}
	if issueKey != "" {
		q.Set("issueKey", issueKey)
	}
	q.Set("maxResults", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	var raw []map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: search users: %w", err)
	}
	out := make([]*jira.JiraUser, len(raw))
	for i, u := range raw {
		out[i] = &jira.JiraUser{}
		out[i].FromAPIResponse(u)
	}
	return out, nil
}

func (f *Fetcher) GetIssueImages(ctx context.Context, issueKey string) (map[string]any, error) {
	attachments, err := f.GetIssueAttachmentsRaw(ctx, issueKey)
	if err != nil {
		return nil, err
	}
	if len(attachments) == 0 {
		return map[string]any{
			"success":   true,
			"issue_key": issueKey,
			"total":     0,
			"images":    0,
			"fetched":   []map[string]any{},
			"failed":    []map[string]any{},
		}, nil
	}

	fetched := []map[string]any{}
	failed := []map[string]any{}

	for _, a := range attachments {
		mime, _ := a["mimeType"].(string)
		fn, _ := a["filename"].(string)
		size, _ := a["size"].(float64)

		// Check size limit.
		if int64(size) > utils.AttachmentMaxBytes {
			failed = append(failed, map[string]any{
				"filename": fn,
				"error":    fmt.Sprintf("Image is %.0f bytes which exceeds the 50 MB inline limit.", size),
			})
			continue
		}

		// Filter to image attachments.
		isImg, resolvedMime := utils.IsImageAttachment(mime, fn)
		if !isImg {
			continue
		}

		contentURL, _ := a["content"].(string)
		if contentURL == "" {
			failed = append(failed, map[string]any{
				"filename": fn, "error": "No download URL",
			})
			continue
		}

		data, err := f.client.GetRaw(ctx, contentURL)
		if err != nil {
			failed = append(failed, map[string]any{
				"filename": fn, "error": fmt.Sprintf("Fetch failed: %v", err),
			})
			continue
		}

		encoded := base64.StdEncoding.EncodeToString(data)
		attID := fmt.Sprintf("%v", a["id"])
		fetched = append(fetched, map[string]any{
			"filename":      fn,
			"mime_type":     resolvedMime,
			"data_base64":   encoded,
			"size":          len(data),
			"attachment_id": attID,
		})
	}

	return map[string]any{
		"success":   true,
		"issue_key": issueKey,
		"total":     len(attachments),
		"images":    len(fetched),
		"fetched":   fetched,
		"failed":    failed,
	}, nil
}

func (f *Fetcher) GetIssueAttachmentsRaw(ctx context.Context, issueKey string) ([]map[string]any, error) {
	var raw map[string]any
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/api/2/issue/%s?fields=attachment", issueKey), &raw); err != nil {
		return nil, fmt.Errorf("jira: attachments: %w", err)
	}
	fields, _ := raw["fields"].(map[string]any)
	if fields == nil {
		return nil, nil
	}
	attList, _ := fields["attachment"].([]any)
	out := make([]map[string]any, len(attList))
	for i, a := range attList {
		out[i] = a.(map[string]any)
	}
	return out, nil
}

// GetIssueAttachmentContents downloads all attachments for a Jira issue
// into memory. Returns a map with success, issue_key, total, attachments
// (each with filename, content_type, size, data), and failed list.
// Mirrors Python's get_issue_attachment_contents.
func (f *Fetcher) GetIssueAttachmentContents(ctx context.Context, issueKey string) (map[string]any, error) {
	const maxBytes = 50 * 1024 * 1024 // ATTACHMENT_MAX_BYTES

	attachments, err := f.GetIssueAttachmentsRaw(ctx, issueKey)
	if err != nil {
		return nil, err
	}
	if len(attachments) == 0 {
		return map[string]any{
			"success":     true,
			"message":     fmt.Sprintf("No attachments found for issue %s", issueKey),
			"attachments": []any{},
			"failed":      []any{},
		}, nil
	}

	fetched := []map[string]any{}
	failed := []map[string]any{}

	for _, a := range attachments {
		filename, _ := a["filename"].(string)
		if filename == "" {
			filename = fmt.Sprintf("%v", a["id"])
		}

		contentURL, _ := a["content"].(string)
		if contentURL == "" {
			failed = append(failed, map[string]any{
				"filename": filename, "error": "No URL available",
			})
			continue
		}

		// Check size limit.
		if size, ok := a["size"].(float64); ok && int64(size) > maxBytes {
			failed = append(failed, map[string]any{
				"filename": filename,
				"error":    fmt.Sprintf("Attachment '%s' is %.0f bytes which exceeds the 50 MB inline limit.", filename, size),
			})
			continue
		}

		data, err := f.client.GetRaw(ctx, contentURL)
		if err != nil {
			failed = append(failed, map[string]any{
				"filename": filename, "error": fmt.Sprintf("Fetch failed: %v", err),
			})
			continue
		}

		contentType, _ := a["mimeType"].(string)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		fetched = append(fetched, map[string]any{
			"filename":     filename,
			"content_type": contentType,
			"size":         len(data),
			"data":         data,
		})
	}

	return map[string]any{
		"success":     true,
		"issue_key":   issueKey,
		"total":       len(attachments),
		"attachments": fetched,
		"failed":      failed,
	}, nil
}

// DownloadAttachments downloads all attachments for an issue and saves
// them to targetDir. Returns a summary list of saved files.
func (f *Fetcher) DownloadAttachments(ctx context.Context, issueKey, targetDir string) ([]map[string]any, error) {
	attachments, err := f.GetIssueAttachmentsRaw(ctx, issueKey)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return nil, fmt.Errorf("jira: create target dir: %w", err)
	}
	results := make([]map[string]any, 0, len(attachments))
	for _, a := range attachments {
		id := fmt.Sprintf("%v", a["id"])
		filename, _ := a["filename"].(string)
		if filename == "" {
			filename = id
		}
		contentURL, _ := a["content"].(string)
		if contentURL == "" {
			results = append(results, map[string]any{
				"id": id, "filename": filename, "error": "No URL available",
			})
			continue
		}
		data, err := f.client.GetRaw(ctx, contentURL)
		if err != nil {
			results = append(results, map[string]any{
				"id": id, "filename": filename, "error": err.Error(),
			})
			continue
		}
		destPath := filepath.Join(targetDir, filename)
		if err := os.WriteFile(destPath, data, 0o644); err != nil {
			results = append(results, map[string]any{
				"id": id, "filename": filename, "error": err.Error(),
			})
			continue
		}
		results = append(results, map[string]any{
			"id":       id,
			"filename": filename,
			"path":     destPath,
			"size":     len(data),
		})
	}
	return results, nil
}

// --- Service Desk (DC-only) ---------------------------------------------

func (f *Fetcher) GetServiceDesk(ctx context.Context, projectKey string) (map[string]any, error) {
	var raw map[string]any
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/servicedeskapi/servicedesk/by-project/%s", projectKey), &raw); err != nil {
		return nil, fmt.Errorf("jira: service desk: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) GetServiceDeskQueues(ctx context.Context, sdID string, start, limit int) (map[string]any, error) {
	var raw map[string]any
	req, err := f.client.NewRequest(ctx, "GET", fmt.Sprintf("/rest/servicedeskapi/servicedesk/%s/queue", sdID), nil)
	if err != nil {
		return nil, fmt.Errorf("jira: sd queues: %w", err)
	}
	q := req.URL.Query()
	q.Set("start", fmt.Sprintf("%d", start))
	q.Set("limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: sd queues: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) GetQueueIssues(ctx context.Context, sdID, queueID string, start, limit int) (map[string]any, error) {
	var raw map[string]any
	req, err := f.client.NewRequest(ctx, "GET", fmt.Sprintf("/rest/servicedeskapi/servicedesk/%s/queue/%s/issue", sdID, queueID), nil)
	if err != nil {
		return nil, fmt.Errorf("jira: queue issues: %w", err)
	}
	q := req.URL.Query()
	q.Set("start", fmt.Sprintf("%d", start))
	q.Set("limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: queue issues: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) GetRequestTypes(ctx context.Context, sdID string, start, limit int) (map[string]any, error) {
	var raw map[string]any
	req, err := f.client.NewRequest(ctx, "GET", fmt.Sprintf("/rest/servicedeskapi/servicedesk/%s/requesttype", sdID), nil)
	if err != nil {
		return nil, fmt.Errorf("jira: request types: %w", err)
	}
	q := req.URL.Query()
	q.Set("start", fmt.Sprintf("%d", start))
	q.Set("limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: request types: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) CreateCustomerRequest(ctx context.Context, sdID, rtID string, fields map[string]any, raiseOnBehalfOf, requestParticipants string, strictOnBehalf bool, attachments []map[string]any) (map[string]any, error) {
	body := map[string]any{"serviceDeskId": sdID, "requestTypeId": rtID, "requestFieldValues": fields}
	if raiseOnBehalfOf != "" {
		body["raiseOnBehalfOf"] = raiseOnBehalfOf
	}
	if strictOnBehalf {
		body["strictOnBehalfOf"] = true
	}
	if requestParticipants != "" {
		parts := []string{}
		for _, p := range strings.Split(requestParticipants, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				parts = append(parts, p)
			}
		}
		body["requestParticipants"] = parts
	}
	if len(attachments) > 0 {
		body["attachments"] = attachments
	}
	var raw map[string]any
	if err := f.client.Post(ctx, "/rest/servicedeskapi/request", body, &raw); err != nil {
		// Fallback: retry without raiseOnBehalfOf when not in strict mode.
		if raiseOnBehalfOf != "" && !strictOnBehalf {
			delete(body, "raiseOnBehalfOf")
			if err2 := f.client.Post(ctx, "/rest/servicedeskapi/request", body, &raw); err2 != nil {
				return nil, fmt.Errorf("jira: create customer request: %w", err2)
			}
			return raw, nil
		}
		return nil, fmt.Errorf("jira: create customer request: %w", err)
	}
	return raw, nil
}

func (f *Fetcher) GetRequestTypeFields(ctx context.Context, sdID, rtID string) (map[string]any, error) {
	var raw map[string]any
	if err := f.client.Get(ctx, fmt.Sprintf("/rest/servicedeskapi/servicedesk/%s/requesttype/%s/field", sdID, rtID), &raw); err != nil {
		return nil, fmt.Errorf("jira: rt fields: %w", err)
	}
	return raw, nil
}

// --- Development Info (dev-status plugin, DC + Cloud) ------------------------

var commonAppTypes = []string{"stash", "bitbucket", "GitHub", "GitLab"}
var appTypeCasing = map[string]string{"github": "GitHub", "gitlab": "GitLab"}

func (f *Fetcher) fetchDevInfo(ctx context.Context, issueKey, issueID, appType, dataType string) (map[string]any, error) {
	req, err := f.client.NewRequest(ctx, "GET", "/rest/dev-status/1.0/issue/detail", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("issueId", issueID)
	q.Set("applicationType", appType)
	if dataType != "" {
		q.Set("dataType", dataType)
	}
	req.URL.RawQuery = q.Encode()
	var raw map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		// 404 / 403 are expected when plugin is not installed or access is denied.
		if he, ok := err.(interface{ StatusCode() int }); ok && (he.StatusCode() == 404 || he.StatusCode() == 403) {
			return map[string]any{
				"issue_key":    issueKey,
				"error":        fmt.Sprintf("Dev-status plugin returned %d — may not be installed or accessible.", he.StatusCode()),
				"detail":       []any{},
				"pullRequests": []any{},
				"branches":     []any{},
				"commits":      []any{},
				"repositories": []any{},
			}, nil
		}
		return nil, fmt.Errorf("jira: dev info %s/%s: %w", issueKey, appType, err)
	}
	return parseDevInfo(raw, issueKey), nil
}

func parseDevInfo(raw map[string]any, issueKey string) map[string]any {
	result := map[string]any{
		"issue_key":    issueKey,
		"detail":       []any{},
		"pullRequests": []any{},
		"branches":     []any{},
		"commits":      []any{},
		"repositories": []any{},
	}
	details, _ := raw["detail"].([]any)
	for _, detail := range details {
		d, ok := detail.(map[string]any)
		if !ok {
			continue
		}
		result["detail"] = append(result["detail"].([]any), d)

		instance, _ := d["_instance"].(map[string]any)
		instanceName := ""
		if instance != nil {
			instanceName, _ = instance["name"].(string)
		}

		for _, pr := range extractList(d, "pullRequests") {
			source, _ := pr["source"].(map[string]any)
			dest, _ := pr["destination"].(map[string]any)
			srcRepo, _ := source["repository"].(map[string]any)
			author, _ := pr["author"].(map[string]any)
			reviewers := []string{}
			for _, r := range extractList(pr, "reviewers") {
				if name, _ := r["name"].(string); name != "" {
					reviewers = append(reviewers, name)
				}
			}
			prEntry := map[string]any{
				"id": pr["id"], "name": pr["name"], "status": pr["status"],
				"url": pr["url"], "lastUpdate": pr["lastUpdate"],
				"instance": instanceName,
			}
			if source != nil {
				prEntry["source"] = source["branch"]
			}
			if dest != nil {
				prEntry["destination"] = dest["branch"]
			}
			if author != nil {
				prEntry["author"] = author["name"]
			}
			prEntry["reviewers"] = reviewers
			if srcRepo != nil {
				prEntry["repository"] = srcRepo["name"]
				prEntry["repositoryUrl"] = srcRepo["url"]
			}
			result["pullRequests"] = append(result["pullRequests"].([]any), prEntry)
		}

		for _, branch := range extractList(d, "branches") {
			result["branches"] = append(result["branches"].([]any), map[string]any{
				"name":                 branch["name"],
				"url":                  branch["url"],
				"createPullRequestUrl": branch["createPullRequestUrl"],
				"instance":             instanceName,
			})
		}

		repos := extractList(d, "repositories")
		for _, repo := range repos {
			repoName, _ := repo["name"].(string)
			repoURL, _ := repo["url"].(string)
			for _, commit := range extractList(repo, "commits") {
				auth, _ := commit["author"].(map[string]any)
				commitEntry := map[string]any{
					"id":              commit["id"],
					"displayId":       commit["displayId"],
					"message":         commit["message"],
					"authorTimestamp": commit["authorTimestamp"],
					"url":             commit["url"],
					"repository":      repoName,
					"repositoryUrl":   repoURL,
				}
				if auth != nil {
					commitEntry["author"] = auth["name"]
				}
				result["commits"] = append(result["commits"].([]any), commitEntry)
			}
			// PRs from repos (fallback)
			for _, pr := range extractList(repo, "pullRequests") {
				pr["repository"] = repoName
				pr["repositoryUrl"] = repoURL
				result["pullRequests"] = append(result["pullRequests"].([]any), pr)
			}
			// Branches from repos (fallback)
			for _, branch := range extractList(repo, "branches") {
				branch["repository"] = repoName
				branch["repositoryUrl"] = repoURL
				result["branches"] = append(result["branches"].([]any), branch)
			}
			if repoName != "" && repoName != "Unknown" {
				repoInfo := map[string]any{"name": repoName, "url": repoURL}
				if av, _ := repo["avatar"].(string); av != "" {
					repoInfo["avatar"] = av
				}
				repos := result["repositories"].([]any)
				found := false
				for _, existing := range repos {
					if em, ok := existing.(map[string]any); ok && em["name"] == repoName {
						found = true
						break
					}
				}
				if !found {
					result["repositories"] = append(repos, repoInfo)
				}
			}
		}
	}
	return result
}

func extractList(m map[string]any, key string) []map[string]any {
	raw, _ := m[key].([]any)
	if raw == nil {
		return nil
	}
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		im, ok := item.(map[string]any)
		if ok {
			out = append(out, im)
		}
	}
	return out
}

func (f *Fetcher) discoverAppTypes(ctx context.Context, issueKey, issueID, dataType string) []string {
	// Try dev-status summary endpoint.
	req, err := f.client.NewRequest(ctx, "GET", "/rest/dev-status/1.0/issue/summary", nil)
	if err != nil {
		return commonAppTypes
	}
	q := req.URL.Query()
	q.Set("issueId", issueID)
	req.URL.RawQuery = q.Encode()
	var raw map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		return commonAppTypes
	}
	summary, _ := raw["summary"].(map[string]any)
	if summary == nil {
		return commonAppTypes
	}
	dataTypes := []string{dataType}
	if dataType == "" {
		dataTypes = []string{"pullrequest", "branch", "repository"}
	}
	appTypes := map[string]struct{}{}
	for _, dt := range dataTypes {
		section, _ := summary[dt].(map[string]any)
		if section == nil {
			continue
		}
		byType, _ := section["byInstanceType"].(map[string]any)
		for appType, instSummary := range byType {
			is, ok := instSummary.(map[string]any)
			if !ok {
				continue
			}
			count, _ := is["count"].(float64)
			if count > 0 {
				if fixed, ok2 := appTypeCasing[strings.ToLower(appType)]; ok2 {
					appType = fixed
				}
				appTypes[appType] = struct{}{}
			}
		}
	}
	result := make([]string, 0, len(appTypes))
	for at := range appTypes {
		result = append(result, at)
	}
	sort.Strings(result)
	if len(result) == 0 {
		return commonAppTypes
	}
	return result
}

// GetIssueDevelopmentInfo returns development information (PRs, branches,
// commits) linked to a Jira issue via the dev-status plugin API.
func (f *Fetcher) GetIssueDevelopmentInfo(ctx context.Context, issueKey, applicationType, dataType string) (map[string]any, error) {
	issue, err := f.GetIssue(ctx, issueKey, []string{})
	if err != nil {
		return nil, fmt.Errorf("jira: dev info: get issue %s: %w", issueKey, err)
	}
	issueID := fmt.Sprintf("%v", issue.ID)
	if issueID == "" || issueID == "0" {
		return nil, fmt.Errorf("jira: dev info: could not get numeric issue ID for %s", issueKey)
	}

	if applicationType != "" {
		return f.fetchDevInfo(ctx, issueKey, issueID, applicationType, dataType)
	}

	appTypes := f.discoverAppTypes(ctx, issueKey, issueID, dataType)
	dataTypes := []string{dataType}
	if dataType == "" {
		dataTypes = []string{"pullrequest", "branch", "repository"}
	}

	merged := map[string]any{
		"issue_key":    issueKey,
		"detail":       []any{},
		"pullRequests": []any{},
		"branches":     []any{},
		"commits":      []any{},
		"repositories": []any{},
	}

	for _, appType := range appTypes {
		for _, dt := range dataTypes {
			result, fetchErr := f.fetchDevInfo(ctx, issueKey, issueID, appType, dt)
			if fetchErr != nil {
				return nil, fetchErr
			}
			if errMsg, _ := result["error"].(string); errMsg != "" {
				if _, hasErr := merged["error"]; !hasErr {
					merged["error"] = errMsg
				}
				break
			}
			for _, key := range []string{"detail", "pullRequests", "branches", "commits"} {
				if items, ok := result[key].([]any); ok {
					merged[key] = append(merged[key].([]any), items...)
				}
			}
			for _, repo := range result["repositories"].([]any) {
				rm, _ := repo.(map[string]any)
				if rm == nil {
					continue
				}
				found := false
				for _, existing := range merged["repositories"].([]any) {
					em, _ := existing.(map[string]any)
					if em != nil && em["name"] == rm["name"] {
						found = true
						break
					}
				}
				if !found {
					merged["repositories"] = append(merged["repositories"].([]any), rm)
				}
			}
		}
		if _, hasErr := merged["error"]; hasErr {
			break
		}
	}
	return merged, nil
}

// GetIssuesDevelopmentInfo returns dev-info for multiple Jira issues.
func (f *Fetcher) GetIssuesDevelopmentInfo(ctx context.Context, issueKeys []string, applicationType, dataType string) ([]map[string]any, error) {
	results := make([]map[string]any, 0, len(issueKeys))
	for _, key := range issueKeys {
		info, err := f.GetIssueDevelopmentInfo(ctx, key, applicationType, dataType)
		if err != nil {
			results = append(results, map[string]any{
				"issue_key":    key,
				"error":        err.Error(),
				"pullRequests": []any{},
				"branches":     []any{},
				"commits":      []any{},
			})
			continue
		}
		results = append(results, info)
	}
	return results, nil
}

// --- Project Analysis --------------------------------------------------------

var childOfPhrases = map[string]struct{}{
	"is child of": {}, "is contained by": {}, "split from": {},
}

func projectKeyFromIssueKey(issueKey string) string {
	if idx := strings.LastIndex(issueKey, "-"); idx >= 0 {
		return issueKey[:idx]
	}
	return issueKey
}

func detectParentKey(epic map[string]any, ownProject string) string {
	if parent, ok := epic["parent"].(map[string]any); ok {
		if pk, _ := parent["key"].(string); pk != "" && projectKeyFromIssueKey(pk) != ownProject {
			return pk
		}
	}
	for _, link := range extractList(epic, "issuelinks") {
		lt, _ := link["type"].(map[string]any)
		if lt == nil {
			continue
		}
		linkName := strings.ToLower(fmt.Sprintf("%v", lt["name"]))
		inwardLabel := strings.ToLower(fmt.Sprintf("%v", lt["inward"]))
		outwardLabel := strings.ToLower(fmt.Sprintf("%v", lt["outward"]))

		for _, dir := range []struct{ key, label string }{
			{"inward_issue", inwardLabel}, {"outward_issue", outwardLabel},
		} {
			target, ok := link[dir.key].(map[string]any)
			if !ok {
				continue
			}
			targetKey, _ := target["key"].(string)
			if targetKey == "" || projectKeyFromIssueKey(targetKey) == ownProject {
				continue
			}
			if _, ok2 := childOfPhrases[dir.label]; ok2 {
				return targetKey
			}
			if _, ok2 := childOfPhrases[linkName]; ok2 {
				return targetKey
			}
		}
	}
	return ""
}

// GetProjectEpicHierarchy groups a project's epics under their
// cross-project parent issues using the parent field and issuelinks.
func (f *Fetcher) GetProjectEpicHierarchy(ctx context.Context, projectKey string, maxEpics int) (map[string]any, error) {
	jql := fmt.Sprintf(`project = "%s" AND issuetype = Epic ORDER BY updated DESC`, projectKey)
	epics := f.fetchIssuesWithLinks(ctx, jql, maxEpics)

	parentKeys := map[string]struct{}{}
	epicToParent := map[string]string{}
	for _, epic := range epics {
		key, _ := epic["key"].(string)
		parentKey := detectParentKey(epic, projectKey)
		epicToParent[key] = parentKey
		if parentKey != "" {
			parentKeys[parentKey] = struct{}{}
		}
	}

	// Batch-fetch parent summaries.
	parentInfo := map[string]map[string]string{}
	if len(parentKeys) > 0 {
		keysList := make([]string, 0, len(parentKeys))
		for k := range parentKeys {
			keysList = append(keysList, k)
		}
		sort.Strings(keysList)
		// Fetch in chunks of 50.
		for i := 0; i < len(keysList); i += 50 {
			end := i + 50
			if end > len(keysList) {
				end = len(keysList)
			}
			chunk := keysList[i:end]
			chunkJQL := "key in (" + strings.Join(chunk, ",") + ")"
			search, sErr := f.searchJQL(ctx, chunkJQL, []string{"summary", "status"}, len(chunk))
			if sErr != nil {
				continue
			}
			issues, _ := search["issues"].([]any)
			for _, iss := range issues {
				im, _ := iss.(map[string]any)
				if im == nil {
					continue
				}
				k, _ := im["key"].(string)
				s, _ := im["summary"].(string)
				st := "Unknown"
				if stObj, ok := im["status"].(map[string]any); ok {
					st, _ = stObj["name"].(string)
				}
				if k != "" {
					parentInfo[k] = map[string]string{"summary": s, "status": st}
				}
			}
		}
	}

	groups := map[string][]map[string]any{}
	unlinked := []map[string]any{}
	for _, epic := range epics {
		key, _ := epic["key"].(string)
		summary, _ := epic["summary"].(string)
		status := ""
		if st, ok := epic["status"].(map[string]any); ok {
			status, _ = st["name"].(string)
		}
		entry := map[string]any{"key": key, "summary": summary, "status": status}
		parentKey := epicToParent[key]
		if parentKey == "" {
			unlinked = append(unlinked, entry)
		} else {
			groups[parentKey] = append(groups[parentKey], entry)
		}
	}

	resultGroups := []map[string]any{}
	if len(unlinked) > 0 {
		resultGroups = append(resultGroups, map[string]any{
			"parent": nil, "group_name": "Unlinked", "epics": unlinked,
		})
	}
	for pk, epicsList := range groups {
		info := parentInfo[pk]
		parentSummary, parentStatus := "", ""
		if info != nil {
			parentSummary = info["summary"]
			parentStatus = info["status"]
		}
		resultGroups = append(resultGroups, map[string]any{
			"parent": map[string]any{
				"key": pk, "summary": parentSummary,
				"project": projectKeyFromIssueKey(pk), "status": parentStatus,
			},
			"epics": epicsList,
		})
	}

	return map[string]any{
		"project_key": projectKey,
		"total_epics": len(epics),
		"groups":      resultGroups,
	}, nil
}

// GetCrossProjectDependencies finds all cross-project issue links for a project.
func (f *Fetcher) GetCrossProjectDependencies(ctx context.Context, projectKey string, maxIssues int) (map[string]any, error) {
	jql := fmt.Sprintf(`project = "%s" ORDER BY updated DESC`, projectKey)
	issues := f.fetchIssuesWithLinks(ctx, jql, maxIssues)

	byProject := map[string]map[string][]map[string]string{}
	totalLinks := 0

	for _, issue := range issues {
		issueKey, _ := issue["key"].(string)
		for _, link := range extractList(issue, "issuelinks") {
			lt, _ := link["type"].(map[string]any)
			linkTypeName := ""
			if lt != nil {
				linkTypeName, _ = lt["name"].(string)
			}
			for _, dir := range []struct{ dir, field string }{
				{"outward", "outward_issue"}, {"inward", "inward_issue"},
			} {
				target, ok := link[dir.field].(map[string]any)
				if !ok {
					continue
				}
				targetKey, _ := target["key"].(string)
				if targetKey == "" || projectKeyFromIssueKey(targetKey) == projectKey {
					continue
				}
				targetProj := projectKeyFromIssueKey(targetKey)
				if byProject[targetProj] == nil {
					byProject[targetProj] = map[string][]map[string]string{}
				}
				byProject[targetProj][linkTypeName] = append(byProject[targetProj][linkTypeName], map[string]string{
					"source":    issueKey,
					"target":    targetKey,
					"direction": dir.dir,
				})
				totalLinks++
			}
		}
	}

	byProjectOut := map[string]any{}
	projKeys := make([]string, 0, len(byProject))
	for k := range byProject {
		projKeys = append(projKeys, k)
	}
	sort.Strings(projKeys)
	for _, proj := range projKeys {
		linkTypes := byProject[proj]
		typeSum := 0
		for _, items := range linkTypes {
			typeSum += len(items)
		}
		byProjectOut[proj] = map[string]any{
			"total_links":  typeSum,
			"by_link_type": linkTypes,
		}
	}

	return map[string]any{
		"project_key":               projectKey,
		"total_issues_scanned":      len(issues),
		"total_cross_project_links": totalLinks,
		"by_project":                byProjectOut,
	}, nil
}

// fetchIssuesWithLinks uses JQL + pagination to fetch issues with link fields.
func (f *Fetcher) fetchIssuesWithLinks(ctx context.Context, jql string, maxIssues int) []map[string]any {
	linkFields := []string{"summary", "status", "issuetype", "issuelinks", "parent"}
	search, err := f.searchJQL(ctx, jql, linkFields, maxIssues)
	if err != nil {
		return nil
	}
	issues, _ := search["issues"].([]any)
	if issues == nil {
		return nil
	}
	result := make([]map[string]any, 0, len(issues))
	for _, iss := range issues {
		im, _ := iss.(map[string]any)
		if im != nil {
			result = append(result, im)
		}
	}
	return result[:minInt(maxIssues, len(result))]
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (f *Fetcher) searchJQL(ctx context.Context, jql string, fields []string, maxResults int) (map[string]any, error) {
	req, err := f.client.NewRequest(ctx, "GET", "/rest/api/2/search", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("jql", jql)
	q.Set("maxResults", fmt.Sprintf("%d", maxResults))
	if len(fields) > 0 {
		q.Set("fields", strings.Join(fields, ","))
	}
	req.URL.RawQuery = q.Encode()
	var raw map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: search JQL: %w", err)
	}
	return raw, nil
}
