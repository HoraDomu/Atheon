# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Atheon is a community-driven pattern matching engine for detecting secrets, PII, and other sensitive data. The engine itself is minimal and domain-agnostic — all intelligence lives in patterns defined as YAML files in `community/`.

## Core Architecture

The project has three main components:

### 1. Engine (`core/`)
- **bundle.go**: Pattern loading and management. Patterns are shipped as a compressed gzip bundle embedded in the binary. Supports dynamic reloading from `~/.atheon/patterns.bundle` for updates.
- **pattern.go**: Pattern interface and registry. Simple interface with `Name()` and `Matches(line string) bool`.
- **runner.go**: Scanning logic with parallel file processing, ignore rules, and category filtering.
- **finding.go**: Data structures for findings and statistics.

### 2. CLI (`main.go`)
Entry point that orchestrates the engine. Handles argument parsing, output formatting (text/JSON), and exit codes for CI integration.

### 3. MCP Server (`cmd/mcp/`)
Model Context Protocol server over stdio. Exposes `scan_string`, `scan_file`, and `scan_dir` tools for AI assistant integration.

### 4. Bundler (`bundler/`)
Compiles YAML patterns from `community/` into the compressed bundle format. This must be run after adding/modifying patterns.

## Development Commands

**Build the CLI:**
```bash
go build -o atheon .
```

**Build the MCP server:**
```bash
go build -o atheon-mcp ./cmd/mcp
```

**Run tests:**
```bash
go test ./...
```

**Rebuild the pattern bundle (required after pattern changes):**
```bash
go run ./bundler
```

**Test a pattern manually:**
```bash
atheon --file <path-to-sample>
atheon --categories=secrets ./
```

**Release build (requires goreleaser):**
```bash
goreleaser release --clean
```

## Pattern Development Workflow

Patterns are defined as YAML files in `community/<category>/`:

```yaml
name: my-service-api-key
match: '\bmsvc_[A-Za-z0-9]{32}\b'
```

1. Create YAML file in appropriate category folder
2. Run `go run ./bundler` to compile patterns into `core/patterns.bundle`
3. Add test cases to `core/bundle_test.go` under the `cases` map
4. Run `go test ./...` to verify
5. Test manually against real samples
6. Commit both the YAML file and updated `core/patterns.bundle`

## Key Implementation Details

**Category Filtering**: Patterns are organized by category and filtered at load time. Within each category, patterns are combined into a single regex for efficient pre-filtering before individual pattern matching.

**Parallel Scanning**: Directory scans use worker pools with 256 concurrent workers. Each file is read independently, and findings are collected into a slice.

**Ignore Rules**: Automatically respects `.gitignore` and `.atheonignore`. Line-level ignores use `atheon:ignore` comments.

**Bundle Loading**: The embedded bundle is loaded at init, but can be overridden by `~/.atheon/patterns.bundle` for pattern updates.

**Exit Codes**: Returns 0 for clean scans, 1 for findings. This makes it CI/CD friendly for hooks and pipelines.

## Pattern Testing Approach

Tests in `core/bundle_test.go` follow a table-driven pattern:

```go
"pattern-name": {
    matches:    []string{"expected positive match 1", "expected positive match 2"},
    nonMatches: []string{"should not match", "also should not match"},
},
```

Each pattern should have comprehensive test coverage for both positive and negative cases to prevent false positives and false negatives.

## Release Process

Releases are automated via GitHub Actions on version tags. GoReleaser builds binaries for multiple platforms, publishes to GitHub Releases, Homebrew, and Scoop. The bundle is automatically included as a release asset.

The release schedule is the 10th and 21st of every month.
