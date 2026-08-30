package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
	srv "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server"
)

// fullMockJira handles ALL DC Jira REST endpoints used by our tools.
func fullMockJira(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		switch {
		case path == "/rest/api/2/issue/PROJ-1":
			if r.URL.Query().Get("fields") == "attachment" {
				base := "http://" + r.Host + "/"
				json.NewEncoder(w).Encode(map[string]any{
					"id": "10001", "key": "PROJ-1",
					"fields": map[string]any{
						"attachment": []any{
							map[string]any{
								"id": "att1", "filename": "test.png",
								"mimeType": "image/png", "size": float64(1024),
								"content": base + "secure/attachment/att1/test.png",
							},
						},
					},
				})
				return
			}
			json.NewEncoder(w).Encode(map[string]any{"id": "10001", "key": "PROJ-1", "fields": map[string]any{"summary": "test", "status": map[string]any{"id": "1", "name": "Open"}, "created": "2024-01-01", "updated": "2024-01-02"}})
		case path == "/rest/api/2/issue/TEST-1":
			json.NewEncoder(w).Encode(map[string]any{"id": "20001", "key": "TEST-1", "fields": map[string]any{"summary": "test2", "status": map[string]any{"id": "1", "name": "Open"}, "created": "2024-01-01", "updated": "2024-01-02"}})
		case path == "/rest/api/2/issue/PROJ-1/watchers":
			json.NewEncoder(w).Encode(map[string]any{"isWatching": true, "watchCount": 2, "watchers": []any{map[string]any{"accountId": "u1"}, map[string]any{"accountId": "u2"}}})
		case path == "/rest/api/2/issue/PROJ-1/transitions":
			json.NewEncoder(w).Encode(map[string]any{"transitions": []map[string]any{{"id": "11", "name": "Done", "to": map[string]any{"name": "Done"}}}})
		case path == "/rest/api/2/issue/PROJ-1/worklog":
			json.NewEncoder(w).Encode(map[string]any{"worklogs": []map[string]any{{"id": "w1", "timeSpent": "1h", "timeSpentSeconds": 3600, "started": "2024-01-01", "created": "2024-01-01", "updated": "2024-01-01"}}})
		case path == "/rest/api/2/issue/PROJ-1/comment" && r.Method == "GET":
			json.NewEncoder(w).Encode(map[string]any{"comments": []map[string]any{{"id": "c1", "body": "hello", "created": "2024-01-01"}}})
		case path == "/rest/api/2/issueLinkType":
			json.NewEncoder(w).Encode(map[string]any{"issueLinkTypes": []map[string]any{{"id": "1", "name": "Blocks", "inward": "blocked by", "outward": "blocks"}}})
		case path == "/rest/api/2/project/search":
			json.NewEncoder(w).Encode(map[string]any{"values": []map[string]any{{"key": "PROJ", "name": "Project", "id": "1"}}})
		case path == "/rest/api/2/project/picker":
			json.NewEncoder(w).Encode(map[string]any{"projects": []map[string]any{{"key": "PROJ", "name": "Project"}}})
		case path == "/rest/api/2/issue/createmeta/PROJ/issuetypes":
			json.NewEncoder(w).Encode(map[string]any{"values": []map[string]any{{"id": "3", "name": "Task", "subtask": false}}})
		case path == "/rest/api/2/issue/createmeta/PROJ/issuetype/3":
			json.NewEncoder(w).Encode(map[string]any{"fields": []map[string]any{{"fieldId": "summary", "name": "Summary", "required": true}}})
		case path == "/rest/api/2/project/PROJ/versions":
			json.NewEncoder(w).Encode([]map[string]any{{"id": "v1", "name": "1.0"}})
		case path == "/rest/api/2/project/PROJ/components":
			json.NewEncoder(w).Encode([]map[string]any{{"id": "c1", "name": "Frontend"}})
		case path == "/rest/agile/1.0/board":
			json.NewEncoder(w).Encode(map[string]any{"values": []map[string]any{{"id": "1", "name": "Board", "type": "scrum"}}})
		case path == "/rest/agile/1.0/board/1/issue":
			json.NewEncoder(w).Encode(map[string]any{"issues": []map[string]any{{"key": "PROJ-1"}}})
		case path == "/rest/agile/1.0/board/1/sprint":
			json.NewEncoder(w).Encode(map[string]any{"values": []map[string]any{{"id": "s1", "name": "Sprint 1"}}})
		case path == "/rest/agile/1.0/sprint/s1/issue":
			json.NewEncoder(w).Encode(map[string]any{"issues": []map[string]any{{"key": "PROJ-1"}}})
		case path == "/rest/api/2/field":
			json.NewEncoder(w).Encode([]map[string]any{{"id": "summary", "name": "Summary", "custom": false}, {"id": "customfield_10010", "name": "Story Points", "custom": true}})
		case path == "/rest/api/2/field/customfield_10010/context":
			json.NewEncoder(w).Encode(map[string]any{"values": []map[string]any{}})
		case path == "/rest/api/2/user/assignable/search":
			json.NewEncoder(w).Encode([]map[string]any{{"accountId": "u1", "displayName": "Alice", "active": true}})
		case path == "/rest/api/2/user":
			json.NewEncoder(w).Encode(map[string]any{"accountId": "u1", "displayName": "Alice", "active": true})
		case path == "/rest/api/2/search":
			json.NewEncoder(w).Encode(map[string]any{"issues": []map[string]any{{"key": "PROJ-1", "id": "10001", "fields": map[string]any{"summary": "test"}}}, "total": 1})
		case path == "/rest/api/2/version":
			json.NewEncoder(w).Encode(map[string]any{"id": "v1", "name": "1.0"})
		case path == "/rest/api/2/version/v1":
			json.NewEncoder(w).Encode(map[string]any{"id": "v1", "name": "1.0"})
		case path == "/rest/api/2/issue/PROJ-1/remotelink":
			json.NewEncoder(w).Encode(map[string]any{"id": "r1"})
		case path == "/rest/api/2/issueLink/100":
			json.NewEncoder(w).Encode(map[string]any{})
		case path == "/rest/api/2/issueLink":
			json.NewEncoder(w).Encode(map[string]any{})
		case path == "/rest/api/2/issue/PROJ-1/comment" && r.Method == "POST":
			json.NewEncoder(w).Encode(map[string]any{"id": "c2", "body": "test"})
		case path == "/rest/api/2/issue/PROJ-1/comment/c1" && r.Method == "PUT":
			json.NewEncoder(w).Encode(map[string]any{"id": "c1", "body": "updated"})
		case path == "/rest/api/2/issue/PROJ-1/assignee":
			json.NewEncoder(w).Encode(map[string]any{"key": "PROJ-1"})
		case path == "/rest/api/2/issue/PROJ-1/transitions" && r.Method == "POST":
			json.NewEncoder(w).Encode(map[string]any{"key": "PROJ-1"})
		case path == "/rest/api/2/issue/PROJ-1/worklog" && r.Method == "POST":
			json.NewEncoder(w).Encode(map[string]any{"id": "w2", "timeSpent": "2h"})
		case path == "/rest/api/2/issue":
			json.NewEncoder(w).Encode(map[string]any{"id": "10002", "key": "PROJ-2"})
		case path == "/rest/agile/1.0/sprint":
			json.NewEncoder(w).Encode(map[string]any{"id": "s2", "name": "Sprint 2"})
		case path == "/rest/agile/1.0/sprint/s1":
			json.NewEncoder(w).Encode(map[string]any{"id": "s1", "name": "Updated Sprint"})
		case path == "/rest/agile/1.0/sprint/s1/issue":
			json.NewEncoder(w).Encode(map[string]any{})
		case path == "/rest/agile/1.0/backlog/issue":
			json.NewEncoder(w).Encode(map[string]any{})
		// Service Desk
		case path == "/rest/servicedeskapi/servicedesk/by-project/SUP":
			json.NewEncoder(w).Encode(map[string]any{"id": "4", "projectKey": "SUP"})
		case path == "/rest/servicedeskapi/servicedesk/4/queue":
			json.NewEncoder(w).Encode(map[string]any{"values": []map[string]any{{"id": "47", "name": "Open"}}})
		case path == "/rest/servicedeskapi/servicedesk/4/queue/47/issue":
			json.NewEncoder(w).Encode(map[string]any{"values": []map[string]any{{"key": "SUP-1"}}})
		case path == "/rest/servicedeskapi/servicedesk/4/requesttype":
			json.NewEncoder(w).Encode(map[string]any{"values": []map[string]any{{"id": "23", "name": "Get help"}}})
		case path == "/rest/servicedeskapi/servicedesk/4/requesttype/23/field":
			json.NewEncoder(w).Encode(map[string]any{"fields": []map[string]any{}})
		case path == "/rest/servicedeskapi/request":
			json.NewEncoder(w).Encode(map[string]any{"issueKey": "SUP-1"})
		case strings.HasPrefix(path, "/secure/attachment/"):
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write([]byte("attachment-content"))
		default:
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		}
	}))
}

// fullMockConf handles ALL DC Confluence REST endpoints.
func fullMockConf(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		switch {
		case path == "/rest/api/content/12345":
			json.NewEncoder(w).Encode(map[string]any{"id": "12345", "type": "page", "title": "Test Page", "status": "current"})
		case path == "/rest/api/content/search":
			json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{{"id": "1", "type": "page", "title": "R1"}}, "size": 1})
		case path == "/rest/api/content/12345/child/comment":
			json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{
				{"id": "c1", "body": map[string]any{"view": map[string]any{"value": "hello"}}, "created": "2024-01-01"},
				{"id": "c2", "body": map[string]any{"view": map[string]any{"value": "reply"}}, "created": "2024-01-02",
					"ancestors": []any{map[string]any{"id": "12345", "type": "page"}, map[string]any{"id": "c1", "type": "comment"}}},
			}})
		case path == "/rest/api/content/12345/label":
			json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{{"prefix": "global", "name": "test"}}})
		case path == "/rest/api/content/12345/child/page":
			json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{{"id": "c1", "title": "Child"}}})
		case path == "/rest/api/content":
			j := map[string]any{"results": []map[string]any{{"id": "1", "title": "Home", "ancestors": []any{}}}, "size": 1}
			if r.URL.Query().Get("spaceKey") == "DEV" {
				json.NewEncoder(w).Encode(j)
				return
			}
			json.NewEncoder(w).Encode(map[string]any{"id": "new1", "type": "page", "title": "New Page"})
		case path == "/rest/api/content/12345/child/attachment":
			json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{{"id": "att1", "title": "file.pdf"}}, "size": 1})
		case path == "/rest/api/content/att1":
			json.NewEncoder(w).Encode(map[string]any{})
		case path == "/rest/api/group/confluence-users/member":
			json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{{"displayName": "Alice", "username": "alice"}}, "size": 1})
		case path == "/rest/api/space":
			json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{{"key": "DEV", "name": "Development"}}})
		default:
			json.NewEncoder(w).Encode(map[string]any{"ok": true})
		}
	}))
}

func TestAllJiraReadTools(t *testing.T) {
	jiraSrv := fullMockJira(t)
	defer jiraSrv.Close()
	confSrv := fullMockConf(t)
	defer confSrv.Close()

	t.Setenv("JIRA_URL", jiraSrv.URL)
	t.Setenv("JIRA_PERSONAL_TOKEN", "pat")
	t.Setenv("CONFLUENCE_URL", confSrv.URL)
	t.Setenv("CONFLUENCE_PERSONAL_TOKEN", "pat")

	app := srv.NewAppState(config.JiraConfigFromEnv(), config.ConfluenceConfigFromEnv(), "http")
	ms := server.NewMCPServer("Test", "1.0.0", server.WithToolCapabilities(true))
	srv.RegisterAllTools(ms, app)
	ts := server.NewTestServer(ms)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cli, _ := client.NewSSEMCPClient(ts.URL + "/sse")
	defer cli.Close()
	cli.Start(ctx)
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{Name: "test", Version: "1.0.0"}
	cli.Initialize(ctx, initReq)

	tools, _ := cli.ListTools(ctx, mcp.ListToolsRequest{})
	names := map[string]bool{}
	for _, tl := range tools.Tools {
		names[tl.Name] = true
	}
	t.Logf("Registered: %d tools", len(tools.Tools))

	jiraReadTools := map[string]map[string]any{
		"jira_get_issue":                    {"issue_key": "PROJ-1"},
		"jira_get_user_profile":             {"user_identifier": "alice"},
		"jira_search":                       {"jql": "project=PROJ"},
		"jira_get_project_issues":           {"project_key": "PROJ"},
		"jira_get_transitions":              {"issue_key": "PROJ-1"},
		"jira_get_worklog":                  {"issue_key": "PROJ-1"},
		"jira_get_issue_watchers":           {"issue_key": "PROJ-1"},
		"jira_get_link_types":               {},
		"jira_get_all_projects":             {},
		"jira_search_projects":              {"query": "proj"},
		"jira_get_project_issue_types":      {"project_key": "PROJ"},
		"jira_get_create_fields":            {"project_key": "PROJ", "issue_type_id": "3"},
		"jira_get_project_versions":         {"project_key": "PROJ"},
		"jira_get_project_components":       {"project_key": "PROJ"},
		"jira_get_project_fields":           {"project_key": "PROJ"},
		"jira_get_agile_boards":             {},
		"jira_get_board_issues":             {"board_id": "1", "jql": "project=PROJ"},
		"jira_get_sprints_from_board":       {"board_id": "1"},
		"jira_get_sprint_issues":            {"sprint_id": "s1"},
		"jira_search_fields":                {"keyword": "story"},
		"jira_get_field_options":            {"field_id": "customfield_10010", "project_key": "PROJ", "issue_type": "Task"},
		"jira_search_assignable_users":      {"query": "alice", "project_key": "PROJ"},
		"jira_get_issue_dates":              {"issue_key": "PROJ-1"},
		"jira_get_issue_sla":                {"issue_key": "PROJ-1"},
		"jira_get_issue_images":             {"issue_key": "PROJ-1"},
		"jira_download_attachments":         {"issue_key": "PROJ-1"},
		"jira_get_service_desk_for_project": {"project_key": "SUP"},
		"jira_get_service_desk_queues":      {"service_desk_id": "4"},
		"jira_get_queue_issues":             {"service_desk_id": "4", "queue_id": "47"},
		"jira_get_request_types":            {"service_desk_id": "4"},
		"jira_get_request_type_fields":      {"service_desk_id": "4", "request_type_id": "23"},
	}

	for name, args := range jiraReadTools {
		t.Run(name, func(t *testing.T) {
			if !names[name] {
				t.Skipf("tool %s not registered", name)
				return
			}
			req := mcp.CallToolRequest{}
			req.Params.Name = name
			req.Params.Arguments = args
			result, err := cli.CallTool(ctx, req)
			if err != nil {
				t.Fatalf("call failed: %v", err)
			}
			if result.IsError {
				t.Fatalf("tool error: %v", result.Content)
			}
			t.Logf("%s ✅", name)
		})
	}
}

func TestAllConfluenceReadTools(t *testing.T) {
	jiraSrv := fullMockJira(t)
	defer jiraSrv.Close()
	confSrv := fullMockConf(t)
	defer confSrv.Close()

	t.Setenv("JIRA_URL", jiraSrv.URL)
	t.Setenv("JIRA_PERSONAL_TOKEN", "pat")
	t.Setenv("CONFLUENCE_URL", confSrv.URL)
	t.Setenv("CONFLUENCE_PERSONAL_TOKEN", "pat")

	app := srv.NewAppState(config.JiraConfigFromEnv(), config.ConfluenceConfigFromEnv(), "http")
	ms := server.NewMCPServer("Test", "1.0.0", server.WithToolCapabilities(true))
	srv.RegisterAllTools(ms, app)
	ts := server.NewTestServer(ms)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cli, _ := client.NewSSEMCPClient(ts.URL + "/sse")
	defer cli.Close()
	cli.Start(ctx)
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{Name: "test", Version: "1.0.0"}
	cli.Initialize(ctx, initReq)

	confReadTools := map[string]map[string]any{
		"confluence_get_page":            {"page_id": "12345"},
		"confluence_search":              {"query": "test"},
		"confluence_get_page_children":   {"parent_id": "12345"},
		"confluence_get_space_page_tree": {"space_key": "DEV"},
		"confluence_get_comments":        {"page_id": "12345"},
		"confluence_get_labels":          {"page_id": "12345"},
		"confluence_search_user":         {"query": "alice"},
		"confluence_get_attachments":     {"content_id": "12345"},
	}

	for name, args := range confReadTools {
		t.Run(name, func(t *testing.T) {
			req := mcp.CallToolRequest{}
			req.Params.Name = name
			req.Params.Arguments = args
			result, err := cli.CallTool(ctx, req)
			if err != nil {
				t.Fatalf("call failed: %v", err)
			}
			if result.IsError {
				t.Fatalf("tool error: %v", result.Content)
			}
			t.Logf("%s ✅", name)
		})
	}
}

func TestAllJiraWriteTools(t *testing.T) {
	jiraSrv := fullMockJira(t)
	defer jiraSrv.Close()
	confSrv := fullMockConf(t)
	defer confSrv.Close()

	t.Setenv("JIRA_URL", jiraSrv.URL)
	t.Setenv("JIRA_PERSONAL_TOKEN", "pat")
	t.Setenv("CONFLUENCE_URL", confSrv.URL)
	t.Setenv("CONFLUENCE_PERSONAL_TOKEN", "pat")

	app := srv.NewAppState(config.JiraConfigFromEnv(), config.ConfluenceConfigFromEnv(), "http")
	ms := server.NewMCPServer("Test", "1.0.0", server.WithToolCapabilities(true))
	srv.RegisterAllTools(ms, app)
	ts := server.NewTestServer(ms)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cli, _ := client.NewSSEMCPClient(ts.URL + "/sse")
	defer cli.Close()
	cli.Start(ctx)
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{Name: "test", Version: "1.0.0"}
	cli.Initialize(ctx, initReq)

	jiraWriteTools := map[string]map[string]any{
		"jira_create_issue":             {"project_key": "PROJ", "summary": "New issue", "issue_type": "Task"},
		"jira_update_issue":             {"issue_key": "PROJ-1", "fields": `{"summary":"Updated"}`},
		"jira_delete_issue":             {"issue_key": "PROJ-1"},
		"jira_assign_issue":             {"issue_key": "PROJ-1", "assignee": "alice"},
		"jira_transition_issue":         {"issue_key": "PROJ-1", "transition_id": "11"},
		"jira_add_worklog":              {"issue_key": "PROJ-1", "time_spent": "1h"},
		"jira_add_watcher":              {"issue_key": "PROJ-1", "user_identifier": "user1"},
		"jira_remove_watcher":           {"issue_key": "PROJ-1", "user_identifier": "user1"},
		"jira_link_to_epic":             {"issue_key": "PROJ-1", "epic_key": "EPIC-1"},
		"jira_create_issue_link":        {"link_type": "Blocks", "inward_issue_key": "PROJ-1", "outward_issue_key": "PROJ-2"},
		"jira_create_remote_issue_link": {"issue_key": "PROJ-1", "url": "http://x.com", "title": "X"},
		"jira_remove_issue_link":        {"link_id": "100"},
		"jira_add_comment":              {"issue_key": "PROJ-1", "body": "test"},
		"jira_edit_comment":             {"issue_key": "PROJ-1", "comment_id": "c1", "body": "updated"},
		"jira_create_version":           {"project_key": "PROJ", "name": "2.0"},
		"jira_update_version":           {"version_id": "v1"},
		"jira_batch_create_versions":    {"project_key": "PROJ", "versions": `[{"name":"1.0"}]`},
		"jira_batch_create_issues":      {"issues": `[{"project_key":"PROJ","summary":"test","issue_type":"Task"}]`},
		"jira_create_sprint":            {"board_id": "1", "name": "Sprint X", "start_date": "2024-01-01", "end_date": "2024-01-14"},
		"jira_update_sprint":            {"sprint_id": "s1", "name": "Sprint Y"},
		"jira_add_issues_to_sprint":     {"sprint_id": "s1", "issue_keys": "PROJ-1,PROJ-2"},
		"jira_move_issues_to_backlog":   {"issue_keys": "PROJ-3,PROJ-4"},
		"jira_create_customer_request":  {"service_desk_id": "4", "request_type_id": "23", "request_field_values": `{"summary":"Help"}`},
	}

	for name, args := range jiraWriteTools {
		t.Run(name, func(t *testing.T) {
			req := mcp.CallToolRequest{}
			req.Params.Name = name
			req.Params.Arguments = args
			result, err := cli.CallTool(ctx, req)
			if err != nil {
				t.Fatalf("call failed: %v", err)
			}
			if result.IsError {
				t.Fatalf("tool error: %v", result.Content)
			}
			t.Logf("%s ✅", name)
		})
	}
}

func TestAllConfluenceWriteTools(t *testing.T) {
	jiraSrv := fullMockJira(t)
	defer jiraSrv.Close()
	confSrv := fullMockConf(t)
	defer confSrv.Close()

	t.Setenv("JIRA_URL", jiraSrv.URL)
	t.Setenv("JIRA_PERSONAL_TOKEN", "pat")
	t.Setenv("CONFLUENCE_URL", confSrv.URL)
	t.Setenv("CONFLUENCE_PERSONAL_TOKEN", "pat")

	app := srv.NewAppState(config.JiraConfigFromEnv(), config.ConfluenceConfigFromEnv(), "http")
	ms := server.NewMCPServer("Test", "1.0.0", server.WithToolCapabilities(true))
	srv.RegisterAllTools(ms, app)
	ts := server.NewTestServer(ms)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cli, _ := client.NewSSEMCPClient(ts.URL + "/sse")
	defer cli.Close()
	cli.Start(ctx)
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{Name: "test", Version: "1.0.0"}
	cli.Initialize(ctx, initReq)

	confWriteTools := map[string]map[string]any{
		"confluence_create_page":       {"space_key": "DEV", "title": "New", "content": "<p>test</p>"},
		"confluence_update_page":       {"page_id": "12345", "title": "Updated", "content": "<p>updated</p>"},
		"confluence_delete_page":       {"page_id": "12345"},
		"confluence_move_page":         {"page_id": "12345", "target_parent_id": "99999"},
		"confluence_add_comment":       {"page_id": "12345", "body": "test"},
		"confluence_reply_to_comment":  {"comment_id": "c1", "body": "reply"},
		"confluence_add_label":         {"page_id": "12345", "name": "test"},
		"confluence_delete_attachment": {"attachment_id": "att1"},
	}

	for name, args := range confWriteTools {
		t.Run(name, func(t *testing.T) {
			req := mcp.CallToolRequest{}
			req.Params.Name = name
			req.Params.Arguments = args
			result, err := cli.CallTool(ctx, req)
			if err != nil {
				t.Fatalf("call failed: %v", err)
			}
			if result.IsError {
				t.Fatalf("tool error: %v", result.Content)
			}
			t.Logf("%s ✅", name)
		})
	}
}
