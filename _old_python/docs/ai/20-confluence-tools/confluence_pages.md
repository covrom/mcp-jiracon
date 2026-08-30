# Confluence tools — `toolset:confluence_pages`

The largest toolset, covering page CRUD, history, diff, sections,
restrictions, and copy. Default-enabled.

## `search` (confluence_search)

**Tags:** `confluence`, `read`, `toolset:confluence_pages`
**Annotations:** `title="Search Content"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:188-277`

### Purpose

Search Confluence content using either simple text or CQL.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `query` | `str` | yes | – | Text or CQL string. Detects simple text by absence of `=`, `~`, `>`, `<`, ` AND `, ` OR `, `currentUser()`. |
| `limit` | `int` | no | `10` | Max results. `ge=1, le=50`. |
| `spaces_filter` | `str \| None` | no | `None` | Comma-separated space keys (overrides `CONFLUENCE_SPACES_FILTER`) |

### Handler logic

If `query` looks like a simple text search:
1. Wrap as `siteSearch ~ "<query>"`.
2. On failure (e.g. Cloud returning 400), fall back to `text ~ "<query>"`.
3. Call `confluence_fetcher.search(query, limit, spaces_filter)`.

Else pass through as CQL.

### Underlying REST call

`GET /rest/api/content/search` with `cql`, `limit`, `expand=content.history,content.version`.

### Variants

- Cloud and DC both support CQL.
- `siteSearch` is Cloud-supported; `text` is the universal fallback.

---

## `get_page` (confluence_get_page)

**Tags:** `confluence`, `read`, `toolset:confluence_pages`
**Annotations:** `title="Get Page"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:280-406`

### Purpose

Fetch a Confluence page by ID (numeric, URL, or tiny link) or by title
within a space.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `page_id` | `str \| None` | no | `None` | Numeric ID, full URL, or tiny link. With `BeforeValidator` that coerces non-None to `str`. |
| `title` | `str \| None` | no | `None` | Exact page title (must be combined with `space_key`) |
| `space_key` | `str \| None` | no | `None` | Space key (e.g. `'DEV'`) |
| `include_metadata` | `bool` | no | `True` | Include metadata (creation date, last update, version, labels) |
| `convert_to_markdown` | `bool` | no | `True` | `True` → markdown, `False` → raw HTML storage |

### Handler logic

1. If `page_id`:
   - Resolve via `_resolve_page_id`.
   - `page = confluence_fetcher.get_page_content(page_id, convert_to_markdown)`.
   - On error → `{"error": "Failed to retrieve page by ID '<id>': ..."}`.
2. Else if `title + space_key`:
   - `page = confluence_fetcher.get_page_by_title(space_key, title, convert_to_markdown)`.
   - On not found → `{"error": "Page with title '<title>' not found in space '<key>'."}`.
3. Else: raise `ValueError("Either 'page_id' OR both 'title' and 'space_key' must be provided.")`.
4. If `include_metadata` → `{"metadata": page.to_simplified_dict()}`.
5. Else → `{"content": {"value": page_object.content}}`.

### Variants

- Tiny links like `https://example.atlassian.net/wiki/x/N4CIO` are
  decoded to numeric IDs via `_decode_confluence_tiny_id`.

---

## `get_page_children` (confluence_get_page_children)

**Tags:** `confluence`, `read`, `toolset:confluence_pages`
**Annotations:** `title="Get Page Children"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:409-506`

### Purpose

Get child pages and (optionally) folders of a parent page.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `parent_id` | `str` | yes | – | Parent page ID |
| `expand` | `str` | no | `"version"` | Expand param |
| `limit` | `int` | no | `25` | Max results. `ge=1, le=50`. |
| `include_content` | `bool` | no | `False` | Include page content |
| `convert_to_markdown` | `bool` | no | `True` | Only relevant when `include_content=True` |
| `start` | `int` | no | `0` | Pagination start |
| `include_folders` | `bool` | no | `True` | Include child folders |

### Handler logic

If `include_content` and `"body"` not in `expand`, append
`body.storage` to expand.

- Cloud OAuth: `GET /api/v2/pages/{page_id}/direct-children` with
  `limit`/`cursor` (cursor translated to start/limit by mixin).
- DC / Basic / PAT: `GET /rest/api/content/{page_id}/child/page` (+ 
  `child/folder` if `include_folders`).

Output: `{parent_id, count, limit_requested, start_requested, results: [page_dict, ...]}`.

---

## `get_space_page_tree` (confluence_get_space_page_tree)

**Tags:** `confluence`, `read`, `toolset:confluence_pages`
**Annotations:** `title="Get Space Page Tree"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:509-558`

### Purpose

Get a flat list of pages in a space with `parent_id` and `depth`
attributes for token-efficient hierarchy traversal.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `space_key` | `str` | yes | – | Space key |
| `limit` | `int` | no | `100` | Max pages. `ge=1, le=1000`. |

### Output

```json
{
  "space_key": "DEV",
  "total_pages": 42,
  "pages": [
    {"id": "1", "title": "Home", "parent_id": null, "position": "...",
     "depth": 0},
    {"id": "2", "title": "Setup", "parent_id": "1", "position": "...",
     "depth": 1}
  ],
  "has_more": false,
  "hint": "Results truncated at 100 pages. Increase limit to see more."
}
```

---

## `create_page` (confluence_create_page)

**Tags:** `confluence`, `write`, `toolset:confluence_pages`
**Annotations:** `title="Create Page"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:679-863`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `space_key` | `str` | yes | – | Target space key |
| `title` | `str` | yes | – | Page title |
| `content` | `str \| None` | no | `None` | Inline body (mutually exclusive with `content_file`) |
| `parent_id` | `str \| None` | no | `None` | Optional parent page ID |
| `content_format` | `str` | no | `"markdown"` | `"markdown"` / `"wiki"` / `"storage"` / `"xhtml"` |
| `enable_heading_anchors` | `bool` | no | `False` | Markdown only |
| `include_content` | `bool` | no | `False` | Whether to include body in response |
| `emoji` | `str \| None` | no | `None` | Page title emoji (e.g. `'📝'`) |
| `content_file` | `str \| None` | no | `None` | Filesystem path to UTF-8 body (workspace-confined via `validate_safe_path`) |
| `page_width` | `str \| None` | no | `None` | `'full-width'` or `'default'` |
| `table_layout` | `str \| None` | no | `None` | `'full-width'` (1800px), `'wide'` (960px), `'default'` (760px). Markdown only. |
| `subtype` | `str \| None` | no | `None` | `'live'` for Confluence Live Doc (Cloud only) |

### Handler logic

1. `@check_write_access`.
2. Validate `content_format`.
3. Resolve body via `_resolve_page_content(content, content_file)`
   (exactly one must be supplied; `content_file` is workspace-confined).
4. `page = confluence_fetcher.create_page(space_key, title, body, parent_id, is_markdown, enable_heading_anchors, content_representation, emoji, subtype, page_width, table_layout)`.
5. Pop `content` from result if `include_content=False`.
6. `{"message": "Page created successfully", "page": result}`.

### Underlying REST call

- Cloud OAuth: `POST /api/v2/pages` with body
  `{spaceId, status, title, body:{representation, value}, parentId?, subtype?}`.
- Otherwise: `POST /rest/api/content` with body
  `{type:"page", space:{key}, title, ancestors?:[{id}], body:{storage:{value, representation}}}`.

---

## `update_page` (confluence_update_page)

**Tags:** `confluence`, `write`, `toolset:confluence_pages`
**Annotations:** `title="Update Page"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:866-1041`

### Input schema

Same as `create_page` but for an existing page (no `subtype`); plus:

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `page_id` | `str` | yes | – | Page ID |
| `is_minor_edit` | `bool` | no | `False` | Minor edit flag |
| `version_comment` | `str \| None` | no | `None` | Version comment |

### Underlying REST call

- Cloud OAuth: `PUT /api/v2/pages/{page_id}` with body
  `{id, status, title, body:{representation, value}, version:{number:N+1, message?}}`.
- Otherwise: `PUT /rest/api/content/{page_id}` with body
  `{id, type, title, body:{storage}, version:{number:N+1, message?}, ancestors?, minor_edit?}`.

---

## `update_page_section` (confluence_update_page_section)

**Tags:** `confluence`, `write`, `toolset:confluence_pages`
**Annotations:** `title="Update Page Section"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:1044-1143`

### Purpose

Replace only the body under a named heading on a page, preserving all
other sections, macros, layouts, and Confluence-specific elements.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `page_id` | `str` | yes | – | Page ID |
| `heading_text` | `str` | yes | – | Exact heading text (case-sensitive) |
| `new_content` | `str` | yes | – | New body content (heading not included) |
| `content_format` | `str` | no | `"markdown"` | `"markdown"` or `"storage"` |
| `is_minor_edit` | `bool` | no | `False` | Minor edit |
| `version_comment` | `str \| None` | no | `None` | Defaults to `""` |

### Handler logic

1. `@check_write_access`.
2. `updated_page = confluence_fetcher.update_page_section(...)`.
3. Pop `content` from the result.
4. `{"message": "Section '<heading>' updated successfully", "page": updated_page}`.

---

## `delete_page` (confluence_delete_page)

**Tags:** `confluence`, `write`, `toolset:confluence_pages`
**Annotations:** `title="Delete Page"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:1146-1188`

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `page_id` | `str` | yes | Page ID |

### Handler logic

1. `@check_write_access`.
2. `result = confluence_fetcher.delete_page(page_id)`.
3. On success → `{"success": True, "message": "Page <id> deleted successfully"}`.
4. On exception → `{"success": False, "message": ..., "error": ...}`.

### Underlying REST call

- Cloud OAuth: `DELETE /api/v2/pages/{page_id}`.
- Otherwise: `DELETE /rest/api/content/{page_id}` (with `status=trashed`).

---

## `move_page` (confluence_move_page)

**Tags:** `confluence`, `write`, `toolset:confluence_pages`
**Annotations:** `title="Move Page"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:1191-1267`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `page_id` | `str` | yes | – | Page to move |
| `target_parent_id` | `str \| None` | no | `None` | Target parent page ID |
| `target_space_key` | `str \| None` | no | `None` | Target space (for cross-space moves) |
| `position` | `str` | no | `"append"` | `"append"`, `"above"`, `"below"` |

### Validation

At least one of `target_parent_id` or `target_space_key` must be
provided. Otherwise raise `ValueError`.

### Underlying REST call

`PUT /rest/api/content/{page_id}/move/{position}/{targetId}` (or
`/move/{position}` when moving to a space root). Cloud OAuth uses the
v2 adapter's v1 fallback path.

---

## `get_page_history` (confluence_get_page_history)

**Tags:** `confluence`, `read`, `toolset:confluence_pages`
**Annotations:** `title="Get Page History"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:1600-1677`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `page_id` | `str` | yes | – | Page ID |
| `version` | `int` | yes | – | Version number. `ge=1`. |
| `convert_to_markdown` | `bool` | no | `True` | Markdown vs raw HTML |

### Underlying REST call

- Cloud OAuth: `GET /api/v2/pages/{page_id}/versions` then
  `GET /api/v2/versions/{versionId}?body-format=storage`.
- Otherwise: `GET /rest/api/content/{page_id}?status=historical&version=N&expand=...`.

---

## `get_page_diff` (confluence_get_page_diff)

**Tags:** `confluence`, `read`, `toolset:confluence_pages`
**Annotations:** `title="Get Page Version Diff"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:1680-1754`

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `page_id` | `str` | yes | Page ID |
| `from_version` | `int` | yes | Source version. `ge=1`. |
| `to_version` | `int` | yes | Target version. `ge=1`. |

### Handler logic

`result = confluence_fetcher.get_page_version_diff(page_id, from_version, to_version)`
→ JSON-serialize. Internally calls `get_page_history` twice and computes
a unified diff client-side.

---

## `get_page_restrictions` (confluence_get_page_restrictions)

**Tags:** `confluence`, `read`, `toolset:confluence_pages`
**Annotations:** `title="Get Page Restrictions"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:2834-2860`

### Underlying REST call

`GET /rest/api/content/{page_id}/restriction/byOperation`

### Output

```json
{
  "read": {"users": [...], "groups": [...]},
  "update": {"users": [...], "groups": [...]}
}
```

Empty lists = unrestricted.

---

## `set_page_restrictions` (confluence_set_page_restrictions)

**Tags:** `confluence`, `write`, `toolset:confluence_pages`
**Annotations:** `title="Set Page Restrictions"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:2863-2941`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `page_id` | `str` | yes | – | Page ID |
| `read_users` | `list[str] \| None` | no | `None` | Account IDs (Cloud) or usernames (DC) |
| `read_groups` | `list[str] \| None` | no | `None` | Group names |
| `edit_users` | `list[str] \| None` | no | `None` | Users allowed to edit |
| `edit_groups` | `list[str] \| None` | no | `None` | Groups allowed to edit |

### Underlying REST call

`PUT /rest/api/content/{page_id}/restriction` with body
`[{operation:"read", restrictions:{user:[...], group:[...]}}, ...]`.

### Variants

- Omitting all parameters (or passing empty lists) removes all
  restrictions.
- User entries use `{"type":"known","accountId":...}` on Cloud,
  `{"type":"known","username":...}` on DC.

---

## `copy_page` (confluence_copy_page)

**Tags:** `confluence`, `write`, `toolset:confluence_pages`
**Annotations:** `title="Copy Page"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:2944-3011`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `source_page_id` | `str` | yes | – | Page to copy |
| `destination_space_key` | `str` | yes | – | Space key for new page |
| `new_title` | `str` | yes | – | Title for the new page |
| `destination_parent_id` | `str \| None` | no | `None` | Parent page ID in destination space |
| `copy_attachments` | `bool` | no | `True` | Cloud only |

### Underlying REST call

- Cloud: `POST /wiki/rest/api/content/{source_page_id}/copy` with body
  `{copyAttachments, copyPermissions, copyProperties, copyLabels,
  pageTitle, destination:{type:"parent_page"|"space", value:idOrKey}}`.
- DC: GET source + POST `/rest/api/content` (manual; no native copy).