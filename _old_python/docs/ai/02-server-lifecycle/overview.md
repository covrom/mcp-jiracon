# Server lifecycle, lifespan, mounting, tool filtering

This document describes how the FastMCP servers are constructed, mounted,
filtered, and torn down. The Go port needs to preserve this lifecycle
behavior exactly.

## Construction

There are three FastMCP instances, all subclasses of `ErrorPreservingFastMCP`
(`servers/error_handling.py`):

| Instance | Variable | File | Tools | Namespace |
| --- | --- | --- | --- | --- |
| Root server | `main_mcp` | `servers/main.py:885` | (none directly; hosts mounted sub-servers and `/healthz`) | – |
| Jira sub-server | `jira_mcp` | `servers/jira.py:41` | 63 | `jira` |
| Confluence sub-server | `confluence_mcp` | `servers/confluence.py:182` | 35 | `confluence` |

```python
jira_mcp       = ErrorPreservingFastMCP(name="Jira MCP Service",       instructions="...")
confluence_mcp = ErrorPreservingFastMCP(name="Confluence MCP Service", instructions="...")
main_mcp       = AtlassianMCP(name="Atlassian MCP", lifespan=main_lifespan, auth=...)
```

`AtlassianMCP` is a `FastMCP` subclass that overrides the listing and call
hooks so it can apply the per-request tool filter.

## Mounting

```python
main_mcp.mount(jira_mcp,       namespace="jira")
main_mcp.mount(confluence_mcp, namespace="confluence")
```

When a sub-server is mounted under a namespace, every tool it exposes
gains that namespace prefix. So `jira_mcp` tool `get_issue` becomes
**`jira_get_issue`** in the global tool list.

## Lifespan

`main_lifespan(app)` is an `@asynccontextmanager`. It runs once on server
start and once on shutdown. During start it:

1. Reads environment variables (`get_available_services`, `is_read_only_mode`,
   `get_enabled_tools`, `get_enabled_toolsets`).
2. Tries to build a `JiraConfig` and a `ConfluenceConfig` from env. Each
   is loaded only if `is_auth_configured()` is true; otherwise logged
   warnings are emitted and the corresponding service becomes unavailable.
3. Builds a `MainAppContext` dataclass and yields it as
   `{"app_lifespan_context": MainAppContext(...)}`.

Tools retrieve the context with:

```python
ctx.request_context.lifespan_context["app_lifespan_context"]
```

`MainAppContext` is defined in `servers/context.py`:

```python
@dataclass(frozen=True)
class MainAppContext:
    full_jira_config: JiraConfig | None = None
    full_confluence_config: ConfluenceConfig | None = None
    read_only: bool = False
    enabled_tools: list[str] | None = None
    enabled_toolsets: set[str] | None = None
```

## Tool filtering

Tool filtering has two layers and runs at two times:

### Listing filter (`_list_tools_mcp`)

When a client calls `tools/list`, the main server:

1. Pulls `ctx` info from the request (read-only flag, enabled tools list,
   enabled toolsets, per-request header-based service availability).
2. Walks every registered tool.
3. Drops tools where:
   - **Service unavailable**: tag has `jira`/`confluence` but the service
     isn't configured.
   - **Toolset disabled**: tool's `toolset:*` tag is not in the enabled
     toolsets.
   - **`ENABLED_TOOLS` denies**: tool name not in `enabled_tools`.
   - **Read-only mode**: tool has a `write` tag while `read_only=True`.
4. Returns the surviving list. Tool names are sanitized
   (`_sanitize_schema_for_compatibility`) before being returned.

### Call filter (`_call_tool_mcp`)

When a client calls `tools/call {name, arguments}`:

1. Same context lookup as listing.
2. If the named tool isn't authorized, raise `NotFoundError(f"Unknown tool: {key}")`.
3. Otherwise, delegate to FastMCP's default `_call_tool_mcp`.

The `_is_tool_authorized` predicate is identical for listing and call,
except for the service-availability check (which is listing-only). This
ensures that a tool hidden from the listing cannot be invoked directly by
name.

## Per-request credential resolution

`get_jira_fetcher(ctx)` / `get_confluence_fetcher(ctx)` resolve a
credential-aware fetcher for the current request. The logic lives in
`servers/dependencies.py` and is documented in detail in
[03-auth/per-request-resolution.md](../03-auth/per-request-resolution.md).
The high-level flow:

```
1. Look up cached fetcher in request.state
   ├─ hit: return it
   └─ miss: continue

2. Try header-PAT auth
   ├─ X-Atlassian-Jira-Url + X-Atlassian-Jira-Personal-Token
   └─ create per-URL fetcher, attach SSRF-safe hook, validate, cache

3. Try Basic auth
   ├─ Authorization: Basic base64(email:api_token)
   └─ create user-specific fetcher, validate, cache

4. Try Bearer (OAuth or PAT)
   ├─ Authorization: Bearer <token> (+ optional X-Atlassian-Cloud-Id)
   └─ disambiguate OAuth vs PAT based on global config

5. Fall back to global fetcher
   ├─ stdio mode: always allowed
   └─ HTTP mode: only if ALLOW_GLOBAL_CRED_FALLBACK=true or auth_type=="external"
```

## Error wrapping

`ErrorPreservingFastMCP` overrides `tool(...)` so that every registered tool
is automatically wrapped with `handle_tool_errors`. The wrapper:

- Lets `ToolError` (FastMCP's own exception) propagate unchanged.
- Catches any other exception, logs the traceback, and raises
  `ToolError(f"Error calling tool '<name>': <detail>")` with `from e`.

`check_write_access` is an additional decorator applied per-tool by the
author of write tools. It raises `ValueError("Cannot <verb> in read-only
mode.")` if read-only mode is on.

## ASGI middleware (`UserTokenMiddleware`)

For HTTP transport, a Starlette/ASGI middleware extracts per-request
credentials from the inbound headers and stores them on
`scope["state"]`:

| Header | Purpose |
| --- | --- |
| `Authorization: Basic …` | Decode to email + API token; `state["user_atlassian_auth_type"] = "basic"`. |
| `Authorization: Bearer …` | OAuth bearer; `state["user_atlassian_auth_type"] = "oauth"`. |
| `Authorization: Token …` | Server/DC PAT; `state["user_atlassian_auth_type"] = "pat"`. |
| `X-Atlassian-Cloud-Id` | Optional tenant override for Cloud OAuth. |
| `X-Atlassian-Jira-Url`, `X-Atlassian-Jira-Personal-Token` | Per-request PAT for Jira. |
| `X-Atlassian-Confluence-Url`, `X-Atlassian-Confluence-Personal-Token` | Per-request PAT for Confluence. |

URLs are validated against SSRF (`validate_url_for_ssrf`) before use.

If the request arrives with no auth header and `ALLOW_GLOBAL_CRED_FALLBACK`
is unset, the request is rejected with `401`.

## OAuth provider (`_build_auth_provider`)

When `ATLASSIAN_OAUTH_PROXY_ENABLE=true`, the root server constructs a
`HardenedOAuthProxy` instance. This is the FastMCP OAuth proxy with DCR
controls:

- `allowed_grant_types` filters incoming `grant_types` (default
  `authorization_code`, `refresh_token`).
- `forced_scopes` overrides the `scope` field on every registration.
- `response_types` is forced to `["code"]`.
- `client_storage` (optional) is a custom async KV backend, loaded via
  `build_oauth_client_storage_from_env` (factory mode).

The proxy's authorize / token endpoints are served alongside the MCP
endpoint. Cloud and Data Center are auto-detected.

## Transport

`main_mcp.http_app(...)` returns a Starlette ASGI app. The current default
is `streamable-http`. Configuration knobs:

| Env var | Effect |
| --- | --- |
| `HOST` | Bind address (default `0.0.0.0`). |
| `PORT` | Listen port (default `3000`). |
| `ATLASSIAN_MCP_PATH` | URL path for the MCP endpoint (default `/mcp`). |
| `IGNORE_HEADER_AUTH` | Skip auth header processing (for GCP / AWS load balancers). |
| `ALLOW_GLOBAL_CRED_FALLBACK` | Allow operator's global credentials to back unauthenticated HTTP requests. |

Stdio is launched by the entry-point script
(`uv run mcp-atlassian`) which calls
`mcp.run(transport="stdio")` (default for the `mcp-atlassian` script).

## Health check

`GET /healthz` returns `{"status": "ok"}`. Implemented as
`@main_mcp.custom_route("/healthz", methods=["GET"])`. Used by Kubernetes
probes.

## Teardown

On shutdown (`finally:` block in `main_lifespan`), the lifespan simply
logs that cleanup is happening. There's no explicit resource release
because Atlassian sessions are lazily built per request and have no
persistent connections.

## Go reimplementation map

| Python concept | Go equivalent |
| --- | --- |
| `AtlassianMCP(FastMCP)` | A struct holding the root MCP server, with method overrides for listing/call filters. |
| `main_mcp.mount(sub, namespace=...)` | A prefix-appender at tool registration time (Go has no built-in mount). |
| `main_lifespan` asynccontextmanager | A constructor + `Shutdown(ctx)` method on the app struct. |
| `UserTokenMiddleware` | A `func(http.Handler) http.Handler` wrapper. |
| `ErrorPreservingFastMCP.tool` | A `RegisterTool` method that auto-wraps handlers with a recover/retry helper. |
| `HardenedOAuthProxy` | Use the Go MCP SDK's OAuth provider; apply grant-type / scope filters in the registration handler. |