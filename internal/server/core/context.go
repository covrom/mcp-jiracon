package core

import (
	"context"

	"log/slog"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/atlassian"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
)

type AppState struct {
	JiraConfig             *config.JiraConfig
	ConfluenceConfig       *config.ConfluenceConfig
	GlobalJiraClient       *atlassian.Client
	GlobalConfluenceClient *atlassian.Client
	ReadOnly               bool
	EnabledTools           []string
	EnabledToolsets        map[string]struct{}
	Transport              string
}

func NewAppState(jiraCfg *config.JiraConfig, confluenceCfg *config.ConfluenceConfig, transport string) *AppState {
	app := &AppState{
		JiraConfig: jiraCfg, ConfluenceConfig: confluenceCfg,
		EnabledToolsets: map[string]struct{}{
			"jira_issues": {}, "jira_fields": {}, "jira_comments": {}, "jira_transitions": {},
			"confluence_pages": {}, "confluence_comments": {},
		},
		Transport: transport,
	}
	if config.IsEnvExtendedTruthy("READ_ONLY_MODE", "false") {
		app.ReadOnly = true
	}
	if v := envEnabled("ENABLED_TOOLS"); len(v) > 0 {
		app.EnabledTools = v
	}
	if jiraCfg != nil && jiraCfg.IsAuthConfigured() {
		if c, err := atlassian.New(jiraCfg); err == nil {
			app.GlobalJiraClient = c
		}
	}
	if confluenceCfg != nil && confluenceCfg.IsAuthConfigured() {
		if c, err := atlassian.NewConfluence(confluenceCfg); err == nil {
			app.GlobalConfluenceClient = c
		}
	}
	slog.Info("app state initialized", "read_only", app.ReadOnly,
		"jira_available", app.GlobalJiraClient != nil, "confluence_available", app.GlobalConfluenceClient != nil)
	return app
}

func envEnabled(name string) []string {
	if l := config.GetHeaderNames(name); len(l) > 0 {
		return l
	}
	return nil
}

type ContextKey int

const (
	CtxJiraClient ContextKey = iota + 1
	CtxConfluenceClient
	CtxAppState
)

func ctxJiraClientVal(ctx context.Context) *atlassian.Client {
	v, _ := ctx.Value(CtxJiraClient).(*atlassian.Client)
	return v
}
func ctxConfClientVal(ctx context.Context) *atlassian.Client {
	v, _ := ctx.Value(CtxConfluenceClient).(*atlassian.Client)
	return v
}
func JiraClientFromContext(ctx context.Context) (*atlassian.Client, bool) {
	v, ok := ctx.Value(CtxJiraClient).(*atlassian.Client)
	return v, ok
}
func ConfluenceClientFromContext(ctx context.Context) (*atlassian.Client, bool) {
	v, ok := ctx.Value(CtxConfluenceClient).(*atlassian.Client)
	return v, ok
}

// WithAppState stores the app state in the context.
func WithAppState(ctx context.Context, app *AppState) context.Context {
	return context.WithValue(ctx, CtxAppState, app)
}

// AppStateFromContext retrieves the app state from the context.
func AppStateFromContext(ctx context.Context) *AppState {
	v, _ := ctx.Value(CtxAppState).(*AppState)
	return v
}
