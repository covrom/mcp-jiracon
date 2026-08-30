# Confluence tools — `toolset:confluence_permissions`

Cloud-only permission checks.

## `check_content_permissions` (confluence_check_content_permissions)

**Tags:** `confluence`, `read`, `toolset:confluence_permissions`
**Annotations:** `title="Check Content Permissions"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:3014-3075`

### Purpose

Check whether a user or group can perform a specific operation on
content (page, blog, comment, attachment).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `content_id` | `str` | yes | – | Content ID |
| `user_identifier` | `str` | yes | – | Account ID (for `subject_type='user'`) or group ID (for `subject_type='group'`) |
| `operation` | `str` | yes | – | One of: `read`, `update`, `delete`, `export`, `purge`, `administer`, `create_or_delete_from_view` |
| `subject_type` | `str` | no | `"user"` | `'user'` or `'group'` |

### Underlying REST call

`POST /wiki/rest/api/content/{content_id}/permission/check` with body:

```json
{
  "operation": "read",
  "subject": {"type": "user", "identifier": "5b10a2844c20165700ede21g"}
}
```

### Output

```json
{"hasPermission": true}
```

---

## `get_space_permissions` (confluence_get_space_permissions)

**Tags:** `confluence`, `read`, `toolset:confluence_permissions`
**Annotations:** `title="Get Space Permissions"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:3078-3130`

### Purpose

List all permission assignments for a space (audit who has access).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `space_id` | `str` | yes | – | Numeric internal space ID (NOT the space key) |
| `limit` | `int` | no | `25` | Max permission entries. `ge=1`. |
| `cursor` | `str \| None` | no | `None` | Pagination cursor |

### Underlying REST call

`GET /api/v2/spaces/{space_id}/permissions?limit=N&cursor=…`

### Output

```json
{"results": [
  {"id": "...", "principal": {...}, "operation": "read", "target": {...}}
]}
```