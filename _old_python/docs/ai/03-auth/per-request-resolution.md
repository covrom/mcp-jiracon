# Authentication and per-request credential resolution

This document is the canonical spec for everything related to Atlassian
auth in `mcp-atlassian`. The Go port must reproduce the same precedence,
header parsing, and validation caching.

## Auth modes

The system supports five auth modes; each Jira and Confluence config picks
exactly one via `JiraConfig.auth_type` / `ConfluenceConfig.auth_type`:

| `auth_type` | Source | Header injected by client | Used on |
| --- | --- | --- | --- |
| `basic` | Username + API token | `Authorization: Basic base64(email:api_token)` | Cloud + DC |
| `pat` | Personal access token (Server/DC) | `Authorization: Bearer <pat>` | DC only |
| `oauth` | OAuth 2.0 access token (3LO) | `Authorization: Bearer <oauth_access_token>` | Cloud + DC |
| `cert` | mTLS client cert + key | (no Authorization; client cert at TLS layer) | DC only |
| `external` | No header; auth is injected by an upstream proxy | (none) | DC only |

For Cloud with **PAT** (the newer Atlassian Cloud PAT feature), use `pat`
with `auth_type=pat` and the same `Authorization: Bearer` header.

## Precedence of credential sources (per request)

When an MCP request arrives, the dependency layer walks four branches in
this exact order:

1. **Header-based PAT**: the request carries
   `X-Atlassian-{Jira|Confluence}-Personal-Token` *and*
   `X-Atlassian-{Jira|Confluence}-Url`. This creates a per-URL fetcher
   with the operator's SSL/proxy settings but no inherited credentials.

2. **Basic auth header**: `Authorization: Basic base64(email:api_token)`.
   Creates a fetcher that inherits SSL/proxy settings from the global
   config but uses the user's email + API token.

3. **Bearer auth header**: `Authorization: Bearer <token>` (with optional
   `X-Atlassian-Cloud-Id`). The dependency layer disambiguates OAuth
   vs PAT based on the global config:
   - If `OAuthConfig.cloud_id` is set → OAuth Cloud.
   - If `OAuthConfig.is_data_center` is `True` → OAuth DC.
   - Otherwise → Server/DC PAT.

4. **Operator's global fetcher**: used when no per-request credential was
   found. In HTTP mode this is blocked by default
   (`ALLOW_GLOBAL_CRED_FALLBACK=false`) to avoid letting unauthenticated
   callers transact as the operator. Stdio mode always falls back.

The walk happens in `servers/dependencies.py::_get_fetcher`. The same
function is used for both Jira and Confluence; a `_ServiceSpec` dataclass
parameterizes the per-service details.

## `UserTokenMiddleware` — header parsing

The ASGI middleware (`servers/main.py::UserTokenMiddleware`) extracts
credentials from inbound HTTP headers *before* the dependency layer runs.
It populates `scope["state"]`:

| Header | Effect on state |
| --- | --- |
| `Authorization: Basic <b64>` | `user_atlassian_email`, `user_atlassian_api_token`, `user_atlassian_auth_type="basic"` |
| `Authorization: Bearer <token>` | `user_atlassian_token`, `user_atlassian_auth_type="oauth"` |
| `Authorization: Token <token>` | `user_atlassian_token`, `user_atlassian_auth_type="pat"` |
| `X-Atlassian-Cloud-Id` | `user_atlassian_cloud_id` |
| `X-Atlassian-Jira-Url` + `X-Atlassian-Jira-Personal-Token` | `atlassian_service_headers` dict (consumed by dependency layer) |
| `X-Atlassian-Confluence-Url` + `X-Atlassian-Confluence-Personal-Token` | same |

Non-conforming `Authorization` headers return 401 with a JSON error body.
URLs are validated against SSRF (`validate_url_for_ssrf`) before being
trusted.

## Passthrough headers

Operators can declare a list of header names to forward to Atlassian on
every request. Configured via:

- `JIRA_PASSTHROUGH_HEADERS` (comma-separated header names)
- `CONFLUENCE_PASSTHROUGH_HEADERS`

The middleware reads these headers from the inbound request and the
dependency layer attaches them to the user's fetcher config.

## OAuth flow

The full 3LO flow is implemented in `utils/oauth.py::OAuthConfig`:

1. **Setup wizard** (`utils/oauth_setup.py`):
   - Reads `ATLASSIAN_OAUTH_CLIENT_ID` / `_SECRET` / `_REDIRECT_URI`.
   - Calls `get_authorization_url(state)` which appends
     `audience=api.atlassian.com&prompt=consent` for Cloud.
   - Opens the user's browser to that URL.
   - Runs a local HTTP server on the redirect port to capture the code.
   - Calls `exchange_code_for_tokens(code)` to POST to the token endpoint.
   - For Cloud, calls `_get_cloud_id()` to resolve the tenant ID via
     `https://api.atlassian.com/oauth/token/accessible-resources`.
   - Persists tokens via keyring (preferred) or a `0600` file at
     `~/.mcp-atlassian/oauth-<client_id>.json`.

2. **Runtime** (`utils/oauth.py::configure_oauth_session`):
   - On each request: ensures the token isn't expired
     (`is_token_expired` uses `TOKEN_EXPIRY_MARGIN=300s`).
   - Refreshes if needed (`refresh_access_token`).
   - Sets `Authorization: Bearer <access_token>` on the requests session.

### Token storage

| Backend | Behavior |
| --- | --- |
| OS keyring | Preferred; key format `oauth-<client_id>-cloud-<cloudId>` or `-dc-<sha256[:8]>`. Legacy `oauth-<client_id>` is also written for backward compat. |
| File fallback | `~/.mcp-atlassian/oauth-<client_id>.json` with file mode `0600`, dir mode `0700`. |

### OAuth error response

If the auth header parsing or token retrieval fails, the middleware
returns 401 with a JSON body `{"error": "<message>"}`. Common messages:

- `"Unauthorized: Empty Bearer token"`
- `"Unauthorized: Empty Token (PAT)"`
- `"Unauthorized: Empty Basic auth credentials"`
- `"Unauthorized: Invalid Basic auth encoding"`
- `"Unauthorized: Invalid Basic auth format. Expected 'email:api_token'"`
- `"Unauthorized: Email or API token is empty"`
- `"Unauthorized: Only 'Bearer <OAuthToken>', 'Token <PAT>', or 'Basic <base64(email:api_token)>' types are supported."`
- `"Authentication required: no Atlassian credentials were provided."`

## Validation cache

When the dependency layer creates a fetcher, it issues a validation call
(`get_current_user_account_id` for Jira, `get_current_user_info` for
Confluence) to confirm the credentials work. This validation is expensive
(a network round-trip), so it's cached in a process-wide TTLCache:

- Key: `(service_name, sha256(scope || credential))`.
- TTL: `MCP_ATLASSIAN_VALIDATION_CACHE_TTL` (default 300 s).
- Max size: `MCP_ATLASSIAN_VALIDATION_CACHE_MAXSIZE` (default 100).
- In-flight deduplication: a `_validation_inflight` map ensures that
  concurrent requests for the same cache key share one validation call.

## SSRF protection

Every per-URL fetcher attaches an SSRF-safe redirect hook to its HTTP
session:

- Allowlist: `JIRA_URL`, `CONFLUENCE_URL`, and
  `MCP_ALLOWED_URL_DOMAINS`.
- Redirect URLs are validated via `validate_url_for_ssrf` before being
  followed. DNS resolution rejects any non-global address.
- The hook raises `ValueError("Redirect blocked (SSRF): …")` on violation.

## Atlassian Cloud vs Server / Data Center detection

`utils/urls.py::is_atlassian_cloud_url(url)` is the single source of truth:

| URL pattern | Returns |
| --- | --- |
| `localhost`, `127.x`, `10.x`, `172.16-31.x`, `192.168.x` | `False` (DC, internal) |
| `*.atlassian.net` | `True` |
| `*.jira.com`, `*.jira-dev.com` | `True` |
| `api.atlassian.com` | `True` |
| `*.atlassian.com` | `True` |
| `*.atlassian-us-gov-mod.net`, `*.atlassian-us-gov.net` | `True` (FedRAMP) |
| Anything else | `False` |

The same function is used inside `JiraConfig.is_cloud` /
`ConfluenceConfig.is_cloud` properties, which in turn drive endpoint
selection.

## `JiraConfig` / `ConfluenceConfig` factory

Both configs are `@dataclass` types with a `from_env()` classmethod that
parses the relevant env vars. `is_auth_configured()` returns `True` only
when at least one auth credential is set.

`JiraConfig.from_env()` reads:

- `JIRA_URL`
- `JIRA_USERNAME` + `JIRA_API_TOKEN` (Cloud basic) or `JIRA_PERSONAL_TOKEN`
  (DC PAT)
- `JIRA_SSL_VERIFY`
- `JIRA_PROXY_*` envs
- `JIRA_PROJECTS_FILTER`, `JIRA_PASSTHROUGH_HEADERS`

And additionally, when OAuth is configured globally:

- `ATLASSIAN_OAUTH_CLIENT_ID` / `_SECRET` / `_REDIRECT_URI` / `_SCOPE`
- `ATLASSIAN_OAUTH_CLOUD_ID` (Cloud)
- `ATLASSIAN_OAUTH_ACCESS_TOKEN` (BYO)
- `JIRA_OAUTH_*` overrides

`ConfluenceConfig.from_env()` is analogous.

## OAuth precedence inside `OAuthConfig`

The auth flow resolves in this order:

1. **Full 3LO OAuth** (`OAuthConfig.from_env` returns `OAuthConfig`).
   Used when `ATLASSIAN_OAUTH_CLIENT_ID` + `_SECRET` are present.
2. **BYO access token** (`BYOAccessTokenOAuthConfig.from_env`). Used when
   `ATLASSIAN_OAUTH_ACCESS_TOKEN` is present.
3. **No OAuth config**: returns `None`. The dependency layer falls back
   to PAT/Basic/etc.

## Token storage security

- **Keyring**: OS-protected (Windows Credential Manager, macOS Keychain,
  GNOME Keyring, KWallet).
- **File fallback**: `~/.mcp-atlassian/` directory mode `0700`, JSON file
  mode `0600`.
- **In-memory**: never persisted to logs (`mask_sensitive` strips
  credentials before logging).

## Go port notes

- The `OAuthConfig` dataclass becomes a Go struct with the same fields.
- Keyring access via `github.com/zalando/go-keyring`. File fallback via
  `os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)`.
- The validation cache is a `sync.Map` or `sync.RWMutex`-guarded map with
  per-key goroutines for in-flight de-dup.
- The SSRF redirect hook is implemented in `http.Client.CheckRedirect`.
- DNS pinning is implemented in `http.Transport.DialContext`.
- `is_atlassian_cloud_url` becomes a function that uses `net/url` + a
  suffix allowlist. `net.IP.IsPrivate()`, `IsLoopback()`, etc. cover the
  non-global check.

## Cloud-ID gateway

For Cloud OAuth 3LO, Atlassian routes every API call through a gateway:

```
https://api.atlassian.com/ex/{jira|confluence}/{cloudId}/rest/api/3/...
https://api.atlassian.com/ex/{jira|confluence}/{cloudId}/wiki/api/v2/...
```

The `{cloudId}` is the tenant ID returned from
`/oauth/token/accessible-resources`. The Go port must select this base
URL whenever `oauth_config.cloud_id != ""` and `is_cloud == true`.

## DCR controls (`HardenedOAuthProxy`)

When `ATLASSIAN_OAUTH_PROXY_ENABLE=true`, the root server exposes an OAuth
proxy with DCR controls:

- `ATLASSIAN_OAUTH_ALLOWED_CLIENT_REDIRECT_URIS` — default allowlist includes
  `http://localhost:*`, `https://chatgpt.com/...`, etc.
- `ATLASSIAN_OAUTH_ALLOWED_GRANT_TYPES` — defaults to
  `["authorization_code", "refresh_token"]`.
- `ATLASSIAN_OAUTH_REQUIRE_CONSENT` — default `true`.
- `ATLASSIAN_OAUTH_SCOPE` — space-separated scopes; if set, replaced on
  every DCR registration.
- `ATLASSIAN_OAUTH_AUDIENCE` — `api.atlassian.com` for Cloud.

For all DCR registrations, `response_types` is forced to `["code"]`.