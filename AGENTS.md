# AGENTS.md — инструкция для AI-агентов

## О проекте

`mcp-atlassian` — MCP-сервер для **Atlassian Jira и Confluence Data Center**.
Написан на Go 1.26. Экспонирует 85 MCP-инструментов с аутентификацией через
Personal Access Token (Bearer).

## Сборка

```bash
go build -o mcp-atlassian ./cmd/mcp-jiracon
```

## Тесты

```bash
go test ./... -count=1        # все тесты
go test ./test/integration -v  # только интеграционные (MCP SSE + mock)
```

Всего 138 тестов, покрывают все 85 инструментов через MCP SSE client.

## Структура

```
cmd/mcp-jiracon/main.go     — точка входа (stdio / HTTP SSE)
internal/
  server/tools_registry.go     — регистрация всех 85 MCP-инструментов (McpTool + JiraHandler/ConfluenceHandler)
  server/middleware.go         — AuthMiddleware (PAT, per-request headers)
  server/core/tool.go          — Tool, Args, Runnable, JiraHandler, ConfluenceHandler, helpers
  server/core/context.go       — AppState, контекст-ключи клиентов
  server/jira/read.go          — Jira read (34)
  server/jira/write.go         — Jira write (23)
  server/confluence/read.go    — Confluence read (15)
  server/confluence/write.go   — Confluence write (13)
  jira/fetcher.go              — базовые методы Jira API v2
  jira/fetcher_ext.go          — расширенные методы (watchers, transitions, agile, SD)
  confluence/fetcher.go        — базовые методы Confluence API v1
  confluence/fetcher_ext.go    — расширенные методы (comments, labels, attachments)
  atlassian/client.go          — HTTP-клиент с PAT Bearer-аутентификацией
  atlassian/transport.go       — транспорт: SSRF guard, retry, rate-limit, TLS skip
  config/jira.go               — JiraConfig: JIRA_URL, JIRA_PERSONAL_TOKEN
  config/confluence.go         — ConfluenceConfig
  config/env.go                — парсинг env-переменных (IsEnvTruthy, GetIntEnv, ...)
  urls/urls.go                 — ValidateURLForSSRF
  preprocessing/jira.go        — Jira wiki ↔ Markdown
  models/...                   — JiraIssue, JiraUser, ConfluencePage, ...
  utils/...                    — logging, media, pagination, toolsets
test/integration/              — MCP-интеграционные тесты (SSE + httptest)
```

## Ключевые файлы

| Файл | Назначение |
| --- | --- |
| `internal/server/jira/{read,write}.go` | Jira-инструменты (34 + 23). Добавлять новые — здесь |
| `internal/server/confluence/{read,write}.go` | Confluence-инструменты (15 + 13). Добавлять новые — здесь |
| `internal/jira/fetcher_ext.go` | Методы Jira API. Добавлять эндпоинты — здесь |
| `internal/confluence/fetcher_ext.go` | Методы Confluence API |
| `internal/atlassian/client.go` | HTTP-клиент. Если нужно изменить аутентификацию |
| `test/integration/all_tools_test.go` | Тесты всех 85 инструментов через MCP |

## Конвенции

- **Аутентификация**: только PAT. `Authorization: Bearer <JIRA_PERSONAL_TOKEN>`
- **Jira API**: `/rest/api/2/...` (v2, DC)
- **Confluence API**: `{base}/rest/api/...` (v1, DC). Метод `V1BaseURL()` = `BaseURL + "/rest/api"`
- **Обработка ошибок**: инструменты возвращают ошибку через `return nil, err`
- **Read-only**: `base.GuardWrite()` проверяет `app.ReadOnly` в каждом write-инструменте
- **Тесты с моками**: `httptest.NewServer` + `client.NewSSEMCPClient`

### Архитектура инструментов

Каждый инструмент — пустая структура (`struct{}`) без эмбеддинга. Контекст
(`*core.Tool`) передаётся явным параметром в `Run`:

```go
type jiraGetIssueTool struct{}

func (jiraGetIssueTool) McpTool() mcp.Tool { ... }

func (jiraGetIssueTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
    if err := base.RequireJira(); err != nil { return nil, err }
    a := base.Args
    issue, err := base.Jira.GetIssueFull(ctx, a.String("issue_key"), ...)
    ...
}
```

Ресиверы без имени и без указателя. Регистрация — через дженерик:

```go
// core/tool.go
func JiraHandler[T Runnable](v T, app *AppState) mcpserver.ToolHandlerFunc
func ConfluenceHandler[T Runnable](v T, app *AppState) mcpserver.ToolHandlerFunc
```

Чтобы добавить новый инструмент:
1. Создать `struct{}` с методами `McpTool()` и `Run(ctx, base *core.Tool)`.
2. Добавить в `ReadTools()` / `WriteTools()` соответствующего файла.
3. Фетчеры доступны через `base.Jira` / `base.Confluence`, аргументы — через `base.Args`.

## Запуск

```bash
export JIRA_URL=https://jira.example.com
export JIRA_PERSONAL_TOKEN=<token>
export CONFLUENCE_URL=https://confluence.example.com
export CONFLUENCE_PERSONAL_TOKEN=<token>
go run ./cmd/mcp-jiracon
```

Старый Python-код: `_old_python/`.
