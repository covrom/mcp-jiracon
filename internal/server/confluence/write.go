package confluence

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	core "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/server/core"

	"github.com/mark3labs/mcp-go/mcp"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/preprocessing"
)

type confluenceCreatePageTool struct{}

func (t confluenceCreatePageTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	a := base.Args
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	contentFormat := a.String("content_format")
	if contentFormat == "" {
		contentFormat = "markdown"
	}
	// Validate content_format
	if contentFormat != "markdown" && contentFormat != "wiki" && contentFormat != "storage" && contentFormat != "xhtml" {
		return nil, fmt.Errorf("invalid content_format: %s. Must be 'markdown', 'wiki', 'storage', or 'xhtml'", contentFormat)
	}
	// Resolve content from either 'content' or 'content_file'
	content := a.String("content")
	contentFile := a.String("content_file")
	if content != "" && contentFile != "" {
		return nil, fmt.Errorf("provide either 'content' or 'content_file', not both")
	}
	if content == "" && contentFile == "" {
		return nil, fmt.Errorf("one of 'content' or 'content_file' must be provided")
	}
	if contentFile != "" {
		data, err := os.ReadFile(contentFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read content_file %s: %w", contentFile, err)
		}
		content = string(data)
	}
	// Determine representation and body based on content format
	var representation string
	var bodyValue string
	if contentFormat == "markdown" {
		enableHeadingAnchors := a.Bool("enable_heading_anchors", false)
		tableLayout := a.String("table_layout")
		prep := preprocessing.NewConfluencePreprocessor()
		bodyValue = prep.MarkdownToStorageWithTableLayout(content, enableHeadingAnchors, tableLayout)
		representation = "storage"
	} else if contentFormat == "xhtml" {
		representation = "storage"
		bodyValue = content
	} else {
		// 'wiki' or 'storage'
		representation = contentFormat
		bodyValue = content
	}
	body := map[string]any{
		"type":  "page",
		"title": a.String("title"),
		"space": map[string]any{"key": a.String("space_key")},
		"body":  map[string]any{representation: map[string]any{"value": bodyValue, "representation": representation}},
	}
	if parentID := a.String("parent_id"); parentID != "" {
		body["ancestors"] = []map[string]any{{"id": parentID}}
	}
	// Optional metadata fields.
	if emoji := a.String("emoji"); emoji != "" {
		body["metadata"] = map[string]any{"properties": map[string]any{"emoji": map[string]any{"value": emoji}}}
	}
	if pw := a.String("page_width"); pw != "" {
		if body["metadata"] == nil {
			body["metadata"] = map[string]any{"properties": map[string]any{}}
		}
		body["metadata"].(map[string]any)["properties"].(map[string]any)["page-width"] = map[string]any{"value": pw}
	}
	includeContent := a.Bool("include_content", false)
	var raw map[string]any
	if err := base.Confluence.Client().Post(ctx, base.Confluence.V1BaseURL()+"/content", body, &raw); err != nil {
		return nil, err
	}
	if !includeContent {
		return mcp.NewToolResultText(core.Stringify(map[string]any{"message": "Page created successfully", "page_id": raw["id"], "title": a.String("title")})), nil
	}
	return mcp.NewToolResultText(core.Stringify(raw)), nil
}

type confluenceUpdatePageTool struct{}

func (t confluenceUpdatePageTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	a := base.Args
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	pageID := a.String("page_id")
	// Fetch current page to get the version number
	currentPage, err := base.Confluence.GetPage(ctx, pageID, nil)
	if err != nil {
		return nil, fmt.Errorf("confluence update: failed to get current page version: %w", err)
	}
	versionNum := int64(1)
	if currentPage.Version != nil {
		if n, ok := currentPage.Version["number"].(float64); ok {
			versionNum = int64(n) + 1
		}
	}
	isMinorEdit := a.Bool("is_minor_edit", false)
	versionComment := a.String("version_comment")
	versionData := map[string]any{"number": versionNum, "minorEdit": isMinorEdit}
	if versionComment != "" {
		versionData["message"] = versionComment
	}
	contentFormat := a.String("content_format")
	if contentFormat == "" {
		contentFormat = "markdown"
	}
	// Validate content_format
	if contentFormat != "markdown" && contentFormat != "wiki" && contentFormat != "storage" && contentFormat != "xhtml" {
		return nil, fmt.Errorf("invalid content_format: %s. Must be 'markdown', 'wiki', 'storage', or 'xhtml'", contentFormat)
	}
	// Resolve content from either 'content' or 'content_file'
	content := a.String("content")
	contentFile := a.String("content_file")
	if content != "" && contentFile != "" {
		return nil, fmt.Errorf("provide either 'content' or 'content_file', not both")
	}
	if content == "" && contentFile == "" {
		return nil, fmt.Errorf("one of 'content' or 'content_file' must be provided")
	}
	if contentFile != "" {
		data, err := os.ReadFile(contentFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read content_file %s: %w", contentFile, err)
		}
		content = string(data)
	}
	// Determine representation and body based on content format
	var representation string
	var bodyValue string
	if contentFormat == "markdown" {
		enableHeadingAnchors := a.Bool("enable_heading_anchors", false)
		tableLayout := a.String("table_layout")
		prep := preprocessing.NewConfluencePreprocessor()
		bodyValue = prep.MarkdownToStorageWithTableLayout(content, enableHeadingAnchors, tableLayout)
		representation = "storage"
	} else if contentFormat == "xhtml" {
		representation = "storage"
		bodyValue = content
	} else {
		// 'wiki' or 'storage'
		representation = contentFormat
		bodyValue = content
	}
	body := map[string]any{
		"id":      pageID,
		"type":    "page",
		"title":   a.String("title"),
		"body":    map[string]any{representation: map[string]any{"value": bodyValue, "representation": representation}},
		"version": versionData,
	}
	if parentID := a.String("parent_id"); parentID != "" {
		body["ancestors"] = []map[string]any{{"id": parentID}}
	}
	// Optional metadata fields.
	if emoji := a.String("emoji"); emoji != "" {
		body["metadata"] = map[string]any{"properties": map[string]any{"emoji": map[string]any{"value": emoji}}}
	}
	if pw := a.String("page_width"); pw != "" {
		if body["metadata"] == nil {
			body["metadata"] = map[string]any{"properties": map[string]any{}}
		}
		body["metadata"].(map[string]any)["properties"].(map[string]any)["page-width"] = map[string]any{"value": pw}
	}
	var raw map[string]any
	if err := base.Confluence.Client().Put(ctx, base.Confluence.V1BaseURL()+"/content/"+url.PathEscape(pageID), body, &raw); err != nil {
		return nil, err
	}
	includeContent := a.Bool("include_content", false)
	if !includeContent {
		if bodyResult, ok := raw["body"].(map[string]any); ok {
			if storage, ok2 := bodyResult["storage"].(map[string]any); ok2 {
				delete(storage, "value")
			}
		}
	}
	return mcp.NewToolResultText(core.Stringify(raw)), nil
}

type confluenceDeletePageTool struct{}

func (t confluenceDeletePageTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	pageID := base.Args.String("page_id")
	if err := base.Confluence.Client().Delete(ctx, base.Confluence.V1BaseURL()+"/content/"+url.PathEscape(pageID)); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"success": true, "message": fmt.Sprintf("Page %s deleted successfully", pageID),
	})), nil
}

type confluenceMovePageTool struct{}

func (t confluenceMovePageTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	pageID := url.PathEscape(base.Args.String("page_id"))
	targetParentID := base.Args.String("target_parent_id")
	targetSpaceKey := base.Args.String("target_space_key")
	position := base.Args.String("position")
	if position == "" {
		position = "append"
	}

	var path string
	body := make(map[string]any)
	if targetParentID == "" {
		path = fmt.Sprintf("%s/content/%s/move/%s", base.Confluence.V1BaseURL(), pageID, position)
		if targetSpaceKey != "" {
			body["destination"] = map[string]any{"type": "page", "space": map[string]any{"key": targetSpaceKey}}
		}
	} else {
		path = fmt.Sprintf("%s/content/%s/move/%s/%s", base.Confluence.V1BaseURL(), pageID, position, url.PathEscape(targetParentID))
		body["destination"] = map[string]any{"id": targetParentID, "type": "page"}
		if targetSpaceKey != "" {
			body["destination"].(map[string]any)["space"] = map[string]any{"key": targetSpaceKey}
		}
	}
	var raw map[string]any
	if err := base.Confluence.Client().Put(ctx, path, body, &raw); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(raw)), nil
}

type confluenceAddCommentTool struct{}

func (t confluenceAddCommentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	// Support both 'body' (Python-compatible) and 'content' (legacy Go) parameter names.
	commentBody := base.Args.String("body")
	if commentBody == "" {
		commentBody = base.Args.String("content")
	}
	c, err := base.Confluence.AddComment(ctx, base.Args.String("page_id"), commentBody)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"success": true, "message": "Comment added successfully", "comment": c.ToSimplifiedDict(),
	})), nil
}

type confluenceReplyToCommentTool struct{}

func (t confluenceReplyToCommentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.ReplyToComment(ctx, base.Args.String("comment_id"), base.Args.String("body"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"success": true, "message": "Reply to comment added successfully", "reply": r,
	})), nil
}

type confluenceAddLabelTool struct{}

func (t confluenceAddLabelTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	r, err := base.Confluence.AddPageLabel(ctx, base.Args.String("page_id"), base.Args.String("name"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceDeleteAttachmentTool struct{}

func (t confluenceDeleteAttachmentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	if err := base.Confluence.DeleteAttachment(ctx, base.Args.String("attachment_id")); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(`{"success":true}`), nil
}

type confluenceUpdatePageSectionTool struct{}

func (t confluenceUpdatePageSectionTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	contentFormat := base.Args.String("content_format")
	if contentFormat == "" {
		contentFormat = "markdown"
	}
	isMinorEdit := base.Args.Bool("is_minor_edit", false)
	versionComment := base.Args.String("version_comment")
	r, err := base.Confluence.UpdatePageSection(ctx,
		base.Args.String("page_id"),
		base.Args.String("heading_text"),
		base.Args.String("new_content"),
		contentFormat,
		isMinorEdit,
		versionComment,
	)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceAddInlineCommentTool struct{}

func (t confluenceAddInlineCommentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	matchCount := 1
	if v := base.Args.Int("text_selection_match_count", 0); v > 0 {
		matchCount = v
	}
	matchIndex := 0
	if v := base.Args.Int("text_selection_match_index", 0); v >= 0 {
		matchIndex = v
	}
	r, err := base.Confluence.AddInlineComment(ctx,
		base.Args.String("page_id"),
		base.Args.String("body"),
		base.Args.String("text_selection"),
		matchCount,
		matchIndex,
	)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"success": true, "message": "Inline comment added successfully", "comment": r,
	})), nil
}

type confluenceUploadAttachmentTool struct{}

func (t confluenceUploadAttachmentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	comment := base.Args.String("comment")
	minorEdit := base.Args.Bool("minor_edit", false)

	// Handle base64 upload.
	if b64Content := base.Args.String("content_base64"); b64Content != "" {
		filename := base.Args.String("filename")
		if filename == "" {
			filename = "upload.bin"
		}
		r, err := base.Confluence.UploadAttachmentBase64(ctx, base.Args.String("content_id"), filename, b64Content, comment, minorEdit)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(core.Stringify(r)), nil
	}

	// File-based upload.
	filePath := base.Args.String("file_path")
	if filePath == "" {
		return nil, fmt.Errorf("either file_path or content_base64+filename must be provided")
	}
	r, err := base.Confluence.UploadAttachment(ctx, base.Args.String("content_id"), filePath, comment, minorEdit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceUploadAttachmentsTool struct{}

func (t confluenceUploadAttachmentsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	comment := base.Args.String("comment")
	minorEdit := base.Args.Bool("minor_edit", false)
	paths := core.SplitComma(base.Args.String("file_paths"))
	r, err := base.Confluence.UploadAttachments(ctx, base.Args.String("content_id"), paths, comment, minorEdit)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(r)), nil
}

type confluenceSetPageRestrictionsTool struct{}

func (t confluenceSetPageRestrictionsTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	// Parse JSON arrays from string params.
	parseJSONStringArray := func(key string) []string {
		s := base.Args.String(key)
		if s == "" {
			return nil
		}
		var arr []string
		if err := json.Unmarshal([]byte(s), &arr); err != nil {
			// Fall back to comma-separated
			return core.SplitComma(s)
		}
		return arr
	}
	r, err := base.Confluence.SetPageRestrictions(ctx,
		base.Args.String("page_id"),
		parseJSONStringArray("read_users"),
		parseJSONStringArray("read_groups"),
		parseJSONStringArray("edit_users"),
		parseJSONStringArray("edit_groups"),
	)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"message": "Page restrictions updated successfully", "restrictions": r,
	})), nil
}

type confluenceCopyPageTool struct{}

func (t confluenceCopyPageTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	copyAttachments := base.Args.Bool("copy_attachments", true)
	if copyAttachments {
		// DC does not support automatic attachment copying.
		copyAttachments = false // signal to fetcher
	}
	r, err := base.Confluence.CopyPage(ctx,
		base.Args.String("source_page_id"),
		base.Args.String("destination_space_key"),
		base.Args.String("new_title"),
		base.Args.String("destination_parent_id"),
		copyAttachments,
	)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{
		"message": "Page copied successfully", "page": r,
	})), nil
}

type confluenceRemoveLabelTool struct{}

func (t confluenceRemoveLabelTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	pageID := base.Args.String("page_id")
	label := base.Args.String("label")
	var err error
	if strings.Contains(label, "/") {
		err = base.Confluence.RemovePageLabel(ctx, pageID, label)
	} else {
		err = base.Confluence.RemovePageLabelByName(ctx, pageID, label)
	}
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"success": true, "message": fmt.Sprintf("Label %q removed from page %s", label, pageID)})), nil
}

type confluenceSetPagePropertyTool struct{}

func (t confluenceSetPagePropertyTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	pageID := base.Args.String("page_id")
	key := base.Args.String("key")
	valueStr := base.Args.String("value")
	var value any
	if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
		value = valueStr
	}
	// Try to get existing property to determine create vs update.
	existing, gerr := base.Confluence.GetContentProperty(ctx, pageID, key)
	if gerr != nil {
		// Property does not exist — create it.
		r, cerr := base.Confluence.CreateContentProperty(ctx, pageID, key, value)
		if cerr != nil {
			return nil, cerr
		}
		return mcp.NewToolResultText(core.Stringify(map[string]any{"success": true, "action": "created", "property": r})), nil
	}
	// Property exists — update it.
	id := int(models.MapInt(existing, "id"))
	versionMap, _ := existing["version"].(map[string]any)
	version := int(models.MapInt(versionMap, "number"))
	r, uerr := base.Confluence.UpdateContentProperty(ctx, pageID, key, id, version, value)
	if uerr != nil {
		return nil, uerr
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"success": true, "action": "updated", "property": r})), nil
}

type confluenceDeletePagePropertyTool struct{}

func (t confluenceDeletePagePropertyTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	pageID := base.Args.String("page_id")
	key := base.Args.String("key")
	if err := base.Confluence.DeleteContentProperty(ctx, pageID, key); err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"success": true, "message": fmt.Sprintf("Property %q deleted from page %s", key, pageID)})), nil
}

type confluenceUpdateAttachmentTool struct{}

func (t confluenceUpdateAttachmentTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	attID := base.Args.String("attachment_id")
	r, err := base.Confluence.UpdateAttachmentMetadata(ctx, attID, base.Args.String("title"), base.Args.String("comment"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"success": true, "attachment": r})), nil
}

type confluenceUpdateAttachmentDataTool struct{}

func (t confluenceUpdateAttachmentDataTool) Run(ctx context.Context, base *core.Tool) (*mcp.CallToolResult, error) {
	if err := base.GuardWrite(); err != nil {
		return nil, err
	}
	if err := base.RequireConfluence(); err != nil {
		return nil, err
	}
	attID := base.Args.String("attachment_id")
	filename := base.Args.String("filename")
	b64 := base.Args.String("content_base64")
	data, derr := base64.StdEncoding.DecodeString(b64)
	if derr != nil {
		return nil, fmt.Errorf("invalid base64 content: %w", derr)
	}
	r, err := base.Confluence.UpdateAttachmentData(ctx, attID, filename, data, base.Args.String("comment"))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(core.Stringify(map[string]any{"success": true, "attachment": r})), nil
}

func (confluenceCreatePageTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_create_page", mcp.WithDescription("Create a Confluence page"),
		mcp.WithString("space_key", mcp.Required(), mcp.Description("The key of the space to create the page in (usually a short uppercase code like 'DEV', 'TEAM', or 'DOC')")),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the page")),
		mcp.WithString("content", mcp.Description("The content of the page. Format depends on content_format parameter. Can be Markdown (default), wiki markup, storage format, or XHTML storage format. Either 'content' or 'content_file' must be provided, but not both.")),
		mcp.WithString("parent_id", mcp.Description("Parent page ID for hierarchical placement")),
		mcp.WithString("content_format", mcp.Description("(Optional) The format of the content parameter. Options: 'markdown' (default), 'wiki', 'storage', or 'xhtml'. Use 'xhtml' when providing Confluence XHTML storage format (same as 'storage'). Wiki format uses Confluence wiki markup syntax")),
		mcp.WithBoolean("enable_heading_anchors", mcp.Description("(Optional) Whether to enable automatic heading anchor generation. Only applies when content_format is 'markdown'")),
		mcp.WithBoolean("include_content", mcp.Description("(Optional) Whether to include page content in the response. Defaults to false since callers already have the content at create time")),
		mcp.WithString("emoji", mcp.Description("(Optional) Page title emoji (icon shown in navigation). Can be any emoji character like '📝', '🚀', '📚'. Set to null/None to remove.")),
		mcp.WithString("content_file", mcp.Description("(Optional) Absolute or relative filesystem path to read the page body from (UTF-8). Use this instead of 'content' when the body is too large to pass comfortably as a tool argument. Mutually exclusive with 'content'.")),
		mcp.WithString("page_width", mcp.Description("(Optional) Page layout width. Options: 'full-width', 'default'. Defaults to null (Confluence default).")),
		mcp.WithString("table_layout", mcp.Description("(Optional) Table width preset applied to all markdown tables. Options: 'full-width' (1800 px), 'wide' (960 px), 'default' (760 px). Only applies when content_format is 'markdown'.")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceUpdatePageTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_update_page", mcp.WithDescription("Update an existing Confluence page"),
		mcp.WithString("page_id", mcp.Required(), mcp.Description("The ID of the page to update")),
		mcp.WithString("title", mcp.Required(), mcp.Description("The new title of the page")),
		mcp.WithString("content", mcp.Description("The new content of the page. Format depends on content_format parameter and may be Markdown (default), wiki markup, storage format, or XHTML storage format. Either 'content' or 'content_file' must be provided, but not both.")),
		mcp.WithBoolean("is_minor_edit", mcp.Description("Whether this is a minor edit")),
		mcp.WithString("version_comment", mcp.Description("Optional comment for this version")),
		mcp.WithString("parent_id", mcp.Description("Optional new parent page ID")),
		mcp.WithString("content_format", mcp.Description("(Optional) The format of the content parameter. Options: 'markdown' (default), 'wiki', 'storage', or 'xhtml'. Use 'xhtml' when providing Confluence XHTML storage format (same as 'storage'). Wiki format uses Confluence wiki markup syntax")),
		mcp.WithBoolean("enable_heading_anchors", mcp.Description("(Optional) Whether to enable automatic heading anchor generation. Only applies when content_format is 'markdown'")),
		mcp.WithBoolean("include_content", mcp.Description("(Optional) Whether to include page content in the response. Defaults to false since callers already have the content at update time")),
		mcp.WithString("emoji", mcp.Description("(Optional) Page title emoji (icon shown in navigation). Can be any emoji character like '📝', '🚀', '📚'. Set to null/None to remove.")),
		mcp.WithString("content_file", mcp.Description("(Optional) Absolute or relative filesystem path to read the new page body from (UTF-8). Use this instead of 'content' when the body is too large to pass comfortably as a tool argument. Mutually exclusive with 'content'.")),
		mcp.WithString("page_width", mcp.Description("(Optional) Page layout width. Options: 'full-width', 'default'. Defaults to null (preserve existing).")),
		mcp.WithString("table_layout", mcp.Description("(Optional) Table width preset applied to all markdown tables. Options: 'full-width' (1800 px), 'wide' (960 px), 'default' (760 px). Only applies when content_format is 'markdown'.")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceDeletePageTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_delete_page", mcp.WithDescription("Delete a Confluence page"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceMovePageTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_move_page", mcp.WithDescription("Move a Confluence page to a new parent or space"),
		mcp.WithString("page_id", mcp.Required()), mcp.WithString("target_parent_id", mcp.Description("Target parent page ID. When omitted the page is moved to the space root.")),
		mcp.WithString("target_space_key", mcp.Description("Space key for cross-space moves")),
		mcp.WithString("position", mcp.Description("Position relative to target parent: 'append' (default), 'above', 'below'")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceAddCommentTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_add_comment", mcp.WithDescription("Add comment to a Confluence page"),
		mcp.WithString("page_id", mcp.Required()), mcp.WithString("body", mcp.Required()),
		mcp.WithString("content", mcp.Description("Alias for body (deprecated, use 'body' instead)")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceReplyToCommentTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_reply_to_comment", mcp.WithDescription("Reply to a Confluence comment"),
		mcp.WithString("comment_id", mcp.Required()), mcp.WithString("body", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceAddLabelTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_add_label", mcp.WithDescription("Add label to Confluence page"),
		mcp.WithString("page_id", mcp.Required()), mcp.WithString("name", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceDeleteAttachmentTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_delete_attachment", mcp.WithDescription("Delete Confluence attachment"),
		mcp.WithString("attachment_id", mcp.Required()),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceUpdatePageSectionTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_update_page_section", mcp.WithDescription("Update a single section of a Confluence page without affecting the rest"),
		mcp.WithString("page_id", mcp.Required()), mcp.WithString("heading_text", mcp.Required()), mcp.WithString("new_content", mcp.Required()),
		mcp.WithBoolean("is_minor_edit", mcp.Description("Whether this is a minor edit")),
		mcp.WithString("version_comment", mcp.Description("Optional comment for this version")),
		mcp.WithString("content_format", mcp.Description("(Optional) Format of new_content. Options: 'markdown' (default) or 'storage' (raw Confluence storage XML).")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceAddInlineCommentTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_add_inline_comment", mcp.WithDescription("Add an inline comment anchored to a text selection on a Confluence page"),
		mcp.WithString("page_id", mcp.Required()), mcp.WithString("body", mcp.Required()), mcp.WithString("text_selection", mcp.Required()),
		mcp.WithNumber("text_selection_match_count", mcp.Description("Total number of times the selected text appears on the page. Defaults to 1.")),
		mcp.WithNumber("text_selection_match_index", mcp.Description("Zero-based index of which occurrence of the text to anchor to. Defaults to 0 (first occurrence).")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceUploadAttachmentTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_upload_attachment", mcp.WithDescription("Upload a file as an attachment to Confluence content"),
		mcp.WithString("content_id", mcp.Required()),
		mcp.WithString("file_path", mcp.Description("Local path to the file to upload")),
		mcp.WithString("content_base64", mcp.Description("Base64-encoded file content (alternative to file_path)")),
		mcp.WithString("filename", mcp.Description("Filename for base64 upload (required if content_base64 is provided)")),
		mcp.WithString("comment", mcp.Description("Comment for the attachment")),
		mcp.WithBoolean("minor_edit", mcp.Description("Whether this is a minor edit")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceUploadAttachmentsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_upload_attachments", mcp.WithDescription("Upload multiple files as attachments to Confluence content"),
		mcp.WithString("content_id", mcp.Required()), mcp.WithString("file_paths", mcp.Required()),
		mcp.WithString("comment", mcp.Description("Comment for all attachments")),
		mcp.WithBoolean("minor_edit", mcp.Description("Whether this is a minor edit")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceSetPageRestrictionsTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_set_page_restrictions", mcp.WithDescription("Set view and edit restrictions on a Confluence page"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithString("read_users", mcp.Description("JSON array of usernames for read access, e.g. '[\"user1\",\"user2\"]'")),
		mcp.WithString("read_groups", mcp.Description("JSON array of group names for read access")),
		mcp.WithString("edit_users", mcp.Description("JSON array of usernames for edit access")),
		mcp.WithString("edit_groups", mcp.Description("JSON array of group names for edit access")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceCopyPageTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_copy_page", mcp.WithDescription("Copy a Confluence page to a new location (DC: manual GET+POST)"),
		mcp.WithString("source_page_id", mcp.Required()), mcp.WithString("destination_space_key", mcp.Required()), mcp.WithString("new_title", mcp.Required()),
		mcp.WithString("destination_parent_id", mcp.Description("(Optional) Parent page ID in the destination space. When omitted the page is created at the space root.")),
		mcp.WithBoolean("copy_attachments", mcp.Description("Whether to copy attachments (DC limitation: attachments are not copied automatically)")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceRemoveLabelTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_remove_label", mcp.WithDescription("Remove a label from a Confluence page"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithString("label", mcp.Description("Label name (e.g., 'mylabel') or full path with prefix (e.g., 'global/mylabel')")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceSetPagePropertyTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_set_page_property", mcp.WithDescription("Create or update a content property on a Confluence page"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithString("key", mcp.Required(), mcp.Description("Property key")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Property value (JSON string)")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceDeletePagePropertyTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_delete_page_property", mcp.WithDescription("Delete a content property from a Confluence page"),
		mcp.WithString("page_id", mcp.Required()),
		mcp.WithString("key", mcp.Required(), mcp.Description("Property key to delete")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceUpdateAttachmentTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_update_attachment", mcp.WithDescription("Update metadata (title, comment) for a Confluence attachment"),
		mcp.WithString("attachment_id", mcp.Required()),
		mcp.WithString("title", mcp.Description("New title for the attachment")),
		mcp.WithString("comment", mcp.Description("Comment for the update")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func (confluenceUpdateAttachmentDataTool) McpTool() mcp.Tool {
	return mcp.NewTool("confluence_update_attachment_data", mcp.WithDescription("Replace the data of an existing Confluence attachment"),
		mcp.WithString("attachment_id", mcp.Required()),
		mcp.WithString("filename", mcp.Required(), mcp.Description("Filename for the new data")),
		mcp.WithString("content_base64", mcp.Required(), mcp.Description("Base64-encoded file content")),
		mcp.WithString("comment", mcp.Description("Comment for the update")),
		mcp.WithDestructiveHintAnnotation(true),
	)
}

func WriteTools() []core.Runnable {
	return []core.Runnable{
		&confluenceCreatePageTool{},
		&confluenceUpdatePageTool{},
		&confluenceDeletePageTool{},
		&confluenceMovePageTool{},
		&confluenceAddCommentTool{},
		&confluenceReplyToCommentTool{},
		&confluenceAddLabelTool{},
		&confluenceDeleteAttachmentTool{},
		&confluenceUpdatePageSectionTool{},
		&confluenceAddInlineCommentTool{},
		&confluenceUploadAttachmentTool{},
		&confluenceUploadAttachmentsTool{},
		&confluenceSetPageRestrictionsTool{},
		&confluenceCopyPageTool{},
		&confluenceRemoveLabelTool{},
		&confluenceSetPagePropertyTool{},
		&confluenceDeletePagePropertyTool{},
		&confluenceUpdateAttachmentTool{},
		&confluenceUpdateAttachmentDataTool{},
	}
}
