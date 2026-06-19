package core

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FindingCallback is called for each finding as it's discovered
type FindingCallback func(finding Finding)

// ScanFileChunked scans a file line-by-line using buffered reading to avoid
// loading the entire file into memory. Each finding is passed to the callback.
// Returns statistics about the scan.
func ScanFileChunked(path string, callback FindingCallback) (*Stats, error) {
	start := time.Now()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bytesCount int64
	var findingsCount int

	scanner := bufio.NewScanner(file)
	// Increase buffer size for very long lines (up to 1MB)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		bytesCount += int64(len(line)) + 1 // +1 for newline

		// Check for ignore directive
		if strings.Contains(line, "atheon:ignore") {
			lineNum++
			continue
		}

		// Check each category scanner
		for _, cs := range activeScanners {
			if !cs.combined.MatchString(line) {
				continue
			}
			for _, p := range cs.patterns {
				if p.Matches(line) {
					callback(Finding{
						Pattern: p.Name(),
						File:    path,
						Line:    lineNum + 1,
						Content: strings.TrimSpace(line),
					})
					findingsCount++
				}
			}
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Stats{
		Files:     1,
		Bytes:     bytesCount,
		ElapsedMs: time.Since(start).Milliseconds(),
	}, nil
}

// ScanDirStreaming scans a directory and streams findings via callback instead
// of buffering all results in memory. Uses concurrency for better performance.
func ScanDirStreaming(root string, callback FindingCallback) (*Stats, error) {
	start := time.Now()
	ignoreMatcher := loadIgnorePatternsMatcher(root)

	// First pass: collect file paths
	var paths []string
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if d.IsDir() {
			if skipDirs[d.Name()] || isIgnored(rel, ignoreMatcher) {
				return filepath.SkipDir
			}
			return nil
		}
		if isIgnored(rel, ignoreMatcher) {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if !binaryExts[ext] {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Scan files with controlled concurrency
	var wg sync.WaitGroup
	sem := make(chan struct{}, 32) // Limit concurrent file operations
	var mu sync.Mutex
	var totalBytes int64
	var filesScanned int

	for _, p := range paths {
		wg.Add(1)
		sem <- struct{}{}
		go func(path string) {
			defer wg.Done()
			defer func() { <-sem }()

			file, err := os.Open(path)
			if err != nil {
				return
			}
			defer file.Close()

			var bytesCount int64
			scanner := bufio.NewScanner(file)
			buf := make([]byte, 0, 64*1024)
			scanner.Buffer(buf, 1024*1024)

			lineNum := 0
			for scanner.Scan() {
				line := scanner.Text()
				bytesCount += int64(len(line)) + 1

				if strings.Contains(line, "atheon:ignore") {
					lineNum++
					continue
				}

				for _, cs := range activeScanners {
					if !cs.combined.MatchString(line) {
						continue
					}
					for _, p := range cs.patterns {
						if p.Matches(line) {
							mu.Lock()
							callback(Finding{
								Pattern: p.Name(),
								File:    path,
								Line:    lineNum + 1,
								Content: strings.TrimSpace(line),
							})
							mu.Unlock()
						}
					}
				}
				lineNum++
			}

			if scanner.Err() == nil {
				mu.Lock()
				totalBytes += bytesCount
				filesScanned++
				mu.Unlock()
			}
		}(p)
	}
	wg.Wait()

	return &Stats{
		Files:     filesScanned,
		Bytes:     totalBytes,
		ElapsedMs: time.Since(start).Milliseconds(),
	}, nil
}

// ScanReader scans from an io.Reader line-by-line, useful for streaming
// data or processing piped input.
func ScanReader(r io.Reader, source string, callback FindingCallback) error {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "atheon:ignore") {
			lineNum++
			continue
		}

		for _, cs := range activeScanners {
			if !cs.combined.MatchString(line) {
				continue
			}
			for _, p := range cs.patterns {
				if p.Matches(line) {
					callback(Finding{
						Pattern: p.Name(),
						File:    source,
						Line:    lineNum + 1,
						Content: strings.TrimSpace(line),
					})
				}
			}
		}
		lineNum++
	}

	return scanner.Err()
}

// ProfiledScan wraps the original ScanFile for comparison in profiling
// This helps identify bottlenecks between the old and new implementations.
func ProfiledScan(path string, callback FindingCallback) (*Stats, error) {
	start := time.Now()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	findings := scanLines(string(data), path)
	for _, f := range findings {
		callback(f)
	}

	return &Stats{
		Files:     1,
		Bytes:     int64(len(data)),
		ElapsedMs: time.Since(start).Milliseconds(),
	}, nil
}
