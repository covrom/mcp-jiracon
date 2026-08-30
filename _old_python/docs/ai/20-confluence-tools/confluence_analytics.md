# Confluence tools — `toolset:confluence_analytics`

Cloud-only view analytics.

## `get_page_views` (confluence_get_page_views)

**Tags:** `confluence`, `read`, `analytics`, `toolset:confluence_analytics`
**Annotations:** `title="Get Page Views"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:1757-1821`

### Purpose

Get view statistics for a Confluence page (total views, last viewed
date).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `page_id` | `str` | yes | – | Page ID |
| `include_title` | `bool` | no | `True` | Fetch and include page title |

### Handler logic

`result = confluence_fetcher.get_page_views(page_id, include_title)` →
JSON-serialize `to_simplified_dict()`.

### Underlying REST call

- Title lookup: `GET /rest/api/content/{page_id}?expand=title`.
- Analytics: `GET /rest/api/analytics/content/{page_id}/views`.

### Output

```json
{
  "title": "...",
  "views": {
    "totalViews": 1234,
    "lastViewed": "2026-01-09T..."
  }
}
```

### Variants

- **Cloud only.** Server/DC instances do not support the Analytics
  namespace. The tool raises `ValueError` (via the underlying client).