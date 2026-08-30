# Jira tools — `toolset:jira_service_desk`

Jira Service Management / Service Desk operations. Most endpoints are
Server/DC only.

## `get_service_desk_for_project` (jira_get_service_desk_for_project)

**Tags:** `jira`, `read`, `toolset:jira_service_desk`
**Annotations:** `title="Get Service Desk For Project"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3371-3409`

### Purpose

Get the Jira Service Desk associated with a project. **Server/DC only**
— raises `NotImplementedError` on Cloud.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `project_key` | `str` | yes | Project key |

### Underlying REST call

`GET /servicedeskapi/servicedesk/by-project/{projectKey}`

---

## `get_service_desk_queues` (jira_get_service_desk_queues)

**Tags:** `jira`, `read`, `toolset:jira_service_desk`
**Annotations:** `title="Get Service Desk Queues"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3412-3455`

### Purpose

List queues for a Service Desk. **Server/DC only.**

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `service_desk_id` | `str` | yes | – | Service desk ID |
| `start_at` | `int` | no | `0` | Pagination offset |
| `limit` | `int` | no | `50` | Max results. `ge=1, le=50`. |

### Underlying REST call

`GET /servicedeskapi/servicedesk/{serviceDeskId}/queue` with `start`, `limit`, `includeCount`.

---

## `get_queue_issues` (jira_get_queue_issues)

**Tags:** `jira`, `read`, `toolset:jira_service_desk`
**Annotations:** `title="Get Queue Issues"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3458-3506`

### Purpose

List issues in a Service Desk queue. **Server/DC only.**

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `service_desk_id` | `str` | yes | – | Service desk ID |
| `queue_id` | `str` | yes | – | Queue ID |
| `start_at` | `int` | no | `0` | Pagination offset |
| `limit` | `int` | no | `50` | Max results. `ge=1`. |

### Underlying REST call

`GET /servicedeskapi/servicedesk/{serviceDeskId}/queue/{queueId}/issue`

---

## `get_request_types` (jira_get_request_types)

**Tags:** `jira`, `read`, `toolset:jira_service_desk`
**Annotations:** `title="Get Request Types"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3509-3545`

### Purpose

List request types for a JSM Service Desk (works on Cloud and DC).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `service_desk_id` | `str` | yes | – | Service desk ID |
| `start_at` | `int` | no | `0` | Pagination offset |
| `limit` | `int` | no | `50` | Max results. `ge=1, le=50`. |

### Underlying REST call

`GET /servicedeskapi/servicedesk/{serviceDeskId}/requesttype`

---

## `get_request_type_fields` (jira_get_request_type_fields)

**Tags:** `jira`, `read`, `toolset:jira_service_desk`
**Annotations:** `title="Get Request Type Fields"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:3548-3578`

### Purpose

Get field definitions for a JSM request type.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `service_desk_id` | `str` | yes | Service desk ID |
| `request_type_id` | `str` | yes | Request type ID |

### Underlying REST call

`GET /servicedeskapi/servicedesk/{serviceDeskId}/requesttype/{requestTypeId}/field`

---

## `create_customer_request` (jira_create_customer_request)

**Tags:** `jira`, `write`, `toolset:jira_service_desk`
**Annotations:** `title="Create Customer Request"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/jira.py:3581-3677`

### Purpose

Create a JSM customer request, optionally on behalf of another user.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `service_desk_id` | `str` | yes | – | Service desk ID |
| `request_type_id` | `str` | yes | – | Request type ID |
| `request_field_values` | `str` | yes | – | JSON object string keyed by field ID |
| `raise_on_behalf_of` | `str \| None` | no | `None` | Jira user identifier |
| `request_participants` | `str \| None` | no | `None` | JSON array or CSV of participants |
| `attachments` | `str \| None` | no | `None` | JSON array of `{filename, mime_type, base64}` |
| `strict_on_behalf` | `bool` | no | `False` | Fail vs retry if `raiseOnBehalfOf` is unsupported |

### Handler logic

1. `@check_write_access`.
2. Parse `request_field_values` JSON.
3. Parse `request_participants`: JSON array first, CSV fallback.
4. Parse `attachments`: JSON array of `{filename, mime_type, base64}`.
5. `result = jira.create_customer_request(...)`.
6. JSON-serialize `to_simplified_dict()`.

### Underlying REST call

`POST /servicedeskapi/request` (Cloud) or `/rest/servicedeskapi/request`
(DC) with body:

```json
{
  "serviceDeskId": "4",
  "requestTypeId": "23",
  "requestFieldValues": {"summary": "...", "description": "..."},
  "raiseOnBehalfOf": "usernameOrAccountId",
  "requestParticipants": ["user1", "user2"]
}
```

### Variants

- `strict_on_behalf=True` → raises if `raiseOnBehalfOf` can't be applied;
  `False` → retries without on-behalf and returns a fallback warning.