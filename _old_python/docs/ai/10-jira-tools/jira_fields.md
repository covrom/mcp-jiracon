# Jira tools — `toolset:jira_fields`

Jira field metadata and option lookups. Default-enabled toolset.

## `search_fields` (jira_search_fields)

**Tags:** `jira`, `read`, `toolset:jira_fields`
**Annotations:** `title="Search Fields"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:860-894`

### Purpose

Fuzzy-search Jira field definitions. Use to find the field ID
(`customfield_XXXXX`) before calling `update_issue`.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `keyword` | `str` | no | `""` | Empty = first `limit` fields in default order |
| `limit` | `int` | no | `10` | Max results. `ge=1`. |
| `refresh` | `bool` | no | `False` | Force-refresh the cached field list |

### Handler logic

`result = jira.search_fields(keyword, limit, refresh)` → JSON.

### Output

```json
[
  {"id": "customfield_10010", "name": "Story Points", "schema": {"type": "number"}, "custom": true},
  ...
]
```

### Underlying REST call

`GET /field` (cached client-side; `refresh=True` refetches).

---

## `get_field_options` (jira_get_field_options)

**Tags:** `jira`, `read`, `toolset:jira_fields`
**Annotations:** `title="Get Field Options"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:975-1075`

### Purpose

Get allowed option values for select/multi-select/radio/checkbox/cascading
custom fields. Cloud uses the Field Context Option API; Server/DC uses
`createmeta`'s `allowedValues`.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `field_id` | `str` | yes | – | Custom field ID (e.g. `'customfield_10001'`) |
| `context_id` | `str \| None` | no | `None` | Field context ID (Cloud only; auto-resolves global) |
| `project_key` | `str \| None` | no | `None` | Required for Server/DC |
| `issue_type` | `str \| None` | no | `None` | Required for Server/DC |
| `contains` | `str \| None` | no | `None` | Case-insensitive substring filter (also matches children in cascading selects) |
| `return_limit` | `int \| None` | no | `None` | Cap on results after filtering. `ge=1`. |
| `values_only` | `bool` | no | `False` | Return compact `["value"]` / `[{value, children}]` |

### Handler logic

1. `options = jira.get_field_options(field_id, context_id, project_key, issue_type)`.
2. Convert to simplified dicts.
3. Apply `contains` filter (case-insensitive, walks `child_options` for
   cascading selects).
4. Apply `return_limit` cap.
5. If `values_only`, return the compact payload:
   - Single-level: `["value1", "value2"]`.
   - Cascading: `[{"value": "parent", "children": ["child1", ...]}]`.
6. Otherwise return the full simplified dict list.

### Underlying REST call

- Cloud: `GET /field/{fieldId}/context/{contextId}/option` with
  `startAt`, `maxResults`.
- Server/DC: uses the createmeta endpoint's `allowedValues`.

### Variants

- `values_only=True` returns the compact payload intended for token
  efficiency when the LLM only needs the option names.