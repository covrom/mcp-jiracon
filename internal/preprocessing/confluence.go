package preprocessing

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
)

// ConfluencePreprocessor converts Markdown to Confluence storage format (XHTML).
type ConfluencePreprocessor struct{}

// NewConfluencePreprocessor creates a new Confluence preprocessor.
func NewConfluencePreprocessor() *ConfluencePreprocessor {
	return &ConfluencePreprocessor{}
}

// MarkdownToStorage converts Markdown content to Confluence storage format (XHTML).
// Confluence storage format accepts HTML/XHTML fragments directly, so we convert
// markdown to HTML which serves as valid storage format content.
func (p *ConfluencePreprocessor) MarkdownToStorage(markdownContent string, enableHeadingAnchors bool) string {
	if strings.TrimSpace(markdownContent) == "" {
		return ""
	}

	parserOptions := []parser.Option{}
	if enableHeadingAnchors {
		parserOptions = append(parserOptions, parser.WithAutoHeadingID())
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.TaskList,
		),
		goldmark.WithParserOptions(parserOptions...),
		goldmark.WithRendererOptions(
			htmlrenderer.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdownContent), &buf); err != nil {
		// Fallback: wrap raw content in a paragraph
		return "<p>" + strings.ReplaceAll(markdownContent, "\n", "<br/>") + "</p>"
	}

	htmlContent := strings.TrimSpace(buf.String())

	return htmlContent
}

// MarkdownToStorageWithTableLayout converts Markdown to Confluence storage format
// with optional table layout (width) attributes applied.
func (p *ConfluencePreprocessor) MarkdownToStorageWithTableLayout(markdownContent string, enableHeadingAnchors bool, tableLayout string) string {
	result := p.MarkdownToStorage(markdownContent, enableHeadingAnchors)
	if tableLayout != "" {
		result = applyTableLayout(result, tableLayout)
	}
	return result
}

// applyTableLayout adds table width/layout attributes for Confluence storage format.
func applyTableLayout(html, layout string) string {
	var width string
	switch layout {
	case "full-width":
		width = "1800"
	case "wide":
		width = "960"
	case "default":
		width = "760"
	default:
		return html
	}

	// Add data attributes to table tags
	return strings.ReplaceAll(html, "<table", "<table data-table-width=\""+width+"\" data-layout=\""+layout+"\"")
}
