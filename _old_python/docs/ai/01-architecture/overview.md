# Architecture overview

This document describes the layered architecture that `mcp-atlassian` actually
implements, end-to-end, so a Go reimplementation can decide which layers to
collapse and which to keep.

## Layered diagram

```
┌────────────────────────────────────────────────────────────────────────────┐
│ Layer 0 — MCP transport                                                    │
│   • Streamable HTTP (FastMCP) / stdio                                       │
│   • Inbound: JSON-RPC over HTTP/POST or stdio lines                        │
│   • Outbound: same, plus SSE for streaming responses                       │
├────────────────────────────────────────────────────────────────────────────┤
│ Layer 1 — Main FastMCP server (servers/main.py)                            │
│   • AtlassianMCP subclass of FastMCP                                        │
│   • Mounts jira_mcp / confluence_mcp sub-servers under namespaces          │
│   • Applies the tool filter (read-only, ENABLED_TOOLS, TOOLSETS,           │
│     service availability)                                                  │
│   • Hosts the UserTokenMiddleware for per-request credential extraction    │
│   • Optionally hosts an OAuth proxy / DCR provider                         │
├────────────────────────────────────────────────────────────────────────────┤
│ Layer 2 — Sub-servers (servers/jira.py, servers/confluence.py)             │
│   • ErrorPreservingFastMCP instances named "Jira MCP Service" and          │
│     "Confluence MCP Service"                                               │
│   • Each tool is an async function decorated with @<svc>_mcp.tool(...)     │
│   • Tools call into a fetcher (jira or confluence) via the request         │
│     context's lifespan state                                               │
├────────────────────────────────────────────────────────────────────────────┤
│ Layer 3 — Dependencies (servers/dependencies.py)                           │
│   • `get_jira_fetcher(ctx)` / `get_confluence_fetcher(ctx)` resolve to a   │
│     credential-aware JiraFetcher / ConfluenceFetcher                       │
│   • Four auth branches: header PAT, basic auth, OAuth/PAT bearer, and    │
│     operator-global fallback                                               │
│   • Cross-request validation cache (SHA-256 of credential × scope)        │
│   • SSRF-safe redirect hook attached to per-request sessions               │
├────────────────────────────────────────────────────────────────────────────┤
│ Layer 4 — Fetchers (jira/__init__.py, confluence/__init__.py)              │
│   • JiraFetcher composes 21 mixins (IssuesMixin, SearchMixin, …)          │
│   • ConfluenceFetcher composes 11 mixins (PagesMixin, CommentsMixin, …)    │
│   • Mixins expose high-level methods that drive one or more REST calls    │
├────────────────────────────────────────────────────────────────────────────┤
│ Layer 5 — Client wrappers (jira/client.py, confluence/client.py)          │
│   • Wrap atlassian-python-api's Jira / Confluence                          │
│   • Inject auth headers (Basic / Bearer / PAT)                              │
│   • Mount the session middleware stack (SSL/NO_PROXY, proxy, retry,       │
│     rate-limit, concurrency, circuit-breaker, SSRF pin)                    │
├────────────────────────────────────────────────────────────────────────────┤
│ Layer 6 — Atlassian REST API                                               │
│   • Jira: /rest/api/2 (Server/DC) or /rest/api/3 (Cloud), plus agile     │
│     and Service Desk endpoints                                             │
│   • Confluence: /wiki/rest/api (v1) or /wiki/api/v2 (v2, Cloud OAuth)     │
│   • Authentication: Basic, Bearer (PAT or OAuth2 access token)            │
└────────────────────────────────────────────────────────────────────────────┘
```

## Module dependencies (what imports what)

```
servers/main.py
  ├── servers/jira.py              (jira_mcp)
  ├── servers/confluence.py        (confluence_mcp)
  ├── servers/oauth_proxy.py       (HardenedOAuthProxy)
  ├── servers/client_storage.py    (build_oauth_client_storage_from_env)
  ├── servers/context.py           (MainAppContext)
  ├── servers/error_handling.py    (ErrorPreservingFastMCP)
  └── servers/dependencies.py      (get_jira_fetcher / get_confluence_fetcher)

servers/jira.py
  ├── servers/dependencies.py
  ├── servers/async_utils.py       (run_jira_fetcher_call)
  ├── servers/error_handling.py
  ├── jira/constants.py            (DEFAULT_READ_JIRA_FIELDS)
  ├── jira/forms_common.py         (convert_datetime_to_timestamp)
  ├── models/jira/                 (JiraAttachment, JiraUser)
  ├── utils/decorators.py          (check_write_access)
  ├── utils/media.py               (ATTACHMENT_MAX_BYTES, fetch_and_encode_attachment, is_image_attachment)
  └── exceptions.py                (MCPAtlassianAuthenticationError)

servers/confluence.py
  ├── servers/dependencies.py
  ├── servers/error_handling.py
  ├── models/confluence/           (ConfluenceAttachment)
  ├── utils/decorators.py          (check_write_access)
  ├── utils/io.py                  (validate_safe_path)
  ├── utils/media.py
  └── exceptions.py

servers/dependencies.py
  ├── jira/{config,__init__}       (JiraConfig, JiraFetcher)
  ├── confluence/{config,__init__} (ConfluenceConfig, ConfluenceFetcher)
  ├── servers/context.py           (MainAppContext)
  ├── utils/env.py                 (get_header_names, is_env_ssl_verify, is_env_truthy)
  ├── utils/oauth.py               (OAuthConfig)
  ├── utils/proxy.py               (get_proxy_settings_from_env)
  └── utils/urls.py                (validate_url_for_ssrf)

jira/__init__.py
  └── jira/{mixin modules}.py      (each mixin composes Jira methods)
       └── jira/client.py          (JiraClient)
            ├── jira/config.py     (JiraConfig)
            └── atlassian-python-api (Jira)

confluence/__init__.py
  └── confluence/{mixin modules}.py
       └── confluence/client.py    (ConfluenceClient)
            ├── confluence/config.py (ConfluenceConfig)
            ├── confluence/v2_adapter.py (ConfluenceV2Adapter)
            └── atlassian-python-api (Confluence)

preprocessing/{base,confluence,jira}.py
  └── Used by Jira/Confluence model conversions and tool handlers

models/{jira,confluence}/*.py
  └── JiraIssue, ConfluencePage, … with .to_simplified_dict() / .from_api_response()
```

## Mixin composition tables

### `JiraFetcher`

```python
class JiraFetcher(
    ProjectsMixin,
    FieldsMixin,
    FieldOptionsMixin,
    FormsApiMixin,
    FormattingMixin,
    TransitionsMixin,
    WorklogMixin,
    EpicsMixin,
    CommentsMixin,
    CustomerRequestsMixin,
    SearchMixin,
    IssuesMixin,
    UsersMixin,
    WatchersMixin,
    BoardsMixin,
    SprintsMixin,
    QueuesMixin,
    AttachmentsMixin,
    LinksMixin,
    MetricsMixin,
    SLAMixin,
    DevelopmentMixin,
    ProjectAnalysisMixin,
):
    pass
```

### `ConfluenceFetcher`

```python
class ConfluenceFetcher(
    SearchMixin,
    SpacesMixin,
    PagesMixin,
    CommentsMixin,
    LabelsMixin,
    UsersMixin,
    AnalyticsMixin,
    AttachmentsMixin,
    TemplatesMixin,
    PermissionsMixin,
    RestrictionsMixin,
):
    pass
```

## What the layers are responsible for

### Layer 0 — MCP transport

`fastmcp` (the underlying library) implements the JSON-RPC framing. Most
production deployments use **streamable HTTP** with optional SSE for
streaming responses. Stdio mode works for single-user CLI integrations. In
HTTP mode, requests arrive on the URL configured via
`ATLASSIAN_MCP_PATH` (defaults to FastMCP's `streamable_http_path`).

### Layer 1 — Main server

Responsibilities:
- Read lifespan environment variables; build `JiraConfig` /
  `ConfluenceConfig`.
- Mount the two sub-servers under namespaces `jira` and `confluence`.
- Run a per-request tool filter:
  - **Listing filter**: hides tools if the service is unavailable, the
    toolset isn't enabled, or read-only mode disables write tools.
  - **Call filter**: enforces the same at call time (a hidden tool cannot
    be invoked directly).
- Host the OAuth provider if `ATLASSIAN_OAUTH_PROXY_ENABLE=true`.

### Layer 2 — Sub-servers

Each sub-server is its own `FastMCP` instance. Tools are decorated with
`@<svc>_mcp.tool(tags={...}, annotations={...})`. Tags drive the
toolset/read-only filter; annotations surface in the MCP `listTools`
response (used by clients to render UI affordances).

A typical tool is a thin async wrapper:

```python
@jira_mcp.tool(tags={"jira", "read", "toolset:jira_issues"},
               annotations={"title": "Get Issue", "readOnlyHint": True})
async def get_issue(ctx, issue_key, fields=DEFAULT_FIELDS, …):
    jira = await get_jira_fetcher(ctx)
    issue = await run_jira_fetcher_call(
        jira.get_issue, issue_key=issue_key, fields=fields, …
    )
    return json.dumps(issue.to_simplified_dict(), indent=2, ensure_ascii=False)
```

### Layer 3 — Dependencies (`servers/dependencies.py`)

This is the heart of multi-tenant credential resolution. Per request:

1. Look up the request state for an existing fetcher (cached by
   `request.state`).
2. If absent, walk the auth branches in order:
   - **header PAT**: `X-Atlassian-Jira-Url` + `X-Atlassian-Jira-Personal-Token`
   - **basic auth**: `Authorization: Basic base64(email:api_token)`
   - **bearer (OAuth or PAT)**: `Authorization: Bearer <token>` plus optional
     `X-Atlassian-Cloud-Id`
3. Cache the resulting fetcher on `request.state` for the lifetime of the
   request.
4. Validation (network round-trip to Atlassian) is deduplicated via a
   SHA-256 credential cache, TTL-controlled by
   `MCP_ATLASSIAN_VALIDATION_CACHE_TTL`.

### Layer 4 — Fetchers

The fetcher classes expose a flat API for tools, composed from mixins.
The mixin tree is a way to organize ~30 kLOC of Jira-specific and
~6 kLOC of Confluence-specific code without inheritance tangles; the Go
port doesn't need to mirror the mixin tree and can use plain methods.

### Layer 5 — Client wrappers

The `JiraClient` / `ConfluenceClient` thin-wrap `atlassian-python-api`'s
`Jira` / `Confluence` classes. They add:

- Auth-header injection (based on `JiraConfig.auth_type` /
  `ConfluenceConfig.auth_type`).
- SSL verification / NO_PROXY / mTLS / SSL-skip middleware.
- Proxy / PAC / WPAD configuration.
- Retry / rate-limit / concurrency / circuit-breaker middleware
  (`utils/http.py`).
- SSRF DNS pinning (`utils/ssrf_adapter.py`).

In Go this entire layer collapses into a `*http.Client` with a custom
`Transport` chain.

### Layer 6 — Atlassian REST API

The actual HTTP calls. See
[01-architecture/jira-rest-api.md](jira-rest-api.md) and
[01-architecture/confluence-rest-api.md](confluence-rest-api.md) for
endpoint-by-endpoint coverage.

## Async / concurrency model

The tool functions are `async def`. Inside, blocking `atlassian-python-api`
calls are offloaded to a worker thread with a configurable cap
(`JIRA_FETCHER_MAX_WORKERS`, default 8) via `run_jira_fetcher_call`. In
Go, the equivalent is just a normal function call — Go's runtime handles
concurrency naturally.

## Where the Go port can collapse layers

| Python layer | Go equivalent |
| --- | --- |
| Layer 1 (`AtlassianMCP`) | A struct that registers tools and applies the tool filter |
| Layer 2 (`jira_mcp`, `confluence_mcp`) | A `Server` per service in the same Go binary |
| Layer 3 (`get_jira_fetcher`) | A single `func (s *Server) jiraFetcher(ctx) (*JiraClient, error)` |
| Layer 4 (mixins) | A single `*JiraFetcher` struct holding all methods |
| Layer 5 (`JiraClient`) | A `*http.Client` + the `*atlassian.Jira` shim (or replace entirely) |
| Layer 6 (REST) | The same — Atlassian doesn't care what language the client is in |

In other words, you don't need the mixin tree or the dependency-injection
machinery; you do need the per-request credential resolution and the
schema-faithful tool layer.