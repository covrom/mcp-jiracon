package confluence

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/htmlconv"
	confmodels "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/confluence"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/preprocessing"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/utils"
	"golang.org/x/net/html"
)

// GetComments returns comments for a page. When inlineOnly is true, only inline
// comments (extensions.location == "inline") are returned and the response includes
// extensions.inlineProperties. Mirrors the old GetPageComments/GetInlineComments pair.
func (f *Fetcher) GetComments(ctx context.Context, pageID string, inlineOnly bool) ([]map[string]any, error) {
	expand := "body.view.value,version,ancestors"
	if inlineOnly {
		expand = "body.view.value,version,ancestors,extensions.inlineProperties"
	}
	q := url.Values{}
	q.Set("expand", expand)
	q.Set("depth", "all")
	var raw struct {
		Results []map[string]any `json:"results"`
	}
	if err := f.getJSON(ctx, f.contentPath(pageID, "child", "comment"), q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get comments %s: %w", pageID, err)
	}
	out := make([]map[string]any, 0, len(raw.Results))
	for _, c := range raw.Results {
		if inlineOnly {
			ext, _ := c["extensions"].(map[string]any)
			if ext == nil || ext["location"] != "inline" {
				continue
			}
		}
		out = append(out, c)
	}
	return out, nil
}

// GetPageComments returns all comments for a page (non-inline filter).
func (f *Fetcher) GetPageComments(ctx context.Context, pageID string) ([]*ConfluenceComment, error) {
	raw, err := f.GetComments(ctx, pageID, false)
	if err != nil {
		return nil, err
	}
	out := make([]*ConfluenceComment, 0, len(raw))
	for _, r := range raw {
		c := &ConfluenceComment{}
		c.FromAPIResponse(r)
		out = append(out, c)
	}
	return out, nil
}

// GetInlineComments returns only inline comments for a page.
func (f *Fetcher) GetInlineComments(ctx context.Context, pageID string) ([]map[string]any, error) {
	return f.GetComments(ctx, pageID, true)
}

func (f *Fetcher) AddComment(ctx context.Context, pageID, body string) (*ConfluenceComment, error) {
	payload := map[string]any{"type": "comment", "container": map[string]any{"id": pageID, "type": "page"}, "body": map[string]any{"storage": map[string]any{"value": body, "representation": "storage"}}}
	var raw map[string]any
	if err := f.client.Post(ctx, f.V1BaseURL()+"/content", payload, &raw); err != nil {
		return nil, fmt.Errorf("confluence: add comment: %w", err)
	}
	c := &ConfluenceComment{}
	c.FromAPIResponse(raw)
	return c, nil
}

func (f *Fetcher) GetPageLabels(ctx context.Context, contentID string, prefix string, start, limit int) ([]map[string]any, error) {
	q := url.Values{}
	if prefix != "" {
		q.Set("prefix", prefix)
	}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	var raw struct {
		Results []map[string]any `json:"results"`
	}
	if err := f.getJSON(ctx, f.contentPath(contentID, "label"), q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get labels: %w", err)
	}
	return raw.Results, nil
}

func (f *Fetcher) AddPageLabel(ctx context.Context, contentID, label string) ([]map[string]any, error) {
	body := []map[string]any{{"prefix": "global", "name": label}}
	var raw struct {
		Results []map[string]any `json:"results"`
	}
	if err := f.client.Post(ctx, fmt.Sprintf("%s/content/%s/label", f.V1BaseURL(), url.PathEscape(contentID)), body, &raw); err != nil {
		return nil, fmt.Errorf("confluence: add label: %w", err)
	}
	return raw.Results, nil
}

// SearchUser on DC uses group member API with client-side filtering.
func (f *Fetcher) SearchUser(ctx context.Context, query string, limit int, groupName string) ([]map[string]any, error) {
	if groupName == "" {
		groupName = "confluence-users"
	}
	term := strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(query), `"`), `user.fullname ~ "`)
	allUsers := []map[string]any{}
	start := 0
	for {
		path := fmt.Sprintf("%s/group/%s/member", f.V1BaseURL(), url.PathEscape(groupName))
		req, err := f.client.NewRequest(ctx, "GET", path, nil)
		if err != nil {
			return nil, err
		}
		q := req.URL.Query()
		q.Set("start", fmt.Sprintf("%d", start))
		q.Set("limit", "200")
		req.URL.RawQuery = q.Encode()
		var raw struct {
			Results []map[string]any `json:"results"`
			Size    int              `json:"size"`
		}
		if err := f.client.Do(req, &raw); err != nil {
			return nil, fmt.Errorf("confluence: search user DC: %w", err)
		}
		for _, u := range raw.Results {
			dn, _ := u["displayName"].(string)
			un, _ := u["username"].(string)
			if term == "" || strings.Contains(strings.ToLower(dn), strings.ToLower(term)) || strings.Contains(strings.ToLower(un), strings.ToLower(term)) {
				allUsers = append(allUsers, u)
			}
		}
		if len(raw.Results) == 0 || len(allUsers) >= limit {
			break
		}
		start += len(raw.Results)
	}
	if len(allUsers) > limit {
		allUsers = allUsers[:limit]
	}
	return allUsers, nil
}

func (f *Fetcher) GetPageChildren(ctx context.Context, pageID string, start, limit int, includeContent, convertToMarkdown, includeFolders bool) ([]map[string]any, error) {
	suffix := "/child/page"
	if includeFolders {
		suffix = "/child"
	}
	path := fmt.Sprintf("%s/content/%s%s", f.V1BaseURL(), url.PathEscape(pageID), suffix)
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("start", fmt.Sprintf("%d", start))
	q.Set("limit", fmt.Sprintf("%d", limit))
	if includeContent {
		q.Set("expand", "body.storage,version,space")
	} else {
		q.Set("expand", "version")
	}
	req.URL.RawQuery = q.Encode()
	var raw struct {
		Results []map[string]any `json:"results"`
	}
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get children: %w", err)
	}
	if includeContent && convertToMarkdown {
		for _, r := range raw.Results {
			body, _ := r["body"].(map[string]any)
			if body == nil {
				continue
			}
			storage, _ := body["storage"].(map[string]any)
			if storage == nil {
				continue
			}
			if v, ok := storage["value"].(string); ok {
				storage["value"] = htmlconv.Convert(v)
				storage["representation"] = "markdown"
			}
		}
	}
	return raw.Results, nil
}

func (f *Fetcher) GetSpacePageTree(ctx context.Context, spaceKey string, limit int) (map[string]any, error) {
	path := f.V1BaseURL() + "/content"
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("spaceKey", spaceKey)
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("expand", "ancestors,version")
	req.URL.RawQuery = q.Encode()
	var raw struct {
		Results []map[string]any `json:"results"`
		Size    int              `json:"size"`
		Links   map[string]any   `json:"_links"`
	}
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("confluence: space tree: %w", err)
	}
	pages := make([]map[string]any, len(raw.Results))
	for i, r := range raw.Results {
		a, _ := r["ancestors"].([]any)
		dep := len(a)
		pid := ""
		if dep > 0 {
			if ancestor, ok := a[dep-1].(map[string]any); ok {
				pid = fmt.Sprintf("%v", ancestor["id"])
			}
		}
		pages[i] = map[string]any{"id": r["id"], "title": r["title"], "parent_id": pid, "depth": dep}
	}
	_, hasNext := raw.Links["next"]
	return map[string]any{"space_key": spaceKey, "total_pages": raw.Size, "has_more": hasNext, "pages": pages}, nil
}

func (f *Fetcher) GetContentAttachments(ctx context.Context, contentID string, start, limit int, filename, mediaType string) (map[string]any, error) {
	q := url.Values{}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if filename != "" {
		q.Set("filename", filename)
	}
	if mediaType != "" {
		q.Set("mediaType", mediaType)
	}
	var raw struct {
		Results []map[string]any `json:"results"`
		Size    int              `json:"size"`
	}
	if err := f.getJSON(ctx, f.contentPath(contentID, "child", "attachment"), q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get attachments: %w", err)
	}
	return map[string]any{"success": true, "attachments": raw.Results, "total": raw.Size}, nil
}

func (f *Fetcher) DeleteAttachment(ctx context.Context, attachmentID string) error {
	return f.client.Delete(ctx, fmt.Sprintf("%s/content/%s", f.V1BaseURL(), url.PathEscape(attachmentID)))
}

func (f *Fetcher) FetchAttachmentContent(ctx context.Context, dlURL string) ([]byte, error) {
	body, err := f.client.GetRaw(ctx, dlURL)
	if err != nil {
		return nil, fmt.Errorf("confluence: fetch attachment: %w", err)
	}
	return body, nil
}

func (f *Fetcher) GetSpaces(ctx context.Context, start, limit int) ([]map[string]any, error) {
	path := f.V1BaseURL() + "/space"
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("start", fmt.Sprintf("%d", start))
	q.Set("limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	var raw struct {
		Results []map[string]any `json:"results"`
	}
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get spaces: %w", err)
	}
	return raw.Results, nil
}

func (f *Fetcher) ReplyToComment(ctx context.Context, commentID, body string) (map[string]any, error) {
	payload := map[string]any{
		"type":      "comment",
		"ancestors": []map[string]any{{"id": commentID, "type": "comment"}},
		"body":      map[string]any{"storage": map[string]any{"value": body, "representation": "storage"}},
	}
	var raw map[string]any
	if err := f.client.Post(ctx, f.V1BaseURL()+"/content", payload, &raw); err != nil {
		return nil, fmt.Errorf("confluence: reply: %w", err)
	}
	return raw, nil
}

type ConfluenceComment struct {
	ID, Created  string
	Body, Author map[string]any
	// ParentID is the id of the direct parent comment (empty for top-level).
	// Derived from the last entry of the API "ancestors" chain.
	ParentID string
}

func (c *ConfluenceComment) FromAPIResponse(d map[string]any) {
	c.ID, c.Created = fmt.Sprintf("%v", d["id"]), fmt.Sprintf("%v", d["created"])
	if b, ok := d["body"].(map[string]any); ok {
		c.Body = b
	}
	if a, ok := d["author"].(map[string]any); ok {
		c.Author = a
	}
	if anc, ok := d["ancestors"].([]any); ok && len(anc) > 0 {
		if last, ok := anc[len(anc)-1].(map[string]any); ok {
			c.ParentID = fmt.Sprintf("%v", last["id"])
		}
	}
}
func (c *ConfluenceComment) ToSimplifiedDict() map[string]any {
	return map[string]any{"id": c.ID, "parent_id": c.ParentID, "body": c.Body, "author": c.Author, "created": c.Created}
}

// --- Page Section Update -----------------------------------------------------

// UpdatePageSection replaces content under a named heading without touching
// the rest of the page. Finds a heading tag whose text matches headingText,
// then replaces everything until the next heading of equal or higher level.
func (f *Fetcher) UpdatePageSection(ctx context.Context, pageID, headingText, newContent, contentFormat string, isMinorEdit bool, versionComment string) (map[string]any, error) {
	// 1. Fetch current page with body.storage.
	page, err := f.GetPage(ctx, pageID, []string{"body.storage", "version", "space"})
	if err != nil {
		return nil, fmt.Errorf("confluence: update section: get page: %w", err)
	}

	storageVal, _ := page.Body["storage"].(map[string]any)
	if storageVal == nil {
		return nil, fmt.Errorf("confluence: update section: page has no storage body")
	}
	rawHTML, _ := storageVal["value"].(string)

	// 2. Parse HTML and find/replace the section.
	replacedHTML, headingLevel, err := replaceSectionHTML(rawHTML, headingText, newContent, contentFormat)
	if err != nil {
		return nil, fmt.Errorf("confluence: update section: %w", err)
	}

	// 3. Update the page with new body.
	versionNum := int64(1)
	if page.Version != nil {
		if n, ok := page.Version["number"].(float64); ok {
			versionNum = int64(n) + 1
		}
	}
	spaceKey, _ := page.Space["key"].(string)

	body := map[string]any{
		"id":    pageID,
		"type":  "page",
		"title": page.Title,
		"space": map[string]any{"key": spaceKey},
		"body": map[string]any{
			"storage": map[string]any{
				"value":          replacedHTML,
				"representation": "storage",
			},
		},
		"version": map[string]any{
			"number":    versionNum,
			"minorEdit": isMinorEdit,
		},
	}
	if versionComment != "" {
		body["version"].(map[string]any)["message"] = versionComment
	}

	var raw map[string]any
	if err := f.client.Put(ctx, f.V1BaseURL()+"/content/"+url.PathEscape(pageID), body, &raw); err != nil {
		return nil, fmt.Errorf("confluence: update section: put: %w", err)
	}

	return map[string]any{
		"message":       fmt.Sprintf("Section '%s' updated successfully", headingText),
		"heading_level": headingLevel,
		"page_id":       pageID,
		"version":       versionNum,
	}, nil
}

// replaceSectionHTML finds a heading element whose inner text equals headingText
// (case-sensitive, trim space), then replaces all siblings until the next heading
// of equal or greater level with newContent. Returns the modified HTML.
func replaceSectionHTML(htmlStr, headingText, newContent, contentFormat string) (string, int, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", 0, fmt.Errorf("parse HTML: %w", err)
	}

	// Find the heading node with matching text.
	headingNode, headingLevel := findHeadingByText(doc, headingText)
	if headingNode == nil {
		return "", 0, fmt.Errorf("heading '%s' not found on the page", headingText)
	}

	// Collect siblings after the heading until next heading of <= level.
	nodesToRemove := []*html.Node{}
	for sib := headingNode.NextSibling; sib != nil; sib = sib.NextSibling {
		if isHeadingOfLevel(sib, headingLevel) {
			break
		}
		nodesToRemove = append(nodesToRemove, sib)
	}

	// Remove collected nodes from parent.
	parent := headingNode.Parent
	for _, n := range nodesToRemove {
		parent.RemoveChild(n)
	}

	// Create new content node(s).
	var newNodes []*html.Node
	if contentFormat == "storage" {
		// Raw storage format — parse as HTML fragment.
		frag, fErr := html.ParseFragment(strings.NewReader(newContent), nil)
		if fErr != nil {
			return "", 0, fmt.Errorf("parse new content HTML: %w", fErr)
		}
		// Insert after the heading.
		for i := len(frag) - 1; i >= 0; i-- {
			newNodes = append([]*html.Node{frag[i]}, newNodes...)
		}
	} else {
		// Markdown — convert to HTML (valid Confluence storage format).
		prep := preprocessing.NewConfluencePreprocessor()
		htmlFragment := prep.MarkdownToStorage(newContent, false)
		frag, fErr := html.ParseFragment(strings.NewReader(htmlFragment), nil)
		if fErr != nil {
			return "", 0, fmt.Errorf("parse markdown content: %w", fErr)
		}
		for i := len(frag) - 1; i >= 0; i-- {
			newNodes = append([]*html.Node{frag[i]}, newNodes...)
		}
	}

	// Insert new nodes after the heading.
	insertAfter := headingNode
	for _, n := range newNodes {
		parent.InsertBefore(n, insertAfter.NextSibling)
		insertAfter = n
	}

	var buf bytes.Buffer
	if rErr := html.Render(&buf, doc); rErr != nil {
		return "", 0, fmt.Errorf("render HTML: %w", rErr)
	}
	return buf.String(), headingLevel, nil
}

func findHeadingByText(n *html.Node, text string) (*html.Node, int) {
	if n.Type == html.ElementNode {
		tag := strings.ToLower(n.Data)
		if len(tag) == 2 && tag[0] == 'h' && tag[1] >= '1' && tag[1] <= '6' {
			innerText := strings.TrimSpace(extractText(n))
			if innerText == text {
				return n, int(tag[1] - '0')
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found, level := findHeadingByText(c, text); found != nil {
			return found, level
		}
	}
	return nil, 0
}

func isHeadingOfLevel(n *html.Node, maxLevel int) bool {
	if n.Type != html.ElementNode {
		return false
	}
	tag := strings.ToLower(n.Data)
	if len(tag) == 2 && tag[0] == 'h' && tag[1] >= '1' && tag[1] <= '6' {
		return int(tag[1]-'0') <= maxLevel
	}
	return false
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(extractText(c))
	}
	return b.String()
}

// --- Inline Comments ---------------------------------------------------------

// AddInlineComment adds an inline comment anchored to a text selection (DC v1).
func (f *Fetcher) AddInlineComment(ctx context.Context, pageID, body, textSelection string, matchCount, matchIndex int) (map[string]any, error) {
	lastFetchTime := fmt.Sprintf("%d", time.Now().UnixMilli())
	serializedHighlights, _ := json.Marshal([][]string{{textSelection}})

	data := map[string]any{
		"type": "comment",
		"container": map[string]any{
			"id":   pageID,
			"type": "page",
		},
		"body": map[string]any{
			"storage": map[string]any{
				"value":          body,
				"representation": "storage",
			},
		},
		"extensions": map[string]any{
			"location": "inline",
			"inlineProperties": map[string]any{
				"originalSelection":    textSelection,
				"numMatches":           matchCount,
				"matchIndex":           matchIndex,
				"lastFetchTime":        lastFetchTime,
				"serializedHighlights": string(serializedHighlights),
			},
		},
	}
	var raw map[string]any
	if err := f.client.Post(ctx, f.V1BaseURL()+"/content", data, &raw); err != nil {
		return nil, fmt.Errorf("confluence: add inline comment: %w", err)
	}
	return raw, nil
}

// --- Page History & Diff -----------------------------------------------------

// GetPageHistory returns a historical version of a page.
func (f *Fetcher) GetPageHistory(ctx context.Context, pageID string, version int, convertToMarkdown bool) (map[string]any, error) {
	path := fmt.Sprintf("%s/content/%s", f.V1BaseURL(), url.PathEscape(pageID))
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("status", "historical")
	q.Set("version", fmt.Sprintf("%d", version))
	q.Set("expand", "body.storage,version,space")
	req.URL.RawQuery = q.Encode()
	var raw map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get page history %s v%d: %w", pageID, version, err)
	}
	page := &confmodels.ConfluencePage{}
	if err := page.FromAPIResponse(raw); err != nil {
		return nil, fmt.Errorf("confluence: parse page history: %w", err)
	}
	out := page.ToSimplifiedDict()
	if convertToMarkdown {
		body, _ := out["body"].(map[string]any)
		if body != nil {
			storage, _ := body["storage"].(map[string]any)
			if storage != nil {
				if v, ok := storage["value"].(string); ok {
					storage["value"] = htmlconv.Convert(v)
					storage["representation"] = "markdown"
				}
			}
		}
	}
	return out, nil
}

// GetPageVersionDiff returns a unified diff between two page versions.
func (f *Fetcher) GetPageVersionDiff(ctx context.Context, pageID string, fromVersion, toVersion int, convertToMarkdown bool) (map[string]any, error) {
	fromPage, err := f.GetPageHistory(ctx, pageID, fromVersion, convertToMarkdown)
	if err != nil {
		return nil, fmt.Errorf("confluence: diff: from version %d: %w", fromVersion, err)
	}
	toPage, err := f.GetPageHistory(ctx, pageID, toVersion, convertToMarkdown)
	if err != nil {
		return nil, fmt.Errorf("confluence: diff: to version %d: %w", toVersion, err)
	}

	getBody := func(raw map[string]any) string {
		body, _ := raw["body"].(map[string]any)
		if body == nil {
			return ""
		}
		storage, _ := body["storage"].(map[string]any)
		if storage == nil {
			return ""
		}
		v, _ := storage["value"].(string)
		return v
	}

	fromText := getBody(fromPage)
	toText := getBody(toPage)
	title, _ := toPage["title"].(string)

	diff := unifiedDiff(fromText, toText, fmt.Sprintf("v%d", fromVersion), fmt.Sprintf("v%d", toVersion))

	return map[string]any{
		"page_id":      pageID,
		"title":        title,
		"from_version": fromVersion,
		"to_version":   toVersion,
		"diff":         diff,
	}, nil
}

func unifiedDiff(a, b, fromFile, toFile string) string {
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")

	// Simple LCS-based diff.
	type edit struct{ del, add int }
	lcs := lcsLength(aLines, bLines)

	var out strings.Builder
	out.WriteString(fmt.Sprintf("--- %s\n", fromFile))
	out.WriteString(fmt.Sprintf("+++ %s\n", toFile))

	i, j := 0, 0
	for idx := range lcs {
		// Skipped lines in a (deletions) and b (additions).
		for i < len(aLines) && (idx >= len(lcs) || lcs[idx].i != i) {
			out.WriteString(fmt.Sprintf("-%s\n", aLines[i]))
			i++
		}
		for j < len(bLines) && (idx >= len(lcs) || lcs[idx].j != j) {
			out.WriteString(fmt.Sprintf("+%s\n", bLines[j]))
			j++
		}
		if idx < len(lcs) && i < len(aLines) && j < len(bLines) {
			out.WriteString(fmt.Sprintf(" %s\n", aLines[i]))
			i++
			j++
		}
		_ = edit{} // suppress unused warning
	}
	for i < len(aLines) {
		out.WriteString(fmt.Sprintf("-%s\n", aLines[i]))
		i++
	}
	for j < len(bLines) {
		out.WriteString(fmt.Sprintf("+%s\n", bLines[j]))
		j++
	}
	return out.String()
}

type lcsPos struct{ i, j int }

func lcsLength(a, b []string) []lcsPos {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] > dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	// Backtrack to collect positions.
	var result []lcsPos
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			result = append([]lcsPos{{i - 1, j - 1}}, result...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}
	return result
}

// --- Attachment Upload -------------------------------------------------------

const maxAttachmentBytes = 50 * 1024 * 1024

// UploadAttachment uploads a file to Confluence content via multipart POST.
func (f *Fetcher) UploadAttachment(ctx context.Context, contentID, filePath, comment string, minorEdit bool) (map[string]any, error) {
	if contentID == "" {
		return map[string]any{"success": false, "error": "No content ID provided"}, nil
	}
	if filePath == "" {
		return map[string]any{"success": false, "error": "No file path provided"}, nil
	}

	if _, stErr := os.Stat(filePath); stErr != nil {
		return map[string]any{"success": false, "error": fmt.Sprintf("File not found: %s", filePath)}, nil
	}

	fileBytes, rdErr := os.ReadFile(filePath)
	if rdErr != nil {
		return map[string]any{"success": false, "error": fmt.Sprintf("Cannot read file: %v", rdErr)}, nil
	}

	filename := filepath.Base(filePath)
	fields := map[string]string{"minorEdit": "true"}
	if !minorEdit {
		fields["minorEdit"] = "false"
	}
	if comment != "" {
		fields["comment"] = comment
	}

	path := fmt.Sprintf("%s/content/%s/child/attachment", f.V1BaseURL(), url.PathEscape(contentID))
	var raw map[string]any
	if err := f.client.PostMultipart(ctx, path, fields, filename, fileBytes, &raw); err != nil {
		// On Server/DC, duplicate filename returns 400. Try update path.
		return map[string]any{"success": false, "error": fmt.Sprintf("Upload failed: %v (duplicate filename may require update attachment flow)", err)}, nil
	}

	return map[string]any{
		"success":    true,
		"content_id": contentID,
		"filename":   filename,
		"size":       len(fileBytes),
		"id":         raw["id"],
	}, nil
}

// UploadAttachments uploads multiple files to Confluence content.
func (f *Fetcher) UploadAttachments(ctx context.Context, contentID string, filePaths []string, comment string, minorEdit bool) (map[string]any, error) {
	uploaded := []map[string]any{}
	failed := []map[string]any{}

	for _, fp := range filePaths {
		result, _ := f.UploadAttachment(ctx, contentID, fp, comment, minorEdit)
		if ok, _ := result["success"].(bool); ok {
			uploaded = append(uploaded, map[string]any{
				"filename": result["filename"],
				"size":     result["size"],
				"id":       result["id"],
			})
		} else {
			failed = append(failed, map[string]any{
				"filename": filepath.Base(fp),
				"error":    result["error"],
			})
		}
	}

	return map[string]any{
		"success":    true,
		"content_id": contentID,
		"total":      len(filePaths),
		"uploaded":   uploaded,
		"failed":     failed,
	}, nil
}

// UploadAttachmentBase64 uploads a file from base64-encoded content.
func (f *Fetcher) UploadAttachmentBase64(ctx context.Context, contentID, filename, b64Content, comment string, minorEdit bool) (map[string]any, error) {
	if contentID == "" {
		return map[string]any{"success": false, "error": "No content ID provided"}, nil
	}
	if filename == "" {
		return map[string]any{"success": false, "error": "No filename provided"}, nil
	}
	data, err := base64.StdEncoding.DecodeString(b64Content)
	if err != nil {
		return map[string]any{"success": false, "error": fmt.Sprintf("Invalid base64: %v", err)}, nil
	}
	if len(data) > maxAttachmentBytes {
		return map[string]any{
			"success": false, "error": fmt.Sprintf("Content is %d bytes which exceeds the 50 MB inline limit", len(data)),
		}, nil
	}

	fields := map[string]string{"minorEdit": "true"}
	if !minorEdit {
		fields["minorEdit"] = "false"
	}
	if comment != "" {
		fields["comment"] = comment
	}

	path := fmt.Sprintf("%s/content/%s/child/attachment", f.V1BaseURL(), url.PathEscape(contentID))
	var raw map[string]any
	if err := f.client.PostMultipart(ctx, path, fields, filename, data, &raw); err != nil {
		return map[string]any{"success": false, "error": fmt.Sprintf("Upload failed: %v", err)}, nil
	}

	return map[string]any{
		"success":    true,
		"content_id": contentID,
		"filename":   filename,
		"size":       len(data),
		"id":         raw["id"],
	}, nil
}

// --- Attachment Download -----------------------------------------------------

// FetchAttachmentMetadata returns metadata for an attachment by ID.
func (f *Fetcher) FetchAttachmentMetadata(ctx context.Context, attachmentID string) (map[string]any, error) {
	var raw map[string]any
	if err := f.client.Get(ctx, f.V1BaseURL()+"/content/"+url.PathEscape(attachmentID), &raw); err != nil {
		return nil, fmt.Errorf("confluence: attachment metadata: %w", err)
	}
	return raw, nil
}

// ResolveAttachmentDownloadURL resolves the download URL for an attachment (DC).
func (f *Fetcher) ResolveAttachmentDownloadURL(downloadURL string) string {
	if downloadURL == "" {
		return ""
	}
	// DC: use the legacy download link as-is.
	if !strings.HasPrefix(downloadURL, "http") {
		downloadURL = f.cfg.URL + downloadURL
	}
	return downloadURL
}

// DownloadAttachmentBase64 downloads an attachment and returns base64-encoded data.
func (f *Fetcher) DownloadAttachmentBase64(ctx context.Context, attachmentID string) (map[string]any, error) {
	// Get metadata.
	meta, err := f.FetchAttachmentMetadata(ctx, attachmentID)
	if err != nil {
		return nil, err
	}

	links, _ := meta["_links"].(map[string]any)
	if links == nil {
		return nil, fmt.Errorf("confluence: no _links in attachment metadata for %s", attachmentID)
	}
	dlURL, _ := links["download"].(string)
	if dlURL == "" {
		return nil, fmt.Errorf("confluence: no download URL for attachment %s", attachmentID)
	}

	dlURL = f.ResolveAttachmentDownloadURL(dlURL)

	filename, _ := meta["title"].(string)
	if filename == "" {
		filename = attachmentID
	}

	ext, _ := meta["extensions"].(map[string]any)
	mimeType := ""
	fileSize := 0
	if ext != nil {
		mimeType, _ = ext["mediaType"].(string)
		if fs, ok := ext["fileSize"].(float64); ok {
			fileSize = int(fs)
		}
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if fileSize > maxAttachmentBytes {
		return map[string]any{
			"success": false, "attachment_id": attachmentID,
			"filename": filename, "file_size": fileSize,
			"error": fmt.Sprintf("Attachment '%s' is %d bytes which exceeds the 50 MB inline limit.", filename, fileSize),
		}, nil
	}

	data, derr := f.client.GetRaw(ctx, dlURL)
	if derr != nil {
		return nil, fmt.Errorf("confluence: download attachment: %w", derr)
	}

	if len(data) > maxAttachmentBytes {
		return map[string]any{
			"success": false, "attachment_id": attachmentID,
			"filename": filename, "file_size": len(data),
			"error": fmt.Sprintf("Attachment '%s' is %d bytes which exceeds the 50 MB inline limit.", filename, len(data)),
		}, nil
	}

	return map[string]any{
		"success":       true,
		"attachment_id": attachmentID,
		"filename":      filename,
		"mime_type":     mimeType,
		"size":          len(data),
		"data":          data,
		"data_base64":   base64.StdEncoding.EncodeToString(data),
	}, nil
}

// DownloadContentAttachmentsBase64 downloads all attachments for a content item.
func (f *Fetcher) DownloadContentAttachmentsBase64(ctx context.Context, contentID string) (map[string]any, error) {
	attResult, err := f.GetContentAttachments(ctx, contentID, 0, 100, "", "")
	if err != nil {
		return nil, err
	}
	attachments, _ := attResult["attachments"].([]map[string]any)
	if attachments == nil {
		return map[string]any{
			"success":    true,
			"content_id": contentID,
			"message":    "No attachments found",
			"downloaded": 0,
			"failed":     []any{},
		}, nil
	}

	fetched := []map[string]any{}
	failed := []map[string]any{}

	for _, att := range attachments {
		attID := fmt.Sprintf("%v", att["id"])
		result, ferr := f.DownloadAttachmentBase64(ctx, attID)
		if ferr != nil {
			filename, _ := att["title"].(string)
			failed = append(failed, map[string]any{"filename": filename, "error": ferr.Error()})
			continue
		}
		if ok, _ := result["success"].(bool); !ok {
			errMsg, _ := result["error"].(string)
			filename, _ := result["filename"].(string)
			failed = append(failed, map[string]any{"filename": filename, "error": errMsg})
			continue
		}
		fetched = append(fetched, result)
	}

	return map[string]any{
		"success":     true,
		"content_id":  contentID,
		"total":       len(attachments),
		"downloaded":  len(fetched),
		"attachments": fetched,
		"failed":      failed,
	}, nil
}

// --- Page Images --------------------------------------------------------------

// GetPageImages returns all image attachments for a page, with base64 data.
func (f *Fetcher) GetPageImages(ctx context.Context, contentID string) (map[string]any, error) {
	attResult, err := f.GetContentAttachments(ctx, contentID, 0, 100, "", "")
	if err != nil {
		return nil, err
	}
	attachments, _ := attResult["attachments"].([]map[string]any)

	fetched := []map[string]any{}
	failed := []map[string]any{}
	imageCount := 0

	for _, att := range attachments {
		ext, _ := att["extensions"].(map[string]any)
		mimeType := ""
		if ext != nil {
			mimeType, _ = ext["mediaType"].(string)
		}
		filename, _ := att["title"].(string)
		if isImg, _ := utils.IsImageAttachment(mimeType, filename); !isImg {
			continue
		}
		imageCount++

		attID := fmt.Sprintf("%v", att["id"])
		result, ferr := f.DownloadAttachmentBase64(ctx, attID)
		if ferr != nil {
			failed = append(failed, map[string]any{"filename": filename, "error": ferr.Error()})
			continue
		}
		if ok, _ := result["success"].(bool); !ok {
			errMsg, _ := result["error"].(string)
			failed = append(failed, map[string]any{"filename": filename, "error": errMsg})
			continue
		}
		fetched = append(fetched, result)
	}

	return map[string]any{
		"success":    true,
		"content_id": contentID,
		"total":      imageCount,
		"images":     len(fetched),
		"fetched":    fetched,
		"failed":     failed,
	}, nil
}

// --- Page Restrictions -------------------------------------------------------

// GetPageRestrictions returns view and edit restrictions for a page.
func (f *Fetcher) GetPageRestrictions(ctx context.Context, pageID string) (map[string]any, error) {
	// Docs default expand includes read/update restrictions for both user and group.
	q := url.Values{}
	q.Set("expand", "update.restrictions.user,read.restrictions.group,read.restrictions.user,update.restrictions.group")
	var raw map[string]any
	if err := f.getJSON(ctx, f.contentPath(pageID, "restriction", "byOperation"), q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get restrictions: %w", err)
	}

	result := map[string]any{
		"read":   map[string]any{"users": []any{}, "groups": []any{}},
		"update": map[string]any{"users": []any{}, "groups": []any{}},
	}

	for _, opKey := range []string{"read", "update"} {
		opData, _ := raw[opKey].(map[string]any)
		if opData == nil {
			continue
		}
		restrictions, _ := opData["restrictions"].(map[string]any)
		if restrictions == nil {
			continue
		}
		users, _ := restrictions["user"].(map[string]any)
		if users != nil {
			for _, u := range extractStringList(users, "results") {
				for _, idKey := range []string{"accountId", "username", "name"} {
					if id, _ := u[idKey].(string); id != "" {
						result[opKey].(map[string]any)["users"] = append(
							result[opKey].(map[string]any)["users"].([]any), id)
						break
					}
				}
			}
		}
		groups, _ := restrictions["group"].(map[string]any)
		if groups != nil {
			for _, g := range extractStringList(groups, "results") {
				if name, _ := g["name"].(string); name != "" {
					result[opKey].(map[string]any)["groups"] = append(
						result[opKey].(map[string]any)["groups"].([]any), name)
				}
			}
		}
	}
	return result, nil
}

func extractStringList(m map[string]any, key string) []map[string]any {
	raw, _ := m[key].([]any)
	if raw == nil {
		return nil
	}
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if im, ok := item.(map[string]any); ok {
			out = append(out, im)
		}
	}
	return out
}

// SetPageRestrictions sets view and edit restrictions on a page (DC uses username).
func (f *Fetcher) SetPageRestrictions(ctx context.Context, pageID string, readUsers, readGroups, editUsers, editGroups []string) (map[string]any, error) {
	if readUsers == nil {
		readUsers = []string{}
	}
	if readGroups == nil {
		readGroups = []string{}
	}
	if editUsers == nil {
		editUsers = []string{}
	}
	if editGroups == nil {
		editGroups = []string{}
	}

	buildUserEntry := func(id string) map[string]any {
		return map[string]any{"type": "known", "username": id}
	}
	buildGroupEntry := func(name string) map[string]any {
		return map[string]any{"type": "group", "name": name}
	}

	readUserObjs := make([]map[string]any, len(readUsers))
	for i, u := range readUsers {
		readUserObjs[i] = buildUserEntry(u)
	}
	readGroupObjs := make([]map[string]any, len(readGroups))
	for i, g := range readGroups {
		readGroupObjs[i] = buildGroupEntry(g)
	}
	editUserObjs := make([]map[string]any, len(editUsers))
	for i, u := range editUsers {
		editUserObjs[i] = buildUserEntry(u)
	}
	editGroupObjs := make([]map[string]any, len(editGroups))
	for i, g := range editGroups {
		editGroupObjs[i] = buildGroupEntry(g)
	}

	payload := []map[string]any{
		{
			"operation": "read",
			"restrictions": map[string]any{
				"user":  readUserObjs,
				"group": readGroupObjs,
			},
		},
		{
			"operation": "update",
			"restrictions": map[string]any{
				"user":  editUserObjs,
				"group": editGroupObjs,
			},
		},
	}

	path := fmt.Sprintf("%s/content/%s/restriction", f.V1BaseURL(), url.PathEscape(pageID))
	var raw map[string]any
	if err := f.client.Put(ctx, path, payload, &raw); err != nil {
		return nil, fmt.Errorf("confluence: set restrictions: %w", err)
	}

	return map[string]any{
		"read":   map[string]any{"users": readUsers, "groups": readGroups},
		"update": map[string]any{"users": editUsers, "groups": editGroups},
	}, nil
}

// --- Copy Page ---------------------------------------------------------------

// CopyPage creates a copy of a page. On DC, fetches the source page body and
// creates a new page manually (no native copy endpoint).
func (f *Fetcher) CopyPage(ctx context.Context, sourcePageID, destSpaceKey, newTitle, destParentID string, copyAttachments bool) (map[string]any, error) {
	// DC path: GET source page body, then POST create.
	source, err := f.GetPage(ctx, sourcePageID, []string{"body.storage", "version", "space"})
	if err != nil {
		return nil, fmt.Errorf("confluence: copy page: get source: %w", err)
	}

	storageVal, _ := source.Body["storage"].(map[string]any)
	bodyValue := ""
	if storageVal != nil {
		bodyValue, _ = storageVal["value"].(string)
	}

	createBody := map[string]any{
		"type":  "page",
		"title": newTitle,
		"space": map[string]any{"key": destSpaceKey},
		"body": map[string]any{
			"storage": map[string]any{
				"value":          bodyValue,
				"representation": "storage",
			},
		},
	}
	if destParentID != "" {
		createBody["ancestors"] = []map[string]any{{"id": destParentID}}
	}

	var raw map[string]any
	if err := f.client.Post(ctx, f.V1BaseURL()+"/content", createBody, &raw); err != nil {
		return nil, fmt.Errorf("confluence: copy page: create: %w", err)
	}
	return raw, nil
}

// GetPageProperties returns content properties for a page.
func (f *Fetcher) GetPageProperties(ctx context.Context, pageID string) ([]map[string]any, error) {
	path := fmt.Sprintf("%s/content/%s/property", f.V1BaseURL(), url.PathEscape(pageID))
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("limit", "50")
	req.URL.RawQuery = q.Encode()
	var raw struct {
		Results []map[string]any `json:"results"`
	}
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get page properties: %w", err)
	}
	return raw.Results, nil
}

// extractEmojiFromProperty extracts an emoji character from a Confluence
// content property value. Mirrors Python's extract_emoji_from_property().
func extractEmojiFromProperty(value any) string {
	switch v := value.(type) {
	case map[string]any:
		if fb, ok := v["fallback"].(string); ok && fb != "" {
			return fb
		}
		if sn, ok := v["shortName"].(string); ok && sn != "" {
			return sn
		}
		if id, ok := v["id"].(string); ok && id != "" {
			if cp, err := strconv.ParseUint(id, 16, 32); err == nil {
				return string(rune(cp))
			}
		}
	case string:
		return v
	}
	return ""
}

// extractPageWidth extracts the page layout width from a content property value.
func extractPageWidth(value any) string {
	switch v := value.(type) {
	case map[string]any:
		if s, ok := v["value"].(string); ok {
			return s
		}
	case string:
		return v
	}
	return ""
}

// --- New read endpoints ------------------------------------------------------

// ScanContent scans for content matching a CQL query (DC /content/scan).
func (f *Fetcher) ScanContent(ctx context.Context, cql string, start, limit int, expand []string) (map[string]any, error) {
	q := url.Values{}
	if cql != "" {
		q.Set("cql", cql)
	}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw struct {
		Results []map[string]any `json:"results"`
		Size    int              `json:"size"`
		Start   int              `json:"start"`
		Limit   int              `json:"limit"`
		Links   map[string]any   `json:"_links"`
	}
	if err := f.getJSON(ctx, f.V1BaseURL()+"/content/scan", q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: scan content: %w", err)
	}
	return map[string]any{
		"success": true, "results": raw.Results, "total": raw.Size,
		"start": raw.Start, "limit": raw.Limit, "has_more": raw.Links["next"] != nil,
	}, nil
}

// GetContentHistory returns the history object for a content item (/content/{id}/history).
func (f *Fetcher) GetContentHistory(ctx context.Context, contentID string, expand []string) (map[string]any, error) {
	q := url.Values{}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw map[string]any
	if err := f.getJSON(ctx, f.contentPath(contentID, "history"), q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get content history %s: %w", contentID, err)
	}
	return raw, nil
}

// GetRestrictionForOperation returns restrictions for a single operation key (read/update).
func (f *Fetcher) GetRestrictionForOperation(ctx context.Context, contentID, operationKey string) (map[string]any, error) {
	q := url.Values{}
	q.Set("expand", "update.restrictions.user,read.restrictions.group,read.restrictions.user,update.restrictions.group")
	var raw map[string]any
	path := f.contentPath(contentID, "restriction", "byOperation", operationKey)
	if err := f.getJSON(ctx, path, q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get restriction %s/%s: %w", contentID, operationKey, err)
	}
	return raw, nil
}

// GetContentChildren returns child items of a given type (page, comment, attachment, ...).
func (f *Fetcher) GetContentChildren(ctx context.Context, contentID, childType string, start, limit int, expand []string) ([]map[string]any, error) {
	q := url.Values{}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw struct {
		Results []map[string]any `json:"results"`
	}
	path := f.contentPath(contentID, "child", childType)
	if err := f.getJSON(ctx, path, q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get children %s/%s: %w", contentID, childType, err)
	}
	return raw.Results, nil
}

// GetDescendants returns all descendants of a content item.
func (f *Fetcher) GetDescendants(ctx context.Context, contentID string, start, limit int, expand []string) (map[string]any, error) {
	q := url.Values{}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw struct {
		Results []map[string]any `json:"results"`
		Size    int              `json:"size"`
	}
	if err := f.getJSON(ctx, f.contentPath(contentID, "descendant"), q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get descendants %s: %w", contentID, err)
	}
	return map[string]any{"success": true, "results": raw.Results, "total": raw.Size}, nil
}

// GetDescendantsOfType returns descendants of a specific type.
func (f *Fetcher) GetDescendantsOfType(ctx context.Context, contentID, descendantType string, start, limit int, expand []string) (map[string]any, error) {
	q := url.Values{}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw struct {
		Results []map[string]any `json:"results"`
		Size    int              `json:"size"`
	}
	path := f.contentPath(contentID, "descendant", descendantType)
	if err := f.getJSON(ctx, path, q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get descendants %s/%s: %w", contentID, descendantType, err)
	}
	return map[string]any{"success": true, "results": raw.Results, "total": raw.Size}, nil
}

// GetContentProperty returns a single content property by key.
func (f *Fetcher) GetContentProperty(ctx context.Context, contentID, key string) (map[string]any, error) {
	var raw map[string]any
	path := f.contentPath(contentID, "property", key)
	if err := f.getJSON(ctx, path, nil, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get property %s/%s: %w", contentID, key, err)
	}
	return raw, nil
}

// SearchEntities performs a global search across Confluence entities (/search).
func (f *Fetcher) SearchEntities(ctx context.Context, cql, excerpt, expand string, start, limit int, includeArchivedSpaces bool) (map[string]any, error) {
	q := url.Values{}
	if cql != "" {
		q.Set("cql", cql)
	}
	if excerpt != "" {
		q.Set("excerpt", excerpt)
	}
	if expand != "" {
		q.Set("expand", expand)
	}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if includeArchivedSpaces {
		q.Set("includeArchivedSpaces", "true")
	}
	var raw struct {
		Results []map[string]any `json:"results"`
		Size    int              `json:"size"`
		Start   int              `json:"start"`
		Limit   int              `json:"limit"`
	}
	if err := f.getJSON(ctx, f.V1BaseURL()+"/search", q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: search entities: %w", err)
	}
	return map[string]any{
		"success": true, "results": raw.Results, "total": raw.Size,
		"start": raw.Start, "limit": raw.Limit,
	}, nil
}

// GetSpace returns a space by its key.
func (f *Fetcher) GetSpace(ctx context.Context, spaceKey string, expand []string) (map[string]any, error) {
	q := url.Values{}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw map[string]any
	path := f.V1BaseURL() + "/space/" + url.PathEscape(spaceKey)
	if err := f.getJSON(ctx, path, q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get space %s: %w", spaceKey, err)
	}
	return raw, nil
}

// GetSpaceContent returns content items in a space filtered by type.
func (f *Fetcher) GetSpaceContent(ctx context.Context, spaceKey, contentType string, depth int, expand []string, start, limit int) (map[string]any, error) {
	q := url.Values{}
	if depth > 0 {
		q.Set("depth", fmt.Sprintf("%d", depth))
	}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	var raw struct {
		Results []map[string]any `json:"results"`
		Size    int              `json:"size"`
	}
	path := f.V1BaseURL() + "/space/" + url.PathEscape(spaceKey) + "/content/" + url.PathEscape(contentType)
	if err := f.getJSON(ctx, path, q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get space content %s/%s: %w", spaceKey, contentType, err)
	}
	return map[string]any{"success": true, "results": raw.Results, "total": raw.Size}, nil
}

// GetSpaceProperties returns properties for a space.
func (f *Fetcher) GetSpaceProperties(ctx context.Context, spaceKey string, expand []string, start, limit int) ([]map[string]any, error) {
	q := url.Values{}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	if start > 0 {
		q.Set("start", fmt.Sprintf("%d", start))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	var raw struct {
		Results []map[string]any `json:"results"`
	}
	path := f.V1BaseURL() + "/space/" + url.PathEscape(spaceKey) + "/property"
	if err := f.getJSON(ctx, path, q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get space properties %s: %w", spaceKey, err)
	}
	return raw.Results, nil
}

// GetSpaceProperty returns a single space property by key.
func (f *Fetcher) GetSpaceProperty(ctx context.Context, spaceKey, key string, expand []string) (map[string]any, error) {
	q := url.Values{}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw map[string]any
	path := f.V1BaseURL() + "/space/" + url.PathEscape(spaceKey) + "/property/" + url.PathEscape(key)
	if err := f.getJSON(ctx, path, q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get space property %s/%s: %w", spaceKey, key, err)
	}
	return raw, nil
}

// GetCurrentUser returns the authenticated user's profile.
func (f *Fetcher) GetCurrentUser(ctx context.Context, expand []string) (map[string]any, error) {
	q := url.Values{}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw map[string]any
	if err := f.getJSON(ctx, f.V1BaseURL()+"/user/current", q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get current user: %w", err)
	}
	return raw, nil
}

// GetUser returns a user by key or username.
func (f *Fetcher) GetUser(ctx context.Context, keyOrUsername string, expand []string) (map[string]any, error) {
	q := url.Values{}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	}
	var raw map[string]any
	path := f.V1BaseURL() + "/user/" + url.PathEscape(keyOrUsername)
	if err := f.getJSON(ctx, path, q, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get user %s: %w", keyOrUsername, err)
	}
	return raw, nil
}

// --- New write endpoints -----------------------------------------------------

// RemovePageLabelByName removes a label from content by its name (?name=).
func (f *Fetcher) RemovePageLabelByName(ctx context.Context, contentID, name string) error {
	q := url.Values{}
	q.Set("name", name)
	path := f.contentPath(contentID, "label")
	req, err := f.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = q.Encode()
	if err := f.client.Do(req, nil); err != nil {
		return fmt.Errorf("confluence: remove label %s/%s: %w", contentID, name, err)
	}
	return nil
}

// RemovePageLabel removes a label from content by its full label path ({prefix}/{name}).
func (f *Fetcher) RemovePageLabel(ctx context.Context, contentID, label string) error {
	path := f.contentPath(contentID, "label", label)
	if err := f.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("confluence: remove label %s/%s: %w", contentID, label, err)
	}
	return nil
}

// CreateContentProperty creates a new content property.
func (f *Fetcher) CreateContentProperty(ctx context.Context, contentID, key string, value any) (map[string]any, error) {
	body := map[string]any{"key": key, "value": value}
	var raw map[string]any
	path := f.contentPath(contentID, "property")
	if err := f.client.Post(ctx, path, body, &raw); err != nil {
		return nil, fmt.Errorf("confluence: create property %s/%s: %w", contentID, key, err)
	}
	return raw, nil
}

// UpdateContentProperty updates an existing content property (requires id and version).
func (f *Fetcher) UpdateContentProperty(ctx context.Context, contentID, key string, id int, version int, value any) (map[string]any, error) {
	body := map[string]any{
		"id":      id,
		"version": map[string]any{"number": version},
		"value":   value,
	}
	var raw map[string]any
	path := f.contentPath(contentID, "property", key)
	if err := f.client.Put(ctx, path, body, &raw); err != nil {
		return nil, fmt.Errorf("confluence: update property %s/%s: %w", contentID, key, err)
	}
	return raw, nil
}

// DeleteContentProperty deletes a content property by key.
func (f *Fetcher) DeleteContentProperty(ctx context.Context, contentID, key string) error {
	path := f.contentPath(contentID, "property", key)
	if err := f.client.Delete(ctx, path); err != nil {
		return fmt.Errorf("confluence: delete property %s/%s: %w", contentID, key, err)
	}
	return nil
}

// UpdateAttachmentMetadata updates metadata (title, comment) for an attachment.
func (f *Fetcher) UpdateAttachmentMetadata(ctx context.Context, attachmentID, title, comment string) (map[string]any, error) {
	body := map[string]any{}
	if title != "" {
		body["title"] = title
	}
	if comment != "" {
		body["comment"] = comment
	}
	var raw map[string]any
	path := f.contentPath(attachmentID)
	if err := f.client.Put(ctx, path, body, &raw); err != nil {
		return nil, fmt.Errorf("confluence: update attachment metadata %s: %w", attachmentID, err)
	}
	return raw, nil
}

// UpdateAttachmentData replaces the data of an attachment via multipart upload.
func (f *Fetcher) UpdateAttachmentData(ctx context.Context, attachmentID, filename string, fileBytes []byte, comment string) (map[string]any, error) {
	fields := map[string]string{}
	if comment != "" {
		fields["comment"] = comment
	}
	path := f.contentPath(attachmentID, "data")
	var raw map[string]any
	if err := f.client.PostMultipart(ctx, path, fields, filename, fileBytes, &raw); err != nil {
		return nil, fmt.Errorf("confluence: update attachment data %s: %w", attachmentID, err)
	}
	return raw, nil
}
