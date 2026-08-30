package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
	srv "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server"
)

func setupBackends(t *testing.T) (*httptest.Server, *httptest.Server) {
	t.Helper()
	j := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/rest/api/2/issue/PROJ-1":
			json.NewEncoder(w).Encode(map[string]any{"id": "10001", "key": "PROJ-1", "fields": map[string]any{"summary": "test"}})
		case "/rest/api/2/user":
			json.NewEncoder(w).Encode(map[string]any{"accountId": "u1", "displayName": "Alice", "active": true})
		case "/rest/api/2/myself":
			json.NewEncoder(w).Encode(map[string]any{"accountId": "me"})
		default:
			w.WriteHeader(404)
		}
	}))
	c := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/rest/api/content/12345":
			json.NewEncoder(w).Encode(map[string]any{"id": "12345", "type": "page", "title": "Test Page", "status": "current"})
		case "/rest/api/content/search":
			json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{{"id": "1", "title": "R1", "type": "page"}}, "size": 1})
		default:
			w.WriteHeader(404)
		}
	}))
	t.Setenv("JIRA_URL", j.URL)
	t.Setenv("CONFLUENCE_URL", c.URL)
	t.Setenv("JIRA_PERSONAL_TOKEN", "test-pat")
	t.Setenv("CONFLUENCE_PERSONAL_TOKEN", "test-pat")
	return j, c
}

func TestMCPServer_ListTools(t *testing.T) {
	j, c := setupBackends(t)
	defer j.Close()
	defer c.Close()
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
	cli.Initialize(ctx, mcp.InitializeRequest{})
	tools, _ := cli.ListTools(ctx, mcp.ListToolsRequest{})
	names := map[string]bool{}
	for _, tl := range tools.Tools {
		names[tl.Name] = true
	}
	for _, n := range []string{"jira_get_issue", "jira_get_user_profile", "jira_search", "confluence_get_page", "confluence_search"} {
		if !names[n] {
			t.Errorf("missing: %s", n)
		}
	}
	t.Logf("tools listed: %d ✅", len(tools.Tools))
}

func TestMCPServer_JiraGetIssue(t *testing.T) {
	j, c := setupBackends(t)
	defer j.Close()
	defer c.Close()
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
	cli.Initialize(ctx, mcp.InitializeRequest{})
	req := mcp.CallToolRequest{}
	req.Params.Name = "jira_get_issue"
	req.Params.Arguments = map[string]any{"issue_key": "PROJ-1"}
	result, _ := cli.CallTool(ctx, req)
	if result.IsError {
		t.Fatal(result.Content)
	}
	var issue map[string]any
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &issue)
	if issue["key"] != "PROJ-1" {
		t.Errorf("key=%v", issue["key"])
	}
	t.Logf("jira_get_issue ✅")
}

func TestMCPServer_ConfluenceGetPage(t *testing.T) {
	j, c := setupBackends(t)
	defer j.Close()
	defer c.Close()
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
	cli.Initialize(ctx, mcp.InitializeRequest{})
	req := mcp.CallToolRequest{}
	req.Params.Name = "confluence_get_page"
	req.Params.Arguments = map[string]any{"page_id": "12345"}
	result, _ := cli.CallTool(ctx, req)
	if result.IsError {
		t.Fatal(result.Content)
	}
	var page map[string]any
	json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &page)
	metadata, ok := page["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'metadata' key in response, got: %v", page)
	}
	if metadata["title"] != "Test Page" {
		t.Errorf("title=%v", metadata["title"])
	}
	t.Logf("confluence_get_page ✅")
}

func TestMCPServer_ReadOnlyBlocksWrite(t *testing.T) {
	_, c := setupBackends(t)
	defer c.Close()
	t.Setenv("READ_ONLY_MODE", "true")
	app := srv.NewAppState(config.JiraConfigFromEnv(), config.ConfluenceConfigFromEnv(), "http")
	if !app.ReadOnly {
		t.Fatal("ReadOnly should be true")
	}
	ms := server.NewMCPServer("Test", "1.0.0", server.WithToolCapabilities(true))
	srv.RegisterAllTools(ms, app)
	ts := server.NewTestServer(ms)
	defer ts.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cli, _ := client.NewSSEMCPClient(ts.URL + "/sse")
	defer cli.Close()
	cli.Start(ctx)
	cli.Initialize(ctx, mcp.InitializeRequest{})
	req := mcp.CallToolRequest{}
	req.Params.Name = "jira_create_issue"
	req.Params.Arguments = map[string]any{"project_key": "PROJ", "summary": "test", "issue_type": "Bug"}
	_, err := cli.CallTool(ctx, req)
	if err != nil {
		t.Logf("write blocked ✅: %v", err)
	} else {
		t.Error("should be blocked")
	}
}

func TestMCPServer_ConfluenceSearch(t *testing.T) {
	j, c := setupBackends(t)
	defer j.Close()
	defer c.Close()
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
	cli.Initialize(ctx, mcp.InitializeRequest{})
	req := mcp.CallToolRequest{}
	req.Params.Name = "confluence_search"
	req.Params.Arguments = map[string]any{"query": "test", "limit": 10}
	result, _ := cli.CallTool(ctx, req)
	if result.IsError {
		t.Fatal(result.Content)
	}
	t.Logf("confluence_search ✅")
}
