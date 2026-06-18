# Integration Guide

This guide covers integrating Atheon into various development workflows, CI/CD pipelines, and tooling ecosystems.

## Development Environment Integration

### IDE Integration

#### VS Code

**Option 1: Tasks**
```json
// .vscode/tasks.json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Atheon Scan",
      "type": "shell",
      "command": "atheon",
      "args": ["./"],
      "problemMatcher": []
    }
  ]
}
```

**Option 2: Extensions**
- Install security scanning extensions that support external tools
- Configure extension to run Atheon on file save

#### JetBrains IDEs (IntelliJ, PyCharm, etc.)

**External Tools:**
1. Settings → Tools → External Tools → Add
2. Configuration:
   - Name: `Atheon Scan`
   - Program: `atheon`
   - Arguments: `$ProjectFileDir$`
   - Working directory: `$ProjectFileDir$`

#### Vim/Neovim

**Quick Command:**
```vim
" Run Atheon on current project
command! Atheon execute '!atheon ' . expand('%:p:h')
```

**Automatic on Write:**
```vim
" Run Atheon after writing files
autocmd BufWritePost * :!atheon %
```

#### Emacs

**Add to init.el:**
```elisp
(defun run-atheon ()
  "Run Atheon on current file"
  (interactive)
  (shell-command (concat "atheon " (buffer-file-name))))

(global-set-key (kbd "C-c a") 'run-atheon)
```

### Git Hooks

#### Pre-commit Hook

**Installation:**
```bash
# Copy hook to git hooks directory
cp hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

**Manual Setup:**
```bash
# .git/hooks/pre-commit
#!/bin/sh
atheon ./
if [ $? -ne 0 ]; then
  echo "❌ Security findings detected! Please fix before committing."
  exit 1
fi
```

#### Pre-push Hook

```bash
# .git/hooks/pre-push
#!/bin/sh
echo "Running security scan before push..."
atheon --categories=secrets ./
if [ $? -ne 0 ]; then
  echo "❌ Secrets detected! Push blocked."
  exit 1
fi
```

#### Post-merge Hook

```bash
# .git/hooks/post-merge
#!/bin/sh
echo "Scanning for security issues after merge..."
atheon --json ./ > post-merge-scan.json
if [ -s post-merge-scan.json ]; then
  echo "⚠️ Security findings detected in merge!"
fi
```

### Hook Manager Integration

#### Husky (Node.js)

```json
// package.json
{
  "husky": {
    "hooks": {
      "pre-commit": "atheon ./",
      "pre-push": "atheon --categories=secrets ./"
    }
  }
}
```

#### Lefthook

```yaml
# .lefthook.yml
pre-commit:
  commands:
    atheon-secrets:
      run: atheon --categories=secrets {files}
    atheon-pii:
      run: atheon --categories=pii {files}

pre-push:
  commands:
    atheon-full:
      run: atheon ./
```

#### pre-commit (Python)

```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: atheon
        name: Atheon Pattern Scanner
        entry: atheon
        language: system
        pass_filenames: true
        args: [--categories=secrets]
```

## CI/CD Integration

### GitHub Actions

#### Basic Security Scan

```yaml
# .github/workflows/security-scan.yml
name: Security Scan
on: [push, pull_request]

jobs:
  atheon:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Atheon
        run: |
          curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-amd64 -o atheon
          chmod +x atheon

      - name: Run Scan
        run: ./atheon --json ./ > findings.json

      - name: Upload Results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: security-findings
          path: findings.json
```

#### With Failure Enforcement

```yaml
name: Security Gate
on: [pull_request]

jobs:
  security-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Atheon
        run: |
          brew tap HoraDomu/atheon
          brew install atheon

      - name: Security Scan
        run: |
          atheon --json ./ > findings.json
          if [ -s findings.json ]; then
            echo "❌ Security findings detected!"
            cat findings.json
            exit 1
          fi
```

#### PR Comments

```yaml
name: Security PR Check
on: [pull_request]

jobs:
  security-comment:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Scan and Comment
        uses: actions/github-script@v7
        with:
          script: |
            const { execSync } = require('child_process');
            try {
              execSync('atheon ./');
              console.log('✅ No security findings');
            } catch (error) {
              const findings = error.stdout.toString();
              // Comment on PR with findings
              await github.rest.issues.createComment({
                issue_number: context.issue.number,
                owner: context.repo.owner,
                repo: context.repo.repo,
                body: `⚠️ Security findings detected:\n\`\`\`\n${findings}\n\`\`\``
              });
              process.exit(1);
            }
```

### GitLab CI

```yaml
# .gitlab-ci.yml
security_scan:
  stage: test
  script:
    - curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-amd64 -o atheon
    - chmod +x atheon
    - ./atheon --json ./ > security-report.json
  artifacts:
    reports:
      security: security-report.json
    paths:
      - security-report.json
  only:
    - merge_requests
    - main
```

### Jenkins Pipeline

```groovy
// Jenkinsfile
pipeline {
    agent any

    stages {
        stage('Security Scan') {
            steps {
                script {
                    sh '''
                        curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-amd64 -o atheon
                        chmod +x atheon
                        ./atheon --json ./ > findings.json
                    '''

                    def findings = readJSON file: 'findings.json'
                    if (findings.size() > 0) {
                        error("Security findings detected: ${findings.size()} issues found")
                    }
                }
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'findings.json', fingerprint: true
        }
    }
}
```

### Azure DevOps

```yaml
# azure-pipelines.yml
trigger:
- main

pr:
- main

jobs:
- job: security_scan
  pool:
    vmImage: 'ubuntu-latest'

  steps:
  - script: |
      curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-amd64 -o atheon
      chmod +x atheon
      ./atheon --json ./ > $(Build.ArtifactStagingDirectory)/findings.json
    displayName: 'Run Atheon Scan'

  - task: PublishBuildArtifacts@1
    inputs:
      pathtoPublish: '$(Build.ArtifactStagingDirectory)/findings.json'
      artifactName: 'security-findings'
```

## DevOps Integration

### Kubernetes

#### Pre-deployment Check

```yaml
# deployment-pipeline.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: atheon-scan-script
data:
  scan.sh: |
    #!/bin/sh
    curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-amd64 -o atheon
    chmod +x atheon
    ./atheon ./k8s/
---
apiVersion: batch/v1
kind: Job
metadata:
  name: security-scan
spec:
  template:
    spec:
      containers:
      - name: scanner
        command: ["/bin/sh", "/scripts/scan.sh"]
        volumeMounts:
        - name: scripts
          mountPath: /scripts
      volumes:
      - name: scripts
        configMap:
          name: atheon-scan-script
      restartPolicy: Never
```

### Docker Integration

#### Multi-stage Build

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
RUN apk --no-cache add curl
RUN curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-musl-arm64 -o atheon
RUN chmod +x atheon

FROM alpine:latest
COPY --from=builder /atheon /usr/local/bin/
RUN apk --no-cache add ca-certificates

# Scan during build
COPY . /app
RUN atheon /app || (echo "Security findings detected!" && exit 1)

CMD ["your-app"]
```

### Terraform

#### State File Scanning

```hcl
# terraform/scripts/security-scan.sh
resource "null_resource" "security_scan" {
  provisioner "local-exec" {
    command = <<EOT
      curl -sL https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-amd64 -o atheon
      chmod +x atheon
      ./atheon --categories=secrets .
    EOT
  }

  triggers = {
    always_run = true
  }
}
```

### Ansible

#### Playbook Integration

```yaml
# security-scan.yml
---
- name: Security Scan with Atheon
  hosts: localhost
  tasks:
    - name: Download Atheon
      get_url:
        url: https://github.com/HoraDomu/Atheon/releases/latest/download/atheim-linux-amd64
        dest: /tmp/atheim
        mode: '0755'

    - name: Run security scan
      command: /tmp/atheim --json ./ > results.json
      register: scan_results
      failed_when: false

    - name: Check for findings
      debug:
        msg: "Security findings detected!"
      when: scan_results.rc != 0
```

## Monitoring & Alerting

### Slack Integration

```bash
#!/bin/bash
# scan-and-alert.sh

findings=$(atheon --json ./)
if [ ! -z "$findings" ]; then
  curl -X POST -H 'Content-type: application/json' \
    --data "{\"text\":\":rotating_light: Atheon found security issues in $(hostname):\"}" \
    $SLACK_WEBHOOK_URL

  # Send findings
  curl -X POST -H 'Content-type: application/json' \
    --data "{\"text\":\`\`\`$findings\`\`\`\"}" \
    $SLACK_WEBHOOK_URL
fi
```

### Email Alerts

```bash
#!/bin/bash
# scan-and-email.sh

atheon --json ./ > findings.json
if [ -s findings.json ]; then
  mail -s "Security Findings Detected" security-team@example.com < findings.json
fi
```

### Prometheus Integration

```bash
#!/bin/bash
# scan-with-metrics.sh

START_TIME=$(date +%s)
atheon ./ > /dev/null 2>&1
EXIT_CODE=$?
END_TIME=$(date +%s)

# Push to Pushgateway
cat <<EOF | curl --data-binary @- http://pushgateway:9091/metrics/job/atheim-scan
# HELP atheon_scan_exit_code Exit code of Atheon scan
# TYPE atheon_scan_exit_code gauge
atheon_scan_exit_code{host="$(hostname)"} $EXIT_CODE
# HELP atheon_scan_duration_seconds Duration of Atheon scan in seconds
# TYPE atheon_scan_duration_seconds gauge
atheon_scan_duration_seconds{host="$(hostname)"} $(($END_TIME - $START_TIME))
EOF
```

## AI/LLM Integration

### MCP Server

#### Claude Desktop Integration

```json
// claude_desktop_config.json
{
  "mcpServers": {
    "atheon": {
      "command": "/usr/local/bin/atheim-mcp"
    }
  }
}
```

#### Cline (VS Code)

```json
// settings.json
{
  "cline.mcpServers": {
    "atheon": {
      "command": "atheon-mcp"
    }
  }
}
```

### Custom LLM Integration

```python
# llm_integration.py
import subprocess
import json

def scan_with_atheon(content):
    """Scan content using Atheon MCP server"""
    # Call MCP server or CLI
    result = subprocess.run(
        ['atheon', '--json', '-'],
        input=content,
        capture_output=True,
        text=True
    )

    if result.returncode == 0:
        return json.loads(result.stdout)
    return []

# Use in LLM context
def enhance_llm_context(code_content):
    findings = scan_with_atheon(code_content)

    if findings:
        warning = f"⚠️ SECURITY WARNING: {len(findings)} potential issues found:\n"
        for finding in findings[:3]:  # Limit for context
            warning += f"- {finding['pattern']} at line {finding['line']}\n"

        return warning + "\n\n" + code_content
    return code_content
```

## Automation Workflows

### Cron Jobs

```bash
# crontab -e
# Daily security scan at 2 AM
0 2 * * * cd /path/to/project && /usr/local/bin/atheim --json ./ > /var/log/atheim-$(date +\%Y\%m\%d).json

# Weekly comprehensive scan
0 3 * * 0 cd /path/to/project && /usr/local/bin/atheim --all --json ./ > /var/log/atheim-weekly-$(date +\%Y\%m\%d).json
```

### Systemd Service

```ini
# /etc/systemd/system/atheim-scan.service
[Unit]
Description=Atheon Security Scanner
After=network.target

[Service]
Type=oneshot
WorkingDirectory=/path/to/project
ExecStart=/usr/local/bin/atheim --json ./ /var/log/atheim-scan.json
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```ini
# /etc/systemd/system/atheim-scan.timer
[Unit]
Description=Run Atheon scan daily
Requires=atheim-scan.service

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
```

### File Watcher Integration

```bash
# Using entr for automatic scanning
find . -name "*.go" -o -name "*.js" -o -name "*.py" | \
  entr -r atheon --categories=secrets ./
```

## Additional Resources

- [Common Use Cases](../getting-started/common-use-cases.md)
- [Pattern Development](../patterns/README.md)
- [Architecture Documentation](../architecture/README.md)