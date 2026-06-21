#!/usr/bin/env bash
# scripts/audit-api-symbols.sh
#
# Verify that every public symbol documented in docs/API.md
# actually exists in the Go source code.
#
# Usage:
#   scripts/audit-api-symbols.sh           check all symbols
#   scripts/audit-api-symbols.sh --staged check only staged changes
#
# Exit code: 0 = clean, 1 = at least one missing symbol.

set -euo pipefail

MODE="${1:-all}"
DOCS_FILE="docs/API.md"

if [ ! -f "$DOCS_FILE" ]; then
    echo "SKIP: $DOCS_FILE not found"
    exit 0
fi

# Extract symbols from the code blocks in API.md
# We look at each "```go" block and the following lines

EXTRACTED_CORE=$(mktemp)
EXTRACTED_BUNDLER=$(mktemp)

awk '
/^```go$/ { in_block=1; next }
/^```$/ { in_block=0 }
in_block && /^func [A-Z]/ {
    sub(/^func /, ""); sub(/\(.*/, "");
    print $0 > "'"$EXTRACTED_CORE"'"
}
in_block && /^type [A-Z]/ {
    sub(/^type /, ""); sub(/ .*/, "");
    print $0 > "'"$EXTRACTED_CORE"'"
}
' "$DOCS_FILE"

# Also extract sentinel errors from var blocks
awk '
in_var && /^    Err[A-Z]/ {
    sub(/[[:space:]]*=.*/, "");
    gsub(/[[:space:]]*/, "");
    print $0 > "'"$EXTRACTED_CORE"'"
}
in_var && /^```$/ { in_var=0 }
/^    Err[A-Z]/ { in_var=1 }
in_block && /^    Err[A-Z]/ {
    sub(/[[:space:]]*=.*/, "");
    print $0 > "'"$EXTRACTED_CORE"'"
}
' "$DOCS_FILE"

# bundler package symbols (Config, Build)
awk '
/^```go$/ { in_block=1; next }
/^```$/ { in_block=0 }
in_block && /^type Config struct/ { print "Config" > "'"$EXTRACTED_BUNDLER"'" }
in_block && /^func Build\(/ { print "Build" > "'"$EXTRACTED_BUNDLER"'" }
' "$DOCS_FILE"

# Dedupe
sort -u "$EXTRACTED_CORE" -o "$EXTRACTED_CORE"
sort -u "$EXTRACTED_BUNDLER" -o "$EXTRACTED_BUNDLER"

FAILED=0

check_core_symbol() {
    local sym="$1"
    case "$sym" in
    # Types - check for type definition
    Finding|Stats|Report|Format|AuditReport|AuditResult|AuditFinding|AuditSummary|Pattern|PatternState)
        grep -qE "^type $sym (struct|interface)" core/*.go 2>/dev/null
        ;;
    # Sentinel errors
    ErrPatternNotFound|ErrBundleDownload|ErrBundleParse|ErrInvalidPattern)
        grep -qE "^[[:space:]]*$sym[[:space:]]*=" core/*.go 2>/dev/null
        ;;
    # Format is a type string
    Format)
        grep -qE "^type Format string" core/*.go 2>/dev/null
        ;;
    # Functions
    Audit|WriteReport|ScanFile|ScanDir|ScanString|ScanEnv|ScanURL|ScanGitRemote|ScanGitRemoteFiles|Register|All|Categories|EnablePattern|DisablePattern|SetPatternEnabled|ListEnabledPatterns|ListDisabledPatterns|EnableAllPatterns|SetActiveCategories|SetBundleDownloadURL|DownloadBundle|Render)
        grep -qE "^func $sym\(" core/*.go 2>/dev/null
        ;;
    # Skip known stdlib names that appear in docs
    fmt|os|regexp|strings|time|error)
        return 0
        ;;
    *)
        grep -qE "^func $sym\(|^type $sym |^var $sym " core/*.go 2>/dev/null
        ;;
    esac
}

check_bundler_symbol() {
    local sym="$1"
    case "$sym" in
    Config)
        grep -qE "^type Config struct" bundler/*.go 2>/dev/null
        ;;
    Build)
        grep -qE "^func Build\(" bundler/*.go 2>/dev/null
        ;;
    *)
        return 1
        ;;
    esac
}

for sym in $(cat "$EXTRACTED_CORE"); do
    if ! check_core_symbol "$sym"; then
        echo "MISSING_SYMBOL: $sym (in core/ package)"
        FAILED=1
    fi
done

for sym in $(cat "$EXTRACTED_BUNDLER"); do
    if ! check_bundler_symbol "$sym"; then
        echo "MISSING_SYMBOL: $sym (in bundler/ package)"
        FAILED=1
    fi
done

rm -f "$EXTRACTED_CORE" "$EXTRACTED_BUNDLER"

if [ "$FAILED" -ne 0 ]; then
    echo "API symbol audit failed"
    exit 1
fi
echo "OK"
exit 0
