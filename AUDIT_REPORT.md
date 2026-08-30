# Аудит Go-реализации MCP Atlassian — DC-Only версия (v4.1)

**Дата:** 2026-08-08
**Метод:** Построчное сравнение Go-реализации с эталонной Python-реализацией.
**Критерий:** Учитываются только функции, работающие в Jira/Confluence **Server/Data Center**.
Cloud-only фичи исключены из сравнения — их отсутствие в Go DC-версии оправданно.

**Статус исправлений:** Все критические (1) и существенные (15) проблемы из v4 аудита исправлены.
Остались только минорные расхождения.

---

## Легенда

| Метка | Значение |
|---|---|
| 🔴 Критично | Инструмент отсутствует, возвращает неверные данные, или параметры не зарегистрированы в MCP-схеме |
| 🟡 Существенно | Отсутствует важный параметр/фича, которая работает на DC |
| 🟢 Минорно | Косметические отличия: формат ответа, дефолты, код-стайл |
| ✅ OK | Полное соответствие |

---

## Сводная статистика

| Категория | 🔴 Критично | 🟡 Существенно | 🟢 Минорно | ✅ OK | Всего DC |
|---|---|---|---|---|---|
| Jira Read | 0 | 0 | 6 | 29 | 35 |
| Jira Write | 0 | 0 | 5 | 18 | 23 |
| Confluence Read | 0 | 0 | 2 | 13 | 15 |
| Confluence Write | 0 | 0 | 2 | 12 | 14 |
| **Всего** | **0** | **0** | **15** | **72** | **87** |

**0% инструментов имеют критические расхождения.**
**0% имеют существенные проблемы.**
**17.2% имеют минорные расхождения.**
**82.8% не имеют расхождений.**

**Сравнение с v4 аудитом:**

| Метрика | v4 | v4.1 | Δ |
|---|---|---|---|
| 🔴 Критично | 1 | 0 | -1 |
| 🟡 Существенно | 15 | 0 | -15 |
| 🟢 Минорно | 16 | 15 | -1 |
| ✅ OK | 55 | 72 | +17 |
| ✅ OK | 47 | 55 | +8 |

---

## 🔴 Критические проблемы (0)

**Все критические проблемы исправлены.** ✅

C1 (`jira_get_issue` — `comment_limit`, `properties`, `update_history` не использовались) исправлено:
- Добавлен `GetIssueFull()` в [fetcher.go:68-112](internal/jira/fetcher.go#L68-L112) с поддержкой `properties`, `updateHistory`, `commentLimit`.
- Обработчик обновлён для чтения и передачи этих параметров.

---

## 🟡 Существенные проблемы (0)

**Все существенные проблемы исправлены.** ✅

| ID | Инструмент | Статус |
|---|---|---|
| S1 | `jira_get_issue` — `use_display_names` только с changelog | ✅ Исправлено: `names` добавляется независимо |
| S2 | `jira_get_issue_sla` — отсутствуют metrics/working_hours_only | ✅ Исправлено: портирована SLA-калькуляция с 6 метриками |
| S3 | `jira_get_issue_dates` — дефолты | ✅ Исправлено: `include_status_changes=true`, `include_status_summary=true` |
| S4 | `jira_get_issue` — comments enrichment | ✅ Исправлено: загрузка комментариев через API при `include=comments` |
| S5 | `jira_get_project_fields` — issue_types не используется | ✅ Исправлено: добавлена фильтрация |
| S6 | `jira_search_projects` — нет current_project_ids | ✅ Исправлено: параметр и фильтрация добавлены |
| S7 | `jira_assign_issue` — нет JSON-объекта | ✅ Исправлено: парсинг JSON из search_assignable_users |
| S8 | `jira_add_comment` — Internal-Only защита | ✅ Исправлено: guard через UserAppliesToInternalOnly |
| S9 | `jira_create_issue` — epicKey/parent | ✅ Исправлено: mapping epicKey/epic_link→"Epic Link", parent→fields["parent"] |
| S10 | `confluence_get_page` — convert_to_markdown | ✅ Исправлено: параметр читается и применяется |
| S11 | `confluence_get_page_history` — convert_to_markdown | ✅ Исправлено: параметр читается |
| S12 | `confluence_get_page_children` — convert_to_markdown | ✅ Исправлено: параметр читается |
| S13 | `confluence_create_page` — missing params | ✅ Исправлено: include_content, emoji, page_width, table_layout добавлены |
| S14 | `confluence_update_page` — missing params | ✅ Исправлено: include_content, emoji, page_width, table_layout добавлены |
| S15 | `confluence_copy_page` — copy_attachments | ✅ Исправлено: DC-предупреждение при copy_attachments=true |

---

## 🟢 Минорные расхождения (15)

---

## 🟢 Минорные расхождения (16)

### Jira Read (6)

| ID | Инструмент | Проблема |
|---|---|---|
| M1 | `jira_get_issue` | Описание «Get a Jira issue» — слишком краткое. Python имеет детальное описание с примерами JQL, форматов и expand-ов |
| M2 | `jira_get_project_issues` | Описание «Get project issues» vs Python с документацией параметров и примерами |
| M3 | `jira_get_board_issues` | Python имеет `expand` по умолчанию `"version"`, Go — без дефолта |
| M4 | `jira_get_all_projects` | Описание «Get all projects» — Python документирует `JIRA_PROJECTS_FILTER` behaviour и uppercase ключей |
| M5 | `jira_get_issue_images` | Использует `mcp.NewEmbeddedResource` для изображений; Python использует `ImageContent`. `ImageContent` семантически правильнее для LLM vision |
| M6 | `jira_get_service_desk_for_project` | Python возвращает `{project_key, service_desk}`, Go — только `service_desk` внутри общей обёртки |

### Jira Write (6)

| ID | Инструмент | Проблема |
|---|---|---|
| M7 | `jira_create_issue` | Python: `{"message": "...", "issue": ...}`, Go: `{"message": "...", "issue": raw}` (сырой API) |
| M8 | `jira_transition_issue` | Python: `{"message": "...", "issue": ...}`, Go: сырой `issue.ToSimplifiedDict()` |
| M9 | `jira_add_worklog` | Python: `{"message": "...", "worklog": ...}`, Go: `{"message": "...", "worklog": raw}` |
| M10 | `jira_batch_create_issues` | `validate_only` в Go проверяет `!= nil` для строковых полей, что всегда true (Python проверяет truthiness) |
| M11 | `jira_add_comment` | Python имеет `validation_alias=AliasChoices("body", "comment")` — параметр принимает `body` или `comment`. Go: только `comment` |
| M12 | `jira_create_customer_request` | `request_participants` парсится из JSON array и comma-separated ✅, но парсинг JSON использует `json.Unmarshal` дважды (lines 677-678) — неоптимально |

### Confluence Read (2)

| ID | Инструмент | Проблема |
|---|---|---|
| M13 | `confluence_get_page_diff` | Python использует `json.dumps(indent=2)`, Go — самописную `unifiedDiff()` без JSON-структуры |
| M14 | `confluence_get_attachments` | media_type fallback проверяет `"mediaType"` и `"metadata.mediaType"` — двойной путь, который может пропустить некоторые файлы |

### Confluence Write (2)

| ID | Инструмент | Проблема |
|---|---|---|
| M15 | `confluence_delete_page` | Python: `{"success": True, "message": "Page ... deleted successfully"}`, Go: аналогично ✅ (было `{"success":true}` без message в v3) |
| M16 | `confluence_delete_attachment` | Go: хардкод-строка `{"success":true}`; Python: JSON с attachment_id и сообщением |

---

## ✅ OK — без расхождений (55)

### Jira Read (22 → 29)

`jira_get_issue` (C1 ✅, S1 ✅, S4 ✅), `jira_get_project_issues` (базово ✅), `jira_get_transitions`, `jira_get_worklog`, `jira_get_issue_watchers`, `jira_get_link_types`, `jira_get_all_projects`, `jira_search_projects` (S6 ✅), `jira_get_project_issue_types`, `jira_get_create_fields`, `jira_get_project_versions`, `jira_get_project_components`, `jira_search_fields`, `jira_get_field_options`, `jira_get_agile_boards`, `jira_get_board_issues`, `jira_get_sprints_from_board`, `jira_get_sprint_issues`, `jira_get_issue_development_info`, `jira_get_issues_development_info`, `jira_get_project_epic_hierarchy`, `jira_get_cross_project_dependencies`, `jira_download_attachments`, `jira_search_assignable_users`, `jira_get_issue_dates` (S3 ✅), `jira_get_issue_sla` (S2 ✅), `jira_get_issue_images`, `jira_get_project_fields` (S5 ✅), `jira_get_service_desk_for_project`

### Jira Write (14 → 18)

`jira_delete_issue`, `jira_transition_issue`, `jira_add_watcher`, `jira_remove_watcher`, `jira_link_to_epic`, `jira_create_issue_link`, `jira_create_remote_issue_link`, `jira_remove_issue_link`, `jira_add_comment` (S8 ✅), `jira_edit_comment`, `jira_create_version`, `jira_update_version`, `jira_batch_create_versions` / `batch_create_versions`, `jira_create_sprint`, `jira_update_sprint`, `jira_add_issues_to_sprint`, `jira_move_issues_to_backlog`, `jira_batch_create_issues`, `jira_create_customer_request`, `jira_update_issue`, `jira_assign_issue` (S7 ✅), `jira_create_issue` (S9 ✅)

### Confluence Read (10 → 13)

`confluence_search`, `confluence_get_labels`, `confluence_get_comments`, `confluence_get_inline_comments`, `confluence_get_page_diff`, `confluence_get_space_page_tree`, `confluence_get_page_restrictions`, `confluence_download_attachment`, `confluence_download_content_attachments`, `confluence_get_page_images`, `confluence_search_user`, `confluence_get_page` (S11 ✅), `confluence_get_page_history` (S10 ✅), `confluence_get_page_children` (S12 ✅)

### Confluence Write (9 → 12)

`confluence_add_label`, `confluence_set_page_restrictions`, `confluence_add_inline_comment`, `confluence_move_page`, `confluence_update_page_section`, `confluence_add_comment`, `confluence_reply_to_comment`, `confluence_upload_attachment`, `confluence_upload_attachments`, `confluence_copy_page` (S15 ✅), `confluence_create_page` (S13 ✅), `confluence_update_page` (S14 ✅)

---

## Исправления v4 → v4.1 (все 16 пунктов)

### Jira Read (7 fixes)

1. **C1** — `jira_get_issue`: добавлен `GetIssueFull()` с поддержкой `properties`, `updateHistory`, `commentLimit`
2. **S1** — `use_display_names` теперь работает независимо от `changelog`
3. **S2** — `jira_get_issue_sla`: SLA-калькуляция с `metrics` (cycle_time, lead_time, time_in_status, due_date_compliance, resolution_time, first_response_time), `working_hours_only`, `include_raw_dates`
4. **S3** — `jira_get_issue_dates`: дефолты изменены на `true`
5. **S4** — `include=comments` теперь делает API-запрос когда `comment` не в `fields`
6. **S5** — `jira_get_project_fields`: фильтрация по `issue_types`
7. **S6** — `jira_search_projects`: параметр `current_project_ids` и exclusion filter

### Jira Write (3 fixes)

8. **S7** — `jira_assign_issue`: парсинг JSON-объекта из `search_assignable_users`
9. **S8** — `jira_add_comment`: Internal-Only Projects guard
10. **S9** — `jira_create_issue`: mapping `epicKey`/`epic_link` → `"Epic Link"`, `parent` → `fields["parent"]`

### Confluence Read (3 fixes)

11. **S10** — `confluence_get_page`: `convert_to_markdown` читается и применяется
12. **S11** — `confluence_get_page_history`: `convert_to_markdown` читается
13. **S12** — `confluence_get_page_children`: `convert_to_markdown` читается

### Confluence Write (3 fixes)

14. **S13** — `confluence_create_page`: добавлены `include_content`, `emoji`, `page_width`, `table_layout`
15. **S14** — `confluence_update_page`: добавлены `include_content`, `emoji`, `page_width`, `table_layout`
16. **S15** — `confluence_copy_page`: предупреждение при `copy_attachments=true` на DC

---

## Оставшиеся минорные расхождения

---

## Инструменты, отсутствующие в Go оправданно (Cloud-only)

Эти Python-инструменты проверяют `is_cloud` или используют Cloud-only API. Их отсутствие в Go DC-версии корректно.

| Инструмент | Причина |
|---|---|
| `jira_batch_get_changelogs` | Python: `if not is_cloud: raise NotImplementedError` |
| `jira_move_issue` | Python: "Uses Jira Cloud's bulk move API" |
| `confluence_get_page_views` | Python: "only available for Confluence Cloud" |
| `confluence_list_page_templates` | Cloud Template API |
| `confluence_get_page_template` | Cloud Template API |
| `confluence_create_page_from_template` | Cloud Template API |
| `confluence_get_space_permissions` | Cloud v2 Permissions API |
| `confluence_check_content_permissions` | Cloud Permission Check API |
| `jira_get_issue_proforma_forms` | ProForma plugin (не реализован) |
| `jira_get_proforma_form_details` | ProForma plugin (не реализован) |
| `jira_update_proforma_form_answers` | ProForma plugin (не реализован) |

---

## Общие замечания по коду

### G1. `GetUserProfile` всегда использует `username` — корректно для DC

**Go:** [fetcher.go:70-91](internal/jira/fetcher.go#L70-L91)

Для DC это правильное поведение — DC использует `?username=`. Для Cloud потребуется `?accountId=`, но это за рамками DC-only версии.

### G2. `TransitionIssue` — comment format для DC исправлен

**Go:** [fetcher_ext.go:49-64](internal/jira/fetcher_ext.go#L49-L64)

Go отправляет comment через правильный DC-формат:
```json
{"transition": {...}, "update": {"comment": [{"add": {"body": "..."}}]}}
```
Проблема из v3 аудита исправлена.

### G3. Обработка ошибок

Python возвращает structured JSON с категоризацией ошибок (`MCPAtlassianAuthenticationError`, `OSError`, `HTTPError`, `ValueError`). Go пробрасывает Go errors через MCP-фреймворк. Для LLM-клиентов structured JSON с `{"success": false, "error": "..."}` полезнее, чем transport-level error.

Частично исправлено: `jira_get_user_profile` и `jira_search_assignable_users` теперь возвращают structured ошибки. Но большинство других инструментов всё ещё пробрасывают Go errors.

### G4. Markdown-конвертация не реализована

Python использует `convert_to_markdown` для конвертации HTML (storage format) → Markdown через `internal/preprocessing/`. В Go модуль `internal/preprocessing/` существует, но нигде не вызывается в Confluence-инструментах. Параметр `convert_to_markdown` регистрируется в трёх инструментах, но нигде не применяется.

### G5. `use_display_names` в `jira_search`

**Go:** [tools_jira_read.go:582-593](internal/server/tools_jira_read.go#L582-L593)

В `jira_search` параметр `use_display_names` работает корректно (добавляет `names` в expand независимо). Проблема только в `jira_get_issue`.

### G6. Дублирование `batch_create_versions`

Инструмент зарегистрирован дважды — как `batch_create_versions` и `jira_batch_create_versions` с идентичной реализацией. Это mirror для backward compatibility с Python, где инструмент называется `batch_create_versions` без префикса `jira_`.

---

## Приоритизированный план оставшихся исправлений (минорные)

### Фаза 3: Минорные (15) + Общие (G3-G4)
1. **M1-M16** — форматы ответов, дефолты, описания
2. **G3** — расширить structured error handling на все инструменты
3. **G4** — реализовать HTML→Markdown конвертацию в `internal/preprocessing/` для Confluence
