package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
	srv "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/utils"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/pkg/version"
)

func main() {
	transport := flag.String("transport", "stdio", "transport: stdio | http")
	addr := flag.String("addr", ":3000", "HTTP listen address")
	flag.Parse()

	logger := utils.SetupLogging(utils.LogLevelFromEnv())
	logger.Info("mcp-atlassian starting", "version", version.Version, "transport", *transport)

	jiraCfg := config.JiraConfigFromEnv()
	confluenceCfg := config.ConfluenceConfigFromEnv()
	app := srv.NewAppState(jiraCfg, confluenceCfg, *transport)

	if err := jiraCfg.Validate(); err != nil {
		logger.Warn("jira config", "err", err)
	}
	if err := confluenceCfg.Validate(); err != nil {
		logger.Warn("confluence config", "err", err)
	}

	ms := mcpserver.NewMCPServer("Atlassian MCP", version.Version, mcpserver.WithToolCapabilities(true))
	srv.RegisterAllTools(ms, app)

	switch *transport {
	case "stdio":
		if err := mcpserver.ServeStdio(ms); err != nil {
			fmt.Fprintln(os.Stderr, "stdio error:", err)
			os.Exit(1)
		}
	case "http":
		sse := mcpserver.NewSSEServer(ms)
		mux := http.NewServeMux()
		mux.Handle("/mcp", srv.AuthMiddleware(app)(sse))
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"ok"}`))
		})
		logger.Info("listening", "addr", *addr)
		if err := http.ListenAndServe(*addr, mux); err != nil {
			fmt.Fprintln(os.Stderr, "http error:", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "unknown transport:", *transport)
		os.Exit(1)
	}
	_ = slog.Default()
}
