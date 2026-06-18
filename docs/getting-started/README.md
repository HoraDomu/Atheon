# Getting Started with Atheon

Welcome to Atheon - a community-driven pattern matching engine for detecting secrets, PII, and other sensitive data.

## Installation

### Package Managers

**macOS / Linux (Homebrew):**
```bash
brew tap HoraDomu/atheon
brew install atheon
```

**Windows (Scoop):**
```bash
scoop bucket add atheon https://github.com/HoraDomu/scoop-atheon
scoop install atheon
```

### Manual Installation

Download the latest binary for your platform from [GitHub Releases](https://github.com/HoraDomu/Atheon/releases/latest). No installation, runtime, or dependencies required.

### Build from Source

```bash
git clone https://github.com/HoraDomu/Atheon.git
cd Atheon
go build -o atheon .
```

## Quick Start

### Basic Usage

```bash
# Scan current directory
atheon ./

# Scan a specific file
atheon --file config.yaml

# Scan environment variables
atheon --env

# Pipe content through stdin
cat file.txt | atheon -
git diff | atheon -
```

### Category Filtering

```bash
# Scan only specific categories
atheon --categories=secrets ./
atheon --categories=secrets,pii ./

# List available categories
atheon list categories

# List all loaded patterns
atheon list
```

### Output Formats

```bash
# Human-readable output (default)
atheon ./

# JSON output for automation
atheon --json ./
```

## Exit Codes

- `0` - No findings detected
- `1` - One or more findings detected

This makes Atheon CI/CD friendly for automated pipelines and git hooks.

## Configuration

### Ignore Rules

**Git Integration:** Automatically respects `.gitignore` files

**Custom Ignores:** Create `.atheonignore` in your project root:
```
test/
*.generated.go
.env
node_modules/
```

**Line-level Ignores:** Add `atheon:ignore` to suppress specific lines:
```yaml
debug_key: sk-fake-key-for-testing  # atheon:ignore
```

## Updating Patterns

```bash
# Download the latest community patterns bundle
atheon update
```

New pattern bundles are released on the 10th and 21st of every month.

## Pre-commit Integration

### Git Hooks

**Simple Pre-commit Hook:**
```bash
# .git/hooks/pre-commit
#!/bin/sh
atheon ./
```

**Category-filtered Hook:**
```bash
# .git/hooks/pre-commit
#!/bin/sh
atheon --categories=secrets ./
```

### Hook Managers

Atheon works with popular hook managers:

**Husky (Node.js):**
```json
{
  "husky": {
    "hooks": {
      "pre-commit": "atheon ./"
    }
  }
}
```

**Lefthook:**
```yaml
pre-commit:
  commands:
    atheon:
      run: atheon ./
```

**pre-commit (Python):**
```yaml
- repo: local
  hooks:
    - id: atheon
      name: Atheon Pattern Scanner
      entry: atheon
      language: system
      pass_filenames: true
```

## MCP Server Integration

Atheon includes an MCP server for AI tool integration:

```json
{
  "mcpServers": {
    "atheon": {
      "command": "atheon-mcp"
    }
  }
}
```

**Available Tools:**
- `scan_string` - Scan text content
- `scan_file` - Scan individual files
- `scan_dir` - Scan directories

All tools accept optional `categories` parameter for filtering.

## Next Steps

- Explore [available patterns](../patterns/README.md)
- Learn about [pattern development](../patterns/creating-patterns.md)
- Set up [automation workflows](../integration/ci-cd.md)
- Join the [community](../community/README.md)