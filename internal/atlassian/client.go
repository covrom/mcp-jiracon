// Atlassian HTTP client — PAT-only authentication.
package atlassian

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/utils"
)

type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func New(cfg *config.JiraConfig) (*Client, error) {
	return NewWith(cfg.BaseURL(), cfg.PersonalToken, cfg.SSLVerify,
		cfg.HTTPProxy, cfg.HTTPSProxy, cfg.NoProxy, cfg.SOCKSProxy,
		cfg.ProxyWPADEnable, cfg.ProxyWPADURL)
}

func NewConfluence(cfg *config.ConfluenceConfig) (*Client, error) {
	return NewWith(cfg.BaseURL(), cfg.PersonalToken, cfg.SSLVerify,
		cfg.HTTPProxy, cfg.HTTPSProxy, cfg.NoProxy, cfg.SOCKSProxy,
		cfg.ProxyWPADEnable, cfg.ProxyWPADURL)
}

func NewWith(baseURL, token string, sslVerify bool,
	httpProxy, httpsProxy, noProxy, socksProxy string,
	wpadEnable bool, wpadURL string,
) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("atlassian: empty base URL")
	}
	transport := buildTransport(baseURL, sslVerify, httpProxy, httpsProxy, noProxy, socksProxy, wpadEnable, wpadURL)
	if token == "" {
		return nil, fmt.Errorf("atlassian: PAT token is required")
	}
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Token:   token,
		HTTP:    &http.Client{Transport: transport, Timeout: 30 * time.Second},
	}, nil
}

func (c *Client) NewRequest(ctx context.Context, method, path string, body []byte) (*http.Request, error) {
	var rdr io.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, path, rdr)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) applyAuth(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", utils.DefaultUserAgent())
}

func (c *Client) Do(req *http.Request, out any) error {
	c.applyAuth(req)
	if req.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if req.URL.Scheme == "" {
		full := c.BaseURL + req.URL.Path
		if req.URL.RawQuery != "" {
			full += "?" + req.URL.RawQuery
		}
		u, err := url.Parse(full)
		if err != nil {
			return err
		}
		req.URL = u
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return &HTTPError{StatusCode: resp.StatusCode, Body: string(b), URL: req.URL.String()}
	}
	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) Get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return fmt.Errorf("http new request: %w", err)
	}
	return c.Do(req, out)
}
func (c *Client) Post(ctx context.Context, path string, body, out any) error {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", path, &buf)
	if err != nil {
		return fmt.Errorf("http new request: %w", err)
	}
	return c.Do(req, out)
}
func (c *Client) Put(ctx context.Context, path string, body, out any) error {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", path, &buf)
	if err != nil {
		return fmt.Errorf("http new request: %w", err)
	}
	return c.Do(req, out)
}
func (c *Client) Delete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("http new request: %w", err)
	}
	return c.Do(req, nil)
}
func (c *Client) PostMultipart(ctx context.Context, path string, fields map[string]string, fileName string, fileBytes []byte, out any) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		w.WriteField(k, v)
	}
	if fileName != "" && fileBytes != nil {
		fw, _ := w.CreateFormFile("file", fileName)
		fw.Write(fileBytes)
	}
	w.Close()
	req, err := http.NewRequestWithContext(ctx, "POST", path, &buf)
	if err != nil {
		return fmt.Errorf("http new request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-Atlassian-Token", "no-check")
	return c.Do(req, out)
}
func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("http new request: %w", err)
	}
	c.applyAuth(req)
	if req.URL.Scheme == "" {
		full := c.BaseURL + req.URL.Path
		if req.URL.RawQuery != "" {
			full += "?" + req.URL.RawQuery
		}
		u, perr := url.Parse(full)
		if perr != nil {
			return nil, perr
		}
		req.URL = u
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: string(b), URL: req.URL.String()}
	}
	return io.ReadAll(resp.Body)
}

type HTTPError struct {
	StatusCode int
	Body, URL  string
}

func (e *HTTPError) Error() string { return fmt.Sprintf("HTTP %d %s: %s", e.StatusCode, e.URL, e.Body) }
