# Confluence tools — `toolset:confluence_templates`

Cloud-only content templates.

## `list_page_templates` (confluence_list_page_templates)

**Tags:** `confluence`, `read`, `toolset:confluence_templates`
**Annotations:** `title="List Page Templates"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:2686-2744`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `space_key` | `str \| None` | no | `None` | Space key (omit for global templates) |
| `limit` | `int` | no | `25` | Max templates. `ge=1, le=200`. |

### Underlying REST call

`GET /rest/api/template/page?limit=N&spaceKey=K`

### Output

```json
{"templates": [
  {"templateId": "...", "name": "...", "templateType": "page", "description": "..."},
  ...
], "total": 7}
```

### Variants

- **Cloud only.** Server/DC raises.

---

## `get_page_template` (confluence_get_page_template)

**Tags:** `confluence`, `read`, `toolset:confluence_templates`
**Annotations:** `title="Get Page Template"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:2747-2779`

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `template_id` | `str` | yes | Template ID |

### Underlying REST call

`GET /rest/api/template/{templateId}`

### Output

```json
{
  "templateId": "...",
  "name": "...",
  "templateType": "page",
  "description": "...",
  "body": "<ac:structured-macro>...</ac:structured-macro>"
}
```

---

## `create_page_from_template` (confluence_create_page_from_template)

**Tags:** `confluence`, `write`, `toolset:confluence_templates`
**Annotations:** `title="Create Page from Template"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:2782-2831`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `space_key` | `str` | yes | – | Target space key |
| `title` | `str` | yes | – | New page title |
| `template_id` | `str` | yes | – | Template ID |
| `parent_id` | `str \| None` | no | `None` | Optional parent page ID |

### Handler logic

`result = confluence_fetcher.create_page_from_template(space_key, title, template_id, parent_id)` →
JSON-serialize.

### Underlying REST call

`GET /rest/api/template/{templateId}` (read template) then
`POST /rest/api/content` (create page with template body).