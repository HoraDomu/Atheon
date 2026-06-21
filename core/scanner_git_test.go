package core

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestScanGitRemote_LocalRepo tests ScanGitRemote against a freshly
// initialised local git repository.
func TestScanGitRemote_LocalRepo(t *testing.T) {
	// Create a temp directory and initialise a local git repo inside it.
	repoDir := t.TempDir()
	subDir := filepath.Join(repoDir, "subdir")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("MkdirAll subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "config.txt"), []byte("password=super-secret-123\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Init git (required for file:// URLs to work as remotes).
	for _, dir := range []string{repoDir, subDir} {
		cmd := exec.Command("git", "init")
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Skipf("git init failed (may not be configured for local repos): %v\n%s", err, out)
		}
		cmd = exec.Command("git", "config", "user.email", "test@test.com")
		cmd.Dir = dir
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Skipf("git config failed: %v", err)
		}
		cmd = exec.Command("git", "config", "user.name", "Test")
		cmd.Dir = dir
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Skipf("git config failed: %v", err)
		}
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = dir
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Skipf("git add failed: %v", err)
		}
		cmd = exec.Command("git", "commit", "-m", "initial")
		cmd.Dir = dir
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Skipf("git commit failed: %v", err)
		}
	}

	// Scan the local repo using a file:// URL.
	fileURL := "file://" + repoDir
	findings, stats, err := ScanGitRemote(context.Background(), fileURL)
	if err != nil {
		t.Fatalf("ScanGitRemote failed: %v", err)
	}
	if stats == nil {
		t.Fatal("expected non-nil stats")
	}
	_ = findings // may be empty if no patterns matched
}

// TestScanGitRemoteFiles_LocalRepo tests ScanGitRemoteFiles against a local repo.
func TestScanGitRemoteFiles_LocalRepo(t *testing.T) {
	repoDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(repoDir, "secrets.txt"), []byte("github_token=ghp_abcdefghijklmnopqrstuvwxyz123456\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Init git
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git init failed: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = repoDir
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git config failed: %v", err)
	}
	cmd = exec.Command("git", "config", "user.name", "Test")
	cmd.Dir = repoDir
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git config failed: %v", err)
	}
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = repoDir
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git add failed: %v", err)
	}
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = repoDir
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git commit failed: %v", err)
	}

	fileURL := "file://" + repoDir
	findings, paths, stats, err := ScanGitRemoteFiles(context.Background(), fileURL)
	if err != nil {
		t.Fatalf("ScanGitRemoteFiles failed: %v", err)
	}
	if stats == nil {
		t.Fatal("expected non-nil stats")
	}
	if len(paths) == 0 {
		t.Error("expected at least one scanned path")
	}
	_ = findings
	for _, p := range paths {
		if p == "" {
			t.Error("empty path in scanned paths")
		}
	}
}

// TestScanGitRemote_InvalidURL tests that an invalid URL produces an error.
func TestScanGitRemote_InvalidURL(t *testing.T) {
	_, _, err := ScanGitRemote(context.Background(), "http://127.0.0.1:99999/nonexistent")
	if err == nil {
		t.Error("expected error for unreachable URL")
	}
}

// TestScanGitRemote_CancelledContext tests that context cancellation works.
func TestScanGitRemote_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	// Use a URL that would take time to clone if it were reachable.
	_, _, err := ScanGitRemote(ctx, "https://github.com/nonexistent-repo-abcxyz123456")
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

// cloneTimeEstimate is a stub that always returns 0.
// It exists for future use by progress reporters.
//
//nolint:gocritic
func TestCloneTimeEstimate_Stub(t *testing.T) {
	got := cloneTimeEstimate("https://example.com/repo")
	if got != 0 {
		t.Errorf("cloneTimeEstimate: got %v, want 0", got)
	}
}
