# Data models

The `models/` package defines the data shapes that flow out of every
tool. The Go port doesn't need Pydantic — Go structs with `json` tags
and explicit unmarshal methods are the right pattern.

## Base classes (`models/base.py`)

### `ApiModel`

```python
class ApiModel(BaseModel):
    @classmethod
    def from_api_response(cls, data: dict, **kwargs) -> cls:
        raise NotImplementedError

    def to_simplified_dict(self) -> dict[str, Any]:
        return self.model_dump(mode="json", exclude_none=True)
```

The contract:

1. **`from_api_response(cls, data, **kwargs) -> cls`**: subclass-defined
   constructor that takes the raw Atlassian REST response dict and
   returns a model instance. This is the *only* sanctioned way to build
   a model from API data.
2. **`to_simplified_dict(self) -> dict`**: subclass-inherited
   serializer. Returns a JSON-safe dict suitable for `json.dumps(...)`.
   `None` fields are dropped.

### `TimestampMixin`

```python
class TimestampMixin:
    @staticmethod
    def format_timestamp(timestamp: str | None) -> str:
        if not timestamp:
            return ""
        try:
            dt = parse_date(timestamp)
            return dt.strftime("%Y-%m-%d %H:%M:%S %Z") if dt else timestamp
        except Exception:
            return timestamp

    @staticmethod
    def is_valid_timestamp(timestamp: str | None) -> bool:
        if not timestamp:
            return False
        return parse_date(timestamp) is not None
```

`format_timestamp` returns `"%Y-%m-%d %H:%M:%S %Z"` on success, the
original string on failure, or `""` for falsy input.

## Jira models

Located in `src/mcp_atlassian/models/jira/`.

| Class | File | Notes |
| --- | --- | --- |
| `JiraIssue` | `issue.py` | The central object. Fields: `id`, `key`, `self`, `fields` (dict of `JiraField` instances). The `to_simplified_dict` returns a flat dict with key fields (summary, status, assignee, etc.) for token efficiency. |
| `JiraComment` | `comment.py` | `id`, `body`, `author`, `created`, `updated`. |
| `JiraWorklog` | `worklog.py` | `id`, `timeSpent`, `timeSpentSeconds`, `started`, `comment`. |
| `JiraUser` | `common.py` | Cloud: `accountId`. DC: `name`/`key`. |
| `JiraStatus`, `JiraPriority`, `JiraIssueType` | `common.py` | Always have `id` + `name`. Defaults: `{name: "Unknown", id: "0"}`. |
| `JiraProject` | `project.py` | `id`, `key`, `name`, `projectTypeKey`. |
| `JiraAttachment` | `attachment.py` | `id`, `filename`, `content_type`, `size`, `url` (download URL). |
| `JiraVersion` | `version.py` | `id`, `name`, `startDate`, `releaseDate`, `archived`, `released`. |
| `JiraComponent` | `component.py` | `id`, `name`, `description`. |
| `JiraTransition` | `transitions.py` | `id`, `name`, `to`, `hasScreen`, `isGlobal`, `fields`. |
| `JiraRemoteLink` | `links.py` | `id`, `object: {url, title, icon, ...}`, `relationship`. |
| `JiraIssueLink` | `links.py` | `id`, `type`, `inwardIssue`, `outwardIssue`. |
| `JiraSprint` | `sprint.py` | `id`, `name`, `state`, `startDate`, `endDate`, `goal`. |
| `JiraBoard` | `board.py` | `id`, `name`, `type`. |
| `JiraChangelog` | `changelog.py` | `id`, `author`, `created`, `items` (list of `{field, fromString, toString}`). |

### Cloud-vs-DC field divergence

| Field | Cloud | DC |
| --- | --- | --- |
| User identifier | `accountId` | `name`/`key` |
| Mention syntax | `[~accountid:...]` | `[~username]` or `[~userkey]` |
| Epic link field | `customfield_12311140` | `Epic Link` (configurable) |
| Default fields in `*all` | includes all custom fields | varies |

## Confluence models

Located in `src/mcp_atlassian/models/confluence/`.

| Class | File | Notes |
| --- | --- | --- |
| `ConfluencePage` | `page.py` | `id`, `title`, `type`, `space`, `version`, `body` (dict: `view`, `storage`, etc.), `history`. |
| `ConfluenceComment` | `comment.py` | `id`, `body`, `author`, `created`. |
| `ConfluenceSpace` | `space.py` | `id`, `key`, `name`, `type`. |
| `ConfluenceAttachment` | `attachment.py` | `id`, `title`, `mediaType`, `fileSize`, `downloadUrl`. |
| `ConfluenceLabel` | `label.py` | `prefix`, `name`, `id`. |
| `ConfluenceUser` | `user.py` | Cloud: `accountId`, `displayName`, `email`. DC: `username`, `userKey`. |
| `ConfluenceVersion` | `version.py` | `number`, `when`, `by`. |
| `ConfluenceRestriction` | `restrictions.py` | `read` / `update` operations with user/group lists. |

## Constants (`models/constants.py`)

```python
EMPTY_STRING = ""
UNKNOWN = "Unknown"
UNASSIGNED = "Unassigned"
NONE_VALUE = "None"

JIRA_DEFAULT_ID = "0"
JIRA_DEFAULT_KEY = "UNKNOWN-0"
JIRA_DEFAULT_STATUS = {"name": UNKNOWN, "id": "0"}
JIRA_DEFAULT_PRIORITY = {"name": "None", "id": "0"}
JIRA_DEFAULT_ISSUE_TYPE = {"name": UNKNOWN, "id": "0"}
JIRA_DEFAULT_PROJECT = "0"

CONFLUENCE_DEFAULT_ID = "0"
CONFLUENCE_DEFAULT_SPACE = {"key": "", "name": UNKNOWN, "id": "0"}
CONFLUENCE_DEFAULT_VERSION = {"number": 0, "when": ""}
DEFAULT_TIMESTAMP = "1970-01-01T00:00:00.000+0000"
```

These are used by model constructors when fields are missing from the
API response. The Go port should preserve these constants so that the
output JSON for malformed API responses is identical.

## Go port notes

- The Go port doesn't need Pydantic; structs with `json` tags work.
- For default values, write custom `UnmarshalJSON` methods that fill in
  the defaults.
- `to_simplified_dict` becomes `ToSimplifiedDict() map[string]any` with
  explicit `nil` skipping.
- `TimestampMixin` becomes a free function:
  ```go
  func FormatTimestamp(s string) string {
      if s == "" { return "" }
      if t, err := ParseDate(s); err == nil {
          return t.Format("2006-01-02 15:04:05 MST")
      }
      return s
  }
  ```
- `parse_date` (utils/date.py) becomes `time.Parse(time.RFC3339, s)` plus
  a few variants; epoch-ms handling via `time.UnixMilli(n)`.