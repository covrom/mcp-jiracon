# Jira tools — `toolset:jira_watchers`

Tools for managing issue watchers.

## `get_issue_watchers` (jira_get_issue_watchers)

**Tags:** `jira`, `read`, `toolset:jira_watchers`
**Annotations:** `title="Get Issue Watchers"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:431-459`

### Purpose

Get the count and list of users watching a Jira issue.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key (e.g. `'PROJ-123'`), pattern `^[A-Z][A-Z0-9_]+-\d+(?:-\d+)*$` |

### Handler logic

1. `jira = await get_jira_fetcher(ctx)`.
2. `result = jira.get_issue_watchers(issue_key)`.
3. JSON-serialize `result`.

### Output

```json
{
  "self": "https://...",
  "isWatching": true,
  "watchCount": 3,
  "watchers": [
    {"accountId": "...", "displayName": "...", "active": true, ...}
  ]
}
```

### Underlying REST call

`GET /issue/{issueKey}/watchers`

### Variants

- Cloud: `accountId` for each watcher.
- DC: `name`/`key`.

---

## `add_watcher` (jira_add_watcher)

**Tags:** `jira`, `write`, `toolset:jira_watchers`
**Annotations:** `title="Add Issue Watcher"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:462-504`

### Purpose

Add a user as a watcher to an issue.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |
| `user_identifier` | `str` | yes | Cloud: `accountId`. Server/DC: `username`. |

### Handler logic

1. `@check_write_access` — raises if read-only mode.
2. `result = jira.add_watcher(issue_key, user_identifier)`.
3. JSON-serialize `result` (typically `None` or success dict).

### Underlying REST call

`POST /issue/{issueKey}/watchers` with body:

- Cloud: `{"accountId": "<id>"}`.
- Server/DC: `{"username": "<name>"}` (string body, not JSON, depending
  on lib version).

### Variants

- HTTP 204 No Content on success (no body).

---

## `remove_watcher` (jira_remove_watcher)

**Tags:** `jira`, `write`, `toolset:jira_watchers`
**Annotations:** `title="Remove Issue Watcher"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:507-555`

### Purpose

Remove a user from watching an issue.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `username` | `str \| None` | no | `None` | Server/DC username |
| `account_id` | `str \| None` | no | `None` | Cloud account ID |

### Handler logic

1. `@check_write_access`.
2. `result = jira.remove_watcher(issue_key, username=username, account_id=account_id)`.
3. JSON-serialize `result`.

### Underlying REST call

`DELETE /issue/{issueKey}/watchers?accountId=…&username=…`

### Variants

- Exactly one of `username` / `account_id` should be set depending on
  the deployment type. The mixin handles both branches.