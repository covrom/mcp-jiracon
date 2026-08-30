package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/confluence"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/jira"
)

// jiraFetcher / confFetcher are aliases so tool files read naturally without
// importing the internal packages in every file.
type (
	jiraFetcher = jira.Fetcher
	confFetcher = confluence.Fetcher
)

// Tool is the base context embedded in every concrete tool struct. It carries
// the app state, the typed request arguments, and the Jira/Confluence fetchers
// resolved from the request context and app configuration.
type Tool struct {
	App  *AppState
	Args Args

	Jira       *jiraFetcher
	Confluence *confFetcher
}

// Runnable is the contract implemented by every tool struct: McpTool provides
// the MCP schema, Run performs the tool's work and returns the MCP result.
type Runnable interface {
	McpTool() mcp.Tool
	Run(ctx context.Context, base *Tool) (*mcp.CallToolResult, error)
}

// JiraHandler returns an MCP handler func for a Jira tool.
func JiraHandler[T Runnable](v T, app *AppState) mcpserver.ToolHandlerFunc {
	return handlerFor(v, app, func(app *AppState, ctx context.Context) *Tool {
		base := &Tool{App: app}
		if c, ok := JiraClientFromContext(ctx); ok && app.JiraConfig != nil {
			base.Jira = jira.New(app.JiraConfig, c)
		} else if app.GlobalJiraClient != nil && app.JiraConfig != nil {
			base.Jira = jira.New(app.JiraConfig, app.GlobalJiraClient)
		}
		return base
	})
}

// ConfluenceHandler returns an MCP handler func for a Confluence tool.
func ConfluenceHandler[T Runnable](v T, app *AppState) mcpserver.ToolHandlerFunc {
	return handlerFor(v, app, func(app *AppState, ctx context.Context) *Tool {
		base := &Tool{App: app}
		if c, ok := ConfluenceClientFromContext(ctx); ok && app.ConfluenceConfig != nil {
			base.Confluence = confluence.New(app.ConfluenceConfig, c)
		} else if app.GlobalConfluenceClient != nil && app.ConfluenceConfig != nil {
			base.Confluence = confluence.New(app.ConfluenceConfig, app.GlobalConfluenceClient)
		}
		return base
	})
}

// handlerFor is the shared implementation for JiraHandler and ConfluenceHandler.
func handlerFor[T Runnable](v T, app *AppState, makeBase func(app *AppState, ctx context.Context) *Tool) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		base := makeBase(app, ctx)
		base.Args = newArgs(args)
		return v.Run(ctx, base)
	}
}

var errNoJiraClient = fmt.Errorf("no Jira client available")
var errNoConfluenceClient = fmt.Errorf("no Confluence client available")

// requireJira returns an error when no Jira client is available.
func (t *Tool) RequireJira() error {
	if t.Jira == nil {
		return errNoJiraClient
	}
	return nil
}

// requireConfluence returns an error when no Confluence client is available.
func (t *Tool) RequireConfluence() error {
	if t.Confluence == nil {
		return errNoConfluenceClient
	}
	return nil
}

// guardWrite rejects write operations when the server runs in read-only mode.
func (t *Tool) GuardWrite() error {
	if t.App.ReadOnly {
		return fmt.Errorf("cannot perform write in read-only mode")
	}
	return nil
}

// Args provides typed access to the raw tool arguments.
type Args struct {
	raw map[string]any
}

func newArgs(raw map[string]any) Args { return Args{raw: raw} }

// Has reports whether the key is present in the arguments.
func (a Args) Has(key string) bool {
	_, ok := a.raw[key]
	return ok
}

// String returns the string value of key ("" when absent or not a string).
func (a Args) String(key string) string {
	v, _ := a.raw[key].(string)
	return v
}

// Bool returns the bool value of key (def when absent or not a bool).
func (a Args) Bool(key string, def bool) bool {
	if v, ok := a.raw[key].(bool); ok {
		return v
	}
	return def
}

// Num returns the numeric value of key and whether it was present as a number.
// Use this when an explicit zero is meaningful (distinct from an absent key).
func (a Args) Num(key string) (float64, bool) {
	v, ok := a.raw[key].(float64)
	return v, ok
}

// Float returns the numeric value of key as float64 (def when absent or zero).
func (a Args) Float(key string, def float64) float64 {
	if v, ok := a.raw[key].(float64); ok && v != 0 {
		return v
	}
	return def
}

// Int returns the numeric value of key as int (def when absent or zero).
func (a Args) Int(key string, def int) int {
	return int(a.Float(key, float64(def)))
}

// List returns the []string value of key (empty when absent).
func (a Args) List(key string) []string {
	raw, _ := a.raw[key].([]any)
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// JSON unmarshals the JSON string at key into out; false when absent or invalid.
func (a Args) JSON(key string, out any) bool {
	s, _ := a.raw[key].(string)
	if s == "" {
		return false
	}
	return json.Unmarshal([]byte(s), out) == nil
}

// Stringify renders a value as indented JSON. Exported for tool sub-packages.
func Stringify(v any) string { b, _ := json.MarshalIndent(v, "", "  "); return string(b) }

// ContainsFold reports whether s contains substr (case-insensitive).
func ContainsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// SplitComma splits a comma-separated string into trimmed non-empty parts.
func SplitComma(s string) []string {
	out := []string{}
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
