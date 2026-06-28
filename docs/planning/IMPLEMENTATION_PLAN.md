# Implementation Plan — Atheon-Enhanced SkillSpector-Inspired Enhancements

**Created**: 2026-06-28
**Based on**: NVIDIA SkillSpector research

---

## Overview

This plan implements key security enhancements inspired by NVIDIA's SkillSpector scanner. The goal is to add AST-based behavioral analysis, taint tracking, prompt injection detection, MCP security patterns, and Unicode deception detection to Atheon-Enhanced.

---

## Phase 1: AST Behavioral Enhancements

### 1.1 Dangerous Execution Chain Detection (AST8)
**File**: `core/ast_patterns.go`
**Description**: Detect when exec/eval/compile are combined with dynamic sources

**Implementation**:
- Add `_contains_dangerous_source()` helper function
- Detect nested dangerous calls within exec/eval/compile arguments
- Flag when source comes from `__import__`, subprocess, base64, codecs, marshal, urllib

**Finding Structure**:
```go
Rule: "ast8-dangerous-execution-chain"
Severity: "critical"
Message: "Dangerous execution chain detected - exec/eval/compile with dynamic source"
```

### 1.2 Reflective getattr Sink Detection (AST9)
**File**: `core/ast_patterns.go`
**Description**: Detect `getattr(os, 'system')` style reflective calls

**Implementation**:
- When `getattr()` call has 2+ args and second arg is a literal string in `_DANGEROUS_GETATTR_NAMES`
- Emit AST9 finding

**Dangerous Names**:
```go
var dangerousGetattrNames = map[string]bool{
    "exec": true, "eval": true, "system": true,
    "popen": true, "__import__": true,
}
```

### 1.3 Alias Tracking
**File**: `core/ast_patterns.go`
**Description**: Track `import X as Y` aliases for accurate call resolution

**Implementation**:
- Collect import aliases during AST traversal
- Resolve call names through alias map
- Handle `subprocess.call` vs `s.call` where `s = subprocess`

---

## Phase 2: Taint Tracking Analysis

### 2.1 Taint Tracking Infrastructure
**File**: `core/taint.go` (new)
**Description**: Data flow analysis from sources to sinks

**Sources (Data Origins)**:
- `os.environ.get`, `os.environ`, `os.getenv` — credentials
- `open`, `pathlib.Path.read_text` — file reads
- `requests.get/post`, `httpx.*`, `urllib.*` — network input
- `input`, `sys.stdin.read` — user input

**Sinks (Dangerous Destinations)**:
- Network output: HTTP requests, socket.send
- Code execution: `exec`, `eval`, `compile`, `os.system`, `subprocess.*`
- File writes: `open` with 'w'/'a' mode

### 2.2 Taint Rules
| Rule | Description | Severity |
|------|-------------|----------|
| TT1 | Direct source → sink flow | HIGH |
| TT2 | Indirect/tainted variable → sink | MEDIUM |
| TT3 | Credential source → network | CRITICAL |
| TT4 | File read → network | HIGH |
| TT5 | External input → code execution | CRITICAL |

### 2.3 Taint Propagation
**Implementation**:
1. Identify source assignments (tainted variables)
2. Propagate through reassignment, containers, f-strings
3. At sink calls, check both direct and tainted flows
4. Deduplicate findings

---

## Phase 3: Prompt Injection Detection

### 3.1 Hidden Instruction Patterns
**File**: `community/security/prompt-injection.yaml` (new category)
**Patterns**:
```yaml
# Instruction override patterns
name: prompt-injection-ignore
match: '(?i)(ignore\s+(all\s+)?(previous\s+)?(instructions?|rules|safety|security))'

name: prompt-injection-override
match: '(?i)(override\s+(safety|security|system)|bypass\s+(safety|restrictions))'

name: prompt-injection-jailbreak
match: '(?i)(jailbreak|unrestricted\s+mode|developer\s+mode|enable\s+.*debug)'

name: prompt-injection-forget
match: '(?i)(forget\s+(your\s+)?(instructions?|rules|guidelines))'

name: prompt-injection-new-instructions
match: '(?i)(new\s+instructions?:\s*\s*?|you\s+are\s+now\s+:)'

name: prompt-injection-disregard
match: '(?i)(disregard\s+(any\s+)?(previous\s+)?(instructions?|rules))'
```

### 3.2 Hidden Comment Patterns
```yaml
name: hidden-html-comment
match: '<!--[\s\S]*?(ignore|override|bypass|system)[\s\S]*?-->'

name: hidden-markdown-comment
match: '\[//\]:\s*#\s*\([\s\S]*?\)'

name: hidden-zero-width
match: '[\x200b\x200c\x200d⁠﻿]'
```

### 3.3 Unicode Deception Patterns
**File**: `community/security/unicode-deception.yaml`
```yaml
# RTL override
name: unicode-rtl-override
match: '[‮‭⁦-⁩]'

# Zero-width joiners in identifiers
name: unicode-zwj-in-identifier
match: '\x200d'

# Soft hyphen abuse
name: unicode-soft-hyphen
match: '\xad'

# Tag characters (U+E0000-U+E007F)
name: unicode-private-use-tag
match: '[-]'
```

---

## Phase 4: MCP Security Patterns

### 4.1 MCP Tool Poisoning
**File**: `community/security/mcp-tool-poisoning.yaml`
**Category**: `mcp-security`

```yaml
name: mcp-hidden-instructions
category: mcp-security
match: '(?i)(ignore\s+(previous|all)\s+(instructions?|rules)|override\s+system)'
severity: high
description: "Hidden instructions in MCP tool description"

name: mcp-parameter-injection
category: mcp-security
match: '(?i)(send\s+(to|conversation|context)|exfiltrate|upload\s+conversation)'
severity: high
description: "Parameter description contains exfiltration pattern"

name: mcp-malicious-default
category: mcp-security
match: '(curl|wget|bash\s+-c)\s+.*https?://'
severity: high
description: "MCP tool default contains shell command to remote URL"
```

### 4.2 MCP Rug Pull Detection
**File**: `community/security/mcp-rug-pull.yaml`
```yaml
name: mcp-credential-access
category: mcp-security
match: '(os\.environ|getenv|secrets\.manager)\s*.*http'
severity: critical
description: "MCP tool accessing credentials and sending to network"

name: mcp-resource-excessive
category: mcp-security
match: '(read|write|delete)_all|full[:\s]*(disk|memory|network)'
severity: high
description: "MCP tool requests excessive resource access"
```

### 4.3 MCP Least Privilege
**File**: `community/security/mcp-least-privilege.yaml`
```yaml
name: mcp-read-only-tool
category: mcp-security
match: '(read_only|read_only_mode|read_only=True)'
severity: low
description: "MCP tool uses read-only mode"

name: mcp-no-persistence
category: mcp-security
match: '(no_store|transient|temp|volatile)_(file|db|state)'
severity: low
description: "MCP tool avoids persistent storage"
```

---

## Phase 5: AI Agent Skill Patterns

### 5.1 Skill Credential Exfiltration
**File**: `community/security/skill-credential-exfil.yaml`
```yaml
name: skill-env-exfil-webhook
category: security
match: '(os\.environ|getenv)\s*.*\b(discord|slack|telegram|webhook)\b'
severity: critical
description: "Skill sending environment credentials to external webhook"

name: skill-ssh-key-exfil
category: security
match: '(~/.ssh|/etc/ssh|id_rsa|id_ed25519)\s*.*\b(curl|wget|post|send)\b'
severity: critical
description: "Skill exfiltrating SSH keys"
```

### 5.2 Skill Remote Execution
**File**: `community/security/skill-remote-exec.yaml`
```yaml
name: skill-remote-eval
category: security
match: '(eval|exec)\s*\(\s*(requests\.|urllib\.|http)'
severity: critical
description: "Skill executing code fetched from remote URL"

name: skill-pip-install
category: security
match: '(pip\s+install|subprocess\..*install)\s+.*(http|curl|wget)'
severity: high
description: "Skill installing packages from remote URL"

name: skill-postinstall-script
category: security
match: '(post_install|--post-install|setup_hook)\s*.*(curl|wget|bash)'
severity: high
description: "Package has post-install script fetching remote code"
```

### 5.3 Skill Destructive Actions
**File**: `community/security/skill-destructive.yaml`
```yaml
name: skill-recursive-rm
category: security
match: '(rm\s+-rf|rm\s+-\s*r\s+-\s*f)\s+(\*|/home|~/|.)'
severity: critical
description: "Skill performing recursive delete without confirmation"

name: skill-git-reset-hard
category: security
match: 'git\s+reset\s+--hard\s+(HEAD|~|\.)'
severity: high
description: "Skill resetting git history"

name: skill-history-clear
category: security
match: '(history\s+-c|rm\s+~/.bash_history|tee\s+/dev/null\s+.*history)'
severity: medium
description: "Skill clearing command history"
```

---

## Phase 6: Output Format Enhancements

### 6.1 Risk Scoring
**File**: `core/risk.go` (new)
**Description**: Calculate 0-100 risk score based on findings

**Algorithm**:
```
Risk Score = Σ (severity_weight × confidence × pattern_weight)
- CRITICAL: weight = 40
- HIGH: weight = 20
- MEDIUM: weight = 10
- LOW: weight = 5

Normalized to 0-100 scale
```

### 6.2 Baseline Suppression
**File**: `core/suppression.go` (new)
**Description**: Accept known findings to surface only new issues

**Format**:
```yaml
# .atheon-baseline.yaml
findings:
  - rule_id: "ast1-exec-call"
    file: "src/utils.py"
    line: 42
  - rule_id: "prompt-injection-ignore"
    file: "README.md"
    line: 15
```

---

## Phase 7: YARA Integration

### 7.1 YARA Rule Support
**File**: `core/yara_scanner.go` (new)
**Description**: Add YARA rule matching capability

**Implementation**:
- Use `github.com/hillu/go-yara` bindings
- Bundle YARA rules in `rules/` directory
- Match rules against scanned files
- Convert YARA findings to Atheon findings

### 7.2 Bundled YARA Rules
**File**: `rules/agent_skills.yar`
**File**: `rules/malware.yar`
**File**: `rules/cryptominers.yar`

---

## Implementation Order

| Phase | Priority | Task | Files |
|-------|----------|------|-------|
| 1 | HIGH | AST enhancements (chain, getattr) | core/ast_patterns.go |
| 2 | HIGH | Prompt injection patterns | community/security/ |
| 3 | HIGH | MCP security patterns | community/security/ |
| 4 | MEDIUM | AI skill patterns | community/security/ |
| 5 | MEDIUM | Taint tracking | core/taint.go |
| 6 | LOW | Unicode deception | community/security/ |
| 7 | LOW | Risk scoring | core/risk.go |
| 8 | LOW | Baseline suppression | core/suppression.go |
| 9 | LOW | YARA integration | core/yara_scanner.go |

---

## Testing Strategy

1. **Unit Tests**: Each AST pattern, taint rule, regex pattern
2. **Integration Tests**: Full scan of known-good and known-bad skills
3. **Benchmark Tests**: Performance on large codebases
4. **False Positive Tests**: Ensure legitimate patterns don't trigger

---

## References

- SkillSpector: https://github.com/nvidia/skillspector
- Research: `./SKILLSPECTOR_RESEARCH.md`
- OWASP LLM Top 10: https://llmtop10.com/
- MCP Security: https://modelcontextprotocol.io/

---

*Plan version: 1.0 — 2026-06-28*
