// JiraConfig — Data Center with PAT-only authentication.
// Reads JIRA_URL and JIRA_PERSONAL_TOKEN from env.

package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

type JiraConfig struct {
	URL                  string
	PersonalToken        string
	SSLVerify            bool
	HTTPProxy            string
	HTTPSProxy           string
	NoProxy              string
	SOCKSProxy           string
	ProxyWPADEnable      bool
	ProxyWPADURL         string
	ProjectsFilter       string
	PassthroughHeaders   []string
	CustomHeaders        map[string]string
	InternalOnlyProjects []string
}

func (c *JiraConfig) IsAuthConfigured() bool {
	return c.URL != "" && c.PersonalToken != ""
}

func JiraConfigFromEnv() *JiraConfig {
	cfg := &JiraConfig{
		URL:             firstNonEmpty(os.Getenv("JIRA_URL"), "https://jira.mts.ru/"),
		PersonalToken:   os.Getenv("JIRA_PERSONAL_TOKEN"),
		SSLVerify:       IsEnvSSLVerify("JIRA_SSL_VERIFY"),
		HTTPProxy:       os.Getenv("JIRA_HTTP_PROXY"),
		HTTPSProxy:      os.Getenv("JIRA_HTTPS_PROXY"),
		NoProxy:         firstNonEmpty(os.Getenv("JIRA_NO_PROXY"), os.Getenv("NO_PROXY"), os.Getenv("no_proxy")),
		SOCKSProxy:      os.Getenv("JIRA_SOCKS_PROXY"),
		ProxyWPADEnable: IsEnvTruthy("JIRA_PROXY_WPAD_ENABLE"),
		ProxyWPADURL:    firstNonEmpty(os.Getenv("JIRA_PROXY_WPAD_URL"), os.Getenv("ATLASSIAN_PROXY_WPAD_URL")),
		ProjectsFilter:  os.Getenv("JIRA_PROJECTS_FILTER"),
		CustomHeaders:   GetCustomHeaders("JIRA_CUSTOM_HEADERS"),
	}
	cfg.PassthroughHeaders = GetHeaderNames("JIRA_PASSTHROUGH_HEADERS")
	if v := os.Getenv("JIRA_INTERNAL_ONLY_PROJECTS"); v != "" {
		cfg.InternalOnlyProjects = splitCSV(v)
	}
	return cfg
}

func (c *JiraConfig) BaseURL() string { return strings.TrimRight(c.URL, "/") }
func (c *JiraConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("JIRA_URL is required")
	}
	if _, err := url.Parse(c.URL); err != nil {
		return fmt.Errorf("JIRA_URL invalid: %w", err)
	}
	if c.PersonalToken == "" {
		slog.Warn("JIRA_PERSONAL_TOKEN not set — Jira tools will be unavailable")
	}
	return nil
}

func splitCSV(v string) []string {
	out := []string{}
	for _, s := range strings.Split(v, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
