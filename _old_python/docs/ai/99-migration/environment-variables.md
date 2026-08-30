# Environment variables reference

Every environment variable the server reads. Use this as the spec for
the Go port's `internal/config` package.

## Service URLs (required)

| Var | Effect |
| --- | --- |
| `JIRA_URL` | Base URL for the Jira instance (Cloud or DC). Required to enable Jira tools. |
| `CONFLUENCE_URL` | Base URL for the Confluence instance. Required to enable Confluence tools. |

## Service-specific auth

### Jira

| Var | Effect |
| --- | --- |
| `JIRA_USERNAME` | Cloud basic auth email |
| `JIRA_API_TOKEN` | Cloud basic auth API token |
| `JIRA_PERSONAL_TOKEN` | DC PAT (Bearer) |
| `JIRA_SSL_VERIFY` | `true` / `false`. Defaults to `true`. |
| `JIRA_PROJECTS_FILTER` | Comma-separated project keys; applied to project-listing tools. |
| `JIRA_PASSTHROUGH_HEADERS` | Comma-separated header names to forward on every request. |
| `JIRA_HTTP_PROXY` / `JIRA_HTTPS_PROXY` / `JIRA_NO_PROXY` / `JIRA_SOCKS_PROXY` | Per-service proxies (override global). |
| `JIRA_PROXY_WPAD_ENABLE` / `JIRA_PROXY_WPAD_URL` | Per-service WPAD/PAC config. |

### Confluence

| Var | Effect |
| --- | --- |
| `CONFLUENCE_USERNAME` | Cloud basic auth email |
| `CONFLUENCE_API_TOKEN` | Cloud basic auth API token |
| `CONFLUENCE_PERSONAL_TOKEN` | DC PAT |
| `CONFLUENCE_SSL_VERIFY` | `true` / `false`. |
| `CONFLUENCE_SPACES_FILTER` | Comma-separated space keys; applied to search results. |
| `CONFLUENCE_PASSTHROUGH_HEADERS` | Comma-separated header names to forward. |
| `CONFLUENCE_HTTP_PROXY` / `CONFLUENCE_HTTPS_PROXY` / `CONFLUENCE_NO_PROXY` / `CONFLUENCE_SOCKS_PROXY` | Per-service proxies. |
| `CONFLUENCE_PROXY_WPAD_ENABLE` / `CONFLUENCE_PROXY_WPAD_URL` | Per-service WPAD/PAC config. |

### mTLS (DC)

| Var | Effect |
| --- | --- |
| `JIRA_CLIENT_CERT` / `CONFLUENCE_CLIENT_CERT` | Path to client cert (PEM) |
| `JIRA_CLIENT_KEY` / `CONFLUENCE_CLIENT_KEY` | Path to client private key. Encrypted keys are refused. |
| `JIRA_CLIENT_KEY_PASSWORD` / `CONFLUENCE_CLIENT_KEY_PASSWORD` | Not allowed (encrypted keys are rejected). |

## OAuth 2.0 (global)

| Var | Effect |
| --- | --- |
| `ATLASSIAN_OAUTH_CLIENT_ID` | OAuth app client ID |
| `ATLASSIAN_OAUTH_CLIENT_SECRET` | OAuth app client secret |
| `ATLASSIAN_OAUTH_REDIRECT_URI` | OAuth redirect URI (default port is `3000`) |
| `ATLASSIAN_OAUTH_SCOPE` | Space-separated scopes |
| `ATLASSIAN_OAUTH_CLOUD_ID` | Cloud tenant ID (required for Cloud OAuth) |
| `ATLASSIAN_OAUTH_ACCESS_TOKEN` | BYO access token (skips full 3LO) |
| `ATLASSIAN_OAUTH_AUDIENCE` | Defaults to `api.atlassian.com` for Cloud |
| `ATLASSIAN_OAUTH_PROXY_ENABLE` | Enables the OAuth proxy / DCR. Defaults to `false`. |
| `ATLASSIAN_OAUTH_ALLOWED_CLIENT_REDIRECT_URIS` | Comma-separated allowlist for DCR redirect URIs |
| `ATLASSIAN_OAUTH_ALLOWED_GRANT_TYPES` | Defaults to `authorization_code,refresh_token` |
| `ATLASSIAN_OAUTH_REQUIRE_CONSENT` | Defaults to `true` |
| `ATLASSIAN_OAUTH_CLIENT_STORAGE_MODE` | `default` or `factory` |
| `ATLASSIAN_OAUTH_CLIENT_STORAGE_FACTORY` | `pkg.mod:callable` |
| `ATLASSIAN_OAUTH_CLIENT_STORAGE_CONFIG_JSON` | JSON config passed to factory |
| `JIRA_OAUTH_*` / `CONFLUENCE_OAUTH_*` | Per-service OAuth overrides |

## Server control

| Var | Effect |
| --- | --- |
| `READ_ONLY_MODE` | `true` / `false`. When `true`, all `write`-tagged tools are blocked. |
| `ENABLED_TOOLS` | Comma-separated tool allowlist. When unset, all tools are available. |
| `TOOLSETS` | Comma-separated toolset allowlist. Tokens: `all`, `default`, or toolset names. |
| `ATLASSIAN_MAX_PAGINATION_LIMIT` | Clamps `limit` arguments to this value (0 = no cap). |
| `JIRA_FETCHER_MAX_WORKERS` | Worker cap for Jira fetcher calls inside the event loop. Default 8. |
| `MCP_ATLASSIAN_VALIDATION_CACHE_TTL` | TTL for the credential validation cache in seconds. Default 300. |
| `MCP_ATLASSIAN_VALIDATION_CACHE_MAXSIZE` | Max entries. Default 100. |
| `MCP_ALLOWED_URL_DOMAINS` | Comma-separated hostnames exempt from SSRF non-global check. |
| `IGNORE_HEADER_AUTH` | Skip auth header processing (for GCP Cloud Run / AWS ALB). |
| `ALLOW_GLOBAL_CRED_FALLBACK` | Allow operator's global creds to back unauthenticated HTTP requests. |
| `HOST` | Bind address. Default `0.0.0.0`. |
| `PORT` | Listen port. Default `3000`. |
| `ATLASSIAN_MCP_PATH` | URL path for MCP endpoint. Default `/mcp`. |
| `PUBLIC_BASE_URL` | Public-facing base URL (used by OAuth proxy). |

## Proxies (global)

| Var | Effect |
| --- | --- |
| `HTTP_PROXY` / `HTTPS_PROXY` / `NO_PROXY` / `SOCKS_PROXY` | Standard proxies. |
| `ATLASSIAN_PROXY_WPAD_ENABLE` | Enable WPAD/PAC discovery. |
| `ATLASSIAN_PROXY_WPAD_URL` | PAC file URL. Default `http://wpad/wpad.dat`. |

## SSL (global)

| Var | Effect |
| --- | --- |
| `SSL_VERIFY` | Global `true` / `false`. Per-service vars override. |

## SSL cert files

| Var | Effect |
| --- | --- |
| `JIRA_CLIENT_CERT` / `CONFLUENCE_CLIENT_CERT` | mTLS client cert |
| `JIRA_CLIENT_KEY` / `CONFLUENCE_CLIENT_KEY` | mTLS private key |
| `JIRA_SSL_CA` / `CONFLUENCE_SSL_CA` | Custom CA bundle |

## Jira-specific

| Var | Effect |
| --- | --- |
| `JIRA_INTERNAL_ONLY_PROJECTS` | Comma-separated project keys. Comments/transitions/links in these projects are blocked from being customer-visible. |

## Jira SLA / metrics

| Var | Effect |
| --- | --- |
| `JIRA_SLA_WORKING_HOURS_ONLY` | `true`/`false`. Use working hours only. |
| `JIRA_SLA_WORKING_HOURS_START` | HH:MM (default `09:00`) |
| `JIRA_SLA_WORKING_HOURS_END` | HH:MM (default `17:00`) |
| `JIRA_SLA_WORKING_DAYS` | Comma-separated day-of-week numbers, 1-7 (default `1,2,3,4,5`) |
| `JIRA_SLA_TIMEZONE` | IANA timezone (default = system) |
| `JIRA_SLA_DEFAULT_METRICS` | Comma-separated metric names (default `cycle_time,time_in_status`) |

## Logging

| Var | Effect |
| --- | --- |
| `MCP_LOG_LEVEL` | Log level (default `WARNING`) |
| `LOG_LEVEL` | Alias of above |
| `MCP_PROXY_DEBUG` | Verbose proxy logging |

## Go port equivalents

The Go port should expose the same env var names verbatim. Operators
have invested in CI pipelines that set these. Changing a name is a
breaking change.