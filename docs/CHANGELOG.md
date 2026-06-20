# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **Reading this file:** versions are listed newest-first. Each release
> section groups changes under `Added`, `Changed`, `Deprecated`,
> `Removed`, `Fixed`, and `Security` so you can scan for what affects
> you. Entries marked **Breaking** change the public API or the
> default CLI behaviour and require a code or config update on your
> side.

## Unreleased

### Added

- `Makefile` — single source of truth for build, test, coverage,
  vet, and four audit gates. Run `make help` for the catalogue;
  `make audit` for the full quality sweep. The audit gates are
  `audit-dead-code`, `audit-nolint`, `audit-fixmes`, and
  `audit-sentinels`.
- `scripts/audit-dead-code.sh` — the implementation shared by
  `make audit-dead-code` and the new pre-commit gate. Detects
  unexported helpers with no callers in the whole tree.
- `.github/SUPPORT.md`, `.github/FUNDING.yml`, `.github/CODEOWNERS`,
  `.github/dependabot.yml`, `.github/ISSUE_TEMPLATE/*.md`, and
  `.github/PULL_REQUEST_TEMPLATE.md` — first-class GitHub
  community-health files so Discussions, sponsor buttons, and
  automatic review requests work out of the box.
- `docs/CHANGELOG.md`, `docs/ROADMAP.md`, `docs/ARCHITECTURE.md`,
  `docs/DESIGN.md`, `docs/TROUBLESHOOTING.md`, `docs/UPGRADE.md`,
  `docs/MIGRATION.md`, `docs/GOVERNANCE.md`, `docs/MAINTAINERS.md`,
  `docs/RELEASE.md`, `docs/PERFORMANCE.md`, `docs/PROJECTS.md`,
  `docs/STANDARDS.md` — consolidated project documentation.
- `docs/audits/DEAD_CODE_AUDIT.md` — Phase 1.1 findings: 0 staticcheck
  issues, 0 go vet issues, 1 dead helper (`contains`) — now removed.
- Global PreToolUse hook
  (`~/.claude/hooks/dead-code-prevention-hook.py`) — wired into
  `~/.claude/settings.json`. Blocks the agent from writing a Go
  file that introduces an unexported helper with no callers.
  Companion to the project-local pre-commit gate.

### Changed

- Moved `CODE_OF_CONDUCT.md` and `SECURITY.md` from the repository
  root to `.github/` so GitHub picks them up automatically. **No
  content change** — only the location moved.
- `.githooks/pre-commit` — the TODO/FIXME gate's helper-extraction
  pipeline was killed by `set -euo pipefail` whenever a commit
  staged no new `.go` files. Added `|| true` so the script runs
  to completion. The new `🪦 Checking for dead code in staged
  files` step also lives here.

### Fixed

- CI: `TestBundleReadFileError` and `TestBundleWriteFileError` now
  skip cleanly on Windows instead of failing, because Windows
  `chmod` cannot restrict the file owner.
- CI: `TestScanFileIgnored` and `TestRunDefaults` use `t.Chdir` so
  the working-directory change is restored via `t.Cleanup` instead
  of a deferred `os.Chdir`, removing a `-race` flake on macOS.
- CI: `TestMain*` integration tests now build the binary once in
  `TestMain` into a `t.TempDir` path and share it across tests,
  eliminating the `executable file not found in $PATH` failure
  when the runner relocates the package working directory.
- CI: security-scan workflow no longer reports the deliberate
  example credentials in `cmd/atheon/testdata/FIXTURES.txt` as
  leaks after `.atheonignore` was extended to cover `testdata/`
  and `fixtures/` directories.
- CI: `sync-stable-clean.yml` detects which layout is checked out
  (root `main.go` vs. `cmd/atheon`) before validating patterns.
- Removed dead `contains` helper in `core/bundle.go` (zero
  production callers; matches upstream #159) and the
  `TestContains` test that was its only consumer.

## 1.0.0 — 2026-06-19

### Added

- Comprehensive Go best-practices improvements: `context.Context`
  propagation through the public API, structured error sentinels
  (`ErrBinaryFile`, `ErrCancelled`, `ErrInvalidPath`, `ErrTooLarge`,
  `ErrSecretInEnv`, `ErrSecretInStdin`), 18 golangci-lint rules,
  godoc `Example*` functions in `core/example_test.go` that double
  as documentation and as runnable tests under `go test -run=Example`.
- Pattern bundle expanded to **179 patterns** across **19
  categories** (secrets, pii, finance, healthcare, code-quality,
  accessibility, networking, cloud, devops, …).

### Fixed

- `core.ScanFile`, `core.ScanDir`, `core.ScanString`, `core.ScanEnv`
  now honor `.atheonignore` and `.gitignore` consistently.
- Self-scan workflow and pre-commit hook no longer false-positive on
  the project's own embedded pattern bundle.

## 0.x — Pre-1.0 Development

The 0.x series predates this changelog. See `git log` and the
[release notes on GitHub](../../releases) for the historical record.
Highlights of the 0.x period include the introduction of the
embedded pattern bundle, the `--categories` and `--all` flags, the
`enable`/`disable` state-persistence flow, and the MCP server entry
point.

---

## How to read this file

- The **Unreleased** section is where contributors add entries as
  they land PRs. Maintainers move them under a versioned heading at
  release time.
- Look for the **Breaking** tag on any `Changed` or `Removed` entry
  — those need a code or config change on your side before upgrading.
- The full commit history is in `git log`; this file exists so you
  don't have to read every commit to find out what changed.
