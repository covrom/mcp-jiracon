# Jira REST API mapping

This document maps each Jira MCP tool (or group of tools) to the underlying
Atlassian REST endpoints. Use it to spec the Go HTTP client.

> **Conventions used in this file**
> - Base URL for Server/DC: `{JIRA_URL}/rest/api/2/...` (or `/rest/api/3/...`
>   if Server version ≥ 9.0 — pick per instance; current code defaults to v2).
> - Base URL for Cloud (Basic / PAT): `{JIRA_URL}/rest/api/3/...` (the lib
>   auto-routes).
> - Base URL for Cloud OAuth: `https://api.atlassian.com/ex/jira/{cloudId}/rest/api/3/...`.
> - `atlassian-python-api` is the wrapper library used; URLs in the table
>   below are what's *actually called* on the wire.

---

## 1. `UsersMixin` (users.py)

| Method (mixin) | Tool(s) | HTTP method | URL template | Query / Body | Cloud vs DC | Auth |
| --- | --- | --- | --- | --- | --- | --- |
| `get_user_profile_by_identifier` | `get_user_profile` | GET | `/user` (Cloud) or `/user?username=...` (DC) | Cloud: query `accountId` from identifier. DC: query `username`/`key`. Expands `groups` and `applicationRoles`. | identical, but identifier field differs | session |
| `search_assignable_users` | `search_assignable_users` | GET | `/user/assignable/search` | `query`, `project` OR `issueKey`, `maxResults`, `startAt` | identical | session |
| `get_current_user_account_id` | (used in fetcher resolution) | GET | `/myself` | – | identical | session |

---

## 2. `SearchMixin` (search.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body | Cloud vs DC |
| --- | --- | --- | --- | --- | --- |
| `search_issues` | `search`, `get_project_issues`, `batch_get_changelogs` (via `_paginate_with_changelogs`) | GET (Cloud) or POST (DC legacy) | `/search` (Cloud) or `/search` (DC POST) | Cloud: `jql`, `startAt`, `maxResults`, `fields`, `expand`, `nextPageToken`. DC: `jql`, `startAt`, `maxResults`, `fields`, `expand`. Body for POST: `{jql, startAt, maxResults, fields, expand}`. | Cloud uses `nextPageToken`; DC uses `startAt`. Body POST only used on DC. |
| `get_project_issues` | `get_project_issues` | GET | `/search` | Same as above with a JQL filter `project=<KEY>`. | identical |
| `get_board_issues` | `get_board_issues` | GET | `/agile/1.0/board/{boardId}/issue` | `jql`, `startAt`, `maxResults`, `fields`, `expand` | identical |
| `get_sprint_issues` | `get_sprint_issues` | GET | `/agile/1.0/sprint/{sprintId}/issue` | `startAt`, `maxResults`, `fields`, `expand` | identical |
| `batch_get_changelogs` | `batch_get_changelogs` | GET | `/issue/{idOrKey}/changelog` (loop) | `startAt`, `maxResults` | **Cloud only** — raises `NotImplementedError` on DC |

---

## 3. `IssuesMixin` (issues.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_issue` | `get_issue` | GET | `/issue/{issueKey}` | `fields`, `expand`, `properties`, `updateHistory` |
| `create_issue` | `create_issue`, `batch_create_issues` | POST | `/issue` (and `/issue/bulk` for batch) | Body: `{fields: {project, summary, issuetype, …}, update: {...}?}` for create. Bulk: `{issueUpdates: [...]}` |
| `update_issue` | `update_issue` | PUT | `/issue/{issueKey}` | Body: `{update: {...}, fields: {...}, historyMetadata?: {...}}`. Supports query `returnFields` (comma-separated). |
| `assign_issue` | `assign_issue` | PUT | `/issue/{issueKey}/assignee` | Body: `{accountId}` (Cloud) or `{name}`/`{key}` (DC). Pass `null` body or `{"accountId": null}` to unassign. |
| `delete_issue` | `delete_issue` | DELETE | `/issue/{issueKey}` | Optional query `deleteSubtasks=true`. |
| `move_issue` | `move_issue` | POST (Cloud bulk) | `/bulk/assign/move` | Body: `{issues: [issueKey], target: {projectKey}}`. Async; tool polls the task. **Cloud only.** |

---

## 4. `CommentsMixin` (comments.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_comments` | (`get_issue` include only) | GET | `/issue/{issueKey}/comment` | `startAt`, `maxResults`, `orderBy` |
| `add_comment` | `add_comment` | POST | `/issue/{issueKey}/comment` | Body: `{body, visibility?}` |
| `edit_comment` | `edit_comment` | PUT | `/issue/{issueKey}/comment/{commentId}` | Body: `{body, visibility?}` |

---

## 5. `WorklogMixin` (worklog.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_worklogs` | `get_worklog` (also `get_issue` include) | GET | `/issue/{issueKey}/worklog` | `startAt`, `maxResults` |
| `add_worklog` | `add_worklog` | POST | `/issue/{issueKey}/worklog` | Body: `{timeSpent, started?, comment?, visibility?, adjustEstimate?, newEstimate?, reduceBy?}` |

---

## 6. `TransitionsMixin` (transitions.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_available_transitions` | `get_transitions` (also `get_issue` include) | GET | `/issue/{issueKey}/transitions` | `expand=transitions.fields` (optional) |
| `transition_issue` | `transition_issue` | POST | `/issue/{issueKey}/transitions` | Body: `{transition: {id: <id>}, fields?: {...}, update?: {...}, historyMetadata?: {...}}` |

---

## 7. `LinksMixin` (links.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_issue_link_types` | `get_link_types` | GET | `/issueLinkType` | – |
| `link_issue_to_epic` | `link_to_epic` | PUT | `/issue/{issueKey}` | Body: `{fields: {"customfield_12311140": <epicKey>}}` (Cloud) or `{"Epic Link": <epicKey>}` (DC). The mixin auto-detects the right field name. |
| `create_issue_link` | `create_issue_link` | POST | `/issueLink` | Body: `{type: {name}, inwardIssue: {key}, outwardIssue: {key}, comment?: {body, visibility?}}` |
| `create_remote_issue_link` | `create_remote_issue_link` | POST | `/issue/{issueKey}/remotelink` | Body: `{object: {url, title, summary?, icon?: {url16x16}}, relationship?}` |
| `remove_issue_link` | `remove_issue_link` | DELETE | `/issueLink/{linkId}` | – |
| `get_remote_issue_links` | (`get_issue` include) | GET | `/issue/{issueKey}/remotelink` | `globalId` optional filter |

---

## 8. `BoardsMixin` (boards.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_all_agile_boards_model` | `get_agile_boards` | GET | `/agile/1.0/board` | `name` (fuzzy), `projectKeyOrId`, `type`, `startAt`, `maxResults` |

---

## 9. `SprintsMixin` (sprints.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_all_sprints_from_board_model` | `get_sprints_from_board` | GET | `/agile/1.0/board/{boardId}/sprint` | `state`, `startAt`, `maxResults` |
| `create_sprint` | `create_sprint` | POST | `/agile/1.0/sprint` (or `/agile/1.0/board/{boardId}/sprint` for board-scoped) | Body: `{name, startDate, endDate, goal?, originBoardId}` |
| `update_sprint` | `update_sprint` | PUT | `/agile/1.0/sprint/{sprintId}` | Body: `{name?, state?, startDate?, endDate?, goal?, completeDate?}` |
| `add_issues_to_sprint` | `add_issues_to_sprint` | POST | `/agile/1.0/sprint/{sprintId}/issue` | Body: `{issues: [key1, key2, …]}` |
| `move_issues_to_backlog` | `move_issues_to_backlog` | POST | `/agile/1.0/backlog/issue` | Body: `{issues: [key1, key2, …]}` |

---

## 10. `FieldsMixin` (fields.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `search_fields` | `search_fields` | GET | `/field` | (cache: refresh on `refresh=True`) |
| `get_create_metadata` | (internal — `get_project_issue_types` callsite) | GET | `/issue/createmeta` | `projectKeys`, `issuetypeNames` |

---

## 11. `FieldOptionsMixin` (field_options.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_field_options` | `get_field_options` | GET | `/field/{fieldId}/context/{contextId}/option` (Cloud) | Cloud path: query `startAt`, `maxResults`. DC: uses `createmeta` `allowedValues` instead. |
| `_resolve_field_context_id` (helper) | – | GET | `/field/{fieldId}/context` | Returns list of `{id, name, isGlobal}`. If `context_id=None`, picks the global one. |

---

## 12. `ProjectsMixin` (projects.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_all_projects` | `get_all_projects` | GET | `/project/search` | `expand`, `query`, `maxResults`, `startAt` |
| `search_projects` | `search_projects` | GET | `/project/picker` | `query`, `maxResults`, `currentProjectId` |
| `get_project_issue_types` | `get_project_issue_types` | GET | `/issue/createmeta/{projectIdOrKey}/issuetypes` | – |
| `get_create_fields` | `get_create_fields` | GET | `/issue/createmeta/{projectIdOrKey}/issuetype/{issueTypeId}` | – |
| `get_project_fields` | `get_project_fields` | GET | `/issue/createmeta/{projectIdOrKey}/issuetypes` then aggregate | – |
| `get_project_versions` | `get_project_versions` | GET | `/project/{projectIdOrKey}/versions` | – |
| `get_project_components` | `get_project_components` | GET | `/project/{projectIdOrKey}/components` | – |
| `create_project_version` | `create_version`, `batch_create_versions` | POST | `/version` | Body: `{name, project, startDate?, releaseDate?, description?}` |
| `update_project_version` | `update_version` | PUT | `/version/{versionId}` | Body: partial: `{name?, description?, startDate?, releaseDate?, archived?, released?}` |

---

## 13. `QueuesMixin` (queues.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_service_desk_for_project` | `get_service_desk_for_project` | GET | `/servicedeskapi/servicedesk/by-project/{projectKey}` | **Server/DC only** — raises `NotImplementedError` on Cloud. |
| `get_service_desk_queues` | `get_service_desk_queues` | GET | `/servicedeskapi/servicedesk/{serviceDeskId}/queue` | `start`, `limit`, `includeCount` (query) |
| `get_queue_issues` | `get_queue_issues` | GET | `/servicedeskapi/servicedesk/{serviceDeskId}/queue/{queueId}/issue` | `start`, `limit` |
| `get_request_types` | `get_request_types` | GET | `/servicedeskapi/servicedesk/{serviceDeskId}/requesttype` | `start`, `limit` |
| `get_request_type_fields` | `get_request_type_fields` | GET | `/servicedeskapi/servicedesk/{serviceDeskId}/requesttype/{requestTypeId}/field` | – |

---

## 14. `CustomerRequestsMixin` (customer_requests.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `create_customer_request` | `create_customer_request` | POST | `/servicedeskapi/request` (Cloud) or `/rest/servicedeskapi/request` (DC) | Body: `{serviceDeskId, requestTypeId, requestFieldValues, raiseOnBehalfOf?, requestParticipants?, attachments?}` |

---

## 15. `AttachmentsMixin` (attachments.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_issue_attachments` | `get_issue_images`, `download_attachments` | GET | `/issue/{issueKey}?fields=attachment` | – |
| `get_issue_attachment_contents` | `download_attachments` | GET (per attachment) | `/attachment/content/{attachmentId}` (binary stream) | – |
| `fetch_attachment_content` | (helper) | GET | relative URL from attachment metadata | – |
| `upload_attachment` | (exposed via `update_issue` `attachments` param) | POST | `/issue/{issueKey}/attachments` | multipart: `file` (binary), `comment` (text). |

---

## 16. `WatchersMixin` (watchers.py)

| Method | Tool(s) | HTTP method | URL template | Query / Body |
| --- | --- | --- | --- | --- |
| `get_issue_watchers` | `get_issue_watchers` (also `get_issue` include) | GET | `/issue/{issueKey}/watchers` | – |
| `add_watcher` | `add_watcher` | POST | `/issue/{issueKey}/watchers` | Body: `accountId` (Cloud) or `username` (DC). For Basic auth on Cloud, JSON body shape is `"accountId": "..."`. For OAuth with no `accountId` field, the legacy string body `"username"` is used. |
| `remove_watcher` | `remove_watcher` | DELETE | `/issue/{issueKey}/watchers?accountId=...&username=...` | Query params (DELETE with query). |

---

## 17. `EpicsMixin` (epics.py)

The Epic operations all work through the regular issue update endpoint
(`/issue/{issueKey}`) with the Epic Link custom field. The mixin adds:

| Method | Tool(s) | HTTP method | URL template | Body |
| --- | --- | --- | --- | --- |
| `link_issue_to_epic` | (shared with `link_to_epic` tool) | PUT | `/issue/{issueKey}` | Body: `{fields: {<epic-field>: <epicKey>}}` |

The Epic field name is detected as:
- `customfield_12311140` on Cloud (configurable).
- `Epic Link` (the `com.atlassian.greenhopper.service.sprint.EpicLinkField` legacy name) on DC.

---

## 18. `MetricsMixin` (metrics.py) & `SLAMixin` (sla.py)

These compute metrics client-side from Jira changelog data (no direct REST
call of their own). They depend on `get_issue` (with `expand=changelog`) and
issue changelog endpoint `/issue/{key}/changelog`. The `metrics` module
adds working-hours filtering via `JIRA_SLA_WORKING_HOURS_*` env vars.

---

## 19. `DevelopmentMixin` (development.py)

| Method | Tool(s) | HTTP method | URL template | Body |
| --- | --- | --- | --- | --- |
| `get_issue_development_info` | `get_issue_development_info`, `get_issues_development_info` | POST | `/rest/devinfo/0.10/bulkByIssue` (Cloud) | Body: `{issueIds: [{id, type: "ISSUE_KEY"}], applicationTypes?, dataTypes?, _fields: ["pullrequests","branches","commits","repositories"]}`. |

The Atlassian Connect "Dev Info" endpoint was deprecated in 2024; Cloud
deployments may need to migrate to the Bitbucket integration's
`/bitbucket/1.0/repositories/{workspace}/{repo}/commits` shape.

---

## 20. `ProjectAnalysisMixin` (project_analysis.py)

| Method | Tool(s) | HTTP method | URL template | Body |
| --- | --- | --- | --- | --- |
| `get_project_epic_hierarchy` | `get_project_epic_hierarchy` | GET | `/search` (paginated) | JQL: `project = <KEY> AND issuetype = Epic`, then `/issue/{key}` for each epic to read parent / inward links. |
| `get_cross_project_dependencies` | `get_cross_project_dependencies` | GET | `/search` + `/issueLink` | Same as above, plus read issue links and group by target project. |

---

## 21. `FormsApiMixin` (forms_api.py)

Uses the **ProForma / Jira Forms REST API** (Cloud-only):

| Method | Tool(s) | HTTP method | URL template | Body |
| --- | --- | --- | --- | --- |
| `get_issue_forms` | `get_issue_proforma_forms` | GET | `/rest/api/3/issue/{issueIdOrKey}/form` | – |
| `get_form_details` | `get_proforma_form_details` | GET | `/rest/api/3/issue/{issueIdOrKey}/form/{formId}` | – |
| `update_form_answers` | `update_proforma_form_answers` | PUT | `/rest/api/3/issue/{issueIdOrKey}/form/{formId}` | Body: `{answers: [{questionId, type, value}, …]}`. The DATE/DATETIME `value` is converted from ISO 8601 string to Unix ms in milliseconds. |

---

## Notes for the Go port

- **`atlassian-python-api` is Python-only.** The Go port must either:
  - (a) call the underlying REST endpoints directly via `net/http`, or
  - (b) vendor a thin shim that wraps the same calls — there's no
    ready-made Go library at this level.
- **OAuth for Jira Cloud uses the gateway URL.**
  `https://api.atlassian.com/ex/jira/{cloudId}/rest/api/3/...`. The
  `JiraConfig.is_cloud` + `oauth_config.cloud_id` selects the base URL.
- **OAuth for DC uses the host URL directly.** The `oauth_config.base_url`
  is appended to the per-method path.
- **Basic auth on Cloud encodes `email:api_token` as Base64.**
  DC encodes `username:password`. PATs on either side use `Bearer <token>`.
- **Form-encoded request bodies** are rare; most endpoints expect
  `application/json`. Multipart for attachments only.
- **Async moves (`move_issue`)** use the bulk-move API and return a task ID
  to poll. The current code polls every 2 s for up to 30 s.
- **The `expand` query parameter** is heavily used; many tools augment the
  user's `expand` with internal expansions (`changelog`, `names`, …) before
  the request.
- **Server version differences** are mostly limited to authentication shape;
  the v2/v3 API gap is small. The Go client should default to v2 but accept
  v3 if the response data shape demands it.
- **No `requests` middleware chain is mandatory** in Go; the equivalent is
  a `http.RoundTripper` chain in the `*http.Client` transport.