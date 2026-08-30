# Jira tools — `toolset:jira_attachments`

Attachment download and inline image extraction.

## `download_attachments` (jira_download_attachments)

**Tags:** `jira`, `read`, `toolset:jira_attachments`
**Annotations:** `title="Download Attachments"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1177-1308`

### Purpose

Download every attachment on a Jira issue. Returns a list of mixed
content blocks: a text summary plus either an `EmbeddedResource` (for
recognized image MIME types) or a `TextContent` carrying a base64
payload (for everything else).

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |

### Handler logic

1. `result = jira.get_issue_attachment_contents(issue_key=issue_key)` →
   dict with `success`, `attachments: [{data, filename, content_type}, …]`,
   `failed`, `issue_key`, `total`.
2. On `result.success == False`, return single `TextContent` with the
   JSON error.
3. For each attachment:
   - If `len(data) > ATTACHMENT_MAX_BYTES` (50 MiB), record in `failed`
     with an explanatory message and continue.
   - Base64-encode the data.
   - Resolve MIME via `is_image_attachment(content_type, filename)`.
   - If image → `EmbeddedResource` with `uri="attachment:///{issue_key}/{filename}"`,
     `mimeType=resolved_mime`, `blob=base64`.
   - If non-image → `TextContent` carrying a JSON envelope with the base64
     payload (this works around an MCP-client limitation where
     `EmbeddedResource` blobs are only forwarded for image MIME types —
     see comment at line ~1262 of `servers/jira.py`).
4. Insert a `TextContent` summary at the head of the result list:
   `{success, issue_key, total, downloaded, failed, message?}`.

### Output

A list of content blocks:
- 1× `TextContent` (summary JSON)
- 0..N× `EmbeddedResource` (images)
- 0..N× `TextContent` (non-image base64 envelopes)

### Underlying REST calls

1. `GET /issue/{issueKey}?fields=attachment` → list of attachment metadata.
2. `GET /attachment/content/{attachmentId}` → binary stream.

### Variants

- The 50 MiB limit is enforced both pre-fetch (if `Content-Length` is
  known and exceeds) and post-fetch (after the bytes have arrived). Files
  exceeding the limit are reported in `failed` rather than raising.

---

## `get_issue_images` (jira_get_issue_images)

**Tags:** `jira`, `read`, `attachments`, `toolset:jira_attachments`
**Annotations:** `title="Get Issue Images"`, `readOnlyHint=True`
**Lines:** `servers/jira.py:1311-1440`

### Purpose

Return only the image attachments on a Jira issue as inline `ImageContent`
for LLM vision models.

### Input schema

| Name | Type | Required | Description |
| --- | --- | --- | --- |
| `ctx` | `Context` | yes | FastMCP context |
| `issue_key` | `str` | yes | Jira issue key |

### Handler logic

1. `attachments = jira.get_issue_attachments(issue_key)`.
2. For each, check `is_image_attachment(content_type, filename)`. If true,
   keep `(att, resolved_mime)`.
3. For each image:
   - If `att.size > ATTACHMENT_MAX_BYTES` → record in `failed`, skip.
   - If `not att.url` → record in `failed`, skip.
   - Fetch via `jira.fetch_attachment_content(url)`.
   - If fetched bytes > 50 MiB → record in `failed`, skip.
   - Else: base64-encode and emit `ImageContent(type="image", data=…, mimeType=resolved_mime)`.
4. Insert a `TextContent` summary at the head: `{success, issue_key, total_images, downloaded, failed}`.

### Output

A list of content blocks:
- 1× `TextContent` (summary)
- 0..N× `ImageContent` (one per image)

### Underlying REST calls

- `GET /issue/{issueKey}?fields=attachment` → metadata.
- `GET /attachment/content/{attachmentId}` → binary stream.

### Variants

- Image MIME types recognized: `image/png`, `image/jpeg`, `image/gif`,
  `image/webp`, `image/svg+xml`, `image/bmp`.
- `application/octet-stream` is rescued by filename-extension lookup via
  `mimetypes.guess_type`.