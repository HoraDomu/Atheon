# Network Scanner

Atheon can scan remote HTTP/S URLs and git repositories in addition to local files and environment variables.

---

## `atheon scan-url <url>`

Fetches a remote URL and scans the response body for pattern matches.

**Usage:**
```bash
atheon scan-url https://example.com/secrets.txt
atheon scan-url --format=sarif https://api.example.com/config
```

**How it works:**
1. Performs an HTTP GET request with a 10-second timeout
2. Does not follow redirects (returns the redirect response as-is)
3. Reads the response body with a 5 MB cap to prevent memory exhaustion
4. Scans only `text/*`, `application/json`, `application/xml`, and `application/javascript` content types; all other types are rejected

**Output:** Uses the same structured reporting pipeline as local scans — `--format=text|json|sarif|html` and `--json` work identically.

**Exit codes:**

| Situation | Exit Code |
|-----------|-----------|
| No findings | 0 |
| One or more findings | 1 |
| Invalid URL, unsupported scheme, non-200 response, or read error | 1 |

**Examples:**
```bash
# Scan a raw file on GitHub
atheon scan-url https://raw.githubusercontent.com/user/repo/main/.env

# Scan an API response as JSON
atheon scan-url --format=json https://api.example.com/config

# Scan with category filter
atheon scan-url --categories=secrets https://example.com/config.json
```

---

## `atheon scan-git <remote_url>`

Shallow-clones a git repository and scans all its files for pattern matches.

**Usage:**
```bash
atheon scan-git https://github.com/user/repo
atheon scan-git git@github.com:user/repo.git
```

**How it works:**
1. Creates a temporary directory
2. Runs `git clone --depth=1 --single-branch` to fetch only the latest commit
3. Scans the cloned tree using the same engine as `atheon <dir>`
4. Removes the temporary clone directory before returning (even on error)

The clone is shallow and single-branch, making it fast even for large repositories.

**Authentication:** `GIT_TERMINAL_PROMPT=0` is set to prevent interactive credential prompts. For private repositories, use SSH (`git@github.com:...`) or configure `netrc` / git credential helpers beforehand.

**Exit codes:** Same as `scan-url`.

**Examples:**
```bash
# Scan a public repository
atheon scan-git https://github.com/user/repo

# SARIF output for GitHub Code Scanning
atheon scan-git --format=sarif https://github.com/user/repo > atheon-results.sarif

# Scan only secrets category
atheon scan-git --categories=secrets https://github.com/user/repo
```

---

## MCP Tools

The `scan_url` and `scan_git` tools are also available via the MCP server (`cmd/mcp/main.go`), enabling AI agents to scan remote resources without a local checkout.

### `scan_url`

```json
{
  "name": "scan_url",
  "arguments": {
    "url": "https://example.com/config",
    "categories": ["secrets"]
  }
}
```

Returns findings as plain text, suitable for ingestion by LLMs.

### `scan_git`

```json
{
  "name": "scan_git",
  "arguments": {
    "remote_url": "https://github.com/user/repo",
    "categories": ["secrets", "pii"]
  }
}
```

Shallow-clones the repository, scans it, and returns findings as plain text. The clone is removed immediately after scanning.

---

## Programmatic API

The network scanner functions are exported from `core/`:

```go
// ScanURL fetches a URL and returns findings, stats, and any error.
func ScanURL(ctx context.Context, rawURL string) ([]Finding, *Stats, error)

// ScanGitRemote shallow-clones a repo and scans it; temp dir is cleaned up.
func ScanGitRemote(ctx context.Context, remoteURL string) ([]Finding, *Stats, error)

// ScanGitRemoteFiles is like ScanGitRemote but also returns scanned file paths.
func ScanGitRemoteFiles(ctx context.Context, remoteURL string) ([]Finding, []string, *Stats, error)
```

All three accept a `context.Context` for cancellation and return `[]core.Finding` that flow through `core.Render` for output.

---

## Limitations

- **Binary responses** are rejected; only text-based content types are scanned
- **Large responses** are capped at 5 MB
- **Redirects** are not followed (prevents credential leakage to untrusted servers)
- **`scan-git`** requires `git` to be installed and in `PATH`; interactive prompts are disabled
- **`scan-git`** uses shallow clone (`--depth=1 --single-branch`) — it does not scan commit history
