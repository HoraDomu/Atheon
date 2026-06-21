package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScanURL(t *testing.T) {
	origClient := scanURLClient
	t.Cleanup(func() { scanURLClient = origClient })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("token=sk-abcdefghijklmnopqrstuvwxyz\npassword=secret123\n"))
	}))
	defer server.Close()

	scanURLClient = server.Client()

	findings, stats, err := ScanURL(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) == 0 {
		t.Error("expected at least one finding")
	}
	if stats == nil {
		t.Error("expected non-nil stats")
		return
	}
	if stats.Files != 1 {
		t.Errorf("files = %d, want 1", stats.Files)
	}
}

func TestScanURLBinaryContentType(t *testing.T) {
	origClient := scanURLClient
	t.Cleanup(func() { scanURLClient = origClient })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("\x89PNG\r\n\x1a\n")) // minimal PNG header
	}))
	defer server.Close()

	scanURLClient = server.Client()

	_, _, err := ScanURL(context.Background(), server.URL)
	if err == nil {
		t.Error("expected error for binary Content-Type")
	}
}

func TestScanURLNon200(t *testing.T) {
	origClient := scanURLClient
	t.Cleanup(func() { scanURLClient = origClient })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	scanURLClient = server.Client()

	_, _, err := ScanURL(context.Background(), server.URL)
	if err == nil {
		t.Error("expected error for 404")
	}
}

func TestScanURLContextCancel(t *testing.T) {
	origClient := scanURLClient
	t.Cleanup(func() { scanURLClient = origClient })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hang until context is cancelled
		select {}
	}))
	defer server.Close()

	scanURLClient = server.Client()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, _, err := ScanURL(ctx, server.URL)
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

func TestScanURLJSONContentType(t *testing.T) {
	origClient := scanURLClient
	t.Cleanup(func() { scanURLClient = origClient })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Use a token long enough to match the openai-api-key pattern
		// (sk- prefix + at least 20 more chars).
		json.NewEncoder(w).Encode(map[string]string{"token": "sk-abcdefghijklmnopqrstuvwxyz"})
	}))
	defer server.Close()

	scanURLClient = server.Client()

	findings, _, err := ScanURL(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) == 0 {
		t.Error("expected findings from JSON content")
	}
}

func TestScanURLUnsupportedScheme(t *testing.T) {
	_, _, err := ScanURL(context.Background(), "ftp://example.com/test")
	if err == nil {
		t.Error("expected error for ftp scheme")
	}
}
