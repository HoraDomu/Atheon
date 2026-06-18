# Community Guide

Join the Atheon community and contribute to building a comprehensive pattern matching engine for everyone.

## Getting Involved

### Contribution Opportunities

**Pattern Contributions:**
- Add missing patterns for services you use
- Improve existing patterns for better accuracy
- Create patterns for new domains (finance, healthcare, legal)

**Code Contributions:**
- Bug fixes and performance improvements
- Documentation improvements
- Test coverage expansion

**Community Support:**
- Help answer questions in issues
- Review pull requests
- Share use cases and integrations

## Communication Channels

### GitHub Issues
- **Bug Reports:** [Create an issue](https://github.com/HoraDomu/Atheon/issues/new?template=bug_report.md)
- **Feature Requests:** [Suggest improvements](https://github.com/HoraDomu/Atheon/issues/new?template=feature_request.md)
- **Pattern Requests:** [Request new patterns](https://github.com/HoraDomu/Atheon/issues/new?template=pattern_request.md)

### Email
- **General Inquiries:** [dommcpro@gmail.com](mailto:dommcpro@gmail.com)
- **Security Issues:** [Report security vulnerabilities](mailto:dommcpro@gmail.com)

## Contribution Guidelines

### Pattern Contributions

1. **Check Existing Patterns:**
   ```bash
   atheon list
   atheon list categories
   ```

2. **Create Pattern File:**
   ```bash
   # community/secrets/my-service.yaml
   name: my-service-api-key
   match: '\bmsvc_[A-Za-z0-9]{32}\b'
   ```

3. **Test Pattern:**
   ```bash
   go run ./bundler
   go test ./...
   ./atheon list | grep my-service-api-key
   ```

4. **Submit Pull Request:**
   - Include pattern description
   - Explain what it detects
   - Provide test cases
   - Show real-world examples

### Code Contributions

**Before Contributing:**
1. Check existing issues for similar work
2. Create an issue to discuss major changes
3. Follow code style guidelines
4. Add tests for new functionality
5. Update documentation

**Pull Request Process:**
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests and documentation
5. Submit pull request with clear description

## Pattern Categories

### Current Categories

**Security/Secrets:**
- API keys and tokens
- Authentication credentials
- Service-specific secrets
- Certificate keys

**PII:**
- Personal identifiers
- Financial information
- Healthcare data
- Contact information

### Requested Categories

**Finance (Requested):**
- Bank account numbers
- Credit card formats
- Routing numbers
- Financial codes

**Healthcare (Requested):**
- Patient IDs
- Medical record numbers
- Insurance numbers
- Health data codes

**Legal (Requested):**
- Contract references
- Case numbers
- Legal document identifiers
- Regulation codes

## Contribution Recognition

### Contributors

All contributors are recognized in:
- [CONTRIBUTORS.md](../../CONTRIBUTORS.md) file
- GitHub release notes
- Project documentation

### Contribution Types

**Pattern Contributors:**
- Pattern authors credited in pattern files
- Listed in contributor documentation
- Included in release notes

**Code Contributors:**
- GitHub attribution for commits
- Listed in contributors documentation
- Recognized in release notes

## Community Guidelines

### Code of Conduct

**Be Respectful:**
- Treat all community members with respect
- Welcome newcomers and help them learn
- Focus on constructive feedback

**Be Collaborative:**
- Work together on improvements
- Share knowledge and experience
- Help review contributions

**Be Professional:**
- Keep discussions focused and productive
- Acknowledge diverse perspectives
- Make decisions based on technical merit

### Communication Style

**Issue Reporting:**
- Provide clear reproduction steps
- Include environment information
- Attach relevant logs/output
- Be patient with volunteer responses

**Pull Requests:**
- Write clear descriptions
- Explain the "why" not just the "what"
- Include testing instructions
- Respond to review feedback

## Pattern Contribution Standards

### Quality Criteria

**Pattern Accuracy:**
- ✅ Detects intended format correctly
- ✅ Minimal false positives
- ✅ Minimal false negatives
- ✅ Handles edge cases properly

**Pattern Specificity:**
- ✅ Clear, descriptive naming
- ✅ Service-specific when applicable
- ✅ Proper use of word boundaries
- ✅ Context-aware matching

**Documentation:**
- ✅ Clear pattern purpose
- ✅ Service/format description
- ✅ Test cases included
- ✅ Real-world examples provided

### Review Process

**Initial Review:**
- Automatic testing (must pass)
- Pattern validation (regex compilation)
- Test coverage verification

**Community Review:**
- Pattern accuracy assessment
- False positive/negative evaluation
- Naming and documentation review
- Integration testing

**Maintainer Review:**
- Final approval decision
- Merge into main branch
- Include in next release
- Update documentation

## Release Process

### Release Schedule

**Automated Releases:**
- **10th and 21st of every month**
- Fully automated via GitHub Actions
- Include all merged contributions

**Release Contents:**
- New patterns from community
- Bug fixes and improvements
- Updated documentation
- Multi-platform binaries

### Version Strategy

**Semantic Versioning:**
- **Major:** Breaking changes
- **Minor:** New features/patterns
- **Patch:** Bug fixes

**Backward Compatibility:**
- Pattern interface stable
- CLI behavior consistent
- Bundle format maintained

## Testing and Validation

### Pattern Testing

**Automated Testing:**
- All patterns must have test cases
- Tests must pass before merge
- Coverage requirements enforced

**Manual Testing:**
- Real-world sample verification
- False positive checking
- Performance validation

### Community Testing

**Beta Testing:**
- Pre-release testing opportunities
- Community feedback collection
- Issue triage and resolution

## Documentation Contributions

### Documentation Types

**User Documentation:**
- Getting started guides
- Common use cases
- Integration examples
- Troubleshooting guides

**Developer Documentation:**
- Architecture explanations
- Code style guidelines
- Testing procedures
- Contribution processes

**Pattern Documentation:**
- Pattern writing guides
- Test case examples
- Best practices
- Common issues

## Recognition and Appreciation

### Contributor Spotlights

**Monthly Highlights:**
- Top pattern contributors
- Notable code contributions
- Community support recognitions

**Project Acknowledgments:**
- README acknowledgments
- Release notes credits
- Documentation references

### Community Impact

**Pattern Contributions:**
- Every pattern benefits all users
- Community-driven growth
- Domain expertise sharing
- Open source collaboration

## Support and Resources

### Getting Help

**Documentation:**
- [Getting Started Guide](../getting-started/README.md)
- [Development Guide](../development/README.md)
- [Pattern Development](../patterns/README.md)
- [Integration Guide](../integration/README.md)

**Community Support:**
- GitHub Issues for questions
- Email for direct contact
- Pattern sharing and discussion

### Learning Resources

**Pattern Development:**
- [Pattern Writing Guide](../patterns/README.md)
- Test case examples
- Best practices documentation
- Community examples

**Integration:**
- CI/CD examples
- Hook configurations
- Automation workflows
- API documentation

## Project Governance

### Maintainer Responsibilities

**Technical Leadership:**
- Code review and quality assurance
- Release management
- Architecture decisions
- Security oversight

**Community Management:**
- Issue triage and response
- Pull request review
- Community guidelines enforcement
- Contributor recognition

### Decision Making

**Technical Decisions:**
- Consensus-based discussion
- Maintainer final authority
- Technical merit focus
- Community input valued

**Community Decisions:**
- Open discussion encouraged
- Contributor input valued
- Transparent decision process
- Documentation of rationale

## Future Directions

### Community Priorities

**Pattern Expansion:**
- New domain coverage
- Improved pattern accuracy
- Better false positive handling
- Performance optimization

**Feature Development:**
- Enhanced automation
- Better integration support
- Improved user experience
- Advanced testing capabilities

### Strategic Goals

**Community Growth:**
- Expand contributor base
- Improve onboarding process
- Enhance documentation
- Strengthen community support

**Technical Excellence:**
- Maintain code quality
- Improve test coverage
- Enhance security
- Optimize performance

## Getting Started

### First Steps

1. **Explore the Project:**
   - Read the [README](../../README.md)
   - Review [existing patterns](../../community/)
   - Check [documentation](../)

2. **Choose Your Contribution:**
   - Add a missing pattern
   - Improve documentation
   - Report an issue
   - Review a pull request

3. **Join the Community:**
   - Star the repository
   - Watch for updates
   - Contribute patterns
   - Share your experience

## Additional Resources

- [Contributing Guidelines](../../CONTRIBUTING.md)
- [Development Guide](../development/README.md)
- [Pattern Development](../patterns/README.md)
- [GitHub Repository](https://github.com/HoraDomu/Atheon)

## Contact

**Questions or suggestions?**
- Email: [dommcpro@gmail.com](mailto:dommcpro@gmail.com)
- GitHub: [https://github.com/HoraDomu/Atheon](https://github.com/HoraDomu/Atheon)

**Welcome to Atheon!**