# Pattern Expansion: Issue #149

## Context
Based on real-world repo benchmarking via Atheon-GitHub-Scanner and Atheon-Benchmark projects.

## High-Priority Pattern Categories to Add

### 1. Infrastructure as Code
- Kubernetes secrets in env vars
- Terraform state files with sensitive data
- Docker config with credentials
- Cloudformation hidden resources

### 2. CI/CD Specific
- GitHub Actions secrets exposure
- GitLab CI variable usage
- Jenkins credentials
- CircleCI context leaks

### 3. API Framework Patterns
- Express.js error handling leaks
- Django DEBUG mode indicators
- Flask SECRET_KEY patterns
- Spring Boot actuator exposure

### 4. Database Patterns
- Redis connection strings with auth
- MongoDB connection strings
- PostgreSQL connection strings
- Database URL patterns in config

### 5. Authentication/Authorization
- JWT token patterns (beyond generic)
- OAuth client secrets
- API key header patterns
- Session token exposures

### 6. Cryptography
- RSA private key blocks
- ECDSA private keys
- Ed25519 keys
- WireGuard config keys

## From Benchmark Findings

The scanner found 2409 security issues in 2000 repositories scanned. Key findings:
- High occurrence of hardcoded credentials in config files
- Environment variable exposure patterns
- API key in URL query parameters
- Bearer token in Authorization headers

## Implementation Plan

1. Create new category folders in `community/`
2. Add YAML patterns for each finding type
3. Validate against existing patterns to avoid duplicates
4. Test with known-good and known-bad samples

## Suggested New Categories
- `community/iac/` - Infrastructure as Code patterns
- `community/ci-cd/` - CI/CD specific patterns
- `community/api-frameworks/` - Web framework specific patterns

---

**Date:** 2026-06-22  
**Branch:** `pr/149-patterns-expansion`
**Source:** Atheon-GitHub-Scanner mass_scan_summary.json (scanned 2000 repos, found 2409 issues)
