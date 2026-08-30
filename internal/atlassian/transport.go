package atlassian

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
)

func buildTransport(baseURL string, sslVerify bool, httpProxy, httpsProxy, noProxy, socksProxy string, wpadEnable bool, wpadURL string) http.RoundTripper {
	// socksProxy, wpadEnable, wpadURL are reserved for future SOCKS/WPAD support.
	_ = socksProxy
	_ = wpadEnable
	_ = wpadURL
	trusted := operatorTrustedHosts(baseURL)
	rt := baseTransport(sslVerify, httpProxy, httpsProxy, noProxy)
	rt = WithRateLimit(rt, config.GetFloatEnv("ATLASSIAN_RATE_LIMIT_RPS", 10))
	rt = WithRetry(rt, config.GetIntEnv("ATLASSIAN_MAX_RETRIES", 3))
	rt = WithSSRFGuard(rt, trusted)
	return rt
}

func baseTransport(sslVerify bool, httpProxy, httpsProxy, noProxy string) http.RoundTripper {
	transport := &http.Transport{
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if !sslVerify {
		transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true}
	} else {
		transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	if proxy := buildProxyFunc(httpProxy, httpsProxy, noProxy); proxy != nil {
		transport.Proxy = proxy
	} else {
		transport.Proxy = http.ProxyFromEnvironment
	}
	return transport
}

func buildProxyFunc(httpProxy, httpsProxy, noProxy string) func(*http.Request) (*url.URL, error) {
	_ = noProxy
	if httpProxy == "" && httpsProxy == "" {
		return nil
	}
	return func(req *http.Request) (*url.URL, error) {
		raw := httpProxy
		if req.URL.Scheme == "https" {
			raw = httpsProxy
		}
		if raw == "" {
			return nil, nil
		}
		return url.Parse(raw)
	}
}

func WithSSRFGuard(base http.RoundTripper, trustedHosts []string) http.RoundTripper {
	trusted := map[string]struct{}{}
	for _, h := range trustedHosts {
		trusted[strings.ToLower(h)] = struct{}{}
	}
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		host := req.URL.Hostname()
		if host == "" {
			return nil, fmt.Errorf("ssrf: empty host")
		}
		// Trusted hosts bypass DNS/IP checks.
		if _, isTrusted := trusted[strings.ToLower(host)]; isTrusted {
			return base.RoundTrip(req)
		}
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, fmt.Errorf("ssrf: lookup %s: %w", host, err)
		}
		for _, ip := range ips {
			if !isAcceptableIP(ip) {
				return nil, fmt.Errorf("ssrf: %s -> non-global %s", host, ip)
			}
		}
		return base.RoundTrip(req)
	})
}

func WithRetry(base http.RoundTripper, maxRetries int) http.RoundTripper {
	if maxRetries <= 0 {
		return base
	}
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		var resp *http.Response
		var err error
		for attempt := 0; attempt <= maxRetries; attempt++ {
			resp, err = base.RoundTrip(req.Clone(req.Context()))
			// Success: no error and not a retryable status (5xx or 429)
			if err == nil && resp.StatusCode < 500 && resp.StatusCode != 429 {
				return resp, nil
			}
			// Retry on errors or retryable status codes
			if attempt == maxRetries {
				break
			}
			delay := retryAfter(resp)
			if delay == 0 {
				delay = backoff(attempt)
			}
			select {
			case <-time.After(delay):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
		return resp, err
	})
}

func WithRateLimit(base http.RoundTripper, rate float64) http.RoundTripper {
	if rate <= 0 {
		return base
	}
	lim := newTokenBucket(rate)
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if !lim.acquire() {
			return nil, fmt.Errorf("atlassian: rate limit exceeded")
		}
		return base.RoundTrip(req)
	})
}

func backoff(attempt int) time.Duration {
	b := time.Duration(1<<uint(attempt)) * time.Second
	if b > 30*time.Second {
		b = 30 * time.Second
	}
	return b
}

func retryAfter(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if n, err := strconv.Atoi(v); err == nil {
		return time.Duration(n) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		return time.Until(t)
	}
	return 0
}

func operatorTrustedHosts(baseURL string) []string {
	out := []string{}
	if baseURL != "" {
		if u, err := url.Parse(baseURL); err == nil && u.Hostname() != "" {
			out = append(out, u.Hostname())
		}
	}
	for _, d := range strings.Split(os.Getenv("MCP_ALLOWED_URL_DOMAINS"), ",") {
		d = strings.TrimSpace(d)
		if d != "" {
			out = append(out, d)
		}
	}
	return out
}

func isAcceptableIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
		return false
	}
	return ip.IsGlobalUnicast()
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

type tokenBucket struct {
	mu                     sync.Mutex
	rate, capacity, tokens float64
	last                   time.Time
}

func newTokenBucket(rps float64) *tokenBucket {
	return &tokenBucket{rate: rps, capacity: rps, tokens: rps, last: time.Now()}
}
func (b *tokenBucket) acquire() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens = minF(b.capacity, b.tokens+elapsed*b.rate)
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
