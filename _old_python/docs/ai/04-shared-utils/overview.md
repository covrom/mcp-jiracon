# Shared utilities

The `utils/` package contains every cross-cutting helper. The Go port
needs equivalents for each. The following index summarizes what each
module does and what the Go equivalent should look like.

| Module | Purpose | Public API | Go equivalent |
| --- | --- | --- | --- |
| `utils/__init__.py` | Re-export surface | `parse_date`, `is_read_only_mode`, `validate_safe_path`, `setup_logging`, `mask_sensitive`, `ATTACHMENT_MAX_BYTES`, `fetch_and_encode_attachment`, `is_image_attachment`, `OAuthConfig`, `NoProxyAdapter`, `SSLIgnoreAdapter`, `is_atlassian_cloud_url`, `validate_url_for_ssrf`, … | A package `utils` re-exporting the same symbols. |
| `utils/env.py` | Env var parsing (truthy/int/float, custom headers, header-name lists) | `is_env_truthy`, `is_env_extended_truthy`, `is_env_ssl_verify`, `get_int_env`, `get_float_env`, `get_custom_headers`, `get_header_names` | `os.LookupEnv` + `strings.EqualFold`. Build a `Config` package with typed accessors. |
| `utils/environment.py` | Service-availability decision per env + headers | `get_available_services(headers: dict \| None)` | A function that walks the precedence waterfall and returns `{"jira": bool, "confluence": bool}`. |
| `utils/io.py` | Read-only mode + safe-path | `is_read_only_mode()`, `validate_safe_path(path, base_dir=None)` | `os.Getenv("READ_ONLY_MODE")` + `filepath.Clean`/`EvalSymlinks`/`IsLocal`. |
| `utils/logging.py` | Logger setup + secret masking | `setup_logging`, `mask_sensitive`, `get_masked_session_headers`, `log_config_param` | `slog.NewTextHandler` + `strings.Repeat("*", n)`. |
| `utils/oauth.py` | OAuth 2.0 logic (3LO + BYO) for Cloud and DC | `OAuthConfig`, `BYOAccessTokenOAuthConfig`, `get_oauth_config_from_env`, `configure_oauth_session` | A struct + methods. Use `golang.org/x/oauth2` for the cloud flow, hand-write the DC form-encoded token exchange. |
| `utils/oauth_setup.py` | Interactive OAuth setup wizard | CLI entry point | CLI program using `net/http` to capture the redirect. |
| `utils/urls.py` | URL helpers, SSRF validation | `is_atlassian_cloud_url`, `validate_url_for_ssrf`, `resolve_relative_url`, `make_ssrf_redirect_hook` | `net.ParseIP` + `IsPrivate()` family; `net.Resolver.LookupIPAddr`; custom `CheckRedirect`. |
| `utils/tools.py` | ENABLED_TOOLS allowlist | `get_enabled_tools`, `should_include_tool` | `strings.Split` + exact-match check. |
| `utils/toolsets.py` | Toolset groups + TOOLSETS env | `JIRA_TOOLSETS`, `CONFLUENCE_TOOLSETS`, `ALL_TOOLSETS`, `get_enabled_toolsets`, `should_include_tool_by_toolset`, `get_toolset_tag` | Mirror the dict structure; a `ResolveToolsets(env string) ([]string, error)` helper. |
| `utils/decorators.py` | `@handle_tool_errors`, `@check_write_access`, `@handle_auth_errors`, `@handle_atlassian_api_errors` | Wraps async tools with error handling and read-only guard; wraps sync client methods with auth/error classification. | Explicit middleware functions or method decorators (Go has no decorators). |
| `utils/proxy.py` | Proxy / WPAD / PAC configuration | `ProxyConfigProtocol`, `get_proxy_settings_from_env`, `get_explicit_proxy_map`, `apply_proxy_configuration` | `http.ProxyFromEnvironment` + custom round-tripper for PAC. |
| `utils/token_verifier.py` | FastMCP `TokenVerifier` for opaque Atlassian tokens | `AtlassianOpaqueTokenVerifier` | Interface from the Go MCP SDK; returns a 30-day-expiry `AccessToken`. |
| `utils/media.py` | MIME detection + attachment fetch + base64 encode | `ATTACHMENT_MAX_BYTES`, `is_image_attachment`, `fetch_and_encode_attachment` | `mime.TypeByExtension` + `encoding/base64`. |
| `utils/ssl.py` | SSL and proxy-bypass adapters | `NoProxyAdapter`, `SSLIgnoreAdapter`, `configure_proxy_bypass`, `configure_ssl_verification` | `http.Transport` with custom `RoundTripper` for proxy bypass + `tls.Config{InsecureSkipVerify: true}` for SSL ignore. |
| `utils/ssrf_adapter.py` | DNS-pinning SSRF adapter | `SsrfPinningAdapter`, `mount_ssrf_pinning` | Custom `DialContext` that resolves once and rejects non-global IPs for untrusted hosts. |
| `utils/date.py` | Mixed-format date parser | `parse_date(date_str)` | `time.Parse(time.RFC3339, s)` + `time.UnixMilli(n)`. |
| `utils/lifecycle.py` | Signal handlers + clean exit | `setup_signal_handlers`, `ensure_clean_exit` | `signal.Notify` channel + `os.Stdout.Sync()`. |
| `utils/user_agent.py` | Versioned User-Agent | `get_default_user_agent()` | Build-time `-ldflags`. |
| `utils/pagination.py` | `limit` clamp | `clamp_limit(requested, context)` | Trivial port. |
| `utils/http.py` | Retry / rate-limit / concurrency / circuit-breaker middleware | `configure_retry`, `configure_concurrency`, `configure_rate_limit`, `configure_circuit_breaker`, `CircuitBreakerOpenError`, `format_rate_limit_error` | Custom `http.RoundTripper` chain. Libraries: `hashicorp/go-retryablehttp`, `golang.org/x/sync/semaphore`, `golang.org/x/time/rate`, `sony/gobreaker`. |

## Decorator error handling — detailed contract

These four decorators are the canonical error-classification surface:

```python
@handle_tool_errors                   # async, applied automatically by ErrorPreservingFastMCP
@check_write_access                   # async, applied per write-tool
@handle_auth_errors(service_name=...) # sync, applied to Jira/Confluence client methods
@handle_atlassian_api_errors(service_name=...)  # sync, applied to Jira/Confluence client methods
```

### `@handle_tool_errors`

Wraps an `async def` tool. Behavior:

- Lets `ToolError` propagate (FastMCP's own exception).
- Catches any other exception. Logs the traceback at ERROR level.
- Raises `ToolError(f"Error calling tool '<name>': <detail>")` chained via
  `from e`. The MCP client receives this as a normal tool error.

### `@check_write_access`

Wraps an `async def` tool (must be a `@tool`). Behavior:

- Pulls `app_lifespan_context` from `ctx.request_context.lifespan_context`.
- If `app_lifespan_ctx.read_only == True` → raises
  `ValueError("Cannot <verb> in read-only mode.")`.
- Otherwise delegates to the wrapped function.

The first arg of the wrapped function must be `ctx: Context`. The
verb in the error message is derived from the wrapped function's name:

```
create_issue     -> "Cannot create in read-only mode."
update_issue     -> "Cannot update in read-only mode."
delete_issue     -> "Cannot delete in read-only mode."
...
```

### `@handle_auth_errors(service_name="Atlassian API")`

Wraps a sync function. Behavior:

- `requests.HTTPError` with response status `401` or `403` →
  `MCPAtlassianAuthenticationError(f"{service_name} authentication failed: HTTP {status}")`.
- Other HTTP errors propagate.
- Other exceptions propagate.

### `@handle_atlassian_api_errors(service_name="Atlassian API")`

Wraps a sync function. Behavior:

- `requests.HTTPError` with response status `401` or `403` →
  `MCPAtlassianAuthenticationError`.
- `KeyError` → `ValueError(f"Error processing {service_name} results: missing key {e}")`.
- `requests.RequestException` → `ValueError(f"Network error during {op}: {e}")`.
- `ValueError`/`TypeError` → `ValueError(f"Error processing {service_name} results: {e}")`.
- Generic `Exception` → `RuntimeError(f"Unexpected error during {op}: {e}")` with full traceback at debug level.

## SSRF — exact rules

`validate_url_for_ssrf(url)` returns `None` on safe; returns a string
error message on unsafe. Rejection rules, in order:

1. Empty string → rejected.
2. Scheme not in `{http, https}` → rejected.
3. Backslash in netloc → rejected (bypasses `requests` URL parser).
4. Hostname in the blocklist (`localhost`, `metadata.google.internal`,
   …) → rejected.
5. If hostname is an IP literal:
   - Private (`10/8`, `172.16/12`, `192.168/16`, `fc00::/7`, …) → rejected.
   - Loopback (`127/8`, `::1`) → rejected.
   - Link-local (`169.254/16`, `fe80::/10`) → rejected.
   - IPv4-mapped IPv6 (`::ffff:0:0/96`) is unwrapped and the embedded IPv4
     is checked.
6. If `MCP_ALLOWED_URL_DOMAINS` is set, the hostname must match (exact or
   `.<domain>` suffix). When matched, skip DNS check.
7. Else resolve via `socket.getaddrinfo`; reject if **any** resolved
   address is non-global.

The DNS-pinning adapter (`utils/ssrf_adapter.py`) provides a defense
against the TOCTOU between middleware validation and the actual connect
call. It:

1. Resolves `host` via `getaddrinfo`.
2. For each candidate address, validates against the SSRF rules.
3. Connects to that exact address — no re-resolution.

Operator-trusted hosts (`JIRA_URL`, `CONFLUENCE_URL`, and
`MCP_ALLOWED_URL_DOMAINS`) are exempt from the non-global check.

## Toolset definitions (verbatim)

| Toolset | Default | Description |
| --- | --- | --- |
| `jira_issues` | yes | Core CRUD on issues |
| `jira_fields` | yes | Field metadata and option lookups |
| `jira_comments` | yes | Comments |
| `jira_transitions` | yes | Status transitions |
| `jira_projects` | no | Projects / versions / components |
| `jira_agile` | no | Agile boards, sprints |
| `jira_links` | no | Issue links (incl. Epic link, remote link) |
| `jira_worklog` | no | Worklogs |
| `jira_attachments` | no | Attachments and images |
| `jira_users` | no | User search and profile |
| `jira_watchers` | no | Issue watchers |
| `jira_service_desk` | no | Service Desk queues and requests (DC mostly) |
| `jira_forms` | no | ProForma forms (Cloud only) |
| `jira_metrics` | no | Issue dates + SLA metrics |
| `jira_development` | no | Dev panel (PRs, branches, commits) |
| `jira_project_analysis` | no | Epic hierarchy + cross-project links |
| `confluence_pages` | yes | Pages CRUD, history, diff, sections |
| `confluence_comments` | yes | Footer + inline comments |
| `confluence_labels` | no | Page labels |
| `confluence_users` | no | User search |
| `confluence_analytics` | no | Page view analytics (Cloud only) |
| `confluence_attachments` | no | Attachments + images |
| `confluence_templates` | no | Page templates (Cloud only) |
| `confluence_permissions` | no | Permission check / space permissions (Cloud only) |

`get_enabled_toolsets()` reads `TOOLSETS` (comma-separated). Tokens:

- `all` (case-insensitive) → all 24 toolset names.
- `default` → the 6 `default=True` ones.
- Unknown names → logged warning, dropped.
- Empty / unset → log deprecation warning, return all toolsets. (v0.22.0
  plans to switch default to `DEFAULT_TOOLSETS`.)