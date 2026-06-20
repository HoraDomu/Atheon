package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	scanURLTimeout   = 10 * time.Second
	scanURLMaxBytes  = 5 << 20 // 5 MB
	scanURLUserAgent = "atheon/0.4.0"
)

// scanURLClient is configurable so tests can inject a mock.
var scanURLClient = &http.Client{
	Timeout: scanURLTimeout,
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse // do not follow redirects
	},
}

// ScanURL fetches a URL and scans the response body for pattern matches.
// It returns all findings, scan stats, and any error encountered.
//
// The response body is streamed and capped at scanURLMaxBytes (5 MB).
// Only text/* content types are scanned; binary types are skipped silently.
// The context controls cancellation; if canceled the scan may be partial.
func ScanURL(ctx context.Context, rawURL string) ([]Finding, *Stats, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, nil, fmt.Errorf("unsupported URL scheme: %s (only http/https supported)", u.Scheme)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", scanURLUserAgent)

	resp, err := scanURLClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	contentType := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))
	if !strings.HasPrefix(contentType, "text/") &&
		!strings.Contains(contentType, "application/json") &&
		!strings.Contains(contentType, "application/xml") &&
		!strings.Contains(contentType, "application/javascript") {
		return nil, nil, fmt.Errorf("unsupported Content-Type: %s (only text/* scanned)", contentType)
	}

	// Wrap with size-limited reader so a malicious server cannot exhaust memory.
	limited := &io.LimitedReader{R: resp.Body, N: scanURLMaxBytes + 1}
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, nil, fmt.Errorf("read failed: %w", err)
	}
	if limited.N <= 0 {
		return nil, nil, fmt.Errorf("response body exceeds %d MB limit", scanURLMaxBytes>>20)
	}

	findings := ScanString(ctx, string(body), u.String())
	return findings, &Stats{
		Files:     1,
		Bytes:     int64(len(body)),
		ElapsedMs: 0, // no timing in ScanURL path yet
	}, nil
}
