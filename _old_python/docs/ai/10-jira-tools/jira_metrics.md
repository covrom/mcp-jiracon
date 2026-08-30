# Jira tools — `toolset:jira_metrics`

Issue date analysis and SLA computations.

## `get_issue_dates` (jira_get_issue_dates)

**Tags:** `jira`, `read`, `metrics`, `toolset:jira_metrics`
**Annotations:** `title="Get Issue Dates"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:4149-4203`

### Purpose

Get date information and status transition history for a Jira issue.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `include_status_changes` | `bool` | no | `True` | Include status change history with timestamps |
| `include_status_summary` | `bool` | no | `True` | Include aggregated time per status |

### Handler logic

`result = jira.get_issue_dates(issue_key, include_created=True, include_updated=True, include_due_date=True, include_resolution_date=True, include_status_changes=…, include_status_summary=…)` →
JSON-serialize.

### Output

```json
{
  "created": "2023-08-01T...",
  "updated": "...",
  "due_date": "...",
  "resolution_date": "...",
  "status_changes": [
    {"from": "Open", "to": "In Progress", "timestamp": "...", "duration_seconds": 86400}
  ],
  "status_summary": [
    {"status": "Open", "total_seconds": 86400},
    {"status": "In Progress", "total_seconds": 259200}
  ]
}
```

### Underlying REST call

`GET /issue/{issueKey}?expand=changelog`. The metrics module parses the
changelog client-side; no separate REST call.

---

## `get_issue_sla` (jira_get_issue_sla)

**Tags:** `jira`, `read`, `metrics`, `sla`, `toolset:jira_metrics`
**Annotations:** `title="Get Issue SLA"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:4206-4284`

### Purpose

Compute SLA metrics for a Jira issue (cycle time, lead time, time in
status, due-date compliance, resolution time, first response time).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `metrics` | `str \| None` | no | `None` | Comma-separated subset of `cycle_time`, `lead_time`, `time_in_status`, `due_date_compliance`, `resolution_time`, `first_response_time`. Default = `cycle_time,time_in_status`. |
| `working_hours_only` | `bool \| None` | no | `None` | Override env `JIRA_SLA_WORKING_HOURS_ONLY` |
| `include_raw_dates` | `bool` | no | `False` | Include raw date values in response |

### Environment variables

| Var | Default | Effect |
| --- | --- | --- |
| `JIRA_SLA_WORKING_HOURS_ONLY` | `false` | Working-hours filtering |
| `JIRA_SLA_WORKING_HOURS_START` | `"09:00"` | Working day start |
| `JIRA_SLA_WORKING_HOURS_END` | `"17:00"` | Working day end |
| `JIRA_SLA_WORKING_DAYS` | `"1,2,3,4,5"` | Mon-Fri (1=Mon, 7=Sun) |
| `JIRA_SLA_TIMEZONE` | system | Timezone for working-hours calc |

### Handler logic

Parse `metrics` CSV → list. Call `jira.get_issue_sla(...)`. JSON-serialize.

### Underlying REST call

`GET /issue/{issueKey}?expand=changelog`. All SLA math is done client-side
in `jira/metrics.py`.

### Variants

- The `metrics` argument defaults come from `JIRA_SLA_DEFAULT_METRICS` env
  var, falling back to `["cycle_time", "time_in_status"]`.