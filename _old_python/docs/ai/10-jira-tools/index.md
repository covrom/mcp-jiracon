# Jira tools — catalog

The Jira MCP server exposes 63 tools organized into 16 toolsets. Each
sub-page in this folder documents one toolset.

The full tool list is also available as a single CSV-ish table at
[../00-overview/tool-catalog.md](../00-overview/tool-catalog.md).

## Toolsets

| Toolset | Default | # tools | Sub-page |
| --- | --- | --- | --- |
| `jira_users` | no | 2 | [jira_users.md](jira_users.md) |
| `jira_watchers` | no | 3 | [jira_watchers.md](jira_watchers.md) |
| `jira_issues` | yes | 10 | [jira_issues.md](jira_issues.md) |
| `jira_transitions` | yes | 2 | [jira_transitions.md](jira_transitions.md) |
| `jira_worklog` | no | 2 | [jira_worklog.md](jira_worklog.md) |
| `jira_attachments` | no | 2 | [jira_attachments.md](jira_attachments.md) |
| `jira_agile` | no | 8 | [jira_agile.md](jira_agile.md) |
| `jira_links` | no | 5 | [jira_links.md](jira_links.md) |
| `jira_comments` | yes | 2 | [jira_comments.md](jira_comments.md) |
| `jira_projects` | no | 10 | [jira_projects.md](jira_projects.md) |
| `jira_service_desk` | no | 6 | [jira_service_desk.md](jira_service_desk.md) |
| `jira_fields` | yes | 2 | [jira_fields.md](jira_fields.md) |
| `jira_forms` | no | 3 | [jira_forms.md](jira_forms.md) |
| `jira_metrics` | no | 2 | [jira_metrics.md](jira_metrics.md) |
| `jira_development` | no | 2 | [jira_development.md](jira_development.md) |
| `jira_project_analysis` | no | 2 | [jira_project_analysis.md](jira_project_analysis.md) |

## Shared argument conventions

Across all Jira tools, several argument patterns repeat:

- **`issue_key`** is always `str` with `pattern=ISSUE_KEY_PATTERN` (regex
  `r"^[A-Z][A-Z0-9_]+-\d+(?:-\d+)*$"`). The Go port should accept this
  pattern at the JSON-schema level.
- **`project_key`** is `str` with `pattern=PROJECT_KEY_PATTERN` (regex
  `r"^[A-Z][A-Z0-9_]+$"`).
- **`user_identifier`** is `str` accepting:
  - email address (Cloud basic),
  - account ID (Cloud OAuth/PAT),
  - username or key (Server/DC).
- **`assignee`** may be a string or a JSON-object string from
  `search_assignable_users` (the tool tries `json.loads` if the string
  starts with `{`).
- **`fields` parameter** is a comma-separated string. Each tool's default
  is `",".join(DEFAULT_READ_JIRA_FIELDS)` =
  `"priority,updated,labels,issuetype,summary,assignee,description,created,reporter,status"`.
  Special values: `"*all"` (include all fields including custom),
  or omit for default.
- **JSON-string arguments** (`fields`, `additional_fields`,
  `request_field_values`, `request_participants`, `attachments`,
  `visibility`, `comment_visibility`) are parsed via
  `_parse_additional_fields`, `_parse_visibility`, etc. (see
  [server file `servers/jira.py` lines 117-247](../src/mcp_atlassian/servers/jira.py)).
- **CSV-style arguments** (`issue_keys`, `issue_ids_or_keys`,
  `components`, `versions` for the bulk helper) split on `,`.
- **`limit` and `start_at`** for pagination:
  - `limit`: typically `ge=1, le=50` for read tools; `ge=1, le=1000` for
    user search.
  - `start_at`: `ge=0`.
  - The optional `ATLASSIAN_MAX_PAGINATION_LIMIT` env var caps the limit
    (`utils/pagination.py::clamp_limit`).
- **Markdown text fields** (`description`, `body`, `comment`, `summary`)
  accept Markdown. Cloud supports the `{expand:Title}...{expand}`
  collapsible-section syntax.

## Shared return shapes

Most tools return `json.dumps(result, indent=2, ensure_ascii=False)`:

```python
return json.dumps(result, indent=2, ensure_ascii=False)
```

The Go port should use `json.MarshalIndent(result, "", "  ")` with HTML
escaping disabled (`SetEscapeHTML(false)`) so that Cyrillic / CJK
characters stay readable.

Binary output tools return:

```python
return [
    TextContent(type="text", text=json.dumps(summary, indent=2, ensure_ascii=False)),
    EmbeddedResource(
        type="resource",
        resource=BlobResourceContents(
            uri=f"attachment:///{issue_key}/{filename}",
            mimeType=resolved_mime_type,
            blob=base64_data,
        ),
    ),
    ...
]
```

or for inline images:

```python
return [
    TextContent(type="text", text=summary_json),
    ImageContent(type="image", data=base64, mimeType=resolved_mime_type),
    ...
]
```

## Error handling

Most read tools wrap network errors in a JSON-string error envelope:

```json
{
  "success": false,
  "error": "<message>",
  "issue_key": "PROJ-123"
}
```

The LLM receives this as a regular tool result and can read the error
message. This is preferred over raising exceptions because it lets the
LLM continue processing in the same agent loop.

A few tools raise `ValueError` directly:
- Tools with `@check_write_access` when `read_only=True`.
- `move_issue` when `config.is_cloud == False`.
- `batch_get_changelogs` when `config.is_cloud == False`.
- All `jira_service_desk_*` tools when `config.is_cloud == True`.
- `update_issue` for invalid input.

`MCPAtlassianAuthenticationError` is raised when the underlying
`atlassian-python-api` client receives a 401/403. The wrapping
`@handle_tool_errors` translates this into a `ToolError` for the MCP
client.

## Special guard rails

- **`JIRA_INTERNAL_ONLY_PROJECTS`** — comma-separated list of project
  keys. If a tool is asked to post a customer-visible comment / transition
  comment / link comment to an issue in one of these projects, the call
  is rejected.
- **`JIRA_FETCHER_MAX_WORKERS`** — caps concurrent Jira fetcher calls
  inside the event loop (`utils/async_utils.py`). Default 8.
- **`JIRA_SLA_WORKING_HOURS_*`** — env vars for working-hours-aware SLA
  computation (`tools/jira_metrics.md`).
- **`JIRA_PROJECTS_FILTER`** — comma-separated project keys. Tools that
  list projects (e.g., `get_all_projects`) filter their results. Some
  tools accept an in-call `projects_filter` override.

## Tool registration decorator reference

```python
@jira_mcp.tool(
    tags={"jira", "<read|write>", "toolset:<name>", ...},
    annotations={"title": "...", "readOnlyHint": True|False, "destructiveHint": True|False},
)
@check_write_access  # only on write tools
async def tool_name(ctx: Context, ...) -> str | list[...]:
    ...
```

Tag semantics enforced at runtime:

| Tag | Effect |
| --- | --- |
| `jira` | Required on every Jira tool. |
| `read` / `write` | Filter for read-only mode. |
| `toolset:<name>` | Group membership (see `jira_toolsets` table above). |
| `attachments`, `metrics`, `sla`, `development`, `analytics` | Free-form grouping hints. |

Annotation semantics:

| Annotation | Effect |
| --- | --- |
| `title` | Surfaced in the MCP `tools/list` response. |
| `readOnlyHint` | Cosmetic; tells the client the tool doesn't mutate. |
| `destructiveHint` | Cosmetic; tells the client the mutation is irreversible. |

## Go port notes

- The decorator pattern is replaced by a `RegisterJiraTool(...)` helper
  function that wraps the handler with the same logic as
  `check_write_access` + `handle_tool_errors`.
- Tags and annotations are passed through to the Go MCP SDK's
  `ToolDefinition.Tags` and `ToolDefinition.Annotations` fields.
- The `JSON-string-or-real-list` argument convention (for fields,
  components, etc.) is implemented in Go as a thin pre-validation step
  before the handler runs. Keep the same precedence: JSON first, CSV
  fallback, error on neither.