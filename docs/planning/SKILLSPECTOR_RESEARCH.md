# SkillSpector Research — NVIDIA AI Agent Skill Security Scanner

**Researched**: 2026-06-28
**Source**: https://github.com/nvidia/skillspector

---

## Executive Summary

SkillSpector is NVIDIA's open-source security scanner for AI agent skills (Claude Code, Codex CLI, Gemini CLI). It uses a two-stage pipeline: fast static analysis + optional LLM semantic evaluation. It detects 68 vulnerability patterns across 17 categories using regex, AST analysis, YARA rules, and OSV.dev CVE lookups.

---

## Key Architectural Insights

### Two-Stage Detection Pipeline
1. **Stage 1 — Static Analysis**: Regex patterns + AST behavioral analysis + OSV.dev CVE lookups
2. **Stage 2 — LLM Semantic Analysis**: Context-aware evaluation, false positive filtering (~87% precision)

### AST Behavioral Patterns (9 patterns)
| Rule | Name | Severity | Confidence |
|------|------|----------|------------|
| AST1 | exec() Call | HIGH | 0.85 |
| AST2 | eval() Call | HIGH | 0.85 |
| AST3 | Dynamic Import (__import__) | MEDIUM | 0.75 |
| AST4 | subprocess Call | MEDIUM | 0.70 |
| AST5 | os.system/exec-family | HIGH | 0.85 |
| AST6 | compile() Call | MEDIUM | 0.65 |
| AST7 | Dynamic getattr() | LOW | 0.50 |
| AST8 | Dangerous Execution Chain | CRITICAL | 0.95 |
| AST9 | Reflective getattr() Sink | HIGH | 0.85 |

### Taint Tracking Analysis (5 rules)
| Rule | Source → Sink | Severity | Confidence |
|------|---------------|----------|------------|
| TT1 | Direct flow | HIGH | 0.80 |
| TT2 | Indirect/tainted flow | MEDIUM | 0.65 |
| TT3 | Credential → Network | CRITICAL | 0.90 |
| TT4 | File read → Network | HIGH | 0.80 |
| TT5 | External input → Code execution | CRITICAL | 0.90 |

### Static Pattern Categories (17)
1. prompt_injection
2. harmful_content
3. data_exfiltration
4. privilege_escalation
5. excessive_agency
6. tool_misuse
7. rogue_agent
8. memory_poisoning
9. system_prompt_leakage
10. anti_refusal
11. agent_snooping
12. output_handling
13. supply_chain
14. ssrf
15. mcp_tool_poisoning
16. mcp_rug_pull
17. mcp_least_privilege

### YARA Rules (5 rules)
1. agent_skill_credential_exfiltration_webhook (CRITICAL)
2. agent_skill_remote_bootstrap_execution (HIGH)
3. agent_skill_prompt_injection_hidden_instructions (HIGH)
4. agent_skill_mcp_tool_poisoning_metadata (HIGH)
5. agent_skill_destructive_autonomous_actions (HIGH)

---

## Advanced Detection Techniques

### Unicode Deception Detection
- Homoglyph confusables (Cyrillic/Greek → Latin)
- RTL/directional overrides (U+202E, U+202D)
- Zero-width characters (U+200B, U+200C, U+200D)
- Invisible formatting chars (U+00AD, U+034F, U+2060)
- Tag characters (U+E0000–U+E007F)
- Mixed-script detection

### Hidden Instruction Detection
- HTML comments with instruction keywords
- Markdown comments `[//]: # (...)`
- Base64-encoded payloads
- Data URIs
- Zero-width joiners in identifiers

### MCP Tool Poisoning (TP1-TP4)
- Hidden instructions in tool descriptions
- Unicode deception
- Parameter description injection
- LLM-based behavior mismatch detection

---

## SkillSpector Finding Structure

```python
class Finding:
    rule_id: str           # e.g., "AST1", "TT3"
    message: str           # Human-readable message
    severity: Severity     # LOW, MEDIUM, HIGH, CRITICAL
    confidence: float      # 0.0-1.0
    file: str             # File path
    start_line: int
    end_line: int | None
    category: str          # Pattern category
    matched_text: str      # First 200 chars of match
    context: str | None   # Surrounding text
    remediation: str | None
    tags: list[str]
```

---

## Gaps in Atheon-Enhanced vs SkillSpector

### Missing AST Capabilities
1. **Dangerous execution chain detection** (exec+eval+compile combo)
2. **Reflective getattr() sink detection** (getattr(os,'system'))
3. **Taint tracking** (source → sink analysis)
4. **Alias tracking** in AST (handling `import subprocess as s`)

### Missing Pattern Categories
1. **MCP tool poisoning** patterns
2. **AI agent skill** specific patterns
3. **Prompt injection** in strings/comments
4. **Hidden instructions** (HTML/markdown comments)
5. **Unicode deception** detection

### Missing Output Formats
1. **SARIF model** (SkillSpector has sarif_models.py)
2. **Risk scoring** (0-100 with severity labels)
3. **Baseline suppression** for re-scans
4. **LLM semantic evaluation** (optional stage 2)

### Missing Integrations
1. **OSV.dev CVE lookups**
2. **YARA rule engine**
3. **MCP server** for runtime guardrails
4. **Multi-skill analysis** (batch scanning)

---

## Recommendations for Atheon-Enhanced

### High Priority (Implement Next)
1. Add dangerous execution chain detection (AST8)
2. Add reflective getattr sink detection (AST9)
3. Add prompt injection patterns
4. Add MCP tool poisoning patterns
5. Add Unicode deception detection

### Medium Priority
1. Add taint tracking analysis
2. Add risk scoring system
3. Add baseline suppression
4. Add YARA rule support

### Lower Priority
1. Add LLM semantic evaluation
2. Add OSV.dev CVE lookups
3. Add MCP server mode
4. Add multi-skill batch scanning

---

## References

- SkillSpector Repo: https://github.com/nvidia/skillspector
- AST Analyzer: `src/skillspector/nodes/analyzers/behavioral_ast.py`
- Taint Tracker: `src/skillspector/nodes/analyzers/behavioral_taint_tracking.py`
- MCP Poisoning: `src/skillspector/nodes/analyzers/mcp_tool_poisoning.py`
- Prompt Injection: `src/skillspector/nodes/analyzers/static_patterns_prompt_injection.py`
- YARA Rules: `src/skillspector/yara_rules/`

---

*This document informs the implementation plan in IMPLEMENTATION_PLAN.md*
