# Confluence REST API mapping

This document maps each Confluence MCP tool to the underlying Atlassian
REST endpoints. Use it to spec the Go HTTP client.

> **Conventions used in this file**
> - Base URL for Server/DC: `{CONFLUENCE_URL}/rest/api/...`
> - Base URL for Cloud (Basic / PAT / mTLS / external): `{CONFLUENCE_URL}/wiki/rest/api/...` (note `/wiki` prefix auto-added by the lib)
> - Base URL for Cloud OAuth: `https://api.atlassian.com/ex/confluence/{cloudId}/wiki/rest/api/...`
> - Cloud v2 endpoint shape: `https://api.atlassian.com/ex/confluence/{cloudId}/wiki/api/v2/...`
>
> See [`utils/oauth.py`](../../src/mcp_atlassian/utils/oauth.py) for the
> URL helper logic; `confluence/_v1_rest_base_url` and `_v2_rest_base_url`
> resolve the actual base per auth mode.

---

## 1. `ConfluenceClient` (client.py)

`ConfluenceClient` is a thin wrapper around `atlassian-python-api`'s
`Confluence` class. The class itself doesn't issue HTTP; it delegates
through the mixins. It validates auth via:

| Method | Verb | URL template | Query |
| --- | --- | --- | --- |
| `_validate_authentication` | GET | `/space?start=0&limit=1` | `start`, `limit` |

Returns list of one space (or empty). Failure raises
`MCPAtlassianAuthenticationError`.

---

## 2. `PagesMixin` (pages.py)

| Method (mixin) | Tool(s) | Auth branch | Verb | URL template | Query / Body |
| --- | --- | --- | --- | --- | --- |
| `get_page_content` | `get_page` | OAuth Cloud | GET | `/api/v2/pages/{page_id}` | `body-format=storage` |
| `get_page_content` | `get_page` | Basic/PAT (Cloud or DC) | GET | `/rest/api/content/{page_id}` | `expand=body.storage,version,space,children.attachment,history` |
| `get_page_ancestors` | (internal, used by some tools) | identical | GET | `/rest/api/content/{page_id}/ancestors` | – |
| `_get_page_emoji` | (internal) | OAuth Cloud | GET | `/api/v2/pages/{page_id}/properties` | – |
| `_get_page_emoji` | (internal) | Basic/PAT (Cloud or DC) | GET | `/rest/api/content/{page_id}/property` | – |
| `_set_single_property` | (internal) | identical | DELETE/PUT/POST | `/rest/api/content/{page_id}/property[/...]` | Body varies. |
| `get_page_by_title` | `get_page` (title path) | identical | GET | `/rest/api/content` | `space={KEY}`, `title={TITLE}`, `expand=body.storage,version` |
| `get_space_pages` | (internal; used by `get_space_page_tree`) | identical | GET | `/rest/api/space/{KEY}/content/page` | `start`, `limit`, `expand=body.storage` |
| `create_page` | `create_page` | OAuth Cloud | POST | `/api/v2/pages` | Body: `{spaceId, status, title, body:{representation, value}, parentId?, subtype?}`. Resolves `spaceId` from key via `/api/v2/spaces?keys={KEY}` first. |
| `create_page` | `create_page` | Basic/PAT (Cloud or DC) | POST | `/rest/api/content` | Body: `{type:"page", space:{key:K}, title, ancestors?:[{id:parent}], body:{storage:{value, representation:"storage"}}}`. Note: DC's storage-format is just XHTML; Cloud's is the same but with macros. |
| `update_page` | `update_page` | OAuth Cloud | PUT | `/api/v2/pages/{page_id}` | Body: `{id, status, title, body:{representation, value}, version:{number:N+1, message?}}`. Reads current version via `/api/v2/pages/{page_id}?body-format=storage`. |
| `update_page` | `update_page` | Basic/PAT (Cloud or DC) | PUT | `/rest/api/content/{page_id}` | Body: `{id, type:"page", title, body:{storage:{...}}, version:{number:N+1, message?}, ancestors?:[{id}], minor_edit?:bool}`. |
| `update_page_section` | `update_page_section` | identical | GET + PUT | `/rest/api/content/{page_id}` (read), then PUT (update with reconstructed body) | No new endpoint; the mixin walks `<h1>`..`<h6>` to find the section, replaces it with the new content, and PUTs the reconstructed page. |
| `get_page_children` | `get_page_children` | OAuth Cloud | GET | `/api/v2/pages/{page_id}/direct-children` | `limit`, `cursor` (cursor-paginated). The mixin emulates `start`/`limit` semantics on top of v2 cursors. |
| `get_page_children` | `get_page_children` | Basic/PAT (Cloud or DC) | GET | `/rest/api/content/{page_id}/child/page` (+ `/child/folder` if `include_folders`) | `start`, `limit`, `expand`. |
| `get_space_page_tree` | `get_space_page_tree` | identical | GET (paginated) | `/rest/api/content` | `spaceKey={K}`, `start`, `limit`, `expand=ancestors`. Loops until exhausted. |
| `delete_page` | `delete_page` | OAuth Cloud | DELETE | `/api/v2/pages/{page_id}` | – |
| `delete_page` | `delete_page` | Basic/PAT (Cloud or DC) | DELETE | `/rest/api/content/{page_id}` | `status=trashed` query (the lib adds it). |
| `get_page_history` | `get_page_history` | OAuth Cloud | GET | `/api/v2/pages/{page_id}/versions` then `/api/v2/versions/{versionId}` | Resolves version number → versionId first. |
| `get_page_history` | `get_page_history` | Basic/PAT (Cloud or DC) | GET | `/rest/api/content/{page_id}` | `status=historical`, `version={N}`, `expand=body.storage,version,space,children.attachment,history`. |
| `move_page` | `move_page` | OAuth Cloud (via v2 adapter) | PUT | `/rest/api/content/{page_id}/move/{position}/{targetId}` (v1 fallback because v2 has no OAuth-friendly move) | – |
| `move_page` | `move_page` | Basic/PAT (Cloud or DC) | PUT | `/rest/api/content/{page_id}/move/{position}/{targetId}` | – |
| `get_page_version_diff` | `get_page_version_diff` | identical | GET (×2) | `/rest/api/content/{page_id}` with `version=N` and `version=M` | Computes unified diff client-side. |
| `copy_page` | `copy_page` | Cloud | POST | `/wiki/rest/api/content/{source_page_id}/copy` | Body: `{copyAttachments, copyPermissions, copyProperties, copyLabels, pageTitle, destination:{type:"parent_page"|"space", value:idOrKey}}` |
| `copy_page` | `copy_page` | DC | GET + POST | `/rest/api/content/{source_page_id}` then `/rest/api/content` | Manual copy: read source, post new page. No native copy endpoint on DC. |

---

## 3. `SearchMixin` (search.py)

| Method | Tool(s) | Verb | URL template | Query |
| --- | --- | --- | --- | --- |
| `search` | `confluence_search` | GET | `/rest/api/content/search` | `cql`, `limit`, `expand=content.history,content.version`. CQL accepts `siteSearch ~ "..."` for text-only queries. |
| `search_user` | `confluence_search_user` (Cloud) | GET | `/rest/api/search/user` | `cql`, `limit` |
| `search_user` | `confluence_search_user` (DC) | GET (paginated) | `/rest/api/group/{group}/member` | `start`, `limit`. Default group: `confluence-users`. |

---

## 4. `SpacesMixin` (spaces.py)

| Method | Tool(s) | Verb | URL template | Query |
| --- | --- | --- | --- | --- |
| `get_spaces` | (none exposed as a tool directly; used internally) | GET | `/rest/api/space` | `start`, `limit` |
| `get_user_contributed_spaces` | (used internally) | GET | `/rest/api/content/search` | `cql=contributor=currentUser() order by lastmodified DESC`, `limit` |

---

## 5. `CommentsMixin` (comments.py)

| Method | Tool(s) | Auth branch | Verb | URL template | Body |
| --- | --- | --- | --- | --- | --- |
| `get_page_comments` | `confluence_get_comments` | identical | GET | `/rest/api/content/{page_id}` then `/rest/api/content/{page_id}/child/comment` | `expand=body.view.value,version`, `depth=all` |
| `add_comment` | `confluence_add_comment` | OAuth Cloud / PAT Cloud | POST | `/api/v2/footer-comments` | `{body:{representation:"storage", value:html}, pageId:...}` |
| `add_comment` | `confluence_add_comment` | Basic / DC | GET + POST | `/rest/api/content/{page_id}` then `/rest/api/content` | `{type:"comment", container:{id, type:"page"}, body:{storage:{value, representation}}}` |
| `reply_to_comment` | `confluence_reply_to_comment` | OAuth Cloud / PAT Cloud | POST | `/api/v2/footer-comments` | `{body:{representation, value}, parentCommentId:...}` |
| `reply_to_comment` | `confluence_reply_to_comment` | Basic / DC | GET + POST | `/rest/api/content/{commentId}` (resolve parent), then `/rest/api/content/` | `{type:"comment", container:{id, type:"page"}, ancestors:[{id:commentId}], body:{storage:{...}}}` |
| `get_inline_comments` | `confluence_get_inline_comments` | OAuth Cloud | GET | `/api/v2/pages/{page_id}/inline-comments` | `body-format=storage`, `status?` |
| `get_inline_comments` | `confluence_get_inline_comments` | Basic / DC | GET (×2) | `/rest/api/content/{page_id}` then `/rest/api/content/{page_id}/child/comment` | `expand=body.view.value,version,extensions.inlineProperties`, `depth=all` |
| `add_inline_comment` | `confluence_add_inline_comment` | OAuth Cloud | POST | `/api/v2/inline-comments` | `{pageId, body:{representation, value}, inlineCommentProperties:{textSelection, textSelectionMatchCount, textSelectionMatchIndex}}` |
| `add_inline_comment` | `confluence_add_inline_comment` | Basic / DC | GET + POST | `/rest/api/content/{page_id}` then `/rest/api/content/` | `{type:"comment", container:{id, type:"page"}, body:{storage:{...}}, extensions:{location:"inline", inlineProperties:{originalSelection, numMatches, matchIndex, lastFetchTime:ms, serializedHighlights:'[["…"]]'}}}` |

---

## 6. `LabelsMixin` (labels.py)

| Method | Tool(s) | Verb | URL template | Body |
| --- | --- | --- | --- | --- |
| `get_page_labels` | `confluence_get_labels` | GET | `/rest/api/content/{page_id}/label` | – |
| `add_page_label` | `confluence_add_label` | POST | `/rest/api/content/{page_id}/label` | Body: `[{"prefix":"global","name":"<label>"}]` |
| (delete label) | (no exposed tool) | DELETE | `/rest/api/content/{page_id}/label/{name}` | – |

---

## 7. `UsersMixin` (users.py)

| Method | Tool(s) | Verb | URL template | Query |
| --- | --- | --- | --- | --- |
| `get_user_details_by_accountid` | (used by preprocessing for user mentions) | GET | `/rest/api/user` | `accountId={id}`, `expand?` (Cloud) |
| `get_user_details_by_username` | (used by preprocessing) | GET | `/rest/api/user` | `username={name}` (DC) |
| `get_user_details_by_userkey` | (used by preprocessing) | GET | `/rest/api/user` | `userKey={key}` (DC legacy) |
| `get_current_user_info` | (used in fetcher resolution) | GET | `/rest/api/user/current` | – |

---

## 8. `AttachmentsMixin` (attachments.py)

| Method | Tool(s) | Auth branch | Verb | URL template | Body / Notes |
| --- | --- | --- | --- | --- | --- |
| `_resolve_attachment_download_url` | (helper) | identical | GET (rewrite) | `{base}/rest/api/content/{content_id}/child/attachment/{attachment_id}/download?version?` | Re-writes legacy `/download/attachments/...` URLs to the modern v1 path. |
| `get_content_attachments` | `confluence_get_attachments` | OAuth Cloud / PAT Cloud | GET | `/api/v2/pages/{page_id}/attachments` | `start`, `limit`, `filename?`, `media-type?`, `sort?` |
| `get_content_attachments` | `confluence_get_attachments` | Basic / DC | GET | `/rest/api/content/{content_id}/child/attachment` | `start`, `limit` |
| `_send_attachment_request` | `confluence_upload_attachment`, `confluence_upload_attachments` | identical | POST | `{base}/rest/api/content/{content_id}/child/attachment` | multipart: `file` (binary), `comment?`, form field `minorEdit=true|false`. Headers: `X-Atlassian-Token: no-check`. On HTTP 400 "file already exists", the mixin does a lookup via `?filename={urlencoded}` and POSTs to `…/attachment/{attachment_id}/data` for a new version. |
| `download_attachment` | `confluence_download_attachment` | OAuth Cloud | GET | `/api/v2/attachments/{attachment_id}` then `_links.download` | binary fetch via same `download` URL. |
| `download_attachment` | `confluence_download_attachment` | Basic / DC | GET | `/rest/api/content/{attachment_id}` | binary fetch via `_links.download`. |
| `fetch_attachment_content` | (helper) | identical | GET | absolute URL (legacy or rewritten) | – |
| `delete_attachment` | `confluence_delete_attachment` | OAuth Cloud | DELETE | `/api/v2/attachments/{attachment_id}` | – |
| `delete_attachment` | `confluence_delete_attachment` | Basic / DC | DELETE | `/rest/api/content/{attachment_id}` | – |

---

## 9. `TemplatesMixin` (templates.py)

All endpoints are Cloud-only.

| Method | Tool(s) | Verb | URL template | Query |
| --- | --- | --- | --- | --- |
| `list_page_templates` | `confluence_list_page_templates` | GET | `/rest/api/template/page` | `limit`, `spaceKey?` |
| `get_page_template` | `confluence_get_page_template` | GET | `/rest/api/template/{templateId}` | – |
| `create_page_from_template` | `confluence_create_page_from_template` | GET + POST | `/rest/api/template/{templateId}` (read template) then `/rest/api/content` (create page with template body) | Body for create: see `create_page` v1 payload |

---

## 10. `PermissionsMixin` (permissions.py)

All endpoints are Cloud-only.

| Method | Tool(s) | Verb | URL template | Body |
| --- | --- | --- | --- | --- |
| `check_content_permissions` | `confluence_check_content_permissions` | POST | `/rest/api/content/{content_id}/permission/check` | `{operation:str, subject:{type:"user|group", identifier:accountIdOrGroupId}}` |
| `get_space_permissions` | `confluence_get_space_permissions` | GET | `/api/v2/spaces/{space_id}/permissions` | `limit`, `cursor?` (note: bare `/api/v2/...` not `/wiki/api/v2/...` when on the OAuth gateway) |

---

## 11. `RestrictionsMixin` (restrictions.py)

| Method | Tool(s) | Verb | URL template | Body |
| --- | --- | --- | --- | --- |
| `get_page_restrictions` | `confluence_get_page_restrictions` | GET | `/rest/api/content/{page_id}/restriction/byOperation` | – |
| `set_page_restrictions` | `confluence_set_page_restrictions` | PUT | `/rest/api/content/{page_id}/restriction` | Body: JSON list `[{"operation":"read","restrictions":{"user":[...],"group":[...]}},{"operation":"update","restrictions":{...}}]`. User entries are `{type:"known", accountId:...}` on Cloud or `{type:"known", username:...}` on DC. Group entries are `{type:"group", name:...}`. |

---

## 12. `AnalyticsMixin` (analytics.py)

All endpoints are Cloud-only.

| Method | Tool(s) | Verb | URL template | Notes |
| --- | --- | --- | --- | --- |
| `get_page_views` (title lookup) | `confluence_get_page_views` | GET | `/rest/api/content/{page_id}` | `expand=title` |
| `_get_page_views_direct` | `confluence_get_page_views` | GET | `{confluence.url}/rest/api/analytics/content/{page_id}/views` | – |
| `batch_get_page_views` | (internal) | – | loops `get_page_views` | – |

---

## 13. `ConfluenceV2Adapter` (v2_adapter.py)

The v2 adapter is used by OAuth Cloud for endpoints where v1 doesn't work
cleanly with bearer tokens. All v2 calls use the Cloud OAuth gateway URL.

| Method | Verb | URL template | Query / Body |
| --- | --- | --- | --- |
| `_get_space_id` | GET | `/api/v2/spaces` | `keys={spaceKey}` |
| `create_page` | POST | `/api/v2/pages` | `{spaceId, status, title, body:{representation, value}, parentId?, subtype?}` |
| `_get_page_version` | GET | `/api/v2/pages/{page_id}` | `body-format=storage` |
| `update_page` | PUT | `/api/v2/pages/{page_id}` | `{id, status, title, body:{representation, value}, version:{number:N+1, message?}}` |
| `_get_space_from_id` | GET | `/api/v2/spaces/{space_id}` | – |
| `get_page` | GET | `/api/v2/pages/{page_id}` | `body-format=storage` |
| `get_page_direct_children` | GET | `/api/v2/pages/{page_id}/direct-children` | `limit?`, `cursor?` |
| `delete_page` | DELETE | `/api/v2/pages/{page_id}` | – |
| `create_footer_comment` | POST | `/api/v2/footer-comments` | `{body:{representation, value}, pageId?|parentCommentId?}` (exactly one) |
| `get_footer_comment` | GET | `/api/v2/footer-comments/{comment_id}` | `body-format=storage` |
| `get_inline_comments` | GET | `/api/v2/pages/{page_id}/inline-comments` | `body-format=storage`, `status?` |
| `create_inline_comment` | POST | `/api/v2/inline-comments` | `{pageId, body:{representation, value}, inlineCommentProperties:{...}}` |
| `move_page` | PUT | `/rest/api/content/{page_id}/move/{position}/{targetId}` | (v1 move URL because OAuth-friendly v2 move isn't available) |
| `get_page_emoji` | GET | `/api/v2/pages/{page_id}/properties` | – |
| `_get_property` | GET | `/api/v2/pages/{page_id}/properties/{property_key}` | – |
| `_set_page_property` (delete) | DELETE | `/api/v2/pages/{page_id}/properties/{property_key}` | – |
| `_set_page_property` (update) | PUT | `/api/v2/pages/{page_id}/properties/{property_key}` | `{key, value, version:{number:N+1}}` |
| `_set_page_property` (create) | POST | `/api/v2/pages/{page_id}/properties` | `{key, value}` |
| `get_page_versions_list` | GET | `/api/v2/pages/{page_id}/versions` | – |
| `get_page_by_version` | GET | `/api/v2/versions/{versionId}` | `body-format=storage` (resolves number → versionId first) |
| `get_page_views` | GET | `/rest/api/analytics/content/{page_id}/views` | – |
| `get_page_attachments` | GET | `/api/v2/pages/{page_id}/attachments` | `start`, `limit`, `filename?`, `media-type?`, `sort?` |
| `get_attachment_by_id` | GET | `/api/v2/attachments/{attachment_id}` | – |
| `delete_attachment` | DELETE | `/api/v2/attachments/{attachment_id}` | – |

---

## Notes for the Go port

- The **`atlassian-python-api` library is the source of truth** for the URL
  patterns above. The Go port must speak the same REST dialect.
- **OAuth Cloud requires the gateway URL prefix.** Detect Cloud OAuth by
  `oauth_config.cloud_id != ""`. Detect Cloud Basic/PAT by URL hostname
  (`*.atlassian.net`) and add `/wiki` to REST paths. Detect DC by URL
  hostname not matching Cloud patterns; use the URL directly.
- **v2 vs v1 selection logic** (Cloud OAuth):
  - **v2** for: pages CRUD + properties + versions, spaces, footer/inline
    comments, attachments, direct-children.
  - **v1** still used for: move (no OAuth-friendly v2), copy, page
    restrictions, and analytics (no v2 equivalent).
- **Cursor vs offset pagination**: v2 uses `cursor`; v1 uses `start`/`limit`.
  For Cloud OAuth page children, the current code emulates v1's
  `start`/`limit` semantics on top of v2 cursors. The Go port should keep
  this emulation so existing callers don't break.
- **Inline comments** carry an empirical DC payload extension with
  `numMatches`/`matchIndex`/`lastFetchTime`/`serializedHighlights` fields.
  The Cloud v2 inline-comments payload uses `textSelection*` instead.
- **Storage format** is XHTML with Confluence-specific macros (`<ac:*>`,
  `<ri:*>`). The `preprocessing/confluence.py` module converts Markdown
  → storage format using `md2conf`; the Go port needs an equivalent.
- **Attachment size cap is 50 MiB** (constant `ATTACHMENT_MAX_BYTES`); both
  the Confluence attachment tools and Jira tools respect it.
- **Streaming downloads** use `requests` Session streaming; in Go, use
  `http.Response.Body` with `io.LimitReader` and check `Content-Length`
  before fully buffering.
- **Server/DC limits**: most v2-only endpoints (templates, permissions,
  analytics, fields context options, bulk move) are Cloud-only and should
  raise `NotImplementedError` on DC. The Go port should mirror this
  surface explicitly.