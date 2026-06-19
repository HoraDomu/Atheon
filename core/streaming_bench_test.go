package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// BenchmarkRegularScanFile benchmarks the original ScanFile implementation
func BenchmarkRegularScanFile(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Some test content\nAWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\nMore content\n"

	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ScanFile(testFile)
	}
}

// BenchmarkChunkedScanFile benchmarks the new chunked scanning implementation
func BenchmarkChunkedScanFile(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Some test content\nAWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\nMore content\n"

	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ScanFileChunked(testFile, func(f Finding) {})
	}
}

// BenchmarkRegularScanFileLarge benchmarks regular scanning on a large file
func BenchmarkRegularScanFileLarge(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "large.txt")

	var content strings.Builder
	for i := 0; i < 10000; i++ {
		content.WriteString("Some content line with some text\n")
		if i%100 == 0 {
			content.WriteString("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n")
		}
	}

	if err := os.WriteFile(testFile, []byte(content.String()), 0o644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ScanFile(testFile)
	}
}

// BenchmarkChunkedScanFileLarge benchmarks chunked scanning on a large file
func BenchmarkChunkedScanFileLarge(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "large.txt")

	var content strings.Builder
	for i := 0; i < 10000; i++ {
		content.WriteString("Some content line with some text\n")
		if i%100 == 0 {
			content.WriteString("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n")
		}
	}

	if err := os.WriteFile(testFile, []byte(content.String()), 0o644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ScanFileChunked(testFile, func(f Finding) {})
	}
}

// BenchmarkRegularScanDir benchmarks regular directory scanning
func BenchmarkRegularScanDir(b *testing.B) {
	tmpDir := b.TempDir()

	// Create 20 test files
	for i := 0; i < 20; i++ {
		file := filepath.Join(tmpDir, filepath.Join("subdir", "file"+string(rune('0'+i%10))+".txt"))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			b.Fatal(err)
		}
		content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ScanDir(tmpDir)
	}
}

// BenchmarkStreamingScanDir benchmarks streaming directory scanning
func BenchmarkStreamingScanDir(b *testing.B) {
	tmpDir := b.TempDir()

	// Create 20 test files
	for i := 0; i < 20; i++ {
		file := filepath.Join(tmpDir, filepath.Join("subdir", "file"+string(rune('0'+i%10))+".txt"))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			b.Fatal(err)
		}
		content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ScanDirStreaming(tmpDir, func(f Finding) {})
	}
}

// BenchmarkMemoryUsageLargeFile compares memory usage for large files
func BenchmarkMemoryUsageLargeFile(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "large.txt")

	var content strings.Builder
	for i := 0; i < 50000; i++ {
		content.WriteString("Some content line with some more text and data\n")
		if i%100 == 0 {
			content.WriteString("AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n")
		}
	}

	if err := os.WriteFile(testFile, []byte(content.String()), 0o644); err != nil {
		b.Fatal(err)
	}

	b.Run("Regular", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = ScanFile(testFile)
		}
	})

	b.Run("Chunked", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = ScanFileChunked(testFile, func(f Finding) {})
		}
	})
}
