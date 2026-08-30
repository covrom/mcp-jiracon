# Confluence tools — `toolset:confluence_labels`

Page label management.

## `get_labels` (confluence_get_labels)

**Tags:** `confluence`, `read`, `toolset:confluence_labels`
**Annotations:** `title="Get Labels"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:593-623`

### Purpose

Get labels for Confluence content (pages, blog posts, attachments).
Works for any content type that supports labels.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `page_id` | `str` | yes | Content ID. Numeric for pages/blogs (`'123456789'`). `att`-prefixed for attachments (`'att123456789'`). |

### Underlying REST call

`GET /rest/api/content/{page_id}/label`

### Output

```json
[
  {"prefix": "global", "name": "draft", "id": "...", "self": "..."},
  {"prefix": "global", "name": "reviewed", ...}
]
```

---

## `add_label` (confluence_add_label)

**Tags:** `confluence`, `write`, `toolset:confluence_labels`
**Annotations:** `title="Add Label"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:626-676`

### Purpose

Add a label to Confluence content. Returns the updated label list.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `page_id` | `str` | yes | Content ID (page or `att<id>` for attachments) |
| `name` | `str` | yes | Label name (lowercase, no spaces recommended) |

### Underlying REST call

`POST /rest/api/content/{page_id}/label` with body
`[{"prefix":"global","name":"<label>"}]`.

### Variants

- If the label already exists, the API is idempotent (no duplicate
  created).