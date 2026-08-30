package confluence

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/atlassian"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
)

func newFetcher(t *testing.T, handler http.HandlerFunc) *Fetcher {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	cfg := &config.ConfluenceConfig{URL: srv.URL, PersonalToken: "test-pat"}
	client, err := atlassian.NewConfluence(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return New(cfg, client)
}

func jSON(v any) string { b, _ := json.Marshal(v); return string(b) }

func TestGetPage(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jSON(map[string]any{"id": "123", "type": "page", "title": "Test"})))
	})
	p, err := f.GetPage(context.Background(), "123", nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Test" {
		t.Errorf("title=%s", p.Title)
	}
}

func TestSearch(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "1", "title": "R1"}}, "size": 1})))
	})
	r, err := f.Search(context.Background(), `text ~ "x"`, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 1 {
		t.Errorf("len=%d", len(r))
	}
}

func TestGetPageComments(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{
			{"id": "c1", "body": map[string]any{"view": map[string]any{"value": "hello"}}, "created": "..."},
			{"id": "c2", "body": map[string]any{"view": map[string]any{"value": "reply"}}, "created": "...",
				"ancestors": []any{map[string]any{"id": "page1", "type": "page"}, map[string]any{"id": "c1", "type": "comment"}}},
		}})))
	})
	comments, err := f.GetPageComments(context.Background(), "123")
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 2 {
		t.Errorf("len=%d", len(comments))
	}
	if comments[0].ParentID != "" {
		t.Errorf("top-level parent_id=%q, want empty", comments[0].ParentID)
	}
	if comments[1].ParentID != "c1" {
		t.Errorf("reply parent_id=%q, want c1", comments[1].ParentID)
	}
	d := comments[1].ToSimplifiedDict()
	if d["parent_id"] != "c1" {
		t.Errorf("dict parent_id=%v, want c1", d["parent_id"])
	}
}

func TestAddComment(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"id": "c1", "body": map[string]any{}})))
	})
	c, err := f.AddComment(context.Background(), "123", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != "c1" {
		t.Errorf("id=%s", c.ID)
	}
}

func TestGetPageLabels(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"prefix": "global", "name": "test"}}})))
	})
	labels, err := f.GetPageLabels(context.Background(), "123", "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 1 {
		t.Errorf("len=%d", len(labels))
	}
}

func TestAddPageLabel(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"name": "test"}}})))
	})
	labels, err := f.AddPageLabel(context.Background(), "123", "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 1 {
		t.Errorf("len=%d", len(labels))
	}
}

func TestSearchUser(t *testing.T) {
	callCount := 0
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/group/") {
			callCount++
			w.Write([]byte(jSON(map[string]any{"results": []map[string]any{}, "size": 0}))) // empty to stop pagination
			return
		}
		// Fallback for non-group paths
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{}, "size": 0})))
	})
	users, err := f.SearchUser(context.Background(), `user.fullname ~ "alice"`, 10, "confluence-users")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 0 {
		t.Errorf("len=%d", len(users))
	}
	_ = callCount
}

func stringsContains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestGetPageChildren(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "c1", "title": "Child"}}})))
	})
	children, err := f.GetPageChildren(context.Background(), "123", 0, 25, false, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(children) != 1 {
		t.Errorf("len=%d", len(children))
	}
}

func TestGetSpacePageTree(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "1", "title": "Home", "ancestors": []any{}}}, "size": 1})))
	})
	tree, err := f.GetSpacePageTree(context.Background(), "DEV", 100)
	if err != nil {
		t.Fatal(err)
	}
	if tree["total_pages"].(int) != 1 {
		t.Errorf("total=%v", tree["total_pages"])
	}
}

func TestGetContentAttachments(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "att1", "title": "file.pdf"}}, "size": 1})))
	})
	atts, err := f.GetContentAttachments(context.Background(), "123", 0, 50, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if atts["total"].(int) != 1 {
		t.Errorf("total=%v", atts["total"])
	}
}

func TestDeleteAttachment(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method=%s", r.Method)
		}
		w.WriteHeader(204)
	})
	if err := f.DeleteAttachment(context.Background(), "att1"); err != nil {
		t.Fatal(err)
	}
}

func TestGetSpaces(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"key": "DEV", "name": "Development"}}})))
	})
	spaces, err := f.GetSpaces(context.Background(), 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(spaces) != 1 {
		t.Errorf("len=%d", len(spaces))
	}
}

func TestScanContent(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/scan" {
			t.Errorf("path=%s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("method=%s", r.Method)
		}
		q := r.URL.Query()
		if q.Get("cql") != "type = page" {
			t.Errorf("cql=%s", q.Get("cql"))
		}
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "1"}}, "size": 1, "start": 0, "limit": 50})))
	})
	r, err := f.ScanContent(context.Background(), "type = page", 0, 50, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r["total"].(int) != 1 {
		t.Errorf("total=%v", r["total"])
	}
}

func TestGetContentHistory(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/123/history" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"createdBy": map[string]any{"displayName": "Alice"}, "createdDate": "2024-01-01"})))
	})
	h, err := f.GetContentHistory(context.Background(), "123", nil)
	if err != nil {
		t.Fatal(err)
	}
	if h["createdBy"] == nil {
		t.Error("createdBy missing")
	}
}

func TestGetRestrictionForOperation(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/123/restriction/byOperation/read" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"restrictions": map[string]any{}})))
	})
	r, err := f.GetRestrictionForOperation(context.Background(), "123", "read")
	if err != nil {
		t.Fatal(err)
	}
	if r["restrictions"] == nil {
		t.Error("restrictions missing")
	}
}

func TestGetContentChildren(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/123/child/page" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "c1"}}})))
	})
	children, err := f.GetContentChildren(context.Background(), "123", "page", 0, 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(children) != 1 {
		t.Errorf("len=%d", len(children))
	}
}

func TestGetDescendants(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/123/descendant" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "d1"}}, "size": 1})))
	})
	r, err := f.GetDescendants(context.Background(), "123", 0, 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r["total"].(int) != 1 {
		t.Errorf("total=%v", r["total"])
	}
}

func TestGetDescendantsOfType(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/123/descendant/page" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "d1"}}, "size": 1})))
	})
	r, err := f.GetDescendantsOfType(context.Background(), "123", "page", 0, 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r["total"].(int) != 1 {
		t.Errorf("total=%v", r["total"])
	}
}

func TestGetContentProperty(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/123/property/my-key" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"key": "my-key", "value": "val", "id": 1, "version": map[string]any{"number": 1}})))
	})
	p, err := f.GetContentProperty(context.Background(), "123", "my-key")
	if err != nil {
		t.Fatal(err)
	}
	if p["key"] != "my-key" {
		t.Errorf("key=%v", p["key"])
	}
}

func TestSearchEntities(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/search" {
			t.Errorf("path=%s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("cql") != "type = page" {
			t.Errorf("cql=%s", q.Get("cql"))
		}
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "1"}}, "size": 1, "start": 0, "limit": 50})))
	})
	r, err := f.SearchEntities(context.Background(), "type = page", "", "", 0, 50, false)
	if err != nil {
		t.Fatal(err)
	}
	if r["total"].(int) != 1 {
		t.Errorf("total=%v", r["total"])
	}
}

func TestGetSpace(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/space/DEV" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"key": "DEV", "name": "Development"})))
	})
	s, err := f.GetSpace(context.Background(), "DEV", nil)
	if err != nil {
		t.Fatal(err)
	}
	if s["key"] != "DEV" {
		t.Errorf("key=%v", s["key"])
	}
}

func TestGetSpaceContent(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/space/DEV/content/page" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"id": "1"}}, "size": 1})))
	})
	r, err := f.GetSpaceContent(context.Background(), "DEV", "page", 0, nil, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if r["total"].(int) != 1 {
		t.Errorf("total=%v", r["total"])
	}
}

func TestGetSpaceProperties(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/space/DEV/property" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"results": []map[string]any{{"key": "k1"}}})))
	})
	props, err := f.GetSpaceProperties(context.Background(), "DEV", nil, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(props) != 1 {
		t.Errorf("len=%d", len(props))
	}
}

func TestGetSpaceProperty(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/space/DEV/property/k1" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"key": "k1", "value": "v1"})))
	})
	p, err := f.GetSpaceProperty(context.Background(), "DEV", "k1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if p["key"] != "k1" {
		t.Errorf("key=%v", p["key"])
	}
}

func TestGetCurrentUser(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/user/current" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"username": "alice", "displayName": "Alice"})))
	})
	u, err := f.GetCurrentUser(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if u["username"] != "alice" {
		t.Errorf("username=%v", u["username"])
	}
}

func TestGetUser(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/user/alice" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"username": "alice", "displayName": "Alice"})))
	})
	u, err := f.GetUser(context.Background(), "alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	if u["username"] != "alice" {
		t.Errorf("username=%v", u["username"])
	}
}

func TestRemovePageLabelByName(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Path != "/rest/api/content/123/label" {
			t.Errorf("path=%s", r.URL.Path)
		}
		if r.URL.Query().Get("name") != "mylabel" {
			t.Errorf("name=%s", r.URL.Query().Get("name"))
		}
		w.WriteHeader(204)
	})
	if err := f.RemovePageLabelByName(context.Background(), "123", "mylabel"); err != nil {
		t.Fatal(err)
	}
}

func TestRemovePageLabel(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Path != "/rest/api/content/123/label/global/mylabel" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.WriteHeader(204)
	})
	if err := f.RemovePageLabel(context.Background(), "123", "global/mylabel"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateContentProperty(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Path != "/rest/api/content/123/property" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"key": "k1", "value": "v1", "id": 1})))
	})
	p, err := f.CreateContentProperty(context.Background(), "123", "k1", "v1")
	if err != nil {
		t.Fatal(err)
	}
	if p["key"] != "k1" {
		t.Errorf("key=%v", p["key"])
	}
}

func TestUpdateContentProperty(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Path != "/rest/api/content/123/property/k1" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"key": "k1", "value": "v2", "id": 1, "version": map[string]any{"number": 2}})))
	})
	p, err := f.UpdateContentProperty(context.Background(), "123", "k1", 1, 1, "v2")
	if err != nil {
		t.Fatal(err)
	}
	if p["value"] != "v2" {
		t.Errorf("value=%v", p["value"])
	}
}

func TestDeleteContentProperty(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Path != "/rest/api/content/123/property/k1" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.WriteHeader(204)
	})
	if err := f.DeleteContentProperty(context.Background(), "123", "k1"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateAttachmentMetadata(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Path != "/rest/api/content/att1" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Write([]byte(jSON(map[string]any{"id": "att1", "title": "new-title"})))
	})
	r, err := f.UpdateAttachmentMetadata(context.Background(), "att1", "new-title", "")
	if err != nil {
		t.Fatal(err)
	}
	if r["title"] != "new-title" {
		t.Errorf("title=%v", r["title"])
	}
}

func TestUpdateAttachmentData(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Path != "/rest/api/content/att1/data" {
			t.Errorf("path=%s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "multipart/form-data; boundary=" && !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("content-type=%s", ct)
		}
		w.Write([]byte(jSON(map[string]any{"id": "att1"})))
	})
	r, err := f.UpdateAttachmentData(context.Background(), "att1", "file.txt", []byte("hello"), "")
	if err != nil {
		t.Fatal(err)
	}
	if r["id"] != "att1" {
		t.Errorf("id=%v", r["id"])
	}
}

func TestUpdatePageSection_Markdown(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/rest/api/content/123":
			w.Write([]byte(jSON(map[string]any{
				"id": "123", "type": "page", "title": "Test Page",
				"space":   map[string]any{"key": "DEV"},
				"version": map[string]any{"number": 1.0},
				"body":    map[string]any{"storage": map[string]any{"value": "<h1>Intro</h1><p>Old text</p><h2>Section A</h2><p>Replace me</p><h2>Section B</h2><p>Keep</p>", "representation": "storage"}},
			})))
		case r.Method == "PUT":
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			storage, _ := body["body"].(map[string]any)["storage"].(map[string]any)
			val, _ := storage["value"].(string)
			// Markdown **bold** and [link](url) must be converted to HTML.
			if !strings.Contains(val, "<strong>bold</strong>") {
				t.Errorf("expected <strong>bold</strong> in %q", val)
			}
			if !strings.Contains(val, `<a href="https://example.com">link</a>`) {
				t.Errorf("expected markdown link converted to <a> in %q", val)
			}
			// Literal markdown syntax must NOT appear.
			if strings.Contains(val, "**bold**") {
				t.Errorf("literal markdown found in %q", val)
			}
			// Surrounding sections must be preserved.
			if !strings.Contains(val, "Intro") || !strings.Contains(val, "Keep") {
				t.Errorf("surrounding content lost: %q", val)
			}
			w.Write([]byte(jSON(map[string]any{"id": "123"})))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	r, err := f.UpdatePageSection(context.Background(), "123", "Section A", "**bold** [link](https://example.com)", "markdown", false, "")
	if err != nil {
		t.Fatal(err)
	}
	if r["heading_level"] != 2 {
		t.Errorf("heading_level=%v", r["heading_level"])
	}
}

func TestUpdatePageSection_Storage(t *testing.T) {
	f := newFetcher(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/rest/api/content/123":
			w.Write([]byte(jSON(map[string]any{
				"id": "123", "type": "page", "title": "Test Page",
				"space":   map[string]any{"key": "DEV"},
				"version": map[string]any{"number": 1.0},
				"body":    map[string]any{"storage": map[string]any{"value": "<h2>Sec</h2><p>old</p>", "representation": "storage"}},
			})))
		case r.Method == "PUT":
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			storage, _ := body["body"].(map[string]any)["storage"].(map[string]any)
			val, _ := storage["value"].(string)
			if !strings.Contains(val, "<custom:tag>raw</custom:tag>") {
				t.Errorf("raw storage not preserved: %q", val)
			}
			w.Write([]byte(jSON(map[string]any{"id": "123"})))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	_, err := f.UpdatePageSection(context.Background(), "123", "Sec", "<custom:tag>raw</custom:tag>", "storage", false, "")
	if err != nil {
		t.Fatal(err)
	}
}
