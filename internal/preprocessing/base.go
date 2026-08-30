// Package preprocessing provides Markdown ↔ storage-format conversion
// for Jira and Confluence. Direct port of Python's
// src/mcp_atlassian/preprocessing/.
//
// All real implementations are in sibling files:
//   - html2md.go   → HTMLToMarkdown (Confluence storage → MD)
//   - confluence.go → MarkdownToStorage (MD → Confluence storage XHTML)
//   - jira.go       → JiraToMarkdown, MarkdownToJira
package preprocessing
