# Jira tools — `toolset:jira_agile`

Boards, sprints, and agile operations.

## `get_agile_boards` (jira_get_agile_boards)

**Tags:** `jira`, `read`, `toolset:jira_agile`
**Annotations:** `title="Get Agile Boards"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1443-1497`

### Purpose

List agile boards (Scrum / Kanban), optionally filtered by name, project,
or type.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `board_name` | `str \| None` | no | `None` | Fuzzy-search by name |
| `project_key` | `str \| None` | no | `None` | Filter to project |
| `board_type` | `str \| None` | no | `None` | `'scrum'` or `'kanban'` |
| `start_at` | `int` | no | `0` | Pagination offset |
| `limit` | `int` | no | `10` | Max results. `ge=1, le=50`. |

### Handler logic

`boards = jira.get_all_agile_boards_model(board_name, project_key, board_type, start, limit)` →
`[b.to_simplified_dict() for b in boards]` → JSON.

### Underlying REST call

`GET /agile/1.0/board` with `name`, `projectKeyOrId`, `type`, `startAt`, `maxResults`.

---

## `get_board_issues` (jira_get_board_issues)

**Tags:** `jira`, `read`, `toolset:jira_agile`
**Annotations:** `title="Get Board Issues"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1500-1577`

### Purpose

Search issues on a specific board using JQL.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `board_id` | `str` | yes | – | Board ID (e.g. `'1001'`) |
| `jql` | `str` | yes | – | JQL filter |
| `fields` | `str` | no | `DEFAULT_READ_JIRA_FIELDS` | Comma-separated fields |
| `start_at` | `int` | no | `0` | Pagination offset |
| `limit` | `int` | no | `10` | Max results. `ge=1, le=50`. |
| `expand` | `str` | no | `'version'` | Expand parameter |

### Handler logic

`jira.get_board_issues(board_id, jql, fields_list, start, limit, expand)` →
JSON-serialize `to_simplified_dict()`.

### Underlying REST call

`GET /agile/1.0/board/{boardId}/issue` with `jql`, `fields`, `startAt`, `maxResults`, `expand`.

---

## `get_sprints_from_board` (jira_get_sprints_from_board)

**Tags:** `jira`, `read`, `toolset:jira_agile`
**Annotations:** `title="Get Sprints from Board"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1580-1617`

### Purpose

List sprints on a board, optionally filtered by state.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `board_id` | `str` | yes | – | Board ID (e.g. `'1000'`) |
| `state` | `str \| None` | no | `None` | `'active'`, `'future'`, `'closed'` |
| `start_at` | `int` | no | `0` | Pagination offset |
| `limit` | `int` | no | `10` | Max results. `ge=1, le=50`. |

### Underlying REST call

`GET /agile/1.0/board/{boardId}/sprint` with `state`, `startAt`, `maxResults`.

---

## `get_sprint_issues` (jira_get_sprint_issues)

**Tags:** `jira`, `read`, `toolset:jira_agile`
**Annotations:** `title="Get Sprint Issues"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1620-1668`

### Purpose

Search issues inside a sprint.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `sprint_id` | `str` | yes | – | Sprint ID (e.g. `'10001'`) |
| `fields` | `str` | no | `DEFAULT_READ_JIRA_FIELDS` | Comma-separated fields |
| `start_at` | `int` | no | `0` | Pagination offset |
| `limit` | `int` | no | `10` | Max results. `ge=1, le=50`. |

### Underlying REST call

`GET /agile/1.0/sprint/{sprintId}/issue` with `fields`, `startAt`, `maxResults`, `expand`.

---

## `create_sprint` (jira_create_sprint)

**Tags:** `jira`, `write`, `toolset:jira_agile`
**Annotations:** `title="Create Sprint"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2854-2897`

### Purpose

Create a new sprint on a board.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `board_id` | `str` | yes | Board ID |
| `name` | `str` | yes | Sprint name |
| `start_date` | `str` | yes | ISO 8601 start date |
| `end_date` | `str` | yes | ISO 8601 end date |
| `goal` | `str \| None` | no | Sprint goal |

### Underlying REST call

`POST /agile/1.0/sprint` (or `/agile/1.0/board/{boardId}/sprint`) with
body `{name, startDate, endDate, goal, originBoardId}`.

---

## `update_sprint` (jira_update_sprint)

**Tags:** `jira`, `write`, `toolset:jira_agile`
**Annotations:** `title="Update Sprint"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2900-2958`

### Purpose

Update an existing sprint's name, state, dates, or goal.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `sprint_id` | `str` | yes | – | Sprint ID |
| `name` | `str \| None` | no | `None` | New name |
| `state` | `str \| None` | no | `None` | `'future'`, `'active'`, `'closed'` |
| `start_date` | `str \| None` | no | `None` | New start date |
| `end_date` | `str \| None` | no | `None` | New end date |
| `goal` | `str \| None` | no | `None` | New goal |

### Underlying REST call

`PUT /agile/1.0/sprint/{sprintId}` with body `{name?, state?, startDate?, endDate?, goal?, completeDate?}`.

---

## `add_issues_to_sprint` (jira_add_issues_to_sprint)

**Tags:** `jira`, `write`, `toolset:jira_agile`
**Annotations:** `title="Add Issues to Sprint"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2961-2995`

### Purpose

Move issues into a sprint.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `sprint_id` | `str` | yes | Sprint ID |
| `issue_keys` | `str` | yes | Comma-separated issue keys |

### Underlying REST call

`POST /agile/1.0/sprint/{sprintId}/issue` with body `{issues: [key1, key2, …]}`.

---

## `move_issues_to_backlog` (jira_move_issues_to_backlog)

**Tags:** `jira`, `write`, `toolset:jira_agile`
**Annotations:** `title="Move Issues to Backlog"`, `readOnlyHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2998-3029`

### Purpose

Move issues to the backlog (removes them from any sprint).

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_keys` | `str` | yes | Comma-separated issue keys |

### Underlying REST call

`POST /agile/1.0/backlog/issue` with body `{issues: [key1, key2, …]}`.