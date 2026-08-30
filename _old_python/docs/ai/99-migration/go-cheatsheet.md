# Go migration cheat-sheet

This document is the entry point for the Go reimplementation. It maps
every major concept in the Python codebase to a Go equivalent and
calls out the libraries to use.

## Library picks

| Concern | Python | Go | Notes |
| --- | --- | --- | --- |
| MCP server | `fastmcp` | `github.com/modelcontextprotocol/go-sdk` (or your own framework) | The official Go MCP SDK is the canonical choice. |
| HTTP client | `requests` | `net/http` | Standard library is fine. |
| OAuth 2.0 | custom + `requests` | `golang.org/x/oauth2` for cloud flow, hand-rolled for DC | |
| Keyring | `keyring` | `github.com/zalando/go-keyring` | Cross-platform. |
| HTML parsing | `beautifulsoup4` | `golang.org/x/net/html` + `github.com/PuerkitoBio/goquery` | For preprocessing. |
| Markdown | `markdownify`, `md2conf` | `github.com/JohannesKaufmann/html-to-markdown` + custom post-processing | |
| JQL/CQL search | (just strings) | (just strings) | No special library needed. |
| JSON | `json` | `encoding/json` | |
| Regex | `re` | `regexp` | RE2 syntax (no lookbehind/lookahead). The Python code uses PCRE-only features sparingly; verify before porting. |
| Env vars | `os.environ` | `os.Getenv` / `os.LookupEnv` | |
| Date parsing | `python-dateutil` | `time.Parse(time.RFC3339, …)` + `time.UnixMilli(…)` | |
| YAML config | `pyyaml` | `gopkg.in/yaml.v3` | Only needed if you keep YAML config files. |
| Async | `asyncio` / `anyio` | native goroutines + channels | No need for an explicit worker cap. |
| Retry | `urllib3.util.Retry` | `github.com/hashicorp/go-retryablehttp` or hand-rolled | |
| Rate limit | token bucket (custom) | `golang.org/x/time/rate` | |
| Concurrency cap | `threading.BoundedSemaphore` | `golang.org/x/sync/semaphore` | |
| Circuit breaker | custom | `github.com/sony/gobreaker` or `github.com/failsafe-go/failsafe-go` | |
| TLS / mTLS | stdlib `ssl` | `crypto/tls` | |

## Module map

| Python module | Go package |
| --- | --- |
| `src/mcp_atlassian/servers/` | `internal/server/` |
| `src/mcp_atlassian/jira/` | `internal/jira/` |
| `src/mcp_atlassian/confluence/` | `internal/confluence/` |
| `src/mcp_atlassian/models/` | `internal/models/` |
| `src/mcp_atlassian/preprocessing/` | `internal/preprocessing/` |
| `src/mcp_atlassian/utils/` | `internal/utils/` |
| `src/mcp_atlassian/exceptions.py` | `internal/atlassian/errors.go` |

## HTTP client design

The Atlassian REST client is the most important layer. Replace the
Python `requests.Session` chain (SSL/NO_PROXY, proxy/PAC, retry,
rate-limit, concurrency, circuit-breaker) with a Go `http.Client` whose
`Transport` is a chain of `http.RoundTripper` instances.

```go
type roundTripperFunc func(http.RoundTripper) http.RoundTripper

func chain(parent http.RoundTripper, mws ...roundTripperFunc) http.RoundTripper {
    for i := len(mws) - 1; i >= 0; i-- {
        parent = mws[i](parent)
    }
    return parent
}

func NewAtlassianTransport(base http.RoundTripper, opts TransportOpts) http.RoundTripper {
    return chain(
        base,
        WithOAuthBearer(opts.AccessToken),
        WithSSRFGuard(opts.AllowedHosts),
        WithRetry(opts.Retry),
        WithRateLimit(opts.Rate),
        WithCircuitBreaker(opts.CB),
    )
}
```

### SSRF DNS pinning

The Python code uses a custom `urllib3` connection class. In Go:

```go
type ssrfDialContext struct {
    base          *net.Resolver
    trustedHosts  map[string]struct{}
}

func (d *ssrfDialContext) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
    host, port, err := net.SplitHostPort(addr)
    if err != nil { return nil, err }
    ips, err := d.base.LookupIPAddr(ctx, host)
    if err != nil { return nil, err }
    for _, ip := range ips {
        if _, trusted := d.trustedHosts[host]; !trusted {
            if !ip.IP.IsGlobalUnicast() || ip.IP.IsPrivate() || ip.IP.IsLoopback() {
                continue   // reject non-global
            }
        }
        return dialOne(network, ip.IP.String(), port)
    }
    return nil, errors.New("ssrf: no acceptable address")
}
```

### Bearer header injection

Build a `transport` that sets `req.Header.Set("Authorization", "Bearer "+
t.token)` if not already set, then delegates to `base.RoundTrip`.

### mTLS

`tls.Config{Certificates: []tls.Certificate{...}}` on the `Transport`. To
detect encrypted PEM keys, look for the `ENCRYPTED` block header and
refuse with `ErrEncryptedKey`.

### Per-request credential resolution

The Python code uses `request.state` to cache a per-request fetcher.
In Go, the MCP SDK typically exposes the request context; use a
`sync.Map` keyed by request ID (or just the credential hash), with
expiration handled by the validation cache.

```go
type fetcherCache struct {
    mu  sync.RWMutex
    fetcherByCred map[string]*JiraFetcher
    expiry        time.Duration
}

func (c *fetcherCache) Get(ctx context.Context, creds *Creds) (*JiraFetcher, error) {
    key := hash(creds)
    c.mu.RLock()
    f, ok := c.fetcherByCred[key]
    c.mu.RUnlock()
    if ok && !f.expired() { return f, nil }
    c.mu.Lock()
    defer c.mu.Unlock()
    if f, ok := c.fetcherByCred[key]; ok && !f.expired() { return f, nil }
    f = buildFetcher(creds)
    c.fetcherByCred[key] = f
    return f, nil
}
```

## OAuth flow

### Cloud OAuth (3LO)

- Use `golang.org/x/oauth2` with the Atlassian endpoints:
  - `Endpoint.AuthURL = "https://auth.atlassian.com/authorize"`
  - `Endpoint.TokenURL = "https://auth.atlassian.com/oauth/token"`
- Add `oauth2.SetAuthURLParam("audience", "api.atlassian.com")` and
  `oauth2.SetAuthURLParam("prompt", "consent")` for the authorize URL.
- After exchange, fetch
  `https://api.atlassian.com/oauth/token/accessible-resources` to resolve
  the cloud ID.
- Persist tokens via `go-keyring` (preferred) or a `0600` JSON file in
  `$HOME/.mcp-atlassian/oauth-<client_id>.json`.

### DC OAuth

- Token URL: `{base_url}/rest/oauth2/latest/token`.
- Authorize URL: `{base_url}/rest/oauth2/latest/authorize`.
- No cloud ID; `base_url` is the identity.

### Token refresh

Schedule refresh 5 minutes before expiry (`TOKEN_EXPIRY_MARGIN`).
`golang.org/x/oauth2` does this automatically with a wrapped `TokenSource`.

## JSON schema / tool registration

The Go MCP SDK uses `mcp.WithDescription(...)`, `mcp.WithString(...)`,
etc. for tool parameter schemas. The Go port should build a per-tool
schema using these helpers.

```go
func GetIssueTool() mcp.Tool {
    return mcp.NewTool("jira_get_issue",
        mcp.WithDescription("Get details of a specific Jira issue..."),
        mcp.WithString("issue_key",
            mcp.Required(),
            mcp.Pattern(`^[A-Z][A-Z0-9_]+-\d+(?:-\d+)*$`),
        ),
        mcp.WithString("fields", mcp.DefaultString("priority,updated,labels,issuetype,summary,assignee,description,created,reporter,status")),
        mcp.WithNumber("comment_limit", mcp.DefaultNumber(10), mcp.Min(0), mcp.Max(100)),
        // ... etc.
    )
}
```

## Markdown / Storage conversion

The `preprocessing/` module is the second-largest concern after the
HTTP layer. Plan to spend meaningful effort here:

- `html-to-markdown` (`github.com/JohannesKaufmann/html-to-markdown`).
- Custom post-processing for:
  - `<ac:image>` → `<ri:attachment>` fix.
  - Task-list rewriting (find `<ul>` of `[ ]/[x]` items).
  - Table-layout attribute injection.
  - User profile macro resolution.
- Jira wiki ↔ Markdown converter: write the regex set as a Go package.
  ~30 substitutions; the placeholder trick `\x00PREFIX<N>\x00` works in
  Go too.

## Models

Models are pure data structs. Translate:

```go
type JiraIssue struct {
    ID     string                 `json:"id"`
    Key    string                 `json:"key"`
    Self   string                 `json:"self"`
    Fields map[string]interface{} `json:"fields"`
}

func (i *JiraIssue) ToSimplifiedDict() map[string]any {
    // ... iterate fields, apply defaults, return map
}
```

Use a custom `UnmarshalJSON` if you need defaults (e.g. `JiraStatus`
defaults to `{name: "Unknown", id: "0"}`).

## Toolsets and filters

Toolsets are an enum-like struct per toolset:

```go
type Toolset struct {
    Name        string
    Description string
    Default    bool
}

var JiraToolsets = map[string]Toolset{
    "jira_issues":   {"Issues CRUD", true},
    "jira_fields":   {"Field metadata", true},
    "jira_comments": {"Comments", true},
    "jira_transitions": {"Status transitions", true},
    // ... etc.
}
```

The filter logic walks every registered tool, drops tools whose service
is unavailable, toolset isn't enabled, or read-only mode blocks write
tags. The same predicate runs at listing time and at call time.

## Read-only mode

`READ_ONLY_MODE=true` blocks every tool whose tag set includes
`"write"`. The Go port should:

1. Read the env var once at startup.
2. Pass the flag to each tool's registration.
3. The tool registration helper checks the flag and either:
   - Drops the tool from the listing (best UX).
   - Adds a guard in the handler that returns an error.

## Tag annotations

Tags are simple strings on the tool registration. The Go port should
mirror the same tags verbatim so the `ENABLED_TOOLS` and `TOOLSETS`
env vars continue to work.

```go
type ToolMeta struct {
    Tags        []string
    ReadOnly    bool
    Destructive bool
    Title       string
}
```

## Status codes & error envelopes

Most tools return JSON-string error envelopes rather than raising
exceptions. Mirror this:

```go
type ErrorEnvelope struct {
    Success   bool   `json:"success"`
    Error     string `json:"error"`
    IssueKey  string `json:"issue_key,omitempty"`
    PageID    string `json:"page_id,omitempty"`
    Other     map[string]any `json:"-"`
}

func (e *ErrorEnvelope) JSON() string {
    b, _ := json.MarshalIndent(e, "", "  ")
    return string(b)
}
```

## Build & deploy

| Python | Go |
| --- | --- |
| `pyproject.toml` + `uv` | `go.mod` + `go build` |
| Dockerfile + `uv run mcp-atlassian` | Dockerfile + `mcp-atlassian` binary |
| Helm chart (existing) | Same chart, swap image |
| `uv.lock` | `go.sum` (transitive deps locked) |
| `pre-commit` (ruff, mypy) | `golangci-lint` |
| `pytest` | standard `testing` + `testify` |

## Open questions for the port

1. **Do you want a separate `jira_mcp` / `confluence_mcp` binary, or a
   single `mcp-atlassian` binary that mounts both?** The current Python
   codebase is the latter. Either is fine in Go.
2. **Do you want to drop the legacy `atlassian-python-api` dependency
   entirely?** The Go port almost certainly should — there's no
   equivalent library, and the REST layer is straightforward.
3. **How much of `preprocessing/jira.py` do you actually need?** Many
   teams only use Markdown→Jira in one direction. Decide upfront whether
   to implement the wiki→Markdown direction.
4. **Do you need `ExternalAuth` mode (operator-configured upstream
   proxy auth)?** This is used in some enterprise deployments; without
   it, the dependency layer is simpler.
5. **Do you need OAuth proxy / DCR?** The current implementation
   supports it via FastMCP's `OAuthProxy`. The Go MCP SDK has its own
   OAuth provider abstraction; use that if needed.

## Migration order (recommended)

1. **Models first.** Translate `models/base.py`, `models/jira/`, and
   `models/confluence/`. These define the output shapes that all tools
   return.
2. **Utils.** `utils/env.py`, `utils/io.py`, `utils/urls.py`, `utils/oauth.py`.
   These are the cross-cutting helpers.
3. **HTTP client.** Build the Atlassian HTTP client with the SSRF guard,
   auth header injection, and per-request credential resolution.
4. **Tool registration helpers.** Build the per-tool registration
   helpers that take a struct definition and emit the MCP tool schema.
5. **One tool category at a time.** Start with the simplest tools
   (`get_user_profile`, `search_assignable_users`, `get_labels`) to
   validate the framework, then add the more complex ones.
6. **Preprocessing.** Implement Markdown ↔ storage conversion last;
   it's the largest source of edge-case bugs.
7. **OAuth setup wizard.** Implement as a CLI subcommand.

## File-by-file checklist

For each file in `src/mcp_atlassian/`, create a corresponding Go file
under `internal/`:

- `jira/` → `internal/jira/*.go` (one file per mixin)
- `confluence/` → `internal/confluence/*.go` (one file per mixin)
- `servers/jira.py` → `internal/server/jira_tools.go`
- `servers/confluence.py` → `internal/server/confluence_tools.go`
- `servers/main.py` → `internal/server/main.go`
- `servers/dependencies.py` → `internal/server/dependencies.go`
- `servers/context.py` → `internal/server/context.go`
- `servers/error_handling.py` → `internal/server/error_handling.go`
- `servers/oauth_proxy.py` → `internal/server/oauth_proxy.go`
- `servers/client_storage.py` → `internal/server/client_storage.go`
- `servers/async_utils.py` → not needed (Go handlers are sync by default)
- `preprocessing/` → `internal/preprocessing/*.go`
- `models/` → `internal/models/*.go`
- `utils/` → `internal/utils/*.go`
- `exceptions.py` → `internal/atlassian/errors.go`