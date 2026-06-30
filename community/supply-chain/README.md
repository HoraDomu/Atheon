# Supply Chain Security Patterns

Detects package management vulnerabilities, dependency confusion attacks, typosquatting, and other supply chain risks.

## Patterns

- `npm-internal-package-mismatch`: npm packages referencing internal registries without explicit scope
- `pypi-internal-package-mismatch`: pip requirements referencing internal PyPI servers
- `pip-extra-index-url`: pip installs from untrusted extra index URLs
- `npm-install-script`: npm packages with potentially malicious install scripts
- `typosquat-common-packages`: typosquatting targets (react, vue, lodash, etc.)

## References

- [CWE-1395: Dependency/Supply Chain](https://cwe.mitre.org/data/definitions/1395.html)
- [OWASP Dependency Confusion](https://owasp.org/www-project-top-ten/2017/A8)
