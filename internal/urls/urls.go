// URL utilities — Data Center only. Cloud detection removed.
package urls

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

var blockedHostnames = map[string]bool{"localhost": true, "metadata.google.internal": true}

func ValidateURLForSSRF(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("empty URL")
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("missing hostname")
	}
	if blockedHostnames[host] {
		return fmt.Errorf("blocked: %s", host)
	}
	if ip := net.ParseIP(host); ip != nil {
		if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
			return fmt.Errorf("non-global IP: %s", ip)
		}
		return nil
	}
	allow := getAllowedDomains()
	if len(allow) > 0 {
		if !hostMatch(host, allow) {
			return fmt.Errorf("host %q not in MCP_ALLOWED_URL_DOMAINS", host)
		}
		return nil
	}
	if err := checkDNS(host); err != nil {
		return err
	}
	return nil
}

func getAllowedDomains() []string {
	v := os.Getenv("MCP_ALLOWED_URL_DOMAINS")
	if v == "" {
		return nil
	}
	out := []string{}
	for _, d := range strings.Split(v, ",") {
		d = strings.TrimSpace(d)
		if d != "" {
			out = append(out, d)
		}
	}
	return out
}

func hostMatch(host string, allow []string) bool {
	h := strings.ToLower(host)
	for _, e := range allow {
		e = strings.ToLower(strings.TrimSpace(e))
		if e == "" {
			continue
		}
		if h == e || strings.HasSuffix(h, "."+e) {
			return true
		}
	}
	return false
}

func checkDNS(host string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	addrs, err := (&net.Resolver{}).LookupIPAddr(ctx, host)
	if err != nil {
		return fmt.Errorf("DNS lookup %s: %w", host, err)
	}
	for _, a := range addrs {
		ip := a.IP
		if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
			return fmt.Errorf("%s → non-global %s", host, ip)
		}
	}
	return nil
}

func ResolveRelativeURL(rel, base string) string {
	if strings.HasPrefix(rel, "/") {
		return strings.TrimRight(base, "/") + rel
	}
	return rel
}
