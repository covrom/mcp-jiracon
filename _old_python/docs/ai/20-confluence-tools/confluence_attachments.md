# Confluence tools — `toolset:confluence_attachments`

Attachment upload, list, download, and deletion.

## `upload_attachment` (confluence_upload_attachment)

**Tags:** `confluence`, `write`, `attachments`, `toolset:confluence_attachments`
**Annotations:** `title="Upload Attachment"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:1827-1966`

### Purpose

Upload a file (or new version of an existing file) to a page/blog.
Accepts the file as a server-readable path (`file_path`) or as a
base64 string (`content_base64` + `filename`).

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `content_id` | `str` | yes | – | Page or blog ID |
| `file_path` | `str \| None` | no | `None` | Absolute or CWD-relative path |
| `content_base64` | `str \| None` | no | `None` | Base64 file content (filesystem-free upload) |
| `filename` | `str \| None` | no | `None` | Required when `content_base64` is used |
| `comment` | `str \| None` | no | `None` | Version comment |
| `minor_edit` | `bool` | no | `False` | Skip watcher notifications |

### Validation

Exactly one of `file_path` / `content_base64` must be supplied.
`content_base64` must be valid base64. `filename` is required when
using `content_base64`.

### Handler logic

If `content_base64`: base64-decode → `confluence_fetcher.upload_attachment_from_content(content_id, filename, content, comment, minor_edit)`.
Else: `confluence_fetcher.upload_attachment(content_id, file_path, comment, minor_edit)`.

### Underlying REST call

`POST /rest/api/content/{content_id}/child/attachment` with multipart
body `file` (binary), `comment?` (text), `minorEdit=true|false`.
Header: `X-Atlassian-Token: no-check`.

If the upload returns 400 "file already exists", the mixin does a
lookup via `?filename=<urlencoded>` and POSTs to
`/child/attachment/{attachment_id}/data` for a new version.

### Output

`{"message": "Attachment uploaded successfully", "attachment": ...}`.

---

## `upload_attachments` (confluence_upload_attachments)

**Tags:** `confluence`, `write`, `attachments`, `toolset:confluence_attachments`
**Annotations:** `title="Upload Multiple Attachments"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:1969-2055`

### Purpose

Bulk upload multiple files in one API call.

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `content_id` | `str` | yes | – | Page or blog ID |
| `file_paths` | `str` | yes | – | Comma-separated file paths |
| `comment` | `str \| None` | no | `None` | Common comment |
| `minor_edit` | `bool` | no | `False` | Skip notifications |

### Handler logic

Split `file_paths` CSV → list. `confluence_fetcher.upload_attachments(content_id, paths, comment, minor_edit)`.

---

## `get_attachments` (confluence_get_attachments)

**Tags:** `confluence`, `read`, `attachments`, `toolset:confluence_attachments`
**Annotations:** `title="Get Content Attachments"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:2058-2160`

### Input schema

| Name | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| `ctx` | `Context` | yes | – | FastMCP context |
| `content_id` | `str` | yes | – | Page or blog ID |
| `start` | `int` | no | `0` | Pagination start |
| `limit` | `int` | no | `50` | Max results. `ge=1, le=100`. |
| `filename` | `str \| None` | no | `None` | Exact-match filename filter |
| `media_type` | `str \| None` | no | `None` | MIME type filter (note: most binaries return `application/octet-stream`) |

### Underlying REST call

- Cloud OAuth/PAT: `GET /api/v2/pages/{page_id}/attachments` with
  `start`, `limit`, `filename?`, `media-type?`, `sort?`.
- Otherwise: `GET /rest/api/content/{content_id}/child/attachment` with
  `start`, `limit`.

---

## `download_attachment` (confluence_download_attachment)

**Tags:** `confluence`, `read`, `attachments`, `toolset:confluence_attachments`
**Annotations:** `title="Download Attachment"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:2163-2312`

### Return type

`TextContent | EmbeddedResource` (single block)

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `attachment_id` | `str` | yes | Attachment ID (`att123456789`) |

### Handler logic

1. Fetch attachment metadata:
   - Cloud OAuth: `GET /api/v2/attachments/{attachment_id}` via
     `confluence_fetcher._v2_adapter.get_attachment_by_id`.
   - Otherwise: `GET /rest/api/content/{attachment_id}`.
2. Extract `_links.download`.
3. Resolve download URL via `_resolve_attachment_download_url`.
4. Determine MIME type: `extensions.mediaType` → `mimetypes.guess_type(filename)` → `application/octet-stream`.
5. Check `extensions.fileSize` against `ATTACHMENT_MAX_BYTES` (50 MiB).
6. Fetch the binary via `confluence_fetcher.fetch_attachment_content(download_url)`.
7. Verify post-fetch size.
8. Base64-encode and emit `EmbeddedResource(uri="attachment:///{id}/{filename}", mimeType=…, blob=…)`.
9. On any error → `TextContent` with JSON envelope `{"success": false, "error": ...}`.

### Variants

- Files larger than 50 MiB return an error `TextContent` rather than
  raising.

---

## `download_content_attachments` (confluence_download_content_attachments)

**Tags:** `confluence`, `read`, `attachments`, `toolset:confluence_attachments`
**Annotations:** `title="Download All Content Attachments"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:2315-2466`

### Return type

`list[TextContent | EmbeddedResource]`

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `content_id` | `str` | yes | Page or blog ID |

### Handler logic

1. Fetch all attachments via `get_content_attachments`.
2. For each: check `file_size`; download if under 50 MiB.
3. Emit `EmbeddedResource` per attachment with URI
   `attachment:///{content_id}/{filename}`.
4. Track `fetched` / `failed` arrays.
5. Insert summary `TextContent` at the head.

---

## `delete_attachment` (confluence_delete_attachment)

**Tags:** `confluence`, `write`, `attachments`, `toolset:confluence_attachments`
**Annotations:** `title="Delete Attachment"`, `destructiveHint=True`
**Decorators:** `@check_write_access`
**Lines:** `servers/confluence.py:2468-2519`

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `attachment_id` | `str` | yes | Attachment ID |

### Underlying REST call

- Cloud OAuth: `DELETE /api/v2/attachments/{attachment_id}`.
- Otherwise: `DELETE /rest/api/content/{attachment_id}`.

> **Warning:** this permanently deletes the attachment and all its
> versions.

---

## `get_page_images` (confluence_get_page_images)

**Tags:** `confluence`, `read`, `attachments`, `toolset:confluence_attachments`
**Annotations:** `title="Get Page Images"`, `readOnlyHint=True`
**Lines:** `servers/confluence.py:2522-2678`

### Return type

`list[TextContent | ImageContent]`

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `content_id` | `str` | yes | Page or blog ID |

### Handler logic

1. Fetch attachment metadata via `get_content_attachments`.
2. Filter to image MIME types via `is_image_attachment(media_type, filename)`.
3. For each image, check size + fetch + base64-encode.
4. Emit `ImageContent(type="image", data=base64, mimeType=…)`.
5. Insert summary `TextContent` at the head.

### Output

```json
{"success": true, "content_id": "...", "total_images": 3, "downloaded": 3, "failed": []}
```

Followed by `ImageContent` blocks for each successfully downloaded image.