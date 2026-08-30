# План переписывания `mcp-atlassian` на Go

> **Целевой стек**: Go 1.22+, [`github.com/mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go), `net/http` (стандартная библиотека).
> **Цель**: один бинарь, тот же внешний контракт — 98 MCP-инструментов с теми же именами, JSON-schema, return-формами, переменными окружения, поведением аутентификации.

---

## 1. Цели и границы

### Что переносим в точности

- **98 MCP-инструментов** (63 Jira + 35 Confluence) с теми же именами, тегами и JSON-schema.
- **Auth**: Basic, Bearer (PAT и OAuth 3LO), mTLS, External — в порядке приоритета.
- **SSRF-защита** с DNS-пиннингом и allowlist для redirect-ов.
- **Multi-tenant** через per-request headers (`X-Atlassian-Jira-Personal-Token`, `X-Atlassian-Confluence-Url` и т.д.).
- **Toolset-фильтрация** (`ENABLED_TOOLS`, `TOOLSETS`, `READ_ONLY_MODE`).
- **Atlassian Cloud OAuth через gateway URL** (`https://api.atlassian.com/ex/jira|confluence/{cloudId}/...`).
- **Переменные окружения** — те же имена, та же семантика.
- **Preprocessing** (Markdown ↔ storage) для Jira и Confluence.

### Что упрощаем

- Нет mixin-иерархии — Go-интерфейсы проще.
- Нет `asyncio`/`anyio` — обычные goroutine + `sync/errgroup`.
- Нет `atlassian-python-api` — пишем свой тонкий клиент (REST тривиален).
- Нет `pickle`/Pydantic — структуры с `json` тегами и `UnmarshalJSON` для дефолтов.
- Нет `requests.Session` chain — `http.RoundTripper` chain.

### Что НЕ делаем в этой фазе

- OAuth proxy / DCR (отложим на phase 9).
- PAC/WPAD-прокси (большинство операторов используют прямой HTTP_PROXY).
- SSL ignore (`SSL_VERIFY=false`) — отметим `TOFO` (требует обсуждения).
- Legacy Jira wiki ↔ Markdown (только Markdown → wiki в первой версии).

---

## 2. Структура проекта

```
mcp-atlassian/
├── go.mod
├── go.sum
├── cmd/
│   └── mcp-atlassian/
│       └── main.go              # entry point: флаги, lifespan, HTTP-server
├── internal/
│   ├── config/                  # env-парсинг + JiraConfig/ConfluenceConfig
│   │   ├── config.go
│   │   ├── env.go               # isEnvTruthy, getIntEnv, getHeaderNames
│   │   └── jira.go              # JiraConfig struct
│   │   └── confluence.go        # ConfluenceConfig struct
│   ├── auth/                    # OAuth 2.0 + token storage + per-request resolution
│   │   ├── oauth.go             # OAuthConfig + BYOAccessToken
│   │   ├── token_storage.go     # keyring +0600 file
│   │   ├── setup.go             # CLI setup wizard
│   │   └── middleware.go        # bearer header injection
│   ├── urls/                    # isCloudURL, validateSSRF, resolveRelative
│   ├── atlassian/               # базовый HTTP client + REST helpers
│   │   ├── client.go            # AtlassianClient: http.Client + middleware chain
│   │   ├── ssrf.go              # DNS-пин dial-context
│   │   ├── retry.go             # RoundTripper
│   │   ├── ratelimit.go         # x/time/rate
│   │   ├── circuit.go           # gobreaker
│   │   └── errors.go            # ErrAuthentication + классификации
│   ├── models/                  # Go-структуры вместо Pydantic
│   │   ├── base.go              # FromAPIResponse + ToSimplifiedDict contract
│   │   ├── jira/                # JiraIssue, JiraComment, JiraWorklog, JiraUser, ...
│   │   └── confluence/          # ConfluencePage, ConfluenceComment, ...
│   ├── preprocessing/
│   │   ├── base.go              # общие HTML-помощники
│   │   ├── confluence.go        # MD → storage (XHTML)
│   │   └── jira.go              # MD ↔ wiki
│   ├── jira/                    # 21 mixin → пакеты с методами
│   │   ├── issues.go
│   │   ├── search.go
│   │   ├── users.go
│   │   ├── comments.go
│   │   ├── worklog.go
│   │   ├── transitions.go
│   │   ├── links.go
│   │   ├── boards.go
│   │   ├── sprints.go
│   │   ├── queues.go
│   │   ├── epics.go
│   │   ├── fields.go
│   │   ├── field_options.go
│   │   ├── projects.go
│   │   ├── metrics.go
│   │   ├── sla.go
│   │   ├── development.go
│   │   ├── project_analysis.go
│   │   ├── attachments.go
│   │   ├── watchers.go
│   │   └── customer_requests.go
│   ├── confluence/              # 11 mixin → пакеты с методами
│   │   ├── pages.go
│   │   ├── search.go
│   │   ├── spaces.go
│   │   ├── comments.go
│   │   ├── labels.go
│   │   ├── users.go
│   │   ├── analytics.go
│   │   ├── attachments.go
│   │   ├── templates.go
│   │   ├── permissions.go
│   │   └── restrictions.go
│   ├── server/                  # главное: mcp-go integration
│   │   ├── server.go            # NewMCPServer + монтирование jira/confluence
│   │   ├── jira_tools.go        # 63 Jira-инструмента (s.AddTool)
│   │   ├── confluence_tools.go  # 35 Confluence-инструментов
│   │   ├── middleware.go        # per-request auth + tool filtering
│   │   ├── fetcher.go           # JiraFetcher/ConfluenceFetcher (per-request cached)
│   │   ├── filter.go            # _is_tool_enabled/_is_tool_authorized
│   │   ├── registry.go          # реестр инструментов + toolsets
│   │   ├── errors.go            # @handle_tool_errors + @check_write_access
│   │   ├── http.go              # ASGI/HTTP middleware: UserTokenMiddleware
│   │   └── oauth_proxy.go       # OAuth proxy (phase 9)
│   └── utils/
│       ├── logging.go           # slog setup + maskSensitive
│       ├── proxy.go             # HTTP_PROXY/HTTPS_PROXY/NO_PROXY
│       ├── ssl.go               # SSLIgnore + mTLS config
│       ├── media.go             # ATTACHMENT_MAX_BYTES + image detection
│       ├── pagination.go        # clamp_limit
│       ├── lifecycle.go         # signal handlers + clean exit
│       ├── user_agent.go        # version string
│       └── toolsets.go          # ToolsetDefinition + filtering
├── pkg/                         # (опционально) публичные пакеты для embed
│   └── version/
│       └── version.go           # -ldflags "-X main.version=..."
├── docs/
│   ├── README.md
│   ├── configuration.md
│   └── tools/                   # перенесённые из Python docs
├── test/                        # интеграционные тесты
│   └── integration/
├── Dockerfile
├── helm/
└── .github/
    └── workflows/
        └── ci.yml
```

### `go.mod`

```go
module github.com/your-org/mcp-atlassian

go 1.22

require (
    github.com/mark3labs/mcp-go v0.20.0             // MCP SDK
    github.com/zalando/go-keyring v0.2.4            // OS keyring
    golang.org/x/oauth2 v0.21.0                    // OAuth 2.0 клиент
    golang.org/x/sync v0.7.0                       // semaphore
    golang.org/x/time v0.5.0                       // rate
    golang.org/x/net v0.27.0                       // html-парсер
    github.com/JohannesKaufmann/html-to-markdown/v2 v2.0.0
    github.com/PuerkitoBio/goquery v1.9.1
    github.com/sony/gobreaker v1.0.0               // circuit breaker
    github.com/hashicorp/go-retryablehttp v0.7.5  // retry
)
```

---

## 3. Архитектурные решения

### 3.1. mcp-go: базовая модель

`mcp-go` предоставляет два главных пакета:

- `github.com/mark3labs/mcp-go/mcp` — типы (`Tool`, `ToolInputSchema`, `CallToolResult` и т.д.).
- `github.com/mark3labs/mcp-go/server` — `MCPServer`, транспорты (`Stdio`, `StreamableHTTP`).

Ключевые API:

```go
import (
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

s := server.NewMCPServer(
    "Atlassian MCP", "1.0.0",
    server.WithToolCapabilities(false),
    server.WithLogging(),
)

tool := mcp.NewTool("jira_get_issue",
    mcp.WithDescription("Get details of a specific Jira issue..."),
    mcp.WithString("issue_key",
        mcp.Required(),
        mcp.Pattern(`^[A-Z][A-Z0-9_]+-\d+(?:-\d+)*$`),
    ),
    mcp.WithString("fields", mcp.DefaultString("priority,updated,...")),
    mcp.WithNumber("comment_limit", mcp.DefaultNumber(10), mcp.Min(0), mcp.Max(100)),
)

s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := req.Params.Arguments.(map[string]any)
    issueKey := args["issue_key"].(string)
    // ... обработка
    return mcp.NewToolResultText(jsonString), nil
})
```

### 3.2. Подход с одной точкой монтирования

Python-версия монтирует `jira_mcp` и `confluence_mcp` в `main_mcp` под namespace'ами. В mcp-go такой функции нет, поэтому мы префиксируем имена инструментов вручную:

```go
type serverSet struct {
    jira       *server.MCPServer
    confluence *server.MCPServer
}

func (s *serverSet) addTool(namespace string, t mcp.Tool, h server.ToolHandlerFunc) {
    t.Name = namespace + "_" + t.Name  // "jira_get_issue", "confluence_get_page"
    switch namespace {
    case "jira":
        s.jira.AddTool(t, h)
    case "confluence":
        s.confluence.AddTool(t, h)
    }
}
```

Или — **проще** — один `MCPServer` с плоским списком инструментов, имена уже префиксированы в коде. Это убирает необходимость в монтировании и делает код линейным.

> **Рекомендация**: плоский подход. Один `MCPServer`, все 98 инструментов зарегистрированы напрямую. Mounting оставлен как legacy-концепция, в Go-порте он не нужен.

### 3.3. Per-request fetcher (бывший `get_jira_fetcher`)

В Python fetcher создаётся через `async def get_jira_fetcher(ctx)`, кладётся в `request.state`, используется в хендлере. В Go с mcp-go запрос приходит в `CallToolRequest`, и у нас нет прямого `request.state`. Решение:

```go
type ctxKey int
const fetcherKey ctxKey = iota

func fetcherFromContext(ctx context.Context) (*JiraFetcher, bool) {
    f, ok := ctx.Value(fetcherKey).(*JiraFetcher)
    return f, ok
}

func contextWithFetcher(ctx context.Context, f *JiraFetcher) context.Context {
    return context.WithValue(ctx, fetcherKey, f)
}

// В middleware:
ctx = contextWithFetcher(r.Context(), f)
```

Или ещё проще: middleware кладёт fetcher прямо в локальную переменную замыкания HTTP-handler'а:

```go
http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
    fetcher := resolveJiraFetcher(r)  // per-request
    ctx := withFetcher(r.Context(), fetcher)
    s.ServeHTTP(w, r.WithContext(ctx))
})
```

### 3.4. Структура HTTP-клиента

```go
type AtlassianClient struct {
    baseURL    string
    isCloud    bool
    authType   string  // "basic" | "pat" | "oauth" | "cert" | "external"
    httpClient *http.Client
    session    *JiraFetcher  // или *ConfluenceFetcher
}

func New(cfg Config) (*AtlassianClient, error) {
    transport := buildTransport(cfg)  // см. ниже
    return &AtlassianClient{
        baseURL:    cfg.URL,
        isCloud:    cfg.IsCloud,
        authType:   cfg.AuthType,
        httpClient: &http.Client{Transport: transport, Timeout: 30 * time.Second},
    }, nil
}

func buildTransport(cfg Config) http.RoundTripper {
    var rt http.RoundTripper = baseRoundTripper(cfg)  // dial context, TLS
    if cfg.AuthType == "basic" || cfg.AuthType == "pat" || cfg.AuthType == "oauth" {
        rt = withAuthHeader(rt, cfg.Credentials)
    }
    rt = withSSRFGuard(rt, cfg.AllowedHosts)
    rt = withRetry(rt, cfg.Retry)
    rt = withRateLimit(rt, cfg.Rate)
    rt = withCircuitBreaker(rt, cfg.CB)
    return rt
}
```

Каждый `with*` — функция-обёртка:

```go
func withAuthHeader(base http.RoundTripper, creds Credentials) http.RoundTripper {
    return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
        req2 := req.Clone(req.Context())
        if creds.Type == "basic" {
            req2.Header.Set("Authorization",
                "Basic "+base64.StdEncoding.EncodeToString([]byte(creds.Email+":"+creds.Token)))
        } else {
            req2.Header.Set("Authorization", "Bearer "+creds.Token)
        }
        return base.RoundTrip(req2)
    })
}
```

### 3.5. SSRF-пиннинг

```go
type ssrfGuard struct {
    base         *http.Transport
    trustedHosts map[string]struct{}
    resolver     *net.Resolver
}

func (g *ssrfGuard) RoundTrip(req *http.Request) (*http.Response, error) {
    host := req.URL.Hostname()
    ips, err := g.resolver.LookupIPAddr(req.Context(), host)
    if err != nil {
        return nil, err
    }
    for _, ip := range ips {
        if _, trusted := g.trustedHosts[host]; !trusted {
            if !ip.IP.IsGlobalUnicast() || ip.IP.IsPrivate() || ip.IP.IsLoopback() {
                return nil, fmt.Errorf("ssrf: %s resolves to non-global IP %s", host, ip.IP)
            }
        }
    }
    return g.base.RoundTrip(req)
}

// Плюс CheckRedirect:
client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        if err := validateSSRF(req.URL.String()); err != nil {
            return err
        }
        return nil
    },
}
```

### 3.6. Per-request auth resolution

```go
type authBranch int
const (
    branchHeaderPAT authBranch = iota
    branchBasic
    branchBearer
    branchGlobalFallback
)

func resolveJiraFetcher(r *http.Request, appState *AppState) (*JiraFetcher, error) {
    // 1. Header-PAT
    if url := r.Header.Get("X-Atlassian-Jira-Url"); url != "" {
        if tok := r.Header.Get("X-Atlassian-Jira-Personal-Token"); tok != "" {
            return buildHeaderPATFetcher(url, tok, appState.NetworkOpts)
        }
    }
    // 2. Basic
    if user, tok, ok := r.BasicAuth(); ok {
        return buildBasicFetcher(appState.JiraConfig, user, tok, appState.NetworkOpts)
    }
    // 3. Bearer
    if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
        token := strings.TrimPrefix(h, "Bearer ")
        cloudID := r.Header.Get("X-Atlassian-Cloud-Id")
        return buildBearerFetcher(appState.JiraConfig, token, cloudID, appState.NetworkOpts)
    }
    // 4. Global fallback (только в stdio или если ALLOW_GLOBAL_CRED_FALLBACK)
    if !inHTTPContext(r) || isEnvTruthy("ALLOW_GLOBAL_CRED_FALLBACK") {
        return appState.GlobalJiraFetcher, nil
    }
    return nil, errors.New("no Atlassian credentials provided")
}
```

### 3.7. Tool-фильтрация

```go
type toolMeta struct {
    tags        []string  // {"jira", "read", "toolset:jira_issues"}
    title       string
    readOnly    bool      // вычислено из тегов
    destructive bool      // вычислено из тегов
}

var tools = map[string]*toolMeta{
    "jira_get_issue": {tags: []string{"jira", "read", "toolset:jira_issues"}, title: "Get Issue"},
    // ... для всех 98
}

func isAuthorized(name string, appState *AppState, r *http.Request) bool {
    meta := tools[name]
    // 1. Toolset
    if !toolsetEnabled(meta.tags, appState.EnabledToolsets) {
        return false
    }
    // 2. ENABLED_TOOLS
    if !toolAllowed(name, appState.EnabledTools) {
        return false
    }
    // 3. Read-only
    if appState.ReadOnly && contains(meta.tags, "write") {
        return false
    }
    return true
}

// Применяется дважды:
// - В `ListTools` фильтрует tools/list
// - В `CallTool` отклоняет с NotFound если не authorized
```

### 3.8. Read-only guard

```go
func checkWriteAccess(ctx context.Context) error {
    if readOnlyFromContext(ctx) {
        return fmt.Errorf("read-only mode: write tools are disabled")
    }
    return nil
}

// Обёртка для хендлера:
func writeTool(name string, h ToolHandler) ToolHandler {
    return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        if err := checkWriteAccess(ctx); err != nil {
            return nil, err
        }
        return h(ctx, req)
    }
}
```

---

## 4. Детальный план по фазам

### Phase 1: Foundation (1–2 дня)

Цель: пустой бинарь, который стартует, читает env, пишет в лог.

**Задачи:**

1. `go mod init github.com/your-org/mcp-atlassian`
2. Создать структуру папок (см. раздел 2).
3. `internal/config/env.go` — портировать `utils/env.py`:
   - `IsEnvTruthy(name string, defaultVal ...string) bool`
   - `GetIntEnv(name string, defaultVal int) int`
   - `GetCustomHeaders(name string) map[string]string`
   - `GetHeaderNames(name string) []string`
4. `internal/utils/logging.go` — `SetupLogging(level slog.Level) *slog.Logger` с теми же масками.
5. `internal/utils/lifecycle.go` — `signal.Notify` + `os.Stdout.Sync()`.
6. `cmd/mcp-jiracon/main.go` — точка входа:
   - Парсит флаги (или env).
   - `slog.Info("starting", "version", version)`.
   - Вызывает `server.Run(ctx)`.

**Done when:** `uv run mcp-atlassian` → `go run ./cmd/mcp-jiracon` запускается и пишет логи.

### Phase 2: HTTP client + Auth (3–4 дня)

Цель: умеем делать авторизованные запросы к Jira и Confluence.

**Задачи:**

1. `internal/atlassian/client.go` — `AtlassianClient` с `*http.Client`.
2. `internal/atlassian/errors.go`:
   - `ErrAuthentication`
   - `IsAuthError(err) bool`
   - `ClassifyHTTPError(resp *http.Response) error` — маппит 401/403 → `ErrAuthentication`.
3. `internal/atlassian/ssrf.go` — SSRF guard с DNS-пиннингом.
4. `internal/atlassian/retry.go` — `WithRetry` (exponential backoff + Retry-After).
5. `internal/atlassian/ratelimit.go` — `x/time/rate` wrapper.
6. `internal/atlassian/circuit.go` — `gobreaker`.
7. `internal/urls/urls.go`:
   - `IsAtlassianCloudURL(rawURL string) bool`
   - `ValidateURLForSSRF(rawURL string) error`
   - `ResolveRelativeURL(url, baseURL string) string`
8. `internal/utils/proxy.go` — поддержка `HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY` (через `http.ProxyFromEnvironment`).
9. `internal/utils/ssl.go` — `tls.Config{Certificates: []tls.Certificate{...}}`, отказ от encrypted keys.
10. `internal/config/jira.go` и `confluence.go` — структуры `JiraConfig`, `ConfluenceConfig`, `FromEnv()`.
11. `internal/auth/oauth.go`:
    - `OAuthConfig` (cloud + DC)
    - `BYOAccessTokenOAuthConfig`
    - `GetOAuthConfigFromEnv()`
    - `ConfigureOAuthSession` — ставит `Authorization: Bearer ...` на запрос.
12. `internal/auth/token_storage.go` — keyring +0600 file fallback.
13. `internal/server/middleware.go` — `UserTokenMiddleware` для HTTP-транспорта:
    - Парсит `Authorization: Basic/Bearer/Token`.
    - Парсит `X-Atlassian-*-Url`, `X-Atlassian-*-Personal-Token`.
    - Валидирует URL через `ValidateURLForSSRF`.
    - Сохраняет credentials в `request.Context()` через typed keys.
14. `internal/server/fetcher.go` — резолвер fetcher'а по request context (4 branches).

**Done when:** простой Go-тест делает `GET /rest/api/3/myself` к Jira Cloud и получает 200.

### Phase 3: Models (2–3 дня)

Цель: умеем парсить ответы Jira и Confluence в Go-структуры.

**Задачи:**

1. `internal/models/base.go`:
   - Интерфейс `FromAPIResponse(data map[string]any) error` — мутирует receiver.
   - `ToSimplifiedDict() map[string]any` — на receiver.
   - `FormatTimestamp(s string) string` — `time.Parse(time.RFC3339, ...)`, fallback на исходную строку.
2. `internal/models/jira/` — портировать каждую модель:
   - `issue.go` — `JiraIssue`, `JiraField` (с default values).
   - `comment.go` — `JiraComment`.
   - `worklog.go` — `JiraWorklog`.
   - `user.go` — `JiraUser` с ветвлением Cloud (`accountId`) vs DC (`name`/`key`).
   - `common.go` — `JiraStatus`, `JiraPriority`, `JiraIssueType` (с дефолтами).
   - `project.go`, `attachment.go`, `version.go`, `component.go`,
     `transition.go`, `link.go`, `sprint.go`, `board.go`,
     `changelog.go`.
3. `internal/models/confluence/`:
   - `page.go` — `ConfluencePage`.
   - `comment.go`, `space.go`, `attachment.go`, `label.go`,
     `user.go`, `version.go`, `restriction.go`.
4. `internal/models/constants.go` — `EmptyString`, `Unknown`, `Unassigned`,
   `JiraDefaultStatus`, `ConfluenceDefaultVersion` и т.д.

**Done when:** `ToSimplifiedDict()` для JiraIssue выдаёт тот же JSON, что и Python-версия (snapshot-тест).

### Phase 4: Internal Jira mixins (5–7 дней)

Цель: 21 mixin реализован как методы на `JiraFetcher`.

**Задачи:**

1. `internal/jira/fetcher.go` — `JiraFetcher` хранит `*atlassian.Client` + per-fetcher options.
2. По одному файлу на mixin (см. структуру выше).
3. Каждая функция делает 1+ HTTP-запрос через клиент, парсит ответ в модель, возвращает модель + ошибку.
4. Cloud/DC-ветвления:
   - URL: `{base}/rest/api/3/...` (Cloud) vs `{base}/rest/api/2/...` (DC).
   - User identifier: `accountId` vs `name`/`key`.
   - Epic link field: `customfield_12311140` vs `Epic Link`.
5. **Async move**: реализовать polling `GET /bulk/assign/move/{taskId}` каждые 2с с таймаутом 30с.

**Done when:** все 21 mixin'а компилируются и имеют unit-тесты с mocked HTTP-клиентом.

### Phase 5: Internal Confluence mixins (3–4 дня)

Цель: 11 mixin'ов Confluence.

**Задачи:**

1. `internal/confluence/fetcher.go` — аналог `ConfluenceFetcher`.
2. `internal/confluence/v2_adapter.go` — обёртка для v2-API (Cloud OAuth).
3. 11 mixin'ов (pages, search, spaces, comments, labels, users,
   analytics, attachments, templates, permissions, restrictions).
4. v1/v2 selection logic:
   - v2: pages CRUD, comments footer/inline, attachments.
   - v1: move, copy, restrictions, analytics.

**Done when:** все 11 mixin'ов компилируются и проходят unit-тесты.

### Phase 6: Read-only Jira tools (3–4 дня)

Цель: первые Jira-инструменты работают end-to-end.

**Задачи:**

1. `internal/server/registry.go` — таблица `tools map[string]*ToolMeta` со всеми 98 инструментами.
2. `internal/server/filter.go` — `isAuthorized(name, appState, r) bool`.
3. `internal/server/jira_tools.go` — реализация 35 read-инструментов:
   - `jira_get_user_profile`, `jira_search_assignable_users`
   - `jira_get_issue_watchers`
   - `jira_get_issue`, `jira_search`, `jira_get_project_issues`
   - `jira_get_transitions`, `jira_get_worklog`
   - `jira_get_field_options`, `jira_search_fields`
   - `jira_get_agile_boards`, `jira_get_board_issues`,
     `jira_get_sprints_from_board`, `jira_get_sprint_issues`
   - `jira_get_link_types`
   - `jira_get_project_issue_types`, `jira_get_create_fields`,
     `jira_get_project_versions`, `jira_get_project_components`,
     `jira_get_all_projects`, `jira_search_projects`, `jira_get_project_fields`
   - `jira_get_service_desk_*` (DC)
   - `jira_get_issue_proforma_forms`, `jira_get_proforma_form_details`
   - `jira_get_issue_dates`, `jira_get_issue_sla`
   - `jira_get_issue_development_info`, `jira_get_issues_development_info`
   - `jira_get_project_epic_hierarchy`, `jira_get_cross_project_dependencies`
   - `jira_get_attachments` (binary → `EmbeddedResource`).
4. Возвращают `mcp.NewToolResultText(jsonString)` или `mcp.NewToolResultError`.
5. Для binary-инструментов: `[]mcp.Content{ mcp.NewToolResultImage(...), ... }`.

**Done when:** каждый инструмент проходит integration-тест с реальным Jira (mock-server).

### Phase 7: Write tools (3–4 дня)

Цель: 28 write-инструментов Jira.

**Задачи:**

1. Обернуть каждый handler в `writeTool(name, h)` (см. 3.8).
2. Реализовать:
   - `jira_add_watcher`, `jira_remove_watcher`
   - `jira_create_issue`, `jira_batch_create_issues`, `jira_update_issue`,
     `jira_assign_issue`, `jira_delete_issue`, `jira_move_issue`
   - `jira_transition_issue`
   - `jira_add_worklog`
   - `jira_link_to_epic`, `jira_create_issue_link`,
     `jira_create_remote_issue_link`, `jira_remove_issue_link`
   - `jira_add_comment`, `jira_edit_comment`
   - `jira_create_version`, `jira_batch_create_versions`,
     `jira_update_version`
   - `jira_create_customer_request`
   - `jira_update_proforma_form_answers`
   - `jira_create_sprint`, `jira_update_sprint`,
     `jira_add_issues_to_sprint`, `jira_move_issues_to_backlog`
3. **JSM internal-only guard**:
   - Проверять `JIRA_INTERNAL_ONLY_PROJECTS` для комментариев / transitions / links.
   - Возвращать ошибку `"Cannot post customer-visible comment in internal-only project"`.
4. `jira_add_comment` с `public` параметром — маршрут через ServiceDesk API.

**Done when:** все 28 write-инструментов работают с реальным Jira и обрабатывают read-only mode.

### Phase 8: Confluence tools (4–5 дней)

Цель: все 35 Confluence-инструментов.

**Задачи:**

1. Реализовать `confluence_tools.go`:
   - `confluence_search`, `confluence_get_page`, `confluence_get_page_children`,
     `confluence_get_space_page_tree`
   - `confluence_create_page`, `confluence_update_page`,
     `confluence_update_page_section`, `confluence_delete_page`,
     `confluence_move_page`, `confluence_copy_page`
   - `confluence_get_comments`, `confluence_add_comment`,
     `confluence_reply_to_comment`, `confluence_get_inline_comments`,
     `confluence_add_inline_comment`
   - `confluence_get_labels`, `confluence_add_label`
   - `confluence_search_user`
   - `confluence_get_page_views`
   - `confluence_upload_attachment`, `confluence_upload_attachments`,
     `confluence_get_attachments`, `confluence_download_attachment`,
     `confluence_download_content_attachments`, `confluence_delete_attachment`,
     `confluence_get_page_images`
   - `confluence_list_page_templates`, `confluence_get_page_template`,
     `confluence_create_page_from_template`
   - `confluence_check_content_permissions`, `confluence_get_space_permissions`
   - `confluence_get_page_restrictions`, `confluence_set_page_restrictions`
2. `_resolve_page_id` — портировать URL-decoder для tiny links.
3. **`content_format` маппинг**:
   - `markdown` → `isMarkdown=true, contentRepresentation=""` (конвертируется MD→storage)
   - `wiki`, `storage` → `isMarkdown=false, contentRepresentation="wiki"|"storage"`
   - `xhtml` → мапится в `storage`
4. Cloud-only: `confluence_get_page_views`, `confluence_*_templates`,
   `confluence_check_content_permissions`, `confluence_get_space_permissions`.
   Возвращать `"Not implemented on Data Center"`.

**Done when:** все 35 Confluence-инструментов проходят integration-тесты.

### Phase 9: Preprocessing (5–7 дней)

Цель: Markdown ↔ storage conversion.

**Задачи:**

1. `internal/preprocessing/base.go`:
   - HTML-парсинг через `golang.org/x/net/html`.
   - User mention resolution (Cloud: `accountId`; DC: `userkey`/`username`).
   - Image rewriting (`<ac:image>` ↔ `<ri:attachment>`).
2. `internal/preprocessing/confluence.go`:
   - MD → HTML через `html-to-markdown` (обратное направление — HTML → MD).
   - Wait, для Confluence нам нужно **MD → storage**, т.е. **MD → HTML**, не наоборот.
   - Использовать `github.com/JohannesKaufmann/html-to-markdown` для storage → MD.
   - Для **MD → storage**: написать свой конвертер на базе `goldmark` + пост-обработка.
3. `internal/preprocessing/jira.go`:
   - MD → Jira wiki: ~30 regex-замен.
   - Jira wiki → MD: обратные regex-замены.
   - Поддержка языков (`VALID_JIRA_LANGUAGES`, `LANGUAGE_MAPPING`).
4. Использовать placeholder-трюк (`\x00PREFIX<N>\x00`) для защиты code-блоков.
5. Smart links, mentions, panel, quote, noformat.

**Done when:** round-trip `jira_create_issue` (MD) → `jira_get_issue` (MD)
сохраняет форматирование, code-блоки, списки, таблицы.

### Phase 10: Server lifecycle + OAuth (3–4 дня)

Цель: реально запускаемый MCP-сервер с stdio и streamable-HTTP транспортами.

**Задачи:**

1. `internal/server/server.go`:
   ```go
   func NewServer(appState *AppState) *server.MCPServer {
       s := server.NewMCPServer("Atlassian MCP", version,
           server.WithToolCapabilities(false),
           server.WithLogging(),
       )
       registerAllTools(s, appState)
       return s
   }
   ```
2. `cmd/mcp-jiracon/main.go`:
   - Парсит env (`READ_ONLY_MODE`, `ENABLED_TOOLS`, `TOOLSETS`, host, port).
   - Строит `AppState` (lifespan).
   - Создаёт `MCPServer`.
   - Если `--transport http`: запускает `server.NewStreamableHTTPServer`.
   - Если `--transport stdio` (default): `server.ServeStdio`.
3. HTTP-middleware chain (применяется **перед** mcp-go handler):
   - `authHeaderMiddleware` — парсит credentials.
   - `toolFilterMiddleware` — отбрасывает tools/list, если hidden.
   - `userTokenMiddleware` — следит за drop/disconnect (ASGI).
4. `/healthz` endpoint.
5. **OAuth proxy (опционально, phase 10b)**: если
   `ATLASSIAN_OAUTH_PROXY_ENABLE=true`, экспонировать `/authorize`,
   `/token` endpoints. Для mcp-go это потребует отдельного HTTP-router'а
   рядом с MCP-router'ом.

**Done when:** `go run ./cmd/mcp-jiracon` запускается, доступен через
stdio (для тестов) и через HTTP (для интеграционных тестов).

### Phase 11: Тесты + документация (3–4 дня)

Цель: всё покрыто тестами, README объясняет как пользоваться.

**Задачи:**

1. **Unit-тесты** для каждого mixin'а с `httptest.Server` mock'ом.
2. **Integration-тесты** в `test/integration/` (build tag `integration`).
3. **E2E-тест** через `mcp-go/client` — стартуем сервер, делаем
   `tools/list`, `tools/call`.
4. README с инструкциями по сборке, конфигурации, deploy'у.
5. CHANGELOG + CONTRIBUTING.
6. Dockerfile (multi-stage, distroless).

**Done when:** `go test ./...` зелёный, integration-тесты с реальным
Jira/Confluence зелёные, README написан.

### Phase 12: Деплой и hardening (2–3 дня)

Цель: production-ready.

**Задачи:**

1. Helm chart (портировать из Python-версии).
2. Graceful shutdown (`signal.Notify`).
3. Structured logging через `slog`.
4. Metrics endpoint (опционально, `prometheus/client_golang`).
5. CI workflow: lint (`golangci-lint`), test, build, push image.
6. Release process: `goreleaser` или скриптованный `git tag`.

---

## 5. Ключевые модули с примерами кода

### 5.1. `internal/server/jira_tools.go`

```go
package server

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"

    "github.com/your-org/mcp-atlassian/internal/jira"
    "github.com/your-org/mcp-atlassian/internal/models/jira"
)

func registerJiraTools(s *server.MCPServer, appState *AppState) {
    s.AddTool(
        mcp.NewTool("jira_get_issue",
            mcp.WithDescription("Get details of a specific Jira issue..."),
            mcp.WithString("issue_key",
                mcp.Required(),
                mcp.Pattern(`^[A-Z][A-Z0-9_]+-\d+(?:-\d+)*$`),
            ),
            mcp.WithString("fields",
                mcp.DefaultString("priority,updated,labels,issuetype,summary,assignee,description,created,reporter,status"),
            ),
            mcp.WithString("expand", mcp.Description("Fields to expand")),
            mcp.WithNumber("comment_limit", mcp.DefaultNumber(10), mcp.Min(0), mcp.Max(100)),
            mcp.WithString("properties", mcp.Description("Comma-separated issue properties")),
            mcp.WithBoolean("update_history", mcp.DefaultNumber(true)),
            mcp.WithString("include", mcp.Description("Sections to inline: all|remote_links|transitions|watchers|changelog|comments|worklogs")),
            mcp.WithBoolean("use_display_names", mcp.DefaultNumber(false)),
        ),
        handleGetIssue(appState),
    )

    // ... остальные 62 инструмента
}

func handleGetIssue(appState *AppState) server.ToolHandlerFunc {
    return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        args := req.Params.Arguments.(map[string]any)
        issueKey := args["issue_key"].(string)
        fields := args["fields"].(string)

        fetcher, ok := jiraFetcherFromContext(ctx)
        if !ok {
            return nil, fmt.Errorf("no Jira fetcher in context")
        }

        issue, err := fetcher.GetIssue(ctx, jira.GetIssueOpts{
            IssueKey:      issueKey,
            Fields:        parseFields(fields),
            Expand:        stringArg(args, "expand"),
            CommentLimit:  intArg(args, "comment_limit", 10),
            Properties:    parseProperties(stringArg(args, "properties")),
            UpdateHistory: boolArg(args, "update_history", true),
            Include:       parseInclude(stringArg(args, "include")),
            UseDisplayNames: boolArg(args, "use_display_names", false),
        })
        if err != nil {
            return nil, fmt.Errorf("jira_get_issue failed: %w", err)
        }

        result := issue.ToSimplifiedDict()
        // Include секции...
        return mcp.NewToolResultText(mustMarshal(result)), nil
    }
}

func mustMarshal(v any) string {
    b, _ := json.MarshalIndent(v, "", "  ")
    return string(b)
}
```

### 5.2. `internal/jira/issues.go`

```go
package jira

import (
    "context"
    "net/http"
    "strconv"

    "github.com/your-org/mcp-atlassian/internal/atlassian"
    "github.com/your-org/mcp-atlassian/internal/models/jira"
)

type Fetcher struct {
    client *atlassian.AtlassianClient
    isCloud bool
}

func New(client *atlassian.AtlassianClient, isCloud bool) *Fetcher {
    return &Fetcher{client: client, isCloud: isCloud}
}

func (f *Fetcher) GetIssue(ctx context.Context, opts GetIssueOpts) (*jiraModels.JiraIssue, error) {
    apiVersion := f.apiVersion()
    path := fmt.Sprintf("/rest/api/%s/issue/%s", apiVersion, opts.IssueKey)

    fields := opts.Fields
    if fields == nil || len(fields) == 0 {
        fields = []string{"summary", "status", "assignee", "description"} // default
    }

    req, _ := http.NewRequestWithContext(ctx, "GET", path, nil)
    q := req.URL.Query()
    q.Set("fields", strings.Join(fields, ","))
    if opts.Expand != "" {
        q.Set("expand", opts.Expand)
    }
    if len(opts.Properties) > 0 {
        q.Set("properties", strings.Join(opts.Properties, ","))
    }
    q.Set("updateHistory", strconv.FormatBool(opts.UpdateHistory))
    req.URL.RawQuery = q.Encode()

    var raw map[string]any
    if err := f.client.Do(req, &raw); err != nil {
        return nil, fmt.Errorf("GET issue: %w", err)
    }

    issue := &jiraModels.JiraIssue{}
    if err := issue.FromAPIResponse(raw); err != nil {
        return nil, fmt.Errorf("parse issue: %w", err)
    }
    return issue, nil
}

func (f *Fetcher) apiVersion() string {
    if f.isCloud {
        return "3"
    }
    return "2"
}
```

### 5.3. `internal/atlassian/client.go`

```go
package atlassian

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type AtlassianClient struct {
    baseURL string
    isCloud bool
    http    *http.Client
}

func New(baseURL string, isCloud bool, opts ...Option) (*AtlassianClient, error) {
    cfg := &clientConfig{
        timeout: 30 * time.Second,
    }
    for _, o := range opts {
        o(cfg)
    }
    if cfg.transport == nil {
        cfg.transport = defaultTransport()
    }
    return &AtlassianClient{
        baseURL: baseURL,
        isCloud: isCloud,
        http:    &http.Client{Timeout: cfg.timeout, Transport: cfg.transport},
    }, nil
}

func (c *AtlassianClient) Do(req *http.Request, out any) error {
    req.URL.Scheme = "https"
    if req.URL.Host == "" {
        req.URL.Host = extractHost(c.baseURL)
    }
    if req.URL.Path[0] != '/' {
        req.URL.Path = "/" + req.URL.Path
    }
    if req.Header.Get("Accept") == "" {
        req.Header.Set("Accept", "application/json")
    }
    if out != nil && req.Body == nil {
        req.Header.Set("Content-Type", "application/json")
    }

    resp, err := c.http.Do(req)
    if err != nil {
        return fmt.Errorf("http: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return &HTTPError{
            StatusCode: resp.StatusCode,
            Body:       string(body),
            URL:        req.URL.String(),
        }
    }
    if out == nil {
        io.Copy(io.Discard, resp.Body)
        return nil
    }
    return json.NewDecoder(resp.Body).Decode(out)
}

func (c *AtlassianClient) Post(ctx context.Context, path string, body, out any) error {
    var buf bytes.Buffer
    if body != nil {
        if err := json.NewEncoder(&buf).Encode(body); err != nil {
            return err
        }
    }
    req, _ := http.NewRequestWithContext(ctx, "POST", path, &buf)
    return c.Do(req, out)
}

func (c *AtlassianClient) Get(ctx context.Context, path string, out any) error {
    req, _ := http.NewRequestWithContext(ctx, "GET", path, nil)
    return c.Do(req, out)
}

// ... Put, Delete, PutMultipart
```

### 5.4. `internal/server/middleware.go` — ASGI/HTTP auth

```go
package server

import (
    "context"
    "encoding/base64"
    "net/http"
    "strings"

    "github.com/your-org/mcp-atlassian/internal/urls"
)

type ctxKey int
const (
    ctxUserEmail ctxKey = iota
    ctxUserToken
    ctxAuthType
    ctxCloudID
)

func AuthMiddleware(appState *AppState) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. Service headers
            jiraURL := r.Header.Get("X-Atlassian-Jira-Url")
            jiraToken := r.Header.Get("X-Atlassian-Jira-Personal-Token")
            confluenceURL := r.Header.Get("X-Atlassian-Confluence-Url")
            confluenceToken := r.Header.Get("X-Atlassian-Confluence-Personal-Token")
            cloudID := r.Header.Get("X-Atlassian-Cloud-Id")

            // SSRF check
            if jiraURL != "" {
                if err := urls.ValidateURLForSSRF(jiraURL); err != nil {
                    http.Error(w, `{"error":"Forbidden: Invalid Jira URL"}`, http.StatusForbidden)
                    return
                }
            }
            if confluenceURL != "" {
                if err := urls.ValidateURLForSSRF(confluenceURL); err != nil {
                    http.Error(w, `{"error":"Forbidden: Invalid Confluence URL"}`, http.StatusForbidden)
                    return
                }
            }

            // 2. Authorization header
            auth := r.Header.Get("Authorization")
            ctx := r.Context()
            switch {
            case strings.HasPrefix(auth, "Basic "):
                payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
                if err != nil {
                    http.Error(w, `{"error":"Invalid Basic auth"}`, http.StatusUnauthorized)
                    return
                }
                email, token, ok := strings.Cut(string(payload), ":")
                if !ok || email == "" || token == "" {
                    http.Error(w, `{"error":"Invalid Basic auth format"}`, http.StatusUnauthorized)
                    return
                }
                ctx = context.WithValue(ctx, ctxUserEmail, email)
                ctx = context.WithValue(ctx, ctxUserToken, token)
                ctx = context.WithValue(ctx, ctxAuthType, "basic")
            case strings.HasPrefix(auth, "Bearer "):
                ctx = context.WithValue(ctx, ctxUserToken, strings.TrimPrefix(auth, "Bearer "))
                ctx = context.WithValue(ctx, ctxAuthType, "oauth")
            case strings.HasPrefix(auth, "Token "):
                ctx = context.WithValue(ctx, ctxUserToken, strings.TrimPrefix(auth, "Token "))
                ctx = context.WithValue(ctx, ctxAuthType, "pat")
            }
            if cloudID != "" {
                ctx = context.WithValue(ctx, ctxCloudID, cloudID)
            }

            // 3. Save service headers for downstream fetcher resolution
            // ... (use a custom struct in context)

            // 4. No credentials → reject (unless ALLOW_GLOBAL_CRED_FALLBACK)
            authType, _ := ctx.Value(ctxAuthType).(string)
            if authType == "" && jiraToken == "" && confluenceToken == "" {
                if !appState.ReadOnly && !appState.AllowGlobalFallback {
                    http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
                    return
                }
            }

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 5.5. `internal/server/fetcher.go` — резолвер fetcher'а

```go
package server

import (
    "context"
    "fmt"

    "github.com/your-org/mcp-atlassian/internal/config"
    "github.com/your-org/mcp-atlassian/internal/confluence"
    "github.com/your-org/mcp-atlassian/internal/jira"
)

type ctxKey int
const (
    ctxJiraFetcher ctxKey = iota + 100
    ctxConfluenceFetcher
)

func WithJiraFetcher(ctx context.Context, f *jira.Fetcher) context.Context {
    return context.WithValue(ctx, ctxJiraFetcher, f)
}

func JiraFetcherFromContext(ctx context.Context) (*jira.Fetcher, bool) {
    f, ok := ctx.Value(ctxJiraFetcher).(*jira.Fetcher)
    return f, ok
}

func ResolveJiraFetcher(ctx context.Context, appState *AppState) (*jira.Fetcher, error) {
    // 1. Header-PAT
    if url := headerValue(ctx, "X-Atlassian-Jira-Url"); url != "" {
        if tok := headerValue(ctx, "X-Atlassian-Jira-Personal-Token"); tok != "" {
            cfg := config.JiraConfig{
                URL:          url,
                AuthType:     "pat",
                PersonalToken: tok,
                IsCloud:      isCloudURL(url),
                ...network opts from appState
            }
            return jira.New(cfg)
        }
    }
    // 2. Basic auth
    if email, ok := ctx.Value(ctxUserEmail).(string); ok && email != "" {
        if token, ok := ctx.Value(ctxUserToken).(string); ok && token != "" {
            cfg := deriveBasicConfig(appState.JiraConfig, email, token)
            return jira.New(cfg)
        }
    }
    // 3. Bearer
    if token, ok := ctx.Value(ctxUserToken).(string); ok && token != "" {
        authType, _ := ctx.Value(ctxAuthType).(string)
        cloudID, _ := ctx.Value(ctxCloudID).(string)
        cfg := deriveBearerConfig(appState.JiraConfig, authType, token, cloudID)
        return jira.New(cfg)
    }
    // 4. Global fallback
    if appState.AllowGlobalFallback || appState.Transport == "stdio" {
        return appState.GlobalJiraFetcher, nil
    }
    return nil, fmt.Errorf("no Jira credentials available")
}
```

---

## 6. Тестирование

### 6.1. Unit-тесты

```go
func TestGetIssue_Success(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Проверяем URL, query params
        if r.URL.Path != "/rest/api/3/issue/PROJ-123" {
            t.Errorf("unexpected path: %s", r.URL.Path)
        }
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, `{"id":"10001","key":"PROJ-123","self":"...","fields":{...}}`)
    }))
    defer srv.Close()

    client := atlassian.New(srv.URL, true, atlassian.WithInsecureSSL())
    fetcher := jira.New(client, true)
    issue, err := fetcher.GetIssue(context.Background(), jira.GetIssueOpts{IssueKey: "PROJ-123"})
    require.NoError(t, err)
    assert.Equal(t, "PROJ-123", issue.Key)
}
```

### 6.2. Integration с mcp-go клиентом

```go
func TestServer_ListTools(t *testing.T) {
    s := server.NewMCPServer("Test", "0.0.0", server.WithToolCapabilities(false))
    appState := newTestAppState()  // с моками
    registerAllTools(s, appState)

    cli := client.NewClient(s)
    result, err := cli.ListTools(context.Background(), mcp.ListToolsRequest{})
    require.NoError(t, err)

    names := []string{}
    for _, tool := range result.Tools {
        names = append(names, tool.Name)
    }
    assert.Contains(t, names, "jira_get_issue")
    assert.Contains(t, names, "confluence_get_page")
}
```

### 6.3. Integration с реальным Atlassian

```go
//go:build integration

func TestJiraIntegration_GetIssue(t *testing.T) {
    if testing.Short() { t.Skip() }
    cfg := config.JiraConfigFromEnv()
    if !cfg.IsAuthConfigured() { t.Skip() }

    client, err := atlassian.New(cfg.URL, cfg.IsCloud,
        atlassian.WithAuth(credsFromCfg(cfg)))
    require.NoError(t, err)

    fetcher := jira.New(client, cfg.IsCloud)
    issue, err := fetcher.GetIssue(context.Background(), jira.GetIssueOpts{
        IssueKey: os.Getenv("TEST_ISSUE_KEY"),
    })
    require.NoError(t, err)
    require.NotNil(t, issue)
}
```

---

## 7. Сводный план-график

| Phase | Содержание | Длительность |
| --- | --- | --- |
| 1 | Foundation (config, logging, lifecycle) | 1–2 дня |
| 2 | HTTP client + auth + middleware | 3–4 дня |
| 3 | Models | 2–3 дня |
| 4 | Jira mixins (21) | 5–7 дней |
| 5 | Confluence mixins (11) | 3–4 дня |
| 6 | Jira read tools (35) | 3–4 дня |
| 7 | Jira write tools (28) | 3–4 дня |
| 8 | Confluence tools (35) | 4–5 дней |
| 9 | Preprocessing | 5–7 дней |
| 10 | Server lifecycle + OAuth | 3–4 дня |
| 11 | Тесты + документация | 3–4 дня |
| 12 | Deploy + hardening | 2–3 дня |
| **Итого** | | **~ 8–10 недель** для одного разработчика |

---

## 8. Риски и открытые вопросы

### Риски

1. **mcp-go API может меняться** — pin to `v0.20.0` (или актуальная) и обновлять по мере стабилизации.
2. **`html-to-markdown` для MD→HTML направления** — возможно придётся писать свой конвертер на базе `goldmark` или `blackfriday`.
3. **OAuth-через-keyring** на Linux-серверах без `gnome-keyring` может ломаться — fallback на0600 файл обязателен.
4. **ASGI middleware parity** — Python-версия делает специальные трюки для ASGI disconnect handling; в Go через `http.Server` это проще.

### Открытые вопросы

1. **Single binary vs separate Jira/Confluence binaries**?
2. **OAuth proxy / DCR** — отдельным сервисом или в том же бинаре?
3. **PAC/WPAD** — поддерживать или выкинуть?
4. **Legacy Jira wiki → MD** — нужен или только MD → wiki?
5. **SSL_VERIFY=false** — оставить как escape-hatch или запретить?

---

## 9. Соответствие контракту Python-версии

| Аспект | Python | Go (наш план) | Соответствие |
| --- | --- | --- | --- |
| Имена инструментов | `jira_*`, `confluence_*` | Те же | ✅ |
| JSON-schema параметров | Pydantic + Annotated | mcp-go `With*` helpers | ✅ Семантически |
| Return shape (text) | `json.dumps(..., indent=2)` | `json.MarshalIndent` без `SetEscapeHTML(false)` | ✅ |
| Return shape (binary) | `list[TextContent\|EmbeddedResource]` | `[]mcp.Content` + `NewToolResultImage`/`NewToolResultBlob` | ✅ |
| Tags + annotations | `tags={…}`, `annotations={…}` | Те же через ToolDef | ✅ |
| Read-only mode | `READ_ONLY_MODE=true` | То же, проверяется в middleware | ✅ |
| Toolset filtering | `TOOLSETS=…` | То же, проверяется в middleware | ✅ |
| `ENABLED_TOOLS` | Comma-separated list | То же | ✅ |
| Cloud OAuth | `https://api.atlassian.com/ex/jira/{cloudId}/…` | То же | ✅ |
| Per-request credentials | Headers `X-Atlassian-*` | Те же | ✅ |
| Multi-tenant | Один бинарь, разные fetchers | То же | ✅ |
| `JIRA_INTERNAL_ONLY_PROJECTS` | Comments guard | Те же проверки | ✅ |
| Tiny-link resolution | `_resolve_page_id` | То же | ✅ |
| 50 MiB attachment cap | `ATTACHMENT_MAX_BYTES` | То же | ✅ |
| `pagingToken` (Cloud) | `nextPageToken` | То же | ✅ |

---

## 10. Начало работы

1. Клонировать репо.
2. `cd docs/ai && ls` — ознакомиться со структурой.
3. Начать с `docs/ai/INDEX.md` → `00-overview/tool-catalog.md` для понимания всех инструментов.
4. Phase 1 → Phase 12 по графику выше.
5. Каждый phase заканчивается работающим `go test ./...`.

После завершения phase 12 у вас будет функциональный Go-порт
`mcp-atlassian`, совместимый по MCP-контракту с Python-версией.