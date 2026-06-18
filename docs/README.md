# Atheon Documentation

Comprehensive documentation for the Atheon pattern matching engine.

## Getting Started

- **[Installation Guide](getting-started/README.md)** - Installation and basic usage
- **[Common Use Cases](getting-started/common-use-cases.md)** - Real-world usage examples

## Core Documentation

### Pattern Development
- **[Pattern Guide](patterns/README.md)** - Creating and contributing patterns
- **[Pattern Reference](reference/pattern-reference.md)** - Built-in patterns reference

### Development
- **[Development Guide](development/README.md)** - Contributing to the codebase
- **[Architecture](architecture/README.md)** - System architecture and design

### Integration
- **[Integration Guide](integration/README.md)** - CI/CD, hooks, and automation

### Reference
- **[CLI Reference](reference/cli-reference.md)** - Command-line interface reference
- **[MCP Reference](reference/mcp-reference.md)** - MCP server documentation

### Community
- **[Community Guide](community/README.md)** - Contributing and community involvement

## Quick Links

### For Users
- [Installation Guide](getting-started/README.md)
- [Common Use Cases](getting-started/common-use-cases.md)
- [CLI Reference](reference/cli-reference.md)

### For Contributors
- [Pattern Development](patterns/README.md)
- [Development Guide](development/README.md)
- [Community Guide](community/README.md)

### For Integrators
- [Integration Guide](integration/README.md)
- [MCP Reference](reference/mcp-reference.md)
- [Architecture](architecture/README.md)

## Documentation Structure

```
docs/
├── getting-started/     # User documentation
│   ├── README.md        # Installation and basic usage
│   └── common-use-cases.md
├── patterns/            # Pattern development
│   └── README.md        # Pattern creation guide
├── development/         # Development documentation
│   └── README.md        # Contributing to codebase
├── architecture/       # Technical documentation
│   └── README.md        # System architecture
├── integration/         # Integration guides
│   └── README.md        # CI/CD and automation
├── community/           # Community resources
│   └── README.md        # Contributing guidelines
├── reference/           # Reference documentation
│   ├── cli-reference.md
│   ├── mcp-reference.md
│   └── pattern-reference.md
└── README.md            # This file
```

## Key Concepts

### Pattern Matching Engine

Atheon is a **community-driven pattern matching engine** for detecting:
- Security issues (API keys, tokens, credentials)
- PII (personal identifiers, sensitive data)
- Compliance issues (regulated data, financial information)
- Custom domain patterns

### Core Features

- **YAML Patterns:** No code required for pattern development
- **Category Filtering:** Scan only what you need
- **Multiple Input Sources:** Files, directories, stdin, environment
- **MCP Server:** AI/LLM integration capability
- **Cross-platform:** Windows, macOS, Linux support

### Development Philosophy

- **Minimal Dependencies:** Only 1 external dependency
- **Simple Architecture:** Clean Go code with clear patterns
- **Community-Driven:** Patterns contributed by users
- **Automation-Ready:** CI/CD friendly exit codes and JSON output

## Getting Help

### Documentation
- Start with [Getting Started Guide](getting-started/README.md)
- Check [Common Use Cases](getting-started/common-use-cases.md)
- Review [FAQ](#frequently-asked-questions)

### Community
- [GitHub Issues](https://github.com/HoraDomu/Atheon/issues)
- [Community Guide](community/README.md)
- Email: [dommcpro@gmail.com](mailto:dommcpro@gmail.com)

### Reference
- [CLI Reference](reference/cli-reference.md)
- [Pattern Reference](reference/pattern-reference.md)
- [MCP Reference](reference/mcp-reference.md)

## Frequently Asked Questions

### What is Atheon?

Atheon is a pattern matching engine for detecting secrets, PII, and other sensitive data in code and text. It's community-driven, meaning anyone can contribute patterns.

### How do I install Atheon?

```bash
# macOS/Linux
brew tap HoraDomu/atheon
brew install atheon

# Windows
scoop bucket add atheon https://github.com/HoraDomu/scoop-atheon
scoop install atheon
```

### How do I create patterns?

Patterns are simple YAML files:

```yaml
# community/secrets/my-service.yaml
name: my-service-api-key
match: '\bmsvc_[A-Za-z0-9]{32}\b'
```

See [Pattern Development Guide](patterns/README.md) for details.

### How do I integrate Atheon into CI/CD?

See [Integration Guide](integration/README.md) for examples with:
- GitHub Actions
- GitLab CI
- Jenkins
- Pre-commit hooks
- And more

### What's the difference between categories?

Categories organize patterns by domain:
- `secrets` - API keys, tokens, credentials
- `pii` - Personal identifiers, sensitive data

Use category filtering to scan only what you need:
```bash
atheon --categories=secrets ./
```

### How do I use the MCP server?

Configure in your AI tool settings:

```json
{
  "mcpServers": {
    "atheon": {
      "command": "atheon-mcp"
    }
  }
}
```

See [MCP Reference](reference/mcp-reference.md) for details.

### How often are patterns updated?

New pattern bundles are released on the **10th and 21st of every month**.

Update with:
```bash
atheon update
```

## Project Status

**Current Version:** See [GitHub Releases](https://github.com/HoraDomu/Atheon/releases/latest)

**Patterns:** 14 built-in patterns
**Categories:** 2 categories (secrets, pii)
**Languages:** Go 1.21+
**Platforms:** Windows, macOS, Linux

## Contributing

We welcome contributions! See:

- **[Pattern Development](patterns/README.md)** - Add new patterns
- **[Development Guide](development/README.md)** - Code contributions
- **[Community Guide](community/README.md)** - Getting involved

## License

MIT with Additional Terms Copyright © 2026 Dominick Yanez

See [LICENSE](../../LICENSE) for complete terms.

## Additional Resources

- **[Main Repository](https://github.com/HoraDomu/Atheon)** - GitHub repository
- **[Issue Tracker](https://github.com/HoraDomu/Atheon/issues)** - Bug reports and feature requests
- **[Contributing Guidelines](../../CONTRIBUTING.md)** - Contribution process
- **[CLAUDE.md](../../CLAUDE.md)** - AI assistant instructions

## Documentation Version

**Last Updated:** 2025-06-17
**Atheon Version:** Latest (see releases)
**Documentation Coverage:** Comprehensive

---

**Ready to get started?** Begin with the [Installation Guide](getting-started/README.md).