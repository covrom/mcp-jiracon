# Tool catalog — quick index

This file lists every MCP tool exposed by `mcp-atlassian`. Each tool is
grouped by its **`toolset:` tag**. Use this as the entry point when looking
up a tool by name.

> **Naming rule** (from `servers/main.py`): each tool is mounted under a
> namespace prefix — `jira.` or `confluence.` — so the fully-qualified MCP
> name is `jira_<function_name>` or `confluence_<function_name>`. The local
> Python function name has no prefix; the prefix is added by the MCP mount.

---

## Jira tools (63 total)

| Tool name (local) | MCP tool name | Toolset | Read/Write |
| --- | --- | --- | --- |
| `get_user_profile` | `jira_get_user_profile` | `jira_users` | R |
| `search_assignable_users` | `jira_search_assignable_users` | `jira_users` | R |
| `get_issue_watchers` | `jira_get_issue_watchers` | `jira_watchers` | R |
| `add_watcher` | `jira_add_watcher` | `jira_watchers` | W |
| `remove_watcher` | `jira_remove_watcher` | `jira_watchers` | W |
| `get_issue` | `jira_get_issue` | `jira_issues` | R |
| `search` | `jira_search` | `jira_issues` | R |
| `get_project_issues` | `jira_get_project_issues` | `jira_issues` | R |
| `create_issue` | `jira_create_issue` | `jira_issues` | W |
| `batch_create_issues` | `jira_batch_create_issues` | `jira_issues` | W |
| `batch_get_changelogs` | `jira_batch_get_changelogs` | `jira_issues` | R |
| `update_issue` | `jira_update_issue` | `jira_issues` | W |
| `assign_issue` | `jira_assign_issue` | `jira_issues` | W |
| `delete_issue` | `jira_delete_issue` | `jira_issues` | W |
| `move_issue` | `jira_move_issue` | `jira_issues` | W |
| `get_transitions` | `jira_get_transitions` | `jira_transitions` | R |
| `transition_issue` | `jira_transition_issue` | `jira_transitions` | W |
| `get_worklog` | `jira_get_worklog` | `jira_worklog` | R |
| `add_worklog` | `jira_add_worklog` | `jira_worklog` | W |
| `download_attachments` | `jira_download_attachments` | `jira_attachments` | R (binary) |
| `get_issue_images` | `jira_get_issue_images` | `jira_attachments` | R (binary) |
| `get_agile_boards` | `jira_get_agile_boards` | `jira_agile` | R |
| `get_board_issues` | `jira_get_board_issues` | `jira_agile` | R |
| `get_sprints_from_board` | `jira_get_sprints_from_board` | `jira_agile` | R |
| `get_sprint_issues` | `jira_get_sprint_issues` | `jira_agile` | R |
| `create_sprint` | `jira_create_sprint` | `jira_agile` | W |
| `update_sprint` | `jira_update_sprint` | `jira_agile` | W |
| `add_issues_to_sprint` | `jira_add_issues_to_sprint` | `jira_agile` | W |
| `move_issues_to_backlog` | `jira_move_issues_to_backlog` | `jira_agile` | W |
| `get_link_types` | `jira_get_link_types` | `jira_links` | R |
| `link_to_epic` | `jira_link_to_epic` | `jira_links` | W |
| `create_issue_link` | `jira_create_issue_link` | `jira_links` | W |
| `create_remote_issue_link` | `jira_create_remote_issue_link` | `jira_links` | W |
| `remove_issue_link` | `jira_remove_issue_link` | `jira_links` | W |
| `add_comment` | `jira_add_comment` | `jira_comments` | W |
| `edit_comment` | `jira_edit_comment` | `jira_comments` | W |
| `get_project_issue_types` | `jira_get_project_issue_types` | `jira_projects` | R |
| `get_create_fields` | `jira_get_create_fields` | `jira_projects` | R |
| `get_project_versions` | `jira_get_project_versions` | `jira_projects` | R |
| `get_project_components` | `jira_get_project_components` | `jira_projects` | R |
| `get_all_projects` | `jira_get_all_projects` | `jira_projects` | R |
| `search_projects` | `jira_search_projects` | `jira_projects` | R |
| `get_project_fields` | `jira_get_project_fields` | `jira_projects` | R |
| `create_version` | `jira_create_version` | `jira_projects` | W |
| `batch_create_versions` | `jira_batch_create_versions` | `jira_projects` | W |
| `update_version` | `jira_update_version` | `jira_projects` | W |
| `get_service_desk_for_project` | `jira_get_service_desk_for_project` | `jira_service_desk` | R |
| `get_service_desk_queues` | `jira_get_service_desk_queues` | `jira_service_desk` | R |
| `get_queue_issues` | `jira_get_queue_issues` | `jira_service_desk` | R |
| `get_request_types` | `jira_get_request_types` | `jira_service_desk` | R |
| `get_request_type_fields` | `jira_get_request_type_fields` | `jira_service_desk` | R |
| `create_customer_request` | `jira_create_customer_request` | `jira_service_desk` | W |
| `search_fields` | `jira_search_fields` | `jira_fields` | R |
| `get_field_options` | `jira_get_field_options` | `jira_fields` | R |
| `get_issue_proforma_forms` | `jira_get_issue_proforma_forms` | `jira_forms` | R |
| `get_proforma_form_details` | `jira_get_proforma_form_details` | `jira_forms` | R |
| `update_proforma_form_answers` | `jira_update_proforma_form_answers` | `jira_forms` | W |
| `get_issue_dates` | `jira_get_issue_dates` | `jira_metrics` | R |
| `get_issue_sla` | `jira_get_issue_sla` | `jira_metrics` | R |
| `get_issue_development_info` | `jira_get_issue_development_info` | `jira_development` | R |
| `get_issues_development_info` | `jira_get_issues_development_info` | `jira_development` | R |
| `get_project_epic_hierarchy` | `jira_get_project_epic_hierarchy` | `jira_project_analysis` | R |
| `get_cross_project_dependencies` | `jira_get_cross_project_dependencies` | `jira_project_analysis` | R |

---

## Confluence tools (35 total)

| Tool name (local) | MCP tool name | Toolset | Read/Write |
| --- | --- | --- | --- |
| `search` | `confluence_search` | `confluence_pages` | R |
| `get_page` | `confluence_get_page` | `confluence_pages` | R |
| `get_page_children` | `confluence_get_page_children` | `confluence_pages` | R |
| `get_space_page_tree` | `confluence_get_space_page_tree` | `confluence_pages` | R |
| `create_page` | `confluence_create_page` | `confluence_pages` | W |
| `update_page` | `confluence_update_page` | `confluence_pages` | W |
| `update_page_section` | `confluence_update_page_section` | `confluence_pages` | W |
| `delete_page` | `confluence_delete_page` | `confluence_pages` | W |
| `move_page` | `confluence_move_page` | `confluence_pages` | W |
| `get_page_history` | `confluence_get_page_history` | `confluence_pages` | R |
| `get_page_diff` | `confluence_get_page_diff` | `confluence_pages` | R |
| `get_page_restrictions` | `confluence_get_page_restrictions` | `confluence_pages` | R |
| `set_page_restrictions` | `confluence_set_page_restrictions` | `confluence_pages` | W |
| `copy_page` | `confluence_copy_page` | `confluence_pages` | W |
| `get_comments` | `confluence_get_comments` | `confluence_comments` | R |
| `add_comment` | `confluence_add_comment` | `confluence_comments` | W |
| `reply_to_comment` | `confluence_reply_to_comment` | `confluence_comments` | W |
| `get_inline_comments` | `confluence_get_inline_comments` | `confluence_comments` | R |
| `add_inline_comment` | `confluence_add_inline_comment` | `confluence_comments` | W |
| `get_labels` | `confluence_get_labels` | `confluence_labels` | R |
| `add_label` | `confluence_add_label` | `confluence_labels` | W |
| `search_user` | `confluence_search_user` | `confluence_users` | R |
| `get_page_views` | `confluence_get_page_views` | `confluence_analytics` | R |
| `upload_attachment` | `confluence_upload_attachment` | `confluence_attachments` | W |
| `upload_attachments` | `confluence_upload_attachments` | `confluence_attachments` | W |
| `get_attachments` | `confluence_get_attachments` | `confluence_attachments` | R |
| `download_attachment` | `confluence_download_attachment` | `confluence_attachments` | R (binary) |
| `download_content_attachments` | `confluence_download_content_attachments` | `confluence_attachments` | R (binary) |
| `delete_attachment` | `confluence_delete_attachment` | `confluence_attachments` | W |
| `get_page_images` | `confluence_get_page_images` | `confluence_attachments` | R (binary) |
| `list_page_templates` | `confluence_list_page_templates` | `confluence_templates` | R |
| `get_page_template` | `confluence_get_page_template` | `confluence_templates` | R |
| `create_page_from_template` | `confluence_create_page_from_template` | `confluence_templates` | W |
| `check_content_permissions` | `confluence_check_content_permissions` | `confluence_permissions` | R |
| `get_space_permissions` | `confluence_get_space_permissions` | `confluence_permissions` | R |

---

## Cross-cutting toolset constants

The toolsets listed above are defined as enum-like entries in
`src/mcp_atlassian/utils/toolsets.py`. Reimplementers should mirror the same
names so the `ENABLED_TOOLSETS` environment variable continues to work.

```
JIRA_TOOLSETS = {
  "jira_issues":          ToolsetDefinition("jira_issues", …, default=True),
  "jira_fields":          ToolsetDefinition("jira_fields", …, default=True),
  "jira_comments":        ToolsetDefinition("jira_comments", …, default=True),
  "jira_transitions":     ToolsetDefinition("jira_transitions", …, default=True),
  "jira_projects":        ToolsetDefinition("jira_projects", …, default=False),
  "jira_agile":           ToolsetDefinition("jira_agile", …, default=False),
  "jira_links":           ToolsetDefinition("jira_links", …, default=False),
  "jira_worklog":         ToolsetDefinition("jira_worklog", …, default=False),
  "jira_attachments":     ToolsetDefinition("jira_attachments", …, default=False),
  "jira_users":           ToolsetDefinition("jira_users", …, default=False),
  "jira_watchers":        ToolsetDefinition("jira_watchers", …, default=False),
  "jira_service_desk":    ToolsetDefinition("jira_service_desk", …, default=False),
  "jira_forms":           ToolsetDefinition("jira_forms", …, default=False),
  "jira_metrics":         ToolsetDefinition("jira_metrics", …, default=False),
  "jira_development":     ToolsetDefinition("jira_development", …, default=False),
  "jira_project_analysis":ToolsetDefinition("jira_project_analysis", …, default=False),
}

CONFLUENCE_TOOLSETS = {
  "confluence_pages":      ToolsetDefinition("confluence_pages", …, default=True),
  "confluence_comments":   ToolsetDefinition("confluence_comments", …, default=True),
  "confluence_labels":     ToolsetDefinition("confluence_labels", …, default=False),
  "confluence_users":      ToolsetDefinition("confluence_users", …, default=False),
  "confluence_analytics":  ToolsetDefinition("confluence_analytics", …, default=False),
  "confluence_attachments":ToolsetDefinition("confluence_attachments", …, default=False),
  "confluence_templates":  ToolsetDefinition("confluence_templates", …, default=False),
  "confluence_permissions":ToolsetDefinition("confluence_permissions", …, default=False),
}
```

---

## Cross-tool constants used in JSON-schema

| Constant | Value | Source |
| --- | --- | --- |
| `ISSUE_KEY_PATTERN` | `r"^[A-Z][A-Z0-9_]+-\d+(?:-\d+)*$"` | `servers/jira.py` |
| `PROJECT_KEY_PATTERN` | `r"^[A-Z][A-Z0-9_]+$"` | `servers/jira.py` |
| `DEFAULT_READ_JIRA_FIELDS` | `priority,updated,labels,issuetype,summary,assignee,description,created,reporter,status` | `jira/constants.py` |
| `ATTACHMENT_MAX_BYTES` | `50 * 1024 * 1024` (50 MiB) | `utils/media.py` |
| `TOKEN_EXPIRY_MARGIN` | `300` seconds | `utils/oauth.py` |
| `_PAGE_ID_PATH_PATTERN` | `(?:^|/)pages/([0-9]+)(?:/|$)` | `servers/confluence.py` |
| `_TINY_LINK_PATH_PATTERN` | `(?:^|/)x/([A-Za-z0-9_-]+)(?:/|$)` | `servers/confluence.py` |

Go reimplementations should re-export these so any consumer can pin them
exactly.