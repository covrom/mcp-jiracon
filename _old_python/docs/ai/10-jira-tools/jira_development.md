# Jira tools — `toolset:jira_development`

Dev panel integration: PRs, branches, commits linked to an issue.

## `get_issue_development_info` (jira_get_issue_development_info)

**Tags:** `jira`, `read`, `development`, `toolset:jira_development`
**Annotations:** `title="Get Issue Development Info"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:4287-4346`

### Purpose

Get development panel data (pull requests, branches, commits,
repositories) linked to a Jira issue from Bitbucket / GitHub / GitLab /
etc.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_key` | `str` | yes | – | Jira issue key |
| `application_type` | `str \| None` | no | `None` | Filter by source-control type. Case-sensitive. Examples: `'stash'` (Bitbucket Server), `'bitbucket'`, `'GitHub'`, `'githube'` (GH Enterprise), `'GitLab'`. |
| `data_type` | `str \| None` | no | `None` | Filter by data type: `'pullrequest'`, `'branch'`, `'repository'`. |

### Underlying REST call

`POST /rest/devinfo/0.10/bulkByIssue` with body:

```json
{
  "issueIds": [{"id": "PROJ-123", "type": "ISSUE_KEY"}],
  "applicationTypes": ["stash"],
  "dataTypes": ["pullrequest"],
  "_fields": ["pullrequests", "branches", "commits", "repositories"]
}
```

### Output

```json
{
  "pullRequests": [...],
  "branches": [...],
  "commits": [...],
  "repositories": [...]
}
```

### Variants

- The Atlassian Connect Dev Info endpoint was deprecated in 2024; new
  Cloud deployments need to use the Bitbucket Cloud integration's API
  directly. Server/DC instances still work as before.

---

## `get_issues_development_info` (jira_get_issues_development_info)

**Tags:** `jira`, `read`, `development`, `toolset:jira_development`
**Annotations:** `title="Get Issues Development Info"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:4349-4410`

### Purpose

Batch version of the above.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `issue_keys` | `str` | yes | – | Comma-separated issue keys |
| `application_type` | `str \| None` | no | `None` | Same as above |
| `data_type` | `str \| None` | no | `None` | Same as above |

### Handler logic

Same as above, but `issueIds` contains every key.