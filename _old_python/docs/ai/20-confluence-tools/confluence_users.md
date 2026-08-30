# Confluence tools — `toolset:confluence_users`

## `search_user` (confluence_search_user)

**Tags:** `confluence`, `read`, `toolset:confluence_users`
**Annotations:** `title="Search User"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:1513-1597`

### Purpose

Search Confluence users via CQL (Cloud) or group-member API (DC).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `query` | `str` | yes | – | CQL string or plain full-name text |
| `limit` | `int` | no | `10` | Max results. `ge=1, le=50`. |
| `group_name` | `str` | no | `"confluence-users"` | Group to search within on DC; ignored on Cloud |

### Handler logic

If `query` doesn't contain CQL operators (`=`, `~`, `>`, `<`,
` AND `, ` OR `, `user.`), wrap it as `user.fullname ~ "<query>"`.

`user_results = confluence_fetcher.search_user(query, limit, group_name)`
→ JSON-serialize `to_simplified_dict()`.

### Underlying REST call

- Cloud: `GET /rest/api/search/user?cql=<cql>&limit=N`.
- DC: `GET /rest/api/group/<group>/member` paginated with `start`, `limit`.

### Output

`[user.to_simplified_dict(), ...]`

Each entry: `{username, userKey, displayName, emailAddress?, type, profilePicture?, ...}`.