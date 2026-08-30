# Таблица соответствия инструментов: Go → Python

Регистрация инструментов в Go разнесена по файлам: [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go), [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go), [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go). Общие хелперы: [`internal/server/tools.go`](internal/server/tools.go). Всего 99 инструментов.
Python-реализации: [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) (63 инструмента Jira) и [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) (35 инструментов Confluence).

## Jira — Чтение (Read)

| Инструмент | Go-файлы | Python-файл |
|---|---|---|
| `jira_get_issue` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_project_issues` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_issue_dates` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_issue_sla` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_issue_images` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_download_attachments` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_service_desk_for_project` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_service_desk_queues` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_queue_issues` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_request_types` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_request_type_fields` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_user_profile` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_search` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_transitions` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_worklog` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_issue_watchers` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_link_types` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_all_projects` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_search_projects` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_project_issue_types` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_create_fields` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_project_versions` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_project_components` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_project_fields` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_agile_boards` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_board_issues` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_sprints_from_board` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_sprint_issues` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_search_fields` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_field_options` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_search_assignable_users` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_issue_development_info` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go), [`internal/jira/fetcher_ext.go`](internal/jira/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_issues_development_info` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go), [`internal/jira/fetcher_ext.go`](internal/jira/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_project_epic_hierarchy` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go), [`internal/jira/fetcher_ext.go`](internal/jira/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_get_cross_project_dependencies` | [`internal/server/tools_jira_read.go`](internal/server/tools_jira_read.go), [`internal/jira/fetcher_ext.go`](internal/jira/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |

## Jira — Запись (Write)

| Инструмент | Go-файлы | Python-файл |
|---|---|---|
| `jira_create_issue` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_update_issue` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_delete_issue` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_assign_issue` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_transition_issue` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_add_worklog` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_add_watcher` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_remove_watcher` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_link_to_epic` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_create_issue_link` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_create_remote_issue_link` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_remove_issue_link` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_add_comment` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_edit_comment` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_create_version` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_update_version` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_batch_create_issues` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_create_customer_request` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_create_sprint` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_update_sprint` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_add_issues_to_sprint` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |
| `jira_move_issues_to_backlog` | [`internal/server/tools_jira_write.go`](internal/server/tools_jira_write.go) | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) |

## Confluence — Чтение (Read)

| Инструмент | Go-файлы | Python-файл |
|---|---|---|
| `confluence_get_page` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_search` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_page_children` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_space_page_tree` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_comments` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_labels` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_search_user` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_attachments` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_inline_comments` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_page_history` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_page_diff` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_download_attachment` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_download_content_attachments` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_page_images` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_get_page_restrictions` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_scan_content` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_get_content_history` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_get_descendants` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_get_space` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_get_space_properties` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_get_current_user` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_search_entities` | [`internal/server/tools_conf_read.go`](internal/server/tools_conf_read.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |

## Confluence — Запись (Write)

| Инструмент | Go-файлы | Python-файл |
|---|---|---|
| `confluence_create_page` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_update_page` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_delete_page` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_move_page` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_add_comment` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_reply_to_comment` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_add_label` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_delete_attachment` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go) |(_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_update_page_section` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_add_inline_comment` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_upload_attachment` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_upload_attachments` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_set_page_restrictions` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_copy_page` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) |
| `confluence_remove_label` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_set_page_property` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_delete_page_property` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_update_attachment` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |
| `confluence_update_attachment_data` | [`internal/server/tools_conf_write.go`](internal/server/tools_conf_write.go), [`internal/confluence/fetcher_ext.go`](internal/confluence/fetcher_ext.go) | — (DC-only, нет в Python) |

---

## Инструменты только в Python (отсутствуют в Go)

Эти инструменты есть в `_old_python`, но ещё не портированы на Go. Все они — **Cloud-only**.

### Jira

| Инструмент | Python-файл | Причина |
|---|---|---|
| `jira_batch_get_changelogs` | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) | Cloud-only: `/rest/api/3/changelog/bulkfetch` |
| `jira_move_issue` | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) | Cloud-only: `/rest/api/3/bulk/issues/move` |
| `jira_get_issue_proforma_forms` | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) | Cloud-only: ProForma/Jira Forms |
| `jira_get_proforma_form_details` | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) | Cloud-only: ProForma/Jira Forms |
| `jira_update_proforma_form_answers` | [`_old_python/src/mcp_atlassian/servers/jira.py`](_old_python/src/mcp_atlassian/servers/jira.py) | Cloud-only: ProForma/Jira Forms |

### Confluence

| Инструмент | Python-файл | Причина |
|---|---|---|
| `confluence_get_page_views` | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) | Cloud-only: Analytics API |
| `confluence_list_page_templates` | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) | Cloud-only: `/rest/api/template/` |
| `confluence_get_page_template` | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) | Cloud-only: `/rest/api/template/` |
| `confluence_create_page_from_template` | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) | Cloud-only: зависит от template API |
| `confluence_check_content_permissions` | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) | Cloud-only: permission check API |
| `confluence_get_space_permissions` | [`_old_python/src/mcp_atlassian/servers/confluence.py`](_old_python/src/mcp_atlassian/servers/confluence.py) | Cloud-only: v2 permissions API |

---

## Сводка

| Категория | Go | Python |
|---|---|---|
| Jira Read | 35 | 40 |
| Jira Write | 23 | 23 |
| Confluence Read | 23 | 17 |
| Confluence Write | 18 | 18 |
| **Всего** | **99** | **98** |

> **Примечание**: Из 13 оставшихся Python-инструментов все являются Cloud-only (5 Jira + 8 Confluence). Go DC-версия содержит 12 дополнительных инструментов, отсутствующих в Python (DC-only API: scan, history, descendants, space, space_properties, current_user, search_entities, remove_label, set_page_property, delete_page_property, update_attachment, update_attachment_data).
