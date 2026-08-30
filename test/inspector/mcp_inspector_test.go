// Package inspector contains black-box integration tests that drive the
// compiled mcp-atlassian binary through the official MCP Inspector CLI
// (npx @modelcontextprotocol/inspector), mirroring how a real MCP host
// connects to the server over stdio.
//
// These tests require Node.js + npx. They are skipped automatically when
// npx is unavailable. The first run downloads the Inspector package, so
// network access is required once.
package inspector

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// requireInspector makes the suite fail (instead of skip) when npx is
// unavailable. Opt in from CI via:
//
//	go test ./test/inspector -inspector-required
var requireInspector = flag.Bool("inspector-required", false, "fail instead of skip when npx is unavailable")

const (
	serverName = "Atlassian MCP"
	wantTools  = 99
)

// binPath is the freshly built mcp-atlassian binary under a temp dir, set in TestMain.
var binPath string

// readTools are the GET-only tools that must be annotated as read-only (and not
// destructive) so MCP hosts do not prompt for confirmation on harmless reads.
var readTools = map[string]bool{
	// Jira read
	"jira_get_issue":                      true,
	"jira_get_project_issues":             true,
	"jira_get_issue_dates":                true,
	"jira_get_issue_sla":                  true,
	"jira_get_issue_images":               true,
	"jira_download_attachments":           true,
	"jira_get_service_desk_for_project":   true,
	"jira_get_service_desk_queues":        true,
	"jira_get_queue_issues":               true,
	"jira_get_request_types":              true,
	"jira_get_request_type_fields":        true,
	"jira_get_user_profile":               true,
	"jira_search":                         true,
	"jira_get_transitions":                true,
	"jira_get_worklog":                    true,
	"jira_get_issue_watchers":             true,
	"jira_get_link_types":                 true,
	"jira_get_all_projects":               true,
	"jira_search_projects":                true,
	"jira_get_project_issue_types":        true,
	"jira_get_create_fields":              true,
	"jira_get_project_versions":           true,
	"jira_get_project_components":         true,
	"jira_get_project_fields":             true,
	"jira_get_agile_boards":               true,
	"jira_get_board_issues":               true,
	"jira_get_sprints_from_board":         true,
	"jira_get_sprint_issues":              true,
	"jira_search_fields":                  true,
	"jira_get_field_options":              true,
	"jira_search_assignable_users":        true,
	"jira_get_issue_development_info":     true,
	"jira_get_issues_development_info":    true,
	"jira_get_project_epic_hierarchy":     true,
	"jira_get_cross_project_dependencies": true,
	// Confluence read
	"confluence_get_page":                     true,
	"confluence_search":                       true,
	"confluence_get_page_children":            true,
	"confluence_get_space_page_tree":          true,
	"confluence_get_comments":                 true,
	"confluence_get_labels":                   true,
	"confluence_search_user":                  true,
	"confluence_get_attachments":              true,
	"confluence_get_inline_comments":          true,
	"confluence_get_page_history":             true,
	"confluence_get_page_diff":                true,
	"confluence_download_attachment":          true,
	"confluence_download_content_attachments": true,
	"confluence_get_page_images":              true,
	"confluence_get_page_restrictions":        true,
}

func TestMain(m *testing.M) {
	// TestMain runs before flag.Parse, so parse explicitly to read custom flags.
	flag.Parse()

	if _, err := exec.LookPath("npx"); err != nil {
		if *requireInspector {
			fmt.Fprintln(os.Stderr, "npx not found; MCP Inspector tests are required (-inspector-required)")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "npx not found; skipping MCP Inspector tests")
		os.Exit(0)
	}

	tmp, err := os.MkdirTemp("", "mcp-atlassian-inspector-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "temp dir:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmp)

	binPath = filepath.Join(tmp, "mcp-atlassian")
	build := exec.Command("go", "build", "-o", binPath, "./cmd/mcp-jiracon")
	build.Dir = moduleRoot()
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build mcp-atlassian: %v\n%s\n", err, out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// moduleRoot returns the repository root (where go.mod lives), derived from
// the location of this file rather than the process working directory.
func moduleRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// inspectorResponse is the top-level shape of the Inspector CLI output when
// run with --format json: either a "result" on success or an "error".
type inspectorResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    json.RawMessage `json:"code"`
		Message string          `json:"message"`
	} `json:"error"`
}

// runInspector launches the Inspector CLI against the compiled binary with the
// given method and server environment, and returns the raw stdout plus stderr.
func runInspector(t *testing.T, method string, env map[string]string, extra ...string) ([]byte, string) {
	t.Helper()

	// The target (command) must come immediately after --cli; options follow.
	args := []string{"--yes", "@modelcontextprotocol/inspector", "--cli", binPath, "--method", method, "--format", "json"}
	for k, v := range env {
		args = append(args, "-e", k+"="+v)
	}
	args = append(args, extra...)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "npx", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("inspector %s failed: %v\nstderr:\n%s\nstdout:\n%s", method, err, stderr.String(), stdout.String())
	}
	return stdout.Bytes(), stderr.String()
}

// parseResponse decodes the Inspector's JSON stdout into an inspectorResponse.
func parseResponse(t *testing.T, stdout []byte) inspectorResponse {
	t.Helper()
	var resp inspectorResponse
	if err := json.Unmarshal(stdout, &resp); err == nil {
		return resp
	}
	// npx may interleave its own notices; fall back to the first JSON object line.
	for _, line := range strings.Split(string(stdout), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line[0] != '{' {
			continue
		}
		if err := json.Unmarshal([]byte(line), &resp); err == nil {
			return resp
		}
	}
	t.Fatalf("no JSON object in inspector stdout:\n%s", stdout)
	return resp
}

func mustResult(t *testing.T, resp inspectorResponse) json.RawMessage {
	t.Helper()
	if resp.Error != nil {
		t.Fatalf("inspector returned error: %s", resp.Error.Message)
	}
	if len(resp.Result) == 0 {
		t.Fatalf("inspector returned empty result")
	}
	return resp.Result
}

func TestInspectorInitialize(t *testing.T) {
	stdout, _ := runInspector(t, "initialize", nil)
	resp := parseResponse(t, stdout)
	var result struct {
		ServerInfo struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"serverInfo"`
	}
	if err := json.Unmarshal(mustResult(t, resp), &result); err != nil {
		t.Fatalf("parse initialize result: %v", err)
	}
	if result.ServerInfo.Name != serverName {
		t.Errorf("server name = %q, want %q", result.ServerInfo.Name, serverName)
	}
	if result.ServerInfo.Version == "" {
		t.Errorf("server version is empty")
	}
}

func TestInspectorListTools(t *testing.T) {
	stdout, _ := runInspector(t, "tools/list", nil)
	resp := parseResponse(t, stdout)
	var result struct {
		Tools []struct {
			Name        string `json:"name"`
			Annotations struct {
				ReadOnlyHint    *bool `json:"readOnlyHint"`
				DestructiveHint *bool `json:"destructiveHint"`
			} `json:"annotations"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(mustResult(t, resp), &result); err != nil {
		t.Fatalf("parse tools/list result: %v", err)
	}
	if len(result.Tools) != wantTools {
		t.Errorf("tool count = %d, want %d", len(result.Tools), wantTools)
	}

	names := map[string]bool{}
	for _, tl := range result.Tools {
		names[tl.Name] = true
	}
	for _, n := range []string{"jira_get_issue", "jira_search", "confluence_get_page", "confluence_search", "jira_create_issue", "confluence_create_page"} {
		if !names[n] {
			t.Errorf("missing tool %q", n)
		}
	}
}

// TestInspectorReadToolsAnnotations asserts that read-only tools are correctly
// annotated so MCP hosts do not prompt for confirmation on harmless reads.
func TestInspectorReadToolsAnnotations(t *testing.T) {
	stdout, _ := runInspector(t, "tools/list", nil)
	resp := parseResponse(t, stdout)
	var result struct {
		Tools []struct {
			Name        string `json:"name"`
			Annotations struct {
				ReadOnlyHint    *bool `json:"readOnlyHint"`
				DestructiveHint *bool `json:"destructiveHint"`
			} `json:"annotations"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(mustResult(t, resp), &result); err != nil {
		t.Fatalf("parse tools/list result: %v", err)
	}

	for _, tl := range result.Tools {
		if !readTools[tl.Name] {
			continue
		}
		if tl.Annotations.ReadOnlyHint == nil || !*tl.Annotations.ReadOnlyHint {
			t.Errorf("%s: readOnlyHint = %v, want true", tl.Name, tl.Annotations.ReadOnlyHint)
		}
		if tl.Annotations.DestructiveHint == nil || *tl.Annotations.DestructiveHint {
			t.Errorf("%s: destructiveHint = %v, want false", tl.Name, tl.Annotations.DestructiveHint)
		}
	}
}

// TestInspectorCallReadTool verifies a real tool call round-trip against a mock
// Jira backend, exercising the full stdio path (initialize -> tools/call).
func TestInspectorCallReadTool(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/rest/api/2/issue/PROJ-1":
			json.NewEncoder(w).Encode(map[string]any{"id": "10001", "key": "PROJ-1", "fields": map[string]any{"summary": "Inspector Mock Issue"}})
		default:
			json.NewEncoder(w).Encode(map[string]any{})
		}
	}))
	defer backend.Close()

	env := map[string]string{
		"JIRA_URL":                  backend.URL,
		"JIRA_PERSONAL_TOKEN":       "test-pat",
		"CONFLUENCE_URL":            backend.URL,
		"CONFLUENCE_PERSONAL_TOKEN": "test-pat",
		"ATLASSIAN_RATE_LIMIT_RPS":  "0",
		"ATLASSIAN_MAX_RETRIES":     "0",
	}

	stdout, _ := runInspector(t, "tools/call", env, "--tool-name", "jira_get_issue", "--tool-args-json", `{"issue_key":"PROJ-1"}`)
	resp := parseResponse(t, stdout)
	var result struct {
		IsError bool `json:"isError"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(mustResult(t, resp), &result); err != nil {
		t.Fatalf("parse tools/call result: %v", err)
	}
	if result.IsError {
		t.Fatalf("tool returned isError: %+v", result.Content)
	}
	if len(result.Content) == 0 || result.Content[0].Text == "" {
		t.Fatalf("empty tool result content")
	}
	if !strings.Contains(result.Content[0].Text, "PROJ-1") {
		t.Errorf("tool result does not contain mock data: %s", result.Content[0].Text)
	}
}
