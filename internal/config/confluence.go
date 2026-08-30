// ConfluenceConfig — Data Center with PAT-only authentication.
// Reads CONFLUENCE_URL and CONFLUENCE_PERSONAL_TOKEN from env.

package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

type ConfluenceConfig struct {
	URL                string
	PersonalToken      string
	SSLVerify          bool
	HTTPProxy          string
	HTTPSProxy         string
	NoProxy            string
	SOCKSProxy         string
	ProxyWPADEnable    bool
	ProxyWPADURL       string
	SpacesFilter       string
	PassthroughHeaders []string
	CustomHeaders      map[string]string
}

func (c *ConfluenceConfig) IsAuthConfigured() bool {
	return c.URL != "" && c.PersonalToken != ""
}

func ConfluenceConfigFromEnv() *ConfluenceConfig {
	cfg := &ConfluenceConfig{
		URL:             firstNonEmpty(os.Getenv("CONFLUENCE_URL"), "https://confluence.mts.ru/"),
		PersonalToken:   os.Getenv("CONFLUENCE_PERSONAL_TOKEN"),
		SSLVerify:       IsEnvSSLVerify("CONFLUENCE_SSL_VERIFY"),
		HTTPProxy:       os.Getenv("CONFLUENCE_HTTP_PROXY"),
		HTTPSProxy:      os.Getenv("CONFLUENCE_HTTPS_PROXY"),
		NoProxy:         firstNonEmpty(os.Getenv("CONFLUENCE_NO_PROXY"), os.Getenv("NO_PROXY"), os.Getenv("no_proxy")),
		SOCKSProxy:      os.Getenv("CONFLUENCE_SOCKS_PROXY"),
		ProxyWPADEnable: IsEnvTruthy("CONFLUENCE_PROXY_WPAD_ENABLE"),
		ProxyWPADURL:    firstNonEmpty(os.Getenv("CONFLUENCE_PROXY_WPAD_URL"), os.Getenv("ATLASSIAN_PROXY_WPAD_URL")),
		SpacesFilter:    os.Getenv("CONFLUENCE_SPACES_FILTER"),
		CustomHeaders:   GetCustomHeaders("CONFLUENCE_CUSTOM_HEADERS"),
	}
	cfg.PassthroughHeaders = GetHeaderNames("CONFLUENCE_PASSTHROUGH_HEADERS")
	return cfg
}

func (c *ConfluenceConfig) BaseURL() string   { return strings.TrimRight(c.URL, "/") }
func (c *ConfluenceConfig) V1BaseURL() string { return c.BaseURL() + "/rest/api" }

func (c *ConfluenceConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("CONFLUENCE_URL is required")
	}
	if _, err := url.Parse(c.URL); err != nil {
		return fmt.Errorf("CONFLUENCE_URL invalid: %w", err)
	}
	if c.PersonalToken == "" {
		slog.Warn("CONFLUENCE_PERSONAL_TOKEN not set")
	}
	return nil
}
