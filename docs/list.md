### 1. Область `api/content`
- `GET` `/rest/api/content`
- `POST` `/rest/api/content`
- `GET` `/rest/api/content/{id}`
- `PUT` `/rest/api/content/{id}` *(в документации также встречается как `{contentId}`)*
- `DELETE` `/rest/api/content/{id}`
- `GET` `/rest/api/content/scan`
- `GET` `/rest/api/content/search`
- `POST` `/rest/api/content/blueprint/instance/{draftId}`
- `PUT` `/rest/api/content/blueprint/instance/{draftId}`
- `GET` `/rest/api/content/{id}/history`
- `GET` `/rest/api/content/{id}/history/{version}/macro/id/{macroId}`
- `GET` `/rest/api/content/{id}/history/{version}/macro/hash/{hash}`
- `GET` `/rest/api/content/{id}/restriction/byOperation`
- `GET` `/rest/api/content/{id}/restriction/byOperation/{operationKey}`

### 2. Область `api/content/{id}/child`
- `GET` `/rest/api/content/{id}/child`
- `GET` `/rest/api/content/{id}/child/{type}`
- `GET` `/rest/api/content/{id}/child/comment`

### 3. Область `api/content/{id}/child/attachment`
- `GET` `/rest/api/content/{id}/child/attachment`
- `POST` `/rest/api/content/{id}/child/attachment`
- `PUT` `/rest/api/content/{id}/child/attachment/{attachmentId}` *(обновление метаданных)*
- `POST` `/rest/api/content/{id}/child/attachment/{attachmentId}/data` *(обновление бинарных данных файла)*
> **Важное примечание:** Эндпоинта `DELETE /rest/api/content/{id}/child/attachment/{attachmentId}` в документации **нет**. Удаление вложения в Confluence выполняется через общий эндпоинт `DELETE /rest/api/content/{id}`, где `{id}` — это ID самого вложения (так как вложение является типом контента).

### 4. Область `api/content/{id}/descendant`
- `GET` `/rest/api/content/{id}/descendant`
- `GET` `/rest/api/content/{id}/descendant/{type}`

### 5. Область `api/content/{id}/label`
- `GET` `/rest/api/content/{id}/label`
- `POST` `/rest/api/content/{id}/label`
- `DELETE` `/rest/api/content/{id}/label` *(с использованием query-параметра `?name=`)*
- `DELETE` `/rest/api/content/{id}/label/{label}`

### 6. Область `api/content/{id}/property`
- `GET` `/rest/api/content/{id}/property`
- `POST` `/rest/api/content/{id}/property`
- `GET` `/rest/api/content/{id}/property/{key}`
- `POST` `/rest/api/content/{id}/property/{key}`
- `PUT` `/rest/api/content/{id}/property/{key}`
- `DELETE` `/rest/api/content/{id}/property/{key}`

### 7. Область `api/search`
- `GET` `/rest/api/search`

### 8. Область `api/space` (строго GET, как вы просили)
- `GET` `/rest/api/space`
- `GET` `/rest/api/space/{spaceKey}`
- `GET` `/rest/api/space/{spaceKey}/content`
- `GET` `/rest/api/space/{spaceKey}/content/{type}`

### 9. Область `api/space/{spaceKey}/property` (строго GET, как вы просили)
- `GET` `/rest/api/space/{spaceKey}/property`
- `GET` `/rest/api/space/{spaceKey}/property/{key}`

### 10. Пользовательские эндпоинты (строго указанные вами GET)
- `GET` `/rest/api/user/current`
- `GET` `/rest/api/user`
- `GET` `/rest/api/user/list`

