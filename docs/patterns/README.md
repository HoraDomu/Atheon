# Pattern Development Guide

Patterns are the heart of Atheon. Every pattern is one YAML file — no Go required, no engine changes, fast to review, and immediately useful to every user once merged.

## Pattern Basics

### Pattern Structure

```yaml
# community/secrets/my-service.yaml
name: my-service-api-key
match: '\bmsvc_[A-Za-z0-9]{32}\b'
```

**Required Fields:**
- `name`: Pattern identifier (lowercase-hyphenated, specific)
- `match`: Valid RE2 regex pattern

**Optional Features:**
- Case sensitivity depends on regex
- Use word boundaries `\b` to prevent partial matches
- Use character classes for precise matching

## Pattern Categories

### Existing Categories

**Security/Secrets:**
- API keys, tokens, credentials
- Authentication secrets
- Service-specific keys

**PII (Personally Identifiable Information):**
- Social security numbers
- Credit card numbers
- Phone numbers
- Personal identifiers

### Creating New Categories

1. **Create folder** in `community/`:
   ```bash
   mkdir community/finance
   ```

2. **Add patterns** to the folder:
   ```bash
   # community/finance/routing-number.yaml
   name: aba-routing-number
   match: '\b\d{9}\b'
   ```

3. **Rebuild bundle:**
   ```bash
   go run ./bundler
   ```

The folder name becomes the category name.

## Writing Patterns

### Pattern Best Practices

**1. Be Specific**
```yaml
# Good: Specific pattern
name: github-pat
match: '\bghp_[0-9a-zA-Z]{36}\b'

# Avoid: Too generic
name: token
match: '\b[a-zA-Z0-9]{32}\b'
```

**2. Use Word Boundaries**
```yaml
# Good: Prevents partial matches
match: '\bapi-key-[A-Za-z0-9]{32}\b'

# Problematic: May match within larger strings
match: 'api-key-[A-Za-z0-9]{32}'
```

**3. Match Service Prefixes**
```yaml
# Good: Service-specific prefix
name: aws-access-key
match: '\bAKIA[0-9A-Z]{16}\b'

# Better: Include format indicator
name: openai-api-key
match: '\bsk-[a-zA-Z0-9]{20}\b'
```

**4. Consider Context**
```yaml
# Good: Environment variable context
name: api-key-env
match: '(?i)API[_-]?KEY\s*[=:]\s*[a-zA-Z0-9]{32}'

# Or: Code assignment context
name: api-key-assignment
match: 'api[_-]?key\s*=\s*["\']?[a-zA-Z0-9]{32}'
```

### Pattern Examples

**API Keys:**
```yaml
# GitHub Personal Access Token
name: github-pat
match: '\bghp_[0-9a-zA-Z]{36}\b'

# AWS Access Key
name: aws-access-key
match: '\bAKIA[0-9A-Z]{16}\b'

# OpenAI API Key
name: openai-api-key
match: '\bsk-[a-zA-Z0-9]{20}\b'
```

**Tokens:**
```yaml
# Slack Bot Token
name: slack-bot-token
match: '\bxoxb-[0-9]{12}-[0-9]{12}-[a-zA-Z0-9]{24}\b'

# Stripe Secret Key
name: stripe-secret-key
match: '\bsk_live_[a-zA-Z0-9]{24}\b'
```

**PII:**
```yaml
# Social Security Number
name: ssn
match: '\b\d{3}-\d{2}-\d{4}\b'

# Credit Card Number
name: credit-card
match: '\b(?:4[0-9]{3}[- ]?){3}[0-9]{4}\b'
```

### Advanced Patterns

**Conditional Matching:**
```yaml
# Match with context
name: database-url
match: '(?i)database[_-]?url\s*=\s*[a-zA-Z0-9:]+@[^/]+/[^\s]+'
```

**Format-Specific:**
```yaml
# JSON key detection
name: json-api-key
match: '"api[_-]?key"\s*:\s*"[a-zA-Z0-9]{32}"'

# YAML key detection
name: yaml-api-key
match: 'api[_-]?key\s*:\s*[a-zA-Z0-9]{32}'
```

## Testing Patterns

### Creating Test Cases

Add test cases to `core/bundle_test.go`:

```go
"github-pat": {
    matches:    []string{
        "token=ghp_" + strings.Repeat("a", 36),
        "Authorization: ghp_" + strings.Repeat("b", 36),
    },
    nonMatches: []string{
        "token=ghp_short",
        "token=" + strings.Repeat("a", 36), // Missing prefix
        "token=github_" + strings.Repeat("a", 36), // Wrong prefix
    },
},
```

### Testing Process

1. **Create pattern YAML file**
2. **Add test cases** to bundle_test.go
3. **Run tests:**
   ```bash
   go test ./...
   ```
4. **Verify manually:**
   ```bash
   ./atheon --file test-sample.txt
   ```

### Manual Testing

Create test files and verify:

```bash
# Create test file
cat > test.txt << EOF
api_key=sk-test12345678901234567890
# Should not match below
api_key=sk-short
EOF

# Test
./atheon --file test.txt

# Should detect openai-api-key pattern
```

## Pattern Submission

### Contribution Workflow

1. **Check existing patterns:**
   ```bash
   ./atheon list
   ./atheon list categories
   ```

2. **Create pattern file:**
   ```bash
   # community/secrets/my-service.yaml
   name: my-service-api-key
   match: '\bmsvc_[A-Za-z0-9]{32}\b'
   ```

3. **Rebuild bundle:**
   ```bash
   go run ./bundler
   ```

4. **Add test cases:**
   ```go
   "my-service-api-key": {
       matches:    []string{"key=msvc_" + strings.Repeat("a", 32)},
       nonMatches: []string{"key=msvc_short", "key=other_" + strings.Repeat("a", 32)},
   },
   ```

5. **Test and verify:**
   ```bash
   go test ./...
   ./atheon list | grep my-service-api-key
   ```

6. **Create pull request** with:
   - Pattern description
   - What it detects
   - Why it matters
   - Test cases used
   - Sample matches/non-matches

### Pull Request Template

```markdown
## Pattern Description
Add detection for MyService API keys

## What it detects
MyService API keys in the format: `msvc_[32 alphanumeric characters]`

## Why it matters
These keys can be used to access MyService APIs and should not be committed to code

## Test cases
**Matches:**
- `API_KEY=msvc_1234567890abcdefghijklmnopqrstuvwxyz`
- `myServiceToken: msvc_abcdefghijklmnopqrstuvwxyz1234567890`

**Non-matches:**
- `API_KEY=msvc_short` (too short)
- `API_KEY=other_1234567890abcdefghijklmnopqrstuvwxyz` (wrong prefix)

## Files changed
- community/secrets/my-service.yaml (new pattern)
- core/patterns.bundle (updated bundle)
- core/bundle_test.go (test cases)
```

## Pattern Quality Criteria

### Review Checklist

**Pattern Correctness:**
- ✅ Detects intended format accurately
- ✅ Minimal false positives
- ✅ Minimal false negatives
- ✅ Handles edge cases

**Pattern Specificity:**
- ✅ Clear, descriptive name
- ✅ Service-specific when applicable
- ✅ Uses word boundaries appropriately
- ✅ Considers context (environment variables, code, configs)

**Testing Coverage:**
- ✅ Positive test cases (what should match)
- ✅ Negative test cases (what shouldn't match)
- ✅ Edge cases covered
- ✅ Real-world samples tested

**Documentation:**
- ✅ Clear pattern purpose
- ✅ Service/format described
- ✅ Test cases documented
- ✅ Contribution notes included

## Common Pattern Issues

### False Positives

**Problem:** Matches unintended content
```yaml
# Problematic: Matches random hex strings
name: hex-string
match: '\b[a-f0-9]{32}\b'
```

**Solution:** Add service-specific prefix
```yaml
# Better: Service-specific detection
name: service-token
match: '\bservicetoken_[a-f0-9]{32}\b'
```

### False Negatives

**Problem:** Misses valid patterns
```yaml
# Problematic: Too strict
name: api-key
match: '\bapi[_-]?key\s*=\s*[a-zA-Z0-9]{32}\b'
```

**Solution:** Handle variations
```yaml
# Better: Multiple formats
name: api-key
match: '\bapi[_-]?key\s*[=:]\s*[a-zA-Z0-9]{32}\b'
```

### ReDoS Prevention

**Problem:** Catastrophic backtracking
```yaml
# Dangerous: Nested quantifiers
match: '(a+)+'
```

**Solution:** Use atomic groups or possessive quantifiers
```yaml
# Safe: Simple character classes
match: '\b[a-zA-Z0-9]{32}\b'
```

## Pattern Categories Guide

### Security/Secrets
- API keys and tokens
- Credentials and authentication secrets
- Service-specific keys
- Certificate keys

### PII
- Personal identifiers
- Financial information
- Healthcare data
- Contact information

### Finance (Future)
- Account numbers
- Routing numbers
- Credit card formats
- Financial codes

### Healthcare (Future)
- Patient IDs
- Medical record numbers
- Insurance numbers
- Health data codes

### Legal (Future)
- Contract references
- Case numbers
- Legal document identifiers
- Regulation codes

## Pattern Maintenance

### Updating Patterns

If a pattern format changes:

1. **Update the YAML** file
2. **Add new test cases**
3. **Run full test suite:**
   ```bash
   go test ./...
   ```
4. **Update bundle:**
   ```bash
   go run ./bundler
   ```
5. **Document changes** in commit message

### Deprecating Patterns

If a pattern is no longer needed:

1. **Add deprecation notice** in pattern name/description
2. **Create issue** to discuss removal
3. **Remove after consensus** and timeline agreement

## Additional Resources

- [RE2 Syntax](https://github.com/google/re2/wiki/Syntax)
- [Regex Tutorial](https://regexone.com/)
- [Security Patterns Guide](https://owasp.org/www-community/attacks/Regular_expression_Denial_of_Service_-_ReDoS)
- [Contributing Guidelines](../development/README.md)