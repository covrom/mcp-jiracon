package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/atlassian"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
)

func newFetcher(t *testing.T, handler http.HandlerFunc) *Fetcher {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	cfg := &config.JiraConfig{URL: srv.URL, PersonalToken: "test-pat"}
	client, err := atlassian.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return New(cfg, client)
}

func jSON(obj any) string { b, _ := json.Marshal(obj); return string(b) }

func TestGetIssue(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/2/issue/PROJ-1" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jSON(map[string]any{"id": "10001", "key": "PROJ-1", "fields": map[string]any{"summary": "test"}})))
	})
	issue, err := f.GetIssue(context.Background(), "PROJ-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if issue.Key != "PROJ-1" {
		t.Errorf("key=%s", issue.Key)
	}
}

func TestGetUserProfile(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("username") != "alice" {
			t.Errorf("username=%s", r.URL.Query().Get("username"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jSON(map[string]any{"accountId": "u1", "displayName": "Alice", "active": true})))
	})
	user, err := f.GetUserProfile(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if user.DisplayName != "Alice" {
		t.Errorf("name=%s", user.DisplayName)
	}
}

// --- Watchers -----------------------------------------------------------

func TestGetIssueWatchers(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if !matchPath(r.URL.Path, "/rest/api/2/issue/PROJ-1/watchers") {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"isWatching": true, "watchCount": 2, "watchers": []any{}})))
	})
	raw, err := f.GetIssueWatchers(context.Background(), "PROJ-1")
	if err != nil {
		t.Fatal(err)
	}
	if raw["watchCount"].(float64) != 2 {
		t.Errorf("wc=%v", raw["watchCount"])
	}
}

func TestAddWatcher(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if name, ok := body["name"].(string); !ok || name != "user1" {
			t.Errorf("body[\"name\"] = %v, want \"user1\"", body["name"])
		}
		if _, ok := body["accountId"]; ok {
			t.Error("body must not contain accountId (DC-only)", body["accountId"])
		}
		w.WriteHeader(204)
	})
	if err := f.AddWatcher(context.Background(), "PROJ-1", "user1"); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveWatcher(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method=%s", r.Method)
		}
		w.WriteHeader(204)
	})
	if err := f.RemoveWatcher(context.Background(), "PROJ-1", "user1"); err != nil {
		t.Fatal(err)
	}
}

// --- Transitions --------------------------------------------------------

func TestGetTransitions(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"transitions": []map[string]any{{"id": "11", "name": "Done"}}})))
	})
	tr, err := f.GetTransitions(context.Background(), "PROJ-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(tr) != 1 || tr[0]["name"] != "Done" {
		t.Errorf("transitions=%v", tr)
	}
}

func TestTransitionIssue(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "10001", "key": "PROJ-1", "fields": map[string]any{}})))
	})
	issue, err := f.TransitionIssue(context.Background(), "PROJ-1", "11", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if issue.Key != "PROJ-1" {
		t.Errorf("key=%s", issue.Key)
	}
}

// --- Worklog ------------------------------------------------------------

func TestGetWorklogs(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"worklogs": []map[string]any{{"id": "w1", "timeSpent": "1h"}}})))
	})
	wl, err := f.GetWorklogs(context.Background(), "PROJ-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(wl) != 1 {
		t.Errorf("len=%d", len(wl))
	}
}

func TestAddWorklog(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "w2", "timeSpent": "30m"})))
	})
	raw, err := f.AddWorklog(context.Background(), "PROJ-1", "30m", "")
	if err != nil {
		t.Fatal(err)
	}
	if raw["timeSpent"] != "30m" {
		t.Errorf("ts=%v", raw["timeSpent"])
	}
}

// --- Links --------------------------------------------------------------

func TestGetIssueLinkTypes(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"issueLinkTypes": []map[string]any{{"name": "Blocks"}}})))
	})
	lt, err := f.GetIssueLinkTypes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(lt) != 1 || lt[0]["name"] != "Blocks" {
		t.Errorf("lt=%v", lt)
	}
}

func TestLinkIssueToEpic(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "10001", "key": "PROJ-1", "fields": map[string]any{}})))
	})
	issue, err := f.LinkIssueToEpic(context.Background(), "PROJ-1", "EPIC-1")
	if err != nil {
		t.Fatal(err)
	}
	if issue.Key != "PROJ-1" {
		t.Errorf("key=%s", issue.Key)
	}
}

func TestCreateIssueLink(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(`{}`))
	})
	_, err := f.CreateIssueLink(context.Background(), "Blocks", "PROJ-1", "PROJ-2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateRemoteIssueLink(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(`{}`))
	})
	_, err := f.CreateRemoteIssueLink(context.Background(), "PROJ-1", "http://x.com", "X", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveIssueLink(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method=%s", r.Method)
		}
		w.WriteHeader(204)
	})
	if err := f.RemoveIssueLink(context.Background(), "100"); err != nil {
		t.Fatal(err)
	}
}

// --- Comments -----------------------------------------------------------

func TestAddComment(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "c1", "body": "hello"})))
	})
	raw, err := f.AddComment(context.Background(), "PROJ-1", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if raw["body"] != "hello" {
		t.Errorf("body=%v", raw["body"])
	}
}

func TestEditComment(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "c1", "body": "updated"})))
	})
	raw, err := f.EditComment(context.Background(), "PROJ-1", "c1", "updated")
	if err != nil {
		t.Fatal(err)
	}
	if raw["body"] != "updated" {
		t.Errorf("body=%v", raw["body"])
	}
}

// --- Projects -----------------------------------------------------------

func TestGetAllProjects(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"values": []map[string]any{{"key": "PROJ", "name": "Project"}}})))
	})
	projs, err := f.GetAllProjects(context.Background(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(projs) != 1 {
		t.Errorf("len=%d", len(projs))
	}
}

func TestSearchProjects(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("query") != "proj" {
			t.Errorf("query=%s", q.Get("query"))
		}
		w.Write([]byte(jSON(map[string]any{"projects": []map[string]any{{"key": "PROJ"}}})))
	})
	projs, err := f.SearchProjects(context.Background(), "proj", 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(projs) != 1 {
		t.Errorf("len=%d", len(projs))
	}
}

func TestGetProjectIssueTypes(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"values": []map[string]any{{"id": "3", "name": "Task"}}})))
	})
	types, err := f.GetProjectIssueTypes(context.Background(), "PROJ")
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 1 {
		t.Errorf("len=%d", len(types))
	}
}

func TestGetCreateFields(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"fields": []map[string]any{{"fieldId": "summary", "name": "Summary", "required": true}}})))
	})
	fields, err := f.GetCreateFields(context.Background(), "PROJ", "3")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Errorf("len=%d", len(fields))
	}
}

func TestGetProjectVersions(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON([]map[string]any{{"id": "v1", "name": "1.0"}})))
	})
	vers, err := f.GetProjectVersions(context.Background(), "PROJ")
	if err != nil {
		t.Fatal(err)
	}
	if len(vers) != 1 {
		t.Errorf("len=%d", len(vers))
	}
}

func TestGetProjectComponents(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON([]map[string]any{{"id": "c1", "name": "Frontend"}})))
	})
	comps, err := f.GetProjectComponents(context.Background(), "PROJ")
	if err != nil {
		t.Fatal(err)
	}
	if len(comps) != 1 {
		t.Errorf("len=%d", len(comps))
	}
}

func TestCreateVersion(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "v1", "name": "1.0"})))
	})
	v, err := f.CreateVersion(context.Background(), "PROJ", "1.0", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if v["name"] != "1.0" {
		t.Errorf("name=%v", v["name"])
	}
}

func TestUpdateVersion(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "v1", "name": "2.0"})))
	})
	v, err := f.UpdateVersion(context.Background(), "v1", "2.0", "", nil, nil, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if v["name"] != "2.0" {
		t.Errorf("name=%v", v["name"])
	}
}

// --- Agile --------------------------------------------------------------

func TestGetAgileBoards(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"values": []map[string]any{{"id": "1", "name": "Board", "type": "scrum"}}})))
	})
	boards, err := f.GetAgileBoards(context.Background(), "", "", "", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(boards) != 1 {
		t.Errorf("len=%d", len(boards))
	}
}

func TestGetBoardIssues(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"issues": []map[string]any{{"key": "PROJ-1"}}})))
	})
	raw, err := f.GetBoardIssues(context.Background(), "1", "project=PROJ", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	issues, _ := raw["issues"].([]any)
	if len(issues) != 1 {
		t.Errorf("len=%d", len(issues))
	}
}

func TestGetSprintsFromBoard(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"values": []map[string]any{{"id": "s1", "name": "Sprint 1"}}})))
	})
	sp, err := f.GetSprintsFromBoard(context.Background(), "1", "active", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(sp) != 1 {
		t.Errorf("len=%d", len(sp))
	}
}

func TestGetSprintIssues(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"issues": []map[string]any{{"key": "PROJ-1"}}})))
	})
	raw, err := f.GetSprintIssues(context.Background(), "s1", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	issues, _ := raw["issues"].([]any)
	if len(issues) != 1 {
		t.Errorf("len=%d", len(issues))
	}
}

func TestCreateSprint(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "s1", "name": "Sprint 1"})))
	})
	s, err := f.CreateSprint(context.Background(), "1", "Sprint 1", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if s["name"] != "Sprint 1" {
		t.Errorf("name=%v", s["name"])
	}
}

func TestUpdateSprint(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "s1", "name": "Updated"})))
	})
	s, err := f.UpdateSprint(context.Background(), "s1", "Updated", "", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if s["name"] != "Updated" {
		t.Errorf("name=%v", s["name"])
	}
}

func TestAddIssuesToSprint(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.WriteHeader(204)
	})
	if err := f.AddIssuesToSprint(context.Background(), "s1", []string{"PROJ-1"}); err != nil {
		t.Fatal(err)
	}
}

func TestMoveIssuesToBacklog(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.WriteHeader(204)
	})
	if err := f.MoveIssuesToBacklog(context.Background(), []string{"PROJ-1"}); err != nil {
		t.Fatal(err)
	}
}

// --- Fields -------------------------------------------------------------

func TestSearchFields(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON([]map[string]any{{"id": "summary", "name": "Summary"}, {"id": "customfield_10010", "name": "Story Points"}})))
	})
	fields, err := f.SearchFields(context.Background(), "story", 10, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 2 {
		t.Errorf("len=%d", len(fields))
	}
}

func TestGetFieldOptions(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/rest/api/2/issue/createmeta/PROJ/issuetypes":
			// Step 1: project issue types.
			w.Write([]byte(jSON(map[string]any{
				"values": []map[string]any{
					{"id": "1", "name": "Bug"},
					{"id": "3", "name": "Task"},
				},
			})))
		case "/rest/api/2/issue/createmeta/PROJ/issuetypes/3":
			// Step 2: createmeta fields with pagination.
			w.Write([]byte(jSON(map[string]any{
				"startAt":    0,
				"maxResults": 50,
				"total":      2,
				"values": []map[string]any{
					{
						"fieldId":       "summary",
						"name":          "Summary",
						"required":      true,
						"allowedValues": nil,
					},
					{
						"fieldId":  "customfield_10001",
						"name":     "Components",
						"required": false,
						"allowedValues": []any{
							map[string]any{"id": "1", "value": "Option A"},
							map[string]any{"id": "2", "value": "Option B", "disabled": true},
						},
					},
				},
			})))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	})
	opts, err := f.GetFieldOptions(context.Background(), "customfield_10001", "PROJ", "Task")
	if err != nil {
		t.Fatal(err)
	}
	if len(opts) != 2 {
		t.Fatalf("expected 2 options, got %d", len(opts))
	}
	if opts[0]["id"] != "1" || opts[0]["value"] != "Option A" {
		t.Errorf("option[0]=%v", opts[0])
	}
	if opts[1]["disabled"] != true {
		t.Errorf("option[1] should be disabled: %v", opts[1])
	}
}

func TestGetFieldOptionsCaseInsensitiveMatch(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/rest/api/2/issue/createmeta/PROJ/issuetypes":
			w.Write([]byte(jSON(map[string]any{
				"values": []map[string]any{
					{"id": "5", "name": "Epic"},
				},
			})))
		case "/rest/api/2/issue/createmeta/PROJ/issuetypes/5":
			w.Write([]byte(jSON(map[string]any{
				"startAt":    0,
				"maxResults": 50,
				"total":      1,
				"values": []map[string]any{
					{
						"fieldId": "customfield_10001",
						"name":    "My Field",
						"allowedValues": []any{
							map[string]any{"id": "10", "value": "X"},
						},
					},
				},
			})))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	})
	// "epic" lower case should match "Epic".
	opts, err := f.GetFieldOptions(context.Background(), "customfield_10001", "PROJ", "epic")
	if err != nil {
		t.Fatal(err)
	}
	if len(opts) != 1 {
		t.Fatalf("expected 1 option, got %d", len(opts))
	}
}

func TestGetFieldOptionsIssueTypeNotFound(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jSON(map[string]any{
			"values": []map[string]any{
				{"id": "1", "name": "Bug"},
			},
		})))
	})
	_, err := f.GetFieldOptions(context.Background(), "customfield_10001", "PROJ", "Story")
	if err == nil {
		t.Fatal("expected error for missing issue type")
	}
}

// --- Service Desk (DC only) ---------------------------------------------

func TestGetServiceDesk(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"id": "4", "projectKey": "SUP"})))
	})
	sd, err := f.GetServiceDesk(context.Background(), "SUP")
	if err != nil {
		t.Fatal(err)
	}
	if sd["id"] != "4" {
		t.Errorf("id=%v", sd["id"])
	}
}

func TestGetRequestTypes(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"values": []map[string]any{{"id": "23", "name": "Get help"}}})))
	})
	rt, err := f.GetRequestTypes(context.Background(), "4", 0, 50)
	if err != nil {
		t.Fatal(err)
	}
	if rt["values"].([]any)[0].(map[string]any)["name"] != "Get help" {
		t.Errorf("rt=%v", rt)
	}
}

func TestGetRequestTypeFields(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"fields": []map[string]any{{"fieldId": "summary"}}})))
	})
	rf, err := f.GetRequestTypeFields(context.Background(), "4", "23")
	if err != nil {
		t.Fatal(err)
	}
	if rf == nil {
		t.Error("nil")
	}
}

func TestSearchAssignableUsers(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON([]map[string]any{{"accountId": "u1", "displayName": "Alice", "active": true}})))
	})
	users, err := f.SearchAssignableUsers(context.Background(), "alice", "PROJ", "", 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Errorf("len=%d", len(users))
	}
}

func matchPath(got, want string) bool { return got == want }
