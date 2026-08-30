// Jira wiki ↔ Markdown converter. DC-only.
package preprocessing

import (
	"fmt"
	"regexp"
	"strings"
)

type JiraPreprocessor struct{ DisableTranslation bool }

func NewJiraPreprocessor(d bool) *JiraPreprocessor { return &JiraPreprocessor{DisableTranslation: d} }

func (p *JiraPreprocessor) JiraToMarkdown(text string) string {
	if p.DisableTranslation {
		return text
	}
	re := regexp.MustCompile(`(?m)^h([1-6])\.\s+(.+)$`)
	text = re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		return strings.Repeat("#", int(parts[1][0]-'0')) + " " + parts[2]
	})
	text = regexp.MustCompile(`\*([^*\n]+)\*`).ReplaceAllString(text, "**$1**")
	text = regexp.MustCompile(`\[([^|]+)\|([^\]]+)\]`).ReplaceAllString(text, "[$1]($2)")
	return text
}

func (p *JiraPreprocessor) MarkdownToJira(text string) string {
	if p.DisableTranslation {
		return text
	}
	re := regexp.MustCompile(`(?m)^(#{1,6})\s+(.*)`)
	text = re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		return "h" + string(rune('0'+len(parts[1]))) + ". " + parts[2]
	})
	text = regexp.MustCompile(`\*\*([^*\n]+)\*\*`).ReplaceAllString(text, "*$1*")
	text = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`).ReplaceAllString(text, "[$1|$2]")
	return text
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Placeholder protector.
type protector struct {
	storage map[string]string
	counter int
}

func (p *protector) extract(text string, re *regexp.Regexp) string {
	if p.storage == nil {
		p.storage = map[string]string{}
	}
	return re.ReplaceAllStringFunc(text, func(m string) string {
		key := fmt.Sprintf("\x00PREFIX%d\x00", p.counter)
		p.counter++
		p.storage[key] = m
		return key
	})
}

func (p *protector) restore(text string) string {
	for k, v := range p.storage {
		text = strings.ReplaceAll(text, k, v)
	}
	return text
}
