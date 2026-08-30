# Jira tools — `toolset:jira_forms`

ProForma forms (Cloud-only). Form IDs are UUIDs.

## `get_issue_proforma_forms` (jira_get_issue_proforma_forms)

**Tags:** `jira`, `read`, `toolset:jira_forms`
**Annotations:** `title="Get Issue Forms"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3894-3946`

### Purpose

List ProForma forms attached to a Jira issue.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |

### Handler logic

1. `forms = jira.get_issue_forms(issue_key)`.
2. `{"success": True, "forms": [f.to_simplified_dict() for f in forms], "count": N}`.
3. On error (auth / network / ValueError "not found") → JSON error.

### Underlying REST call

`GET /rest/api/3/issue/{issueIdOrKey}/form`

---

## `get_proforma_form_details` (jira_get_proforma_form_details)

**Tags:** `jira`, `read`, `toolset:jira_forms`
**Annotations:** `title="Get Form Details"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3949-4016`

### Purpose

Get details (ADF design structure) for a specific ProForma form.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |
| `form_id` | `str` | yes | Form UUID |

### Underlying REST call

`GET /rest/api/3/issue/{issueIdOrKey}/form/{formId}`

---

## `update_proforma_form_answers` (jira_update_proforma_form_answers)

**Tags:** `jira`, `write`, `toolset:jira_forms`
**Annotations:** `title="Update Form Answers"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:4019-4146`

### Purpose

Update ProForma form answers.

### Known limitation

> **DATETIME fields** lose their time component — Jira Forms API stores
> only the date. Workaround: read the underlying custom field ID from the
> form design's `jiraField` property and use `update_issue` to write the
> field with an ISO 8601 string (e.g. `"2026-01-09T11:50:00-08:00"`).

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |
| `form_id` | `str` | yes | Form UUID |
| `answers` | `list[dict]` | yes | Each answer: `{questionId, type, value}`. Types: TEXT, NUMBER, DATE, DATETIME, SELECT, MULTI_SELECT, CHECKBOX. |

### Handler logic

1. `@check_write_access`.
2. For each answer, if `type` is DATE/DATETIME and `value` is a string,
   convert via `convert_datetime_to_timestamp(value, type)`:
   - ISO 8601 string → Unix ms (UTC).
   - Integer Unix ms → pass through.
3. `result = jira.update_form_answers(issue_key, form_id, processed_answers)`.
4. JSON-serialize success envelope.

### Underlying REST call

`PUT /rest/api/3/issue/{issueIdOrKey}/form/{formId}` with body
`{answers: [...]}` where each answer's `value` is in Unix milliseconds
for DATE/DATETIME fields.