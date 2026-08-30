# Confluence tools — `toolset:confluence_comments`

Page-level (footer) comments and inline (text-anchored) comments.
Default-enabled toolset.

## `get_comments` (confluence_get_comments)

**Tags:** `confluence`, `read`, `toolset:confluence_comments`
**Annotations:** `title="Get Comments"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:561-590`

### Purpose

Get footer comments for a Confluence page.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `page_id` | `str` | yes | Page ID |

### Underlying REST call

`GET /rest/api/content/{page_id}/child/comment?expand=body.view.value,version&depth=all`

### Output

`[comment.to_simplified_dict(), ...]`

---

## `add_comment` (confluence_add_comment)

**Tags:** `confluence`, `write`, `toolset:confluence_comments`
**Annotations:** `title="Add Comment"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:1270-1318`

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `page_id` | `str` | yes | Page ID |
| `body` | `str` | yes | Markdown comment body |

### Underlying REST call

- Cloud OAuth/PAT: `POST /api/v2/footer-comments` with
  `{body:{representation:"storage", value:html}, pageId}`.
- Otherwise: GET `/rest/api/content/{page_id}` then POST
  `/rest/api/content` with
  `{type:"comment", container:{id, type:"page"}, body:{storage:{value, representation}}}`.

---

## `reply_to_comment` (confluence_reply_to_comment)

**Tags:** `confluence`, `write`, `toolset:confluence_comments`
**Annotations:** `title="Reply to Comment"`, `destructiveHint=False`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:1321-1371`

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `comment_id` | `str` | yes | Parent comment ID |
| `body` | `str` | yes | Markdown reply body |

### Underlying REST call

- Cloud OAuth/PAT: `POST /api/v2/footer-comments` with
  `{body:{representation, value}, parentCommentId}`.
- DC: GET `/rest/api/content/{commentId}` (resolve parent) then POST
  `/rest/api/content/` with
  `{type:"comment", container:{id, type:"page"}, ancestors:[{id:commentId}], body:{storage:{...}}}`.

---

## `get_inline_comments` (confluence_get_inline_comments)

**Tags:** `confluence`, `read`, `toolset:confluence_comments`
**Annotations:** `title="Get Inline Comments"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:1374-1415`

### Purpose

Get inline (text-anchored) comments for a page.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `page_id` | `str` | yes | Page ID |

### Underlying REST call

- Cloud OAuth: `GET /api/v2/pages/{page_id}/inline-comments?body-format=storage&status=`.
- DC: GET page + `/rest/api/content/{page_id}/child/comment?expand=body.view.value,version,extensions.inlineProperties&depth=all`.

### Output

```json
{
  "success": true,
  "page_id": "...",
  "count": 2,
  "comments": [comment.to_simplified_dict(), ...]
}
```

---

## `add_inline_comment` (confluence_add_inline_comment)

**Tags:** `confluence`, `write`, `toolset:confluence_comments`
**Annotations:** `title="Add Inline Comment"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:1418-1510`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `page_id` | `str` | yes | – | Page ID |
| `body` | `str` | yes | – | Markdown comment body |
| `text_selection` | `str` | yes | – | Exact text snippet on the page |
| `text_selection_match_count` | `int` | no | `1` | Total occurrences. `ge=1`. |
| `text_selection_match_index` | `int` | no | `0` | Zero-based index. `ge=0`. |

### Underlying REST call

- Cloud OAuth: `POST /api/v2/inline-comments` with
  `{pageId, body:{representation, value}, inlineCommentProperties:{textSelection, textSelectionMatchCount, textSelectionMatchIndex}}`.
- DC: GET page + POST `/rest/api/content/` with
  `{type:"comment", container:{id, type:"page"}, body:{storage:{...}}, extensions:{location:"inline", inlineProperties:{originalSelection, numMatches, matchIndex, lastFetchTime:ms, serializedHighlights:'[["..."]]'}}}`.
  The DC payload's exact field names (`numMatches`/`matchIndex`/
  `lastFetchTime`/`serializedHighlights`) were empirically derived from
  Confluence DC 8.x.