package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/atlassian"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
	core "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server/core"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/urls"
)

func AuthMiddleware(app *AppState) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Per-request header PAT (takes precedence).
			jiraURL := r.Header.Get("X-Atlassian-Jira-Url")
			jiraToken := r.Header.Get("X-Atlassian-Jira-Personal-Token")
			confluenceURL := r.Header.Get("X-Atlassian-Confluence-Url")
			confluenceToken := r.Header.Get("X-Atlassian-Confluence-Personal-Token")

			// Also accept Authorization: Bearer <token> as PAT.
			if jiraToken == "" && confluenceToken == "" {
				auth := r.Header.Get("Authorization")
				if strings.HasPrefix(auth, "Bearer ") {
					tok := strings.TrimPrefix(auth, "Bearer ")
					if jiraURL != "" {
						jiraToken = tok
					}
					if confluenceURL != "" {
						confluenceToken = tok
					}
				}
			}

			if jiraURL != "" {
				if err := urls.ValidateURLForSSRF(jiraURL); err != nil {
					writeError(w, 403, "Invalid Jira URL - "+err.Error())
					return
				}
			}
			if confluenceURL != "" {
				if err := urls.ValidateURLForSSRF(confluenceURL); err != nil {
					writeError(w, 403, "Invalid Confluence URL - "+err.Error())
					return
				}
			}

			// Build per-request PAT clients.
			if jiraURL != "" && jiraToken != "" {
				cfg := &config.JiraConfig{URL: jiraURL, PersonalToken: jiraToken}
				if app.JiraConfig != nil {
					cfg.SSLVerify, cfg.HTTPProxy, cfg.HTTPSProxy, cfg.NoProxy = app.JiraConfig.SSLVerify, app.JiraConfig.HTTPProxy, app.JiraConfig.HTTPSProxy, app.JiraConfig.NoProxy
				}
				if c, err := atlassian.New(cfg); err == nil {
					ctx = context.WithValue(ctx, core.CtxJiraClient, c)
				}
			}
			if confluenceURL != "" && confluenceToken != "" {
				cfg := &config.ConfluenceConfig{URL: confluenceURL, PersonalToken: confluenceToken}
				if app.ConfluenceConfig != nil {
					cfg.SSLVerify, cfg.HTTPProxy, cfg.HTTPSProxy, cfg.NoProxy = app.ConfluenceConfig.SSLVerify, app.ConfluenceConfig.HTTPProxy, app.ConfluenceConfig.HTTPSProxy, app.ConfluenceConfig.NoProxy
				}
				if c, err := atlassian.NewConfluence(cfg); err == nil {
					ctx = context.WithValue(ctx, core.CtxConfluenceClient, c)
				}
			}

			// Fallback: global PAT from env.
			jiraCtx, _ := core.JiraClientFromContext(ctx)
			confCtx, _ := core.ConfluenceClientFromContext(ctx)
			hasJira := jiraCtx != nil || app.GlobalJiraClient != nil
			hasConf := confCtx != nil || app.GlobalConfluenceClient != nil
			hasCreds := jiraToken != "" || confluenceToken != ""
			if !hasCreds && !hasJira && !hasConf && app.Transport == "http" {
				writeError(w, 401, "Authentication required: provide Bearer token or X-Atlassian-*-Personal-Token headers")
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body, _ := json.Marshal(map[string]string{"error": message})
	w.Write(body)
}
