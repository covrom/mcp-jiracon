package server

import (
	mcpserver "github.com/mark3labs/mcp-go/server"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
	conf "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server/confluence"
	core "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server/core"
	jira "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server/jira"
)

// AppState is the per-server application state (alias for core.AppState).
type AppState = core.AppState

// Runnable is the contract implemented by every tool struct.
type Runnable = core.Runnable

// NewAppState builds the application state from configuration.
func NewAppState(jiraCfg *config.JiraConfig, confluenceCfg *config.ConfluenceConfig, transport string) *AppState {
	return core.NewAppState(jiraCfg, confluenceCfg, transport)
}

// RegisterAllTools registers every Jira and Confluence tool on the server.
// Each tool's MCP schema is obtained via McpTool(), and its handler is built
// by the generic JiraHandler / ConfluenceHandler, which resolve the fetchers
// from app state and the per-request context at call time.
func RegisterAllTools(s *mcpserver.MCPServer, app *AppState) {
	for _, t := range jira.ReadTools() {
		s.AddTool(t.McpTool(), core.JiraHandler(t, app))
	}
	for _, t := range jira.WriteTools() {
		s.AddTool(t.McpTool(), core.JiraHandler(t, app))
	}
	for _, t := range conf.ReadTools() {
		s.AddTool(t.McpTool(), core.ConfluenceHandler(t, app))
	}
	for _, t := range conf.WriteTools() {
		s.AddTool(t.McpTool(), core.ConfluenceHandler(t, app))
	}
}
