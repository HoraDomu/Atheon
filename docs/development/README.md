# Development Guide

This guide covers development workflows, testing, and contribution practices for the Atheon codebase.

## Project Structure

```
Atheon/
├── main.go              # CLI entry point
├── core/                # Core engine
│   ├── bundle.go        # Pattern loading and management
│   ├── finding.go       # Data structures
│   ├── pattern.go       # Pattern interface and registry
│   ├── runner.go        # Scanning logic
│   ├── patterns.bundle  # Embedded pattern bundle
│   └── bundle_test.go   # Pattern tests
├── cmd/
│   └── mcp/             # MCP server
│       └── main.go      # MCP implementation
├── bundler/             # Pattern bundling tool
│   └── main.go          # Bundle compiler
├── community/           # Pattern definitions
│   ├── secrets/         # Security patterns
│   └── pii/             # PII patterns
└── docs/                # Documentation
```

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git

### Setup

```bash
# Clone the repository
git clone https://github.com/HoraDomu/Atheon.git
cd Atheon

# Build the CLI
go build -o atheon .

# Build the MCP server
go build -o atheon-mcp ./cmd/mcp

# Run tests
go test ./...

# Run with race detector
go test -race ./...
```

## Development Workflow

### Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/my-change
   ```

2. **Make your changes**
   - Follow Go conventions and idioms
   - Add tests for new functionality
   - Update documentation as needed

3. **Test locally**
   ```bash
   go test ./...
   go test -race ./...
   go vet ./...
   ```

4. **Build and verify**
   ```bash
   go build -o atheon .
   ./atheon list
   ./atheon --file test/fixtures.txt
   ```

### Code Standards

**Go Conventions:**
- Exported names: `PascalCase`
- Unexported names: `camelCase`
- Acronyms: `URL` not `Url` for exported, `url` not `URL` for unexported
- No unnecessary comments unless the "why" is non-obvious

**Error Handling:**
- Always check errors
- Use `fmt.Errorf` with `%w` for error wrapping
- Provide context in error messages

**Code Quality:**
```bash
# Format code
go fmt ./...

# Run linter (if available)
golangci-lint run

# Static analysis
go vet ./...
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestPattern ./core

# Run with race detector
go test -race ./...

# Benchmark tests
go test -bench=. ./...
```

### Writing Tests

Tests should be focused and meaningful:

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "simple case",
            input:    "test",
            expected: "result",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.expected {
                t.Errorf("MyFunction() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### Pattern Testing

Add pattern tests to `core/bundle_test.go`:

```go
"my-pattern": {
    matches:    []string{"positive case 1", "positive case 2"},
    nonMatches: []string{"negative case 1", "negative case 2"},
},
```

## Building

### Development Builds

```bash
# Build CLI
go build -o atheon .

# Build MCP server
go build -o atheon-mcp ./cmd/mcp

# Build with race detector
go build -race -o atheon-debug .
```

### Production Builds

Production builds are handled by GoReleaser. For manual testing:

```bash
# Build for current platform
go build -ldflags="-s -w" -o atheon .

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o atheon-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o atheon-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o atheon-windows-amd64.exe .
```

## Pattern Development

### Creating Patterns

1. **Create YAML file** in appropriate category:
   ```bash
   # community/secrets/my-service.yaml
   name: my-service-api-key
   match: '\bmsvc_[A-Za-z0-9]{32}\b'
   ```

2. **Rebuild bundle:**
   ```bash
   go run ./bundler
   ```

3. **Add test cases** to `core/bundle_test.go`:
   ```go
   "my-service-api-key": {
       matches:    []string{"token=msvc_" + strings.Repeat("a", 32)},
       nonMatches: []string{"token=msvc_short", "token=other_" + strings.Repeat("a", 32)},
   },
   ```

4. **Test and verify:**
   ```bash
   go test ./...
   ./atheon list | grep my-service-api-key
   ```

5. **Commit both files:**
   ```bash
   git add community/secrets/my-service.yaml core/patterns.bundle
   git commit -m "Add my-service-api-key pattern"
   ```

## Debugging

### Debug Mode

```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o atheon-debug .

# Use with debugger
dlv debug ./atheon-debug
```

### Verbose Output

```bash
# Run with verbose logging (if available)
./atheon --help

# Enable race detection
go test -race ./...
```

### Common Issues

**Pattern not matching:**
```bash
# Verify pattern is loaded
./atheon list | grep pattern-name

# Test with simple case
echo "test input" | ./atheon -
```

**Build failures:**
```bash
# Clean and rebuild
go clean -cache
go mod tidy
go build -o atheon .
```

## Release Process

Releases are automated via GitHub Actions on version tags (10th and 21st of every month):

1. **Update version** in code if needed
2. **Create git tag:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. **GitHub Actions** automatically:
   - Runs GoReleaser
   - Builds for all platforms
   - Creates GitHub Release
   - Publishes to Homebrew and Scoop

## Continuous Integration

GitHub Actions workflow runs on:
- Pull requests
- Push to main branch
- Tags for releases

**CI checks:**
- Go tests
- Go vet
- Build verification
- Bundle generation

## Contributing Guidelines

### What to Contribute

**Good contributions:**
- New patterns for missing domains
- Bug fixes and performance improvements
- Documentation improvements
- Test coverage expansion

**Not in scope:**
- Core engine changes without clear justification
- Adding dependencies without strong rationale
- Breaking changes to pattern interface

### Pull Request Process

1. **Fork and branch**
2. **Make changes** following code standards
3. **Add tests** for new functionality
4. **Update documentation** as needed
5. **Create PR** with:
   - Clear description of changes
   - Rationale for the contribution
   - Test results
   - Any breaking changes or migration notes

### Code Review Criteria

- **Correctness:** Does it work as intended?
- **Testing:** Are there adequate tests?
- **Documentation:** Is it documented?
- **Style:** Does it follow Go conventions?
- **Impact:** Are there breaking changes?

## Performance Considerations

### Pattern Matching

- Patterns are pre-compiled and cached
- Category filtering reduces unnecessary matching
- Combined regex for efficient pre-filtering

### File Scanning

- Parallel processing with 256 workers
- File size limits prevent memory exhaustion
- Binary file detection and skipping

### Optimization Tips

- Use category filtering for large repositories
- Create `.atheonignore` for large exclusions
- Use `--file` for single files instead of directory scans

## Security Considerations

### Input Validation

- Validate file paths before reading
- Limit file sizes to prevent DoS
- Sanitize regex patterns to prevent ReDoS

### Secret Handling

- Never log actual secret values
- Redact output by default
- No persistent storage of findings

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Pattern Development Guide](../patterns/creating-patterns.md)