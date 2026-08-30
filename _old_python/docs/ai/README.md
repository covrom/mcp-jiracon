# `mcp-atlassian` — AI-facing documentation

> **Audience**: an LLM / engineering agent that must understand every MCP tool the
> Python `mcp-atlassian` server exposes, in order to **reimplement** the same
> surface in another language (target reference: **Go**).
>
> **Goal of this folder**: complete, prescriptive specs — prompts, tool-schema,
> handler semantics, downstream Atlassian REST calls, Cloud-vs-DC differences,
> parameter boundaries, and common variants. After reading this folder you should
> be able to write a Go server that, given the same Atlassian configuration,
> presents an equivalent set of MCP tools to an LLM client.

---

## How to read this folder

The folder is split into three concentric layers:

```
docs/ai/
├── 00-overview/        high-level map of the codebase and tool surface
├── 01-architecture/    layered architecture: servers → fetchers → mixins → HTTP
├── 02-server-lifecycle/ ASGI / lifespan / per-request credential resolution
├── 03-auth/            OAuth 2.0, PAT, Basic, multi-tenant header auth
├── 04-shared-utils/    env, logging, IO, media, validation, decorators
├── 05-models/          Pydantic-v2-style data models + simplified dicts
├── 06-preprocessing/   ADF / storage → Markdown conversion
├── 10-jira-tools/      one file per Jira tool category (read, write, etc.)
├── 20-confluence-tools/ one file per Confluence tool category
└── 99-migration/       Go-rewrite cheat-sheet, terminology map, library picks
```

Read **00-overview** first, then **01-architecture**. The tool catalogs
(`10-` and `20-`) are reference; jump straight to a tool's file when you need
to reimplement it.

---

## What `mcp-atlassian` actually is

`mcp-atlassian` is an [MCP](https://modelcontextprotocol.io/) server that gives an
LLM the ability to read & write **Atlassian Jira** and **Confluence** through a
single mounted pair of FastMCP sub-servers.

At runtime it exposes **98 MCP tools**:

| Group | Count | File |
| --- | --- | --- |
| Jira tools | **63** | `src/mcp_atlassian/servers/jira.py` |
| Confluence tools | **35** | `src/mcp_atlassian/servers/confluence.py` |

Each tool is a thin async wrapper:

1. Resolves a per-request credential-aware fetcher (`get_jira_fetcher(ctx)` or
   `get_confluence_fetcher(ctx)`).
2. Calls a method on the fetcher (which composes ~20 mixins, each holding
   client + high-level methods that hit Atlassian REST).
3. Serializes the return value to **JSON string** (or returns a
   `list[TextContent | EmbeddedResource | ImageContent]` for binary content) and
   hands it back to the LLM.

The Atlassian-side code is built on top of the third-party
[`atlassian-python-api`](https://atlassian-python-api.readthedocs.io/) package,
which is just a thin wrapper over Atlassian's REST APIs. So the "real" work is
on Atlassian's HTTP endpoints; `mcp-atlassian` adds:

- MCP tool framing (prompts, JSON-schema, return shaping).
- Multi-tenant credential resolution (PAT / OAuth / Basic, per-request).
- Markdown ↔ ADF/storage conversion for Jira descriptions and Confluence pages.
- Toolset gating (`toolset:jira_issues`, `toolset:confluence_pages`, …).
- Read-only mode enforcement.
- Cross-version unification (Cloud + Server/DC) with explicit branches.

---

## Mental model for reimplementation

```
LLM (client)
    │  MCP call: tools/call {name, arguments}
    ▼
Atlassian MCP server (FastMCP)        ── 02-server-lifecycle
    │  tool filter (auth, read-only, toolset, service availability)
    ▼
@tool(...) function                    ── 10- / 20- tool catalogs
    │  fetcher = await get_jira_fetcher(ctx)   ── 03-auth
    ▼
JiraFetcher / ConfluenceFetcher       ── 01-architecture
    │  mixin method (IssuesMixin.create_issue, PagesMixin.get_page, …)
    ▼
Jira / Confluence client (atlassian-python-api or your own)
    │  HTTP request with auth headers (Basic / Bearer / PAT)
    ▼
Atlassian REST API  (Cloud or DC)
```

The Go reimplementation can collapse a few of these layers (you don't need a
mixin tree if Go interfaces do the same job) but the **external contract** —
MCP tool names, JSON-schema, return shape, error semantics — must match.

---

## Tool-name convention

All tool names follow `{service}_{action}_{target}` (e.g.
`jira_create_issue`, `confluence_get_page`). The two sub-servers are mounted on
the main `Atlassian MCP` server under namespaces `jira` and `confluence`
respectively, so the fully-qualified MCP name is e.g. `jira_get_issue` and
`confluence_search`. See [02-server-lifecycle/mounting.md](02-server-lifecycle/mounting.md).

---

## Tag system

Each `@<svc>_mcp.tool(...)` decorator carries a set of `tags`:

| Tag | Meaning |
| --- | --- |
| `jira` / `confluence` | which service owns the tool (used by the listing filter) |
| `read` / `write` | whether the tool mutates; write tools are blocked in `READ_ONLY_MODE=true` |
| `toolset:<name>` | which toolset the tool belongs to (see `ENABLED_TOOLSETS`) |
| `<other>` (e.g. `attachments`, `metrics`, `sla`, `development`, `analytics`) | free-form grouping hints, not enforced |

Example:

```python
@jira_mcp.tool(
    tags={"jira", "write", "toolset:jira_issues"},
    annotations={"title": "Create Issue", "destructiveHint": False},
)
```

The Go equivalent should keep this tag set verbatim because operators use it to
gate behavior via `READ_ONLY_MODE`, `ENABLED_TOOLS`, `ENABLED_TOOLSETS`.

---

## Return-shape conventions

Almost every tool returns a **JSON string**:

```python
return json.dumps(result, indent=2, ensure_ascii=False)
```

The exceptions are binary outputs that need to bypass the JSON wrapper:

| Tool(s) | Returns |
| --- | --- |
| `jira_download_attachments`, `confluence_download_attachment`, `confluence_download_content_attachments` | `list[TextContent | EmbeddedResource]` — base64 blobs |
| `jira_get_issue_images`, `confluence_get_page_images` | `list[TextContent | ImageContent]` — inline images |

Errors are usually wrapped as `{"success": false, "error": "...", ...}` and
returned as a JSON string (not raised). A few tools raise `ValueError` for
hard-fail cases (missing arg, read-only mode, OAuth scope errors). Tools that
fail on a *network/API* call log the error and return `{"success": false, ...}`
so the LLM can continue.

The Go implementation should preserve both styles: most tools → JSON-string
return; the binary tools → multi-content blocks.

---

## Cloud vs Server / Data Center

Every Atlassian endpoint has slight differences:

- Auth header construction: Cloud uses `Authorization: Basic email:token`,
  DC uses `Authorization: Basic username:password` or `Bearer <pat>`.
- Base URL shape: Cloud is `*.atlassian.net`, DC is host-specific.
- API versioning: Cloud prefers `/rest/api/3/...`, DC `/rest/api/2/...`.
- Endpoint features: e.g. Cloud has the bulk-move endpoint, field-context
  options, ProForma forms; DC has Service Desk queues, dashboard analytics is
  Cloud-only.

`config.is_cloud` is the discriminator everywhere. The Go reimplementation
needs an equivalent discriminator on its `JiraConfig` / `ConfluenceConfig`.

---

## Where to start

- [00-overview/repository-map.md](00-overview/repository-map.md) — what file
  does what.
- [00-overview/tool-catalog.md](00-overview/tool-catalog.md) — quick index of
  every tool with a one-line description.
- [01-architecture/overview.md](01-architecture/overview.md) — layered diagram
  and module map.
- [99-migration/go-cheatsheet.md](99-migration/go-cheatsheet.md) — terminology
  and library picks if you intend to write the Go port right now.