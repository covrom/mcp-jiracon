# Jira tools — `toolset:jira_transitions`

Status workflow operations. Default-enabled toolset.

## `get_transitions` (jira_get_transitions)

**Tags:** `jira`, `read`, `toolset:jira_transitions`
**Annotations:** `title="Get Transitions"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1119-1145`

### Purpose

List the available status transitions for a Jira issue (so the LLM knows
which `transition_id` to pass to `transition_issue`).

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |

### Handler logic

`jira.get_available_transitions(issue_key)` → JSON-serialize.

### Output

```json
[
  {
    "id": "11",
    "name": "To Do",
    "to": {"self": "...", "id": "10000", "name": "To Do", "statusCategory": {...}},
    "hasScreen": false,
    "isGlobal": true,
    "isInitial": false,
    "isAvailable": true,
    "fields": {}
  },
  ...
]
```

### Underlying REST call

`GET /issue/{issueKey}/transitions`

---

## `transition_issue` (jira_transition_issue)

**Tags:** `jira`, `write`, `toolset:jira_transitions`
**Annotations:** `title="Transition Issue"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2768-2851`

### Purpose

Move an issue to a new status via a known transition.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `transition_id` | `str` | yes | – | Transition ID from `get_transitions` |
| `fields` | `str \| None` | no | `None` | JSON string of fields required by the transition (e.g. `{"resolution": {"name": "Fixed"}}`) |
| `comment` | `str \| None` | no | `None` | Markdown comment visible in history |

### Handler logic

1. `@check_write_access`.
2. Validate both required args.
3. If `comment` is non-empty and the issue's project is in
   `JIRA_INTERNAL_ONLY_PROJECTS` → reject (link / transition comments
   may be customer-visible on JSM and cannot be forced internal).
4. Parse `fields` JSON.
5. `issue = jira.transition_issue(issue_key, transition_id, fields=update_fields, comment=comment)`.
6. `{"message": "Issue <key> transitioned successfully", "issue": issue.to_simplified_dict() if issue else None}`.

### Underlying REST call

`POST /issue/{issueKey}/transitions` with body:

```json
{
  "transition": {"id": "11"},
  "fields": {"resolution": {"name": "Fixed"}},
  "update": {},
  "historyMetadata": null
}
```

### Variants

- Some transitions require fields (e.g. resolution). Use `get_transitions`
  first and inspect the `fields` map.
- The `comment` is rejected for `JIRA_INTERNAL_ONLY_PROJECTS` projects
  because it cannot be marked internal on JSM.