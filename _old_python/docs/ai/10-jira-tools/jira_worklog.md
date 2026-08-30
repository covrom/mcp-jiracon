# Jira tools — `toolset:jira_worklog`

Time-tracking operations.

## `get_worklog` (jira_get_worklog)

**Tags:** `jira`, `read`, `toolset:jira_worklog`
**Annotations:** `title="Get Worklog"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1148-1174`

### Purpose

Get all worklog entries for a Jira issue.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |

### Handler logic

`jira.get_worklogs(issue_key)` → JSON-serialize `{"worklogs": [...]}`.

### Output

```json
{
  "worklogs": [
    {
      "id": "100023",
      "self": "https://...",
      "author": {...},
      "comment": "...",
      "created": "2023-08-01T12:00:00.000+0000",
      "updated": "2023-08-01T12:00:00.000+0000",
      "started": "2023-08-01T12:00:00.000+0000",
      "timeSpent": "1h 30m",
      "timeSpentSeconds": 5400,
      "issueId": "10001"
    }
  ]
}
```

### Underlying REST call

`GET /issue/{issueKey}/worklog`

---

## `add_worklog` (jira_add_worklog)

**Tags:** `jira`, `write`, `toolset:jira_worklog`
**Annotations:** `title="Add Worklog"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2435-2508`

### Purpose

Log work against a Jira issue. Optionally updates the original or
remaining estimate.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `time_spent` | `str` | yes | – | Jira duration format (`'1h 30m'`, `'1d'`, `'30m'`, `'4h'`) |
| `comment` | `str \| None` | no | `None` | Markdown comment |
| `started` | `str \| None` | no | `None` | ISO 8601 start time. Default = current time. |
| `original_estimate` | `str \| None` | no | `None` | New original estimate (e.g. `'1d'`) |
| `remaining_estimate` | `str \| None` | no | `None` | New remaining estimate (e.g. `'4h'`) |

### Handler logic

1. `@check_write_access`.
2. `worklog = jira.add_worklog(issue_key, time_spent, comment, started, original_estimate, remaining_estimate)`.
3. `{"message": "Worklog added successfully", "worklog": worklog}`.

### Underlying REST call

`POST /issue/{issueKey}/worklog` with body:

```json
{
  "timeSpent": "1h 30m",
  "started": "2023-08-01T12:00:00.000+0000",
  "comment": "Markdown comment",
  "adjustEstimate": "auto" | "new" | "manual",
  "newEstimate": "...",
  "reduceBy": "..."
}
```

The `adjustEstimate`/`newEstimate`/`reduceBy` fields are computed by
the mixin based on whether `original_estimate` / `remaining_estimate`
are provided.

### Variants

- The Jira duration format is parsed by `utils/jira/constants.py`:
  - `Nd` for days
  - `Nh` for hours
  - `Nm` for minutes
  - combinations like `1h 30m`