# MCP Server Reference

Reference for the Atheon Model Context Protocol (MCP) server implementation.

## Overview

Atheon includes an MCP server implementation that allows AI tools (like Claude Desktop, Cline, etc.) to directly invoke pattern scanning capabilities.

## Installation

### Building the MCP Server

```bash
# Build from source
go build -o atheon-mcp ./cmd/mcp

# Or use package manager
brew install atheon  # Includes atheon-mcp
```

### Verifying Installation

```bash
atheon-mcp --version
```

## Configuration

### Claude Desktop

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

### Cline (VS Code)

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

### With Category Filtering

```json
{
  "mcpServers": {
    "atheon": {
      "command": "/usr/local/bin/atheim-mcp",
      "env": {
        "ATHEON_CATEGORIES": "secrets,pii"
      }
    }
  }
}
```

## Available Tools

### `scan_string`

Scan text content for pattern matches.

**Parameters:**
```json
{
  "content": "string to scan",
  "source": "optional-source-name",
  "categories": ["secrets", "pii"]
}
```

**Example Usage:**
```json
{
  "name": "scan_string",
  "arguments": {
    "content": "api_key=sk-1234567890abcdefghijklmnopqrstuvwxyz",
    "source": "test-input",
    "categories": ["secrets"]
  }
}
```

**Response:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "openai-api-key  test-input:1\n\n1 finding(s)"
    }
  ]
}
```

---

### `scan_file`

Scan a specific file for pattern matches.

**Parameters:**
```json
{
  "path": "/path/to/file",
  "categories": ["secrets", "pii"]
}
```

**Example Usage:**
```json
{
  "name": "scan_file",
  "arguments": {
    "path": "/path/to/config.yaml",
    "categories": ["secrets"]
  }
}
```

**Response:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "github-pat  /path/to/config.yaml:47\n\n1 finding(s)"
    }
  ]
}
```

---

### `scan_dir`

Scan a directory for pattern matches.

**Parameters:**
```json
{
  "path": "/path/to/directory",
  "categories": ["secrets", "pii"]
}
```

**Example Usage:**
```json
{
  "name": "scan_dir",
  "arguments": {
    "path": "/path/to/project",
    "categories": ["secrets", "pii"]
  }
}
```

**Response:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "openai-api-key  /path/to/project/config.yaml:47\ngithub-pat  /path/to/project/.env:3\n\n2 finding(s)\nscanned 14 file(s)  22.1 KB  3ms"
    }
  ]
}
```

---

## Protocol Details

### Server Information

**Protocol Version:** 2024-11-05
**Server Name:** atheon
**Server Version:** 1.0.0

### Initialize Response

```json
{
  "protocolVersion": "2024-11-05",
  "capabilities": {
    "tools": {}
  },
  "serverInfo": {
    "name": "atheon",
    "version": "1.0.0"
  }
}
```

### Tool List Response

```json
{
  "tools": [
    {
      "name": "scan_string",
      "description": "Scan a string for pattern matches",
      "inputSchema": {
        "type": "object",
        "properties": {
          "content": {
            "type": "string"
          },
          "source": {
            "type": "string"
          },
          "categories": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "categories to scan (omit for all)"
          }
        },
        "required": ["content"]
      }
    },
    {
      "name": "scan_file",
      "description": "Scan a file for pattern matches",
      "inputSchema": {
        "type": "object",
        "properties": {
          "path": {
            "type": "string"
          },
          "categories": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "categories to scan (omit for all)"
          }
        },
        "required": ["path"]
      }
    },
    {
      "name": "scan_dir",
      "description": "Scan a directory for pattern matches",
      "inputSchema": {
        "type": "object",
        "properties": {
          "path": {
            "type": "string"
          },
          "categories": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "categories to scan (omit for all)"
          }
        },
        "required": ["path"]
      }
    }
  ]
}
```

## Usage Examples

### Claude Desktop Integration

**Scenario:** AI assistant scans code during generation

```json
// claude_desktop_config.json
{
  "mcpServers": {
    "atheon": {
      "command": "/usr/local/bin/atheim-mcp",
      "env": {
        "ATHEON_CATEGORIES": "secrets,pii"
      }
    }
  }
}
```

**AI Conversation:**
```
User: Create a Node.js API with environment variable configuration

AI: [Generates code and internally calls atheon-mcp]

    MCP Call: scan_string with generated code
    Response: No findings detected

    [Returns code to user]
```

### Code Review Workflow

**Scenario:** Real-time security scanning during development

```json
{
  "name": "scan_dir",
  "arguments": {
    "path": "/workspace/current-project",
    "categories": ["secrets"]
  }
}
```

**AI Response:**
```
I've scanned your current project for security issues. Here's what I found:

⚠️ Security Findings:
- github-pat detected in .env:3
- openai-api-key detected in config.yaml:47

Recommendation:
1. Remove these secrets from your code
2. Use environment variables for secret management
3. Rotate exposed credentials immediately
```

### Pre-commit Integration

**Scenario:** AI assistant enforces security before commits

```json
{
  "name": "scan_string",
  "arguments": {
    "content": "[git diff content]",
    "source": "git-diff",
    "categories": ["secrets", "pii"]
  }
}
```

**AI Response:**
```
I've analyzed your git diff for potential security issues:

✅ No security findings detected. Your changes are safe to commit.
```

## Error Handling

### Invalid Parameters

```json
// Request
{
  "name": "scan_file",
  "arguments": {}
}

// Response
{
  "error": {
    "code": -32602,
    "message": "invalid params"
  }
}
```

### File Not Found

```json
// Request
{
  "name": "scan_file",
  "arguments": {
    "path": "/nonexistent/file.txt"
  }
}

// Response
{
  "error": {
    "code": -32603,
    "message": "file not found: /nonexistent/file.txt"
  }
}
```

### Unknown Tool

```json
// Request
{
  "name": "unknown_tool",
  "arguments": {}
}

// Response
{
  "error": {
    "code": -32601,
    "message": "unknown tool: unknown_tool"
  }
}
```

## Advanced Configuration

### Custom Category Sets

```json
{
  "mcpServers": {
    "atheon-security": {
      "command": "/usr/local/bin/atheim-mcp",
      "env": {
        "ATHEON_CATEGORIES": "secrets"
      }
    },
    "atheon-pii": {
      "command": "/usr/local/bin/atheim-mcp",
      "env": {
        "ATHEON_CATEGORIES": "pii"
      }
    },
    "atheon-full": {
      "command": "/usr/local/bin/atheim-mcp",
      "env": {
        "ATHEON_ALL": "1"
      }
    }
  }
}
```

### Working Directory Configuration

```json
{
  "mcpServers": {
    "atheon": {
      "command": "/usr/local/bin/atheim-mcp",
      "cwd": "/path/to/project",
      "env": {
        "ATHEON_CATEGORIES": "secrets"
      }
    }
  }
}
```

## Performance Considerations

### Optimization Tips

1. **Use Category Filtering:** Only scan relevant categories
2. **Target Specific Paths:** Scan specific files rather than entire projects
3. **Cache Results:** Avoid repeated scans of unchanged content
4. **Use Scan String:** For small content, prefer `scan_string` over file operations

### Typical Performance

- **scan_string:** <10ms for typical code snippets
- **scan_file:** 50-200ms depending on file size
- **scan_dir:** 100-500ms depending on project size

## Integration Patterns

### Pre-generation Scanning

AI tools can scan generated code before presenting to users:

```
1. User requests code generation
2. AI generates code
3. AI calls scan_string
4. If findings: AI removes/fixes issues
5. AI presents clean code to user
```

### Real-time Feedback

Provide immediate security feedback during development:

```
1. User makes changes to files
2. AI tool calls scan_dir
3. AI presents findings immediately
4. User fixes issues before commit
```

### Pipeline Integration

Integrate into development workflows:

```
1. Pre-commit: AI scans staged files
2. Pre-push: AI scans entire repository
3. PR Review: AI scans changed files
4. CI/CD: AI scans before deployment
```

## Troubleshooting

### Server Not Starting

```bash
# Check if binary exists
which atheon-mcp

# Test directly
echo '{}' | atheon-mcp

# Check permissions
ls -l $(which atheon-mcp)
```

### Connection Issues

```json
// Verify MCP server configuration
{
  "mcpServers": {
    "atheon": {
      "command": "/full/path/to/atheim-mcp"
    }
  }
}
```

### Pattern Not Found

```bash
# Verify patterns are loaded
atheon list

# Update patterns
atheon update
```

## Additional Resources

- [Getting Started Guide](../getting-started/README.md)
- [Integration Guide](../integration/README.md)
- [CLI Reference](cli-reference.md)
- [Pattern Development](../patterns/README.md)