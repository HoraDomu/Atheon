package core

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAudit(t *testing.T) {
	// Run Audit against the current directory — it should complete without error.
	report, err := Audit(context.Background(), ".")
	if err != nil {
		t.Fatalf("Audit failed: %v", err)
	}
	if report == nil {
		t.Fatal("expected non-nil report")
	}
	if report.Version == "" {
		t.Error("report should have a version")
	}
	if report.Root == "" {
		t.Error("report should have a root")
	}
	if report.Summary.Total == 0 {
		t.Error("summary should have total > 0")
	}
}

func TestAudit_Timestamp(t *testing.T) {
	report, err := Audit(context.Background(), ".")
	if err != nil {
		t.Fatalf("Audit failed: %v", err)
	}
	if report.GeneratedAt.IsZero() {
		t.Error("GeneratedAt should be set")
	}
}

func TestWriteReport(t *testing.T) {
	report := &AuditReport{
		Version:     "1.0",
		GeneratedAt: now(),
		Root:        t.TempDir(),
		ElapsedMs:   10,
		Results: []AuditResult{
			{
				Check:    "test-check",
				Passed:   true,
				Findings: []AuditFinding{},
			},
			{
				Check:  "failing-check",
				Passed: false,
				Findings: []AuditFinding{
					{File: "main.go", Line: 10, Message: "test finding", Severity: "warning"},
				},
			},
		},
		Summary: AuditSummary{Total: 2, Passed: 1, Failed: 1},
	}

	dir := filepath.Join(t.TempDir(), "audit-report-test")
	if err := WriteReport(report, dir); err != nil {
		t.Fatalf("WriteReport failed: %v", err)
	}

	jsonPath := filepath.Join(dir, "REPORT.json")
	if _, err := os.Stat(jsonPath); err != nil {
		t.Errorf("expected REPORT.json, got error: %v", err)
	}

	mdPath := filepath.Join(dir, "REPORT.md")
	if _, err := os.Stat(mdPath); err != nil {
		t.Errorf("expected REPORT.md, got error: %v", err)
	}
}

func TestAtoiSafe(t *testing.T) {
	cases := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"456abc", 456},
		{"abc", 0},
		{"", 0},
		{"  42  ", 42}, // spaces are skipped, digits accumulate
	}
	for _, tc := range cases {
		got := atoiSafe(tc.input)
		if got != tc.expected {
			t.Errorf("atoiSafe(%q): got %d, want %d", tc.input, got, tc.expected)
		}
	}
}

func TestRunNolintCheck(t *testing.T) {
	res := runNolintCheck(".")
	if res.Check != "nolint" {
		t.Errorf("expected check name 'nolint', got %q", res.Check)
	}
	// Should not panic, regardless of findings
	_ = res.Findings
}

func TestRunTodoFixmeCheck(t *testing.T) {
	res := runTodoFixmeCheck(".")
	if res.Check != "todo-fixme" {
		t.Errorf("expected check name 'todo-fixme', got %q", res.Check)
	}
	_ = res.Findings
}

func TestRunVetCheck(t *testing.T) {
	res := runVetCheck(".")
	if res.Check != "go-vet" {
		t.Errorf("expected check name 'go-vet', got %q", res.Check)
	}
	_ = res.Findings
}

func TestRunSentinelCheck(t *testing.T) {
	res := runSentinelCheck(".")
	if res.Check != "sentinel-errors" {
		t.Errorf("expected check name 'sentinel-errors', got %q", res.Check)
	}
	_ = res.Findings
}

// TestWriteReport_CreatesDirectory verifies WriteReport creates the directory.
func TestWriteReport_CreatesDirectory(t *testing.T) {
	report := &AuditReport{
		Version:     "1.0",
		GeneratedAt: now(),
		Root:        t.TempDir(),
		ElapsedMs:   1,
		Results:     []AuditResult{},
		Summary:     AuditSummary{Total: 0, Passed: 0, Failed: 0},
	}
	// Directory should not exist yet
	dir := filepath.Join(t.TempDir(), "does-not-exist", "subdir")
	if err := WriteReport(report, dir); err != nil {
		t.Fatalf("WriteReport failed: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("WriteReport should create directory: %v", err)
	}
}

// TestAudit_NolintFindings exercises the nolint check against this very file
// which has a nolint:gocritic comment.
func TestAudit_NolintFindings(t *testing.T) {
	// Run nolint check against this directory.
	res := runNolintCheck(".")
	// This file (audit_test.go) is the only source, but it may not have
	// any nolint comments — the check should not panic either way.
	_ = res.Findings
}

// TestAudit_TodoFixmeFindings exercises the todo-fixme check.
func TestAudit_TodoFixmeFindings(t *testing.T) {
	res := runTodoFixmeCheck(".")
	_ = res.Findings
}

// now is a test helper that returns the current time.
func now() time.Time { return time.Now() }
