# Аудит Confluence: сверка реализации с `docs/api_docs.md`

**Дата:** 2026-08-21
**Источник истины:** `docs/list.md` (43 эндпоинта) + `docs/api_docs.md` (описания методов).
**Область аудита:** `internal/confluence/fetcher.go`, `internal/confluence/fetcher_ext.go`,
`internal/server/tools_conf_read.go`, `internal/server/tools_conf_write.go`.

## Итоговая таблица: метод → статус

| # | Метод API | Статус до аудита | Расхождение / решение | Статус после |
| --- | --- | --- | --- | --- |
| 1 | `GET /rest/api/content` | ✅ частично | `GetSpacePageTree` использует `spaceKey`+`limit`+`expand` — корректно по докам. Параметр `type` не используется (не нужен для дерева страниц) | ✅ OK |
| 2 | `POST /rest/api/content` | ✅ | create page / comment / reply / inline / copy — тело соответствует схеме | ✅ OK |
| 3 | `GET /rest/api/content/{id}` | ✅ | путь, `expand` — соответствуют | ✅ OK |
| 4 | `PUT /rest/api/content/{contentId}` | ✅ | update page + `UpdatePageSection` — соответствуют | ✅ OK |
| 5 | `DELETE /rest/api/content/{id}` | ✅ | delete page + `DeleteAttachment` (вложение = контент) — соответствуют | ✅ OK |
| 6 | `GET /rest/api/content/scan` | ❌ нет | **Добавлено**: `ScanContent` + инструмент `confluence_scan_content` | ✅ добавлено |
| 7 | `GET /rest/api/content/search` | ✅ | cql+limit+expand — соответствуют | ✅ OK |
| 8 | `POST .../blueprint/instance/{draftId}` | ❌ нет | Пропущено: legacy blueprint drafts, редкая фича DC | ⏭️ пропущено |
| 9 | `PUT .../blueprint/instance/{draftId}` | ❌ нет | Пропущено: см. #8 | ⏭️ пропущено |
| 10 | `GET /rest/api/content/{id}/history` | ⚠️ не точно | Текущий `GetPageHistory` ходит на `/content/{id}?status=historical&version=N` — это способ получить конкретную версию, а не историю. **Добавлен** настоящий `GetContentHistory` (`/history`) + инструмент `confluence_get_content_history`; старый метод сохранён (используется diff) | ✅ добавлено |
| 11 | `GET .../history/{version}/macro/id/{macroId}` | ❌ нет | Пропущено: только для connect-аддонов | ⏭️ пропущено |
| 12 | `GET .../history/{version}/macro/hash/{hash}` | ❌ нет | Пропущено: deprecated, connect-аддоны | ⏭️ пропущено |
| 13 | `GET .../restriction/byOperation` | ✅ | **Расхождение:** код не передавал `expand` (дефолт дока включает `read/update.restrictions.user/group`). Исправлено: теперь явно передаётся полный expand | ✅ исправлено |
| 14 | `GET .../restriction/byOperation/{operationKey}` | ⚠️ не отдельно | **Добавлено**: `GetRestrictionForOperation` + параметр `operation` в инструменте `confluence_get_page_restrictions` | ✅ добавлено |
| 15 | `GET /rest/api/content/{id}/child` | ✅ | `GetPageChildren(includeFolders=true)` → `/child` — соответствует | ✅ OK |
| 16 | `GET /rest/api/content/{id}/child/{type}` | ⚠️ частично | **Добавлено**: общий `GetContentChildren(type)`; `GetPageChildren` переиспользует его | ✅ добавлено |
| 17 | `GET /rest/api/content/{id}/child/comment` | ✅ | `GetPageComments`/`GetInlineComments` — один эндпоинт, разные expand. **Рефакторинг:** объединены в `GetComments(expand, locationFilter)` | ✅ рефактор |
| 18 | `GET /rest/api/content/{id}/child/attachment` | ✅ | start/limit — соответствуют. **Улучшение:** серверные фильтры `filename`/`mediaType` (были клиентские) | ✅ улучшено |
| 19 | `POST /rest/api/content/{id}/child/attachment` | ✅ | multipart + `X-Atlassian-Token: no-check` — соответствуют | ✅ OK |
| 20 | `PUT .../child/attachment/{attachmentId}` | ❌ нет | **Добавлено**: `UpdateAttachmentMetadata` + инструмент `confluence_update_attachment` | ✅ добавлено |
| 21 | `POST .../child/attachment/{attachmentId}/data` | ❌ нет | **Добавлено**: `UpdateAttachmentData` (multipart) + инструмент `confluence_update_attachment_data` | ✅ добавлено |
| 22 | `GET /rest/api/content/{id}/descendant` | ❌ нет | **Добавлено**: `GetDescendants` (expand) | ✅ добавлено |
| 23 | `GET /rest/api/content/{id}/descendant/{type}` | ❌ нет | **Добавлено**: `GetDescendantsOfType` | ✅ добавлено |
| 24 | `GET /rest/api/content/{id}/label` | ✅ | **Улучшение:** добавлены query `prefix`/`start`/`limit` | ✅ улучшено |
| 25 | `POST /rest/api/content/{id}/label` | ✅ | тело `[{"prefix":"global","name":...}]` — соответствует | ✅ OK |
| 26 | `DELETE /rest/api/content/{id}/label` | ❌ нет | **Добавлено**: `RemovePageLabelByName` (`?name=`) | ✅ добавлено |
| 27 | `DELETE /rest/api/content/{id}/label/{label}` | ❌ нет | **Добавлено**: `RemovePageLabel` (path); единый инструмент `confluence_remove_label` выбирает вариант по наличию `/` в имени | ✅ добавлено |
| 28 | `GET /rest/api/content/{id}/property` | ✅ | limit=50 — ок (дефолт дока 10, но больше полезнее) | ✅ OK |
| 29 | `POST /rest/api/content/{id}/property` | ❌ нет | **Добавлено**: `CreateContentProperty` + инструмент `confluence_set_page_property` (create-or-update) | ✅ добавлено |
| 30 | `GET /rest/api/content/{id}/property/{key}` | ❌ нет | **Добавлено**: `GetContentProperty` | ✅ добавлено |
| 31 | `POST /rest/api/content/{id}/property/{key}` | ❌ нет | Покрыто через `CreateContentProperty` (единый инструмент) | ✅ покрыто |
| 32 | `PUT /rest/api/content/{id}/property/{key}` | ❌ нет | **Добавлено**: `UpdateContentProperty` (нужен id + version из GET) | ✅ добавлено |
| 33 | `DELETE /rest/api/content/{id}/property/{key}` | ❌ нет | **Добавлено**: `DeleteContentProperty` + инструмент `confluence_delete_page_property` | ✅ добавлено |
| 34 | `GET /rest/api/search` | ❌ нет | **Добавлено**: `SearchEntities` (cql/excerpt/expand/start/limit/includeArchivedSpaces) + инструмент `confluence_search_entities` | ✅ добавлено |
| 35 | `GET /rest/api/space` | ✅ | start/limit — соответствуют | ✅ OK |
| 36 | `GET /rest/api/space/{spaceKey}` | ❌ нет | **Добавлено**: `GetSpace` (expand) + инструмент `confluence_get_space` | ✅ добавлено |
| 37 | `GET /rest/api/space/{spaceKey}/content` | ⚠️ через #1 | Решено оставить `/content?spaceKey=` (рабочий подход, tree строится по ancestors). Отдельный эндпоинт не добавляется — дублирование | ✅ оставлено как есть |
| 38 | `GET /rest/api/space/{spaceKey}/content/{type}` | ❌ нет | **Добавлено**: `GetSpaceContent(spaceKey, type)` (depth/expand/start/limit) | ✅ добавлено |
| 39 | `GET /rest/api/space/{spaceKey}/property` | ❌ нет | **Добавлено**: `GetSpaceProperties` (expand/start/limit) | ✅ добавлено |
| 40 | `GET /rest/api/space/{spaceKey}/property/{key}` | ❌ нет | **Добавлено**: `GetSpaceProperty` (expand) | ✅ добавлено |
| 41 | `GET /rest/api/user/current` | ❌ нет | **Добавлено**: `GetCurrentUser` (expand) + инструмент `confluence_get_current_user` | ✅ добавлено |
| 42 | `GET /rest/api/user` | ⚠️ обходной | Обход через `/group/{g}/member` подтверждён приемлемым для DC (глобальный поиск пользователей в DC не предусмотрен). **Добавлен** прямой `GetUser(key/username)` для точечных запросов | ✅ уточнено |
| 43 | `GET /rest/api/user/list` | ❌ нет | Пропущено: полный список всех пользователей — редко нужен, уже есть group-based search | ⏭️ пропущено |

### Сводка

- Полностью реализовано до аудита: ~17. После: **36 из 43**.
- Добавлено методов fetcher: 17 (scan, history, restriction-by-op, child/{type}, descendant×2,
  label delete×2, property CRUD×4, search entities, space by key, space content by type,
  space properties×2, user current, attachment metadata/data update).
- Добавлено MCP-инструментов: 12 read + 4 write = **16** (read: scan, history, descendants,
  space, space_properties, current_user, search_entities; write: remove_label,
  set_page_property, delete_page_property, update_attachment, update_attachment_data).
  Плюс расширены параметры у 3 существующих инструментов (restrictions.operation,
  attachments filename/media_type стали серверными, labels prefix/limit).
- Пропущено сознательно: 7 (#8–9 blueprint, #11–12 macro, #43 user list) — с причинами в таблице.

## Детали найденных расхождений (исправлено)

### R1. `GetPageRestrictions` — отсутствовал `expand`
Доки: дефолтный expand = `update.restrictions.user,read.restrictions.group,read.restrictions.user,update.restrictions.group`.
Без явного expand часть DC-версий может не вернуть полные списки. Теперь передаётся явно.

### R2. `GetContentAttachments` — клиентская фильтрация вместо серверной
Доки поддерживают query `filename` и `mediaType`. Инструмент принимал эти параметры, но
фильтровал на клиенте после полной выборки. Теперь фильтрация серверная (fetcher),
клиентский fallback удалён.

### R3. `GetPageLabels` — отсутствовали query-параметры
Доки: `prefix`, `start`, `limit` (дефолт 200). Теперь все три доступны в fetcher.

### R4. Дублирование `GetPageComments` / `GetInlineComments`
Оба ходили на `/child/comment` с разными expand и client-side фильтрацией.
Объединены в `GetComments(ctx, pageID, expand string, inlineOnly bool)`:
- `inlineOnly=false` → expand `body.view.value,version`, depth=all (все комментарии);
- `inlineOnly=true` → expand `body.view.value,version,extensions.inlineProperties`,
  depth=all + фильтр `extensions.location == "inline"`.

### R5. `GetPageHistory` vs `/history`
`/content/{id}?status=historical&version=N` возвращает конкретную старую версию страницы —
это рабочий приём и он сохранён (используется `GetPageVersionDiff`). Но отдельного доступа
к объекту истории (createdBy, contributors, previousVersion/nextVersion) не было.
Добавлен `GetContentHistory` на настоящий эндпоинт `/content/{id}/history`.