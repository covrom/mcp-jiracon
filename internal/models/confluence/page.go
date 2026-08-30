package confluence

import (
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/models/constants"
)

// ConfluencePage represents a Confluence page or blog post.
type ConfluencePage struct {
	ID            string           `json:"id"`
	Type          string           `json:"type"`
	Status        string           `json:"status,omitempty"`
	Title         string           `json:"title"`
	Space         map[string]any   `json:"space,omitempty"`
	Body          map[string]any   `json:"body,omitempty"`
	Version       map[string]any   `json:"version,omitempty"`
	History       map[string]any   `json:"history,omitempty"`
	Subtype       string           `json:"subtype,omitempty"`
	Content       string           `json:"content,omitempty"`
	ContentFormat string           `json:"content_format,omitempty"`
	Created       string           `json:"created,omitempty"`
	Updated       string           `json:"updated,omitempty"`
	Author        string           `json:"author,omitempty"`
	URL           string           `json:"url,omitempty"`
	Emoji         string           `json:"emoji,omitempty"`
	PageWidth     string           `json:"page_width,omitempty"`
	Ancestors     []map[string]any `json:"ancestors,omitempty"`
	Attachments   []map[string]any `json:"attachments,omitempty"`
}

// FromAPIResponse populates from a Confluence content dict.
func (p *ConfluencePage) FromAPIResponse(data map[string]any) error {
	p.ID = models.AsID(models.MapAny(data, "id"))
	p.Type = models.MapString(data, "type")
	if p.Type == "" {
		p.Type = "page"
	}
	p.Status = models.MapString(data, "status")
	if p.Status == "" {
		p.Status = "current"
	}
	p.Title = models.MapString(data, "title")
	p.Subtype = models.MapString(data, "subtype")
	p.Space = models.MapMap(data, "space")
	p.Body = models.MapMap(data, "body")
	p.Version = models.MapMap(data, "version")
	p.History = models.MapMap(data, "history")

	// Extract ancestors
	if anc := models.MapSlice(data, "ancestors"); anc != nil {
		for _, a := range anc {
			if am, ok := a.(map[string]any); ok {
				p.Ancestors = append(p.Ancestors, am)
			}
		}
	}

	// Extract attachments from children.attachment.results
	if children := models.MapMap(data, "children"); children != nil {
		if att := models.MapMap(children, "attachment"); att != nil {
			if results := models.MapSlice(att, "results"); results != nil {
				for _, r := range results {
					if rm, ok := r.(map[string]any); ok {
						p.Attachments = append(p.Attachments, rm)
					}
				}
			}
		}
	}

	// Extract created/updated from history
	if hist := models.MapMap(data, "history"); hist != nil {
		p.Created = models.MapString(hist, "createdDate")
		lastUpdated := models.MapMap(hist, "lastUpdated")
		if lastUpdated != nil {
			p.Updated = models.MapString(lastUpdated, "when")
		}
		// Fall back to version date if no history updated
		if p.Updated == "" && p.Version != nil {
			p.Updated = models.MapString(p.Version, "when")
		}
		// Fall back to history.createdBy if no top-level author
		if createdBy := models.MapMap(hist, "createdBy"); createdBy != nil {
			p.Author = models.MapString(createdBy, "displayName")
		}
	}

	return nil
}

// ToSimplifiedDict returns a JSON-friendly map matching Python's format.
func (p *ConfluencePage) ToSimplifiedDict() map[string]any {
	out := map[string]any{
		"id":      p.ID,
		"title":   p.Title,
		"type":    p.Type,
		"created": models.FormatTimestamp(p.Created),
		"updated": models.FormatTimestamp(p.Updated),
		"url":     p.URL,
	}
	if p.Subtype != "" {
		out["subtype"] = p.Subtype
	}
	if p.Space != nil {
		spaceKey := models.MapString(p.Space, "key")
		spaceName := models.MapString(p.Space, "name")
		if spaceName == "" {
			spaceName = constants.Unknown
		}
		out["space"] = map[string]any{"key": spaceKey, "name": spaceName}
	}
	if p.Author != "" {
		out["author"] = p.Author
	}
	if p.Version != nil {
		out["version"] = models.MapInt(p.Version, "number")
	}
	if len(p.Attachments) > 0 {
		atts := make([]map[string]any, 0, len(p.Attachments))
		for _, a := range p.Attachments {
			atts = append(atts, simplifyAttachment(a))
		}
		out["attachments"] = atts
	}
	if p.Content != "" && p.ContentFormat != "" {
		out["content"] = map[string]any{"value": p.Content, "format": p.ContentFormat}
	}
	if len(p.Ancestors) > 0 {
		ancs := make([]map[string]any, 0, len(p.Ancestors))
		for _, a := range p.Ancestors {
			if id := models.AsID(models.MapAny(a, "id")); id != "" {
				ancs = append(ancs, map[string]any{"id": id, "title": models.MapString(a, "title")})
			}
		}
		out["ancestors"] = ancs
	}
	if p.Emoji != "" {
		out["emoji"] = p.Emoji
	}
	if p.PageWidth != "" {
		out["page_width"] = p.PageWidth
	}
	return out
}

// simplifyAttachment converts a raw attachment map to the simplified format.
func simplifyAttachment(data map[string]any) map[string]any {
	result := map[string]any{
		"id":         models.AsID(models.MapAny(data, "id")),
		"type":       models.MapString(data, "type"),
		"status":     models.MapString(data, "status"),
		"title":      models.MapString(data, "title"),
		"media_type": "",
		"file_size":  0,
	}
	if ext := models.MapMap(data, "extensions"); ext != nil {
		result["media_type"] = models.MapString(ext, "mediaType")
		result["file_size"] = models.MapInt(ext, "fileSize")
	}
	if links := models.MapMap(data, "_links"); links != nil {
		if dl := models.MapString(links, "download"); dl != "" {
			result["download_url"] = dl
		}
	}
	if v := models.MapMap(data, "version"); v != nil {
		if n := models.MapInt(v, "number"); n != 0 {
			result["version_number"] = n
		}
		if w := models.MapString(v, "when"); w != "" {
			result["version_when"] = w
		}
	}
	if c := models.MapString(data, "created"); c != "" {
		result["created"] = c
	}
	if v := models.MapMap(data, "version"); v != nil {
		if by := models.MapMap(v, "by"); by != nil {
			if dn := models.MapString(by, "displayName"); dn != "" {
				result["author_display_name"] = dn
			}
		}
	}
	return result
}

var _ = constants.ConfluenceDefaultID
