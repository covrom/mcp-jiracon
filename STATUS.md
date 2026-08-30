# Статус реализации

**99 инструментов в Go (71 Python DC + 28 DC-only). 160+ тестов, 0 ошибок.**

## Инструменты

| Группа | Python DC | Go | Статус |
| --- | --- | --- | --- |
| Jira read | 31 | 35 | ✅ |
| Jira write | 22 | 23 | ✅ |
| Confluence read | 8 | 23 | ✅ (+15 DC-only) |
| Confluence write | 10 | 18 | ✅ (+8 DC-only) |
| **Всего** | **71** | **99** | **✅** |

## Слои

| Слой | Статус |
| --- | --- |
| Config (env, JiraConfig, ConfluenceConfig) | ✅ |
| URLs (ValidateURLForSSRF) | ✅ |
| Utils (logging, media, pagination, toolsets) | ✅ |
| HTTP client (PAT Bearer, Get/Put/Post/Delete/GetRaw/Multipart) | ✅ |
| Transport (SSRF guard, retry, rate-limit, TLS skip) | ✅ |
| Auth middleware (PAT Bearer + per-request headers) | ✅ |
| Models (JiraIssue, JiraUser, JiraStatus, ConfluencePage, ...) | ✅ |
| Preprocessing (Jira wiki ↔ Markdown) | ✅ |
| MCP server (mark3labs/mcp-go v0.57.0 SSE + stdio) | ✅ |
| Read-only guard | ✅ |
| Dockerfile (multi-stage) | ✅ |
| Helm chart | ✅ |
| CI workflow | — |

## Тесты

```
160+ passing, 0 skipped, 0 failed

internal/atlassian      — PAT auth, TLS verify, Confluence PAT
internal/jira           — 37 fetcher tests
internal/confluence     — 30 fetcher tests (12 original + 18 new)
internal/preprocessing  — Jira wiki ↔ MD round-trip
test/integration        — MCP tool tests (every registered tool via SSE client)
test/inspector          — 99 tools verified via MCP inspector
```

## Старый Python-код

`_old_python/` — оригинальная Python-реализация.
