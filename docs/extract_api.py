#!/usr/bin/env python3
"""Извлекает описание API для методов и эндпоинтов из docs.atlassian.com.html.

Список нужных методов берётся из list.md (формат: `- METHOD /path`). Для каждого
эндпоинта скрипт находит соответствующий блок в HTML-документации Confluence REST API
и формирует структурированное описание: метод, путь, заголовок, deprecated, описание,
path/query-параметры, тело запроса (схема/пример) и ответы (статусы + схема/пример).

Правила структуры method-body взяты из docs/struct.md:
  * секции <h5>Request</h5> / <h5>Responses</h5>; подзаголовки <h6>query parameters</h6>,
    <h6>path parameters</h6>, <h6>Schema</h6>, <h6>Example</h6>;
  * значения по умолчанию «спрятаны» внутри ячейки описания (<code>/<em>);
  * JSON-схемы могут содержать вложенные HTML-теги (<a href="#definitions/...">) —
    они очищаются до чистого JSON;
  * статусы ответов оформлены как <span class="aui-lozenge">Status NNN</span>
    (fallback — паттерн <strong>Status</strong> <strong>NNN</strong>);
  * между HTTP-методом и путём стоит &nbsp;.

Особенность разметки: у части методов всё содержимое (описание, Request, Responses)
вложено в один <p> внутри method-body, а у других секции — прямые дети. Скрипт
рекурсивно раскрывает такие обёртки, чтобы корректно разбить содержимое на секции.

Использование:
    python3 extract_api.py                 # markdown в api_docs.md
    python3 extract_api.py --json          # JSON в api_docs.json
    python3 extract_api.py -o out.md       # свой файл вывода
    python3 extract_api.py --list list.md --html docs.atlassian.com.html

Зависимости: beautifulsoup4 (pip install beautifulsoup4)
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from urllib.parse import unquote

from bs4 import BeautifulSoup, NavigableString, Tag

# ---------------------------------------------------------------------------
# Нормализация путей
# ---------------------------------------------------------------------------

#: Замещения имён плейсхолдеров, встречающихся в документации, на канонические
#: имена из list.md. Ключ — «METHOD PATH», значение — замена.
PATH_ALIASES = {
    "PUT /rest/api/content/{contentId}": "PUT /rest/api/content/{id}",
}


def _strip_ws(s: str | None) -> str:
    return re.sub(r"\s+", " ", s or "").strip()


def normalize_path(path: str) -> str:
    """Приводит путь к каноническому виду: без пробелов, с нормализацией плейсхолдеров."""
    path = unquote(path.replace("\xa0", " ").replace("&nbsp;", " "))
    return re.sub(r"\s+", "", path)


def canonical_key(method: str, path: str) -> str:
    """Канонический ключ «METHOD PATH» для сопоставления списка и документа."""
    key = f"{method.upper()} {normalize_path(path)}"
    return PATH_ALIASES.get(key, key)


# ---------------------------------------------------------------------------
# Разбор list.md
# ---------------------------------------------------------------------------

_ENDPOINT_RE = re.compile(r"^-\s*`?([A-Z]+)`?\s+`(/rest/api/\S+)`")


def parse_list_md(text: str) -> list[dict]:
    """Возвращает список {'method', 'path'} в порядке следования в файле."""
    items: list[dict] = []
    for line in text.splitlines():
        m = _ENDPOINT_RE.match(line.strip())
        if not m:
            continue
        items.append({"method": m.group(1).upper(), "path": m.group(2)})
    return items


# ---------------------------------------------------------------------------
# Вспомогательные функции извлечения
# ---------------------------------------------------------------------------

def _text(node) -> str:
    """Чистый текст узла (пробелы схлопнуты)."""
    if node is None:
        return ""
    return re.sub(r"\s+", " ", node.get_text(" ", strip=True)).strip()


def _clean_json(raw: str) -> str:
    """Убирает из JSON-схемы вложенные HTML-теги (например <a href="#definitions/...">),
    возвращая валидный JSON."""
    raw = re.sub(r"<\s*a\b[^>]*>", "", raw)   # открывающие теги <a ...>
    raw = re.sub(r"</\s*a\s*>", "", raw)      # закрывающие </a>
    return raw.strip()


def _schema_of(block) -> str:
    """JSON-схема внутри блока representation-doc (если есть)."""
    pre = block.find("pre")
    if pre:
        code = pre.find("code")
        return _clean_json((code.get_text() if code else pre.get_text()))
    return ""


def _example_of(block) -> str:
    """Текст примера (<h6>Example</h6>) внутри блока, если он есть."""
    h6 = block.find("h6")
    if h6 and "example" in h6.get_text(strip=True).lower():
        pre = h6.find_next("pre")
        if pre:
            code = pre.find("code")
            return (code.get_text() if code else pre.get_text()).strip()
    return ""


def _parse_param_table(table) -> list[dict]:
    """Парсит таблицу параметров (path/query): parameter / type / description.

    Значение по умолчанию «спрятано» внутри ячейки типа (<tt>/<em>) или описания
    (<code>Default: ...</code>).
    """
    params: list[dict] = []
    for tr in table.find_all("tr"):
        tds = tr.find_all("td")
        if len(tds) < 3:
            continue  # строка заголовков
        name_el = tds[0].find("code")
        name = _strip_ws(name_el.get_text()) if name_el else _strip_ws(tds[0].get_text())
        if not name:
            continue
        type_el = tds[1].find("em")
        ptype = _strip_ws(type_el.get_text()) if type_el else _strip_ws(tds[1].get_text())
        default_el = tds[1].find("tt")
        default = _strip_ws(default_el.get_text()) if default_el else None
        desc = _strip_ws(tds[2].get_text())
        if default is None:
            dm = re.search(r"Default(?:\s+value)?:\s*(.+)$", desc)
            if dm:
                default = _strip_ws(dm.group(1))
        params.append({"name": name, "type": ptype, "default": default, "description": desc})
    return params


def _flatten(body) -> list:
    """Все содержимые узлы method-body в документеном порядке.

    Рекурсивно раскрывает обёртки (например, когда всё содержимое вложено в один <p>),
    чтобы получить плоскую последовательность элементов. Возвращает только элементы
    (текстовые узлы отбрасываются).
    """
    result: list = []
    for child in body.children:
        if isinstance(child, Tag):
            # Раскрываем «обёрточные» контейнеры, не несущие смыслового класса.
            classes = set(child.get("class") or [])
            if child.name in ("div", "p", "section") and not classes:
                result.extend(_flatten(child))
            else:
                result.append(child)
        elif isinstance(child, NavigableString) and child.strip():
            result.append(NavigableString(child.strip()))
    return result


def _is_h5(node, label: str) -> bool:
    return isinstance(node, Tag) and node.name == "h5" and label in node.get_text(strip=True).lower()


def _split_sections(nodes: list) -> tuple[list, list, list]:
    """Разбивает узлы на (описание, Request, Responses) по заголовкам h5."""
    desc: list = []
    request: list = []
    responses: list = []
    current = desc
    for n in nodes:
        if _is_h5(n, "request"):
            current = request
            continue
        if _is_h5(n, "response"):
            current = responses
            continue
        current.append(n)
    return desc, request, responses


def _extract_description(desc_nodes: list) -> str:
    """Описание метода: текст узлов до Request/Responses.

    Включает абзацы, списки и прямой текст; исключает таблицы, схемы, примеры URI
    и блоки представления.
    """
    parts: list[str] = []
    for n in desc_nodes:
        if isinstance(n, Tag):
            if n.name in ("table", "pre", "h5", "h6"):
                continue
            classes = set(n.get("class") or [])
            if "representation-doc" in classes or "exampleRequests" in classes:
                continue
            txt = _text(n)
        elif isinstance(n, NavigableString):
            txt = n.strip()
        else:
            continue
        if txt:
            parts.append(txt)
    return " ".join(parts).strip()


def _extract_request(request_nodes: list) -> dict:
    """Разбирает секцию Request: path/query-параметры и тело запроса."""
    result: dict = {"path_parameters": [], "query_parameters": [], "body": None}
    mode = None
    for n in request_nodes:
        if not isinstance(n, Tag):
            continue
        if n.name == "h6":
            t = n.get_text(strip=True).lower()
            if "path" in t:
                mode = "path"
            elif "query" in t:
                mode = "query"
            elif "schema" in t or "example" in t:
                mode = "body"
            continue
        if n.name == "table":
            params = _parse_param_table(n)
            if mode == "path":
                result["path_parameters"].extend(params)
            else:
                result["query_parameters"].extend(params)
        elif n.name == "div" and "representation-doc" in (n.get("class") or []):
            schema = _schema_of(n)
            example = _example_of(n)
            if schema or example:
                result["body"] = {"schema": schema or None, "example": example or None}
    return result


_STATUS_LOZENGE_RE = re.compile(r"(\d{3})")


def _extract_responses(responses_nodes: list) -> list[dict]:
    """Разбирает секцию Responses: статусы, контент-тип, описание, схема/пример.

    Узлы секции могут быть как отдельные <li>, так и обёртка <ul>, содержащая их.
    Поэтому ищем li.representation в каждом узле (включая самого себя).
    """
    responses: list[dict] = []
    for node in responses_nodes:
        if not isinstance(node, Tag):
            continue
        reps = [node] if (node.name == "li" and "representation" in (node.get("class") or [])) \
            else node.find_all("li", class_="representation")
        for rep in reps:
            name_span = rep.find("span", class_="representation-name")
            status = None
            ctype = None
            if name_span:
                lozenge = name_span.find("span", class_="aui-lozenge")
                if lozenge:
                    sm = _STATUS_LOZENGE_RE.search(lozenge.get_text())
                    status = sm.group(1) if sm else None
                ctype_el = name_span.find("i")
                ctype = _strip_ws(ctype_el.get_text()) if ctype_el else None
            doc = rep.find("div", class_="representation-doc")
            # Описание ответа — только прямые текстовые узлы representation-doc;
            # блоки схем/примеров (representation-doc-block) не включаем, чтобы не
            # дублировать заголовки «Schema»/«Example» и сами JSON-схемы.
            desc_parts: list[str] = []
            if doc is not None:
                for c in doc.children:
                    if isinstance(c, NavigableString) and c.strip():
                        desc_parts.append(c.strip())
            desc = re.sub(r"\s+", " ", " ".join(desc_parts)).strip()
            schema = _schema_of(doc) if doc else ""
            example = _example_of(doc) if doc else ""
            responses.append({
                "status": status,
                "content_type": ctype,
                "description": desc,
                "schema": schema or None,
                "example": example or None,
            })
    return responses


def _extract_methods(html: str) -> dict[str, dict]:
    """Возвращает словарь: канонический ключ -> данные метода."""
    soup = BeautifulSoup(html, "html.parser")
    methods: dict[str, dict] = {}

    for res_div in soup.find_all("div", class_="resource"):
        h3 = res_div.find("h3")
        if h3 is None:
            continue
        resource_id = h3.get("id", "")
        res_p = h3.find_next_sibling("p")
        resource_desc = _text(res_p) if res_p else ""

        methods_div = res_div.find("div", class_="methods")
        if methods_div is None:
            continue

        for m_div in methods_div.find_all("div", class_="method", recursive=False):
            h4 = m_div.find("h4")
            if h4 is None:
                continue
            code = h4.find("code")
            if code is None:
                continue
            raw = code.get_text().replace("\xa0", " ")
            m = re.match(r"^\s*(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+(\S+)\s*$", raw)
            if not m:
                continue
            method, path = m.group(1), m.group(2)
            key = canonical_key(method, path)

            title_el = h4.find("a")
            title = _strip_ws(title_el.get_text()) if title_el else ""
            deprecated = h4.find("span", class_="deprecated") is not None

            body = m_div.find("div", class_="method-body")
            entry: dict = {
                "resource": resource_id,
                "resource_description": resource_desc,
                "title": title,
                "method": method,
                "path": path,
                "deprecated": deprecated,
                "description": "",
                "path_parameters": [],
                "query_parameters": [],
                "body": None,
                "responses": [],
            }
            if body is not None:
                desc_nodes, req_nodes, resp_nodes = _split_sections(_flatten(body))
                entry["description"] = _extract_description(desc_nodes)
                entry.update(_extract_request(req_nodes))
                entry["responses"] = _extract_responses(resp_nodes)

            methods[key] = entry
    return methods


# ---------------------------------------------------------------------------
# Вывод
# ---------------------------------------------------------------------------

def _render_params(lines: list[str], params: list[dict], heading: str) -> None:
    lines.append(f"### {heading}")
    lines.append("")
    lines.append("| Параметр | Тип | По умолчанию | Описание |")
    lines.append("| --- | --- | --- | --- |")
    for p in params:
        d = (p["description"] or "").replace("|", "\\|")
        default = f"`{p['default']}`" if p["default"] else "—"
        lines.append(f"| `{p['name']}` | {p['type'] or '—'} | {default} | {d or '—'} |")
    lines.append("")


def _render_code_block(lines: list[str], label: str, content: str) -> None:
    lines.append(f"<details><summary>{label}</summary>")
    lines.append("")
    lines.append("```json")
    lines.append(content)
    lines.append("```")
    lines.append("")
    lines.append("</details>")
    lines.append("")


def render_markdown(items: list[dict], found: dict[str, dict]) -> str:
    lines: list[str] = ["# Confluence REST API — описание методов", ""]
    missing: list[str] = []
    for it in items:
        key = canonical_key(it["method"], it["path"])
        data = found.get(key)
        if data is None:
            missing.append(f"- `{it['method']} {it['path']}`")
            continue

        head = f"## {data['method']} `{data['path']}`"
        if data["deprecated"]:
            head += "  ⚠️ deprecated"
        lines.append(head)
        lines.append("")
        if data["title"]:
            lines.append(f"**{data['title']}**")
            lines.append("")
        if data["description"]:
            lines.append(data["description"])
            lines.append("")
        if data["path_parameters"]:
            _render_params(lines, data["path_parameters"], "Path-параметры")
        if data["query_parameters"]:
            _render_params(lines, data["query_parameters"], "Query-параметры")
        if data["body"]:
            lines.append("### Тело запроса")
            lines.append("")
            if data["body"].get("schema"):
                _render_code_block(lines, "Схема (JSON)", data["body"]["schema"])
            if data["body"].get("example"):
                _render_code_block(lines, "Пример", data["body"]["example"])
        if data["responses"]:
            lines.append("### Ответы")
            lines.append("")
            for r in data["responses"]:
                head = f"**{r['status'] or '?'}**"
                if r["content_type"]:
                    head += f" — `{r['content_type']}`"
                lines.append(head)
                if r["description"]:
                    lines.append("")
                    lines.append(r["description"])
                if r["schema"]:
                    lines.append("")
                    _render_code_block(lines, "Схема (JSON)", r["schema"])
                if r["example"]:
                    lines.append("")
                    _render_code_block(lines, "Пример", r["example"])
            lines.append("")

    lines.append("---")
    lines.append("")
    lines.append(f"Найдено {len(items) - len(missing)} из {len(items)} эндпоинтов.")
    if missing:
        lines.append("")
        lines.append("Не найдены в документе:")
        lines.extend(missing)
    lines.append("")
    return "\n".join(lines)


def main(argv: list[str] | None = None) -> int:
    here = Path(__file__).resolve().parent
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--html", default=str(here / "docs.atlassian.com.html"),
                    help="путь к HTML-документации (по умолчанию рядом со скриптом)")
    ap.add_argument("--list", default=str(here / "list.md"),
                    help="файл со списком эндпоинтов (по умолчанию рядом со скриптом)")
    ap.add_argument("-o", "--output", default=None,
                    help="файл вывода (по умолчанию api_docs.md / api_docs.json)")
    ap.add_argument("--json", action="store_true", help="вывести JSON вместо markdown")
    args = ap.parse_args(argv)

    html_path = Path(args.html)
    list_path = Path(args.list)
    if not html_path.is_file():
        print(f"Ошибка: не найден HTML-файл {html_path}", file=sys.stderr)
        return 1
    if not list_path.is_file():
        print(f"Ошибка: не найден список {list_path}", file=sys.stderr)
        return 1

    items = parse_list_md(list_path.read_text(encoding="utf-8"))
    if not items:
        print("Ошибка: в списке нет ни одного эндпоинта", file=sys.stderr)
        return 1

    found = _extract_methods(html_path.read_text(encoding="utf-8"))

    if args.json:
        payload = {
            "endpoints": [found.get(canonical_key(i["method"], i["path"])) for i in items],
            "missing": [f"{i['method']} {i['path']}" for i in items
                        if found.get(canonical_key(i["method"], i["path"])) is None],
        }
        out_path = Path(args.output) if args.output else here / "api_docs.json"
        out_path.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
    else:
        md = render_markdown(items, found)
        out_path = Path(args.output) if args.output else here / "api_docs.md"
        out_path.write_text(md, encoding="utf-8")

    found_n = sum(1 for i in items if found.get(canonical_key(i["method"], i["path"])))
    print(f"Готово: {found_n}/{len(items)} эндпоинтов -> {out_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())