package jira

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	core "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server/core"

	"github.com/mark3labs/mcp-go/mcp"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/utils"
)

// extractStatusChanges parses the Jira changelog for status transitions.
func extractStatusChanges(changelog map[string]any) []map[string]any {
	histories, _ := changelog["histories"].([]any)
	if histories == nil {
		return nil
	}
	var changes []map[string]any
	for _, h := range histories {
		hm, ok := h.(map[string]any)
		if !ok {
			continue
		}
		items, _ := hm["items"].([]any)
		for _, item := range items {
			im, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if field, _ := im["field"].(string); field == "status" {
				changes = append(changes, map[string]any{
					"created": hm["created"],
					"from":    im["fromString"],
					"to":      im["toString"],
					"author":  hm["author"],
				})
			}
		}
	}
	return changes
}

// summarizeStatusTimes computes total time spent in each status from changelog.
func summarizeStatusTimes(statusChanges []map[string]any) map[string]float64 {
	if len(statusChanges) < 2 {
		return nil
	}
	summary := map[string]float64{}
	for i := 0; i < len(statusChanges)-1; i++ {
		fromTime, _ := statusChanges[i]["created"].(string)
		toTime, _ := statusChanges[i+1]["created"].(string)
		status, _ := statusChanges[i]["to"].(string)

		t1, err1 := time.Parse(time.RFC3339, fromTime)
		t2, err2 := time.Parse(time.RFC3339, toTime)
		if err1 == nil && err2 == nil && t2.After(t1) {
			summary[status] += t2.Sub(t1).Hours()
		}
	}
	return summary
}

type jiraGetIssueTool struct{}

func (t jiraGetIssueTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	a := base.Args
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	var fields []string
	if fs := a.String("fields"); fs != "" {
		fields = core.SplitComma(fs)
	} else {
		fields = []string{"priority", "updated", "labels", "issuetype", "summary", "assignee", "description", "created", "reporter", "status"}
	}
	includeSet := map[string]bool{}
	for _, sec := range core.SplitComma(a.String("include")) {
		includeSet[strings.ToLower(sec)] = true
	}
	useDisplayNames := a.Bool("use_display_names", false)
	// Build expand string.
	expand := a.String("expand")
	if includeSet["changelog"] {
		if expand != "" {
			expand = expand + ",changelog"
		} else {
			expand = "changelog"
		}
	}
	// use_display_names adds 'names' to expand independently of changelog.
	if useDisplayNames {
		if expand != "" {
			expand = expand + ",names"
		} else {
			expand = "names"
		}
	}
	commentLimit := 10
	if v, ok := a.Num("comment_limit"); ok {
		commentLimit = int(v)
	}
	properties := a.String("properties")
	updateHistory := a.Bool("update_history", true)

	issue, err := base.Jira.GetIssueFull(ctx, a.String("issue_key"), fields, expand, properties, updateHistory, commentLimit)
	if err != nil {
		return nil, err
	}
	result := issue.ToSimplifiedDict()
	// Inline enrichment sections (DC-compatible).
	if includeSections := a.String("include"); includeSections != "" {
		for _, sec := range core.SplitComma(includeSections) {
			switch strings.ToLower(sec) {
			case "transitions":
				trans, tErr := base.Jira.GetTransitions(ctx, issue.Key)
				if tErr == nil {
					result["transitions"] = trans
				}
			case "watchers":
				w, wErr := base.Jira.GetIssueWatchers(ctx, issue.Key)
				if wErr == nil {
					result["watchers"] = w
				}
			case "comments":
				// Fetch comments via API when not already in fields.
				if !strings.Contains(strings.Join(fields, ","), "comment") {
					commentsReq, cErr := base.Jira.Client().NewRequest(ctx, "GET", "/rest/api/2/issue/"+issue.Key+"/comment", nil)
					if cErr == nil {
						var commentsRaw map[string]any
						if cErr2 := base.Jira.Client().Do(commentsReq, &commentsRaw); cErr2 == nil {
							result["comments"] = commentsRaw
						}
					}
				}
			case "worklogs":
				wl, wlErr := base.Jira.GetWorklogs(ctx, issue.Key)
				if wlErr == nil {
					simplified := make([]map[string]any, len(wl))
					for i, w := range wl {
						simplified[i] = w.ToSimplifiedDict()
					}
					result["worklogs"] = simplified
				}
			}
		}
	}
	// Changelog enrichment.
	if includeSet["changelog"] && issue.Changelog != nil {
		result["changelogs"] = issue.Changelog
	}
	// Remote links enrichment.
	if includeSet["remote_links"] {
		req2, rlErr := base.Jira.Client().NewRequest(ctx, "GET", "/rest/api/2/issue/"+issue.Key+"/remotelink", nil)
		if rlErr == nil {
			var raw map[string]any
			if rlErr2 := base.Jira.Client().Do(req2, &raw); rlErr2 == nil {
				result["remote_links"] = raw
			}
		}
	}
	return mcp.NewToolResultText(core.Stringify(result)), nil
}

type jiraGetProjectIssuesTool struct{}

func (t jiraGetProjectIssuesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	start := base.Args.Int("start_at", 0)
	limit := base.Args.Int("limit", 10)
	r, err := base.Jira.GetProjectIssues(ctx, base.Args.String("project_key"), start, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetIssueDatesTool struct{}

func (t jiraGetIssueDatesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	includeStatusChanges := base.Args.Bool("include_status_changes", true)
	includeStatusSummary := base.Args.Bool("include_status_summary", true)
	// Fetch issue with changelog expansion if needed.
	expand := ""
	if includeStatusChanges || includeStatusSummary {
		expand = "changelog"
	}
	issue, err := base.Jira.GetIssueWithExpand(ctx, base.Args.String("issue_key"), nil, expand)
	if err != nil {
		return nil, err
	}
	created, _ := issue.Fields["created"].(string)
	updated, _ := issue.Fields["updated"].(string)
	resolutiondate, _ := issue.Fields["resolutiondate"].(string)
	duedate, _ := issue.Fields["duedate"].(string)

	// Compute approximate cycle time: time from creation to resolution.
	var cycleTimeHours float64
	if created != "" && resolutiondate != "" {
		if ct, err := time.Parse(time.RFC3339, created); err == nil {
			if rt, err2 := time.Parse(time.RFC3339, resolutiondate); err2 == nil {
				cycleTimeHours = rt.Sub(ct).Hours()
			}
		}
	}
	// Compute approximate lead time: time from creation to last update.
	var leadTimeHours float64
	if created != "" && updated != "" {
		if ct, err := time.Parse(time.RFC3339, created); err == nil {
			if ut, err2 := time.Parse(time.RFC3339, updated); err2 == nil {
				leadTimeHours = ut.Sub(ct).Hours()
			}
		}
	}

	result := map[string]any{
		"key":              issue.Key,
		"created":          created,
		"updated":          updated,
		"duedate":          duedate,
		"resolutiondate":   resolutiondate,
		"cycle_time_hours": cycleTimeHours,
		"cycle_time_days":  cycleTimeHours / 24,
		"lead_time_hours":  leadTimeHours,
		"lead_time_days":   leadTimeHours / 24,
	}

	// Changelog-based status tracking.
	if (includeStatusChanges || includeStatusSummary) && issue.Changelog != nil {
		statusChanges := extractStatusChanges(issue.Changelog)
		if includeStatusChanges {
			result["status_changes"] = statusChanges
		}
		if includeStatusSummary {
			result["status_summary"] = summarizeStatusTimes(statusChanges)
		}
	}
	return mcp.NewToolResultText(core.Stringify(result)), nil
}

type jiraGetIssueSLATool struct{}

func (t jiraGetIssueSLATool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	// Parse metrics parameter.
	metricsReq := base.Args.String("metrics")
	if metricsReq == "" {
		metricsReq = "cycle_time,time_in_status"
	}
	metricsSet := map[string]bool{}
	for _, m := range core.SplitComma(metricsReq) {
		metricsSet[strings.TrimSpace(m)] = true
	}

	workingHoursOnly := base.Args.Bool("working_hours_only", false)
	includeRawDates := base.Args.Bool("include_raw_dates", false)

	// Fetch issue with changelog for time-in-status calculation.
	expand := ""
	if metricsSet["time_in_status"] || metricsSet["first_response_time"] {
		expand = "changelog"
	}
	issue, err := base.Jira.GetIssueWithExpand(ctx, base.Args.String("issue_key"), nil, expand)
	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"key": issue.Key,
	}
	if includeRawDates {
		result["created"] = issue.Fields["created"]
		result["updated"] = issue.Fields["updated"]
		result["resolutiondate"] = issue.Fields["resolutiondate"]
		result["duedate"] = issue.Fields["duedate"]
	}

	created, _ := issue.Fields["created"].(string)
	updated, _ := issue.Fields["updated"].(string)
	resolutiondate, _ := issue.Fields["resolutiondate"].(string)
	duedate, _ := issue.Fields["duedate"].(string)

	// Helper: parse RFC3339 time.
	parseTime := func(s string) (time.Time, bool) {
		t, err := time.Parse(time.RFC3339, s)
		return t, err == nil
	}

	// Helper: working hours between two times.
	workingHoursBetween := func(t1, t2 time.Time) float64 {
		if !workingHoursOnly {
			return t2.Sub(t1).Hours()
		}
		// Simple weekday-based: exclude Sat/Sun, count 8h/day.
		total := 0.0
		current := t1
		for current.Before(t2) {
			weekday := current.Weekday()
			if weekday != time.Saturday && weekday != time.Sunday {
				dayEnd := time.Date(current.Year(), current.Month(), current.Day(), 18, 0, 0, 0, current.Location())
				dayStart := time.Date(current.Year(), current.Month(), current.Day(), 9, 0, 0, 0, current.Location())
				if current.Before(dayStart) {
					current = dayStart
				}
				if dayEnd.After(t2) {
					dayEnd = t2
				}
				if current.Before(dayEnd) {
					total += dayEnd.Sub(current).Hours()
				}
			}
			current = time.Date(current.Year(), current.Month(), current.Day()+1, 0, 0, 0, 0, current.Location())
		}
		return total
	}

	// Metrics calculation.
	if metricsSet["cycle_time"] {
		if ct, ok1 := parseTime(created); ok1 {
			if rt, ok2 := parseTime(resolutiondate); ok2 {
				hours := workingHoursBetween(ct, rt)
				result["cycle_time_hours"] = hours
				result["cycle_time_days"] = hours / 24
			}
		}
	}

	if metricsSet["lead_time"] {
		if ct, ok1 := parseTime(created); ok1 {
			if ut, ok2 := parseTime(updated); ok2 {
				hours := workingHoursBetween(ct, ut)
				result["lead_time_hours"] = hours
				result["lead_time_days"] = hours / 24
			}
		}
	}

	if metricsSet["resolution_time"] {
		if rt, ok := parseTime(resolutiondate); ok {
			if ct, ok2 := parseTime(created); ok2 {
				hours := workingHoursBetween(ct, rt)
				result["resolution_time_hours"] = hours
				result["resolution_time_days"] = hours / 24
			}
		}
	}

	if metricsSet["due_date_compliance"] {
		var isOverdue bool
		var overdueDays float64
		if duedate != "" {
			if dd, err := time.Parse("2006-01-02", duedate); err == nil {
				now := time.Now().UTC().Truncate(24 * time.Hour)
				ddDate := dd.UTC().Truncate(24 * time.Hour)
				if resolutiondate == "" {
					if now.After(ddDate) {
						isOverdue = true
						overdueDays = now.Sub(ddDate).Hours() / 24
					}
				} else {
					if rt, ok := parseTime(resolutiondate); ok {
						if rt.After(ddDate) {
							isOverdue = true
							overdueDays = rt.Sub(ddDate).Hours() / 24
						}
					}
				}
			}
		}
		result["is_overdue"] = isOverdue
		result["overdue_days"] = overdueDays
	}

	if metricsSet["time_in_status"] && issue.Changelog != nil {
		statusChanges := extractStatusChanges(issue.Changelog)
		result["status_changes"] = statusChanges
		result["status_summary_hours"] = summarizeStatusTimes(statusChanges)
		if workingHoursOnly && len(statusChanges) >= 2 {
			// Recalculate status summary with working hours.
			whSummary := map[string]float64{}
			for i := 0; i < len(statusChanges)-1; i++ {
				fromTime, _ := statusChanges[i]["created"].(string)
				toTime, _ := statusChanges[i+1]["created"].(string)
				status, _ := statusChanges[i]["to"].(string)
				if t1, ok1 := parseTime(fromTime); ok1 {
					if t2, ok2 := parseTime(toTime); ok2 {
						if t2.After(t1) {
							whSummary[status] += workingHoursBetween(t1, t2)
						}
					}
				}
			}
			result["status_summary_working_hours"] = whSummary
		}
	}

	if metricsSet["first_response_time"] {
		// First response = time from creation to first comment or status change, whichever comes first.
		firstResponseAt := ""
		if issue.Changelog != nil {
			statusChanges := extractStatusChanges(issue.Changelog)
			if len(statusChanges) > 0 {
				firstResponseAt, _ = statusChanges[0]["created"].(string)
			}
		}
		if firstResponseAt != "" {
			if ct, ok1 := parseTime(created); ok1 {
				if fr, ok2 := parseTime(firstResponseAt); ok2 {
					hours := workingHoursBetween(ct, fr)
					result["first_response_time_hours"] = hours
					result["first_response_time_days"] = hours / 24
				}
			}
		}
	}

	// Include raw field data.
	result["timeoriginalestimate"] = issue.Fields["timeoriginalestimate"]
	result["timeestimate"] = issue.Fields["timeestimate"]
	result["timespent"] = issue.Fields["timespent"]

	if progress, ok := issue.Fields["progress"]; ok {
		result["progress"] = progress
	}
	return mcp.NewToolResultText(core.Stringify(result)), nil
}

type jiraGetIssueImagesTool struct{}

func (t jiraGetIssueImagesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetIssueImages(ctx, base.Args.String("issue_key"))
	if err != nil {
		return nil, err
	}
	contents := []mcp.Content{mcp.NewTextContent(core.Stringify(map[string]any{
		"success":      r["success"],
		"issue_key":    r["issue_key"],
		"total":        r["total"],
		"total_images": r["images"],
		"downloaded":   r["images"],
		"failed":       r["failed"],
	}))}
	fetched, _ := r["fetched"].([]map[string]any)
	for _, img := range fetched {
		b64, _ := img["data_base64"].(string)
		mime, _ := img["mime_type"].(string)
		if b64 == "" {
			continue
		}
		contents = append(contents, mcp.NewImageContent(b64, mime))
	}
	return &mcp.CallToolResult{Content: contents}, nil
}

type jiraDownloadAttachmentsTool struct{}

func (t jiraDownloadAttachmentsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetIssueAttachmentContents(ctx, base.Args.String("issue_key"))
	if err != nil {
		return nil, err
	}

	// Failure from fetcher — return summary text only.
	ok, _ := r["success"].(bool)
	if !ok {
		summary, _ := json.MarshalIndent(r, "", "  ")
		return mcp.NewToolResultText(string(summary)), nil
	}

	issueKey := models.MapString(r, "issue_key")
	if issueKey == "" {
		issueKey = base.Args.String("issue_key")
	}
	total := int(models.MapInt(r, "total"))
	attachmentsList, _ := r["attachments"].([]map[string]any)
	failedList, _ := r["failed"].([]map[string]any)
	if failedList == nil {
		failedList = []map[string]any{}
	}

	contents := []mcp.Content{mcp.NewTextContent(core.Stringify(map[string]any{
		"success":    true,
		"issue_key":  issueKey,
		"total":      total,
		"downloaded": 0, // updated below
		"failed":     failedList,
	}))}

	downloaded := 0
	for _, att := range attachmentsList {
		dataBytes, _ := att["data"].([]byte)
		if dataBytes == nil {
			continue
		}
		filename, _ := att["filename"].(string)
		contentType, _ := att["content_type"].(string)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		if int64(len(dataBytes)) > utils.AttachmentMaxBytes {
			failedList = append(failedList, map[string]any{
				"filename": filename,
				"error":    fmt.Sprintf("Attachment '%s' is %d bytes which exceeds the 50 MB inline limit.", filename, len(dataBytes)),
			})
			continue
		}

		encoded := base64.StdEncoding.EncodeToString(dataBytes)
		isImg, resolvedMIME := utils.IsImageAttachment(contentType, filename)
		downloaded++

		uri := "attachment:///" + issueKey + "/" + filename
		if isImg {
			contents = append(contents, mcp.NewEmbeddedResource(mcp.BlobResourceContents{
				URI:      uri,
				MIMEType: resolvedMIME,
				Blob:     encoded,
			}))
		} else {
			// Non-image → TextContent with base64 payload.
			// Many MCP clients only forward EmbeddedResource blobs
			// whose mimeType is a recognized image format (#1419).
			contents = append(contents, mcp.NewTextContent(core.Stringify(map[string]any{
				"success":   true,
				"issue_key": issueKey,
				"filename":  filename,
				"mime_type": resolvedMIME,
				"encoding":  "base64",
				"content":   encoded,
			})))
		}
	}

	// Update the summary text with the actual download count.
	contents[0] = mcp.NewTextContent(core.Stringify(map[string]any{
		"success":    true,
		"issue_key":  issueKey,
		"total":      total,
		"downloaded": downloaded,
		"failed":     failedList,
	}))

	return &mcp.CallToolResult{Content: contents}, nil
}

type jiraGetServiceDeskForProjectTool struct{}

func (t jiraGetServiceDeskForProjectTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	projectKey := base.Args.String("project_key")
	r, err := base.Jira.GetServiceDesk(ctx, projectKey)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"project_key":  projectKey,
		"service_desk": r,
	})), nil
}

type jiraGetServiceDeskQueuesTool struct{}

func (t jiraGetServiceDeskQueuesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	start := base.Args.Int("start_at", 0)
	limit := base.Args.Int("limit", 50)
	r, err := base.Jira.GetServiceDeskQueues(ctx, base.Args.String("service_desk_id"), start, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetQueueIssuesTool struct{}

func (t jiraGetQueueIssuesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	start := base.Args.Int("start_at", 0)
	limit := base.Args.Int("limit", 50)
	r, err := base.Jira.GetQueueIssues(ctx, base.Args.String("service_desk_id"), base.Args.String("queue_id"), start, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetRequestTypesTool struct{}

func (t jiraGetRequestTypesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	start := base.Args.Int("start_at", 0)
	limit := base.Args.Int("limit", 50)
	r, err := base.Jira.GetRequestTypes(ctx, base.Args.String("service_desk_id"), start, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetRequestTypeFieldsTool struct{}

func (t jiraGetRequestTypeFieldsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetRequestTypeFields(ctx, base.Args.String("service_desk_id"), base.Args.String("request_type_id"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetUserProfileTool struct{}

func (t jiraGetUserProfileTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	userIdentifier := base.Args.String("user_identifier")
	user, err := base.Jira.GetUserProfile(ctx, userIdentifier)
	if err != nil {
		return mcp.NewToolResultText(core.Stringify(map[string]any{
			"success":         false,
			"error":           err.Error(),
			"user_identifier": userIdentifier,
		})), nil
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"success": true,
		"user":    user.ToSimplifiedDict(),
	})), nil
}

type jiraSearchTool struct{}

func (t jiraSearchTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	req2, err := base.Jira.Client().NewRequest(ctx, "GET", "/rest/api/2/search", nil)
	if err != nil {
		return nil, fmt.Errorf("jira: create search request: %w", err)
	}
	q := req2.URL.Query()
	q.Set("jql", base.Args.String("jql"))
	if fields := base.Args.String("fields"); fields != "" {
		q.Set("fields", fields)
	} else {
		q.Set("fields", "priority,updated,labels,issuetype,summary,assignee,description,created,reporter,status")
	}
	limit := base.Args.Int("limit", 10)
	q.Set("maxResults", strconv.Itoa(limit))
	startAt := base.Args.Int("start_at", 0)
	if startAt > 0 {
		q.Set("startAt", strconv.Itoa(startAt))
	}
	expandVal := base.Args.String("expand")
	useDisplayNames := base.Args.Bool("use_display_names", false)
	if useDisplayNames {
		if expandVal != "" {
			expandVal = expandVal + ",names"
		} else {
			expandVal = "names"
		}
	}
	if expandVal != "" {
		q.Set("expand", expandVal)
	}
	req2.URL.RawQuery = q.Encode()
	var raw map[string]any
	if err := base.Jira.Client().Do(req2, &raw); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(raw)), nil
}

type jiraGetTransitionsTool struct{}

func (t jiraGetTransitionsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetTransitions(ctx, base.Args.String("issue_key"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetWorklogTool struct{}

func (t jiraGetWorklogTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	wl, err := base.Jira.GetWorklogs(ctx, base.Args.String("issue_key"))
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, len(wl))
	for i, w := range wl {
		out[i] = w.ToSimplifiedDict()
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"worklogs": out})), nil
}

type jiraGetIssueWatchersTool struct{}

func (t jiraGetIssueWatchersTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetIssueWatchers(ctx, base.Args.String("issue_key"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetLinkTypesTool struct{}

func (t jiraGetLinkTypesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetIssueLinkTypes(ctx)
	if err != nil {
		return nil, err
	}
	// Apply name filter if provided.
	if nameFilter := base.Args.String("name_filter"); nameFilter != "" {
		filtered := []map[string]any{}
		for _, lt := range r {
			if name, _ := lt["name"].(string); name != "" {
				if core.ContainsFold(name, nameFilter) {
					filtered = append(filtered, lt)
				}
			}
		}
		r = filtered
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetAllProjectsTool struct{}

func (t jiraGetAllProjectsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	includeArchived := base.Args.Bool("include_archived", false)
	r, err := base.Jira.GetAllProjects(ctx, includeArchived)
	if err != nil {
		return nil, err
	}
	// Normalize project keys to uppercase.
	for _, p := range r {
		if k, ok := p["key"].(string); ok {
			p["key"] = strings.ToUpper(k)
		}
	}
	// Apply projects_filter if configured.
	if base.App.JiraConfig != nil && base.App.JiraConfig.ProjectsFilter != "" {
		allowedKeys := map[string]struct{}{}
		for _, pk := range strings.Split(base.App.JiraConfig.ProjectsFilter, ",") {
			allowedKeys[strings.TrimSpace(strings.ToUpper(pk))] = struct{}{}
		}
		filtered := make([]map[string]any, 0, len(r))
		for _, p := range r {
			if k, ok := p["key"].(string); ok {
				if _, allowed := allowedKeys[k]; allowed {
					filtered = append(filtered, p)
				}
			}
		}
		r = filtered
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraSearchProjectsTool struct{}

func (t jiraSearchProjectsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	maxResults := base.Args.Int("max_results", 20)
	r, err := base.Jira.SearchProjects(ctx, base.Args.String("query"), maxResults)
	if err != nil {
		return nil, err
	}
	// Apply current_project_ids exclusion filter.
	if excludeIDs := base.Args.String("current_project_ids"); excludeIDs != "" {
		exclude := map[string]bool{}
		for _, id := range core.SplitComma(excludeIDs) {
			exclude[id] = true
		}
		filtered := make([]map[string]any, 0, len(r))
		for _, p := range r {
			if id, _ := p["id"].(string); id != "" && exclude[id] {
				continue
			}
			// Also check numeric ID.
			if idFloat, _ := p["id"].(float64); exclude[fmt.Sprintf("%.0f", idFloat)] {
				continue
			}
			filtered = append(filtered, p)
		}
		r = filtered
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetProjectIssueTypesTool struct{}

func (t jiraGetProjectIssueTypesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetProjectIssueTypes(ctx, base.Args.String("project_key"))
	if err != nil {
		return nil, err
	}
	// Transform to compact format.
	compact := make([]map[string]any, 0, len(r))
	for _, it := range r {
		subtask := false
		if st, ok := it["subtask"].(bool); ok {
			subtask = st
		}
		compact = append(compact, map[string]any{
			"id":                it["id"],
			"name":              it["name"],
			"description":       it["description"],
			"subtask":           subtask,
			"untranslated_name": it["untranslatedName"],
		})
	}
	return mcp.NewToolResultText(core.Stringify(compact)), nil
}

type jiraGetCreateFieldsTool struct{}

func (t jiraGetCreateFieldsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetCreateFields(ctx, base.Args.String("project_key"), base.Args.String("issue_type_id"))
	if err != nil {
		return nil, err
	}
	// Transform to compact format.
	compact := make([]map[string]any, 0, len(r))
	for _, field := range r {
		required := false
		if reqVal, ok := field["required"].(bool); ok {
			required = reqVal
		}
		compact = append(compact, map[string]any{
			"field_id": field["fieldId"],
			"name":     field["name"],
			"required": required,
			"schema":   field["schema"],
		})
	}
	return mcp.NewToolResultText(core.Stringify(compact)), nil
}

type jiraGetProjectVersionsTool struct{}

func (t jiraGetProjectVersionsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetProjectVersions(ctx, base.Args.String("project_key"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetProjectComponentsTool struct{}

func (t jiraGetProjectComponentsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetProjectComponents(ctx, base.Args.String("project_key"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetProjectFieldsTool struct{}

func (t jiraGetProjectFieldsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetProjectFields(ctx, base.Args.String("project_key"))
	if err != nil {
		return nil, err
	}
	// Apply schema_type filter if provided.
	if schemaType := base.Args.String("schema_type"); schemaType != "" {
		filtered := []map[string]any{}
		for _, field := range r {
			if schema, ok := field["schema"].(map[string]any); ok {
				if typ, _ := schema["type"].(string); typ != "" {
					if core.ContainsFold(typ, schemaType) {
						filtered = append(filtered, field)
					}
				}
			}
		}
		r = filtered
	}
	// Apply issue_types filter if provided.
	if issueTypes := base.Args.String("issue_types"); issueTypes != "" {
		allowedTypes := map[string]bool{}
		for _, id := range core.SplitComma(issueTypes) {
			allowedTypes[id] = true
		}
		filtered := []map[string]any{}
		for _, field := range r {
			issueTypeIDs, _ := field["issue_type_ids"].([]any)
			if issueTypeIDs == nil {
				filtered = append(filtered, field)
				continue
			}
			for _, itID := range issueTypeIDs {
				if idStr, ok := itID.(string); ok && allowedTypes[idStr] {
					filtered = append(filtered, field)
					break
				}
			}
		}
		r = filtered
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetAgileBoardsTool struct{}

func (t jiraGetAgileBoardsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	startAt := base.Args.Int("start_at", 0)
	limit := base.Args.Int("limit", 10)
	r, err := base.Jira.GetAgileBoards(ctx,
		base.Args.String("board_name"),
		base.Args.String("project_key"),
		base.Args.String("board_type"),
		startAt, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetBoardIssuesTool struct{}

func (t jiraGetBoardIssuesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	startAt := base.Args.Int("start_at", 0)
	limit := base.Args.Int("limit", 10)
	r, err := base.Jira.GetBoardIssues(ctx, base.Args.String("board_id"), base.Args.String("jql"), startAt, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetSprintsFromBoardTool struct{}

func (t jiraGetSprintsFromBoardTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	startAt := base.Args.Int("start_at", 0)
	limit := base.Args.Int("limit", 10)
	r, err := base.Jira.GetSprintsFromBoard(ctx, base.Args.String("board_id"), base.Args.String("state"), startAt, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetSprintIssuesTool struct{}

func (t jiraGetSprintIssuesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	startAt := base.Args.Int("start_at", 0)
	limit := base.Args.Int("limit", 10)
	r, err := base.Jira.GetSprintIssues(ctx, base.Args.String("sprint_id"), startAt, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraSearchFieldsTool struct{}

func (t jiraSearchFieldsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	limit := base.Args.Int("limit", 10)
	refresh := base.Args.Bool("refresh", false)
	r, err := base.Jira.SearchFields(ctx, base.Args.String("keyword"), limit, refresh)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetFieldOptionsTool struct{}

func (t jiraGetFieldOptionsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetFieldOptions(ctx,
		base.Args.String("field_id"),
		base.Args.String("project_key"),
		base.Args.String("issue_type"),
	)
	if err != nil {
		return nil, err
	}
	// Apply contains filter.
	if contains := base.Args.String("contains"); contains != "" {
		lower := strings.ToLower(contains)
		filtered := []map[string]any{}
		for _, opt := range r {
			val, _ := opt["value"].(string)
			if strings.Contains(strings.ToLower(val), lower) {
				filtered = append(filtered, opt)
				continue
			}
			// Check child options for cascading selects.
			if children, ok := opt["child_options"].([]map[string]any); ok {
				for _, child := range children {
					cv, _ := child["value"].(string)
					if strings.Contains(strings.ToLower(cv), lower) {
						filtered = append(filtered, opt)
						break
					}
				}
			}
		}
		r = filtered
	}
	// Apply return_limit.
	if rl := base.Args.Int("return_limit", 0); rl > 0 && rl < len(r) {
		r = r[:rl]
	}
	// Apply values_only compact format.
	if base.Args.Bool("values_only", false) {
		compact := make([]any, 0, len(r))
		for _, opt := range r {
			val := opt["value"]
			if children, ok := opt["child_options"].([]map[string]any); ok {
				childVals := make([]any, len(children))
				for ci, c := range children {
					childVals[ci] = c["value"]
				}
				compact = append(compact, map[string]any{"value": val, "children": childVals})
			} else {
				compact = append(compact, val)
			}
		}
		return mcp.NewToolResultText(core.Stringify(compact)), nil
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraSearchAssignableUsersTool struct{}

func (t jiraSearchAssignableUsersTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	projectKey := base.Args.String("project_key")
	issueKey := base.Args.String("issue_key")
	if projectKey == "" && issueKey == "" {
		return mcp.NewToolResultText(core.Stringify(map[string]any{
			"success": false,
			"error":   "exactly one of project_key or issue_key must be provided",
		})), nil
	}
	if projectKey != "" && issueKey != "" {
		return mcp.NewToolResultText(core.Stringify(map[string]any{
			"success": false,
			"error":   "exactly one of project_key or issue_key must be provided (not both)",
		})), nil
	}
	limit := 20
	if v, ok := base.Args.Num("limit"); ok {
		if n := int(v); n >= 1 && n <= 1000 {
			limit = n
		}
	}
	r, err := base.Jira.SearchAssignableUsers(ctx, base.Args.String("query"), projectKey, issueKey, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"success": true,
		"count":   len(r),
		"users":   r,
	})), nil
}

type jiraGetIssueDevelopmentInfoTool struct{}

func (t jiraGetIssueDevelopmentInfoTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetIssueDevelopmentInfo(ctx,
		base.Args.String("issue_key"),
		base.Args.String("application_type"),
		base.Args.String("data_type"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetIssuesDevelopmentInfoTool struct{}

func (t jiraGetIssuesDevelopmentInfoTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	keys := core.SplitComma(base.Args.String("issue_keys"))
	r, err := base.Jira.GetIssuesDevelopmentInfo(ctx, keys,
		base.Args.String("application_type"),
		base.Args.String("data_type"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetProjectEpicHierarchyTool struct{}

func (t jiraGetProjectEpicHierarchyTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetProjectEpicHierarchy(ctx, base.Args.String("project_key"), 200)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type jiraGetCrossProjectDependenciesTool struct{}

func (t jiraGetCrossProjectDependenciesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireJira(); err != nil {
		return nil, err
	}
	r, err := base.Jira.GetCrossProjectDependencies(ctx, base.Args.String("project_key"), 200)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

func (jiraGetIssueTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_issue", mcp.WithDescription("Get a Jira issue"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to return")),
		mcp.WithString("expand", mcp.Description("Fields to expand (e.g., 'renderedFields,transitions,changelog')")),
		mcp.WithNumber("comment_limit", mcp.Description("Maximum number of comments to include")),
		mcp.WithString("properties", mcp.Description("A comma-separated list of issue properties to return")),
		mcp.WithBoolean("update_history", mcp.Description("Whether to update the issue view history for the requesting user")),
		mcp.WithString("include", mcp.Description("Comma-separated sections to inline: transitions, watchers, changelog, comments, worklogs, remote_links")),
		mcp.WithBoolean("use_display_names", mcp.Description("When true, custom field keys use human-readable display names")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetProjectIssuesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_project_issues", mcp.WithDescription("Get project issues"),
		mcp.WithString("project_key", mcp.Required()),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination (0-based)")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results (1-50)")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetIssueDatesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_issue_dates", mcp.WithDescription("Get issue dates with cycle/lead time and optional changelog-based status history"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithBoolean("include_status_changes", mcp.DefaultBool(true), mcp.Description("Include status change history with timestamps and durations")),
		mcp.WithBoolean("include_status_summary", mcp.DefaultBool(true), mcp.Description("Include aggregated time spent in each status")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetIssueSLATool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_issue_sla", mcp.WithDescription("Get issue SLA metrics"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithString("metrics", mcp.Description("Comma-separated list of SLA metrics to calculate. Available: cycle_time, lead_time, time_in_status, due_date_compliance, resolution_time, first_response_time. Defaults to 'cycle_time,time_in_status'.")),
		mcp.WithBoolean("working_hours_only", mcp.Description("Calculate using working hours only (excludes weekends). Defaults to false.")),
		mcp.WithBoolean("include_raw_dates", mcp.Description("Include raw date values in the response")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetIssueImagesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_issue_images", mcp.WithDescription("Get issue images as inline ImageContent for LLM vision"), mcp.WithString("issue_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraDownloadAttachmentsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_download_attachments", mcp.WithDescription("Download issue attachments"), mcp.WithString("issue_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetServiceDeskForProjectTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_service_desk_for_project", mcp.WithDescription("Get SD for project"), mcp.WithString("project_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetServiceDeskQueuesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_service_desk_queues", mcp.WithDescription("Get SD queues"),
		mcp.WithString("service_desk_id", mcp.Required()),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Description("Maximum number of results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetQueueIssuesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_queue_issues", mcp.WithDescription("Get queue issues"),
		mcp.WithString("service_desk_id", mcp.Required()),
		mcp.WithString("queue_id", mcp.Required()),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Description("Maximum number of results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetRequestTypesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_request_types", mcp.WithDescription("Get request types"),
		mcp.WithString("service_desk_id", mcp.Required()),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Description("Maximum number of results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetRequestTypeFieldsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_request_type_fields", mcp.WithDescription("Get request type fields"), mcp.WithString("service_desk_id", mcp.Required()), mcp.WithString("request_type_id", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetUserProfileTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_user_profile", mcp.WithDescription("Get a Jira user"), mcp.WithString("user_identifier", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraSearchTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_search", mcp.WithDescription("Search Jira via JQL"),
		mcp.WithString("jql", mcp.Required()),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to return")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results")),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithString("expand", mcp.Description("Fields to expand (e.g., 'renderedFields,transitions')")),
		mcp.WithString("projects_filter", mcp.Description("Comma-separated list of project keys to filter by")),
		mcp.WithBoolean("use_display_names", mcp.Description("When true, custom field keys use human-readable display names")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetTransitionsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_transitions", mcp.WithDescription("Get issue transitions"), mcp.WithString("issue_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetWorklogTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_worklog", mcp.WithDescription("Get issue worklog"), mcp.WithString("issue_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetIssueWatchersTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_issue_watchers", mcp.WithDescription("Get issue watchers"), mcp.WithString("issue_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetLinkTypesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_link_types", mcp.WithDescription("Get link types"),
		mcp.WithString("name_filter", mcp.Description("Filter link types by name (fuzzy match)")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetAllProjectsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_all_projects", mcp.WithDescription("Get all projects"),
		mcp.WithBoolean("include_archived", mcp.Description("Whether to include archived projects")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraSearchProjectsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_search_projects", mcp.WithDescription("Search projects"),
		mcp.WithString("query", mcp.Required()),
		mcp.WithNumber("max_results", mcp.DefaultNumber(20), mcp.Description("Maximum number of results")),
		mcp.WithString("current_project_ids", mcp.Description("Comma-separated list of project IDs to exclude from results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetProjectIssueTypesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_project_issue_types", mcp.WithDescription("Get project issue types"), mcp.WithString("project_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetCreateFieldsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_create_fields", mcp.WithDescription("Get create fields"), mcp.WithString("project_key", mcp.Required()), mcp.WithString("issue_type_id", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetProjectVersionsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_project_versions", mcp.WithDescription("Get project versions"), mcp.WithString("project_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetProjectComponentsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_project_components", mcp.WithDescription("Get project components"), mcp.WithString("project_key", mcp.Required()), mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false))
}

func (jiraGetProjectFieldsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_project_fields", mcp.WithDescription("Get project fields"),
		mcp.WithString("project_key", mcp.Required()),
		mcp.WithString("schema_type", mcp.Description("Filter fields by schema type (e.g., 'option', 'array', 'user')")),
		mcp.WithString("issue_types", mcp.Description("Comma-separated list of issue type IDs to filter by applicability")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetAgileBoardsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_agile_boards", mcp.WithDescription("Get agile boards"),
		mcp.WithString("board_name", mcp.Description("Filter boards by name (fuzzy search)")),
		mcp.WithString("project_key", mcp.Description("Filter boards by project key")),
		mcp.WithString("board_type", mcp.Description("Filter boards by type ('scrum' or 'kanban')")),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination (0-based)")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results (1-50)")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetBoardIssuesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_board_issues", mcp.WithDescription("Get board issues"),
		mcp.WithString("board_id", mcp.Required()),
		mcp.WithString("jql", mcp.Required()),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to return")),
		mcp.WithString("expand", mcp.Description("Fields to expand (e.g., 'renderedFields,transitions')")),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetSprintsFromBoardTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_sprints_from_board", mcp.WithDescription("Get sprints from board"),
		mcp.WithString("board_id", mcp.Required()),
		mcp.WithString("state", mcp.Description("Filter sprints by state ('active', 'future', 'closed')")),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetSprintIssuesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_sprint_issues", mcp.WithDescription("Get sprint issues"),
		mcp.WithString("sprint_id", mcp.Required()),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to return")),
		mcp.WithNumber("start_at", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraSearchFieldsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_search_fields", mcp.WithDescription("Search fields by keyword with fuzzy match"),
		mcp.WithString("keyword", mcp.Description("Keyword for fuzzy search. If empty, lists the first 'limit' available fields.")),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results")),
		mcp.WithBoolean("refresh", mcp.Description("Whether to force refresh the field list")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetFieldOptionsTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_field_options", mcp.WithDescription("Get allowed options for a custom field (select, multi-select, radio, checkbox). Requires project_key and issue_type for DC."),
		mcp.WithString("field_id", mcp.Required()),
		mcp.WithString("project_key", mcp.Description("(Optional for Cloud) Project key for DC context")),
		mcp.WithString("issue_type", mcp.Description("(Optional for Cloud) Issue type name for DC context")),
		mcp.WithString("contains", mcp.Description("Case-insensitive substring filter on option values")),
		mcp.WithNumber("return_limit", mcp.Description("Maximum number of results to return (after filtering)")),
		mcp.WithBoolean("values_only", mcp.Description("If true, return only value strings in compact format")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraSearchAssignableUsersTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_search_assignable_users", mcp.WithDescription("Search assignable users in a project or issue"),
		mcp.WithString("query", mcp.Required()),
		mcp.WithString("project_key"),
		mcp.WithString("issue_key"),
		mcp.WithInteger("limit"),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetIssueDevelopmentInfoTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_issue_development_info", mcp.WithDescription("Get development info (PRs, branches, commits) for a Jira issue via dev-status plugin"),
		mcp.WithString("issue_key", mcp.Required()),
		mcp.WithString("application_type", mcp.Description("Filter by application type (e.g., 'github', 'bitbucket', 'gitlab')")),
		mcp.WithString("data_type", mcp.Description("Filter by data type (e.g., 'repository', 'branch', 'pullrequest')")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetIssuesDevelopmentInfoTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_issues_development_info", mcp.WithDescription("Get development info for multiple Jira issues"),
		mcp.WithString("issue_keys", mcp.Required()),
		mcp.WithString("application_type", mcp.Description("Filter by application type (e.g., 'github', 'bitbucket', 'gitlab')")),
		mcp.WithString("data_type", mcp.Description("Filter by data type (e.g., 'repository', 'branch', 'pullrequest')")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetProjectEpicHierarchyTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_project_epic_hierarchy", mcp.WithDescription("Group a project's epics under their cross-project parent issues"),
		mcp.WithString("project_key", mcp.Required()),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (jiraGetCrossProjectDependenciesTool) McpTool() mcp.Tool {
	return mcp.NewTool("jira_get_cross_project_dependencies", mcp.WithDescription("Find all cross-project issue links for a project"),
		mcp.WithString("project_key", mcp.Required()),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func ReadTools() []core.Runnable {
	return []core.Runnable{
		&jiraGetIssueTool{},
		&jiraGetProjectIssuesTool{},
		&jiraGetIssueDatesTool{},
		&jiraGetIssueSLATool{},
		&jiraGetIssueImagesTool{},
		&jiraDownloadAttachmentsTool{},
		&jiraGetServiceDeskForProjectTool{},
		&jiraGetServiceDeskQueuesTool{},
		&jiraGetQueueIssuesTool{},
		&jiraGetRequestTypesTool{},
		&jiraGetRequestTypeFieldsTool{},
		&jiraGetUserProfileTool{},
		&jiraSearchTool{},
		&jiraGetTransitionsTool{},
		&jiraGetWorklogTool{},
		&jiraGetIssueWatchersTool{},
		&jiraGetLinkTypesTool{},
		&jiraGetAllProjectsTool{},
		&jiraSearchProjectsTool{},
		&jiraGetProjectIssueTypesTool{},
		&jiraGetCreateFieldsTool{},
		&jiraGetProjectVersionsTool{},
		&jiraGetProjectComponentsTool{},
		&jiraGetProjectFieldsTool{},
		&jiraGetAgileBoardsTool{},
		&jiraGetBoardIssuesTool{},
		&jiraGetSprintsFromBoardTool{},
		&jiraGetSprintIssuesTool{},
		&jiraSearchFieldsTool{},
		&jiraGetFieldOptionsTool{},
		&jiraSearchAssignableUsersTool{},
		&jiraGetIssueDevelopmentInfoTool{},
		&jiraGetIssuesDevelopmentInfoTool{},
		&jiraGetProjectEpicHierarchyTool{},
		&jiraGetCrossProjectDependenciesTool{},
	}
}
