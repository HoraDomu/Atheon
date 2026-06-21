#!/usr/bin/env bash
# scripts/lint-community-patterns.sh
#
# Lint all community pattern YAML files for structural validity.
# Checks each *.yaml file under community/:
#   1. Has a top-level "name:" field
#   2. Has a top-level "match:" field
#   3. The regex in "match:" compiles successfully with Go's regexp engine
#
# Usage:
#   scripts/lint-community-patterns.sh           lint all patterns
#   scripts/lint-community-patterns.sh --staged   lint only staged files
#
# Exit code: 0 = clean, 1 = at least one lint error.

set -euo pipefail

MODE="${1:-all}"

case "$MODE" in
--staged)
    YAMLS=$(git diff --cached --name-only --diff-filter=AM 2>/dev/null \
        | grep '\.yaml$' | grep 'community/' || true)
    ;;
all)
    YAMLS=$(find community -name '*.yaml' 2>/dev/null || true)
    ;;
*)
    echo "usage: $0 [--staged]" >&2
    exit 2
    ;;
esac

if [ -z "$YAMLS" ]; then
    echo "OK (no pattern files to lint)"
    exit 0
fi

FAILED=0

for yaml in $YAMLS; do
    [ -f "$yaml" ] || continue

    # Check for required name field
    if ! grep -qE '^name:' "$yaml"; then
        echo "MISSING_NAME: $yaml"
        FAILED=1
        continue
    fi

    # Check for required match field
    if ! grep -qE '^match:' "$yaml"; then
        echo "MISSING_MATCH: $yaml"
        FAILED=1
        continue
    fi

    # Extract the regex value (handles quoted and unquoted values)
    MATCH_VAL=$(grep '^match:' "$yaml" | sed -E 's/^match:\s*//' | head -1)
    # Remove surrounding quotes if present
    MATCH_VAL=$(echo "$MATCH_VAL" | sed 's/^["\x27]//' | sed 's/["\x27]$//')

    if [ -z "$MATCH_VAL" ]; then
        echo "EMPTY_MATCH: $yaml"
        FAILED=1
        continue
    fi

    # Validate regex using Go's regexp engine
    # Write a small Go program to a temp file and run it
    LINTDIR=$(mktemp -d)
    cat > "$LINTDIR/lint.go" <<'GOEOF'
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	raw := os.Args[1]
	// Strip optional /.../ delimiters sometimes used in docs
	raw = strings.TrimPrefix(raw, "/")
	raw = strings.TrimSuffix(raw, "/")
	if _, err := regexp.Compile(raw); err != nil {
		fmt.Fprintf(os.Stderr, "INVALID_REGEX: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
GOEOF

    if ! go run "$LINTDIR/lint.go" "$MATCH_VAL" 2>&1; then
        echo "BAD_REGEX: $yaml: match=$MATCH_VAL"
        FAILED=1
    fi

    rm -rf "$LINTDIR"
done

if [ "$FAILED" -ne 0 ]; then
    echo "Pattern lint failed"
    exit 1
fi
echo "OK"
exit 0
