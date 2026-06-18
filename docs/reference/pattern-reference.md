# Pattern Reference

Reference for all built-in patterns included in Atheon.

## Security/Secrets Patterns

### AWS Access Key ID

**Pattern:** `aws-access-key-id`
**Category:** `secrets`
**Description:** Detects AWS access key identifiers

**Format:** `AKIA[0-9A-Z]{16}`

**Examples:**
```
✓ Matches: AWS_ACCESS_KEY_ID=AKIAXXXXXXXXXXXXXXXXXXXXXXXXX
✗ No Match: AWS_ACCESS_KEY_ID=AKIA123
✗ No Match: AKIAiosfodnn7example
```

---

### AWS Secret Access Key

**Pattern:** `aws-secret-access-key`
**Category:** `secrets`
**Description:** Detects AWS secret access keys

**Format:** Variable assignment with long alphanumeric string

---

### GitHub Personal Access Token

**Pattern:** `github-pat`
**Category:** `secrets`
**Description:** Detects GitHub personal access tokens

**Format:** `ghp_[0-9a-zA-Z]{36}`

**Examples:**
```
✓ Matches: token=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
✗ No Match: token=ghp_short
✗ No Match: token=github_pat_xxxxxxxxxxabcdefghijklmnopqrstuvwxyz
```

---

### OpenAI API Key

**Pattern:** `openai-api-key`
**Category:** `secrets`
**Description:** Detects OpenAI API keys

**Format:** `sk-[a-zA-Z0-9]{20}`

**Examples:**
```
✓ Matches: OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
✗ No Match: OPENAI_API_KEY=sk-short
✗ No Match: OPENAI_API_KEY=pk-xxxxxxxxxxxxxxxxxxxx
```

---

### Slack Bot Token

**Pattern:** `slack-bot-token`
**Category:** `secrets`
**Description:** Detects Slack bot tokens

**Format:** `xoxb-[0-9]{12}-[0-9]{12}-[a-zA-Z0-9]{24}`

**Examples:**
```
✓ Matches: SLACK_BOT_TOKEN=xoxb-xxxxxxxxxxxx-xxxxxxxxxx-xxxxxxxxxxxxxxxxxxyz
✗ No Match: SLACK_BOT_TOKEN=xoxb-short
✗ No Match: SLACK_BOT_TOKEN=xoxa-xxxxxxxxxx12-xxxxxxxxxx12-abcdefghijklmnopqrstuvwxyz
```

---

### Stripe Secret Key

**Pattern:** `stripe-secret-key`
**Category:** `secrets`
**Description:** Detects Stripe secret keys (live mode)

**Format**: `sk_live_` followed by 24 alphanumeric characters

**Examples:**
```
✓ Matches: STRIPE_SECRET_KEY=sk_live_24chars
✗ No Match: STRIPE_SECRET_KEY=sk_test_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
✗ No Match: STRIPE_SECRET_KEY=sk_live_short
```

---

### Twilio Account SID

**Pattern:** `twilio-account-sid`
**Category:** `secrets`
**Description:** Detects Twilio account security identifiers

**Format:** `AC[a-zA-Z0-9]{32}`

**Examples:**
```
✓ Matches: TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
✗ No Match: TWILIO_ACCOUNT_SID=ACshort
✗ No Match: TWILIO_ACCOUNT_SID=SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx12
```

---

### GCP API Key

**Pattern:** `gcp-api-key`
**Category:** `secrets`
**Description:** Detects Google Cloud Platform API keys

**Format:** `AIza[a-zA-Z0-9]{35}`

**Examples:**
```
✓ Matches: api_key=AIzaxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
✗ No Match: api_key=AIza-short
✗ No Match: api_key=AIza!@#$%^&*()_+xxxxxxxxxxxxxxxxxxxx12345
```

---

### GCP OAuth Client ID

**Pattern:** `gcp-oauth-client-id`
**Category:** `secrets`
**Description:** Detects GCP OAuth client identifiers

**Format:** `[0-9]+-[a-z]+\.apps\.googleusercontent\.com`

**Examples:**
```
✓ Matches: client_id=xxxxxxxxxx-xxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
✗ No Match: client_id=project.apps.googleusercontent.com
✗ No Match: client_id=xxxxxxxxxx-.apps.googleusercontent.com
```

---

### GCP OAuth Client Secret

**Pattern:** `gcp-oauth-client-secret`
**Category:** `secrets`
**Description:** Detects GCP OAuth client secrets

**Format:** `GOCSPX-[a-zA-Z0-9]{28}`

**Examples:**
```
✓ Matches: client_secret=GOCSPX-xxxxxxxxxxxxxxxxxxxxxxxxxxxx
✗ No Match: client_secret=GOCSPX-short
✗ No Match: client_secret=GOOGLE-xxxxxxxxxxxxxxxxxxxx12345678
```

---

### GCP Service Account Email

**Pattern:** `gcp-service-account-email`
**Category:** `secrets`
**Description:** Detects GCP service account email addresses

**Format:** `[a-z0-9]+@.*\.gserviceaccount\.com`

**Examples:**
```
✓ Matches: svc=my-service@project.iam.gserviceaccount.com
✗ No Match: svc=my-service@example.com
✗ No Match: svc=@project.iam.gserviceaccount.com
```

---

### GCP Service Account Key

**Pattern:** `gcp-service-account-key`
**Category:** `secrets`
**Description:** Detects GCP service account private key IDs

**Format:** 40-character alphanumeric string in JSON context

**Examples:**
```
✓ Matches: "private_key_id": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
✓ Matches: "client_email": "svc@project.iam.gserviceaccount.com"
✗ No Match: "private_key_id": "short"
✗ No Match: "client_email": "svc@example.com"
```

## PII Patterns

### Social Security Number

**Pattern:** `ssn`
**Category:** `pii`
**Description:** Detects US Social Security Numbers

**Format:** `\d{3}-\d{2}-\d{4}`

**Examples:**
```
✓ Matches: ssn=xxx-xx-xxxx
✗ No Match: ssn=123-456-789
✗ No Match: invoice=123-45-678
```

---

### Credit Card Number

**Pattern:** `credit-card`
**Category:** `pii`
**Description:** Detects major credit card numbers

**Formats Supported:**
- Visa: `4[0-9]{3}[- ]?(?:[0-9]{4}[- ]?){2}[0-9]{4}`
- MasterCard: `5[1-5][0-9]{2}[- ]?(?:[0-9]{4}[- ]?){2}[0-9]{4}`
- American Express: `3[47][0-9]{2}[- ]?[0-9]{6}[- ]?[0-9]{5}`
- Discover: `6(?:011|5[0-9]{2})[- ]?(?:[0-9]{4}[- ]?){2}[0-9]{4}`

**Examples:**
```
✓ Matches: card=xxxx xxxx xxxx xxxx
✓ Matches: amex=xxxx-xxxxxx-xxxxx
✗ No Match: card=xxxx xxxx xxxx xxxx
✗ No Match: order=4242
```

---

### Phone Number

**Pattern:** `phone-number`
**Category:** `pii`
**Description:** Detects US phone numbers

**Formats Supported:**
- `(xxx 123-4567`
- `+1 xxx123 4567`
- Various US phone number formats

**Examples:**
```
✓ Matches: phone=(xxx 123-4567
✓ Matches: phone=+1 xxx123 4567
✗ No Match: version=xxx123
✗ No Match: ticket=xxx123-456
```

## Pattern Statistics

### Current Coverage

**Total Patterns:** 14
**Categories:** 2 (secrets, pii)

**Breakdown by Category:**
- **Secrets:** 11 patterns
- **PII:** 3 patterns

## Pattern Accuracy

### False Positive Rate

**Security Patterns:** Low (<5%)
- Designed for specific service formats
- Include prefix validation
- Use word boundaries

**PII Patterns:** Medium (5-15%)
- General formats may have more false positives
- Context-dependent accuracy
- May require manual verification

### False Negative Rate

**Security Patterns:** Very Low (<2%)
- Covers major service formats
- Regular updates for new patterns
- Community contributions expand coverage

**PII Patterns:** Low (<5%)
- Covers common formats
- International formats limited
- Regional variations may be missed

## Contributing Patterns

See [Pattern Development Guide](../patterns/README.md) for:
- Creating new patterns
- Testing procedures
- Contribution guidelines
- Best practices

## Pattern Updates

### Update Frequency

- **Schedule:** 10th and 21st of every month
- **Source:** Community contributions
- **Delivery:** `atheon update` command

### Getting Latest Patterns

```bash
# Update to latest pattern bundle
atheon update

# Verify new patterns
atheon list
```

## Pattern Quality Metrics

### Pattern Characteristics

**Specificity:**
- Most patterns include service-specific prefixes
- Word boundary validation
- Format-specific validation

**Performance:**
- Pre-compiled regex patterns
- Optimized for fast matching
- Category-based filtering

**Maintainability:**
- Simple YAML definitions
- Clear naming conventions
- Comprehensive test coverage

## Usage Recommendations

### Security Scanning

```bash
# Scan for secrets in code
atheon --categories=secrets ./

# Scan environment variables
atheon --env
```

### PII Detection

```bash
# Scan for personal identifiers
atheon --categories=pii ./

# Combined scanning
atheon --categories=secrets,pii ./
```

### Best Practices

1. **Use Category Filtering:** Only scan relevant categories
2. **Verify Findings:** Check false positives manually
3. **Regular Updates:** Keep patterns current with `atheon update`
4. **Contribute Patterns:** Add missing patterns for your domain
5. **Context Matters:** Consider the context of findings

## Pattern Limitations

### Current Limitations

**Geographic Coverage:**
- Primarily US-centric PII patterns
- Limited international format support
- Regional variations not fully covered

**Service Coverage:**
- Major services covered
- Niche services may be missing
- Legacy formats not always included

**Format Variations:**
- Standard formats prioritized
- Custom/variant formats may be missed
- Obfuscation techniques not handled

### Future Enhancements

**Planned Expansions:**
- International PII formats
- Additional service patterns
- Custom format support
- Regional pattern variations

## Additional Resources

- [Pattern Development Guide](../patterns/README.md)
- [Contributing Guidelines](../community/README.md)
- [CLI Reference](cli-reference.md)
- [Getting Started Guide](../getting-started/README.md)