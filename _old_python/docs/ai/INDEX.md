# `docs/ai/` — table of contents

This folder is the AI-facing specification for the `mcp-atlassian`
Python codebase. It exists to enable a faithful reimplementation in
another language (the reference target is Go).

> If you're an LLM agent asked to "reimplement `mcp-atlassian` in Go",
> **start with [00-overview/](./00-overview/)** and read top-to-bottom.

## Layer 0 — Orientation

| File | What's in it |
| --- | --- |
| [README.md](README.md) | High-level orientation: what `mcp-atlassian` is, mental model, return-shape conventions |
| [00-overview/repository-map.md](00-overview/repository-map.md) | Every file in `src/mcp_atlassian/` and what it does |
| [00-overview/tool-catalog.md](00-overview/tool-catalog.md) | Index of all 98 MCP tools (63 Jira + 35 Confluence) with toolset grouping |

## Layer 1 — Architecture

| File | What's in it |
| --- | --- |
| [01-architecture/overview.md](01-architecture/overview.md) | Layered diagram, module dependencies, mixin composition |
| [01-architecture/jira-rest-api.md](01-architecture/jira-rest-api.md) | Per-mixin mapping to Jira REST endpoints |
| [01-architecture/confluence-rest-api.md](01-architecture/confluence-rest-api.md) | Per-mixin mapping to Confluence REST endpoints (v1 + v2) |
| [01-architecture/error-handling.md](01-architecture/error-handling.md) | 5-layer error handling model; per-pattern recipes |

## Layer 2 — Server lifecycle

| File | What's in it |
| --- | --- |
| [02-server-lifecycle/overview.md](02-server-lifecycle/overview.md) | Server construction, mounting, lifespan, tool filter, ASGI middleware |

## Layer 3 — Authentication

| File | What's in it |
| --- | --- |
| [03-auth/per-request-resolution.md](03-auth/per-request-resolution.md) | Auth modes, header parsing, OAuth 2.0, validation cache, SSRF protection |

## Layer 4 — Shared utilities

| File | What's in it |
| --- | --- |
| [04-shared-utils/overview.md](04-shared-utils/overview.md) | Per-module catalog of `utils/` with public APIs and Go equivalents |

## Layer 5 — Data models

| File | What's in it |
| --- | --- |
| [05-models/overview.md](05-models/overview.md) | Pydantic-style data model classes; `to_simplified_dict` / `from_api_response` contract; constants |

## Layer 6 — Preprocessing

| File | What's in it |
| --- | --- |
| [06-preprocessing/overview.md](06-preprocessing/overview.md) | Markdown ↔ ADF/storage converters for Jira and Confluence |

## Layer 7 — Jira tools

| File | Toolset | # tools |
| --- | --- | --- |
| [10-jira-tools/index.md](10-jira-tools/index.md) | – | 1 (overview) |
| [10-jira-tools/jira_users.md](10-jira-tools/jira_users.md) | `jira_users` | 2 |
| [10-jira-tools/jira_watchers.md](10-jira-tools/jira_watchers.md) | `jira_watchers` | 3 |
| [10-jira-tools/jira_issues.md](10-jira-tools/jira_issues.md) | `jira_issues` | 10 |
| [10-jira-tools/jira_transitions.md](10-jira-tools/jira_transitions.md) | `jira_transitions` | 2 |
| [10-jira-tools/jira_worklog.md](10-jira-tools/jira_worklog.md) | `jira_worklog` | 2 |
| [10-jira-tools/jira_attachments.md](10-jira-tools/jira_attachments.md) | `jira_attachments` | 2 |
| [10-jira-tools/jira_agile.md](10-jira-tools/jira_agile.md) | `jira_agile` | 8 |
| [10-jira-tools/jira_links.md](10-jira-tools/jira_links.md) | `jira_links` | 5 |
| [10-jira-tools/jira_comments.md](10-jira-tools/jira_comments.md) | `jira_comments` | 2 |
| [10-jira-tools/jira_projects.md](10-jira-tools/jira_projects.md) | `jira_projects` | 10 |
| [10-jira-tools/jira_service_desk.md](10-jira-tools/jira_service_desk.md) | `jira_service_desk` | 6 |
| [10-jira-tools/jira_fields.md](10-jira-tools/jira_fields.md) | `jira_fields` | 2 |
| [10-jira-tools/jira_forms.md](10-jira-tools/jira_forms.md) | `jira_forms` | 3 |
| [10-jira-tools/jira_metrics.md](10-jira-tools/jira_metrics.md) | `jira_metrics` | 2 |
| [10-jira-tools/jira_development.md](10-jira-tools/jira_development.md) | `jira_development` | 2 |
| [10-jira-tools/jira_project_analysis.md](10-jira-tools/jira_project_analysis.md) | `jira_project_analysis` | 2 |

## Layer 8 — Confluence tools

| File | Toolset | # tools |
| --- | --- | --- |
| [20-confluence-tools/index.md](20-confluence-tools/index.md) | – | 1 (overview) |
| [20-confluence-tools/confluence_pages.md](20-confluence-tools/confluence_pages.md) | `confluence_pages` | 14 |
| [20-confluence-tools/confluence_comments.md](20-confluence-tools/confluence_comments.md) | `confluence_comments` | 5 |
| [20-confluence-tools/confluence_labels.md](20-confluence-tools/confluence_labels.md) | `confluence_labels` | 2 |
| [20-confluence-tools/confluence_users.md](20-confluence-tools/confluence_users.md) | `confluence_users` | 1 |
| [20-confluence-tools/confluence_analytics.md](20-confluence-tools/confluence_analytics.md) | `confluence_analytics` | 1 |
| [20-confluence-tools/confluence_attachments.md](20-confluence-tools/confluence_attachments.md) | `confluence_attachments` | 7 |
| [20-confluence-tools/confluence_templates.md](20-confluence-tools/confluence_templates.md) | `confluence_templates` | 3 |
| [20-confluence-tools/confluence_permissions.md](20-confluence-tools/confluence_permissions.md) | `confluence_permissions` | 2 |

## Layer 9 — Migration

| File | What's in it |
| --- | --- |
| [99-migration/go-cheatsheet.md](99-migration/go-cheatsheet.md) | Go reimplementation cheat-sheet: libraries, module map, HTTP client design, OAuth, JSON schema, recommended migration order |
| [99-migration/environment-variables.md](99-migration/environment-variables.md) | Every env var the server reads; required for parity |

---

## Quick lookup

| "I need to know..." | Read this |
| --- | --- |
| what files exist in `src/mcp_atlassian/` | [00-overview/repository-map.md](00-overview/repository-map.md) |
| what tools exist | [00-overview/tool-catalog.md](00-overview/tool-catalog.md) |
| how a single Jira tool works (e.g. `get_issue`) | [10-jira-tools/jira_issues.md](10-jira-tools/jira_issues.md) |
| how a single Confluence tool works (e.g. `get_page`) | [20-confluence-tools/confluence_pages.md](20-confluence-tools/confluence_pages.md) |
| what Atlassian REST endpoint a tool calls | [01-architecture/jira-rest-api.md](01-architecture/jira-rest-api.md) / [confluence-rest-api.md](01-architecture/confluence-rest-api.md) |
| how auth is resolved per request | [03-auth/per-request-resolution.md](03-auth/per-request-resolution.md) |
| what env vars exist | [99-migration/environment-variables.md](99-migration/environment-variables.md) |
| which Go libraries to use | [99-migration/go-cheatsheet.md](99-migration/go-cheatsheet.md#library-picks) |
| how errors are classified | [01-architecture/error-handling.md](01-architecture/error-handling.md) |
| how Markdown is converted | [06-preprocessing/overview.md](06-preprocessing/overview.md) |
| how the model classes serialize output | [05-models/overview.md](05-models/overview.md) |