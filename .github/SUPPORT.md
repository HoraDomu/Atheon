# Support

Get help with Atheon-Enhanced.

## Getting Help

### GitHub Issues

For bugs, feature requests, and questions, please [open an issue](https://github.com/aliasfoxkde/Atheon-Enhanced/issues/new/choose).

### Discussions

For questions and community discussions, visit the [GitHub Discussions](https://github.com/aliasfoxkde/Atheon-Enhanced/discussions) page.

## FAQ

### General

**Q: What is Atheon-Enhanced?**
A: Atheon-Enhanced is a feature-rich pattern matching engine for secrets detection, AI-generated code identification, and quality enforcement. It's a fork of [HoraDomu/Atheon](https://github.com/HoraDomu/Atheon) with additional patterns and features.

**Q: What's the difference between Atheon and Atheon-Enhanced?**
A: The enhanced version includes 327+ patterns across 23 categories (vs 57 in upstream), streaming API for memory-efficient scanning, MCP server integration, AI-generated code detection, and more. See [docs/reports/FEATURE_COMPARISON.md](docs/reports/FEATURE_COMPARISON.md) for details.

**Q: Is Atheon-Enhanced production-ready?**
A: This is an experimental testing fork. For production use, consider the [official HoraDomu/Atheon](https://github.com/HoraDomu/Atheon) project which has more conservative, battle-tested patterns.

### Installation

**Q: How do I install Atheon-Enhanced?**
```bash
go install github.com/aliasfoxkde/Atheon@latest
```

Or build from source:
```bash
git clone https://github.com/aliasfoxkde/Atheon-Enhanced.git
cd Atheon-Enhanced
go build -o atheon ./cmd/atheon
```

**Q: What Go version is required?**
A: Go 1.22 or later is recommended.

### Usage

**Q: How do I scan a project?**
```bash
atheon ./my-project
```

**Q: How do I use MCP integration?**
See [docs/integrations/mcp.md](docs/integrations/mcp.md) for setup instructions.

**Q: Can I use Atheon as a library?**
Yes. See the library usage examples in [README.md](README.md#library-usage).

### Troubleshooting

**Q: Patterns aren't being detected**
- Check that the category is enabled: `atheon list --category <name>`
- Verify the pattern is enabled: `atheon list --enabled`
- Try with `--all` flag to include disabled patterns

**Q: MCP server won't start**
- Ensure port 8080 is available
- Check logs for error messages
- Verify configuration in `config/profiles/mcp-integration.json`

**Q: High memory usage**
- Use streaming mode: `--streaming`
- Scan directories individually instead of large trees
- Adjust chunk size in configuration

## Community Resources

- [GitHub Discussions](https://github.com/aliasfoxkde/Atheon-Enhanced/discussions) - Community Q&A
- [GitHub Wiki](https://github.com/aliasfoxkde/Atheon-Enhanced/wiki) - User documentation
- [Project Roadmap](docs/ROADMAP.md) - Upcoming features

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to contribute to the project.

## Security

For security vulnerabilities, please see [SECURITY.md](SECURITY.md).

## Support Channels

| Channel | Purpose |
|---------|---------|
| [GitHub Issues](https://github.com/aliasfoxkde/Atheon-Enhanced/issues) | Bug reports and feature requests |
| [GitHub Discussions](https://github.com/aliasfoxkde/Atheon-Enhanced/discussions) | Questions and community help |
| [Security Advisories](https://github.com/aliasfoxkde/Atheon-Enhanced/security/advisories) | Private vulnerability reporting |
