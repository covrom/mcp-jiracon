# План: проверка и рефакторинг реализации методов Confluence

## Цель

Пройтись по всем методам Confluence REST API v1 (DC) из `docs/list.md`, сверить их с
реализацией в `internal/confluence/` и `internal/server/tools_conf_*.go`, выявить
расхождения с документацией (`docs/api_docs.md`) и провести рефакторинг, чтобы код
точно соответствовал описанному API.

## Источник истины для API

- **Список методов**: `docs/list.md` (40 эндпоинтов, 10 областей).
- **Описание каждого метода**: `docs/api_docs.md` — каждый метод оформлен как H2-заголовок
  вида `## METHOD \`path\`` (например `## GET \`/rest/api/content/{id}\``). Описание метода
  начинается с этого заголовка и заканчивается строкой перед следующим H2-заголовком.

### Как вырезать часть документа на метод (скриптом)

Документ устроен так: раздел метода = от его H2 до следующего H2. Для точного сохранения
описания использовать `sed`/`awk` по номерам строк или по заголовкам:

```bash
# Вариант A: по номеру строки заголовка (известен из grep -nE '^## ')
sed -n '68,99p' docs/api_docs.md          # GET /rest/api/content/{id} (строки 68..99)

# Вариант B: awk между двумя H2-заголовками (от METHOD до следующего ##)
awk '/^## GET `\/rest\/api\/content\/{id}`$/{f=1; print; next} /^## /{if(f) exit} f' docs/api_docs.md

# Вариант C: получить номера всех H2 разом для построения карты разделов
grep -nE '^## ' docs/api_docs.md
```

Карта разделов (H2 → строка начала), полученная из `grep -nE '^## '` — см. ниже в разделе
«Методы и их расположение в api_docs.md».

## Текущая реализация (что есть)

| Файл | Назначение |
| --- | --- |
| `internal/confluence/fetcher.go` | Базовые методы: `GetPage`, `GetPageByTitle`, `GetPageContent`, `Search` |
| `internal/confluence/fetcher_ext.go` | Расширенные: комментарии, лейблы, дети, вложения, spaces, ограничения, свойства, история/diff, копирование |
| `internal/server/tools_conf_read.go` | MCP read-инструменты (15) |
| `internal/server/tools_conf_write.go` | MCP write-инструменты (13) |

## Методы и их расположение в `api_docs.md`

| # | Метод | Строка H2 | Реализовано? | Где (fetcher) |
| --- | --- | --- | --- | --- |
| 1 | `GET /rest/api/content` | 3 | ✅ частично | `GetSpacePageTree` (spaceKey+limit+expand) |
| 2 | `POST /rest/api/content` | 40 | ✅ | create page / add comment / reply / inline / copy |
| 3 | `GET /rest/api/content/{id}` | 68 | ✅ | `GetPage` |
| 4 | `PUT /rest/api/content/{contentId}` | 100 | ✅ | update page / `UpdatePageSection` |
| 5 | `DELETE /rest/api/content/{id}` | 150 | ✅ | delete page / `DeleteAttachment` |
| 6 | `GET /rest/api/content/scan` | 177 | ❌ нет | — |
| 7 | `GET /rest/api/content/search` | 197 | ✅ | `Search` (cql+limit+expand) |
| 8 | `POST /rest/api/content/blueprint/instance/{draftId}` | 232 | ❌ нет | — |
| 9 | `PUT /rest/api/content/blueprint/instance/{draftId}` | 259 | ❌ нет | — |
| 10 | `GET /rest/api/content/{id}/history` | 286 | ⚠️ не точно | `GetPageHistory` использует `status=historical&version=N` на `/content/{id}`, а не `/history` |
| 11 | `GET .../history/{version}/macro/id/{macroId}` | 316 | ❌ нет | — |
| 12 | `GET .../history/{version}/macro/hash/{hash}` | 344 | ❌ нет | — |
| 13 | `GET /rest/api/content/{id}/restriction/byOperation` | 374 | ✅ | `GetPageRestrictions` |
| 14 | `GET .../restriction/byOperation/{operationKey}` | 401 | ⚠️ не отдельно | только byOperation (без ключа операции) |
| 15 | `GET /rest/api/content/{id}/child` | 430 | ✅ | `GetPageChildren` (includeFolders=true) |
| 16 | `GET /rest/api/content/{id}/child/{type}` | 465 | ⚠️ частично | `GetPageChildren` хардкодит `/child/page`; общий `{type}` отсутствует |
| 17 | `GET /rest/api/content/{id}/child/comment` | 498 | ✅ | `GetPageComments`, `GetInlineComments` |
| 18 | `GET /rest/api/content/{id}/child/attachment` | 533 | ✅ | `GetContentAttachments` |
| 19 | `POST /rest/api/content/{id}/child/attachment` | 567 | ✅ | `UploadAttachment`, `UploadAttachmentBase64` |
| 20 | `PUT .../child/attachment/{attachmentId}` | 606 | ❌ нет | — (обновление метаданных вложения) |
| 21 | `POST .../child/attachment/{attachmentId}/data` | 647 | ❌ нет | — (обновление бинарных данных) |
| 22 | `GET /rest/api/content/{id}/descendant` | 679 | ❌ нет | — |
| 23 | `GET /rest/api/content/{id}/descendant/{type}` | 712 | ❌ нет | — |
| 24 | `GET /rest/api/content/{id}/label` | 744 | ✅ | `GetPageLabels` |
| 25 | `POST /rest/api/content/{id}/label` | 776 | ✅ | `AddPageLabel` |
| 26 | `DELETE /rest/api/content/{id}/label` | 810 | ❌ нет | — (по query `?name=`) |
| 27 | `DELETE /rest/api/content/{id}/label/{label}` | 834 | ❌ нет | — |
| 28 | `GET /rest/api/content/{id}/property` | 857 | ✅ | `GetPageProperties` |
| 29 | `POST /rest/api/content/{id}/property` | 889 | ❌ нет | — |
| 30 | `GET /rest/api/content/{id}/property/{key}` | 929 | ❌ нет | — |
| 31 | `POST /rest/api/content/{id}/property/{key}` | 959 | ❌ нет | — |
| 32 | `PUT /rest/api/content/{id}/property/{key}` | 999 | ❌ нет | — |
| 33 | `DELETE /rest/api/content/{id}/property/{key}` | 1047 | ❌ нет | — |
| 34 | `GET /rest/api/search` | 1062 | ❌ нет | — (есть только `/content/search`) |
| 35 | `GET /rest/api/space` | 1098 | ✅ | `GetSpaces` |
| 36 | `GET /rest/api/space/{spaceKey}` | 1133 | ❌ нет | — |
| 37 | `GET /rest/api/space/{spaceKey}/content` | 1163 | ⚠️ через #1 | `GetSpacePageTree` использует `/content?spaceKey=`, а не `/space/{key}/content` |
| 38 | `GET /rest/api/space/{spaceKey}/content/{type}` | 1196 | ❌ нет | — |
| 39 | `GET /rest/api/space/{spaceKey}/property` | 1229 | ❌ нет | — |
| 40 | `GET /rest/api/space/{spaceKey}/property/{key}` | 1261 | ❌ нет | — |
| 41 | `GET /rest/api/user/current` | 1291 | ❌ нет | — |
| 42 | `GET /rest/api/user` | 1321 | ⚠️ обходной | `SearchUser` идёт через `/group/{g}/member`, а не `/user` |
| 43 | `GET /rest/api/user/list` | 1356 | ❌ нет | — |

> Итого: полностью реализовано ~17, частично/не точно ~6, отсутствует ~20.

## Этапы работы

### Этап 0. Подготовка и инфраструктура
1. Зафиксировать карту разделов `api_docs.md` (H2 → строка) скриптом `grep -nE '^## '` — уже сделано выше.
2. Написать вспомогательный скрипт/функцию извлечения раздела метода по заголовку
   (варианты A/B/C из «Как вырезать часть документа»), чтобы при проверке каждого метода
   подтягивать его точное описание (query-параметры, тело запроса, ответы) без искажений.
3. Уточнить scope: какие из отсутствующих методов реально нужны (см. вопрос ниже).
   Кандидаты на пропуск (Cloud-only / редкие): blueprint instance (#8–9), macro id/hash (#11–12),
   descendant (#22–23), space property (#39–40), user list (#43). Решить по каждому.

### Этап 1. Аудит существующих методов (read-only)
Для каждого реализованного метода (#1–5, 7, 13, 15–19, 24–25, 28, 35, 42):
1. Вырезать описание из `api_docs.md` скриптом.
2. Сверить с кодом: путь, HTTP-метод, query-параметры, тело запроса, обработка ответов.
3. Зафиксировать расхождения в чек-листе (таблица: метод → параметр → ожидание из дока → что в коде → статус).
4. Особое внимание:
   - **#10 history** — сверить, корректно ли используется `status=historical&version` вместо `/history`.
   - **#16 child/{type}** — проверить, нужен ли общий тип кроме `page`/`comment`/`attachment`.
   - **#37 space content** — решить, оставить `/content?spaceKey=` или добавить `/space/{key}/content`.
   - **#42 user** — подтвердить, что DC-обход через group member приемлем.

### Этап 2. Рефакторинг существующего кода
По итогам аудита:
1. Вынести повторяющийся паттерн «NewRequest + query + Do + decode» в маленький хелпер
   (если это сократит объём без потери читаемости).
2. Привести имена/подписи методов к единому стилю (сейчас смесь `GetPage*`, `GetContent*`).
3. Убрать дублирование: например `GetPageComments` и `GetInlineComments` ходят на один и тот же
   `/child/comment` с разными expand — объединить в один метод с параметром.
4. Проверить обработку ошибок: везде `return nil, err` с контекстом `confluence: <операция>: %w`.
5. Не менять публичный контракт MCP-инструментов (имена, параметры) без необходимости —
   это ломает интеграционные тесты.

### Этап 3. Добавление отсутствующих методов (по scope)
Для каждого выбранного к добавлению метода:
1. Вырезать описание из `api_docs.md`.
2. Реализовать метод в `fetcher_ext.go` (или новый файл по области: `content.go`, `label.go`,
   `property.go`, `space.go`, `user.go` — если решим разбить большой `fetcher_ext.go`).
3. При необходимости добавить MCP-инструмент в `tools_conf_read.go` / `tools_conf_write.go`.
4. Добавить интеграционный тест в `test/integration/all_tools_test.go` (mock httptest).
5. Пропустить write-методы, если scope строго read-only (уточнить).

Кандидаты (при положительном решении):
- `GET /content/scan` (#6)
- `GET /content/{id}/history` (#10) — переделать текущий
- `GET /content/{id}/restriction/byOperation/{operationKey}` (#14)
- `GET /content/{id}/child/{type}` (#16) — общий
- `PUT /content/{id}/child/attachment/{attachmentId}` (#20)
- `POST /content/{id}/child/attachment/{attachmentId}/data` (#21)
- `GET /content/{id}/descendant[/{type}]` (#22–23)
- `DELETE /content/{id}/label[/{label}]` (#26–27)
- `POST/GET/PUT/DELETE /content/{id}/property[/{key}]` (#29–33)
- `GET /search` (#34)
- `GET /space/{spaceKey}` (#36)
- `GET /space/{spaceKey}/content[/{type}]` (#37–38)
- `GET /user/current` (#41)
- `GET /user`, `GET /user/list` (#42–43)

### Этап 4. Тесты и верификация
1. `go build ./...` — сборка.
2. `go vet ./...` — статический анализ.
3. `go test ./... -count=1` — все 138 тестов (добавленные новые тоже).
4. `go test ./test/integration -v` — интеграционные (MCP SSE + mock).
5. Для каждого нового/изменённого метода — отдельный тест с `httptest.NewServer`,
   проверяющий: путь, метод, query/body, разбор ответа.

### Этап 5. Документация и отчёт
1. Обновить `TOOLS.md` / `STATUS.md`, если добавлены инструменты.
2. Свести итоговый отчёт: таблица «метод → статус (было/стало) → тест», расхождения,
   которые были исправлены, и методы, сознательно пропущенные (с причиной).

## Принципы рефакторинга
- Сохранять имена и параметры MCP-инструментов (не ломать интеграционные тесты и клиентов).
- Только PAT-аутентификация, Jira/Confluence DC (v1/v2 API), без Cloud-only фич.
- Read-only guard: каждый write-инструмент проходит `guardWrite(app)`.
- Ошибки — через `return nil, err` с контекстом.
- Минимальные изменения: не переписывать то, что работает и соответствует докам.
- Каждый новый метод покрыт тестом.

## Открытые вопросы (решить до старта)
1. **Scope**: добавлять ВСЕ отсутствующие методы или только те, что реально нужны?
   (blueprint, macro, descendant, space property, user list — вероятно, пропустить.)
2. **Read-only vs write**: в `list.md` есть write-методы (POST/PUT/DELETE label, property,
   attachment data). Добавлять write-инструменты или только read?
3. **Разбиение файлов**: оставить всё в `fetcher_ext.go` (уже 1180 строк) или разбить
   по областям (`content.go`, `label.go`, `property.go`, `space.go`, `user.go`)?
4. **История (#10)**: переделать `GetPageHistory` на настоящий `/history`-эндпоинт или
   оставить текущий рабочий подход (`status=historical`)?