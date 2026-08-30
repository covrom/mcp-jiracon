# Preprocessing — Markdown / ADF / Storage conversion

The `preprocessing/` package handles content-format conversions so that
Jira and Confluence can round-trip user-provided Markdown through their
respective storage formats.

## Modules

| Module | Direction | Purpose |
| --- | --- | --- |
| `preprocessing/base.py` | – | Shared HTML/Markdown helpers; user-mention + user-profile macro resolution; `<time datetime>` preservation; `<ac:image>` → `<img>` rewrite |
| `preprocessing/confluence.py` | Markdown → Confluence storage (XHTML) | Uses `md2conf` + task-list + table-layout post-processing |
| `preprocessing/jira.py` | Markdown ↔ Jira wiki markup | Smart links, mention normalization, language normalization, smart links |

## `BasePreprocessor` (`preprocessing/base.py`)

### Class

```python
class BasePreprocessor:
    def __init__(self, base_url: str = "") -> None: ...
    def process_html_content(
        self,
        html_content: str,
        space_key: str = "",
        confluence_client: ConfluenceClient | None = None,
        content_id: str = "",
        attachments: list | None = None,
    ) -> tuple[str, str]:
        # Returns (processed_html, processed_markdown).
```

### Algorithm

1. Parse `html_content` with BeautifulSoup (`html.parser`).
2. Process user mentions and user profile macros.
3. For Confluence `<time>` elements, copy the `datetime` attribute into
   the element body so markdownify doesn't drop it.
4. Convert `<ac:image>` to `<img>` with cross-page awareness.
5. `str(soup)` → `md(...)` (markdownify) → `(processed_html, processed_markdown)`.

### User mention resolution

Priority: `ri:account-id` (Cloud) → `ri:userkey` (DC) → `ri:username` (DC).
When resolution succeeds → `@DisplayName`. When it fails →
`[User Profile: <id>]` placeholder.

### Image resolution

For `<ri:attachment>`:
1. Look up by filename in the provided `attachments` list.
2. Fallback: `f"{confluence_url}/download/attachments/{content_id}/{filename}"`.
3. External `ri:url` → preserved.
4. Unknown source → `[unsupported image]`.

Preserves `ac:width`, `ac:height`, `ac:alt` attributes.

## `ConfluencePreprocessor` (`preprocessing/confluence.py`)

### Class

```python
class ConfluencePreprocessor(BasePreprocessor):
    def __init__(self, base_url: str) -> None: ...
    def markdown_to_confluence_storage(
        self,
        markdown_content: str,
        *,
        enable_heading_anchors: bool = False,
        apply_task_lists: bool = True,
        table_layout: str | None = None,
    ) -> str: ...
```

### Table width presets

| `table_layout` value | Pixels |
| --- | --- |
| `full-width` | 1800 |
| `wide` | 960 |
| `default` | 760 |

The chosen width is injected into every `<table>` via `data-table-width`
and `data-layout` attributes.

### Conversion pipeline

1. `markdown_to_html(...)` via `md2conf.api`.
2. `_fix_attachment_images` (round 1) on the HTML.
3. If `apply_task_lists=False`, protect `<li>[ ]/[x]` markers with
   private-use placeholders before md2conf consumes them; restore after.
4. Set up `ConfluenceStorageFormatConverter` with `force_valid_url=False`,
   `heading_anchors=…`, `render_mermaid=False`, empty
   `ConfluencePageCollection` / `ConfluenceUserCollection`.
5. `converter.visit(root)`; convert back to string.
6. `_fix_attachment_images` (round 2), `_restore_task_list_markers`,
   `_normalize_task_list_bodies`.
7. If `apply_task_lists=True`: `_apply_task_lists`, then optionally
   `_apply_table_layout`.

### Task list rewriting

Finds top-level `<ul>` whose every direct `<li>` starts with
`^\s*[ ]` or `^\s*[x]`. Rewrites to:

```xml
<ac:task-list>
  <ac:task>
    <ac:task-status>incomplete</ac:task-status>
    <ac:task-body>...</ac:task-body>
  </ac:task>
  ...
</ac:task-list>
```

Nested or mixed lists are left untouched.

### Attachment image fixing

Turns `<img src="bare.ext">` into:

```xml
<ac:image ac:alt="...">
  <ri:attachment ri:filename="bare.ext" />
</ac:image>
```

Leaves `http://`, `https://`, `data:`, `/...`, `#...` URLs alone.

## `JiraPreprocessor` (`preprocessing/jira.py`)

### Class

```python
class JiraPreprocessor(BasePreprocessor):
    VALID_JIRA_LANGUAGES: set[str]   # 47 entries + aliases
    LANGUAGE_MAPPING: dict[str, str] # 12 mappings

    def __init__(
        self,
        base_url: str = "",
        disable_translation: bool = False,
    ) -> None: ...
    def clean_jira_text(self, text: str) -> str: ...
    def jira_to_markdown(self, input_text: str) -> str: ...
    def markdown_to_jira(self, input_text: str) -> str: ...
```

### `clean_jira_text`

- `[~accountid:<id>]` → `User:<id>` (placeholder; real name lookup is
  injected by callers when they have a Jira client).
- `[text|url|smart-link]` → issue browse URL or Confluence page title.
- Unless `disable_translation`: applies `jira_to_markdown` then
  `_convert_html_to_markdown`.

### `jira_to_markdown` (Jira wiki → Markdown)

Conversion table:

| Jira wiki | Markdown |
| --- | --- |
| `{code[:lang]}...{code}` | fenced ``` ``` ``` blocks |
| `{noformat}...{noformat}` | fenced ` ``` ` blocks (no language) |
| `{{...}}` (inline code) | `` `...` `` |
| `bq.` | `> ` (blockquote line) |
| `*x*` / `_x_` | `**x**` / `*x*` |
| Multi-level `*`/`-`/`+`/`#` lists | nested `- ` / `1. ` |
| `h1.`-`h6.` | `#`-`######` |
| `??cite??` | `<cite>` |
| `+ins+` | `<ins>` |
| `^sup^` | `<sup>` |
| `~sub~` | `<sub>` |
| `-strikethrough-` | `~~strikethrough~~` |
| `{quote}...{quote}` | `> ` prefixed lines |
| `{panel[:title]}...{panel}` | bold title + content |
| `!src|alt=alt,...params!` | `![alt](src)` |
| `[label|url]` | `[label](url)` |
| `{color:#xxx}...{color}` | `<span style="color:#xxx">` |
| `||` tables | Markdown tables |

The function uses private-use placeholders (`\x00PREFIX<N>\x00`) to
protect fenced/inline code blocks from downstream regex transforms.

### `markdown_to_jira` (Markdown → Jira wiki)

Conversion table:

| Markdown | Jira wiki |
| --- | --- |
| `===` underlined headers | `h1.` |
| `---` underlined headers | `h2.` |
| `# ` / `## ` / ... | `h1.` / `h2.` / ... (space required) |
| Intraword underscores | escaped `snake\_case` |
| `**x**` | `*x*` |
| `*x*` | `*x*` |
| Multi-level `-` / `+` / `*` | `*` / `**` / `***` |
| Multi-level `1.` | `#` / `##` / `###` |
| `<cite>` | `??cite??` |
| `<del>` | `-strikethrough-` |
| `<ins>` | `+ins+` |
| `<sup>` | `^sup^` |
| `<sub>` | `~sub~` |
| `<span style="color:#xxx">` | `{color:#xxx}` |
| `~~strike~~` | `-strike-` |
| `![](src)` | `!src!` |
| `![alt](src)` | `!src|alt=alt!` |
| `[label](url)` | `[label|url]` |
| `<url>` | `[url]` |
| Markdown tables | `||` header + `|` rows |

### Language normalization

`VALID_JIRA_LANGUAGES` (47 entries):

```
actionscript, actionscript3, ada, applescript, bash (sh), c, c# (csharp, cs),
c++ (cpp), css, sass, less, coldfusion, delphi, diff (patch), erlang (erl),
go, groovy, haskell, html, xml, java, javafx, javascript (js), json, lua, nyan,
objc (objective-c), perl, php, powershell (ps1), python (py), r, rainbow,
ruby (rb), scala, sql, swift, visualbasic (vb), yaml (yml), none
```

`LANGUAGE_MAPPING` (12 entries):

```
dockerfile, docker → bash
typescript, ts, tsx, jsx → javascript
kotlin, kt → java
makefile, make, cmake → bash
```

### Cloud-vs-DC differences

- Mentions: Cloud uses `[~accountid:<id>]`, older DC uses `[~username]` or
  `[~userkey]`. The regex in `jira_to_markdown` accepts both.

## Go port notes

### HTML → Markdown

The Go port needs an HTML-to-Markdown library. Options:

- [`github.com/JohannesKaufmann/html-to-markdown`](https://github.com/JohannesKaufmann/html-to-markdown) — most direct equivalent of `markdownify`.
- [`github.com/yuin/goldmark`](https://github.com/yuin/goldmark) — Markdown ↔ HTML.

For HTML parsing use `golang.org/x/net/html` + a small wrapper, or
`github.com/PuerkitoBio/goquery` for jQuery-like traversal.

### Markdown → Confluence storage

There's no direct `md2conf` equivalent in Go. The Go port should:

1. Convert Markdown to HTML with a library.
2. Post-process the HTML:
   - Fix attachment images (`<ac:image>` + `<ri:attachment>`).
   - Apply task-list rewriting.
   - Apply table-layout attribute injection.
3. Wrap the result in a Confluence storage-format JSON payload
   (`{body:{storage:{value, representation}}}`).

### Markdown ↔ Jira wiki

The Go port needs a Jira-wiki ↔ Markdown converter. The current Python
implementation uses ~30 regex substitutions. The Go port can use the
same regexes via `regexp/syntax` (compatible with PCRE syntax).

The private-use placeholder trick (`\x00PREFIX<N>\x00`) maps directly to
Go strings — `\x00` is a valid rune in Go and doesn't conflict with
UTF-8.

### User profile macros

For Confluence user-profile-macro resolution, define a Go interface:

```go
type UserResolver interface {
    GetUserByAccountID(ctx context.Context, accountID string) (*User, error)
    GetUserByUsername(ctx context.Context, username string) (*User, error)
    GetUserByUserKey(ctx context.Context, userKey string) (*User, error)
}
```

The preprocessor uses these to resolve mentions to `@DisplayName` or
the placeholder when resolution fails.

### Storage format details

Confluence storage format is XHTML with namespace macros:

- `<ac:structured-macro ac:name="profile">...</ac:structured-macro>`
- `<ac:link><ri:page ri:content-title="..." /></ac:link>`
- `<ac:image ac:alt="..."><ri:attachment ri:filename="..." /></ac:image>`
- `<ac:task-list><ac:task>...</ac:task></ac:task-list>`

These must be preserved exactly when round-tripping.