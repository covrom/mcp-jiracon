# Jira tools — `toolset:jira_users`

Tools for looking up users, fetching profiles, and searching assignable
users.

## `get_user_profile` (jira_get_user_profile)

**Tags:** `jira`, `read`, `toolset:jira_users`
**Annotations:** `title="Get User Profile"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:250-308`

### Purpose

Retrieve profile information for a specific Jira user by any of the
supported identifier formats.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context (implicit) |
| `user_identifier` | `str` | yes | – | Email (Cloud), display name, account ID (`accountid:...`), or username / key for Server/DC |

### Handler logic

1. `jira = await get_jira_fetcher(ctx)`.
2. `user = jira.get_user_profile_by_identifier(user_identifier)`.
3. Catch exceptions:
   - `ValueError` with `"not found"` → warning log; return error JSON.
   - `MCPAtlassianAuthenticationError` → error JSON.
   - `OSError` / `HTTPError` → error JSON.
   - Other → log traceback at ERROR; return generic error JSON.
4. On success: `{"success": True, "user": user.to_simplified_dict()}`.

### Output

JSON string. `user` shape includes accountId / name / emailAddress /
displayName / active / avatarUrls / groups / applicationRoles.

### Underlying REST call

`GET /user` with `accountId`, `username`, or `key` query param depending
on identifier format. Expand `groups,applicationRoles`.

### Variants

- For Cloud: `accountId` is canonical. Tool auto-detects `accountid:` prefix.
- For DC: prefer `username`, fall back to `key` (legacy).

---

## `search_assignable_users` (jira_search_assignable_users)

**Tags:** `jira`, `read`, `toolset:jira_users`
**Annotations:** `title="Search Assignable Users"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:311-428`

### Purpose

Search Jira users who can be assigned issues in a given project or issue.
Resolves partial names / emails to concrete identifiers (accountId / name
/ key).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context (implicit) |
| `query` | `str` | yes | – | Display name, username, or email substring (partial, case-insensitive) |
| `project_key` | `str \| None` | no | `None` | Scope to project (e.g. `'DT'`). Required if `issue_key` is unset. |
| `issue_key` | `str \| None` | no | `None` | Scope to issue (e.g. `'DT-779'`). Required if `project_key` is unset. |
| `limit` | `int` | no | `20` | Max users to return. `ge=1, le=1000`. |

### Handler logic

1. Validate: exactly one of `project_key` or `issue_key` must be set.
   - If both or neither → return `{"success": False, "error": "..."}`
     JSON.
2. `users = jira.search_assignable_users(query=query, project_key=…, issue_key=…, limit=limit)`.
3. `{"success": True, "count": N, "users": [u.to_simplified_dict() for u in users]}`.

### Underlying REST call

`GET /user/assignable/search` with `query`, `project` OR `issueKey`,
`maxResults`, `startAt`.

### Variants

- Use this tool to resolve a partial user name before assigning an issue.
- The bot user pattern: an operator's bot account often lacks the global
  "Browse Users" permission, so `/user/search` returns empty. This
  `/user/assignable/search` endpoint works without that permission because
  it scopes by project/issue.

### Why no write variant

Jira REST does not expose a user-creation endpoint. User provisioning
happens via Atlassian admin or SCIM, not via REST.