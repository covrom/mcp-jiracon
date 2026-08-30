# Error handling patterns

The codebase uses a small set of error-handling conventions. The Go
port must mirror all of them.

## Layers of error handling

### Layer 1: Atlassian REST client

`utils/decorators.py` provides two sync decorators used to wrap
Jira/Confluence client methods:

- `@handle_auth_errors(service_name="Atlassian API")` — catches
  `requests.HTTPError` with status 401/403 and raises
  `MCPAtlassianAuthenticationError`. Other errors propagate.
- `@handle_atlassian_api_errors(service_name="Atlassian API")` — same
  auth handling plus key/value/network/unexpected error classification:

```python
@handle_atlassian_api_errors(service_name="Jira")
def create_issue(...):
    ...
```

Mapping:

| Exception | Translated to |
| --- | --- |
| `HTTPError` with 401/403 | `MCPAtlassianAuthenticationError("Atlassian API authentication failed: HTTP 401")` |
| `KeyError` | `ValueError("Error processing Atlassian API results: missing key <e>")` |
| `RequestException` | `ValueError("Network error during <op>: <e>")` |
| `ValueError`/`TypeError` | `ValueError("Error processing Atlassian API results: <e>")` |
| Generic `Exception` | `RuntimeError("Unexpected error during <op>: <e>")` |

### Layer 2: FastMCP tool wrapper

`utils/decorators.py::handle_tool_errors` is auto-applied by
`ErrorPreservingFastMCP.tool(...)`. It wraps each `async def` tool:

```python
def handle_tool_errors(func):
    @functools.wraps(func)
    async def wrapper(*args, **kwargs):
        try:
            return await func(*args, **kwargs)
        except ToolError:
            raise
        except Exception as e:
            logger.exception(...)
            raise ToolError(f"Error calling tool '<name>': <e>") from e
    return wrapper
```

So any unhandled exception in a tool is caught here and translated to
a `ToolError`. The MCP client sees it as a structured error.

### Layer 3: Read-only guard

`@check_write_access` is applied per-tool by tool authors:

```python
@check_write_access
async def create_issue(ctx, ...):
    ...
```

Behavior:

- Pull `app_lifespan_context` from `ctx.request_context.lifespan_context`.
- If `read_only == True`, raise `ValueError("Cannot <verb> in read-only mode.")`.
- Otherwise delegate.

The verb is derived from the function name (`create_issue` →
`"Cannot create in read-only mode."`).

### Layer 4: JSM internal-only guard

A handful of tools check the issue's project against
`JIRA_INTERNAL_ONLY_PROJECTS`. If matched, certain mutations are
rejected:

- `add_comment(public=True)` — public comment in an internal-only project.
- `add_comment(public=False)` outside an internal-only project is treated
  as `public=None` (most MCP clients auto-fill `false`; routing through
  the ServiceDesk API would 403 for non-JSM issues).
- `transition_issue(comment=...)` — transition comment in an internal-only project.
- `create_issue_link(comment=...)` — link comment in an internal-only project.

### Layer 5: Tool-specific JSON error envelopes

Most tools wrap network errors in a JSON-string error envelope:

```json
{"success": false, "error": "<msg>", "issue_key": "..."}
```

Returned as a regular tool result rather than raised, so the LLM
agent loop continues.

## Exception hierarchy

```
Exception
└── MCPAtlassianAuthenticationError
```

There's a single shared exception in the entire codebase
(`exceptions.py`). The Go port should follow the same minimalism:

```go
var ErrAuthentication = errors.New("atlassian: authentication failed")

func IsAuthError(err error) bool {
    return errors.Is(err, ErrAuthentication)
}
```

## Per-tool error handling recipes

### Read tool on a missing issue

```python
# Returns an error envelope, doesn't raise
try:
    issue = jira.get_issue(key)
except Exception as e:
    return json.dumps({"success": false, "error": str(e), "issue_key": key})
```

Go equivalent:

```go
func getIssue(ctx context.Context, key string) (string, error) {
    jira, err := getJiraFetcher(ctx)
    if err != nil { return "", err }
    issue, err := jira.GetIssue(key)
    if err != nil {
        return jsonString(map[string]any{
            "success": false,
            "error": err.Error(),
            "issue_key": key,
        }), nil
    }
    return jsonString(issue.ToSimplifiedDict()), nil
}
```

### Write tool in read-only mode

```python
# @check_write_access raises
@check_write_access
async def create_issue(ctx, ...):
    ...
```

Go equivalent:

```go
func createIssue(ctx context.Context, ...) (string, error) {
    if readOnly {
        return "", fmt.Errorf("Cannot create in read-only mode.")
    }
    // ...
}
```

### Cloud-only tool on DC

```python
if not jira.config.is_cloud:
    raise NotImplementedError(
        "Batch get issue changelogs is only available on Jira Cloud."
    )
```

Go equivalent:

```go
if !jira.Config.IsCloud {
    return "", ErrNotImplementedOnDC
}
```

Define `ErrNotImplementedOnDC` in your errors package.

## How to translate @check_write_access in Go

The simplest approach: every tool handler checks a global `readOnly`
flag read from the lifespan context:

```go
func readOnlyFromContext(ctx context.Context) bool {
    appState := appStateFromContext(ctx)
    return appState.ReadOnly
}

func checkWriteAccess(ctx context.Context) error {
    if readOnlyFromContext(ctx) {
        return fmt.Errorf("read-only mode: write tools are disabled")
    }
    return nil
}
```

Alternatively, do this at registration time: skip registering write tools
when read-only mode is enabled. This means the tool never appears in
`tools/list` — the cleanest UX.

## Pattern: `run_jira_fetcher_call`

`run_jira_fetcher_call(func, *args, **kwargs)` runs a sync blocking call
in a worker thread. In Go this is unnecessary — handlers are
straightforward sync functions (or wrap `sync/errgroup` if you need
parallelism).

## MCP `ToolError`

`mcp.ToolError` is the FastMCP / Python MCP SDK's structured error type.
The Go MCP SDK has its own equivalent. When you raise / return it, the
MCP client receives a JSON-RPC error response with the message.

## Logging

`utils/logging.py::setup_logging` configures the root logger once at
server start. The Go port should use `slog` (Go 1.21+) with the same
log level env var.

`mask_sensitive` strips credentials from log messages. The Go port
should use `slog.LogAttrs` with structured attributes rather than
formatting strings into log messages, so credential fields don't leak
through accidental formatting.