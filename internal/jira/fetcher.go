package jira

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/atlassian"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
	jmodels "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/jira"
)

type Fetcher struct {
	client *atlassian.Client
	cfg    *config.JiraConfig
}

func New(cfg *config.JiraConfig, c *atlassian.Client) *Fetcher { return &Fetcher{client: c, cfg: cfg} }
func (f *Fetcher) Client() *atlassian.Client                   { return f.client }
func (f *Fetcher) APIVersion() string                          { return "2" }
func (f *Fetcher) NewRequest(ctx context.Context, method, path string) (*http.Request, error) {
	return f.client.NewRequest(ctx, method, path, nil)
}
func (f *Fetcher) Do(req *http.Request, out any) error { return f.client.Do(req, out) }

var defaultReadFields = []string{"priority", "updated", "labels", "issuetype", "summary", "assignee", "description", "created", "reporter", "status", "versions"}

func (f *Fetcher) GetIssue(ctx context.Context, issueKey string, fields []string) (*jmodels.JiraIssue, error) {
	return f.GetIssueWithExpand(ctx, issueKey, fields, "")
}

// GetIssueWithExpand fetches a Jira issue with optional expand parameter.
func (f *Fetcher) GetIssueWithExpand(ctx context.Context, issueKey string, fields []string, expand string) (*jmodels.JiraIssue, error) {
	if issueKey == "" {
		return nil, fmt.Errorf("jira: issue key required")
	}
	path := fmt.Sprintf("/rest/api/2/issue/%s", issueKey)
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if len(fields) > 0 {
		fs := []string{}
		for _, f := range fields {
			if f != "" {
				fs = append(fs, f)
			}
		}
		if len(fs) > 0 {
			q.Set("fields", strings.Join(fs, ","))
		}
	}
	if expand != "" {
		q.Set("expand", expand)
	}
	req.URL.RawQuery = q.Encode()
	var raw map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: get issue %s: %w", issueKey, err)
	}
	issue := &jmodels.JiraIssue{}
	if err := issue.FromAPIResponse(raw); err != nil {
		return nil, fmt.Errorf("jira: parse issue %s: %w", issueKey, err)
	}
	return issue, nil
}

// GetIssueFull fetches a Jira issue with all optional parameters:
// properties (comma-separated), updateHistory, commentLimit.
func (f *Fetcher) GetIssueFull(ctx context.Context, issueKey string, fields []string, expand, properties string, updateHistory bool, commentLimit int) (*jmodels.JiraIssue, error) {
	if issueKey == "" {
		return nil, fmt.Errorf("jira: issue key required")
	}
	path := fmt.Sprintf("/rest/api/2/issue/%s", issueKey)
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if len(fields) > 0 {
		fs := []string{}
		for _, f := range fields {
			if f != "" {
				fs = append(fs, f)
			}
		}
		if len(fs) > 0 {
			q.Set("fields", strings.Join(fs, ","))
		}
	}
	if expand != "" {
		q.Set("expand", expand)
	}
	if properties != "" {
		q.Set("properties", properties)
	}
	if !updateHistory {
		q.Set("updateHistory", "false")
	}
	req.URL.RawQuery = q.Encode()
	var raw map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: get issue %s: %w", issueKey, err)
	}
	issue := &jmodels.JiraIssue{}
	if err := issue.FromAPIResponse(raw); err != nil {
		return nil, fmt.Errorf("jira: parse issue %s: %w", issueKey, err)
	}
	// Apply comment_limit client-side: trim comments to the requested limit.
	if commentLimit > 0 {
		if fieldsRaw, ok := raw["fields"].(map[string]any); ok {
			if comment, ok := fieldsRaw["comment"].(map[string]any); ok {
				if comments, ok := comment["comments"].([]any); ok && len(comments) > commentLimit {
					comment["comments"] = comments[:commentLimit]
					comment["maxResults"] = float64(commentLimit)
				}
			}
		}
	}
	return issue, nil
}

func (f *Fetcher) GetUserProfile(ctx context.Context, identifier string) (*jmodels.JiraUser, error) {
	if identifier == "" {
		return nil, fmt.Errorf("jira: user identifier required")
	}
	path := "/rest/api/2/user"
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("username", identifier)
	req.URL.RawQuery = q.Encode()
	var raw map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("jira: get user %s: %w", identifier, err)
	}
	user := &jmodels.JiraUser{}
	if err := user.FromAPIResponse(raw); err != nil {
		return nil, fmt.Errorf("jira: parse user: %w", err)
	}
	return user, nil
}

func (f *Fetcher) IsInternalOnlyProject(projectKey string) bool {
	for _, p := range f.cfg.InternalOnlyProjects {
		if strings.EqualFold(p, projectKey) {
			return true
		}
	}
	return false
}

func (f *Fetcher) UserAppliesToInternalOnly(ctx context.Context, issueKey string) (bool, error) {
	issue, err := f.GetIssue(ctx, issueKey, []string{"project"})
	if err != nil {
		return false, err
	}
	proj, ok := issue.Fields["project"].(map[string]any)
	if !ok {
		return false, nil
	}
	key, ok := proj["key"].(string)
	if !ok {
		return false, nil
	}
	return f.IsInternalOnlyProject(key), nil
}
