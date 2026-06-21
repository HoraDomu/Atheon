# API Reference

The public Go API of the `atheon` module. Anything not
documented here is internal and may change without notice.

> **Stability:** the v1 API is stable. Breaking changes
> require a major version bump.

## Package layout

```
github.com/aliasfoxkde/Atheon
├── core/      — scanner engine, public API
├── bundler/   — pattern bundle builder (library)
└── cmd/       — CLI and MCP server (not library API)
```

The CLI and the MCP server live under `cmd/` and import
from `core/`; they are not part of the library API.

## `core` package

The scanner engine. All scan entry points take a
`context.Context` as their first argument. Functions that
mutate pattern state are also concurrency-safe.

### Types

#### `Finding`

```go
type Finding struct {
    Pattern string // pattern name that matched
    File    string // source path, or "env:KEY" for env scans
    Line    int    // 1-indexed line number; 0 for env scans
    Content string // trimmed matching line, or matching env value
}
```

A single pattern match produced by any `Scan*` function.
`File` is the path as supplied to the scan function.

#### `Stats`

```go
type Stats struct {
    Files      int     // number of files actually scanned
    Bytes      int64   // total bytes scanned
    ElapsedMs  int64   // wall-clock duration in milliseconds
    WalkErrors []error // per-file read errors collected by ScanDir
}
```

Aggregate counters returned alongside findings. `Files`
excludes binary files and skipped directories; it counts
only files whose contents were scanned. `WalkErrors`
collects read failures encountered during directory walks
(permission changes, TOCTOU races, broken symlinks) so
callers that care about every skipped file can inspect
them.

#### `Pattern`

```go
type Pattern interface {
    Name() string    // stable, human-readable identifier
    Category() string // grouping label (e.g. "secrets", "pii")
    Enabled() bool   // whether the pattern is currently active
    Matches(line string) bool
}
```

The interface implemented by all scanner patterns, both
those loaded from the embedded bundle and those registered
programmatically via `Register`.

#### `Format`

```go
type Format string

const (
    FormatText  Format = "text"
    FormatJSON  Format = "json"
    FormatSARIF Format = "sarif"
    FormatHTML  Format = "html"
)
```

Output format selector passed to `Render`. `FormatSARIF`
produces SARIF 2.1.0 compatible with GitHub Code Scanning.

#### `Report`

```go
type Report struct {
    Version     string    `json:"version"`
    GeneratedAt time.Time `json:"generatedAt"`
    ScanType   string    `json:"scanType"`   // "file", "dir", "string", "env", "url", "git"
    Target     string    `json:"target,omitempty"`
    Stats      Stats     `json:"stats"`
    Findings   []Finding `json:"findings"`
    Errors     []error  `json:"errors,omitempty"`
}
```

Canonical report structure. All four output formats are
derived from this type. The `Errors` field collects
scan-level errors (as opposed to per-file `WalkErrors`,
which live in `Stats`).

#### `AuditReport`

```go
type AuditReport struct {
    Version     string        `json:"version"`
    GeneratedAt time.Time    `json:"generatedAt"`
    Root        string        `json:"root"`
    ElapsedMs   int64         `json:"elapsedMs"`
    Results     []AuditResult `json:"results"`
    Summary     AuditSummary  `json:"summary"`
}
```

Complete output of an `Audit` run.

#### `AuditResult`

```go
type AuditResult struct {
    Check    string         `json:"check"`
    Passed   bool           `json:"passed"`
    Findings []AuditFinding  `json:"findings,omitempty"`
}
```

Outcome of a single audit check.

#### `AuditFinding`

```go
type AuditFinding struct {
    File     string `json:"file,omitempty"`
    Line     int    `json:"line,omitempty"`
    Message  string `json:"message"`
    Severity string `json:"severity,omitempty"` // "error", "warning", "info"
}
```

A single item produced by an audit check.

### Sentinel Errors

```go
var (
    ErrPatternNotFound = errors.New("pattern not found")
    ErrBundleDownload  = errors.New("bundle download failed")
    ErrBundleParse     = errors.New("bundle parse failed")
    ErrInvalidPattern  = errors.New("invalid pattern")
)
```

Callers should use `errors.Is` to compare against these
rather than string-matching error messages.

### Functions

#### ScanFile

```go
func ScanFile(ctx context.Context, path string) ([]Finding, *Stats, error)
```

Read and scan a single file. Honors `.atheonignore` and
`.gitignore` when the file is under the current working
directory. Returns `ctx.Err()` if the context is canceled
before the read completes.

```go
findings, stats, err := core.ScanFile(ctx, "/srv/app/.env")
```

#### ScanDir

```go
func ScanDir(ctx context.Context, root string) ([]Finding, *Stats, error)
```

Recursively scan every non-binary, non-ignored file under
`root`. Uses one worker per CPU (capped at a sensible
maximum). Honors ignore files at `root`. The context
controls worker cancellation; if canceled mid-walk the
goroutines exit and `ScanDir` returns `ctx.Err()` after
the `WaitGroup` drains.

```go
findings, stats, err := core.ScanDir(ctx, "/srv/app")
```

#### ScanString

```go
func ScanString(ctx context.Context, content, source string) ([]Finding, *Stats, error)
```

Scan an in-memory string. `source` is recorded as the
`File` field on each finding (e.g. the editor buffer
name). The context is checked between lines.

```go
findings, _, _ := core.ScanString(ctx, source, "main.go")
```

#### ScanEnv

```go
func ScanEnv(ctx context.Context) []Finding
```

Scan the current process environment variables for matches
against the active patterns. Each finding uses `"env:KEY"`
as its `File` and the matching value as `Content`; `Line`
is zero. The context is checked between iterations.

```go
findings := core.ScanEnv(ctx)
```

#### ScanURL

```go
func ScanURL(ctx context.Context, rawURL string) ([]Finding, *Stats, error)
```

Fetch a remote URL and scan its response body. Only
`text/*`, `application/json`, `application/xml`, and
`application/javascript` content types are scanned; binary
types return an error. Response body is capped at 5 MB.
Redirects are not followed.

#### ScanGitRemote

```go
func ScanGitRemote(ctx context.Context, remoteURL string) ([]Finding, *Stats, error)
```

Shallow-clone a git repository (`--depth=1 --single-branch`)
into a temporary directory, scan all files, then remove the
clone. The context controls cancellation. Requires `git` in
`PATH`.

#### ScanGitRemoteFiles

```go
func ScanGitRemoteFiles(ctx context.Context, remoteURL string) ([]Finding, []string, *Stats, error)
```

Like `ScanGitRemote` but also returns the absolute paths
of every scanned file.

#### Render

```go
func Render(r Report, format Format) string
```

Render a `Report` to the requested format. Dispatches to
internal renderers for text, JSON, SARIF, and HTML output.

```go
output := core.Render(report, core.FormatJSON)
```

#### Register

```go
func Register(p Pattern)
```

Add `p` to the active registry. Safe to call before or
after bundle load. External patterns survive bundle loads
and appear alongside bundle patterns in `All()` calls.

#### All

```go
func All() []Pattern
```

Return a sorted snapshot of all registered patterns ordered
by `Name()`. The returned slice is owned by the caller.

#### EnablePattern

```go
func EnablePattern(name string) bool
```

Enable the named pattern, rebuild the active scanner set,
and persist the new state. Returns `false` if no pattern
with that name exists.

#### DisablePattern

```go
func DisablePattern(name string) bool
```

Disable the named pattern, rebuild the active scanner set,
and persist the new state. Returns `false` if no pattern
with that name exists.

#### SetActiveCategories

```go
func SetActiveCategories(cats []string)
```

Restrict subsequent scans to the named categories. A nil
or empty slice means "all categories". Rebuilds the internal
pre-filter regexes used by `ScanFile` and `ScanDir`.

#### ListEnabledPatterns

```go
func ListEnabledPatterns() []string
```

Return the names of every currently-enabled pattern, in
bundle order.

#### ListDisabledPatterns

```go
func ListDisabledPatterns() []string
```

Return the names of every currently-disabled pattern, in
bundle order.

#### EnableAllPatterns

```go
func EnableAllPatterns()
```

Enable every pattern in the bundle and rebuild the active
scanner set.

#### Categories

```go
func Categories() []string
```

Return the unique, sorted list of category labels in the
current bundle.

#### DownloadBundle

```go
func DownloadBundle(ctx context.Context) error
```

Fetch the latest pattern bundle from the configured URL,
compare it against the in-memory bundle, print a summary of
added/removed patterns, and persist the new bundle to
`~/.atheon/patterns.bundle`. The context controls the HTTP
request lifecycle.

#### SetBundleDownloadURL

```go
func SetBundleDownloadURL(url string) func()
```

Swap the upstream URL used by `DownloadBundle`. Returns a
restore function to reset the URL after tests or overrides.

#### Audit

```go
func Audit(ctx context.Context, root string) (*AuditReport, error)
```

Run all audit checks against `root` and return a structured
report. Checks: `nolint`, `todo-fixme`, `go-vet`,
`sentinel-errors`. The `dead-code` check is a stub;
enforcement is via `staticcheck` in CI. The context is
propagated to all subprocess calls (`grep`, `go vet`).

#### WriteReport

```go
func WriteReport(r *AuditReport, dir string) error
```

Write `r` as both `dir/REPORT.json` (pretty-printed JSON)
and `dir/REPORT.md` (GitHub-flavored Markdown with a
findings table). Creates `dir` if it does not exist.

## Bundler CLI

The `bundler` directory is a CLI tool, not a library. It is invoked
with `go run ./bundler/main.go <source-dir> <output-file>` and is
documented in `docs/ARCHITECTURE.md`. Its internals are not part of
the public API.

## What's not part of the API

- Anything under `internal/` is not importable.
- The CLI's flag set, exit codes, and output format are
  documented in `CLI.md`.
- The MCP server's tool surface is documented in `MCP.md`.
- The pattern YAML schema is documented in `PATTERNS.md`.

## Version history

- **v0.4** — Added `ScanURL`, `ScanGitRemote`,
  `ScanGitRemoteFiles`, `Report`, `Render`, `Format`,
  `Audit`, `WriteReport`, `AuditReport`, `AuditResult`,
  `AuditFinding`, `ScanEnv`, `EnablePattern`,
  `DisablePattern`, `SetPatternEnabled`,
  `SetActiveCategories`, `ListEnabledPatterns`,
  `ListDisabledPatterns`, `EnableAllPatterns`,
  `Register`, `All`, `Categories`,
  `DownloadBundle`, `SetBundleDownloadURL`.
- **v0.3 and earlier** — `ScanFile`, `ScanDir`,
  `ScanString`, `Finding`, `Stats`, `Pattern` interface,
  sentinel errors.
