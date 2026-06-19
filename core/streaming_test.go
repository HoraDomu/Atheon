package core

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestScanFileChunked(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Some test content\nAWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\nMore content\n"

	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var findings []Finding
	stats, err := ScanFileChunked(testFile, func(f Finding) {
		findings = append(findings, f)
	})

	if err != nil {
		t.Fatalf("ScanFileChunked failed: %v", err)
	}

	if stats.Files != 1 {
		t.Errorf("expected 1 file, got %d", stats.Files)
	}

	if len(findings) == 0 {
		t.Error("expected to find AWS key pattern")
	}

	foundAWSKey := false
	for _, finding := range findings {
		if finding.Pattern == "aws-access-key" {
			foundAWSKey = true
			if finding.File != testFile {
				t.Errorf("expected file %s, got %s", testFile, finding.File)
			}
		}
	}

	if !foundAWSKey {
		t.Error("did not find expected aws-access-key pattern")
	}
}

func TestScanFileChunkedLargeFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "large.txt")

	// Create a file with many lines
	var content strings.Builder
	for i := 0; i < 10000; i++ {
		content.WriteString("Some content line\n")
		if i%100 == 0 {
			content.WriteString("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n")
		}
	}

	if err := os.WriteFile(testFile, []byte(content.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	var findings []Finding
	stats, err := ScanFileChunked(testFile, func(f Finding) {
		findings = append(findings, f)
	})

	if err != nil {
		t.Fatalf("ScanFileChunked failed: %v", err)
	}

	if stats.Files != 1 {
		t.Errorf("expected 1 file, got %d", stats.Files)
	}

	if len(findings) < 100 {
		t.Errorf("expected at least 100 findings, got %d", len(findings))
	}
}

func TestScanFileChunkedWithIgnore(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE // atheon:ignore\nAWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"

	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var findings []Finding
	_, err := ScanFileChunked(testFile, func(f Finding) {
		findings = append(findings, f)
	})

	if err != nil {
		t.Fatalf("ScanFileChunked failed: %v", err)
	}

	// Only the second line should be found
	if len(findings) != 1 {
		t.Errorf("expected 1 finding (first line ignored), got %d", len(findings))
	}
}

func TestScanDirStreaming(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple test files
	for i := 0; i < 5; i++ {
		file := filepath.Join(tmpDir, filepath.Join("sub", "file"+string(rune('0'+i))+".txt"))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	var findings []Finding
	mu := &sync.Mutex{} // Note: sync.Mutex is in the test, not the package
	stats, err := ScanDirStreaming(tmpDir, func(f Finding) {
		mu.Lock()
		findings = append(findings, f)
		mu.Unlock()
	})

	if err != nil {
		t.Fatalf("ScanDirStreaming failed: %v", err)
	}

	if stats.Files < 5 {
		t.Errorf("expected at least 5 files, got %d", stats.Files)
	}

	if len(findings) < 5 {
		t.Errorf("expected at least 5 findings, got %d", len(findings))
	}
}

func TestScanReader(t *testing.T) {
	content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\napi_key=sk-test1234567890abcdef\n"

	var findings []Finding
	err := ScanReader(strings.NewReader(content), "stdin", func(f Finding) {
		findings = append(findings, f)
	})

	if err != nil {
		t.Fatalf("ScanReader failed: %v", err)
	}

	if len(findings) < 2 {
		t.Errorf("expected at least 2 findings, got %d", len(findings))
	}

	// Check findings came from "stdin"
	for _, f := range findings {
		if f.File != "stdin" {
			t.Errorf("expected file 'stdin', got '%s'", f.File)
		}
	}
}

func TestScanReaderWithIgnore(t *testing.T) {
	content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE // atheon:ignore\nAWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"

	var findings []Finding
	err := ScanReader(strings.NewReader(content), "stdin", func(f Finding) {
		findings = append(findings, f)
	})

	if err != nil {
		t.Fatalf("ScanReader failed: %v", err)
	}

	// Only the second line should be found
	if len(findings) != 1 {
		t.Errorf("expected 1 finding (first line ignored), got %d", len(findings))
	}
}

func TestScanReaderEmpty(t *testing.T) {
	var findings []Finding
	err := ScanReader(strings.NewReader(""), "empty", func(f Finding) {
		findings = append(findings, f)
	})

	if err != nil {
		t.Fatalf("ScanReader failed: %v", err)
	}

	if len(findings) != 0 {
		t.Errorf("expected 0 findings for empty input, got %d", len(findings))
	}
}

func TestProfiledScan(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"

	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var findings []Finding
	stats, err := ProfiledScan(testFile, func(f Finding) {
		findings = append(findings, f)
	})

	if err != nil {
		t.Fatalf("ProfiledScan failed: %v", err)
	}

	if stats.Files != 1 {
		t.Errorf("expected 1 file, got %d", stats.Files)
	}

	if len(findings) == 0 {
		t.Error("expected to find AWS key pattern")
	}
}

func TestStreamingVsRegularConsistency(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	content := "Some content\n"
	content += "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"
	content += "api_key=sk-test1234567890abcdef\n"
	content += "export TOKEN=AKIAIOSFODNN7EXAMPLE\n"

	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// Get results using regular ScanFile
	regularFindings, _, _ := ScanFile(testFile)

	// Get results using streaming ScanFileChunked
	var streamingFindings []Finding
	_, err := ScanFileChunked(testFile, func(f Finding) {
		streamingFindings = append(streamingFindings, f)
	})

	if err != nil {
		t.Fatalf("ScanFileChunked failed: %v", err)
	}

	// Both should find the same number of patterns
	if len(regularFindings) != len(streamingFindings) {
		t.Errorf("regular scan found %d findings, streaming found %d",
			len(regularFindings), len(streamingFindings))
	}

	// Check pattern names match
	regularPatterns := make(map[string]bool)
	streamingPatterns := make(map[string]bool)
	for _, f := range regularFindings {
		regularPatterns[f.Pattern] = true
	}
	for _, f := range streamingFindings {
		streamingPatterns[f.Pattern] = true
	}

	for pattern := range regularPatterns {
		if !streamingPatterns[pattern] {
			t.Errorf("streaming scan missing pattern: %s", pattern)
		}
	}
}
