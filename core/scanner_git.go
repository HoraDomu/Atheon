package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ScanGitRemote shallow-clones a git repository and scans it for pattern
// matches. The clone is placed in a temporary directory that is removed before
// ScanGitRemote returns (even on error). The context controls cancellation.
//
// Example: ScanGitRemote(ctx, "https://github.com/user/repo")
func ScanGitRemote(ctx context.Context, remoteURL string) ([]Finding, *Stats, error) {
	tmpDir, err := os.MkdirTemp("", "atheon-git-*")
	if err != nil {
		return nil, nil, fmt.Errorf("create temp dir: %w", err)
	}
	// Clean up the clone no matter what.
	defer os.RemoveAll(tmpDir)

	// Shallow clone: --depth 1 minimises download size and clone time.
	// --single-branch avoids fetching full history.
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", "--single-branch", remoteURL, tmpDir)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if err := cmd.Run(); err != nil {
		return nil, nil, fmt.Errorf("git clone failed: %w (is git installed?)", err)
	}

	findings, stats, err := ScanDir(ctx, tmpDir)
	if err != nil {
		return findings, stats, err
	}
	if stats != nil {
		// Stamp the elapsed time for the clone+scan.
		stats.ElapsedMs = 0 // not currently tracked per-phase
	}
	return findings, stats, nil
}

// ScanGitRemoteFiles is like ScanGitRemote but returns the absolute paths of
// all scanned files alongside the findings. This is useful for tooling that
// needs to know which files were examined.
func ScanGitRemoteFiles(ctx context.Context, remoteURL string) ([]Finding, []string, *Stats, error) {
	tmpDir, err := os.MkdirTemp("", "atheon-git-*")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", "--single-branch", remoteURL, tmpDir)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if err := cmd.Run(); err != nil {
		return nil, nil, nil, fmt.Errorf("git clone failed: %w (is git installed?)", err)
	}

	// Collect file paths before scanning so callers can see what was examined.
	var scannedPaths []string
	_ = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			scannedPaths = append(scannedPaths, path)
		}
		return nil
	})

	findings, stats, err := ScanDir(ctx, tmpDir)
	if err != nil {
		return findings, scannedPaths, stats, err
	}
	return findings, scannedPaths, stats, nil
}
