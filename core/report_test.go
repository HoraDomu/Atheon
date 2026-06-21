package core

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestRenderText_NoFindings(t *testing.T) {
	r := Report{Version: "test", Findings: nil, Stats: Stats{}}
	out := Render(r, FormatText)
	if out == "" {
		t.Error("expected non-empty text output")
	}
	if !strings.Contains(out, "no findings") {
		t.Errorf("expected 'no findings' in output: %s", out)
	}
}

func TestRenderText_WithFindings(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "dir",
		Target:   "/src",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "config.txt", Line: 3, Content: "AKIAIOSFODNN7EXAMPLE"},
		},
		Stats: Stats{Files: 10, Bytes: 1024, ElapsedMs: 50},
	}
	out := Render(r, FormatText)
	if !strings.Contains(out, "aws-access-key") {
		t.Errorf("expected pattern in output: %s", out)
	}
	if !strings.Contains(out, "config.txt:3") {
		t.Errorf("expected location in output: %s", out)
	}
	if !strings.Contains(out, "1 finding(s)") {
		t.Errorf("expected summary in output: %s", out)
	}
	if !strings.Contains(out, "10 file(s)") {
		t.Errorf("expected file count in output: %s", out)
	}
}

func TestRenderJSON_Shape(t *testing.T) {
	r := Report{
		Version:  "test",
		Findings: []Finding{{Pattern: "aws-access-key", File: "f", Line: 1, Content: "AKIAIOSFODNN7EXAMPLE"}},
	}
	out := Render(r, FormatJSON)
	var items []map[string]any
	if err := json.Unmarshal([]byte(out), &items); err != nil {
		t.Fatalf("output is not valid JSON: %s", out)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if items[0]["pattern"] != "aws-access-key" {
		t.Errorf("unexpected pattern field: %v", items[0]["pattern"])
	}
	// Content must be redacted
	if items[0]["match"] == "AKIAIOSFODNN7EXAMPLE" {
		t.Error("match should be redacted, got raw value")
	}
}

func TestRenderJSON_Empty(t *testing.T) {
	r := Report{Version: "test", Findings: []Finding{}}
	out := Render(r, FormatJSON)
	var items []map[string]any
	if err := json.Unmarshal([]byte(out), &items); err != nil {
		t.Fatalf("output is not valid JSON: %s", out)
	}
	if len(items) != 0 {
		t.Errorf("expected empty array, got %d items", len(items))
	}
}

func TestRenderSARIF_Structure(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "dir",
		Target:   "/src",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "config.txt", Line: 5, Content: "AKIAIOSFODNN7EXAMPLE"},
		},
		Stats: Stats{Files: 1, Bytes: 100, ElapsedMs: 10},
	}
	out := Render(r, FormatSARIF)

	var log sarifLog
	if err := json.Unmarshal([]byte(out), &log); err != nil {
		t.Fatalf("output is not valid SARIF JSON: %s", out)
	}

	if log.Version != "2.1.0" {
		t.Errorf("expected SARIF version 2.1.0, got %s", log.Version)
	}
	if len(log.Runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(log.Runs))
	}
	run := log.Runs[0]
	if run.Tool.Driver.Name != "Atheon" {
		t.Errorf("expected tool name 'Atheon', got %s", run.Tool.Driver.Name)
	}
	if run.Tool.Driver.Version != "test" {
		t.Errorf("expected tool version 'test', got %s", run.Tool.Driver.Version)
	}
	if len(run.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(run.Results))
	}
	res := run.Results[0]
	if res.RuleID != "aws-access-key" {
		t.Errorf("expected ruleId 'aws-access-key', got %s", res.RuleID)
	}
	if res.Level != "warning" {
		t.Errorf("expected level 'warning', got %s", res.Level)
	}
	if len(res.Locations) != 1 {
		t.Errorf("expected 1 location, got %d", len(res.Locations))
	}
	loc := res.Locations[0]
	if loc.PhysicalLocation.ArtifactLocation.URI != "config.txt" {
		t.Errorf("expected URI 'config.txt', got %s", loc.PhysicalLocation.ArtifactLocation.URI)
	}
	if loc.PhysicalLocation.Region.StartLine != 5 {
		t.Errorf("expected startLine 5, got %d", loc.PhysicalLocation.Region.StartLine)
	}
}

func TestRenderSARIF_RuleDeduplication(t *testing.T) {
	// Two findings with the same pattern should produce one rule entry
	// and both results should reference the same ruleIndex.
	r := Report{
		Version: "test",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "a.txt", Line: 1, Content: "key1"},
			{Pattern: "aws-access-key", File: "b.txt", Line: 2, Content: "key2"},
			{Pattern: "openai-api-key", File: "c.txt", Line: 3, Content: "key3"},
		},
	}
	out := Render(r, FormatSARIF)
	var log sarifLog
	if err := json.Unmarshal([]byte(out), &log); err != nil {
		t.Fatalf("not valid JSON: %s", out)
	}
	run := log.Runs[0]
	if len(run.Tool.Driver.Rules) != 2 {
		t.Errorf("expected 2 unique rules, got %d", len(run.Tool.Driver.Rules))
	}
	if len(run.Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(run.Results))
	}
	// Both aws-access-key results should have the same ruleIndex.
	var awsIdx int
	for i, res := range run.Results {
		if res.RuleID == "aws-access-key" {
			awsIdx = res.RuleIndex
			if i > 0 && run.Results[i-1].RuleID == "aws-access-key" {
				if run.Results[i-1].RuleIndex != awsIdx {
					t.Error("duplicate pattern should share ruleIndex")
				}
			}
		}
	}
}

func TestRenderHTML_Structure(t *testing.T) {
	r := Report{
		Version:     "v99.0.0",
		GeneratedAt: time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC),
		ScanType:    "dir",
		Target:      "/src",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "config.txt", Line: 5, Content: "AKIAIOSFODNN7EXAMPLE"},
		},
		Stats: Stats{Files: 10, Bytes: 2048, ElapsedMs: 123},
	}
	out := Render(r, FormatHTML)
	for _, needle := range []string{
		"<!DOCTYPE html>",
		"<title>Atheon Scan Report</title>",
		`<span class="summary-value">dir</span>`,
		`<span class="summary-value">/src</span>`,
		`<span class="summary-value">1</span>`,     // findings count
		`<span class="summary-value">10</span>`,    // files scanned
		`<span class="summary-value">123ms</span>`, // elapsed
		`Generated by Atheon v99.0.0`,
		"aws-access-key",
		"config.txt:5",
	} {
		if !strings.Contains(out, needle) {
			t.Errorf("expected %q in HTML output", needle)
		}
	}
	// No findings case
	r2 := Report{Version: "test", Findings: nil, GeneratedAt: time.Now()}
	out2 := Render(r2, FormatHTML)
	if !strings.Contains(out2, "No secrets detected") {
		t.Error("expected 'No secrets detected' in empty HTML")
	}
}

// TestRedact_ShortStrings covers the len <= 8 branch in redact.
func TestRedact_ShortStrings(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"", "***"},
		{"abc", "***"},
		{"12345678", "***"}, // exactly 8
	}
	for _, tc := range cases {
		got := redact(tc.input)
		if got != tc.expected {
			t.Errorf("redact(%q): got %q, want %q", tc.input, got, tc.expected)
		}
	}
}

// TestFormatBytes_AllRanges covers all branches in formatBytes.
func TestFormatBytes_AllRanges(t *testing.T) {
	cases := []struct {
		b        int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1 << 20, "1.0 MB"},
		{5 << 20, "5.0 MB"},
		{10 << 20, "10.0 MB"},
	}
	for _, tc := range cases {
		got := formatBytes(tc.b)
		if got != tc.expected {
			t.Errorf("formatBytes(%d): got %q, want %q", tc.b, got, tc.expected)
		}
	}
}

// TestRenderText_NoStats covers the branch where Stats.Files == 0.
func TestRenderText_NoStats(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "string",
		Findings: []Finding{},
		Stats:    Stats{Files: 0, Bytes: 0, ElapsedMs: 0},
	}
	out := Render(r, FormatText)
	if strings.Contains(out, "file(s)") {
		t.Errorf("expected no file count when Files=0: %s", out)
	}
	if !strings.Contains(out, "no findings") {
		t.Errorf("expected 'no findings': %s", out)
	}
}

// TestRenderText_MultipleFindings verifies multiple findings listed.
func TestRenderText_MultipleFindings(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "dir",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "a.txt", Line: 1, Content: "AKIAIOSFODNN7EXAMPLE"},
			{Pattern: "openai-api-key", File: "b.txt", Line: 2, Content: "sk-abcdefghijklmnopqrstuvwxyz"},
		},
		Stats: Stats{Files: 10, Bytes: 2048, ElapsedMs: 15},
	}
	out := Render(r, FormatText)
	if !strings.Contains(out, "2 finding(s)") {
		t.Errorf("expected '2 finding(s)': %s", out)
	}
}

// TestRenderHTML_Stats verifies the stats section in HTML.
func TestRenderHTML_Stats(t *testing.T) {
	r := Report{
		Version:     "test",
		GeneratedAt: time.Now(),
		ScanType:    "dir",
		Target:      "/src",
		Findings:    []Finding{},
		Stats:       Stats{Files: 42, Bytes: 123456, ElapsedMs: 99},
	}
	out := Render(r, FormatHTML)
	for _, needle := range []string{"42", "99ms"} {
		if !strings.Contains(out, needle) {
			t.Errorf("expected %q in HTML: %s", needle, out)
		}
	}
}

// TestRenderSARIF_NoFindings covers empty findings in SARIF.
func TestRenderSARIF_NoFindings(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "string",
		Findings: []Finding{},
	}
	out := Render(r, FormatSARIF)
	var log sarifLog
	if err := json.Unmarshal([]byte(out), &log); err != nil {
		t.Fatalf("not valid JSON: %s", out)
	}
	if len(log.Runs[0].Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(log.Runs[0].Results))
	}
}

// TestRenderSARIF_ThreeFindings verifies rule deduplication with three findings.
func TestRenderSARIF_ThreeFindings(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "dir",
		Findings: []Finding{
			{Pattern: "aws-access-key", File: "a.txt", Line: 1, Content: "key1"},
			{Pattern: "aws-access-key", File: "b.txt", Line: 2, Content: "key2"},
			{Pattern: "openai-api-key", File: "c.txt", Line: 3, Content: "key3"},
		},
	}
	out := Render(r, FormatSARIF)
	var log sarifLog
	if err := json.Unmarshal([]byte(out), &log); err != nil {
		t.Fatalf("not valid JSON: %s", out)
	}
	run := log.Runs[0]
	if len(run.Tool.Driver.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(run.Tool.Driver.Rules))
	}
	if len(run.Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(run.Results))
	}
}

// TestRenderText_StatsFilesPositive covers the scanned-bytes branch.
func TestRenderText_StatsFilesPositive(t *testing.T) {
	r := Report{
		Version:  "test",
		ScanType: "dir",
		Findings: []Finding{},
		Stats:    Stats{Files: 3, Bytes: 2048, ElapsedMs: 7},
	}
	out := Render(r, FormatText)
	if !strings.Contains(out, "3 file(s)") {
		t.Errorf("expected file count in output: %s", out)
	}
	if !strings.Contains(out, "2.0 KB") {
		t.Errorf("expected KB bytes in output: %s", out)
	}
	if !strings.Contains(out, "7ms") {
		t.Errorf("expected ms in output: %s", out)
	}
}
