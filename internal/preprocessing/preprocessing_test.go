package preprocessing

import (
	"strings"
	"testing"
)

func TestJiraToMarkdown_Bold(t *testing.T) {
	p := NewJiraPreprocessor(false)
	got := p.JiraToMarkdown("this is *bold* text")
	if !strings.Contains(got, "**bold**") {
		t.Errorf("expected **bold**, got: %s", got)
	}
}

func TestJiraToMarkdown_Links(t *testing.T) {
	p := NewJiraPreprocessor(false)
	got := p.JiraToMarkdown("[text|http://example.com]")
	if !strings.Contains(got, "[text](http://example.com)") {
		t.Errorf("got: %s", got)
	}
}

func TestMarkdownToJira_Headings(t *testing.T) {
	p := NewJiraPreprocessor(false)
	got := p.MarkdownToJira("# Hello\n## World")
	if !strings.Contains(got, "h1. Hello") {
		t.Errorf("got: %s", got)
	}
	if !strings.Contains(got, "h2. World") {
		t.Errorf("got: %s", got)
	}
}

func TestMarkdownToJira_Links(t *testing.T) {
	p := NewJiraPreprocessor(false)
	got := p.MarkdownToJira("[text](http://example.com)")
	if !strings.Contains(got, "[text|http://example.com]") {
		t.Errorf("got: %s", got)
	}
}

func TestConfluenceMarkdownToStorage_Empty(t *testing.T) {
	p := NewConfluencePreprocessor()
	got := p.MarkdownToStorage("", false)
	if got != "" {
		t.Errorf("expected empty string, got: %s", got)
	}
}

func TestConfluenceMarkdownToStorage_Bold(t *testing.T) {
	p := NewConfluencePreprocessor()
	got := p.MarkdownToStorage("**bold**", false)
	if !strings.Contains(got, "<strong>") {
		t.Errorf("expected <strong> in output, got: %s", got)
	}
}

func TestConfluenceMarkdownToStorage_Headings(t *testing.T) {
	p := NewConfluencePreprocessor()
	got := p.MarkdownToStorage("# Hello\n## World", false)
	if !strings.Contains(got, "<h1>") || !strings.Contains(got, "<h2>") {
		t.Errorf("expected h1 and h2 tags, got: %s", got)
	}
}

func TestConfluenceMarkdownToStorage_HeadingAnchors(t *testing.T) {
	p := NewConfluencePreprocessor()
	got := p.MarkdownToStorage("# Hello World", true)
	if !strings.Contains(got, "id=") {
		t.Errorf("expected heading id attribute, got: %s", got)
	}
}

func TestConfluenceMarkdownToStorage_TableLayout(t *testing.T) {
	p := NewConfluencePreprocessor()
	md := "| A | B |\n|---|---|\n| 1 | 2 |"
	got := p.MarkdownToStorageWithTableLayout(md, false, "full-width")
	if !strings.Contains(got, "data-layout=\"full-width\"") {
		t.Errorf("expected data-layout attribute, got: %s", got)
	}
}

func TestConfluenceMarkdownToStorageWithTableLayout_NoLayout(t *testing.T) {
	p := NewConfluencePreprocessor()
	got := p.MarkdownToStorageWithTableLayout("**bold**", false, "")
	if !strings.Contains(got, "<strong>") {
		t.Errorf("expected <strong>, got: %s", got)
	}
	if strings.Contains(got, "data-layout") {
		t.Errorf("should not have data-layout when empty, got: %s", got)
	}
}
