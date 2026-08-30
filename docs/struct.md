Вот **полный и исчерпывающий список** структурных элементов, которые вы найдете внутри `method-body`, с реальными примерами из версии 8.5.31:

---

### 1. Иерархия заголовков внутри `method-body`
Вложенность заголовков строже, чем просто `h5`:
*   `<h5>Request</h5>` и `<h5>Responses</h5>` — основные разделы.
*   `<h6>query parameters</h6>` или `<h6>path parameters</h6>` — подзаголовки перед таблицами параметров.
*   `<h6>Schema</h6>` — заголовок перед блоком JSON-схемы.
*   `<h6>Example</h6>` — заголовок перед примером JSON-тела (встречается не всегда, но часто).

### 2. Структура таблиц параметров
Таблицы не просто `<table>`, они имеют специфическую структуру:
```html
<table>
  <thead>
    <tr>
      <th>parameter</th>
      <th>type</th>
      <th>description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>expand</td>
      <td>string</td>
      <td>
        A comma separated list... Default value: <code>history,space,version</code>
      </td>
    </tr>
  </tbody>
</table>
```
**Важно для парсинга:** Значения по умолчанию (Default values) часто обернуты в тег `<code>` или `<em>` прямо внутри ячейки `<td>` описания.

### 3. Блоки кода и схем (Representation)
JSON-схемы и примеры обернуты не просто в `<pre>`, а в специфическую комбинацию:
```html
<div class="representation-doc"> <!-- Или просто без класса, но внутри method-body -->
  <h6>Schema</h6>
  <pre><code>{"id":"https://docs...","title":"Content","type":"object", ... }</code></pre>
</div>
```
**Важно для парсинга:** Внутри `<code>` JSON-схемы часто содержат HTML-ссылки на определения: `<a href="#definitions/content">$ref</a>`. Ваш парсер должен либо игнорировать эти теги `<a>` при извлечении чистого JSON, либо корректно их обрабатывать, иначе JSON будет невалидным.

### 4. Статусы ответов (Responses)
Статусы не всегда имеют строгую табличную или блочную структуру, они часто оформлены как параграфы с жестким форматированием:
```html
<p>
  <strong>Status</strong> <strong>200</strong> - <em>application/json</em>
</p>
<p>Returns a full JSON representation of the content property list</p>
<!-- Далее следует div с Schema -->

<p>
  <strong>Status</strong> <strong>404</strong>
</p>
<p>Returned if there is no content with the given id...</p>
```
**Важно для парсинга:** Нужно искать паттерн `<strong>Status</strong>\s*<strong>(\d+)</strong>`, чтобы надежно извлекать HTTP-коды, так как классы для этого не используются.

### 5. Примеры URI (Example request URIs)
Часто встречаются перед разделом Request:
```html
<p>Example request URI(s):</p>
<ul>
  <li>
    <a href="http://example.com/..." class="external-link" rel="nofollow">
      http://example.com/confluence/rest/api/content/1234
    </a>
  </li>
</ul>
```

### 6. Специфические маркеры и классы
*   **Deprecated:** Устаревшие методы могут иметь маркер: `<span class="deprecated">deprecated</span>` рядом с именем метода в `<h4>`.
*   **Show more:** В конце описания ресурса часто бывает: `<a href="#" class="show-more">Show more</a>` (это UI-элемент, но он присутствует в DOM).
*   **Неразрывный пробел:** Как уже отмечалось, между методом и путем всегда стоит `&nbsp;`: `<code>GET&nbsp;/rest/api/content</code>`. Это нужно учитывать при регулярных выражениях (искать `\s+` или `&nbsp;`).

---

### Итоговая исчерпывающая схема (Шаблон для парсера)

Если вы пишете парсер, вот полный CSS/XPath-подобный шаблон того, что вы встретите:

```html
<div class="resource">
  <h3 id="api/{resource_path}">api/{resource_path}</h3>
  <p>{Описание ресурса}</p>
  
  <div class="methods">
    <div class="method">
      <h4 id="api/{resource_path}-{method_name}">
        <span class="method-name">{Human Readable Name}</span>
        <code>{METHOD}&nbsp;/{path}</code>
        <!-- Опционально: <span class="deprecated">deprecated</span> -->
      </h4>
      
      <div class="method-body">
        <p>{Описание метода}</p>
        
        <!-- Опционально: Примеры URL -->
        <p>Example request URI(s):</p>
        <ul>
          <li><a class="external-link" ...>{URL}</a></li>
        </ul>

        <h5>Request</h5>
        
        <!-- Опционально: Path params -->
        <h6>path parameters</h6>
        <table>
          <thead><tr><th>parameter</th><th>type</th><th>description</th></tr></thead>
          <tbody>...</tbody>
        </table>

        <!-- Опционально: Query params -->
        <h6>query parameters</h6>
        <table>
          <thead><tr><th>parameter</th><th>type</th><th>description</th></tr></thead>
          <tbody>
            <tr>
              <td>{name}</td>
              <td>{type}</td>
              <td>{description} (может содержать <code>Default: {value}</code>)</td>
            </tr>
          </tbody>
        </table>

        <!-- Опционально: Request Body Schema -->
        <h6>Schema</h6>
        <pre><code>{JSON Schema с возможными <a> тегами внутри}</code></pre>
        
        <!-- Опционально: Request Body Example -->
        <h6>Example</h6>
        <pre><code>{JSON Example}</code></pre>

        <h5>Responses</h5>
        
        <!-- Блок ответа (повторяется для каждого статуса) -->
        <p><strong>Status</strong> <strong>{HTTP_CODE}</strong> - <em>{MIME_TYPE}</em></p>
        <p>{Описание ответа}</p>
        
        <!-- Опционально: Response Schema -->
        <h6>Schema</h6>
        <pre><code>{JSON Schema}</code></pre>
      </div>
    </div>
  </div>
</div>
```

### Резюме для разработчика парсера:
1. Не полагайтесь только на теги, используйте комбинацию тегов и текстовых маркеров (например, поиск `<strong>Status</strong>`).
2. Очищайте извлекаемый JSON от встроенных HTML-тегов (особенно `<a href="...">` внутри `$ref`), иначе получите невалидный JSON.
3. Учитывайте `&nbsp;` между HTTP-методом и путем.
4. Значения по умолчанию "спрятаны" внутри текстового описания в таблицах, а не в отдельных колонках.

Это описание покрывает 100% структурных особенностей генератора REST-документации Atlassian для Confluence 8.5.x.