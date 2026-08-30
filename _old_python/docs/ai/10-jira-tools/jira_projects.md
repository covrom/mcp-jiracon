# Jira tools — `toolset:jira_projects`

Project metadata and version management.

## `get_project_issue_types` (jira_get_project_issue_types)

**Tags:** `jira`, `read`, `toolset:jira_projects`
**Annotations:** `title="Get Project Issue Types"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3032-3073`

### Purpose

Get the issue types available in a project (e.g. `'Bug'`, `'Task'`,
`'Story'`, `'Epic'`, `'Subtask'`). Use the returned issue type IDs with
`get_create_fields` to discover what fields each type requires.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `project_key` | `str` | yes | Project key |

### Output

```json
[
  {"id": "10001", "name": "Bug", "description": "A bug...", "subtask": false, "untranslatedName": "Bug"},
  {"id": "10002", "name": "Task", "description": "A task...", "subtask": false},
  ...
]
```

### Underlying REST call

`GET /issue/createmeta/{projectKey}/issuetypes`

---

## `get_create_fields` (jira_get_create_fields)

**Tags:** `jira`, `read`, `toolset:jira_projects`
**Annotations:** `title="Get Create Fields"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3076-3126`

### Purpose

Get metadata of all fields (required and optional) for creating an issue
of a specific type in a specific project.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `project_key` | `str` | yes | Project key |
| `issue_type_id` | `str` | yes | Issue type ID (from `get_project_issue_types`) |

### Output

```json
[
  {"field_id": "summary", "name": "Summary", "required": true, "schema": {...}},
  {"field_id": "customfield_10010", "name": "Story Points", "required": false, "schema": {...}}
]
```

### Underlying REST call

`GET /issue/createmeta/{projectKey}/issuetype/{issueTypeId}`

---

## `get_project_versions` (jira_get_project_versions)

**Tags:** `jira`, `read`, `toolset:jira_projects`
**Annotations:** `title="Get Project Versions"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3129-3146`

### Purpose

Get all fix versions (releases) for a project.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `project_key` | `str` | yes | Project key |

### Underlying REST call

`GET /project/{projectKey}/versions`

---

## `get_project_components` (jira_get_project_components)

**Tags:** `jira`, `read`, `toolset:jira_projects`
**Annotations:** `title="Get Project Components"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3149-3166`

### Underlying REST call

`GET /project/{projectKey}/components`

---

## `get_all_projects` (jira_get_all_projects)

**Tags:** `jira`, `read`, `toolset:jira_projects`
**Annotations:** `title="Get All Projects"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3169-3234`

### Purpose

List all projects accessible to the current user. Keys are uppercased.
Honors `JIRA_PROJECTS_FILTER` env var (filter applied to results).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `include_archived` | `bool` | no | `False` | Include archived projects |

### Handler logic

1. `projects = jira.get_all_projects(include_archived)`.
2. Catch auth/network/ValueError → return `{"success": false, "error": ...}` JSON.
3. Uppercase all project keys.
4. Apply `JIRA_PROJECTS_FILTER` filter (comma-separated keys, uppercased).
5. JSON-serialize.

### Underlying REST call

`GET /project/search` with `expand`, `query`, `maxResults`, `startAt`.

---

## `search_projects` (jira_search_projects)

**Tags:** `jira`, `read`, `toolset:jira_projects`
**Annotations:** `title="Search Projects"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3237-3319`

### Purpose

Search projects by name or key prefix using the picker endpoint (more
efficient than fetching all projects).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `query` | `str` | yes | – | Name or key prefix |
| `max_results` | `int` | no | `20` | Max results. `ge=1, le=50`. |
| `current_project_ids` | `str \| None` | no | `None` | Comma-separated project IDs to exclude |

### Underlying REST call

`GET /project/picker` with `query`, `maxResults`, `currentProjectId`.

---

## `get_project_fields` (jira_get_project_fields)

**Tags:** `jira`, `read`, `toolset:jira_projects`
**Annotations:** `title="Get Project Fields"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3322-3368`

### Purpose

Get the fields available on issues in a project (deduplicated across
the project's issue types).

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `project_key` | `str` | yes | Project key |

### Underlying REST call

`GET /issue/createmeta/{projectKey}/issuetypes` then aggregate.

---

## `create_version` (jira_create_version)

**Tags:** `jira`, `write`, `toolset:jira_projects`
**Annotations:** `title="Create Version"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:3680-3734`

### Purpose

Create a new fix version in a project.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `project_key` | `str` | yes | – | Project key |
| `name` | `str` | yes | – | Version name |
| `start_date` | `str \| None` | no | `None` | `YYYY-MM-DD` |
| `release_date` | `str \| None` | no | `None` | `YYYY-MM-DD` |
| `description` | `str \| None` | no | `None` | Description |

### Underlying REST call

`POST /version` with body `{name, project, startDate?, releaseDate?, description?}`.

---

## `batch_create_versions` (jira_batch_create_versions)

**Tags:** `jira`, `write`, `toolset:jira_projects`
**Annotations:** `title="Batch Create Versions"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Name override:** `name="batch_create_versions"`
**Lines:** `servers/jira.py:3737-3818`

### Purpose

Create multiple versions in a project, with per-item try/catch.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `project_key` | `str` | yes | Project key |
| `versions` | `str` | yes | JSON array of `{name, startDate?, releaseDate?, description?}` |

### Handler logic

1. `@check_write_access`.
2. Parse `versions` JSON.
3. For each item, try `create_project_version`; on success → `{success: True, version: ...}`; on failure → `{success: False, error: ..., input: v}`.

### Underlying REST call

`POST /version` (one per item).

---

## `update_version` (jira_update_version)

**Tags:** `jira`, `write`, `toolset:jira_projects`
**Annotations:** `title="Update Version"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:3821-3891`

### Purpose

Update an existing fix version (archive/unarchive, rename, change
release date).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `version_id` | `str` | yes | – | Numeric ID (e.g. `'10001'`) |
| `name` | `str \| None` | no | `None` | New name |
| `description` | `str \| None` | no | `None` | New description |
| `start_date` | `str \| None` | no | `None` | New start date (`YYYY-MM-DD`) |
| `release_date` | `str \| None` | no | `None` | New release date (`YYYY-MM-DD`) |
| `archived` | `bool \| None` | no | `None` | Archive flag |
| `released` | `bool \| None` | no | `None` | Released flag |

### Underlying REST call

`PUT /version/{versionId}` with partial body.