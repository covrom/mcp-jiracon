# Confluence tools — catalog

The Confluence MCP server exposes 35 tools organized into 8 toolsets.
Each sub-page in this folder documents one toolset.

The full tool list is also available as a single CSV-ish table at
[../00-overview/tool-catalog.md](../00-overview/tool-catalog.md).

## Toolsets

| Toolset | Default | # tools | Sub-page |
| --- | --- | --- | --- |
| `confluence_pages` | yes | 14 | [confluence_pages.md](confluence_pages.md) |
| `confluence_comments` | yes | 5 | [confluence_comments.md](confluence_comments.md) |
| `confluence_labels` | no | 2 | [confluence_labels.md](confluence_labels.md) |
| `confluence_users` | no | 1 | [confluence_users.md](confluence_users.md) |
| `confluence_analytics` | no | 1 | [confluence_analytics.md](confluence_analytics.md) |
| `confluence_attachments` | no | 7 | [confluence_attachments.md](confluence_attachments.md) |
| `confluence_templates` | no | 3 | [confluence_templates.md](confluence_templates.md) |
| `confluence_permissions` | no | 2 | [confluence_permissions.md](confluence_permissions.md) |

## Shared argument conventions

- **`page_id`** is `str`. Accepts numeric IDs (e.g. `'123456789'`),
  full page URLs (e.g.
  `'https://example.atlassian.net/wiki/spaces/TEAM/pages/123456789/Page+Title'`),
  and tiny-link URLs (e.g.
  `'https://example.atlassian.net/wiki/x/N4CIO'`). The server resolves
  the input via `_resolve_page_id` before issuing the API call.
- **`space_key`** is `str` (e.g. `'DEV'`, `'TEAM'`, `'DOC'`).
- **`content_id`** is `str` — generic content ID (page or blog post).
- **`comment_id`** is `str` — parent comment ID for replies.
- **`attachment_id`** is `str` with `att`-prefix (e.g. `'att123456789'`).
- **`body`** / **`content`** accept Markdown (when `content_format =
  "markdown"`), wiki markup, storage XML, or XHTML (mapped to storage).
- **`title`** is `str` (page title).
- **`is_markdown`** is `bool` (auto-set when `content_format =
  "markdown"`).
- **`expand`** is `str` (Jira-style expand parameter, comma-separated).

### Tiny-link encoding

Confluence encodes a 64-bit page ID into a tiny-link form
(`/x/<token>`) using URL-safe base64 with padding stripped and a
single trailing `A` removed. The token is 11 chars or fewer and
matches `^[A-Za-z0-9_-]{1,11}$`.

The `_resolve_page_id` helper accepts the tiny-link form and decodes
the page ID. The Go port should preserve this convention.

## Shared return shapes

Most tools return `json.dumps(result, indent=2, ensure_ascii=False)`:

```python
return json.dumps(result, indent=2, ensure_ascii=False)
```

Error envelopes are typically one of:

```json
{"success": false, "error": "...", "page_id": "..."}
{"error": "Failed to ...", "page_id": "..."}
```

Binary output tools (attachment downloads, page images) return a list
of `TextContent` + `EmbeddedResource` / `ImageContent` blocks. See
[confluence_attachments.md](confluence_attachments.md) for details.

## Error handling

- `_resolve_page_id` silently returns the original input if the URL
  can't be decoded (warning logged).
- Page-not-found → returns `{"error": "Page with title '...' not
  found in space '...'."}` rather than raising.
- Auth errors raise `MCPAtlassianAuthenticationError`. The wrapping
  `ErrorPreservingFastMCP` translates this into a `ToolError`.
- Most write tools catch exceptions and return `{"success": false,
  ...}` envelopes so the LLM can continue.

## Content format mapping

The `content_format` parameter on `create_page` / `update_page` /
`update_page_section`:

| Value | `is_markdown` | `content_representation` | Storage format |
| --- | --- | --- | --- |
| `"markdown"` | `True` | `None` (converted server-side) | Markdown → XHTML storage |
| `"wiki"` | `False` | `"wiki"` | Confluence wiki markup |
| `"storage"` | `False` | `"storage"` | XHTML storage (verbatim) |
| `"xhtml"` | `False` | `"storage"` (mapped) | XHTML storage (verbatim) |

For `"markdown"`, the `preprocessing/confluence.py` module applies the
`md2conf` converter plus task-list and table-layout post-processing.

For `update_page_section`, only `"markdown"` or `"storage"` are valid.

## Cloud-only endpoints

These tools raise `NotImplementedError` on Server/DC (the mixin layer
returns it explicitly):

- All `confluence_templates` tools (list / get / create from template).
- All `confluence_permissions` tools (check content permissions / space
  permissions).
- `confluence_get_page_views` (analytics).
- `confluence_copy_page` with `copy_attachments=True` (Cloud has native
  copy; DC falls back to manual read+create).

The Go port should mirror this behavior.

## OAuth-only endpoints

The `ConfluenceV2Adapter` is used when `config.auth_type == "oauth"`
and `is_cloud == True`. It bypasses the v1 REST API for:

- Pages CRUD (create / update / get)
- Pages properties (emoji, width)
- Versions
- Spaces
- Footer / inline comments
- Attachments (list / get / delete)

But it falls back to v1 for `move_page`, `copy_page`, page
restrictions, and analytics.

## Tool registration decorator reference

```python
@confluence_mcp.tool(
    tags={"confluence", "<read|write>", "toolset:<name>", ...},
    annotations={"title": "...", "readOnlyHint": True|False, "destructiveHint": True|False},
)
@check_write_access  # only on write tools
async def tool_name(ctx: Context, ...) -> str | list[...]:
    ...
```

Same conventions as Jira. See
[../10-jira-tools/index.md#tool-registration-decorator-reference](../10-jira-tools/index.md#tool-registration-decorator-reference).