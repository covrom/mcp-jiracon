# mcp-atlassian-go

MCP-сервер для **Atlassian Jira и Confluence Data Center** с аутентификацией
через Personal Access Token (PAT). Экспонирует **71 MCP-инструмент** (53 Jira +
18 Confluence). Полностью совместим с Python-версией `mcp-atlassian` для DC.

## Требования

- Go 1.22+
- Jira Data Center / Server (любая версия с REST API v2)
- Confluence Data Center / Server (любая версия с REST API v1)
- Personal Access Token для каждого сервиса

## Быстрый старт

### 1. Установка

```bash
git clone git@gitlab.services.mts.ru:rstsovan/mcp-atlassian.git
cd mcp-atlassian/mcp-atlassian-go
go build -o mcp-atlassian ./cmd/mcp-jiracon
```

### 2. Настройка

```bash
# Jira Data Center
export JIRA_URL="https://jira.example.com"
export JIRA_PERSONAL_TOKEN="your-jira-pat-token"

# Confluence Data Center
export CONFLUENCE_URL="https://confluence.example.com"
export CONFLUENCE_PERSONAL_TOKEN="your-confluence-pat-token"
```

### 3. Запуск

```bash
# Через stdio (стандартный MCP транспорт)
./mcp-atlassian

# Через HTTP SSE (для отладки и тестирования)
./mcp-atlassian -transport=http -addr=:3000
```

## Переменные окружения

### Обязательные

| Переменная | Описание |
| --- | --- |
| `JIRA_URL` | URL Jira Data Center. Пример: `https://jira.example.com` |
| `JIRA_PERSONAL_TOKEN` | Personal Access Token для Jira |
| `CONFLUENCE_URL` | URL Confluence Data Center. Пример: `https://confluence.example.com` |
| `CONFLUENCE_PERSONAL_TOKEN` | Personal Access Token для Confluence |

### Безопасность и сеть

| Переменная | По умолчанию | Описание |
| --- | --- | --- |
| `JIRA_SSL_VERIFY` | `true` | Проверка TLS-сертификата Jira. `false` для самоподписанных сертификатов |
| `CONFLUENCE_SSL_VERIFY` | `true` | Проверка TLS-сертификата Confluence |
| `JIRA_HTTP_PROXY` | — | HTTP/HTTPS прокси для Jira |
| `JIRA_HTTPS_PROXY` | — | HTTPS прокси для Jira |
| `JIRA_NO_PROXY` | — | Исключения прокси для Jira |
| `CONFLUENCE_HTTP_PROXY` | — | HTTP/HTTPS прокси для Confluence |
| `CONFLUENCE_HTTPS_PROXY` | — | HTTPS прокси для Confluence |
| `CONFLUENCE_NO_PROXY` | — | Исключения прокси для Confluence |
| `MCP_ALLOWED_URL_DOMAINS` | — | Разрешённые домены (SSRF-защита) |

Также поддерживаются стандартные переменные: `HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, `SOCKS_PROXY`, `SSL_CA_FILE`.

### Режимы работы

| Переменная | По умолчанию | Описание |
| --- | --- | --- |
| `READ_ONLY_MODE` | `false` | `true` — запрещает все write-инструменты |
| `JIRA_PROJECTS_FILTER` | — | Фильтр проектов (через запятую). Например: `PROJ,DEV` |
| `CONFLUENCE_SPACES_FILTER` | — | Фильтр пространств (через запятую) |
| `MCP_LOG_LEVEL` | `INFO` | Уровень логирования: `DEBUG`, `INFO`, `WARN`, `ERROR` |

## Аутентификация

Сервер использует **только PAT-аутентификацию**. Токен передаётся в заголовке
каждого HTTP-запроса к Jira/Confluence:

```
Authorization: Bearer <PERSONAL_TOKEN>
```

### Multi-tenant режим

Для per-request аутентификации (разные учётные данные для разных
пользователей) используются HTTP-заголовки:

```bash
curl -H 'X-Atlassian-Jira-Url: https://jira.example.com' \
     -H 'X-Atlassian-Jira-Personal-Token: <pat>' \
     -H 'X-Atlassian-Confluence-Url: https://confluence.example.com' \
     -H 'X-Atlassian-Confluence-Personal-Token: <pat>' \
     http://localhost:3000/mcp
```

## Инструменты

Сервер экспонирует **71 MCP-инструмент**:

### Jira (53 инструмента)

**Read (31):**
`get_issue`, `search` (JQL), `get_project_issues`, `get_user_profile`, `search_assignable_users`,
`get_issue_watchers`, `get_transitions`, `get_worklog`, `get_link_types`, `get_all_projects`,
`search_projects`, `get_project_issue_types`, `get_create_fields`, `get_project_versions`,
`get_project_components`, `get_project_fields`, `get_agile_boards`, `get_board_issues`,
`get_sprints_from_board`, `get_sprint_issues`, `search_fields`, `get_field_options`,
`get_issue_dates`, `get_issue_sla`, `get_issue_images`, `download_attachments`,
`get_service_desk_for_project`, `get_service_desk_queues`, `get_queue_issues`,
`get_request_types`, `get_request_type_fields`

**Write (22):**
`create_issue`, `update_issue`, `delete_issue`, `assign_issue`, `transition_issue`,
`add_worklog`, `add_watcher`, `remove_watcher`, `link_to_epic`, `create_issue_link`,
`create_remote_issue_link`, `remove_issue_link`, `add_comment`, `edit_comment`,
`create_version`, `update_version`, `batch_create_issues`,
`create_sprint`, `update_sprint`, `add_issues_to_sprint`, `move_issues_to_backlog`,
`create_customer_request`

### Confluence (18 инструментов)

**Read (8):**
`get_page`, `search` (CQL), `get_page_children`, `get_space_page_tree`,
`get_comments`, `get_labels`, `search_user` (DC: через group member API),
`get_attachments`

**Write (10):**
`create_page`, `update_page`, `delete_page`, `move_page`,
`add_comment`, `reply_to_comment`, `add_label`,
`delete_attachment`

## Сборка и запуск

```bash
# Установка зависимостей
go mod download

# Сборка
go build -ldflags "-s -w" -o mcp-atlassian ./cmd/mcp-jiracon

# Тесты
go test ./... -count=1

# Линтинг
go vet ./...
```

### Docker

```bash
docker build -t mcp-atlassian .
docker run \
  -e JIRA_URL=https://jira.example.com \
  -e JIRA_PERSONAL_TOKEN=<token> \
  -e CONFLUENCE_URL=https://confluence.example.com \
  -e CONFLUENCE_PERSONAL_TOKEN=<token> \
  mcp-atlassian
```

### Helm

```bash
helm install mcp-atlassian ./helm \
  --set jira.url=https://jira.example.com \
  --set jira.personalToken=<token> \
  --set confluence.url=https://confluence.example.com \
  --set confluence.personalToken=<token>
```

## Настройка MCP-клиента

### Claude Desktop

В файле `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "atlassian": {
      "command": "/path/to/mcp-atlassian/mcp-atlassian-go/mcp-atlassian",
      "env": {
        "JIRA_URL": "https://jira.example.com",
        "JIRA_PERSONAL_TOKEN": "<token>",
        "CONFLUENCE_URL": "https://confluence.example.com",
        "CONFLUENCE_PERSONAL_TOKEN": "<token>"
      }
    }
  }
}
```

### Любой MCP-клиент

Сервер поддерживает транспорт **stdio** (по умолчанию) и **HTTP SSE** (флаг `-transport=http`).

## Архитектура

```
MCP-клиент (LLM)
    │  stdio или HTTP SSE
    ▼
Atlassian MCP Server (mark3labs/mcp-go)
    │  AuthMiddleware: Bearer <PAT> + per-request headers
    ▼
RegisterAllTools → McpTool() + JiraHandler/ConfluenceHandler
    │  base.Jira / base.Confluence из app state + request context
    ▼
JiraFetcher / ConfluenceFetcher
    │  AtlassianClient.Get/Post/Put/Delete
    ▼
Jira REST API v2 / Confluence REST API v1 (Data Center)
```

### Архитектура инструментов

Каждый инструмент — пустая структура без эмбеддинга. Контекст (`*core.Tool`)
передаётся явным параметром в `Run`:

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

Регистрация — через дженерик `JiraHandler` / `ConfluenceHandler`:

```go
s.AddTool(t.McpTool(), core.JiraHandler(t, app))
s.AddTool(t.McpTool(), core.ConfluenceHandler(t, app))
```

## Тесты

138 тестов покрывают все 71 DC-инструмент через MCP SSE client с mock-бэкендами:

```
ok  atlassian        — PAT auth, TLS verify, Confluence PAT
ok  confluence       — 12 fetcher tests
ok  jira             — 38 fetcher tests
ok  preprocessing    — 4 Jira wiki ↔ MD tests
ok  test/integration — 79 MCP tool tests (все 71 инструмент + read-only)
```

## Лицензия

Apache 2.0