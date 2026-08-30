# Jira tools — `toolset:jira_project_analysis`

Higher-order analyses over issues in a project.

## `get_project_epic_hierarchy` (jira_get_project_epic_hierarchy)

**Tags:** `jira`, `read`, `toolset:jira_project_analysis`
**Annotations:** `title="Get Project Epic Hierarchy"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:4413-4459`

### Purpose

Fetch all epics in a project, detect their parent issue (via the
`parent` field or inward issue links), and group them by parent.
Epics with no detected parent appear under `"Unlinked"`.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `project_key` | `str` | yes | – | Project key |
| `max_epics` | `int` | no | `200` | Max epics to fetch. `ge=1, le=500`. |

### Output

```json
{
  "project_key": "PROJ",
  "epics": [
    {
      "epic": {"key": "PROJ-100", "summary": "...", ...},
      "parent": {"key": "OTHER-1", "summary": "...", "project_key": "OTHER", "via": "parent_field"|"inward_link"},
      "children": [...]    // other epics under the same parent
    },
    {
      "epic": {...},
      "parent": null,
      "parent_group": "Unlinked"
    }
  ]
}
```

### Underlying REST calls

- `GET /search?jql=project=KEY AND issuetype=Epic` (paginated up to
  `max_epics`).
- For each epic, `GET /issue/{key}?fields=parent,issuelinks` to find the
  parent / inward link.

---

## `get_cross_project_dependencies` (jira_get_cross_project_dependencies)

**Tags:** `jira`, `read`, `toolset:jira_project_analysis`
**Annotations:** `title="Get Cross-Project Dependencies"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:4462-4507`

### Purpose

Scan issues in a project, extract every issue link whose target
belongs to a different project, and group them by target project and
link type.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `project_key` | `str` | yes | – | Project key |
| `max_issues` | `int` | no | `200` | Max issues to scan. `ge=1, le=500`. |

### Output

```json
{
  "project_key": "PROJ",
  "dependencies": {
    "OTHER": {
      "Blocks": [
        {"from": "PROJ-1", "to": "OTHER-100", "summary": "..."}
      ],
      "Relates to": [...]
    }
  },
  "total_links": 42
}
```

### Underlying REST calls

- `GET /search?jql=project=KEY ORDER BY created DESC` (paginated up to
  `max_issues`).
- For each issue, inspect its `issuelinks` to find cross-project
  relationships.