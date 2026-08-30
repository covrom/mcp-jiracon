# Jira tools — `toolset:jira_issues`

The core CRUD surface for issues. This toolset is enabled by default.

## `get_issue` (jira_get_issue)

**Tags:** `jira`, `read`, `toolset:jira_issues`
**Annotations:** `title="Get Issue"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:558-735`

### Purpose

Fetch a Jira issue with full detail, including Epic link / parent
relationships. Supports inline enrichment sections via `include` to
reduce follow-up calls.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `fields` | `str` | no | `priority,updated,labels,issuetype,summary,assignee,description,created,reporter,status` | Comma-separated fields, or `"*all"` for all fields including custom |
| `expand` | `str \| None` | no | `None` | Jira `expand` query param (e.g. `'renderedFields'`, `'transitions'`, `'changelog'`) |
| `comment_limit` | `int` | no | `10` | Max comments. `ge=0, le=100`. 0 suppresses comments. |
| `properties` | `str \| None` | no | `None` | Comma-separated issue properties |
| `update_history` | `bool` | no | `True` | Update the requester's view history |
| `include` | `str \| None` | no | `None` | Comma-separated inline sections: `all`, `remote_links`, `transitions`, `watchers`, `changelog`, `comments`, `worklogs` |
| `use_display_names` | `bool` | no | `False` | Replace custom-field keys with human-readable display names |

### Handler logic

1. `fields_list` is a list if `fields` is a CSV; `None` or `*all`
   passes through.
2. If `comments` is in `include`, add `"comment"` to the fields list
   (unless `*all`).
3. If `changelog` is in `include`, append `"changelog"` to `expand`.
4. If `use_display_names`, append `"names"` to `expand`.
5. Call `await run_jira_fetcher_call(jira.get_issue, ...)`.
6. If `use_display_names`, `issue.to_display_name_dict(...)`; else
   `issue.to_simplified_dict()`.
7. Set `result["comments"] = []` and `result["changelogs"] = []` defaults.
8. For each requested `include` section:
   - `remote_links` → `jira.get_remote_issue_links(issue_key)`; on error → `[]`.
   - `transitions` → `jira.get_available_transitions(issue_key)`; on error → `[]`.
   - `watchers` → `jira.get_issue_watchers(issue_key)`; on error → `{}`.
   - `worklogs` → `jira.get_worklogs(issue_key)`; on error → `[]`.
9. JSON-serialize `result`.

### Output

```json
{
  "id": "10001",
  "key": "PROJ-123",
  "self": "https://...",
  "fields": { ... },
  "renderedFields": null,
  "changelog": { "startAt": 0, "maxResults": 100, "total": 5, "histories": [...] },
  "transitions": [...],
  "watchers": {"watchCount": 3, "watchers": [...]},
  "comments": [...],
  "worklogs": [...],
  "remote_links": [...],
  "comments": [...],
  "changelogs": [...]
}
```

### Underlying REST calls

- `GET /issue/{issueKey}?fields=&expand=&properties=&updateHistory=`
- For each `include` section, a separate endpoint:
  - `remote_links` → `GET /issue/{issueKey}/remotelink`
  - `transitions` → `GET /issue/{issueKey}/transitions`
  - `watchers` → `GET /issue/{issueKey}/watchers`
  - `worklogs` → `GET /issue/{issueKey}/worklog`
- `changelog` is included via the `expand=changelog` parameter on the
  main issue fetch.

### Variants

- `use_display_names=True` causes the response to use `"Story Points"`
  instead of `"customfield_10243"` (the underlying Atlassian expansion
  `names` provides the mapping).

---

## `search` (jira_search)

**Tags:** `jira`, `read`, `toolset:jira_issues`
**Annotations:** `title="Search Issues"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:738-857`

### Purpose

Search Jira issues using JQL with pagination.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `jql` | `str` | yes | – | JQL query string (e.g. `'project = PROJ AND status = "In Progress"'`) |
| `fields` | `str` | no | `DEFAULT_READ_JIRA_FIELDS` | Comma-separated fields, or `"*all"` |
| `limit` | `int` | no | `10` | Max results. `ge=1`. |
| `start_at` | `int` | no | `0` | Pagination offset (Server/DC). `ge=0`. |
| `projects_filter` | `str \| None` | no | `None` | Comma-separated project keys, overrides `JIRA_PROJECTS_FILTER` |
| `expand` | `str \| None` | no | `None` | Expand parameter |
| `page_token` | `str \| None` | no | `None` | Cloud pagination token from a previous response |
| `use_display_names` | `bool` | no | `False` | Translate custom field keys to display names |

### Handler logic

1. Convert `fields` CSV → list (or pass `*all`).
2. If `use_display_names`, append `"names"` to `expand`.
3. Call `await run_jira_fetcher_call(jira.search_issues, jql, fields, limit, start, expand, projects_filter, page_token)`.
4. `result = search_result.to_simplified_dict()` (or
   `to_display_name_dict()`).
5. JSON-serialize `result`.

### Output

```json
{
  "issues": [...],
  "startAt": 0,
  "maxResults": 10,
  "total": 42,
  "nextPageToken": "..."   // Cloud only, when there are more pages
}
```

### Underlying REST calls

- Cloud: `GET /search?jql=&fields=&startAt=&maxResults=&expand=&nextPageToken=`
- Server/DC: `POST /search` with body `{"jql", "startAt", "maxResults", "fields", "expand"}`.

### Variants

- `page_token` is Cloud-only.
- `projects_filter` overrides the `JIRA_PROJECTS_FILTER` env var for a
  single call.

---

## `get_project_issues` (jira_get_project_issues)

**Tags:** `jira`, `read`, `toolset:jira_issues`
**Annotations:** `title="Get Project Issues"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1078-1116`

### Purpose

Get issues for a given project (paginated). Equivalent to a JQL search
with `project = <KEY>`.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `project_key` | `str` | yes | – | Project key (pattern `^[A-Z][A-Z0-9_]+$`) |
| `limit` | `int` | no | `10` | Max results. `ge=1, le=50`. |
| `start_at` | `int` | no | `0` | Pagination offset |

### Handler logic

`jira.get_project_issues(project_key=…, start=start_at, limit=limit)` →
JSON-serialize `to_simplified_dict()`.

---

## `create_issue` (jira_create_issue)

**Tags:** `jira`, `write`, `toolset:jira_issues`
**Annotations:** `title="Create Issue"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:1706-1817`

### Purpose

Create a new Jira issue (and optionally link it to an Epic or set it as
a Subtask of an existing issue).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `project_key` | `str` | yes | – | Project key |
| `summary` | `str` | yes | – | Issue summary |
| `issue_type` | `str` | yes | – | Issue type (e.g. `'Task'`, `'Bug'`, `'Story'`, `'Epic'`, `'Subtask'`) |
| `assignee` | `str \| None` | no | `None` | Email, display name, or account ID |
| `description` | `str \| None` | no | `None` | Issue description (Markdown) |
| `components` | `str \| None` | no | `None` | Comma-separated component names |
| `additional_fields` | `str \| None` | no | `None` | JSON string of additional fields (priority, labels, parent, fixVersions, custom fields, epicKey / epic_link) |

### Handler logic

1. `@check_write_access`.
2. Split `components` CSV → list.
3. Parse `additional_fields` JSON → dict.
4. `issue = jira.create_issue(project_key, summary, issue_type, description, assignee, components_list, **extra_fields)`.
5. JSON-serialize `{"message": "Issue created successfully", "issue": issue.to_simplified_dict()}`.

### Underlying REST call

`POST /issue` with body:

```json
{
  "fields": {
    "project": {"key": "PROJ"},
    "summary": "...",
    "issuetype": {"name": "Task"},
    "assignee": {"accountId": "..."},
    "description": "...markdown...",
    "components": [{"name": "Frontend"}, {"name": "API"}],
    "priority": {"name": "High"},
    "labels": ["frontend", "urgent"],
    "parent": {"key": "PROJ-456"},
    "customfield_10010": "...",
    "fixVersions": [{"id": "10020"}],
    "<epic-field>": "EPIC-123"
  }
}
```

### Variants

- For Epic link: use `additional_fields = '{"epicKey": "EPIC-123"}'`
  (the mixin auto-detects Cloud vs DC and sets the right field name).
- For Subtask: use `additional_fields = '{"parent": "PROJ-123"}'`.

---

## `batch_create_issues` (jira_batch_create_issues)

**Tags:** `jira`, `write`, `toolset:jira_issues`
**Annotations:** `title="Batch Create Issues"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:1820-1889`

### Purpose

Create multiple issues in a single batch.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issues` | `str` | yes | – | JSON array of issue objects (each with `project_key`, `summary`, `issue_type`, optional `description`, `assignee`, `components`) |
| `validate_only` | `bool` | no | `False` | If true, only validates without creating |

### Handler logic

1. `@check_write_access`.
2. Parse `issues` JSON → list (raises `ValueError` if not a list).
3. `created = jira.batch_create_issues(issues_list, validate_only=validate_only)`.
4. `{"message": "Issues created successfully", "issues": [i.to_simplified_dict() for i in created]}`.

### Underlying REST call

`POST /issue/bulk` with body `{"issueUpdates": [...]}` (Cloud).
Server/DC: same endpoint, sometimes paginated.

---

## `batch_get_changelogs` (jira_batch_get_changelogs)

**Tags:** `jira`, `read`, `toolset:jira_issues`
**Annotations:** `title="Batch Get Changelogs"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1892-1970`

### Purpose

Get changelogs (audit history) for multiple issues in one batch.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_ids_or_keys` | `str` | yes | – | Comma-separated list (e.g. `'PROJ-123,PROJ-124'`) |
| `fields` | `str \| None` | no | `None` | Comma-separated field names to filter changelogs |
| `limit` | `int` | no | `-1` | Per-issue cap. `-1` = all. |

### Handler logic

1. Raise `NotImplementedError` if `not jira.config.is_cloud` — Server/DC
   doesn't expose this endpoint.
2. Parse CSV → list of keys, optional field filter list.
3. `issues_with_changelogs = jira.batch_get_changelogs(keys_list, fields_list)`.
4. For each, slice `changelogs` to `limit`.
5. JSON-serialize `[{issue_id, changelogs: [...]}]`.

### Underlying REST call

`GET /issue/{idOrKey}/changelog` in a loop (one request per issue) with
optional `fields` query param.

---

## `update_issue` (jira_update_issue)

**Tags:** `jira`, `write`, `toolset:jira_issues`
**Annotations:** `title="Update Issue"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:1973-2126`

### Purpose

Update an existing issue's fields, components, attachments. Supports
Epic linking via `additional_fields.epicKey` / `epic_link`.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `fields` | `str` | yes | – | JSON string of fields to update (`assignee`, `summary`, `description`, etc.) |
| `additional_fields` | `str \| None` | no | `None` | JSON string of additional (custom / epic) fields |
| `components` | `str \| None` | no | `None` | Comma-separated component names |
| `attachments` | `str \| None` | no | `None` | JSON array string or comma-separated file paths |
| `return_fields` | `str` | no | `"*all"` | Comma-separated field list, or `*all` for the full updated issue |

### Handler logic

1. `@check_write_access`.
2. Parse `fields` JSON, `additional_fields` JSON.
3. Parse `components` CSV.
4. Parse `attachments`: JSON array first, CSV fallback.
5. Merge: `{**fields, **extra_fields, components?, attachments?}`.
6. `issue = jira.update_issue(issue_key=…, return_fields=…, **updates)`.
7. JSON-serialize `{"message": "Issue updated successfully", "issue": result}`.
   Includes `attachment_results` from the underlying issue's `custom_fields`
   if present.

### Underlying REST call

`PUT /issue/{issueKey}` with body `{fields: {...}, update?: {...}}`.
File attachments are uploaded via `POST /issue/{issueKey}/attachments`
multipart.

### Variants

- `return_fields="*all"` is the default and returns the full updated issue
  (potentially huge). Token-saving: pass a small list like
  `"status,duedate"` or just `"status"`.

---

## `assign_issue` (jira_assign_issue)

**Tags:** `jira`, `write`, `toolset:jira_issues`
**Annotations:** `title="Assign Issue"`, `readOnlyHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2129-2193`

### Purpose

Assign an issue via the dedicated `PUT /issue/{key}/assignee` endpoint,
which is more reliable than setting the assignee field via `update_issue`
(some Jira configurations silently ignore the latter).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `assignee` | `str \| None` | no | `None` | User identifier (email, display name, account ID), or a JSON object string from `search_assignable_users`. `None` / `""` → unassign. |

### Handler logic

1. `@check_write_access`.
2. If `assignee` starts with `{`, attempt `json.loads` → dict.
3. `issue = jira.assign_issue(issue_key=…, assignee=parsed)`.
4. JSON-serialize `{"message": "Issue <key> assigned successfully", "issue": result}`.

### Underlying REST call

`PUT /issue/{issueKey}/assignee` with body:

```json
{"accountId": "..."}   // Cloud
{"name": "..."}        // Server/DC
```

Pass `null` (or `-1` for account ID) to unassign.

---

## `delete_issue` (jira_delete_issue)

**Tags:** `jira`, `write`, `toolset:jira_issues`
**Annotations:** `title="Delete Issue"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2196-2227`

### Purpose

Delete an existing issue.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |

### Handler logic

1. `@check_write_access`.
2. `jira.delete_issue(issue_key)` (raises on failure).
3. `{"message": "Issue <key> has been deleted successfully."}`.

### Underlying REST call

`DELETE /issue/{issueKey}?deleteSubtasks=true` (default).

---

## `move_issue` (jira_move_issue)

**Tags:** `jira`, `write`, `toolset:jira_issues`
**Annotations:** `title="Move Issue to Project"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2230-2299`

### Purpose

Move an issue to a different project (Cloud-only bulk-move API). The
issue keeps its current issue type and may be assigned a new key in
the target project. Polls until confirmed or times out after 30 seconds.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Source issue key |
| `target_project_key` | `str` | yes | Target project key |

### Handler logic

1. `@check_write_access`.
2. `result = await asyncio.to_thread(jira.move_issue, issue_key, target_project_key)`.
3. `{"message": "...", "issue": result.to_simplified_dict()}`.

### Underlying REST call

`POST /bulk/assign/move` with body `{"issues": [...], "target": {"projectKey": ...}}`.
Then poll `GET /bulk/assign/move/{taskId}` until done.

### Variants

- **Cloud only**. On Server/DC, the lib raises `NotImplementedError`.