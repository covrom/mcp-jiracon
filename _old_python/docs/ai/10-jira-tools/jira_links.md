# Jira tools — `toolset:jira_links`

Issue link management (incl. Epic link and remote links).

## `get_link_types` (jira_get_link_types)

**Tags:** `jira`, `read`, `toolset:jira_links`
**Annotations:** `title="Get Link Types"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1671-1703`

### Purpose

List all Jira issue link types (e.g. `'Blocks'`, `'Relates to'`,
`'Duplicate'`), optionally filtered by a case-insensitive name substring.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `name_filter` | `str \| None` | no | `None` | Substring filter |

### Handler logic

`link_types = jira.get_issue_link_types()` →
`[lt.to_simplified_dict() for lt in link_types]` →
optional substring filter → JSON.

### Output

```json
[
  {"id": "10000", "name": "Blocks", "inward": "is blocked by", "outward": "blocks", "self": "..."},
  {"id": "10001", "name": "Duplicate", "inward": "is duplicated by", "outward": "duplicates", "self": "..."}
]
```

### Underlying REST call

`GET /issueLinkType`

---

## `link_to_epic` (jira_link_to_epic)

**Tags:** `jira`, `write`, `toolset:jira_links`
**Annotations:** `title="Link to Epic"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2511-2552`

### Purpose

Link an existing issue to an Epic.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | The issue to link |
| `epic_key` | `str` | yes | The target epic |

### Handler logic

1. `@check_write_access`.
2. `issue = jira.link_issue_to_epic(issue_key, epic_key)`.
3. `{"message": "Issue <key> has been linked to epic <epic>.", "issue": issue.to_simplified_dict()}`.

### Underlying REST call

`PUT /issue/{issueKey}` with body:

```json
{"fields": {"customfield_12311140": "EPIC-123"}}   // Cloud
{"fields": {"Epic Link": "EPIC-123"}}              // DC (greenhopper legacy)
```

The mixin auto-detects which field to use by inspecting the issue's
project.

---

## `create_issue_link` (jira_create_issue_link)

**Tags:** `jira`, `write`, `toolset:jira_links`
**Annotations:** `title="Create Issue Link"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2555-2648`

### Purpose

Create a link between two issues, with optional comment and visibility.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `link_type` | `str` | yes | – | Link type name (e.g. `'Duplicate'`, `'Blocks'`, `'Relates to'`) |
| `inward_issue_key` | `str` | yes | – | The inward issue key |
| `outward_issue_key` | `str` | yes | – | The outward issue key |
| `comment` | `str \| None` | no | `None` | Comment text |
| `comment_visibility` | `str \| None` | no | `None` | JSON string `{"type": "group", "value": "jira-users"}` |

### Handler logic

1. `@check_write_access`.
2. Validate required args.
3. If `comment` is set and either issue belongs to a project in
   `JIRA_INTERNAL_ONLY_PROJECTS` → reject (link comments may be
   customer-visible on JSM).
4. Build the link data dict.
5. `result = jira.create_issue_link(link_data)` → JSON-serialize.

### Underlying REST call

`POST /issueLink` with body:

```json
{
  "type": {"name": "Blocks"},
  "inwardIssue": {"key": "PROJ-1"},
  "outwardIssue": {"key": "PROJ-2"},
  "comment": {
    "body": "...",
    "visibility": {"type": "group", "value": "jira-users"}
  }
}
```

---

## `create_remote_issue_link` (jira_create_remote_issue_link)

**Tags:** `jira`, `write`, `toolset:jira_links`
**Annotations:** `title="Create Remote Issue Link"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2651-2736`

### Purpose

Create a web/Confluence remote link on an issue (the URL appears in the
issue's "Links" section).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Issue key |
| `url` | `str` | yes | – | Target URL |
| `title` | `str` | yes | – | Display title |
| `summary` | `str \| None` | no | `None` | Link description |
| `relationship` | `str \| None` | no | `None` | e.g. `'causes'`, `'relates to'`, `'documentation'` |
| `icon_url` | `str \| None` | no | `None` | 16×16 icon URL |

### Underlying REST call

`POST /issue/{issueKey}/remotelink` with body:

```json
{
  "object": {
    "url": "https://example.com/page",
    "title": "Documentation",
    "summary": "Optional summary",
    "icon": {"url16x16": "https://example.com/favicon.ico", "title": "..."}
  },
  "relationship": "documentation"
}
```

---

## `remove_issue_link` (jira_remove_issue_link)

**Tags:** `jira`, `write`, `toolset:jira_links`
**Annotations:** `title="Remove Issue Link"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:2739-2765`

### Purpose

Delete an issue link by ID.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `link_id` | `str` | yes | Link ID |

### Underlying REST call

`DELETE /issueLink/{linkId}`