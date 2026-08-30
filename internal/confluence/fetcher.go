package confluence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/atlassian"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/htmlconv"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models"
	confmodels "gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/confluence"
)

type Fetcher struct {
	client *atlassian.Client
	cfg    *config.ConfluenceConfig
}

func New(cfg *config.ConfluenceConfig, c *atlassian.Client) *Fetcher {
	return &Fetcher{client: c, cfg: cfg}
}
func (f *Fetcher) Client() *atlassian.Client { return f.client }
func (f *Fetcher) V1BaseURL() string         { return f.cfg.BaseURL() + "/rest/api" }

func (f *Fetcher) GetPage(ctx context.Context, pageID string, expand []string) (*confmodels.ConfluencePage, error) {
	if pageID == "" {
		return nil, fmt.Errorf("confluence: page ID required")
	}
	path := fmt.Sprintf("%s/content/%s", f.V1BaseURL(), url.PathEscape(pageID))
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if len(expand) > 0 {
		q := req.URL.Query()
		q.Set("expand", strings.Join(expand, ","))
		req.URL.RawQuery = q.Encode()
	}
	var raw map[string]any
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("confluence: get page %s: %w", pageID, err)
	}
	page := &confmodels.ConfluencePage{}
	if err := page.FromAPIResponse(raw); err != nil {
		return nil, fmt.Errorf("confluence: parse page: %w", err)
	}
	return page, nil
}

// GetPageByTitle finds a page by its title and space key.
func (f *Fetcher) GetPageByTitle(ctx context.Context, spaceKey, title string, expand []string) (*confmodels.ConfluencePage, error) {
	title = strings.ReplaceAll(title, `"`, `\"`)
	cql := fmt.Sprintf(`space = "%s" AND title = "%s"`, spaceKey, title)
	results, err := f.Search(ctx, cql, 2)
	if err != nil {
		return nil, fmt.Errorf("confluence: search page by title: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("confluence: page with title %q not found in space %q", title, spaceKey)
	}
	pageID := models.AsID(models.MapAny(results[0], "id"))
	if pageID == "" {
		return nil, fmt.Errorf("confluence: page found but has no ID")
	}
	return f.GetPage(ctx, pageID, expand)
}

// GetPageContent fetches a page by ID with full metadata and processed content.
// Mirrors Python's get_page_content().
func (f *Fetcher) GetPageContent(ctx context.Context, pageID string, convertToMarkdown bool) (*confmodels.ConfluencePage, error) {
	expand := []string{"body.storage", "version", "space", "children.attachment", "history"}
	page, err := f.GetPage(ctx, pageID, expand)
	if err != nil {
		return nil, err
	}
	f.processPageContent(page, convertToMarkdown)
	return page, nil
}

// GetPageByTitleContent finds a page by title and returns it with processed content.
// Mirrors Python's get_page_by_title().
func (f *Fetcher) GetPageByTitleContent(ctx context.Context, spaceKey, title string, convertToMarkdown bool) (*confmodels.ConfluencePage, error) {
	expand := []string{"body.storage", "version", "space", "children.attachment", "history"}
	page, err := f.GetPageByTitle(ctx, spaceKey, title, expand)
	if err != nil {
		return nil, err
	}
	f.processPageContent(page, convertToMarkdown)
	return page, nil
}

// processPageContent extracts body storage value, sets Content/ContentFormat,
// fetches emoji and page width from content properties, and builds the URL.
func (f *Fetcher) processPageContent(page *confmodels.ConfluencePage, convertToMarkdown bool) {
	// Extract storage HTML
	var storageHTML string
	if page.Body != nil {
		if storage, ok := page.Body["storage"].(map[string]any); ok {
			storageHTML, _ = storage["value"].(string)
		}
	}

	// Convert storage HTML to markdown when requested.
	if convertToMarkdown {
		page.Content = htmlconv.Convert(storageHTML)
		page.ContentFormat = "markdown"
	} else {
		page.Content = storageHTML
		page.ContentFormat = "storage"
	}

	// Fetch emoji and page width from content properties
	page.Emoji = f.getPageEmoji(page.ID)
	page.PageWidth = f.getPageWidth(page.ID)

	// Build URL (DC format)
	spaceKey := ""
	if page.Space != nil {
		spaceKey, _ = page.Space["key"].(string)
	}
	baseURL := f.cfg.BaseURL()
	if spaceKey != "" {
		page.URL = fmt.Sprintf("%s/spaces/%s/pages/%s", baseURL, spaceKey, page.ID)
	} else {
		page.URL = fmt.Sprintf("%s/pages/viewpage.action?pageId=%s", baseURL, page.ID)
	}
}

// getPageEmoji fetches the page title emoji from content properties.
func (f *Fetcher) getPageEmoji(pageID string) string {
	props, err := f.GetPageProperties(context.Background(), pageID)
	if err != nil {
		return ""
	}
	for _, prop := range props {
		key, _ := prop["key"].(string)
		if key == "emoji-title-published" || key == "emoji-title-draft" {
			return extractEmojiFromProperty(prop["value"])
		}
	}
	return ""
}

// getPageWidth fetches the page layout width from content properties.
func (f *Fetcher) getPageWidth(pageID string) string {
	props, err := f.GetPageProperties(context.Background(), pageID)
	if err != nil {
		return ""
	}
	for _, prop := range props {
		key, _ := prop["key"].(string)
		if key == "content-appearance-published" || key == "content-appearance-draft" {
			return extractPageWidth(prop["value"])
		}
	}
	return ""
}

func (f *Fetcher) Search(ctx context.Context, cql string, limit int) ([]map[string]any, error) {
	path := f.V1BaseURL() + "/content/search"
	req, err := f.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("cql", cql)
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	q.Set("expand", "content.history,content.version")
	req.URL.RawQuery = q.Encode()
	var raw struct {
		Results []map[string]any `json:"results"`
		Size    int              `json:"size"`
	}
	if err := f.client.Do(req, &raw); err != nil {
		return nil, fmt.Errorf("confluence: search: %w", err)
	}
	return raw.Results, nil
}

// --- Shared helpers ----------------------------------------------------------

// doJSON performs a request with query params and optional JSON body, decoding the response.
func (f *Fetcher) doJSON(ctx context.Context, method, path string, query url.Values, body, out any) error {
	req, err := f.client.NewRequest(ctx, method, path, nil)
	if err != nil {
		return err
	}
	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}
	if body != nil {
		b := mustMarshal(body)
		req.Body = io.NopCloser(bytes.NewReader(b))
		req.ContentLength = int64(len(b))
		req.Header.Set("Content-Type", "application/json")
	}
	return f.client.Do(req, out)
}

func mustMarshal(v any) []byte { b, _ := json.Marshal(v); return b }

// getJSON is a convenience wrapper for GET requests without a body.
func (f *Fetcher) getJSON(ctx context.Context, path string, query url.Values, out any) error {
	return f.doJSON(ctx, http.MethodGet, path, query, nil, out)
}

// contentPath builds a /rest/api/content/{id}[/suffix...] path with escaped segments.
func (f *Fetcher) contentPath(contentID string, suffix ...string) string {
	p := fmt.Sprintf("%s/content/%s", f.V1BaseURL(), url.PathEscape(contentID))
	for _, s := range suffix {
		p += "/" + url.PathEscape(s)
	}
	return p
}
