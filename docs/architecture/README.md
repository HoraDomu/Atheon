# Architecture Documentation

This document describes the architecture and design decisions behind Atheon's pattern matching engine.

## System Overview

Atheon implements a **deterministic pattern matching engine** with a plugin architecture for community-driven pattern development. The system is designed around a simple principle: any pattern that can return true/false on text input can be used to scan for security, compliance, or any domain-specific issues.

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                           │
│                    (main.go)                                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Core Engine                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Bundle     │  │   Pattern    │  │   Runner     │     │
│  │   Manager    │  │   Registry   │  │   Scanner    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Pattern Bundle                            │
│              (Embedded + Runtime Loadable)                   │
└─────────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. CLI Layer (`main.go`)

**Purpose:** User interface and command orchestration

**Key Responsibilities:**
- Argument parsing and validation
- Output formatting (text/JSON)
- Exit code management
- User command routing

**Design Decisions:**
- **Exit Code Strategy:** Returns 1 for findings, 0 for clean (CI/CD friendly)
- **Output Redaction:** Automatically redacts sensitive content in output
- **Flexible Input:** Supports files, directories, stdin, and environment variables

**Command Flow:**
```
User Input → Argument Parsing → Category Setup → Core Engine Call → Output Formatting → Exit
```

### 2. Core Engine (`core/`)

#### Bundle Manager (`bundle.go`)

**Purpose:** Pattern loading, lifecycle management, and distribution

**Key Features:**
- **Embedded Bundle:** Patterns shipped in binary via `go:embed`
- **Runtime Updates:** Override with `~/.atheon/patterns.bundle`
- **Category Filtering:** Load only required categories for performance
- **HTTP Updates:** Download latest bundle from GitHub releases

**Architecture Pattern:** Singleton with lazy initialization

```go
// Bundle loading sequence
init() → loadBundle() → SetActiveCategories() → activeScanners
```

**Performance Optimization:**
- Pre-compile all regex patterns at bundle load
- Combine patterns per category into single regex for pre-filtering
- Cache compiled patterns for lifetime of process

#### Pattern Registry (`pattern.go`)

**Purpose:** Pattern interface definition and registration

**Interface Design:**
```go
type Pattern interface {
    Name() string
    Matches(line string) bool
}
```

**Design Rationale:**
- **Minimal Interface:** Only methods needed for pattern matching
- **Extensibility:** New pattern types can implement interface without engine changes
- **Go Idioms:** Interface with single implementations

**Registry Pattern:**
- Global registration at init time
- Sorted output for consistent display
- Thread-safe read-only access after initialization

#### Scanner Runner (`runner.go`)

**Purpose:** File scanning orchestration and parallel processing

**Key Features:**
- **Parallel Processing:** 256 concurrent workers for file scanning
- **Smart Filtering:** Automatic skipping of binary files and ignored directories
- **Integration Support:** Respects `.gitignore` and `.atheonignore`
- **Memory Efficient:** Processes files independently without caching

**Scanning Pipeline:**
```
Directory Traversal → File Filtering → Parallel Processing → Pattern Matching → Result Aggregation
```

**Performance Characteristics:**
- **Scalability:** Handles large codebases efficiently
- **Memory:** Linear scaling with file count, not file size
- **Speed:** Parallel processing with worker pool pattern

**Concurrency Model:**
```go
// Semaphore-based worker pool
sem := make(chan struct{}, 256)
for _, path := range paths {
    sem <- struct{}{}
    go func(p string) {
        defer func() { <-sem }()
        // Process file
    }(path)
}
```

#### Data Structures (`finding.go`)

**Purpose:** Result representation and statistics

**Design:**
- **Simple Structs:** Minimal overhead for result storage
- **Clear Semantics:** File, line, pattern, content
- **Statistics Tracking:** Files scanned, bytes processed, timing

### 3. MCP Server (`cmd/mcp/main.go`)

**Purpose:** Model Context Protocol integration for AI tools

**Architecture:**
- **Protocol Handler:** JSON-RPC 2.0 over stdin/stdout
- **Tool Registration:** Dynamic tool list generation
- **Category Support:** Per-request category filtering

**Tool Design:**
```go
"scan_string" → ScanString() → textResult()
"scan_file"   → ScanFile()   → textResult()
"scan_dir"    → ScanDir()    → textResult()
```

**Integration Pattern:**
- Single binary deployment
- No external dependencies
- Stateless request handling

### 4. Pattern Bundler (`bundler/main.go`)

**Purpose:** Compile YAML patterns into embedded bundle

**Build Pipeline:**
```
YAML Files → Pattern Validation → JSON Serialization → GZIP Compression → Bundle Output
```

**Design Decisions:**
- **JSON Format:** Easy parsing and compatibility
- **GZIP Compression:** Reduce bundle size
- **Category Extraction:** Folder structure determines categories
- **Error Handling:** Fail on invalid patterns

**Usage:**
```bash
go run ./bundler [community-dir] [output-path]
```

## Data Flow

### Pattern Loading Flow

```
1. Application Start
   ↓
2. Check for ~/.atheon/patterns.bundle
   ↓
3. Load embedded bundle if not found
   ↓
4. Decompress and parse JSON
   ↓
5. Compile regex patterns
   ↓
6. Register patterns
   ↓
7. Build category scanners
   ↓
8. Ready for scanning
```

### Scanning Flow

```
1. User invokes scan
   ↓
2. Parse categories and arguments
   ↓
3. Set active categories
   ↓
4. Traverse directory/file
   ↓
5. Filter ignored files
   ↓
6. Read file contents
   ↓
7. Split into lines
   ↓
8. For each line and category:
   - Check combined regex (pre-filter)
   - Check individual patterns (match)
   ↓
9. Collect findings
   ↓
10. Format output
   ↓
11. Return exit code
```

## Design Patterns

### 1. Plugin Architecture

**Pattern:** Interface-based extensibility

**Implementation:**
```go
type Pattern interface {
    Name() string
    Matches(line string) bool
}

// Any type implementing this can be registered
func Register(p Pattern) {
    registry = append(registry, p)
}
```

**Benefits:**
- Community contributions without code changes
- Domain-specific pattern implementations
- Clear extension points

### 2. Strategy Pattern

**Pattern:** Category-based scanning strategies

**Implementation:**
```go
type categoryScanner struct {
    combined *regexp.Regexp    // Pre-filter strategy
    patterns []Pattern         // Detailed matching strategy
}
```

**Benefits:**
- Performance optimization through pre-filtering
- Category-specific optimization
- Flexible matching strategies

### 3. Builder Pattern

**Pattern:** Bundle construction and loading

**Implementation:**
```go
// YAML → JSON → GZIP → Binary Bundle
func loadBundle(data []byte) error {
    r, _ := gzip.NewReader(bytes.NewReader(data))
    json.NewDecoder(r).Decode(&defs)
    // Build pattern objects
}
```

**Benefits:**
- Separation of pattern definition from loading
- Multiple serialization formats
- Build-time optimization

### 4. Singleton Pattern

**Pattern:** Global pattern registry

**Implementation:**
```go
var registry []Pattern

func All() []Pattern {
    return registry
}
```

**Benefits:**
- Single source of truth
- Consistent pattern availability
- Memory efficiency

## Performance Considerations

### Memory Management

**Strategies:**
- **No File Caching:** Read and process files immediately
- **Streaming Parsing:** JSON decoder for bundle loading
- **Pattern Pre-compilation:** Compile once, use many times
- **Worker Pool:** Limit concurrent goroutines to 256

### CPU Optimization

**Techniques:**
- **Parallel File Processing:** Utilize multiple cores
- **Category Pre-filtering:** Reduce pattern checks
- **Regex Caching:** Compile patterns at load time
- **Combined Regex:** Single pre-filter per category

### I/O Efficiency

**Approaches:**
- **Binary Detection:** Skip known binary file types
- **Ignore Rules:** Skip unnecessary directories
- **File Size Limits:** Prevent memory exhaustion
- **Batch Processing:** Process multiple files in parallel

## Security Architecture

### Input Validation

**Strategies:**
- **File Path Validation:** Prevent path traversal
- **Regex Validation:** Check pattern compilation
- **Size Limits:** Prevent memory exhaustion
- **Type Checking:** Validate input parameters

### Secret Handling

**Principles:**
- **No Logging:** Never log actual secret values
- **Output Redaction:** Mask sensitive content in display
- **No Persistence:** Don't store findings permanently
- **Memory Clearing:** Clear sensitive data after processing

### Threat Model

**Considerations:**
- **Malicious Patterns:** Validate regex to prevent ReDoS
- **Large Files:** Size limits to prevent DoS
- **Path Traversal:** Validate all file paths
- **Code Injection:** No eval or dynamic code execution

## Extensibility Points

### 1. Custom Pattern Types

```go
type MyPattern struct {
    name string
    config MyConfig
}

func (p *MyPattern) Name() string { return p.name }
func (p *MyPattern) Matches(line string) bool {
    // Custom matching logic
    return myMatcher.Match(line)
}
```

### 2. Category Extensions

```bash
# Create new category folder
mkdir community/mycategory

# Add patterns
# community/mycategory/pattern1.yaml
# community/mycategory/pattern2.yaml
```

### 3. Output Formats

```go
// Add new output format in main.go
func printMyFormat(findings []core.Finding) {
    // Custom formatting logic
}
```

## Technology Choices

### Go Language

**Rationale:**
- **Performance:** Compiled language with efficient execution
- **Concurrency:** Built-in goroutines for parallel processing
- **Cross-platform:** Easy cross-compilation for multiple platforms
- **Standard Library:** Excellent regex and file I/O support
- **Single Binary:** Easy deployment without dependencies

### RE2 Regex Engine

**Rationale:**
- **Performance:** Linear time complexity
- **Safety:** No catastrophic backtracking
- **Predictability:** Consistent performance across patterns
- **Go Integration:** Native Go regex implementation

### YAML for Patterns

**Rationale:**
- **Readability:** Human-friendly format
- **Simplicity:** Easy to write and maintain
- **No Code:** Non-programmers can contribute patterns
- **Structure:** Supports key-value pairs cleanly

### JSON for Bundle

**Rationale:**
- **Parsing:** Fast and efficient JSON parsing in Go
- **Compression:** Compresses well with GZIP
- **Compatibility:** Universal format support
- **Type Safety:** Clear data structure mapping

## Future Architecture Considerations

### Potential Improvements

1. **Streaming Architecture:** Process large files without full loading
2. **Pattern Caching:** Cache compiled patterns across runs
3. **Index Generation:** Pre-index code for faster repeated scans
4. **Plugin System:** Dynamic pattern loading without recompilation
5. **API Layer:** HTTP API for remote scanning services

### Scalability Enhancements

1. **Distributed Scanning:** Split scanning across multiple machines
2. **Incremental Scanning:** Only scan changed files
3. **Pattern Optimization:** Automatically optimize pattern performance
4. **Memory Profiling:** Identify and fix memory bottlenecks

## Documentation Maintenance

**Architecture Updates:**
- Update this document when making significant architectural changes
- Document design decisions and trade-offs
- Include performance characteristics for new features
- Maintain diagrams and flowcharts

**Related Documentation:**
- [Development Guide](../development/README.md)
- [Pattern Development](../patterns/README.md)
- [Integration Guide](../integration/README.md)