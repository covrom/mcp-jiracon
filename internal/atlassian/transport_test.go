package atlassian

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
)

func TestNewWithPAT(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token-123" {
			t.Errorf("missing PAT header: %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer srv.Close()

	cfg := &config.JiraConfig{URL: srv.URL, PersonalToken: "test-token-123"}
	c, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	var raw map[string]any
	if err := c.Get(context.Background(), "/path", &raw); err != nil {
		t.Fatal(err)
	}
	if raw["ok"] != true {
		t.Error("unexpected response")
	}
	t.Log("PAT auth header verified ✅")
}

func TestNewWithoutPAT(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	cfg := &config.JiraConfig{URL: srv.URL, PersonalToken: ""}
	_, err := New(cfg)
	if err == nil {
		t.Error("expected error for empty PAT")
	}
	t.Logf("empty PAT rejected: %v", err)
}

func TestTLSInsecureVerify(t *testing.T) {
	tlsSrv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	}))
	defer tlsSrv.Close()

	t.Run("verify=false succeeds on self-signed", func(t *testing.T) {
		cfg := &config.JiraConfig{URL: tlsSrv.URL, PersonalToken: "tok", SSLVerify: false}
		c, err := New(cfg)
		if err != nil {
			t.Fatal(err)
		}
		var raw map[string]any
		if err := c.Get(context.Background(), "/", &raw); err != nil {
			t.Fatal(err)
		}
		t.Log("self-signed TLS with verify=false ✅")
	})

	// Env var integration: JIRA_SSL_VERIFY=false → SSLVerify=false.
	t.Setenv("JIRA_URL", tlsSrv.URL)
	t.Setenv("JIRA_PERSONAL_TOKEN", "tok")
	t.Setenv("JIRA_SSL_VERIFY", "false")
	cfg := config.JiraConfigFromEnv()
	if cfg.SSLVerify {
		t.Fatal("SSLVerify should be false")
	}
	c, _ := New(cfg)
	var raw map[string]any
	if err := c.Get(context.Background(), "/", &raw); err != nil {
		t.Fatal(err)
	}
	t.Logf("JIRA_SSL_VERIFY=false → SSLVerify=false ✅")
}

func TestConfluencePAT(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer conf-token" {
			t.Errorf("unexpected auth: %s", r.Header.Get("Authorization"))
		}
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cfg := &config.ConfluenceConfig{URL: srv.URL, PersonalToken: "conf-token"}
	c, err := NewConfluence(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Get(context.Background(), "/path", nil); err != nil {
		t.Fatal(err)
	}
	t.Log("Confluence PAT auth ✅")
}
