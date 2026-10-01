# Сравнение спецификаций: docs/9.0.0.swagger.v3.json ↔ docs/api_docs.md

Дата сравнения: 2026-10-01. Confluence REST API (v1, `/rest/api`).

**Объём:** api_docs.md описывает 43 метода, swagger v3 — 111 операций (85 путей).
Все 43 метода из api_docs.md присутствуют в swagger; **68 методов есть только в swagger** (полный список в конце).

Каждый из 43 общих методов сравнивался отдельным агентом по категориям:
описание / параметры / тело запроса / ответы.

---

## 1. Системные отличия (повторяются почти во всех методах)

| Категория | api_docs.md | swagger v3 |
|---|---|---|
| Path-параметры (`id`, `key`, `spaceKey`, `type`…) | **не документируются вообще** (секций «Path-параметры» нет ни в одном разделе) | описаны явно (string, `required: true`, с описаниями) |
| Типы `start`/`limit` | `int` + default (25/50/100/200) | **`type: string`, без default** (сквозной артефакт конвертации) |
| Схемы ответов-списков | голый **array** (Swagger-2 стиль, `#/definitions/*`, `docs.atlassian.com/jira/REST/schema/...`) | **объект-обёртка** пагинации: `pageRequest`, `nextCursor`/`prevCursor`, `results`, `start`, `limit`, `next`/`self`/`prev` |
| Вложенные объекты (`space`, `history`, `version`, `container`, `contributors`) | **массивы** `$ref` | одиночные объекты `$ref` + доп. поля `links`, `position`, `historyRef`/`spaceRef`/`containerRef`/`versionRef` |
| Схемы ошибок 4xx | только текст | `application/json` + `$ref RestError` (`statusCode`, `data: ValidationResult`, `message`, `reason`) |
| Метаданные | Javadoc-артефакты `{@link ...}`, `{@see ...}` | `operationId`, `tags`, примеры URI («Example request URI(s)») |
| Default-значения параметров | в колонке «По умолчанию» таблицы | только текстом внутри description, формального `default` в схеме нет |
| Content-type тел запросов | не указан | `application/json` / `multipart/form-data` |
| required полей схем | не указано | как правило тоже нет (`required` отсутствует) |
| Enum | нигде формально не задан | тоже (допустимые значения только текстом) |

---

## 2. Отличия по методам (специфичные, сверх системных)

### Content

#### GET `/rest/api/content`
- Отличий сверх системных нет (набор query-параметров идентичен: type, spaceKey, title, status, postingDay, expand, start, limit).

#### POST `/rest/api/content`
- Код ответа в md — плейсхолдер `**?**` (в swagger: 200 + 404).
- Default `expand`: в таблице md — `body.storage,history,space,container.history,container.version,version,ancestors`, в тексте того же md и в swagger — `history,space,version` (три противоречащих варианта).
- Default `status`: `current` в md-таблице, `[current]` текстом в swagger.
- Дефект swagger: описание requestBody «new content to be created» у метода обновления/создания без required.

#### GET `/rest/api/content/{id}`
- `version`: **int (md) vs string (swagger)**.
- Summary: «Get content by id» vs «Get content by ID».

#### PUT `/rest/api/content/{contentId}`
- md содержит «Updating a draft is not currently supported.» — в swagger фразы нет.
- `conflictPolicy`: default `abort` есть только в md-таблице; в swagger — только в тексте (с незакрытым тегом `<code>`).
- В swagger к описанию приложены 3 примера JSON-тел.

#### DELETE `/rest/api/content/{id}`
- Только системные (в swagger у 200 пустая схема, а не Content).

#### GET `/rest/api/content/scan`
- Описание `status=any`: md включает `draft`, swagger — нет (+ добавляет «Does not support 'historical' status for now»).
- Ответ в md — `**?**`; в swagger 200 (пагинационная обёртка) + 404.

#### GET `/rest/api/content/search`
- **Код ошибки при невалидном/отсутствующем CQL: 400 (md) vs 404 (swagger)**.
- Описание в md обрезано: нет ссылки на CQL-документацию и примеров URI.
- `expand` в swagger дополнен «Default value: history,space,version».

#### POST `/rest/api/content/blueprint/instance/{draftId}`
- Ответ в md — `**?**`; в swagger 200 + схема Content.
- Defaults (`status=draft`, `expand=body.storage,history,space,version,ancestors`) есть только в md.

#### PUT `/rest/api/content/blueprint/instance/{draftId}`
- Аналогично POST: ответ `**?**` vs 200; defaults только в md.
- Дефект swagger: описание requestBody начинается с «he content…» (опечатка).

### History / Macro

#### GET `/rest/api/content/{id}/history`
- `previousVersion`/`nextVersion`/`lastUpdated`: **массивы полных объектов version (md) vs Reference-объекты `{expanded, idProperties}` (swagger)**.
- `contributors`: массив детализированных объектов (md) vs `ReferenceContributors` (swagger).
- `createdBy`: развёрнутый `person` (md) vs компактный `Person {profilePicture, displayName, type}` (swagger).
- `createdDate` в swagger имеет `format: date-time` + example.
- Swagger в примерах URI упоминает `cql`, `cqlcontext`, `start`, `limit`, но **не объявляет их как параметры**.

#### GET `/rest/api/content/{id}/history/{version}/macro/id/{macroId}`
- В md нет секции параметров вовсе (даже query).
- Дефект swagger: описание `version` скопировано с hash-варианта («the version of the content which the **hash** belongs»).
- Схема `parameters` макроса в swagger утеряна: нет структуры `value`, `MacroParameterInstance` объявлен, но не референсится; в md схема полная (patternProperties + additionalProperties: false).

#### GET `/rest/api/content/{id}/history/{version}/macro/hash/{hash}`
- **`deprecated: true` есть только в swagger** (у sibling-метода по macroId флага нет).
- В md остался Javadoc-артефакт `{@link #getContentById(...)}`.
- Те же потери схемы `parameters`, что и у macro/id.

### Restrictions

#### GET `/rest/api/content/{id}/restriction/byOperation`
- Default `expand`: полный список (md) vs `group` (swagger).
- Схема 200 в swagger — `MockRestrictionsResponse`: обёртка со свойством `restrictions` + `get_links` (вместо `_links`), не формализована (только examples). В md — прямая мапа operation → content-restriction.

#### GET `/rest/api/content/{id}/restriction/byOperation/{operationKey}`
- `limit`: default **100** только в md.
- Схема: md имеет свойство `content`; swagger — `get_links`/`get_expandable` (вместо `_links`/`_expandable`) и `restrictions` как объект с PageResponseObject (`user`/`group`), а не массив subject-restrictions.
- `operation` в md — $ref на operation-key, в swagger — просто string.
- Дефект swagger: опечатка «Aa comma separated list…» в описании expand.

### Child / Attachment

#### GET `/rest/api/content/{id}/child`
- `parentVersion`: int в обоих, default 0 в обоих, но описание есть только в md.
- Схема md — map «тип контента → массив» (`patternProperties`), swagger — generic-пагинация (map не отражена).

#### GET `/rest/api/content/{id}/child/{type}`
- Только системные (в т.ч. массив vs пагинационная обёртка).

#### GET `/rest/api/content/{id}/child/comment`
- В md длиннее описание `expand` (упоминание extensions.inlineProperties / extensions.resolution).

#### GET `/rest/api/content/{id}/child/attachment`
- **Код успеха: 200 (md) vs 201 (swagger)** — вероятная ошибка swagger.
- Default `limit=50` только в md.
- Дефект swagger: описание `mediaType` — копипаст описания `expand`.

#### POST `/rest/api/content/{id}/child/attachment`
- `allowDuplicated`: **boolean + default false (md) vs string без default (swagger)**.
- Multipart-тело (`file`, `comment`, **`minorEdit`**, **`hidden`** — `MockAttachmentRequest`) есть только в swagger; в md секции «Тело запроса» нет вообще.
- Схема 200: массив Content (md) vs одиночный объект Content (swagger).
- В md описание дополнено примерами загрузки файла.

#### PUT `/rest/api/content/{id}/child/attachment/{attachmentId}`
- В md нет секции параметров вовсе (нет даже query).
- В swagger есть тексты описаний 200 и ошибок; в обоих коды 200/400/403/404/409 совпадают.

#### POST `/rest/api/content/{id}/child/attachment/{attachmentId}/data`
- Multipart-тело только в swagger (та же `MockAttachmentRequest`).
- Текст описания в md обрывается на «...with no comment:» (утеряны curl-примеры).
- В md нет секции path-параметров и тела.

### Descendant

#### GET `/rest/api/content/{id}/descendant`
- Схема swagger (плоская пагинация) **противоречит собственному описанию** «map keyed by content type» (как в md).

#### GET `/rest/api/content/{id}/descendant/{type}`
- **Схема 200 в swagger — одиночный объект Content ($ref), а не массив**.
- Default `limit=25` в md-таблице противоречит собственному тексту md «(optional, default: site limit)».
- В swagger описание опускает фразу «limited to a single descendant type».

### Label

#### GET `/rest/api/content/{id}/label`
- У Label в swagger лишнее поле `label` (string); в md — {prefix, name, id} + additionalProperties: false.
- Default `limit=200` только в md.
- В swagger есть неиспользуемый компонент `PageResponseLabel`.

#### POST `/rest/api/content/{id}/label`
- **Схема тела в md — просто `string` (явная ошибка)**; в swagger — объект Label; фактический формат (см. docs/conf_api_audit.md) — массив Label, не описанный ни там, ни там.
- Описание 200 в swagger существенно длиннее (existing and added labels, empty list).

#### DELETE `/rest/api/content/{id}/label`
- Дефект swagger: мусорные хвосты в текстах 403/404 («…permission to the content.permission to the content.»).
- Формулировка 204: «Empty response…» (md) vs «returns no response…» (swagger).

#### DELETE `/rest/api/content/{id}/label/{label}`
- Дефект swagger: лишняя фраза «The body is the json representation of the list.» (копипаст из POST) + хвосты в 400/403.

### Property (content)

#### GET `/rest/api/content/{id}/property`
- Дефект swagger: описание `limit` — копипаст «the limit of the number of **labels** to return».
- Default `expand=version.` с висячей точкой (в обоих, по-разному оформлен).
- Типы полей JsonContentProperty: `id` object (md) vs string (swagger); `value` произвольный object (md) vs `JsonString{value: string}` (swagger); `version`/`content` массивы (md) vs объекты (swagger); `contentRef` только в swagger.

#### POST `/rest/api/content/{id}/property`
- Те же типовые расхождения JsonContentProperty (id/value/version/content/contentRef).
- `contentRef` ($ref ReferenceContent) есть только в swagger.

#### GET `/rest/api/content/{id}/property/{key}`
- В swagger есть сомнительный query-параметр `limit` (описание про «labels», для единичного свойства) — в md его нет.

#### POST `/rest/api/content/{id}/property/{key}`
- **Swagger обеднён**: нет summary/description, нет описания у path-параметра `key`, вместо кодов 200/400/403/413 — единственный `default`-ответ.

#### PUT `/rest/api/content/{id}/property/{key}`
- Query-параметр `expand` есть **только в swagger**.
- Дефект md: описание 200 говорит про «a piece of content» вместо content property.
- Дефект swagger: опечатка «eturned» в 404.

#### DELETE `/rest/api/content/{id}/property/{key}`
- Только системные (204/404, тексты совпадают).

### Search

#### GET `/rest/api/search`
- В md обрезано описание (нет ссылки на CQL и примеров).
- **В swagger ошибочно объявлен requestBody (application/json, type: string) у GET-метода**.
- Типы параметров: `start`/`limit` int и `includeArchivedSpaces` boolean (md) vs все string (swagger); defaults (`excerpt=highlight`, `limit=25`, `includeArchivedSpaces=false`) только в md.
- SearchResult: в swagger на 5 полей больше — `entityRef`, `resultParentRef`, `resultGlobalContainerRef`, `resourceType`, `entity`, `entityType`; контейнеры — одиночные объекты, не массивы; `lastModified` с format: date-time.

### Space

#### GET `/rest/api/space`
- В swagger лишний параметр `spaceKeySingle` (string, без описания) — артефакт repeatable `spaceKey`.
- `favourite`/`hasRetentionPolicy`: **boolean (md) vs string (swagger)**.
- **В swagger ошибочно объявлен requestBody у GET-метода** (application/json, type: string).
- **Схема 200: массив Space (md) vs одиночный объект Space (swagger)** — ни в одном нет пагинационной обёртки, хотя md называет схему «Page Response of Space».
- Дефект swagger: описание `spaceKey` — копипаст «the key of the space to update».

#### GET `/rest/api/space/{spaceKey}`
- Space в swagger дополнительно содержит `status` и `links`; `icon`/`homepage`/`retentionPolicy` — массивы (md) vs Reference-объекты (swagger); `description` — структурированная (md) vs безликий object (swagger); `id` int64 + example (swagger).

#### GET `/rest/api/space/{spaceKey}/content`
- Дефекты swagger-копипастов: описание `expand` — «on the space» (вместо «on each piece of content»), `spaceKey` — «to update», опечатка «he start point».
- Default `depth=all` только в md.

#### GET `/rest/api/space/{spaceKey}/content/{type}`
- Допустимые значения `type` (page, blogpost) указаны **только в swagger** (в path-параметре).
- Default `depth=all`, `limit=25` только в md.

#### GET `/rest/api/space/{spaceKey}/property`
- JsonSpaceProperty: `value` object (md) vs `JsonString` (swagger); `version`/`space` массивы (md) vs объекты (swagger); `spaceRef` только в swagger.
- В swagger есть неиспользуемый компонент `PageResponseJsonSpaceProperty`.
- Дефект swagger: опечатка «he start point» у `start`.

#### GET `/rest/api/space/{spaceKey}/property/{key}`
- Описание в md неверное: «Returns a paginated list of space properties» для метода, возвращающего одно свойство (копипаст из GET /property).
- Query-параметры `start`/`limit` есть только в swagger.
- Текст 404 в swagger дополнен «or no property with the given key».
- У path-параметра `key` в swagger нет описания.

### User

#### GET `/rest/api/user/current`
- Описание `expand` в md пустое («—»).
- Person в swagger урезан до 3 полей (`profilePicture`, `displayName`, `type`); в md — полный anyOf (anonymous/known-user/unknown-user/user) с `username`, `status`, Icon с required [width, height, isDefault].

#### GET `/rest/api/user`
- Дефект swagger: описание `username` — копипаст из watcher-эндпоинта («userName of the user to create the new watcher for»).
- Та же урезанная схема Person в swagger vs полный anyOf в md.

#### GET `/rest/api/user/list`
- Summary: «Get users» (md) vs «Get registered users» (swagger).
- **Внутреннее противоречие md**: default `limit` = 100 в таблице, но «Default value is 200» в тексте описания.
- Defaults (`start=0`, `limit=100`) только в md.

---

## 3. Дефекты каждого из источников (сводно)

### Явные дефекты swagger (9.0.0.swagger.v3.json)
1. `start`/`limit` (и boolean-параметры) везде типизированы как string, без default.
2. GET `/content/{id}/child/attachment` возвращает **201** вместо 200.
3. requestBody у GET-методов: `/search`, `/space`.
4. Копипаст-описания: `mediaType` (attachments), `username` (user, из watcher), `spaceKey` («to update» в 3 местах), `expand` («on the space»), `limit` («labels» у content property), `version` у macro/id (из hash-варианта).
5. Мусорные хвосты в текстах ошибок label-методов («…permission to the content.permission to the content.»), опечатки («he start point», «Aa comma», «eturned», «he content»).
6. Схема 200 у `descendant/{type}` — одиночный объект вместо массива/мапы; у `/space` — одиночный Space вместо списка.
7. `POST /content/{id}/property/{key}` без summary/description/кодов ответов (один default).
8. Утеряна структура `parameters` макроса в MacroInstance (MacroParameterInstance не референсится).
9. Поля `get_links`/`get_expandable` вместо ожидаемых `_links`/`_expandable`.
10. Лишний параметр `spaceKeySingle` у GET /space.

### Явные дефекты api_docs.md

> **Проверено по исходнику (docs/docs.atlassian.com.html, официальная документация
> Atlassian ~8.x):** все перечисленные ниже «дефекты md» присутствуют и в самом
> HTML-оригинале — экстрактор `extract_api.py` воспроизвёл их корректно, не исказив
> источник. Это дефекты самой старой документации Atlassian, часть которых исправлена
> в swagger 9.0 (массивы→объекты, `**?**`→200, path-параметры и RestError добавлены).

1. Path-параметры не документированы ни для одного метода (в HTML нет ни одного `<h6>path parameters</h6>`).
2. Тело POST `/content/{id}/label` описано как `string` (в HTML `"type":"string"`).
3. Код ответа `**?**` у 4 методов (POST /content, GET /content/scan, POST/PUT blueprint) — статус в HTML реально отсутствует у этих representation.
4. Отсутствуют схемы ошибок (RestError) и content-type тел (content-type только в ответах).
5. Противоречия внутри файла: default expand у POST /content (3 варианта), limit у user/list (100 vs 200), limit=25 vs «site limit» у descendant/{type}.
6. Устаревшие Jira-style схемы: массивы вместо объектов для space/history/version/container; ссылки на `docs.atlassian.com/jira/REST/schema/...`.
7. Javadoc-артефакты `{@link ...}`/`{@see ...}` (34 вхождения в HTML); обрыв текста на «...with no comment:» у attachment data — тоже в оригинале.
8. Неверное описание 200 у PUT property/{key} («content» вместо property) и у GET space property/{key} («list» вместо одного свойства) — копипасты из HTML.
9. macro/id: секция параметров отсутствует и в HTML (нет query/path таблиц и тела).
10. Обрезанные описания: search заканчивается на «For example:», user/current expand без описания — так и в HTML.

---

## 4. Главные практические выводы

1. **Схемы ответов несовместимы по корню**: md описывает голые массивы, swagger — пагинационные обёртки с `results`. Для клиента Confluence 9.0 ближе swagger, но его generic-обёртки местами противоречат собственным описаниям методов (descendant, space content).
2. Для кодогенерации/валидации брать за основу swagger (path-параметры, RestError, content-type), но типы `start`/`limit`/boolean-параметров править на int/boolean вручную и проверять спорные схемы живыми запросами.
3. Критичные для клиента расхождения, требующие проверки на живом инстансе: код ошибки CQL (400 vs 404), код успеха attachments (200 vs 201), форма ответов-списков (массив vs обёртка), состав полей Person/Label/Content.

---

## 5. 68 методов, присутствующих только в swagger (нет в api_docs.md)

> **Обновлено:** все перечисленные ниже swagger-only методы удалены из
> `docs/9.0.0.swagger.v3.json` (2026-10-01). В swagger осталось ровно 43 операции —
> множество методов из api_docs.md. Раздел сохранён как исторический список удалённых.

**admin (7):**
- POST /rest/api/admin/group — create
- DELETE /rest/api/admin/group/{groupName}
- POST /rest/api/admin/user — createUser
- DELETE /rest/api/admin/user/{username}
- POST /rest/api/admin/user/{username}/password — changePassword
- PUT /rest/api/admin/user/{username}/disable
- PUT /rest/api/admin/user/{username}/enable

**audit / accessmode (2):**
- GET /rest/api/audit — getAuditRecords
- GET /rest/api/accessmode — getAccessModeStatus

**backup-restore (12):**
- POST /rest/api/backup-restore/backup/site — createSiteBackupJob
- POST /rest/api/backup-restore/backup/space — createSpaceBackupJob
- POST /rest/api/backup-restore/restore/site — createSiteRestoreJob
- POST /rest/api/backup-restore/restore/site/upload
- POST /rest/api/backup-restore/restore/space — createSpaceRestoreJob
- POST /rest/api/backup-restore/restore/space/upload
- GET /rest/api/backup-restore/jobs — findJobs
- GET /rest/api/backup-restore/jobs/{jobId} — getJob
- GET /rest/api/backup-restore/jobs/{jobId}/download — downloadBackupFile
- GET /rest/api/backup-restore/restore/files — getFiles
- PUT /rest/api/backup-restore/jobs/clear-queue — cancelAllQueuedJobs
- PUT /rest/api/backup-restore/jobs/{jobId}/cancel — cancelJob

**content (прочее) (7):**
- POST /rest/api/contentbody/convert/{to} — convert
- GET /rest/api/content/{contentId}/watchers — index
- PUT /rest/api/content/{id}/restriction — updateRestrictions
- DELETE /rest/api/content/{id}/version/{versionNumber} — deleteContentHistory
- POST /rest/api/content/{id}/child/attachment/{attachmentId}/move — move
- DELETE /rest/api/content/{id}/child/attachment/{attachmentId} — removeAttachment
- DELETE /rest/api/content/{id}/child/attachment/{attachmentId}/version/{version} — removeAttachmentVersion

**group (3):**
- GET /rest/api/group — getGroups
- GET /rest/api/group/{groupName} — getGroup
- GET /rest/api/group/{groupName}/member — getMembers

**longtask (2):**
- GET /rest/api/longtask — getTasks
- GET /rest/api/longtask/{id} — getTask

**space (20):**
- POST /rest/api/space — createSpace
- POST /rest/api/space/_private — createPrivateSpace
- PUT /rest/api/space/{spaceKey} — update
- DELETE /rest/api/space/{spaceKey}
- PUT /rest/api/space/{spaceKey}/archive
- PUT /rest/api/space/{spaceKey}/restore
- GET /rest/api/space/{spaceKey}/watchers
- GET /rest/api/space/{spaceKey}/labels
- GET /rest/api/space/{spaceKey}/labels/popular
- GET /rest/api/space/{spaceKey}/labels/recent
- GET /rest/api/space/{spaceKey}/labels/{labelName}/related
- POST /rest/api/space/{spaceKey}/property — create
- POST /rest/api/space/{spaceKey}/property/{key} — create
- PUT /rest/api/space/{spaceKey}/property/{key} — update
- DELETE /rest/api/space/{spaceKey}/property/{key}

**user (11):**
- GET /rest/api/user/anonymous — getAnonymous
- GET /rest/api/user/memberof — getGroups
- POST /rest/api/user/current/password — changePassword
- GET /rest/api/user/watch/content/{contentId} — isWatchingContent
- POST /rest/api/user/watch/content/{contentId} — addContentWatcher
- DELETE /rest/api/user/watch/content/{contentId} — removeContentWatcher
- GET /rest/api/user/watch/space/{spaceKey} — isWatchingSpace
- POST /rest/api/user/watch/space/{spaceKey} — addSpaceWatch
- DELETE /rest/api/user/watch/space/{spaceKey} — removeSpaceWatch
- PUT /rest/api/user/{username}/group/{groupName} — update
- DELETE /rest/api/user/{username}/group/{groupName}

**webhooks (9):**
- GET /rest/api/webhooks — findWebhooks
- POST /rest/api/webhooks — createWebhook
- GET /rest/api/webhooks/{webhookId} — getWebhook
- PUT /rest/api/webhooks/{webhookId} — updateWebhook
- DELETE /rest/api/webhooks/{webhookId} — deleteWebhook
- GET /rest/api/webhooks/{webhookId}/latest — getLatestInvocation
- GET /rest/api/webhooks/{webhookId}/statistics — getStatistics
- GET /rest/api/webhooks/{webhookId}/statistics/summary — getStatisticsSummary
- POST /rest/api/webhooks/test — testWebhook
