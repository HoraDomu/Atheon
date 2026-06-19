package core

import (
	"os/exec"
	"strings"
)

// GitChangedFiles returns files changed since the given commit
func GitChangedFiles(baseCommit string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", baseCommit, "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var changedFiles []string
	for _, file := range files {
		if file != "" {
			changedFiles = append(changedFiles, strings.TrimSpace(file))
		}
	}

	return changedFiles, nil
}

// GitStagedFiles returns files staged for commit
func GitStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var stagedFiles []string
	for _, file := range files {
		if file != "" {
			stagedFiles = append(stagedFiles, strings.TrimSpace(file))
		}
	}

	return stagedFiles, nil
}

// FrameworkDetection detects the framework used in the project
func FrameworkDetection() string {
	// Check for package.json (Node.js/JavaScript)
	if _, err := exec.Command("test", "-f", "package.json").Output(); err == nil {
		return "nodejs"
	}

	// Check for requirements.txt (Python)
	if _, err := exec.Command("test", "-f", "requirements.txt").Output(); err == nil {
		return "python"
	}

	// Check for go.mod (Go)
	if _, err := exec.Command("test", "-f", "go.mod").Output(); err == nil {
		return "go"
	}

	// Check for pom.xml (Java/Maven)
	if _, err := exec.Command("test", "-f", "pom.xml").Output(); err == nil {
		return "maven"
	}

	return "unknown"
}

// SeverityLevel represents the severity of a finding
type SeverityLevel int

const (
	SeverityLow SeverityLevel = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// String returns the string representation of the severity
func (s SeverityLevel) String() string {
	switch s {
	case SeverityLow:
		return "LOW"
	case SeverityMedium:
		return "MEDIUM"
	case SeverityHigh:
		return "HIGH"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// PatternSeverity defines severity levels for different pattern types
var PatternSeverity = map[string]SeverityLevel{
	// Critical severity
	"openai-api-key":                     SeverityCritical,
	"aws-access-key-id":                 SeverityCritical,
	"github-personal-token":             SeverityCritical,
	"azure-client-secret":               SeverityCritical,
	"gcp-api-key":                       SeverityCritical,

	// High severity
	"azure-service-principal-secret":    SeverityHigh,
	"azure-storage-account-key":          SeverityHigh,
	"stripe-live-key":                    SeverityHigh,

	// Medium severity
	"console-log":                       SeverityMedium,
	"debug-print":                       SeverityMedium,
	"todo-comment":                      SeverityMedium,

	// Low severity
	"placeholder-code":                  SeverityLow,
	"deprecated-function":               SeverityLow,
	"fixme-comment":                     SeverityLow,
}