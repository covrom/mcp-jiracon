package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	core "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server/core"

	"github.com/mark3labs/mcp-go/mcp"
)

type jiraCreateIssueTool struct{}

func (t jiraCreateIssueTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	a := base.Args
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	fields := map[string]any{
		"project":   map[string]any{"key": a.String("project_key")},
		"summary":   a.String("summary"),
		"issuetype": map[string]any{"name": a.String("issue_type")},
	}
	if assignee := a.String("assignee"); assignee != "" {
		fields["assignee"] = map[string]any{"name": assignee}
	}
	if desc := a.String("description"); desc != "" {
		fields["description"] = desc
	}
	if comps := a.String("components"); comps != "" {
		compList := []map[string]any{}
		for _, name := range core.SplitComma(comps) {
			compList = append(compList, map[string]any{"name": name})
		}
		fields["components"] = compList
	}
	if af := a.String("additional_fields"); af != "" {
		var extra map[string]any
		if err := json.Unmarshal([]byte(af), &extra); err != nil {
			return nil, fmt.Errorf("invalid JSON in additional_fields: %w", err)
		}
		// Handle special keys: epicKey/epic_link → "Epic Link", parent → fields["parent"].
		for k, v := range extra {
			if k == "project" || k == "summary" || k == "issuetype" ||
				k == "assignee" || k == "description" || k == "components" {
				continue
			}
			if k == "epicKey" || k == "epic_link" {
				fields["Epic Link"] = v
				continue
			}
			if k == "parent" {
				fields["parent"] = v
				continue
			}
			fields[k] = v
		}
	}
	body := map[string]any{"fields": fields}
	var raw map[string]any
	if err := base.Jira.Client().Post(ctx, "/rest/api/2/issue", body, &raw); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"message": "Issue created successfully", "issue": raw})), nil
}

type jiraUpdateIssueTool struct{}

func (t jiraUpdateIssueTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	a := base.Args
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	var fields map[string]any
	if err := json.Unmarshal([]byte(a.String("fields")), &fields); err != nil {
		return nil, fmt.Errorf("invalid JSON in fields: %w", err)
	}
	if af := a.String("additional_fields"); af != "" {
		var extra map[string]any
		if err := json.Unmarshal([]byte(af), &extra); err != nil {
			return nil, fmt.Errorf("invalid JSON in additional_fields: %w", err)
		}
		for k, v := range extra {
			fields[k] = v
		}
	}
	if comps := a.String("components"); comps != "" {
		compList := []map[string]any{}
		for _, name := range core.SplitComma(comps) {
			compList = append(compList, map[string]any{"name": name})
		}
		fields["components"] = compList
	}
	body := map[string]any{"fields": fields}
	var raw map[string]any
	k := a.String("issue_key")
	if err := base.Jira.Client().Put(ctx, "/rest/api/2/issue/"+k, body, &raw); err != nil {
		return nil, err
	}
	// Attachments
	attPaths := []string{}
	if attRaw := a.String("attachments"); attRaw != "" {
		var jsonPaths []string
		if json.Unmarshal([]byte(attRaw), &jsonPaths) == nil {
			attPaths = jsonPaths
		} else {
			attPaths = core.SplitComma(attRaw)
		}
	}
	attResults := []map[string]any{}
	for _, p := range attPaths {
		data, err := os.ReadFile(p)
		if err != nil {
			attResults = append(attResults, map[string]any{"path": p, "error": err.Error()})
			continue
		}
		var att map[string]any
		if err := base.Jira.Client().PostMultipart(ctx,
			"/rest/api/2/issue/"+k+"/attachments",
			nil,
			filepath.Base(p), data, &att); err != nil {
			attResults = append(attResults, map[string]any{"path": p, "error": err.Error()})
			continue
		}
		attResults = append(attResults, map[string]any{"path": p, "result": att})
	}
	resp := map[string]any{"update": raw}
	if len(attResults) > 0 {
		resp["attachment_results"] = attResults
	}
	return mcp.NewToolResultText(core.Stringify(resp)), nil
}

type jiraDeleteIssueTool struct{}

func (t jiraDeleteIssueTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	if err := base.Jira.Client().Delete(ctx, "/rest/api/2/issue/"+base.Args.String("issue_key")); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(`{"success":true}`), nil
}

type jiraAssignIssueTool struct{}

func (t jiraAssignIssueTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	k := base.Args.String("issue_key")
	a := base.Args.String("assignee")

	// Try to parse assignee as a JSON object from jira_search_assignable_users.
	if a != "" && strings.HasPrefix(a, "{") {
		var parsed map[string]any
		if err := json.Unmarshal([]byte(a), &parsed); err == nil {
			if name, ok := parsed["name"].(string); ok && name != "" {
				a = name
			} else if key, ok := parsed["key"].(string); ok && key != "" {
				a = key
			} else if accountID, ok := parsed["accountId"].(string); ok && accountID != "" {
				a = accountID
			}
		}
	}
	var body any = map[string]any{"name": a}
	if a == "" {
		// Unassign: set to null (send empty string name)
		body = map[string]any{"name": nil}
	}
	var raw map[string]any
	if err := base.Jira.Client().Put(ctx, "/rest/api/2/issue/"+k+"/assignee", body, &raw); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(raw)), nil
}

type jiraTransitionIssueTool struct{}

func (t jiraTransitionIssueTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	k := base.Args.String("issue_key")
	tid := base.Args.String("transition_id")
	var fields map[string]any
	if s := base.Args.String("fields"); s != "" {
		if err := json.Unmarshal([]byte(s), &fields); err != nil {
			return nil, fmt.Errorf("invalid JSON in fields: %w", err)
		}
	}
	issue, err := base.Jira.TransitionIssue(ctx, k, tid, fields, base.Args.String("comment"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(issue)), nil
}

type jiraAddWorklogTool struct{}

func (t jiraAddWorklogTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	comment := base.Args.String("comment")
	started := base.Args.String("started")
	origEstimate := base.Args.String("original_estimate")
	remainingEstimate := base.Args.String("remaining_estimate")
	r, err := base.Jira.AddWorklogExtended(ctx, base.Args.String("issue_key"), base.Args.String("time_spent"), comment, started, origEstimate, remainingEstimate)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"message": "Worklog added successfully", "worklog": r})), nil
}

type jiraAddWatcherTool struct{}

func (t jiraAddWatcherTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	if err := base.Jira.AddWatcher(ctx, base.Args.String("issue_key"), base.Args.String("user_identifier")); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(`{"success":true}`), nil
}

type jiraRemoveWatcherTool struct{}

func (t jiraRemoveWatcherTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	userID := base.Args.String("user_identifier")
	if userID == "" {
		userID = base.Args.String("username")
	}
	if userID == "" {
		return nil, fmt.Errorf("user_identifier or username is required")
	}
	if err := base.Jira.RemoveWatcher(ctx, base.Args.String("issue_key"), userID); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(`{"success":true}`), nil
}

type jiraLinkToEpicTool struct{}

func (t jiraLinkToEpicTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	issueKey := base.Args.String("issue_key")
	epicKey := base.Args.String("epic_key")
	issue, err := base.Jira.LinkIssueToEpic(ctx, issueKey, epicKey)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"message": fmt.Sprintf("Issue %s linked to epic %s", issueKey, epicKey), "issue": issue.ToSimplifiedDict()})), nil
}

type jiraCreateIssueLinkTool struct{}

func (t jiraCreateIssueLinkTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.CreateIssueLink(ctx, base.Args.String("link_type"), base.Args.String("inward_issue_key"), base.Args.String("outward_issue_key"))
	if err != nil {
		return nil, err
	}
	// Add comment if provided.
	if comment := base.Args.String("comment"); comment != "" {
		visJSON := base.Args.String("comment_visibility")
		issueKey := base.Args.String("inward_issue_key")
		commentBody := map[string]any{"body": comment}
		if visJSON != "" {
			var vis map[string]any
			if json.Unmarshal([]byte(visJSON), &vis) == nil {
				commentBody["visibility"] = vis
			}
		}
		base.Jira.Client().Post(ctx, "/rest/api/2/issue/"+issueKey+"/comment", commentBody, nil)
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraCreateRemoteIssueLinkTool struct{}

func (t jiraCreateRemoteIssueLinkTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.CreateRemoteIssueLink(ctx, base.Args.String("issue_key"), base.Args.String("url"), base.Args.String("title"), base.Args.String("summary"), base.Args.String("relationship"), base.Args.String("icon_url"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraRemoveIssueLinkTool struct{}

func (t jiraRemoveIssueLinkTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	if err := base.Jira.RemoveIssueLink(ctx, base.Args.String("link_id")); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(`{"success":true}`), nil
}

type jiraAddCommentTool struct{}

func (t jiraAddCommentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	a := base.Args
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	// Guard: reject comments on internal-only projects.
	issueKey := a.String("issue_key")
	if applies, checkErr := base.Jira.UserAppliesToInternalOnly(ctx, issueKey); checkErr == nil && applies {
		return nil, fmt.Errorf("jira: adding comments to issues in internal-only projects is not allowed")
	}

	// Support both 'body' (Python-compatible) and 'comment' (legacy Go) parameter names.
	commentBody := a.String("body")
	if commentBody == "" {
		commentBody = a.String("comment")
	}
	format := a.String("comment_format")
	var body map[string]any
	if format == "adf" {
		// Try to parse as already-formatted ADF JSON
		var adf map[string]any
		if err := json.Unmarshal([]byte(commentBody), &adf); err == nil {
			body = map[string]any{"body": adf}
		} else {
			// Wrap plain text in a simple ADF paragraph structure
			body = map[string]any{
				"body": map[string]any{
					"version": 1,
					"type":    "doc",
					"content": []map[string]any{
						{
							"type": "paragraph",
							"content": []map[string]any{
								{
									"type": "text",
									"text": commentBody,
								},
							},
						},
					},
				},
			}
		}
	} else {
		body = map[string]any{"body": commentBody}
	}
	// Apply visibility if provided.
	if visJSON := a.String("visibility"); visJSON != "" {
		var vis map[string]any
		if err := json.Unmarshal([]byte(visJSON), &vis); err == nil {
			body["visibility"] = vis
		}
	}
	var raw map[string]any
	// Check if routing through ServiceDesk API for JSM requests.
	isPublicComment := false
	usePublic := false
	if a.Has("public") {
		usePublic = true
		isPublicComment = a.Bool("public", false)
	}
	endpoint := fmt.Sprintf("/rest/api/2/issue/%s/comment", issueKey)
	if usePublic {
		// Route through JSM ServiceDesk API for request comments.
		body["public"] = isPublicComment
		endpoint = fmt.Sprintf("/rest/servicedeskapi/request/%s/comment", issueKey)
	}
	if err := base.Jira.Client().Post(ctx, endpoint, body, &raw); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(raw)), nil
}

type jiraEditCommentTool struct{}

func (t jiraEditCommentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.EditComment(ctx, base.Args.String("issue_key"), base.Args.String("comment_id"), base.Args.String("body"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraCreateVersionTool struct{}

func (t jiraCreateVersionTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.CreateVersion(ctx, base.Args.String("project_key"), base.Args.String("name"),
		base.Args.String("start_date"), base.Args.String("release_date"), base.Args.String("description"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraUpdateVersionTool struct{}

func (t jiraUpdateVersionTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	name := base.Args.String("name")
	desc := base.Args.String("description")
	sd := base.Args.String("start_date")
	rd := base.Args.String("release_date")
	var archived, released *bool
	if base.Args.Has("archived") {
		v := base.Args.Bool("archived", false)
		archived = &v
	}
	if base.Args.Has("released") {
		v := base.Args.Bool("released", false)
		released = &v
	}
	r, err := base.Jira.UpdateVersion(ctx, base.Args.String("version_id"), name, desc, archived, released, sd, rd)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraBatchCreateVersionsTool struct{}

func (t jiraBatchCreateVersionsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	var vers []map[string]any
	if err := json.Unmarshal([]byte(base.Args.String("versions")), &vers); err != nil {
		return nil, fmt.Errorf("invalid JSON in versions: %w", err)
	}
	results := []map[string]any{}
	pk := base.Args.String("project_key")
	for _, v := range vers {
		nm, _ := v["name"].(string)
		sd, _ := v["startDate"].(string)
		rd, _ := v["releaseDate"].(string)
		des, _ := v["description"].(string)
		r, err := base.Jira.CreateVersion(ctx, pk, nm, sd, rd, des)
		if err != nil {
			results = append(results, map[string]any{"error": err.Error()})
			continue
		}
		results = append(results, map[string]any{"version": r})
	}
	return mcp.NewToolResultText(core.Stringify(results)), nil
}

type jiraBatchCreateIssuesTool struct{}

func (t jiraBatchCreateIssuesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	var items []map[string]any
	if err := json.Unmarshal([]byte(base.Args.String("issues")), &items); err != nil {
		return nil, fmt.Errorf("invalid JSON in issues: %w", err)
	}

	// If validate_only, return validated summary.
	if base.Args.Bool("validate_only", false) {
		validated := make([]map[string]any, 0, len(items))
		for _, item := range items {
			entry := map[string]any{
				"project_key": item["project_key"],
				"summary":     item["summary"],
				"issue_type":  item["issue_type"],
				"valid":       item["project_key"] != nil && item["summary"] != nil && item["issue_type"] != nil,
			}
			validated = append(validated, entry)
		}
		return mcp.NewToolResultText(core.Stringify(map[string]any{"validated": validated, "total": len(items)})), nil
	}

	results := []map[string]any{}
	for _, item := range items {
		pk, _ := item["project_key"].(string)
		summary, _ := item["summary"].(string)
		it, _ := item["issue_type"].(string)
		fields := map[string]any{
			"project":   map[string]any{"key": pk},
			"summary":   summary,
			"issuetype": map[string]any{"name": it},
		}
		if desc, _ := item["description"].(string); desc != "" {
			fields["description"] = desc
		}
		if assignee, _ := item["assignee"].(string); assignee != "" {
			fields["assignee"] = map[string]any{"name": assignee}
		}
		if comps, _ := item["components"].([]any); comps != nil {
			compObjs := make([]map[string]any, 0, len(comps))
			for _, c := range comps {
				if cs, ok := c.(string); ok {
					compObjs = append(compObjs, map[string]any{"name": cs})
				}
			}
			if len(compObjs) > 0 {
				fields["components"] = compObjs
			}
		}
		body := map[string]any{"fields": fields}
		var raw map[string]any
		if err := base.Jira.Client().Post(ctx, "/rest/api/2/issue", body, &raw); err != nil {
			results = append(results, map[string]any{"error": err.Error()})
			continue
		}
		results = append(results, map[string]any{"issue": raw})
	}
	return mcp.NewToolResultText(core.Stringify(results)), nil
}

type jiraCreateCustomerRequestTool struct{}

func (t jiraCreateCustomerRequestTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	a := base.Args
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	var fields map[string]any
	if err := json.Unmarshal([]byte(a.String("request_field_values")), &fields); err != nil {
		return nil, fmt.Errorf("invalid JSON in request_field_values: %w", err)
	}
	// Parse request_participants: try JSON array first, fall back to comma-separated.
	participants := a.String("request_participants")
	parsedParticipants := participants
	var jsonParts []string
	if json.Unmarshal([]byte(participants), &jsonParts) == nil {
		parsedParticipants = strings.Join(jsonParts, ",")
	}
	strictOnBehalf := a.Bool("strict_on_behalf", false)
	// Parse attachments if provided.
	var attachments []map[string]any
	if attachmentsRaw := a.String("attachments"); attachmentsRaw != "" {
		if json.Unmarshal([]byte(attachmentsRaw), &attachments) != nil {
			// If it's not valid JSON, leave it as nil.
			attachments = nil
		}
	}
	r, err := base.Jira.CreateCustomerRequest(ctx, a.String("service_desk_id"), a.String("request_type_id"), fields, a.String("raise_on_behalf_of"), parsedParticipants, strictOnBehalf, attachments)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraCreateSprintTool struct{}

func (t jiraCreateSprintTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.CreateSprint(ctx, base.Args.String("board_id"), base.Args.String("name"),
		base.Args.String("start_date"), base.Args.String("end_date"),
		base.Args.String("goal"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraUpdateSprintTool struct{}

func (t jiraUpdateSprintTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.UpdateSprint(
		ctx,
		base.Args.String("sprint_id"),
		base.Args.String("name"),
		base.Args.String("state"),
		base.Args.String("start_date"),
		base.Args.String("end_date"),
		base.Args.String("goal"),
	)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraAddIssuesToSprintTool struct{}

func (t jiraAddIssuesToSprintTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	sprintID := base.Args.String("sprint_id")
	issueKeys := core.SplitComma(base.Args.String("issue_keys"))
	if err := base.Jira.AddIssuesToSprint(ctx, sprintID, issueKeys); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"message": fmt.Sprintf("Successfully added %d issue(s) to sprint", len(issueKeys)), "sprint_id": sprintID, "issue_keys": issueKeys})), nil
}

type jiraMoveIssuesToBacklogTool struct{}

func (t jiraMoveIssuesToBacklogTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	issueKeys := core.SplitComma(base.Args.String("issue_keys"))
	if err := base.Jira.MoveIssuesToBacklog(ctx, issueKeys); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"message": fmt.Sprintf("Successfully moved %d issue(s) to backlog", len(issueKeys)), "issue_keys": issueKeys})), nil
}

func (jiraCreateIssueTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_create_issue", mcp.WithDescription("Create a Jira issue"),
		mcp.WithString("project_key", mcp.Required()),
		mcp.WithString("summary", mcp.Required()),
		mcp.WithString("issue_type", mcp.Required()),
		mcp.WithString("assignee", mcp.Description("Assignee's user identifier (email, display name, or account ID)")),
		mcp.WithString("description", mcp.Description("Issue description")),
		mcp.WithString("components", mcp.Description("Comma-separated list of component names")),
		mcp.WithString("additional_fields", mcp.Description("JSON string of additional fields to set")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraUpdateIssueTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_update_issue", mcp.WithDescription("Update a Jira issue"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithString("fields", mcp.Required()),
		mcp.WithString("additional_fields", mcp.Description("Additional JSON fields to merge into the update")),
		mcp.WithString("components", mcp.Description("Comma-separated component names (replaces existing components)")),
		mcp.WithString("attachments", mcp.Description("JSON array or comma-separated list of file paths to attach")),
		mcp.WithString("return_fields", mcp.Description("Comma-separated fields to return (default *all). Use to save tokens.")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraDeleteIssueTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_delete_issue", mcp.WithDescription("Delete a Jira issue"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraAssignIssueTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_assign_issue", mcp.WithDescription("Assign a Jira issue (or unassign if assignee is empty)"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithString("assignee", mcp.Description("Assignee's user identifier. Leave empty to unassign.")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraTransitionIssueTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_transition_issue", mcp.WithDescription("Transition a Jira issue"),
		mcp.WithString("issue_key", mcp.Required()), mcp.WithString("transition_id", mcp.Required()),
		mcp.WithString("fields"), mcp.WithString("comment"),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraAddWorklogTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_add_worklog", mcp.WithDescription("Add worklog to a Jira issue"),
		mcp.WithString("issue_key", mcp.Required()), mcp.WithString("time_spent", mcp.Required()),
		mcp.WithString("comment", mcp.Description("Comment for the worklog in Markdown format")),
		mcp.WithString("started", mcp.Description("Start time in ISO format (e.g., '2023-08-01T12:00:00.000+0000')")),
		mcp.WithString("original_estimate", mcp.Description("New value for the original estimate (e.g., '1h 30m')")),
		mcp.WithString("remaining_estimate", mcp.Description("New value for the remaining estimate")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraAddWatcherTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_add_watcher", mcp.WithDescription("Add a watcher to a Jira issue"),
		mcp.WithString("issue_key", mcp.Required()), mcp.WithString("user_identifier", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraRemoveWatcherTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_remove_watcher", mcp.WithDescription("Remove a watcher from a Jira issue"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithString("user_identifier", mcp.Required(), mcp.Description("User identifier (username for Server/DC, account ID for Cloud)")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraLinkToEpicTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_link_to_epic", mcp.WithDescription("Link an issue to an epic"),
		mcp.WithString("issue_key", mcp.Required()), mcp.WithString("epic_key", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraCreateIssueLinkTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_create_issue_link", mcp.WithDescription("Create an issue link"),
		mcp.WithString("link_type", mcp.Required()), mcp.WithString("inward_issue_key", mcp.Required()), mcp.WithString("outward_issue_key", mcp.Required()),
		mcp.WithString("comment", mcp.Description("Comment to add to the link")),
		mcp.WithString("comment_visibility", mcp.Description("JSON for comment visibility, e.g. '{\"type\":\"group\",\"value\":\"jira-users\"}'")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraCreateRemoteIssueLinkTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_create_remote_issue_link", mcp.WithDescription("Create a remote issue link"),
		mcp.WithString("issue_key", mcp.Required()), mcp.WithString("url", mcp.Required()), mcp.WithString("title", mcp.Required()),
		mcp.WithString("summary", mcp.Description("Summary description of the remote link")),
		mcp.WithString("relationship", mcp.Description("Link relationship (e.g. 'clones', 'relates to')")),
		mcp.WithString("icon_url", mcp.Description("Icon URL for the remote link (16x16)")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraRemoveIssueLinkTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_remove_issue_link", mcp.WithDescription("Remove an issue link"),
		mcp.WithString("link_id", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraAddCommentTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_add_comment", mcp.WithDescription("Add a comment to a Jira issue"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithString("body", mcp.Required(), mcp.Description("Comment text (Markdown or plain text)")),
		mcp.WithString("comment", mcp.Description("Alias for body (deprecated, use 'body' instead)")),
		mcp.WithString("comment_format", mcp.Description("Format: 'text' (plain text, default) or 'adf' (Atlassian Document Format)")),
		mcp.WithBoolean("public", mcp.Description("For JSM requests: set comment visibility (public/internal). Routes through ServiceDesk API.")),
		mcp.WithString("visibility", mcp.Description(`JSON for comment visibility, e.g. {"type":"group","value":"jira-users"}`)),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraEditCommentTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_edit_comment", mcp.WithDescription("Edit a Jira comment"),
		mcp.WithString("issue_key", mcp.Required()), mcp.WithString("comment_id", mcp.Required()), mcp.WithString("body", mcp.Required()),
		mcp.WithString("visibility", mcp.Description("JSON for comment visibility")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraCreateVersionTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_create_version", mcp.WithDescription("Create a fix version"),
		mcp.WithString("project_key", mcp.Required()), mcp.WithString("name", mcp.Required()),
		mcp.WithString("start_date", mcp.Description("Start date (YYYY-MM-DD)")),
		mcp.WithString("release_date", mcp.Description("Release date (YYYY-MM-DD)")),
		mcp.WithString("description", mcp.Description("Description of the version")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraUpdateVersionTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_update_version", mcp.WithDescription("Update a fix version"),
		mcp.WithString("version_id", mcp.Required()),
		mcp.WithString("name"),
		mcp.WithString("description"),
		mcp.WithString("start_date", mcp.Description("Start date (YYYY-MM-DD)")),
		mcp.WithString("release_date", mcp.Description("Release date (YYYY-MM-DD)")),
		mcp.WithBoolean("archived"),
		mcp.WithBoolean("released"),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraBatchCreateVersionsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_batch_create_versions", mcp.WithDescription("Batch create versions"),
		mcp.WithString("project_key", mcp.Required()), mcp.WithString("versions", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraBatchCreateIssuesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_batch_create_issues", mcp.WithDescription("Batch create issues"),
		mcp.WithString("issues", mcp.Required()),
		mcp.WithBoolean("validate_only", mcp.Description("If true, only validates the issues without creating them")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraCreateCustomerRequestTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_create_customer_request", mcp.WithDescription("Create JSM customer request"),
		mcp.WithString("service_desk_id", mcp.Required()), mcp.WithString("request_type_id", mcp.Required()), mcp.WithString("request_field_values", mcp.Required()),
		mcp.WithString("raise_on_behalf_of", mcp.Description("User identifier to raise request on behalf of")),
		mcp.WithString("request_participants", mcp.Description("JSON array or comma-separated list of participant identifiers")),
		mcp.WithBoolean("strict_on_behalf", mcp.Description("If true, validate that the raise_on_behalf_of user exists")),
		mcp.WithString("attachments", mcp.Description(`JSON array of base64-encoded files: [{"filename":"...","mime_type":"...","base64":"..."}]`)),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraCreateSprintTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_create_sprint", mcp.WithDescription("Create a sprint"),
		mcp.WithString("board_id", mcp.Required()), mcp.WithString("name", mcp.Required()),
		mcp.WithString("start_date", mcp.Required()), mcp.WithString("end_date", mcp.Required()),
		mcp.WithString("goal", mcp.Description("Sprint goal")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraUpdateSprintTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_update_sprint", mcp.WithDescription("Update a sprint"),
		mcp.WithString("sprint_id", mcp.Required()),
		mcp.WithString("name"), mcp.WithString("state"),
		mcp.WithString("start_date"), mcp.WithString("end_date"),
		mcp.WithString("goal"),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraAddIssuesToSprintTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_add_issues_to_sprint", mcp.WithDescription("Add issues to sprint"),
		mcp.WithString("sprint_id", mcp.Required()), mcp.WithString("issue_keys", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (jiraMoveIssuesToBacklogTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_move_issues_to_backlog", mcp.WithDescription("Move to backlog"),
		mcp.WithString("issue_keys", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func WriteTools() []core.Runnable {
	return []core.Runnable{
		&jiraCreateIssueTool{},
		&jiraUpdateIssueTool{},
		&jiraDeleteIssueTool{},
		&jiraAssignIssueTool{},
		&jiraTransitionIssueTool{},
		&jiraAddWorklogTool{},
		&jiraAddWatcherTool{},
		&jiraRemoveWatcherTool{},
		&jiraLinkToEpicTool{},
		&jiraCreateIssueLinkTool{},
		&jiraCreateRemoteIssueLinkTool{},
		&jiraRemoveIssueLinkTool{},
		&jiraAddCommentTool{},
		&jiraEditCommentTool{},
		&jiraCreateVersionTool{},
		&jiraUpdateVersionTool{},
		&jiraBatchCreateVersionsTool{},
		&jiraBatchCreateIssuesTool{},
		&jiraCreateCustomerRequestTool{},
		&jiraCreateSprintTool{},
		&jiraUpdateSprintTool{},
		&jiraAddIssuesToSprintTool{},
		&jiraMoveIssuesToBacklogTool{},
	}
}
