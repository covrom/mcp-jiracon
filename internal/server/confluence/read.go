package confluence

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	core "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server/core"

	"github.com/mark3labs/mcp-go/mcp"

	confmodels "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/confluence"
)

var numericRe = regexp.MustCompile(`^\d+$`)
var pageIDPathRe = regexp.MustCompile(`(?:^|/)pages/([0-9]+)(?:/|$)`)
var tinyLinkPathRe = regexp.MustCompile(`(?:^|/)x/([A-Za-z0-9_-]+)(?:/|$)`)
var tinyIDPatternRe = regexp.MustCompile(`^[A-Za-z0-9_-]{1,11}$`)

const maxConfluencePageID = int64(1<<63 - 1)

// resolvePageID resolves a page ID from various input formats.
// Mirrors Python's _resolve_page_id().
func resolvePageID(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	if numericRe.MatchString(input) {
		return input, nil
	}

	parsed, err := url.Parse(input)
	if err != nil {
		return input, nil
	}

	// Full URL: extract numeric ID from /pages/NNN path segment.
	if m := pageIDPathRe.FindStringSubmatch(parsed.Path); m != nil {
		return m[1], nil
	}

	// Query parameter: ?pageId=NNN
	queryIDs := parsed.Query()["pageId"]
	if len(queryIDs) == 1 && numericRe.MatchString(queryIDs[0]) {
		return queryIDs[0], nil
	}

	// Tiny link: /x/AbCdE -> decode base64-encoded Confluence ID.
	if m := tinyLinkPathRe.FindStringSubmatch(parsed.Path); m != nil {
		encoded := m[1]
		resolved := decodeTinyLink(encoded)
		if resolved > 0 {
			return fmt.Sprintf("%d", resolved), nil
		}
	}

	return input, nil
}

// decodeTinyLink decodes a Confluence tiny-link identifier to a numeric page ID.
// Mirrors Python's _decode_confluence_tiny_id(). Returns 0 when invalid.
func decodeTinyLink(encoded string) int64 {
	if !tinyIDPatternRe.MatchString(encoded) {
		return 0
	}
	standard := strings.ReplaceAll(encoded, "-", "/")
	standard = strings.ReplaceAll(standard, "_", "+")
	// Pad with "A" characters to length 11, then add "=" for standard base64.
	padded := standard
	for len(padded) < 11 {
		padded += "A"
	}
	padded += "="
	decoded, err := base64.StdEncoding.DecodeString(padded)
	if err != nil || len(decoded) != 8 {
		return 0
	}
	id := int64(binary.LittleEndian.Uint64(decoded))
	if id <= 0 || id > maxConfluencePageID {
		return 0
	}
	// Validate round-trip.
	if encodeTinyLink(uint64(id)) != encoded {
		return 0
	}
	return id
}

// encodeTinyLink encodes a uint64 ID as a Confluence tiny-link token.
func encodeTinyLink(id uint64) string {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, id)
	// Strip trailing zero bytes.
	for len(buf) > 0 && buf[len(buf)-1] == 0 {
		buf = buf[:len(buf)-1]
	}
	encoded := base64.StdEncoding.EncodeToString(buf)
	// Remove padding.
	encoded = strings.TrimRight(encoded, "=")
	// Reverse Confluence's encoding: replace "/" with "-", "+" with "_"
	encoded = strings.ReplaceAll(encoded, "/", "-")
	encoded = strings.ReplaceAll(encoded, "+", "_")
	return encoded
}

type confluenceGetPageTool struct{}

func (t confluenceGetPageTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	a := base.Args
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	pageID := a.String("page_id")
	title := a.String("title")
	spaceKey := a.String("space_key")
	includeMetadata := a.Bool("include_metadata", true)
	convertToMarkdown := a.Bool("convert_to_markdown", true)

	var page *confmodels.ConfluencePage
	var fetchErr error
	if pageID != "" {
		resolvedID, resolveErr := resolvePageID(pageID)
		if resolveErr != nil {
			return mcp.NewToolResultText(core.Stringify(map[string]any{"error": fmt.Sprintf("Failed to retrieve page by ID '%s': %v", pageID, resolveErr)})), nil
		}
		page, fetchErr = base.Confluence.GetPageContent(ctx, resolvedID, convertToMarkdown)
	} else if title != "" && spaceKey != "" {
		page, fetchErr = base.Confluence.GetPageByTitleContent(ctx, spaceKey, title, convertToMarkdown)
	} else {
		return mcp.NewToolResultText(core.Stringify(map[string]any{"error": "Either 'page_id' OR both 'title' and 'space_key' must be provided."})), nil
	}
	if fetchErr != nil {
		return mcp.NewToolResultText(core.Stringify(map[string]any{"error": fmt.Sprintf("Failed to retrieve page: %v", fetchErr)})), nil
	}
	if page == nil {
		return mcp.NewToolResultText(core.Stringify(map[string]any{"error": "Page not found with the provided identifiers."})), nil
	}

	if includeMetadata {
		return mcp.NewToolResultText(core.Stringify(map[string]any{"metadata": page.ToSimplifiedDict()})), nil
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"content": map[string]any{"value": page.Content}})), nil
}

type confluenceSearchTool struct{}

func (t confluenceSearchTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	limit := base.Args.Int("limit", 10)
	query := base.Args.String("query")
	searchType := base.Args.String("search_type")
	spacesFilter := base.Args.String("spaces_filter")

	// Build CQL — mirror Python logic.
	// If the query contains no CQL operators, wrap as text search.
	cql := query
	isSimple := !strings.Contains(query, "=") && !strings.Contains(query, "~") &&
		!strings.Contains(query, ">") && !strings.Contains(query, "<") &&
		!strings.Contains(query, " AND ") && !strings.Contains(query, " OR ") &&
		!strings.Contains(query, "currentUser()") &&
		!strings.Contains(query, "user.") && !strings.Contains(query, "memberOf")
	if searchType == "text" || (searchType == "" && isSimple) {
		// DC-only: use text search (siteSearch is Cloud-only, fails on DC).
		cql = fmt.Sprintf(`text ~ "%s"`, strings.ReplaceAll(query, `"`, `\"`))
	}
	if spacesFilter != "" {
		spaces := strings.Split(spacesFilter, ",")
		for i, s := range spaces {
			spaces[i] = fmt.Sprintf(`"%s"`, strings.TrimSpace(s))
		}
		cql = fmt.Sprintf(`(%s) AND space IN (%s)`, cql, strings.Join(spaces, ","))
	}
	r, err := base.Confluence.Search(ctx, cql, limit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetPageChildrenTool struct{}

func (t confluenceGetPageChildrenTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	start := base.Args.Int("start", 0)
	limit := base.Args.Int("limit", 25)
	includeContent := base.Args.Bool("include_content", false)
	includeFolders := base.Args.Bool("include_folders", true)
	convertToMarkdown := base.Args.Bool("convert_to_markdown", true)
	r, err := base.Confluence.GetPageChildren(ctx, base.Args.String("parent_id"), start, limit, includeContent, convertToMarkdown, includeFolders)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetSpacePageTreeTool struct{}

func (t confluenceGetSpacePageTreeTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	limit := base.Args.Int("limit", 100)
	r, err := base.Confluence.GetSpacePageTree(ctx, base.Args.String("space_key"), limit)
	if err != nil {
		return nil, err
	}
	if hasMore, _ := r["has_more"].(bool); hasMore {
		r["hint"] = "Tree truncated: increase 'limit' to fetch more pages"
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetCommentsTool struct{}

func (t confluenceGetCommentsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	format := base.Args.String("format")
	comments, err := base.Confluence.GetPageComments(ctx, base.Args.String("page_id"))
	if err != nil {
		return nil, err
	}
	if format == "raw" {
		rawComments := make([]map[string]any, len(comments))
		for i, c := range comments {
			rawComments[i] = map[string]any{
				"id":        c.ID,
				"parent_id": c.ParentID,
				"body":      c.Body,
				"author":    c.Author,
				"created":   c.Created,
			}
		}
		return mcp.NewToolResultText(core.Stringify(rawComments)), nil
	}
	out := make([]map[string]any, len(comments))
	for i, c := range comments {
		out[i] = c.ToSimplifiedDict()
	}
	return mcp.NewToolResultText(core.Stringify(out)), nil
}

type confluenceGetLabelsTool struct{}

func (t confluenceGetLabelsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.GetPageLabels(ctx, base.Args.String("page_id"), base.Args.String("prefix"), base.Args.Int("start", 0), base.Args.Int("limit", 0))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceSearchUserTool struct{}

func (t confluenceSearchUserTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	limit := base.Args.Int("limit", 10)
	groupName := base.Args.String("group_name")
	if groupName == "" {
		groupName = "confluence-users"
	}
	// Strip CQL wrapper if present, pass raw name for DC group search.
	query := base.Args.String("query")
	if strings.HasPrefix(query, "user.fullname ~ ") {
		query = strings.Trim(query[len("user.fullname ~ "):], `"`)
	}
	r, err := base.Confluence.SearchUser(ctx, query, limit, groupName)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetAttachmentsTool struct{}

func (t confluenceGetAttachmentsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	start := base.Args.Int("start", 0)
	limit := base.Args.Int("limit", 50)
	filterFilename := base.Args.String("filename")
	filterMediaType := base.Args.String("media_type")
	r, err := base.Confluence.GetContentAttachments(ctx, base.Args.String("content_id"), start, limit, filterFilename, filterMediaType)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetInlineCommentsTool struct{}

func (t confluenceGetInlineCommentsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	resultFormat := base.Args.String("result_format")
	comments, err := base.Confluence.GetInlineComments(ctx, base.Args.String("page_id"))
	if err != nil {
		return nil, err
	}
	if resultFormat == "raw" {
		return mcp.NewToolResultText(core.Stringify(comments)), nil
	}
	// Simplified format: extract key fields from each comment.
	simplified := make([]map[string]any, len(comments))
	for i, c := range comments {
		body := ""
		if b, _ := c["body"].(map[string]any); b != nil {
			if s, _ := b["storage"].(map[string]any); s != nil {
				body, _ = s["value"].(string)
			}
			if body == "" {
				if v, _ := b["view"].(map[string]any); v != nil {
					body, _ = v["value"].(string)
				}
			}
		}
		simplified[i] = map[string]any{
			"id":        c["id"],
			"parent_id": parentIDOfComment(c),
			"body":      body,
			"created":   c["created"],
			"author":    c["author"],
		}
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"success":  true,
		"page_id":  base.Args.String("page_id"),
		"count":    len(comments),
		"comments": simplified,
	})), nil
}

// parentIDOfComment extracts the direct parent comment id from a raw API
// comment map (last entry of the "ancestors" chain; empty for top-level).
func parentIDOfComment(c map[string]any) string {
	anc, _ := c["ancestors"].([]any)
	if len(anc) == 0 {
		return ""
	}
	last, _ := anc[len(anc)-1].(map[string]any)
	if last == nil {
		return ""
	}
	return fmt.Sprintf("%v", last["id"])
}

type confluenceGetPageHistoryTool struct{}

func (t confluenceGetPageHistoryTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	version := base.Args.Int("version", 0)
	convertToMarkdown := base.Args.Bool("convert_to_markdown", true)
	r, err := base.Confluence.GetPageHistory(ctx, base.Args.String("page_id"), version, convertToMarkdown)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetPageDiffTool struct{}

func (t confluenceGetPageDiffTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	fromV := base.Args.Int("from_version", 0)
	toV := base.Args.Int("to_version", 0)
	convertToMarkdown := base.Args.Bool("convert_to_markdown", true)
	r, err := base.Confluence.GetPageVersionDiff(ctx, base.Args.String("page_id"), fromV, toV, convertToMarkdown)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceDownloadAttachmentTool struct{}

func (t confluenceDownloadAttachmentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.DownloadAttachmentBase64(ctx, base.Args.String("attachment_id"))
	if err != nil {
		return nil, err
	}
	if ok, _ := r["success"].(bool); !ok {
		return mcp.NewToolResultText(core.Stringify(r)), nil
	}
	b64, _ := r["data_base64"].(string)
	mime, _ := r["mime_type"].(string)
	filename, _ := r["filename"].(string)
	aid, _ := r["attachment_id"].(string)
	if b64 == "" {
		return mcp.NewToolResultText(core.Stringify(r)), nil
	}
	decoded, _ := base64.StdEncoding.DecodeString(b64)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(core.Stringify(map[string]any{
				"success": true, "attachment_id": aid,
				"filename": filename, "mime_type": mime, "size": len(decoded),
			})),
			mcp.NewEmbeddedResource(mcp.BlobResourceContents{
				URI:      "attachment:///" + aid + "/" + filename,
				MIMEType: mime,
				Blob:     b64,
			}),
		},
	}, nil
}

type confluenceDownloadContentAttachmentsTool struct{}

func (t confluenceDownloadContentAttachmentsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.DownloadContentAttachmentsBase64(ctx, base.Args.String("content_id"))
	if err != nil {
		return nil, err
	}
	contents := []mcp.Content{mcp.NewTextContent(core.Stringify(map[string]any{
		"success": r["success"], "content_id": r["content_id"],
		"total": r["total"], "downloaded": r["downloaded"],
		"failed": r["failed"],
	}))}
	attachments, _ := r["attachments"].([]map[string]any)
	for _, att := range attachments {
		b64, _ := att["data_base64"].(string)
		mime, _ := att["mime_type"].(string)
		filename, _ := att["filename"].(string)
		aid, _ := att["attachment_id"].(string)
		if b64 == "" {
			continue
		}
		contents = append(contents, mcp.NewEmbeddedResource(mcp.BlobResourceContents{
			URI:      "attachment:///" + aid + "/" + filename,
			MIMEType: mime,
			Blob:     b64,
		}))
	}
	return &mcp.CallToolResult{Content: contents}, nil
}

type confluenceGetPageImagesTool struct{}

func (t confluenceGetPageImagesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.GetPageImages(ctx, base.Args.String("content_id"))
	if err != nil {
		return nil, err
	}
	contents := []mcp.Content{mcp.NewTextContent(core.Stringify(map[string]any{
		"success":      r["success"],
		"content_id":   r["content_id"],
		"total":        r["total"],
		"total_images": r["total"],
		"downloaded":   r["images"],
		"failed":       r["failed"],
	}))}
	fetched, _ := r["fetched"].([]map[string]any)
	for _, img := range fetched {
		b64, _ := img["data_base64"].(string)
		mime, _ := img["mime_type"].(string)
		if b64 == "" {
			continue
		}
		contents = append(contents, mcp.NewImageContent(b64, mime))
	}
	return &mcp.CallToolResult{Content: contents}, nil
}

type confluenceGetPageRestrictionsTool struct{}

func (t confluenceGetPageRestrictionsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	pageID := base.Args.String("page_id")
	op := base.Args.String("operation")
	var (
		r   map[string]any
		err error
	)
	if op != "" {
		r, err = base.Confluence.GetRestrictionForOperation(ctx, pageID, op)
	} else {
		r, err = base.Confluence.GetPageRestrictions(ctx, pageID)
	}
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceScanContentTool struct{}

func (t confluenceScanContentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.ScanContent(ctx, base.Args.String("cql"), base.Args.Int("start", 0), base.Args.Int("limit", 0), base.Args.List("expand"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetContentHistoryTool struct{}

func (t confluenceGetContentHistoryTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.GetContentHistory(ctx, base.Args.String("content_id"), base.Args.List("expand"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetDescendantsTool struct{}

func (t confluenceGetDescendantsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	contentID := base.Args.String("content_id")
	typ := base.Args.String("type")
	start := base.Args.Int("start", 0)
	limit := base.Args.Int("limit", 0)
	expand := base.Args.List("expand")
	var (
		r   map[string]any
		err error
	)
	if typ != "" {
		r, err = base.Confluence.GetDescendantsOfType(ctx, contentID, typ, start, limit, expand)
	} else {
		r, err = base.Confluence.GetDescendants(ctx, contentID, start, limit, expand)
	}
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetSpaceTool struct{}

func (t confluenceGetSpaceTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.GetSpace(ctx, base.Args.String("space_key"), base.Args.List("expand"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetSpacePropertiesTool struct{}

func (t confluenceGetSpacePropertiesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.GetSpaceProperties(ctx, base.Args.String("space_key"), base.Args.List("expand"), base.Args.Int("start", 0), base.Args.Int("limit", 0))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceGetCurrentUserTool struct{}

func (t confluenceGetCurrentUserTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.GetCurrentUser(ctx, base.Args.List("expand"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceSearchEntitiesTool struct{}

func (t confluenceSearchEntitiesTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	incArchived := base.Args.Bool("include_archived_spaces", false)
	r, err := base.Confluence.SearchEntities(ctx, base.Args.String("cql"), base.Args.String("excerpt"), base.Args.String("expand"), base.Args.Int("start", 0), base.Args.Int("limit", 0), incArchived)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

func (confluenceGetPageTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_page", mcp.WithDescription("Get a Confluence page by ID, or by title and space key"),
		mcp.WithString("page_id", mcp.Description("Confluence page ID (numeric), page URL (with /pages/), or tiny link URL (with /x/)")),
		mcp.WithString("title", mcp.Description("The exact title of the page (use with space_key)")),
		mcp.WithString("space_key", mcp.Description("The key of the space (use with title)")),
		mcp.WithBoolean("include_metadata", mcp.DefaultBool(true), mcp.Description("Whether to include page metadata")),
		mcp.WithBoolean("convert_to_markdown", mcp.DefaultBool(true), mcp.Description("Convert content to markdown (true) or keep raw HTML (false)")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceSearchTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_search", mcp.WithDescription("Search Confluence content using simple terms or CQL"),
		mcp.WithString("query", mcp.Required()),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results (1-50)")),
		mcp.WithString("spaces_filter", mcp.Description("Comma-separated list of space keys to filter by")),
		mcp.WithString("search_type", mcp.Description("'cql' to pass query as CQL, 'text' for simple text search (default: auto-detect)")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetPageChildrenTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_page_children", mcp.WithDescription("Get child pages of a Confluence page"),
		mcp.WithString("parent_id", mcp.Required()),
		mcp.WithString("expand", mcp.Description("Fields to expand (e.g., 'version,body.storage')")),
		mcp.WithNumber("limit", mcp.DefaultNumber(25), mcp.Description("Maximum number of child pages (1-50)")),
		mcp.WithBoolean("include_content", mcp.Description("Whether to include page content in the response")),
		mcp.WithBoolean("include_folders", mcp.DefaultBool(true), mcp.Description("Whether to include child folders")),
		mcp.WithBoolean("convert_to_markdown", mcp.DefaultBool(true), mcp.Description("Convert content to markdown if include_content is true")),
		mcp.WithNumber("start", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination (0-based)")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetSpacePageTreeTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_space_page_tree", mcp.WithDescription("Get space page tree"),
		mcp.WithString("space_key", mcp.Required()),
		mcp.WithNumber("limit", mcp.DefaultNumber(100), mcp.Description("Maximum number of pages")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetCommentsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_comments", mcp.WithDescription("Get page comments as a flat list; each comment has parent_id (empty for top-level) so replies can be re-nested"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithString("format", mcp.Description("Response format: 'simplified' (default) or 'raw'")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetLabelsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_labels", mcp.WithDescription("Get page labels"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithString("prefix", mcp.Description("Filter by label prefix (e.g., 'global')")),
		mcp.WithNumber("start", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithNumber("limit", mcp.DefaultNumber(200), mcp.Description("Maximum number of results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceSearchUserTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_search_user", mcp.WithDescription("Search Confluence users"),
		mcp.WithString("query", mcp.Required()),
		mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Description("Maximum number of results")),
		mcp.WithString("group_name", mcp.Description("Group name to search within (default 'confluence-users')")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetAttachmentsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_attachments", mcp.WithDescription("Get attachments for Confluence content"),
		mcp.WithString("content_id", mcp.Required()),
		mcp.WithString("expand", mcp.Description("Fields to expand (e.g., 'version,extensions')")),
		mcp.WithNumber("start", mcp.DefaultNumber(0), mcp.Description("Starting index for pagination")),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Description("Maximum number of results")),
		mcp.WithString("filename", mcp.Description("Filter by exact filename match")),
		mcp.WithString("media_type", mcp.Description("Filter by media type (e.g., 'image/png')")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetInlineCommentsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_inline_comments", mcp.WithDescription("Get inline comments for a Confluence page; each comment has parent_id (empty for top-level) so replies can be re-nested"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithString("result_format", mcp.Description("Response format: 'simplified' (default) or 'raw'")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetPageHistoryTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_page_history", mcp.WithDescription("Get a historical version of a Confluence page"),
		mcp.WithString("page_id", mcp.Required()), mcp.WithNumber("version", mcp.Required()),
		mcp.WithBoolean("convert_to_markdown", mcp.DefaultBool(true), mcp.Description("Convert content to markdown (true) or keep raw HTML (false)")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetPageDiffTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_page_diff", mcp.WithDescription("Get unified diff between two versions of a Confluence page"),
		mcp.WithString("page_id", mcp.Required()), mcp.WithNumber("from_version", mcp.Required()), mcp.WithNumber("to_version", mcp.Required()),
		mcp.WithBoolean("convert_to_markdown", mcp.DefaultBool(true), mcp.Description("Convert content to markdown (true) or keep raw HTML (false)")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceDownloadAttachmentTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_download_attachment", mcp.WithDescription("Download a Confluence attachment as base64-encoded resource"),
		mcp.WithString("attachment_id", mcp.Required()),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceDownloadContentAttachmentsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_download_content_attachments", mcp.WithDescription("Download all attachments for Confluence content as base64 resources"),
		mcp.WithString("content_id", mcp.Required()),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetPageImagesTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_page_images", mcp.WithDescription("Get all images attached to a Confluence page as base64"),
		mcp.WithString("content_id", mcp.Required()),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetPageRestrictionsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_page_restrictions", mcp.WithDescription("Get view and edit restrictions for a Confluence page"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithString("operation", mcp.Description("Restriction operation key: 'read' or 'update'. Omit to get all operations")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceScanContentTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_scan_content", mcp.WithDescription("Scan for Confluence content matching a CQL query"),
		mcp.WithString("cql", mcp.Description("CQL query string")),
		mcp.WithNumber("start", mcp.DefaultNumber(0), mcp.Description("Starting index")),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Description("Maximum number of results")),
		mcp.WithArray("expand", mcp.Description("Fields to expand")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetContentHistoryTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_content_history", mcp.WithDescription("Get the history object for a Confluence content item"),
		mcp.WithString("content_id", mcp.Required()),
		mcp.WithArray("expand", mcp.Description("Fields to expand")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetDescendantsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_descendants", mcp.WithDescription("Get all descendants of a Confluence content item"),
		mcp.WithString("content_id", mcp.Required()),
		mcp.WithString("type", mcp.Description("Descendant type filter (e.g., 'page'). Omit for all types")),
		mcp.WithNumber("start", mcp.DefaultNumber(0), mcp.Description("Starting index")),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Description("Maximum number of results")),
		mcp.WithArray("expand", mcp.Description("Fields to expand")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetSpaceTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_space", mcp.WithDescription("Get a Confluence space by its key"),
		mcp.WithString("space_key", mcp.Required()),
		mcp.WithArray("expand", mcp.Description("Fields to expand")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetSpacePropertiesTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_space_properties", mcp.WithDescription("Get properties for a Confluence space"),
		mcp.WithString("space_key", mcp.Required()),
		mcp.WithArray("expand", mcp.Description("Fields to expand")),
		mcp.WithNumber("start", mcp.DefaultNumber(0), mcp.Description("Starting index")),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Description("Maximum number of results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceGetCurrentUserTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_get_current_user", mcp.WithDescription("Get the authenticated Confluence user's profile"),
		mcp.WithArray("expand", mcp.Description("Fields to expand")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func (confluenceSearchEntitiesTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_search_entities", mcp.WithDescription("Search across all Confluence entities using CQL"),
		mcp.WithString("cql", mcp.Description("CQL query string")),
		mcp.WithString("excerpt", mcp.Description("Excerpt strategy: 'default', 'none', or 'indexed'")),
		mcp.WithString("expand", mcp.Description("Fields to expand (comma-separated)")),
		mcp.WithNumber("start", mcp.DefaultNumber(0), mcp.Description("Starting index")),
		mcp.WithNumber("limit", mcp.DefaultNumber(50), mcp.Description("Maximum number of results")),
		mcp.WithBoolean("include_archived_spaces", mcp.Description("Include archived spaces in results")),
		mcp.WithReadOnlyHintAnnotation(true), mcp.WithDestructiveHintAnnotation(false),
	)
}

func ReadTools() []core.Runnable {
	return []core.Runnable{
		&confluenceGetPageTool{},
		&confluenceSearchTool{},
		&confluenceGetPageChildrenTool{},
		&confluenceGetSpacePageTreeTool{},
		&confluenceGetCommentsTool{},
		&confluenceGetLabelsTool{},
		&confluenceSearchUserTool{},
		&confluenceGetAttachmentsTool{},
		&confluenceGetInlineCommentsTool{},
		&confluenceGetPageHistoryTool{},
		&confluenceGetPageDiffTool{},
		&confluenceDownloadAttachmentTool{},
		&confluenceDownloadContentAttachmentsTool{},
		&confluenceGetPageImagesTool{},
		&confluenceGetPageRestrictionsTool{},
		&confluenceScanContentTool{},
		&confluenceGetContentHistoryTool{},
		&confluenceGetDescendantsTool{},
		&confluenceGetSpaceTool{},
		&confluenceGetSpacePropertiesTool{},
		&confluenceGetCurrentUserTool{},
		&confluenceSearchEntitiesTool{},
	}
}
