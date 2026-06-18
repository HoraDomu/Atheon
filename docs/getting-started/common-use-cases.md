# Common Use Cases for Atheon

## Security Scanning

### Preventing Secret Leaks

**Scenario:** Developers accidentally committing API keys, tokens, or credentials.

**Solution:** Pre-commit hook integration
```bash
# .git/hooks/pre-commit
#!/bin/sh
atheon --categories=secrets ./
if [ $? -ne 0 ]; then
  echo "❌ Secrets detected! Remove them before committing."
  exit 1
fi
```

### Environment Variable Scanning

**Scenario:** CI/CD pipelines with injected secrets at runtime.

**Solution:** Environment scanning
```bash
# In CI/CD pipeline
atheon --env
```

### Incident Response

**Scenario:** Investigating potential secret exposure in repository history.

**Solution:** Scan git history
```bash
# Scan all commits
git log --all --pretty=format:%H | while read commit; do
  git show $commit | atheon - && echo "Clean: $commit" || echo "Found in: $commit"
done
```

## Compliance & PII Detection

### Healthcare (HIPAA)

**Scenario:** Medical records containing patient identifiers.

**Solution:** PII category scanning
```bash
atheon --categories=pii ./patient-records/
```

### Financial Compliance

**Scenario:** Code containing credit card numbers or account details.

**Solution:** Targeted scanning
```bash
atheon --categories=finance ./payment-processing/
```

### GDPR Compliance

**Scenario:** EU citizen data in code repositories.

**Solution:** Comprehensive PII scanning
```bash
atheon --categories=pii,secrets ./eu-data/
```

## Development Workflows

### Code Review Automation

**Scenario:** PR reviewers manually checking for security issues.

**Solution:** CI integration
```yaml
# .github/workflows/security-scan.yml
name: Security Scan
on: [pull_request]
jobs:
  atheon:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Scan for secrets
        run: |
          curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-amd64 -o atheon
          chmod +x atheon
          ./atheon --json ./ > findings.json
          if [ -s findings.json ]; then
            echo "Security findings detected!"
            cat findings.json
            exit 1
          fi
```

### Local Development

**Scenario:** Developers want instant feedback on security issues.

**Solution:** File watcher integration
```bash
# Using entr (file watcher)
find . -name "*.go" -o -name "*.js" | entr -r atheon ./
```

### IDE Integration

**Scenario:** Real-time feedback while coding.

**Solution:** MCP server with AI tools
```json
{
  "mcpServers": {
    "atheon": {
      "command": "atheon-mcp",
      "env": {
        "ATHEON_CATEGORIES": "secrets,pii"
      }
    }
  }
}
```

## Automation & Pipelines

### CI/CD Quality Gates

**Scenario:** Prevent deployment of code with security issues.

**Solution:** Pipeline integration
```yaml
# GitLab CI
security_scan:
  stage: test
  script:
    - atheon --json ./ > findings.json
  artifacts:
    reports:
      security: findings.json
  only:
    - merge_requests
    - main
```

### DevOps Automation

**Scenario:** Infrastructure-as-code with secrets.

**Solution:** Terraform/Ansible scanning
```bash
atheon --categories=secrets ./terraform/
atheon --categories=secrets ./ansible/
```

### Log Analysis

**Scenario:** Detecting secrets in application logs.

**Solution:** Log scanning pipeline
```bash
# Scan rotated logs
find /var/log/app -name "*.log" | while read log; do
  echo "Scanning $log"
  atheon --file "$log"
done
```

## Data Migration & Cleanup

### Legacy Code Scanning

**Scenario:** Old repositories with unknown security issues.

**Solution:** Comprehensive scanning
```bash
# Deep scan with reporting
atheon --json --all ./legacy-repo/ > security-report.json
```

### Database Dump Analysis

**Scenario:** Database exports containing sensitive data.

**Solution:** SQL file scanning
```bash
atheon --categories=pii,secrets ./dumps/
```

### Configuration Validation

**Scenario:** Kubernetes configs with secrets.

**Solution:** YAML scanning
```bash
atheon --categories=secrets ./k8s/
```

## Monitoring & Alerting

### Continuous Scanning

**Scenario:** Regular security scans of active repositories.

**Solution:** Cron job integration
```bash
# Weekly security scans
0 2 * * 0 cd /repos && atheon --json ./ > /var/log/atheon-$(date +%Y%m%d).json
```

### Slack Notifications

**Scenario:** Team wants instant alerts on findings.

**Solution:** Integration with notification tools
```bash
#!/bin/bash
findings=$(atheon --json ./)
if [ ! -z "$findings" ]; then
  curl -X POST -H 'Content-type: application/json' \
    --data "{\"text\":\"🚨 Atheon found security issues!\"}" \
    $SLACK_WEBHOOK_URL
fi
```

## Advanced Patterns

### Custom Domain Detection

**Scenario:** Industry-specific pattern requirements.

**Solution:** Create custom patterns
```yaml
# community/secrets/your-service.yaml
name: your-service-api-key
match: '\byourservice_[A-Za-z0-9]{32}\b'
```

### Multi-language Projects

**Scenario:** Different patterns for different languages.

**Solution:** Directory-based scanning
```bash
# Scan backend
atheon --categories=secrets ./backend/

# Scan frontend
atheon --categories=secrets,pii ./frontend/
```

## Troubleshooting Common Issues

### False Positives

**Issue:** Test fixtures triggering patterns.

**Solution:** Ignore files
```bash
# .atheonignore
test/
fixtures/
*_test.go
```

### Performance Issues

**Issue:** Large repositories taking too long.

**Solution:** Category filtering
```bash
atheon --categories=secrets ./  # Only scan what you need
```

### Missing Patterns

**Issue:** Expected pattern not detected.

**Solution:** Update and verify
```bash
atheon update
atheon list | grep pattern-name
```

## Best Practices

1. **Start with category filtering** - Use only the categories you need
2. **Integrate early** - Add to pre-commit hooks before code review
3. **Automate alerts** - CI/CD integration for team awareness
4. **Regular updates** - Run `atheon update` monthly
5. **Custom patterns** - Contribute domain-specific patterns back to community