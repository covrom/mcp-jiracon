package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.services.mts.ru/rstsovan/mcp-atlassian/internal/config"
)

// IsReadOnlyMode returns whether the server is in read-only mode.
// Reads READ_ONLY_MODE; default false.
func IsReadOnlyMode() bool {
	return config.IsEnvExtendedTruthy("READ_ONLY_MODE", "false")
}

// ValidateSafePath ensures the given path resolves inside the current
// working directory (or base_dir if provided) to defeat path traversal.
// Mirrors utils.io.validate_safe_path.
func ValidateSafePath(p string, baseDir ...string) (string, error) {
	base := ""
	if len(baseDir) > 0 && baseDir[0] != "" {
		base = baseDir[0]
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getwd: %w", err)
		}
		base = wd
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(base, p)
	}
	resolved, err := filepath.EvalSymlinks(p)
	if err != nil {
		// EvalSymlinks fails for non-existent files; fall back to
		// Abs+Clean to still produce an absolute, cleaned path.
		abs, absErr := filepath.Abs(p)
		if absErr != nil {
			return "", absErr
		}
		if !strings.HasPrefix(abs, base+string(filepath.Separator)) && abs != base {
			return "", fmt.Errorf("path %q escapes base directory %q", p, base)
		}
		return abs, nil
	}
	if !strings.HasPrefix(resolved, base+string(filepath.Separator)) && resolved != base {
		return "", fmt.Errorf("path %q escapes base directory %q", p, base)
	}
	return resolved, nil
}

// ReadFileLines reads a file line-by-line and returns a slice of
// trimmed non-empty lines. Useful for parsing CSV-like config files.
func ReadFileLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	out := []string{}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			out = append(out, line)
		}
	}
	return out, sc.Err()
}
