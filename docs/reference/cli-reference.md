# CLI Reference

Complete reference for Atheon command-line interface commands and options.

## Commands

### `atheon <path>`

Scan a directory or file for pattern matches.

**Usage:**
```bash
atheon <path>
```

**Examples:**
```bash
# Scan current directory
atheon ./

# Scan specific directory
atheon /path/to/project

# Scan specific file
atheon config.yaml
```

**Exit Codes:**
- `0` - No findings detected
- `1` - One or more findings detected

---

### `atheon --file <path>`

Explicitly scan a single file.

**Usage:**
```bash
atheon --file <path>
```

**Examples:**
```bash
atheon --file config/app.yaml
atheon --file /path/to/secret.txt
```

---

### `atheon --env`

Scan all environment variables for patterns.

**Usage:**
```bash
atheon --env
```

**Examples:**
```bash
# Scan environment variables
atheon --env

# Use in CI/CD for runtime secret detection
atheon --env
```

**Output:**
- Findings reported with format `env:VARIABLE_NAME`

---

### `atheon --json <path>`

Output findings in JSON format for automation.

**Usage:**
```bash
atheon --json <path>
```

**Examples:**
```bash
atheon --json ./
atheon --json --file config.yaml
```

**Output Format:**
```json
[
  {
    "pattern": "openai-api-key",
    "file": "config/app.yaml",
    "line": 47,
    "match": "# debug key: sk-..."
  }
]
```

---

### `atheon --categories=<cats> <path>`

Scan specific pattern categories only.

**Usage:**
```bash
atheon --categories=<category1,category2> <path>
```

**Examples:**
```bash
# Scan only secrets
atheon --categories=secrets ./

# Scan multiple categories
atheon --categories=secrets,pii ./

# Scan all categories (default)
atheon --all ./
```

**Available Categories:**
- `secrets` - API keys, tokens, credentials
- `pii` - Personal identifiers, sensitive data

---

### `atheon --all <path>`

Scan all categories (default behavior).

**Usage:**
```bash
atheon --all <path>
```

**Note:** This is the default behavior if no category filter is specified.

---

### `atheon list`

List every loaded pattern.

**Usage:**
```bash
atheon list
```

**Output:**
```
aws-access-key-id
aws-secret-access-key
credit-card
github-pat
openai-api-key
phone-number
slack-bot-token
ssn
stripe-secret-key
twilio-account-sid
...

14 pattern(s) loaded
```

---

### `atheon list categories`

List available pattern categories.

**Usage:**
```bash
atheon list categories
```

**Output:**
```
secrets
pii
```

---

### `atheon update`

Download the latest patterns bundle from GitHub releases.

**Usage:**
```bash
atheon update
```

**Output:**
```
downloading patterns bundle...
patterns updated.
```

**Note:** Bundle is saved to `~/.atheon/patterns.bundle`

---

### `atheon -` / `atheon --stdin`

Read from stdin for pipe integration.

**Usage:**
```bash
cat file.txt | atheon -
git diff | atheon -
echo "test content" | atheon -
```

**Examples:**
```bash
# Scan git diff
git diff | atheon -

# Scan file content
cat config.yaml | atheon -

# Scan command output
curl -s https://api.example.com | atheon -
```

---

### `atheon --help`

Display help message.

**Usage:**
```bash
atheon --help
```

---

## Options

### Global Options

**`--json`**
- Output findings in JSON format
- Useful for automation and CI/CD integration

**`--categories=<cats>`**
- Filter by pattern categories
- Comma-separated list of category names

**`--all`**
- Scan all categories (default behavior)

**`--file <path>`**
- Explicitly scan a single file

**`--env`**
- Scan environment variables

**`-` / `--stdin`**
- Read from stdin

---

## Output Formats

### Human-Readable Output (Default)

```
openai-api-key  config/app.yaml:47
  # de****a8f3c

github-pat  .env:3
  GITHUB_TOKEN=ghp_****

2 finding(s)
scanned 14 file(s)  22.1 KB  3ms
```

### JSON Output

```json
[
  {
    "pattern": "openai-api-key",
    "file": "config/app.yaml",
    "line": 47,
    "match": "# debug key: sk-..."
  },
  {
    "pattern": "github-pat",
    "file": ".env",
    "line": 3,
    "match": "GITHUB_TOKEN=ghp_..."
  }
]
```

---

## Exit Codes

### 0 - No Findings

```bash
$ atheon ./clean-project/
no findings.
$ echo $?
0
```

### 1 - Findings Detected

```bash
$ atheon ./project-with-secrets/
openai-api-key  config.yaml:47
  # de****a8f3c

1 finding(s)
$ echo $?
1
```

### Usage in Scripts

```bash
#!/bin/bash
atheon ./
if [ $? -ne 0 ]; then
  echo "Security findings detected!"
  exit 1
fi
```

---

## Environment Variables

### `ATHEON_CATEGORIES`

Set default categories for scanning.

```bash
export ATHEON_CATEGORIES="secrets,pii"
atheon ./  # Uses exported categories
```

### `ATHEON_ALL`

Scan all categories by default.

```bash
export ATHEON_ALL=1
atheon ./  # Scans all categories
```

---

## Configuration Files

### `.atheonignore`

Ignore specific files and directories.

```
# Comments start with #
test/
fixtures/
*_test.go
node_modules/
.env
*.generated.go
```

### `.gitignore`

Automatically respected for directory scans.

---

## Performance Options

### File Filtering

Atheon automatically skips:
- Binary files (`.png`, `.jpg`, `.pdf`, `.zip`, etc.)
- Ignored directories (`.git`, `node_modules`, `vendor`, etc.)
- Files matching ignore patterns

### Category Filtering

Use categories to improve performance:

```bash
# Faster: Only scan secrets
atheon --categories=secrets ./

# Slower: Scan all categories
atheon --all ./
```

---

## Error Handling

### Invalid Path

```bash
$ atheon /nonexistent/path
error: path not found: /nonexistent/path
```

### Permission Denied

```bash
$ atheon /root/
error: permission denied
```

### Invalid Regex (Pattern Error)

```bash
# Pattern loading errors are shown at startup
atheon: skipping "invalid-pattern": invalid regex pattern
```

---

## Advanced Usage

### Combining Options

```bash
# JSON output with category filter
atheon --json --categories=secrets ./

# Scan specific file with JSON output
atheon --json --file config.yaml

# Environment scan with categories
atheon --env --categories=secrets,pii
```

### Pipeline Integration

```bash
# Scan git diff
git diff | atheon -

# Scan multiple directories
find . -maxdepth 2 -type d | while read dir; do
  atheon "$dir"
done

# Scan all commits
git log --all --oneline | while read commit rest; do
  git show "$commit" | atheon - && echo "✓ $commit" || echo "✗ $commit"
done
```

---

## Troubleshooting

### Pattern Not Matching

```bash
# Verify pattern is loaded
atheon list | grep pattern-name

# Test with simple input
echo "test input" | atheon -

# Check category filtering
atheon --all ./  # Remove category filter
```

### Performance Issues

```bash
# Use category filtering
atheon --categories=secrets ./

# Create ignore file
cat > .atheonignore << EOF
test/
node_modules/
EOF
```

### Update Issues

```bash
# Manually download bundle
curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/patterns.bundle \
  -o ~/.atheon/patterns.bundle
```

---

## Examples

### Development Workflow

```bash
# Pre-commit hook
#!/bin/sh
atheon --categories=secrets ./
if [ $? -ne 0 ]; then
  echo "Remove secrets before committing"
  exit 1
fi
```

### CI/CD Integration

```yaml
# GitHub Actions
- name: Security Scan
  run: |
    atheon --json ./ > findings.json
    if [ -s findings.json ]; then
      echo "Security findings detected!"
      exit 1
    fi
```

### Environment Scanning

```bash
# CI/CD pipeline
atheon --env
if [ $? -ne 0 ]; then
  echo "Secrets in environment variables!"
  exit 1
fi
```

---

## Additional Resources

- [Getting Started Guide](../getting-started/README.md)
- [Common Use Cases](../getting-started/common-use-cases.md)
- [Integration Guide](../integration/README.md)
- [Pattern Development](../patterns/README.md)