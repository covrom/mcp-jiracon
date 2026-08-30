# Repository map

This is the human-readable directory map of `mcp-atlassian`. Every entry
includes the file path, what it contains, and which tool catalog section
references it.

> **Reading order for reimplementation**: `servers/` (tool definitions) →
> `models/` (output shapes) → `jira/` & `confluence/` (REST clients) →
> `preprocessing/` (Markdown/ADF translation) → `utils/` (cross-cutting) →
> `exceptions.py` (errors).

## Top-level

```
mcp-atlassian/
├── pyproject.toml          uv-managed project metadata + tool config (ruff, mypy, pytest)
├── uv.lock                 Frozen dep graph (Python ≥ 3.10, FastMCP, atlassian-python-api, …)
├── Dockerfile              Container image recipe (the official MCP deployable)
├── helm/                   Kubernetes chart for the MCP server
├── docs/                   Public docs (mkdocs) — what users see
│   ├── tools/              Markdown docs of each tool (one file per tool)
│   ├── guides/             How-tos
│   ├── advanced/           Auth, OAuth setup, transport details
│   ├── authentication.mdx  Auth modes explained
│   ├── compatibility.mdx   Cloud vs DC support matrix
│   ├── configuration.mdx   All environment variables
│   ├── http-transport.mdx  Streamable HTTP + auth headers
│   └── tools-reference.mdx Master list of tool names
├── scripts/                OAuth setup helpers (`oauth_setup.py`, etc.)
├── smithery.yaml           Smithery deployment descriptor
├── tests/                  pytest unit + integration suite
│   ├── unit/               Patches HTTP — fast
│   └── integration/        Hits a real Atlassian instance
└── src/mcp_atlassian/      The actual library
```

## `src/mcp_atlassian/` library layout

```
src/mcp_atlassian/
├── __init__.py             Version + package-level exports
├── exceptions.py           MCPAtlassianAuthenticationError (subclass of Exception)
├── servers/                FastMCP server instances + tool definitions
│   ├── __init__.py         Re-exports main_mcp
│   ├── main.py             Top-level server: ASGI app, OAuth proxy, lifespan, mount jira/confluence
│   ├── jira.py             All 63 Jira tools (@jira_mcp.tool(...))
│   ├── confluence.py       All 35 Confluence tools (@confluence_mcp.tool(...))
│   ├── dependencies.py     get_jira_fetcher / get_confluence_fetcher + per-request credential resolution
│   ├── context.py          MainAppContext dataclass (lifespan payload)
│   ├── error_handling.py   ErrorPreservingFastMCP class — wraps FastMCP with tool error preservation
│   ├── async_utils.py      run_jira_fetcher_call (offload sync methods to a thread)
│   ├── oauth_proxy.py      HardenedOAuthProxy — DCR-enabled OAuth proxy provider
│   └── client_storage.py   OAuth client registration storage
│
├── jira/                   JiraFetcher + 21 mixins + JiraClient (wraps atlassian-python-api Jira)
│   ├── __init__.py         JiraFetcher class composing all mixins
│   ├── client.py           JiraClient (atlassian.Jira wrapper with auth-aware headers)
│   ├── config.py           JiraConfig (dataclass) + from_env() factory + is_auth_configured()
│   ├── constants.py        DEFAULT_READ_JIRA_FIELDS, JQL examples, time formats
│   ├── protocols.py        Shared mixin base class
│   ├── utils.py            Small Jira-only helpers
│   ├── attachments.py      AttachmentsMixin
│   ├── boards.py           BoardsMixin
│   ├── comments.py         CommentsMixin
│   ├── customer_requests.py CustomerRequestsMixin (Service Desk Cloud + DC)
│   ├── development.py      DevelopmentMixin (Dev panel)
│   ├── epics.py            EpicsMixin
│   ├── field_options.py    FieldOptionsMixin (cascading select, etc.)
│   ├── fields.py           FieldsMixin (Jira field metadata + search)
│   ├── forms.py            FormsMixin (legacy wrapper)
│   ├── forms_api.py        FormsApiMixin (ProForma forms REST API)
│   ├── forms_common.py     Shared forms helpers
│   ├── formatting.py       FormattingMixin (Markdown helpers)
│   ├── issues.py           IssuesMixin (CRUD on issues)
│   ├── links.py            LinksMixin (issue links)
│   ├── metrics.py          MetricsMixin (issue dates, SLA computation)
│   ├── project_analysis.py ProjectAnalysisMixin (epic hierarchy, cross-project links)
│   ├── projects.py         ProjectsMixin
│   ├── queues.py           QueuesMixin (Server/DC Service Desk queues)
│   ├── search.py           SearchMixin (JQL)
│   ├── sla.py              SLAMixin (cycle/lead time, due date compliance)
│   ├── sprints.py          SprintsMixin
│   ├── transitions.py      TransitionsMixin
│   ├── users.py            UsersMixin (assignable search, profile)
│   ├── watchers.py         WatchersMixin
│   └── worklog.py          WorklogMixin
│
├── confluence/             ConfluenceFetcher + 11 mixins + ConfluenceClient
│   ├── __init__.py         ConfluenceFetcher class
│   ├── client.py           ConfluenceClient (atlassian.Confluence wrapper + auth)
│   ├── config.py           ConfluenceConfig (dataclass) + from_env()
│   ├── constants.py        Toolset names, etc.
│   ├── protocols.py        Mixin base
│   ├── utils.py            Confluence-only helpers
│   ├── analytics.py        AnalyticsMixin (Cloud-only page views)
│   ├── attachments.py      AttachmentsMixin (upload/download/list/delete)
│   ├── comments.py         CommentsMixin (footer + inline comments)
│   ├── labels.py           LabelsMixin
│   ├── pages.py            PagesMixin (CRUD + history + diff + sections)
│   ├── permissions.py      PermissionsMixin (Cloud-only permission check / space perms)
│   ├── restrictions.py     RestrictionsMixin (view/edit restrictions)
│   ├── search.py           SearchMixin (CQL)
│   ├── spaces.py           SpacesMixin
│   ├── templates.py        TemplatesMixin (Cloud-only content templates)
│   ├── users.py            UsersMixin
│   └── v2_adapter.py       ConfluenceV2Adapter — used by OAuth Cloud to talk to /api/v2
│
├── models/                 Pydantic v2 data models
│   ├── __init__.py         Re-exports
│   ├── base.py             ApiModel (from_api_response + to_simplified_dict)
│   ├── constants.py        Common constants
│   ├── jira/               JiraIssue, JiraComment, JiraWorklog, JiraUser, JiraAttachment, …
│   └── confluence/         ConfluencePage, ConfluenceComment, ConfluenceAttachment, ConfluenceSpace, …
│
├── preprocessing/          Content-format conversion
│   ├── __init__.py         Re-exports
│   ├── base.py             Preprocessor base + apply_markdown_transformations()
│   ├── markdown.py         md → ADF/storage conversion helpers
│   ├── html.py             storage → markdown / HTML sanitization
│   ├── confluence.py       Confluence-specific transformation pipeline
│   └── jira.py             Jira-specific transformation (md ↔ ADF)
│
└── utils/                  Cross-cutting helpers
    ├── __init__.py         Re-exports
    ├── env.py              is_env_truthy / get_int_env / get_header_names / is_env_ssl_verify
    ├── environment.py      get_available_services — decides which services are usable
    ├── io.py               is_read_only_mode / validate_safe_path
    ├── logging.py          mask_sensitive / setup_logging
    ├── oauth.py            OAuthConfig dataclass + cloud vs DC URL resolvers
    ├── urls.py             is_atlassian_cloud_url / validate_url_for_ssrf / resolve_relative_url
    ├── tools.py            get_enabled_tools / should_include_tool — ENABLED_TOOLS gating
    ├── toolsets.py         Toolset tag management
    ├── decorators.py       @check_write_access, error wrappers
    ├── proxy.py            Proxy config resolution + transport adapters
    ├── token_verifier.py   AtlassianOpaqueTokenVerifier (Bearer token introspection)
    ├── media.py            ATTACHMENT_MAX_BYTES, fetch_and_encode_attachment, is_image_attachment
    ├── pagination.py       clamp_limit — protects against huge `limit` values
    ├── http.py             Retry / concurrency / rate-limit / circuit-breaker session helpers
    ├── user_agent.py       Custom User-Agent for Atlassian WAFs
    ├── ssl.py              SSL verification adapters
    ├── ssrf_adapter.py     NoProxyAdapter / SSLIgnoreAdapter
    ├── date.py             Date parsing helpers
    └── oauth_setup.py      Standalone CLI for OAuth setup wizard
```

## File counts at a glance

| Group | Count | Combined LOC |
| --- | --- | --- |
| Jira mixin modules | 21 | ~26 000 |
| Confluence mixin modules | 11 | ~5 700 |
| Tool functions | 63 + 35 = 98 | jira.py ≈ 4 500, confluence.py ≈ 3 100 |
| Models | ~50 | ~6 000 |

(All LOC figures approximate, taken from `wc -l` at doc-generation time.)