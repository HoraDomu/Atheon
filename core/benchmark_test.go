package core

import (
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkScanFile benchmarks scanning a single file
func BenchmarkScanFile(b *testing.B) {
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

// BenchmarkScanFileLarge benchmarks scanning a larger file
func BenchmarkScanFileLarge(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "large.txt")

	// Create a file with many lines and multiple patterns
	var content string
	for i := 0; i < 10000; i++ {
		content += "Some content line with some more text\n"
		if i%100 == 0 {
			content += "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n"
		}
		if i%100 == 1 {
			content += "api_key=sk-test1234567890abcdef\n"
		}
	}

	if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ScanFile(testFile)
	}
}

// BenchmarkScanString benchmarks scanning a string in memory
func BenchmarkScanString(b *testing.B) {
	content := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE api_key=sk-test1234567890abcdef"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ScanString(content, "test-source")
	}
}

// BenchmarkScanDir benchmarks directory scanning
func BenchmarkScanDir(b *testing.B) {
	tmpDir := b.TempDir()

	// Create multiple files
	for i := 0; i < 10; i++ {
		file := filepath.Join(tmpDir, filepath.Join("subdir", "sub"+string(rune('0'+i))+".txt"))
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
