# Jira tools — `toolset:jira_comments`

Comment management. Default-enabled toolset.

## `add_comment` (jira_add_comment)

**Tags:** `jira`, `write`, `toolset:jira_comments`
**Annotations:** `title="Add Comment"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2302-2385`

### Purpose

Add a comment to a Jira issue. Supports Jira Service Management
(JSM) internal / customer-visible flag.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Issue key |
| `body` | `str` | yes | – | Markdown comment body. Alias: `comment` |
| `visibility` | `str \| None` | no | `None` | JSON string `{"type": "group", "value": "jira-users"}` |
| `public` | `bool \| None` | no | `None` | JSM/Service Desk: `True` = customer-visible, `False` = internal/agent-only. Uses ServiceDesk API. |

### Handler logic

1. `@check_write_access`.
2. Parse `visibility` JSON → dict.
3. **JSM internal-only guard**:
   - If `public == False` AND the issue's project is in
     `JIRA_INTERNAL_ONLY_PROJECTS`, honor `public=False` → call
     `add_comment` with `public=False`.
   - If `public == False` AND project is NOT internal-only, treat as
     `public=None` (most MCP clients auto-fill `false`; routing through
     ServiceDesk API would 403 for non-JSM issues).
   - Otherwise pass through.
4. `result = jira.add_comment(issue_key, body, visibility_dict, public=public_value)`.
5. JSON-serialize `result`.

### Underlying REST call

- Regular comment: `POST /issue/{issueKey}/comment` with body
  `{body, visibility?}`.
- ServiceDesk (JSM) comment: `POST /rest/servicedeskapi/request/{issueKey}/comment`
  with body `{body, public}`.

### Variants

- Body alias `comment` is accepted via Pydantic `AliasChoices`.
- The guard against `JIRA_INTERNAL_ONLY_PROJECTS` ensures that operators
  who marked certain projects as internal-only cannot accidentally post
  customer-visible comments.

---

## `edit_comment` (jira_edit_comment)

**Tags:** `jira`, `write`, `toolset:jira_comments`
**Annotations:** `title="Edit Comment"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2388-2432`

### Purpose

Edit an existing comment on a Jira issue.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Issue key |
| `comment_id` | `str` | yes | – | The comment ID |
| `body` | `str` | yes | – | Updated Markdown body |
| `visibility` | `str \| None` | no | `None` | JSON string `{"type": "group", "value": "jira-users"}` |

### Handler logic

1. `@check_write_access`.
2. Parse `visibility` JSON.
3. `result = jira.edit_comment(issue_key, comment_id, body, visibility_dict)`.
4. JSON-serialize `result`.

### Underlying REST call

`PUT /issue/{issueKey}/comment/{commentId}` with body `{body, visibility?}`.

### Variants

- For JSM internal-only projects, editing a public comment is blocked
  server-side by the underlying client.